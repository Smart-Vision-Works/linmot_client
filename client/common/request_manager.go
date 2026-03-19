package client_common

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"
	protocol_motion_control "github.com/Smart-Vision-Works/linmot_client/protocol/motion_control"
	protocol_rtc "github.com/Smart-Vision-Works/linmot_client/protocol/rtc"
	transport "github.com/Smart-Vision-Works/linmot_client/transport"
)

const defaultMaxRetries = 6

// RequestManagerConfig holds configuration options for the RequestManager
type RequestManagerConfig struct {
	MaxRetries     int           // Maximum number of retry attempts (-1 => use default 6, 0 => unlimited, >0 => fixed limit)
	DefaultTimeout time.Duration // Default timeout for requests that don't specify one
}

const maxInvalidStatusRejects = 3

// StatusPacketSnapshot captures the most recent invalid status packet for debugging.
type StatusPacketSnapshot struct {
	PacketHex  string
	ReqBits    uint32
	RepBits    uint32
	Len        int
	StatusWord uint16
	StateVar   uint16
	ErrorCode  uint16
	WarnWord   uint16
	Plausible  bool
	Reason     string
	FromAddr   string
	SentHex    string
	SentAt     time.Time
	RxAt       time.Time
	AgeMs      int64
}

// RequestManager manages reliable send/receive operations with the LinMot drive
type RequestManager struct {
	transport                 transport.Client
	timerCycle                time.Duration
	config                    RequestManagerConfig
	rtcCount                  *RtcCounter
	currentRequest            *pendingRequest
	txChannel                 chan *pendingRequest
	pendingRequests           map[uint8]*pendingRequest
	pendingMCRequests         map[uint8]*pendingRequest // Separate map for MC requests (1-4 counter range)
	pendingStatusRequest      *pendingRequest
	lastInvalidStatusSnapshot *StatusPacketSnapshot
	txTicker                  *time.Ticker
	ctx                       context.Context
	cancel                    context.CancelFunc
	wg                        sync.WaitGroup
	mu                        sync.RWMutex
	lastDriveCmdCount         uint8         // Last cmdCount seen from drive (tracked for debugging/testing only; not used in resync logic)
	lastObservedDriveCmdCount uint8         // Last cmdCount observed from any packet (used for seeding/resync beacons)
	rtcSeeded                 bool          // Track if RTC counter was seeded from drive
	resyncAttempted           bool          // Track if we've attempted resync for current request
	requestGate               chan struct{} // Single-flight gate: capacity 1, seeded with one token
	// TEST-ONLY: beforeFirstSend hook for deterministic testing.
	// Never set in production code. Called in startRequest() before first send (if set).
	beforeFirstSend func(*pendingRequest)

	// lastStateVarLowNibble: Last seen StateVarLow low nibble (0-15, but we only use 0-4 range).
	// Used for MC counter calculation matching linudp.cs: CountNibble = (StateVarLow & 0xF) + 1; if >4 wrap to 1.
	lastStateVarLowNibble uint8
	// lastStateVarLow: Full low byte for debugging (optional, never read).
	lastStateVarLow uint8
	// debug: Enable debug logging (controlled by -linmot_debug flag)
	// Uses atomic operations for race-safe reads in hot paths (no mutex dependency)
	debug atomic.Bool
	// udpBindOnce: ensure UDP bind/config logs only once per client.
	udpBindOnce sync.Once
	// rtcResyncCount: Total number of RTC resyncs that have occurred (pure instrumentation, not gated by debug)
	// Uses atomic operations for race-safe increments in hot paths (no mutex dependency)
	rtcResyncCount atomic.Uint64
	// routerTrace: ring buffer of recent router traces (debug only).
	routerTraceMu     sync.Mutex
	routerTrace       []routerTraceEntry
	routerTraceIndex  int
	routerTraceFilled bool
}

// StatusObservation captures a parsed status frame with metadata.
type StatusObservation struct {
	Timestamp              time.Time
	Status                 protocol_common.Status
	DriveCmdCount          uint8
	ClientExpectedCmdCount uint8
}

type routerTraceEntry struct {
	timestamp time.Time
	entry     string
}

const routerTraceBufferSize = 100

// NewRequestManager creates a new RequestManager with default configuration
func NewRequestManager(transport transport.Client, timerCycle time.Duration) *RequestManager {
	return NewRequestManagerWithConfig(transport, timerCycle, RequestManagerConfig{
		MaxRetries:     -1,
		DefaultTimeout: 20 * time.Millisecond,
	})
}

// NewRequestManagerWithConfig creates a new RequestManager with custom configuration
func NewRequestManagerWithConfig(transport transport.Client, timerCycle time.Duration, config RequestManagerConfig) *RequestManager {
	// Set defaults for zero values
	if config.MaxRetries == -1 {
		config.MaxRetries = defaultMaxRetries
	}
	if config.DefaultTimeout == 0 {
		config.DefaultTimeout = 20 * time.Millisecond
	}

	rm := &RequestManager{
		transport:             transport,
		timerCycle:            timerCycle,
		config:                config,
		rtcCount:              NewRTCCounter(),
		txChannel:             make(chan *pendingRequest, 10),
		pendingRequests:       make(map[uint8]*pendingRequest),
		pendingMCRequests:     make(map[uint8]*pendingRequest),
		requestGate:           make(chan struct{}, 1),
		lastStateVarLowNibble: 0, // Indicates no StateVar seen yet
		lastStateVarLow:       0,
		routerTrace:           make([]routerTraceEntry, routerTraceBufferSize),
		// debug defaults to false (atomic.Bool zero value is false)
	}
	// Seed the gate with one token
	rm.requestGate <- struct{}{}
	return rm
}

// SetDebug enables or disables debug logging for RequestManager operations.
// Uses atomic operations for race-safe updates (no mutex dependency).
func (requestManager *RequestManager) SetDebug(enabled bool) {
	requestManager.debug.Store(enabled)
	if debugSetter, ok := requestManager.transport.(transport.DebugSetter); ok {
		debugSetter.SetDebug(enabled)
	}
	if enabled {
		requestManager.logUDPBindOnce()
	}
}

// DebugEnabled reports whether debug logging is enabled.
func (requestManager *RequestManager) DebugEnabled() bool {
	return requestManager.debug.Load()
}

// TransportClient returns the underlying transport client used for I/O operations.
// Returns nil if the transport is not available.
func (requestManager *RequestManager) TransportClient() transport.Client {
	return requestManager.transport
}

// LastInvalidStatusSnapshot returns the last invalid status packet snapshot, if any.
func (requestManager *RequestManager) LastInvalidStatusSnapshot() (StatusPacketSnapshot, bool) {
	requestManager.mu.RLock()
	defer requestManager.mu.RUnlock()
	if requestManager.lastInvalidStatusSnapshot == nil {
		return StatusPacketSnapshot{}, false
	}
	return *requestManager.lastInvalidStatusSnapshot, true
}

func (requestManager *RequestManager) logUDPBindOnce() {
	if !requestManager.debug.Load() {
		return
	}
	requestManager.udpBindOnce.Do(func() {
		info, ok := requestManager.udpInfo()
		if !ok {
			fmt.Printf("[UDP_BIND] transport does not expose UDP info\n")
			return
		}
		boundPort := "unknown"
		localIP := "unknown"
		if info.LocalAddr != "" {
			if host, port, err := net.SplitHostPort(info.LocalAddr); err == nil {
				boundPort = port
				localIP = host
			}
		}
		fmt.Printf("[UDP_BIND] local_ip=%s local=%s remote=%s (drive)\n", localIP, info.LocalAddr, info.RemoteAddr)
		fmt.Printf("[UDP_CONFIG] linmot_master_port=%d bound_port=%s\n", info.MasterPort, boundPort)
	})
}

func (requestManager *RequestManager) udpInfo() (transport.UDPInfo, bool) {
	provider, ok := requestManager.transport.(transport.UDPInfoProvider)
	if !ok {
		return transport.UDPInfo{}, false
	}
	return provider.UDPInfo(), true
}

func (requestManager *RequestManager) recordRouterTrace(entry string) {
	if !requestManager.debug.Load() {
		return
	}
	requestManager.routerTraceMu.Lock()
	requestManager.routerTrace[requestManager.routerTraceIndex] = routerTraceEntry{
		timestamp: time.Now(),
		entry:     entry,
	}
	requestManager.routerTraceIndex = (requestManager.routerTraceIndex + 1) % len(requestManager.routerTrace)
	if requestManager.routerTraceIndex == 0 {
		requestManager.routerTraceFilled = true
	}
	requestManager.routerTraceMu.Unlock()
}

func (requestManager *RequestManager) snapshotRouterTrace() []routerTraceEntry {
	requestManager.routerTraceMu.Lock()
	defer requestManager.routerTraceMu.Unlock()

	var entries []routerTraceEntry
	if !requestManager.routerTraceFilled {
		entries = append(entries, requestManager.routerTrace[:requestManager.routerTraceIndex]...)
		return entries
	}

	entries = append(entries, requestManager.routerTrace[requestManager.routerTraceIndex:]...)
	entries = append(entries, requestManager.routerTrace[:requestManager.routerTraceIndex]...)
	return entries
}

func (requestManager *RequestManager) dumpRouterTrace(reason string) {
	if !requestManager.debug.Load() {
		return
	}
	entries := requestManager.snapshotRouterTrace()
	if len(entries) == 0 {
		return
	}
	fmt.Printf("[ROUTER_TRACE] %s (last %d entries)\n", reason, len(entries))
	for _, entry := range entries {
		fmt.Printf("[ROUTER_TRACE] %s %s\n", entry.timestamp.Format(time.RFC3339Nano), entry.entry)
	}
}

// rtcRequestWithCmdCode is an interface for RTC requests that have a command code.
// This works for both base types (RTCGetParamRequest, RTCSetParamRequest) and
// wrapper types that embed them (e.g., RestartDriveRequest, StartGettingUPIDListRequest).
type rtcRequestWithCmdCode interface {
	CmdCode() uint8
}

// rtcRequestWithUPID is an interface for RTC requests that have a UPID.
// This works for parameter access requests (RTCGetParamRequest, RTCSetParamRequest)
// and their wrappers.
type rtcRequestWithUPID interface {
	UPID() uint16
}

// extractRTCCommandCode extracts the command code from an RTC request.
// Returns 0 if the request is not an RTC request or doesn't have a command code.
func extractRTCCommandCode(req protocol_common.Request) uint8 {
	if rtcReq, ok := req.(rtcRequestWithCmdCode); ok {
		return rtcReq.CmdCode()
	}
	return 0
}

// extractRTCUPID extracts the UPID from an RTC request.
// Returns 0 if the request is not an RTC request or doesn't have a UPID.
func extractRTCUPID(req protocol_common.Request) uint16 {
	if rtcReq, ok := req.(rtcRequestWithUPID); ok {
		return rtcReq.UPID()
	}
	return 0
}

// updateLastStateVarLow extracts and stores the low nibble from StateVar.
// Called whenever we parse a packet containing StateVar (Status or MC responses).
// This matches linudp.cs behavior: MC counter is derived from last-seen StateVarLow.
// Thread-safe: protected by rm.mu.
func (rm *RequestManager) updateLastStateVarLow(stateVar uint16) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	lowByte := uint8(stateVar & 0xFF)
	rm.lastStateVarLow = lowByte // Debug-only, never read
	rm.lastStateVarLowNibble = lowByte & 0x0F
}

// nextMCCountFromLastStateVar calculates the next MC counter based on last-seen StateVarLow.
// Matches linudp.cs exactly: CountNibble = (StateVarLow & 0xF) + 1; if >4 wrap to 1.
// Returns 1 if no StateVar has been observed yet (lastStateVarLowNibble == 0).
// Thread-safe: protected by rm.mu (read lock).
func (rm *RequestManager) nextMCCountFromLastStateVar() uint8 {
	rm.mu.RLock()
	nibble := rm.lastStateVarLowNibble & 0x0F
	rm.mu.RUnlock()

	// If no StateVar seen yet, default to 1 (matches linudp.cs)
	if nibble == 0 {
		return 1
	}

	// CountNibble = (StateVarLow & 0xF) + 1; if >4 wrap to 1
	next := nibble + 1
	if next > 4 {
		next = 1
	}
	return next
}

// Start starts the background send and receive goroutines.
func (requestManager *RequestManager) Start() {
	requestManager.mu.Lock()
	if requestManager.cancel != nil {
		requestManager.mu.Unlock()
		return
	}
	requestManager.ctx, requestManager.cancel = context.WithCancel(context.Background())
	requestManager.txTicker = time.NewTicker(requestManager.timerCycle)
	requestManager.mu.Unlock()

	requestManager.wg.Add(2)
	go requestManager.txLoop()
	go requestManager.rxLoop()
}

// Stop stops the background goroutines, waits for them to finish, and closes the transport.
// To prevent deadlock, the transport is closed BEFORE waiting for goroutines to finish.
// This ensures that any blocking RecvPacket() calls are immediately unblocked.
func (requestManager *RequestManager) Stop() error {
	var errClose error
	if requestManager.cancel != nil {
		requestManager.cancel()
	}
	if requestManager.txTicker != nil {
		requestManager.txTicker.Stop()
	}
	// Close transport BEFORE wg.Wait() to unblock any RecvPacket() calls.
	// This prevents deadlock: rxLoop may be blocked in RecvPacket(), and closing
	// the transport causes the underlying I/O to return immediately.
	if requestManager.transport != nil {
		errClose = requestManager.transport.Close()
	}
	requestManager.wg.Wait()
	requestManager.mu.Lock()
	requestManager.cancel = nil
	requestManager.ctx = nil
	requestManager.txTicker = nil
	requestManager.mu.Unlock()
	return errClose
}

// SendRequestAndReceive performs a request-response cycle using the background channels.
func SendRequestAndReceive[T protocol_common.Response](requestManager *RequestManager, ctx context.Context, request protocol_common.Request) (T, error) {
	var zero T

	response, err := requestManager.roundTripBase(ctx, request)
	if err != nil {
		return zero, err
	}

	result, ok := response.(T)
	if !ok {
		switch typed := response.(type) {
		case *protocol_common.MonitoringStatusResponse:
			var target T
			if _, ok := any(target).(*protocol_common.StatusResponse); ok {
				converted := protocol_common.NewStatusResponse(typed.Status())
				return any(converted).(T), nil
			}
		}
		return zero, errors.Errorf("unexpected response type, expected %T but got %T", zero, response)
	}
	return result, nil
}

func (requestManager *RequestManager) SendRequest(ctx context.Context, request protocol_common.Request) error {
	_, err := requestManager.roundTripBase(ctx, request)
	return err
}

// removePending removes a pending request from the appropriate map (status, MC, or RTC).
// This is the safe way to clean up any request type, handling cmdCount=0 for RTC requests that
// were canceled before startRequest() assigned a cmdCount.
func (requestManager *RequestManager) removePending(request *pendingRequest) {
	requestManager.mu.Lock()
	defer requestManager.mu.Unlock()
	requestManager.removePendingLocked(request)
}

// roundTripBase performs the non-generic parts of a request-response cycle.
func (requestManager *RequestManager) roundTripBase(ctx context.Context, request protocol_common.Request) (protocol_common.Response, error) {
	// Calculate timeout durations before gate acquisition
	perAttemptTimeout := request.OperationTimeout()
	if perAttemptTimeout <= 0 {
		perAttemptTimeout = requestManager.config.DefaultTimeout
	}
	// Total timeout is per-attempt timeout * max retries
	totalTimeout := perAttemptTimeout * time.Duration(requestManager.config.MaxRetries)

	// Acquire gate BEFORE creating/submitting request (single-flight gating)
	// This ensures only one request is in-flight at a time
	select {
	case <-requestManager.requestGate:
		// Gate acquired - defer release
		defer func() {
			requestManager.requestGate <- struct{}{}
		}()
	case <-ctx.Done():
		// Caller canceled while waiting for gate
		return nil, ctx.Err()
	case <-requestManager.ctx.Done():
		// Manager shutdown while waiting for gate
		return nil, requestManager.ctx.Err()
	}

	// Gate acquired - now submit the request
	pendingRequest, err := requestManager.submitRequest(ctx, request, totalTimeout, perAttemptTimeout)
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		// Use removePending() which handles all request types (status, MC, RTC)
		// and safely handles cmdCount=0 for RTC requests canceled before startRequest()
		requestManager.removePending(pendingRequest)
	}

	// CRITICAL: Wait for sentCh to be closed before starting timeout timer.
	// This ensures startRequest() has allocated counter, stamped deadline, and attempted send.
	// Deadline is stamped BEFORE send (matching C#: timeout starts after counter increment, before send).
	var timeoutDuration time.Duration
	select {
	case <-pendingRequest.sentCh:
		// Request processed: counter allocated, deadline stamped, send attempted - compute timeout duration ONCE
		timeoutDuration = time.Until(pendingRequest.deadline)
		if timeoutDuration < 0 {
			timeoutDuration = 0
		}
	case <-ctx.Done():
		cleanup()
		return nil, ctx.Err()
	case <-requestManager.ctx.Done():
		cleanup()
		return nil, requestManager.ctx.Err()
	case err := <-pendingRequest.errorChannel:
		// Error occurred before first send (e.g., context canceled in startRequest)
		cleanup()
		return nil, err
	}

	// Now wait for response with timeout based on deadline (which is now set)
	timer := time.NewTimer(timeoutDuration)
	defer timer.Stop()
	select {
	case response := <-pendingRequest.responseChannel:
		return response, nil

	case err := <-pendingRequest.errorChannel:
		return nil, err

	case <-ctx.Done():
		cleanup()
		return nil, ctx.Err()

	case <-timer.C:
		requestManager.mu.Lock()
		attempts := int(pendingRequest.sendCounter)
		requestManager.mu.Unlock()
		cleanup()
		return nil, requestManager.newTimeoutError(attempts, timeoutDuration)
	}
}

// submitRequest submits a request to the background channel.
// totalTimeout and perAttemptTimeout are calculated in roundTripBase() before gate acquisition.
func (requestManager *RequestManager) submitRequest(ctx context.Context, request protocol_common.Request, totalTimeout time.Duration, perAttemptTimeout time.Duration) (*pendingRequest, error) {
	_, isStatusRequest := request.(*protocol_common.StatusRequest)
	_, isConnectivityProbeRequest := request.(*protocol_common.ConnectivityProbeRequest)
	_, isMonitoringStatusRequest := request.(*protocol_common.MonitoringStatusRequest)
	// Control Word requests also don't use RTC counters
	_, isControlWordRequest := request.(interface{ GetControlWord() uint16 })
	// Motion Control requests use their own counter system
	_, isMCRequest := request.(protocol_motion_control.MCRequest)

	// ConnectivityProbeRequest is treated as a status request (no RTC counter)
	isAnyStatusRequest := isStatusRequest || isConnectivityProbeRequest

	// Only RTC requests (not Status, ConnectivityProbe, MonitoringStatus, Control Word, or MC) get RTC counters
	needsRTCCounter := !isAnyStatusRequest && !isMonitoringStatusRequest && !isControlWordRequest && !isMCRequest

	var cmdCount uint8
	var mcCount uint8
	if needsRTCCounter {
		// cmdCount will be allocated in startRequest() (tx goroutine) right before first send
		// This matches C# linudp.cs behavior: increment before sending, not when queuing
		cmdCount = 0 // Indicates "assign in startRequest()"
	} else if isMCRequest {
		// mcCount will be allocated in startRequest() (tx goroutine) right before first send
		// This matches C# linudp.cs behavior: counter selected at send-time, not queue-time
		mcCount = 0 // Indicates "assign in startRequest()"
	}

	responseChannel := make(chan protocol_common.Response, 1)
	errorChannel := make(chan error, 1)

	// Extract UPID from request (0 if not a parameter request)
	originalUPID := extractRTCUPID(request)

	pendingRequest := newPendingRequest(
		request,
		cmdCount,
		mcCount,
		originalUPID,
		responseChannel,
		errorChannel,
		ctx, // Store caller's context for cancel-before-send check
		totalTimeout,
		perAttemptTimeout,
		isAnyStatusRequest || isMonitoringStatusRequest || isControlWordRequest, // Status, ConnectivityProbe, MonitoringStatus, and Control Word are non-RTC/non-MC
		isMCRequest,
	)

	select {
	case requestManager.txChannel <- pendingRequest:
		return pendingRequest, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-requestManager.ctx.Done():
		return nil, requestManager.ctx.Err()
	}
}

// txLoop is the background send goroutine.
func (requestManager *RequestManager) txLoop() {
	defer requestManager.wg.Done()

	for {
		select {
		case <-requestManager.ctx.Done():
			return

		case request := <-requestManager.txChannel:
			requestManager.startRequest(request)

		case <-requestManager.txTicker.C:
			requestManager.handleTick()
		}
	}
}

func (requestManager *RequestManager) startRequest(request *pendingRequest) {
	// CRITICAL: Check if caller's context was canceled BEFORE allocating cmdCount,
	// BEFORE adding to pending maps, BEFORE sending packet.
	// This prevents queued-but-not-yet-sent requests from being transmitted after cancellation.
	select {
	case <-request.requestCtx.Done():
		// Caller canceled - complete with error without consuming rtcCount or adding to maps
		requestManager.completeRequestWithError(request, request.requestCtx.Err())
		return
	default:
		// Context not canceled, proceed
	}

	// Test hook: pause before first send (for deterministic testing)
	if requestManager.beforeFirstSend != nil {
		requestManager.beforeFirstSend(request)
		// Re-check context after hook (hook may have canceled it)
		select {
		case <-request.requestCtx.Done():
			requestManager.completeRequestWithError(request, request.requestCtx.Err())
			return
		default:
		}
	}

	// Allocate cmdCount for RTC requests if not yet assigned (matches C# linudp.cs behavior)
	// This happens right before first send, preventing cmdCount from being "burned" by queued requests
	if !request.isStatusRequest && !request.isMCRequest && request.commandCount == 0 {
		requestManager.mu.Lock()
		if !requestManager.rtcSeeded && isValidRTCCounter(requestManager.lastObservedDriveCmdCount) {
			seed := nextRTCCounter(requestManager.lastObservedDriveCmdCount)
			requestManager.rtcCount.Set(seed)
			requestManager.rtcSeeded = true
			if requestManager.debug.Load() {
				fmt.Printf("[RTC_SEED] drive=%d -> client starts at next(drive)=%d\n", requestManager.lastObservedDriveCmdCount, seed)
			}
		}
		requestManager.mu.Unlock()
		request.commandCount = requestManager.rtcCount.Next()
	}

	// Allocate mcCounter for MC requests if not yet assigned (matches C# linudp.cs behavior)
	// This happens right before first send, ensuring counter uses latest StateVarLow value
	if request.isMCRequest && request.mcCounter == 0 {
		mcCount := requestManager.nextMCCountFromLastStateVar()
		request.mcCounter = mcCount
		// Set the counter in the underlying MC request
		if mcReq, ok := request.request.(protocol_motion_control.MCRequest); ok {
			mcReq.SetCounter(mcCount)
		}
	}

	// CRITICAL: Stamp deadline AFTER counter allocation but BEFORE send (matches C# linudp.cs).
	// C# behavior (decompiled_linudp_csharp_lib.cs):
	//   - Line 6520: Counter incremented inside lock
	//   - Line 6533: Lock released
	//   - Line 6537: timeOutTime = getTimeOutTime() called (timeout starts)
	//   - Send happens later in timer loop (async)
	// Our Go equivalent: stamp deadline after counter allocation, before send.
	now := time.Now()
	requestManager.mu.Lock()
	request.lastSendTime = now
	request.sendCounter = 1
	request.deadline = now.Add(request.totalTimeout)
	requestManager.currentRequest = request
	requestManager.resyncAttempted = false // Reset resync flag for new request
	requestManager.addPendingLocked(request)
	requestManager.mu.Unlock()

	// CRITICAL: Signal sentCh immediately after deadline is stamped and request is "live".
	// This unblocks roundTripBase to start its timeout timer, even if SendPacket is slow/blocked.
	// Timeout budget starts when deadline is stamped, NOT when SendPacket completes.
	request.signalFirstSendComplete()

	// Send packet (may fail, but counter and deadline are already set, matching C# behavior)
	if err := requestManager.sendRequestPacket(request, request.sendCounter); err != nil {
		// Send failed - complete with error
		// Note: Counter is NOT rolled back (matches C#: counter consumed even on send failure)
		// Note: sentCh was already closed above, so roundTripBase timeout timer is already running
		requestManager.mu.Lock()
		requestManager.removePendingLocked(request)
		if requestManager.currentRequest == request {
			requestManager.currentRequest = nil
		}
		requestManager.mu.Unlock()
		// Complete with error (idempotent)
		request.tryCompleteWithError(err)
		return
	}

	// Send succeeded - request is already "live" (deadline stamped, sentCh closed above)
}

func (requestManager *RequestManager) sendRequestPacket(request *pendingRequest, attempt uint) error {
	var packet []byte
	var err error

	switch e := request.request.(type) {
	case protocol_rtc.RTCPacketWritable:
		packet, err = e.WriteRtcPacket(request.commandCount)
	case protocol_common.PacketWritable:
		packet, err = e.WritePacket()
	default:
		return errors.Errorf("request type %T is not writeable", e)
	}

	if err != nil {
		return errors.WithMessage(err, "failed to write request packet")
	}

	// Command-table save-to-flash (0x80) is allowed.
	// MC must be stopped before calling this (enforced by caller).
	cmdCode := extractRTCCommandCode(request.request)

	// Debug logging for select request types when debug enabled
	if requestManager.debug.Load() {
		reqBits, repBits, headerErr := protocol_common.ReadPacketHeader(packet)
		requestKind := request.kind()
		requestType := fmt.Sprintf("%T", request.request)

		requestManager.logUDPBindOnce()
		localAddr := ""
		remoteAddr := ""
		if info, ok := requestManager.udpInfo(); ok {
			localAddr = info.LocalAddr
			remoteAddr = info.RemoteAddr
		}

		txCmdCount := "NA"
		txCmdCountValue := uint8(0)
		if _, ok := request.request.(protocol_rtc.RTCPacketWritable); ok {
			if cmdCount, err := protocol_rtc.ExtractRTCCommandCount(packet); err == nil {
				txCmdCountValue = cmdCount
				txCmdCount = fmt.Sprintf("%d(0x%02X)", cmdCount, cmdCount)
			}
		}

		cmdCodeStr := "NA"
		if cmdCode != 0 {
			cmdCodeStr = fmt.Sprintf("0x%02X", cmdCode)
		}

		if headerErr == nil {
			fmt.Printf("[UDP_TX] local=%s -> remote=%s len=%d kind=%s reqBits=0x%08X repBits=0x%08X cmdCount=%s cmdCode=%s\n",
				localAddr, remoteAddr, len(packet), requestKind, reqBits, repBits, txCmdCount, cmdCodeStr)
		} else {
			fmt.Printf("[UDP_TX] local=%s -> remote=%s len=%d kind=%s headerErr=%v cmdCount=%s cmdCode=%s\n",
				localAddr, remoteAddr, len(packet), requestKind, headerErr, txCmdCount, cmdCodeStr)
		}

		if requestKind == "rtc" && headerErr == nil {
			fmt.Printf("[RTC_TX] ts=%s cmd=%s cc=%d attempt=%d reqBits=0x%08X repBits=0x%08X\n",
				time.Now().Format(time.RFC3339Nano), cmdCodeStr, txCmdCountValue, attempt, reqBits, repBits)
		} else if requestKind == "rtc" {
			fmt.Printf("[RTC_TX] ts=%s cmd=%s cc=%d attempt=%d headerErr=%v\n",
				time.Now().Format(time.RFC3339Nano), cmdCodeStr, txCmdCountValue, attempt, headerErr)
		}

		if requestKind == "status" {
			if headerErr == nil {
				fmt.Printf("[ROUTER_TX] kind=%s type=%s reqBits=0x%08X repBits=0x%08X rtcCmdCount=%s\n",
					requestKind, requestType, reqBits, repBits, txCmdCount)
			} else {
				fmt.Printf("[ROUTER_TX] kind=%s type=%s headerErr=%v rtcCmdCount=%s\n",
					requestKind, requestType, headerErr, txCmdCount)
			}
		}

		if requestKind == "status" {
			statusHex := hex.EncodeToString(packet)
			request.statusTxHex = statusHex
			if headerErr == nil {
				fmt.Printf("[STATUS_TX] packet=%s reqBits=0x%08X repBits=0x%08X local=%s remote=%s\n",
					statusHex, reqBits, repBits, localAddr, remoteAddr)
			} else {
				fmt.Printf("[STATUS_TX] packet=%s headerErr=%v local=%s remote=%s\n",
					statusHex, headerErr, localAddr, remoteAddr)
			}
		}

		if cmdCode == 0x35 || cmdCode == 0x36 {
			hexLen := len(packet)
			if hexLen > 64 {
				hexLen = 64
			}
			txHex := hex.EncodeToString(packet[:hexLen])
			rtcCounter := request.commandCount
			cmdName := "StopMC"
			if cmdCode == 0x36 {
				cmdName = "StartMC"
			}
			fmt.Printf("[MC_DEBUG_TX] %s: packet=%s (len=%d), rtcCounter=%d, cmdCode=0x%02X, pendingKey=%d\n",
				cmdName, txHex, len(packet), rtcCounter, cmdCode, rtcCounter)

			if txCmdCountValue != 0 && txCmdCountValue != rtcCounter {
				fmt.Printf("[MC_DEBUG_TX] %s: txRTCCommandCountMismatch extracted=%s pendingKey=%d\n",
					cmdName, txCmdCount, rtcCounter)
			}
		}
	}

	ctx, cancel := context.WithTimeout(requestManager.ctx, 100*time.Millisecond)
	defer cancel()

	return requestManager.transport.SendPacket(ctx, packet)
}

// requestKind is inferred from the flags on the request.
func (request *pendingRequest) kind() string {
	switch {
	case request.isStatusRequest:
		return "status"
	case request.isMCRequest:
		return "mc"
	default:
		return "rtc"
	}
}

func (requestManager *RequestManager) addPendingLocked(request *pendingRequest) {
	switch request.kind() {
	case "status":
		requestManager.pendingStatusRequest = request
	case "mc":
		requestManager.pendingMCRequests[request.mcCounter] = request
	case "rtc":
		requestManager.pendingRequests[request.commandCount] = request
	}
}

func (requestManager *RequestManager) removePendingLocked(request *pendingRequest) {
	switch request.kind() {
	case "status":
		if requestManager.pendingStatusRequest == request {
			requestManager.pendingStatusRequest = nil
		}
	case "mc":
		if cur, ok := requestManager.pendingMCRequests[request.mcCounter]; ok && cur == request {
			delete(requestManager.pendingMCRequests, request.mcCounter)
		}
	case "rtc":
		if cur, ok := requestManager.pendingRequests[request.commandCount]; ok && cur == request {
			delete(requestManager.pendingRequests, request.commandCount)
		}
	}
}

func (requestManager *RequestManager) handleTick() {
	// Read currentRequest with lock to avoid race conditions
	requestManager.mu.RLock()
	currentReq := requestManager.currentRequest
	requestManager.mu.RUnlock()

	if currentReq == nil {
		return
	}

	now := time.Now()
	ev := requestManager.evaluateCurrentRequest(currentReq, now)

	if !ev.exists {
		// Response handler already cleared it; just move on.
		requestManager.mu.Lock()
		requestManager.currentRequest = nil
		requestManager.mu.Unlock()
		return
	}

	if ev.expired || ev.retriesExhausted {
		// Request expired without a response; surface timeout and clear pending state.
		requestManager.completeRequestWithTimeout(ev.request)
		requestManager.mu.Lock()
		requestManager.currentRequest = nil
		requestManager.mu.Unlock()
		return
	}

	if !ev.shouldResend {
		// Still waiting either for response or per-attempt timeout; nothing to do this tick.
		return
	}

	// We should resend once.
	if err := requestManager.sendRequestPacket(ev.request, ev.request.sendCounter+1); err != nil {
		// Treat send failures as fatal for this request.
		requestManager.completeRequestWithError(ev.request, errors.Wrap(err, "send request"))
		requestManager.mu.Lock()
		requestManager.currentRequest = nil
		requestManager.mu.Unlock()
		return
	}

	requestManager.markRequestResent(ev.request)
}

type requestEvaluation struct {
	request          *pendingRequest
	exists           bool
	shouldResend     bool
	expired          bool
	retriesExhausted bool
}

// evaluateCurrentRequest inspects the current pending request under a read lock
// and classifies what should happen on this tick.
func (requestManager *RequestManager) evaluateCurrentRequest(current *pendingRequest, now time.Time) requestEvaluation {
	ev := requestEvaluation{}

	if current == nil {
		return ev
	}

	requestManager.mu.RLock()
	defer requestManager.mu.RUnlock()

	var (
		req    *pendingRequest
		exists bool
	)

	if current.isMCRequest {
		req, exists = requestManager.pendingMCRequests[current.mcCounter]
	} else if current.isStatusRequest {
		req, exists = requestManager.pendingStatusRequest, requestManager.pendingStatusRequest != nil
	} else {
		req, exists = requestManager.pendingRequests[current.commandCount]
	}

	ev.request = req
	ev.exists = exists

	if !exists || req == nil {
		return ev
	}

	maxRetries := requestManager.config.MaxRetries

	// Use stored perAttemptTimeout from request (calculated before gate acquisition)
	perAttemptTimeout := req.perAttemptTimeout
	if perAttemptTimeout <= 0 {
		perAttemptTimeout = requestManager.config.DefaultTimeout
	}

	ev.expired = now.After(req.deadline)
	ev.retriesExhausted = maxRetries > 0 && req.sendCounter >= uint(maxRetries)

	// Only consider resending if we haven't already expired or exhausted retries.
	if !ev.expired && !ev.retriesExhausted {
		ev.shouldResend = now.Sub(req.lastSendTime) >= perAttemptTimeout &&
			(maxRetries == 0 || req.sendCounter < uint(maxRetries))
	}

	return ev
}

// completeRequestWithTimeout sends a timeout error and removes the request from the pending maps.
// Caller must hold no locks; this method manages locking itself.
func (requestManager *RequestManager) completeRequestWithTimeout(req *pendingRequest) {
	requestManager.mu.Lock()

	totalTimeout := req.totalTimeout
	if totalTimeout <= 0 {
		// Fallback: calculate from deadline if durations not set (shouldn't happen)
		totalTimeout = time.Until(req.deadline)
		if totalTimeout < 0 {
			totalTimeout = 0
		}
	}

	timeoutErr := requestManager.newTimeoutError(int(req.sendCounter), totalTimeout)

	// Use idempotent completion (only first completion wins)
	req.tryCompleteWithError(timeoutErr)

	if req.isMCRequest {
		delete(requestManager.pendingMCRequests, req.mcCounter)
	} else if req.isStatusRequest {
		requestManager.pendingStatusRequest = nil
	} else {
		delete(requestManager.pendingRequests, req.commandCount)
	}

	// Ensure sentCh is closed to prevent deadlocks (safe even if already closed)
	req.signalFirstSendComplete()

	kind := req.kind()
	cmdCount := req.commandCount
	mcCounter := req.mcCounter
	attempts := req.sendCounter

	requestManager.mu.Unlock()

	if requestManager.debug.Load() {
		requestManager.dumpRouterTrace(fmt.Sprintf("timeout kind=%s cmdCount=%d mcCounter=%d attempts=%d", kind, cmdCount, mcCounter, attempts))
	}
}

// completeRequestWithError completes the request with a non-timeout error and removes it from maps.
func (requestManager *RequestManager) completeRequestWithError(req *pendingRequest, err error) {
	requestManager.mu.Lock()
	defer requestManager.mu.Unlock()

	// Use idempotent completion (only first completion wins)
	req.tryCompleteWithError(err)

	if req.isMCRequest {
		delete(requestManager.pendingMCRequests, req.mcCounter)
	} else if req.isStatusRequest {
		requestManager.pendingStatusRequest = nil
	} else {
		delete(requestManager.pendingRequests, req.commandCount)
	}

	// Ensure sentCh is closed to prevent deadlocks (safe even if already closed)
	req.signalFirstSendComplete()
}

// newTimeoutError creates a RequestTimeoutError with connection info for diagnostics.
func (requestManager *RequestManager) newTimeoutError(attempts int, timeout time.Duration) *protocol_common.RequestTimeoutError {
	err := &protocol_common.RequestTimeoutError{
		Attempts: attempts,
		Timeout:  timeout,
	}
	// Try to get connection info from transport (if it implements ConnectionInfo)
	if connInfo, ok := requestManager.transport.(transport.ConnectionInfo); ok {
		err.LocalAddr = connInfo.LocalAddr()
		err.RemoteAddr = connInfo.RemoteAddr()
	}
	return err
}

// markRequestResent updates sendCounter and lastSendTime after a successful resend.
func (requestManager *RequestManager) markRequestResent(req *pendingRequest) {
	requestManager.mu.Lock()
	defer requestManager.mu.Unlock()

	req.sendCounter++
	req.lastSendTime = time.Now()
}

// rxLoop is the background receive goroutine.
func (requestManager *RequestManager) rxLoop() {
	defer requestManager.wg.Done()

	// Log once when rx goroutine starts (sanity check to confirm loop is running).
	if requestManager.debug.Load() {
		fmt.Printf("[RX_LOOP_START] rx goroutine started\n")
	}

	for {
		// Exit promptly if the manager context has been canceled.
		if requestManager.ctx.Err() != nil {
			return
		}

		if err := requestManager.rxOnce(); err != nil {
			// If the manager context was canceled, exit immediately.
			// This ensures rxLoop exits promptly when Stop() is called.
			if requestManager.ctx.Err() != nil {
				return
			}
			// Timeouts on the receive context are expected; continue looping.
			if err == context.DeadlineExceeded {
				continue
			}
			// For context.Canceled from rxOnce's timeout context, continue (not manager shutdown).
			if err == context.Canceled {
				continue
			}
			requestManager.recordRouterTrace(fmt.Sprintf("rxLoop transport error: %v", err))
			continue
		}
	}
}

// rxOnce receives a single packet (with timeout), parses it, and routes it.
// It returns an error only for transport / context issues; parse and routing
// errors are treated as "skip this packet".
func (requestManager *RequestManager) rxOnce() error {
	ctx, cancel := context.WithTimeout(requestManager.ctx, 1*time.Second)
	var packet []byte
	var rxFrom string
	var err error
	if recvWithAddr, ok := requestManager.transport.(transport.RecvPacketWithAddr); ok {
		packet, rxFrom, err = recvWithAddr.RecvPacketWithAddr(ctx)
	} else {
		packet, err = requestManager.transport.RecvPacket(ctx)
	}
	cancel()

	if err != nil {
		return err
	}
	receivedAt := time.Now()

	// RAW receive logging - cannot be skipped by gating logic
	if requestManager.debug.Load() {
		reqBits, repBits, _ := protocol_common.ReadPacketHeader(packet)
		headLen := len(packet)
		if headLen > 16 {
			headLen = 16
		}
		head16 := hex.EncodeToString(packet[:headLen])
		fromAddr := rxFrom
		if fromAddr == "" {
			fromAddr = "unknown"
		}
		fmt.Printf("[UDP_RX_RAW] from=%s len=%d head16=%s reqBits=0x%08X repBits=0x%08X\n",
			fromAddr, len(packet), head16, reqBits, repBits)
	}

	if requestManager.debug.Load() {
		requestManager.logUDPBindOnce()
		localAddr := ""
		remoteAddr := ""
		if info, ok := requestManager.udpInfo(); ok {
			localAddr = info.LocalAddr
			remoteAddr = info.RemoteAddr
		}

		reqBits, repBits, headerErr := protocol_common.ReadPacketHeader(packet)
		rxCmdCount := "NA"
		rxCmdCountValue := uint8(0)
		if cmdCount, err := protocol_rtc.ExtractRTCCommandCount(packet); err == nil {
			rxCmdCountValue = cmdCount
			rxCmdCount = fmt.Sprintf("%d(0x%02X)", cmdCount, cmdCount)
		}

		headLen := len(packet)
		if headLen > 64 {
			headLen = 64
		}
		tailLen := len(packet)
		if tailLen > 16 {
			tailLen = 16
		}
		rxHead := hex.EncodeToString(packet[:headLen])
		rxTail := hex.EncodeToString(packet[len(packet)-tailLen:])

		cmdCodeByteStr := "NA"
		hasMeaningfulStr := "NA"
		pendingReq := (*pendingRequest)(nil)
		if rxCmdCountValue != 0 {
			requestManager.mu.RLock()
			pendingReq = requestManager.pendingRequests[rxCmdCountValue]
			requestManager.mu.RUnlock()
		}
		if rtcOffset, err := protocol_common.ResponseBlockOffsetFromPacket(packet, protocol_common.RespBitRTCReplyData); err == nil {
			rtcSize := protocol_common.BlockSizes[protocol_common.RespBitRTCReplyData]
			if rtcOffset+rtcSize <= len(packet) {
				cmdOrStatus := packet[rtcOffset+protocol_rtc.RTCDataOffsetCmdOrStatus]
				cmdCodeByteStr = fmt.Sprintf("0x%02X", cmdOrStatus)
			}
			hasMeaningfulStr = fmt.Sprintf("%t", hasMeaningfulRTCReply(packet, pendingReq))
		}

		fromAddr := remoteAddr
		if rxFrom != "" {
			fromAddr = rxFrom
		}
		if headerErr == nil {
			fmt.Printf("[UDP_RX] remote=%s -> local=%s len=%d reqBits=0x%08X repBits=0x%08X extractedCmdCount=%s rtcCmdCodeByte=%s hasMeaningfulRTC=%s rx_head64=%s rx_tail16=%s\n",
				fromAddr, localAddr, len(packet), reqBits, repBits, rxCmdCount, cmdCodeByteStr, hasMeaningfulStr, rxHead, rxTail)
		} else {
			fmt.Printf("[UDP_RX] remote=%s -> local=%s len=%d headerErr=%v extractedCmdCount=%s rtcCmdCodeByte=%s hasMeaningfulRTC=%s rx_head64=%s rx_tail16=%s\n",
				fromAddr, localAddr, len(packet), headerErr, rxCmdCount, cmdCodeByteStr, hasMeaningfulStr, rxHead, rxTail)
		}
	}

	response, err := requestManager.parseLinMotResponse(packet, rxFrom, receivedAt)
	if err != nil {
		// Handle RTC parse errors specially
		requestManager.handleRTCParseError(err)
		return nil
	}
	if response == nil {
		// Parsing failed or response type not recognized; skip this packet.
		return nil
	}

	requestManager.routeResponse(response)
	return nil
}

// isValidRTCCounter returns true only for valid RTC counter values (1..14).
// Rejects 0 and 15 (invalid values that could indicate corruption or protocol violation).
func isValidRTCCounter(c uint8) bool {
	return c >= 1 && c <= 14
}

// prevRTCCounter returns the previous counter in the 1..14 ring.
// If c == 1, returns 14; otherwise returns c-1.
func prevRTCCounter(c uint8) uint8 {
	if c == 1 {
		return 14
	}
	return c - 1
}

// nextRTCCounter returns the next counter in the 1..14 ring.
// If c == 14, returns 1; otherwise returns c+1.
func nextRTCCounter(c uint8) uint8 {
	if c == 14 {
		return 1
	}
	return c + 1
}

// shouldResync determines if resync should occur based on counter validity and mismatch.
// Only resyncs when:
// - Both counters are valid (1..14)
// - driveCmdCount != clientCmdCount (mismatch detected)
// - driveCmdCount != prev(clientCmdCount) (not a stale previous echo)
// This prevents false positives from stale previous cmdCount echoes while allowing
// legitimate resync when the drive is truly desynced.
func shouldResync(driveCmdCount, clientCmdCount uint8) bool {
	if !isValidRTCCounter(driveCmdCount) || !isValidRTCCounter(clientCmdCount) {
		return false // Invalid counters => no resync
	}
	if driveCmdCount == clientCmdCount {
		return false // Counters match => no resync needed
	}
	if driveCmdCount == prevRTCCounter(clientCmdCount) {
		return false // Drive is echoing previous counter (stale echo) => no resync
	}
	return true // Valid mismatch that's not a stale echo => resync
}

// hasMeaningfulRTCReply returns true if the RTC reply payload appears meaningful.
// It treats a reply as meaningful when it matches a pending request or the cmd/status byte is non-zero.
func hasMeaningfulRTCReply(packet []byte, pendingReq *pendingRequest) bool {
	if pendingReq != nil {
		return true
	}
	rtcOffset, err := protocol_common.ResponseBlockOffsetFromPacket(packet, protocol_common.RespBitRTCReplyData)
	if err != nil {
		return false
	}
	rtcSize := protocol_common.BlockSizes[protocol_common.RespBitRTCReplyData]
	if rtcOffset+rtcSize > len(packet) {
		return false
	}
	cmdOrStatus := packet[rtcOffset+protocol_rtc.RTCDataOffsetCmdOrStatus]
	return cmdOrStatus != 0
}

// isStatusTelegram returns true if the packet is a status telegram suitable for validateConnectivity.
// A status telegram must have:
// - reqBits == 0 (no request flags set - key correlation signal)
// - repBits includes all Standard status bits (can be a superset, e.g., with monitoring or RTC reply)
// - pktLen >= 26 (minimum status packet size)
func isStatusTelegram(reqBits, repBits uint32, pktLen int) (bool, string) {
	const minStatusLen = 26
	standardMask := protocol_common.ResponseFlags.Standard // 0x0000007F (bits 0-6)

	if reqBits != 0 {
		return false, "nonzero_reqbits"
	}
	// Allow repBits to be a superset (e.g., 0x7F, 0xFF, 0x1FF) as long as it includes all Standard bits
	if (repBits & standardMask) != standardMask {
		return false, "missing_standard_bits"
	}
	if pktLen < minStatusLen {
		return false, "too_short"
	}
	return true, "ok"
}

// isControlWordTelegram returns true if the packet shape matches a control-word response.
// Control-word responses include reqBits with ControlWord set and at least StatusWord in repBits.
func isControlWordTelegram(reqBits, repBits uint32, pktLen int) (bool, string) {
	const minControlWordRespLen = 10 // Header (8) + StatusWord (2)

	if reqBits&protocol_common.RequestFlags.ControlWord == 0 {
		return false, "missing_controlword_reqbit"
	}
	if repBits&protocol_common.RespBitStatusWord == 0 {
		return false, "missing_statusword_bit"
	}
	if pktLen < minControlWordRespLen {
		return false, "too_short"
	}
	return true, "ok"
}

// acceptsPendingStatusPacket decides whether a packet can satisfy the current pending
// status-like request (Status/Monitoring/ConnectivityProbe/ControlWord).
func (requestManager *RequestManager) acceptsPendingStatusPacket(reqBits, repBits uint32, pktLen int) (bool, string) {
	requestManager.mu.RLock()
	pending := requestManager.pendingStatusRequest
	requestManager.mu.RUnlock()
	if pending == nil {
		return true, "no_pending_status"
	}

	// Control-word requests expect control-word-shaped responses (reqBits bit0 set).
	if _, isControlWordRequest := pending.request.(interface{ GetControlWord() uint16 }); isControlWordRequest {
		return isControlWordTelegram(reqBits, repBits, pktLen)
	}

	// All other status-like requests require standard status telegram shape.
	return isStatusTelegram(reqBits, repBits, pktLen)
}

func isPlausibleStatusPacket(reqBits uint32, repBits uint32, pktLen int, sw, sv, err, warn uint16) (bool, string) {
	_ = reqBits
	const minStatusLen = 26
	statusRepMask := protocol_common.ResponseFlags.Standard

	if pktLen < minStatusLen {
		return false, "too_short"
	}
	if (repBits & statusRepMask) != statusRepMask {
		return false, "missing_status_bits"
	}
	if sw == 0 && sv == 0 && err == 0 && warn == 0 {
		return false, "all_zero_status"
	}
	return true, "ok"
}

func (requestManager *RequestManager) evaluateStatusPacket(reqBits, repBits uint32, packet []byte, rxFrom string, rxAt time.Time) (*protocol_common.Status, bool, string) {
	status, err := protocol_common.ReadStatusDynamic(packet)
	if err != nil {
		status = &protocol_common.Status{}
	}

	sw := uint16(0)
	sv := uint16(0)
	errWord := uint16(0)
	warn := uint16(0)
	if status != nil {
		sw = status.StatusWord
		sv = status.StateVar
		errWord = status.ErrorCode
		warn = status.WarnWord
	}

	plausible, reason := isPlausibleStatusPacket(reqBits, repBits, len(packet), sw, sv, errWord, warn)
	if requestManager.debug.Load() {
		fromAddr := rxFrom
		if fromAddr == "" {
			fromAddr = "unknown"
		}
		packetHex := hex.EncodeToString(packet)
		isStatus, statusReason := isStatusTelegram(reqBits, repBits, len(packet))
		if isStatus {
			fmt.Printf("[STATUS_RX] from=%s len=%d reqBits=0x%08X repBits=0x%08X sw=0x%04X sv=0x%04X err=0x%04X warn=0x%04X plausible=%v reason=%s packet=%s\n",
				fromAddr, len(packet), reqBits, repBits, sw, sv, errWord, warn, plausible, reason, packetHex)
		} else {
			// Determine packet kind for logging
			kind := "unknown"
			if reqBits&protocol_common.RequestFlags.RTCCommand != 0 {
				kind = "rtc"
			} else if reqBits&protocol_common.RequestFlags.MotionControl != 0 {
				kind = "motion_control"
			} else if reqBits&protocol_common.RequestFlags.ControlWord != 0 {
				kind = "control_word"
			}
			fmt.Printf("[NONSTATUS_IGNORED] kind=%s reqBits=0x%08X repBits=0x%08X len=%d reason=%s\n",
				kind, reqBits, repBits, len(packet), statusReason)
		}
	}

	return status, plausible, reason
}

func (requestManager *RequestManager) shouldRejectStatusPacket(reqBits, repBits uint32, packet []byte, status *protocol_common.Status, plausible bool, reason string, rxFrom string, rxAt time.Time) bool {
	if plausible {
		return false
	}

	// Only treat as invalid status telegram if it's actually a status telegram
	isStatus, _ := isStatusTelegram(reqBits, repBits, len(packet))
	if !isStatus {
		// Not a status telegram, ignore it (don't count as invalid telemetry)
		return false
	}

	sw := uint16(0)
	sv := uint16(0)
	errWord := uint16(0)
	warn := uint16(0)
	if status != nil {
		sw = status.StatusWord
		sv = status.StateVar
		errWord = status.ErrorCode
		warn = status.WarnWord
	}

	var pending *pendingRequest
	var rejects uint8
	requestManager.mu.Lock()
	pending = requestManager.pendingStatusRequest
	if pending != nil {
		pending.invalidStatusCount++
		rejects = pending.invalidStatusCount
		sentHex := pending.statusTxHex
		sentAt := pending.lastSendTime
		ageMs := int64(-1)
		if !sentAt.IsZero() && !rxAt.IsZero() {
			ageMs = rxAt.Sub(sentAt).Milliseconds()
		}
		fromAddr := rxFrom
		if fromAddr == "" {
			fromAddr = "unknown"
		}
		snapshot := StatusPacketSnapshot{
			PacketHex:  hex.EncodeToString(packet),
			ReqBits:    reqBits,
			RepBits:    repBits,
			Len:        len(packet),
			StatusWord: sw,
			StateVar:   sv,
			ErrorCode:  errWord,
			WarnWord:   warn,
			Plausible:  plausible,
			Reason:     reason,
			FromAddr:   fromAddr,
			SentHex:    sentHex,
			SentAt:     sentAt,
			RxAt:       rxAt,
			AgeMs:      ageMs,
		}
		requestManager.lastInvalidStatusSnapshot = &snapshot
		if rejects >= maxInvalidStatusRejects {
			requestManager.pendingStatusRequest = nil
		}
	}
	requestManager.mu.Unlock()

	if pending != nil && rejects >= maxInvalidStatusRejects {
		pending.tryCompleteWithError(ErrInvalidStatusTelegram)
	}

	return pending != nil
}

// tryDeliverStatusLikeFromPacket delivers a status or monitoring response if present in the packet.
// This allows status handling to occur even when RTC reply data is present.
func (requestManager *RequestManager) tryDeliverStatusLikeFromPacket(reqBits, repBits uint32, packet []byte, status *protocol_common.Status, plausible bool, reason string, rxFrom string, rxAt time.Time) {
	// Harden validateConnectivity: if there's a pending status request, only accept strict status telegrams
	requestManager.mu.Lock()
	hasPendingStatusRequest := requestManager.pendingStatusRequest != nil
	requestManager.mu.Unlock()
	if hasPendingStatusRequest {
		ok, statusReason := requestManager.acceptsPendingStatusPacket(reqBits, repBits, len(packet))
		if !ok {
			// Not a status telegram, ignore it (don't deliver to pending status request)
			if requestManager.debug.Load() {
				fromAddr := rxFrom
				if fromAddr == "" {
					fromAddr = "unknown"
				}
				fmt.Printf("[PENDING_STATUS_IGNORED] reqBits=0x%08X repBits=0x%08X len=%d from=%s reason=%s\n",
					reqBits, repBits, len(packet), fromAddr, statusReason)
			}
			return
		}
	}

	var response protocol_common.Response
	var err error

	if repBits&protocol_common.RespBitMonitoringChannel != 0 {
		if requestManager.shouldRejectStatusPacket(reqBits, repBits, packet, status, plausible, reason, rxFrom, rxAt) {
			return
		}
		response, err = protocol_common.ReadMonitoringStatusResponse(packet)
	} else if repBits&protocol_common.RespBitStatusWord != 0 {
		if requestManager.shouldRejectStatusPacket(reqBits, repBits, packet, status, plausible, reason, rxFrom, rxAt) {
			return
		}
		response, err = protocol_common.ReadStatusResponse(packet)
	}
	if err != nil || response == nil {
		return
	}

	requestManager.deliverStatusLikeResponse(response)
}

// parseLinMotResponse inspects the packet header and parses the appropriate
// response type (Status, Monitoring, RTC, or Motion Control).
func (requestManager *RequestManager) parseLinMotResponse(packet []byte, rxFrom string, rxAt time.Time) (protocol_common.Response, error) {
	reqBits, repBits, err := protocol_common.ReadPacketHeader(packet)
	if err != nil {
		return nil, err
	}

	requestManager.mu.Lock()
	currentReq := requestManager.currentRequest
	currentIsStatus := false
	currentIsMC := false
	currentKind := "none"
	pendingRTCCount := len(requestManager.pendingRequests)
	if currentReq != nil {
		currentIsStatus = currentReq.isStatusRequest
		currentIsMC = currentReq.isMCRequest
		currentKind = currentReq.kind()
	}
	requestManager.mu.Unlock()

	traceWanted := requestManager.debug.Load() && (repBits&protocol_common.RespBitRTCReplyData != 0 ||
		repBits&protocol_common.RespBitMonitoringChannel != 0 || (currentReq != nil && currentReq.isMCRequest))
	traceCmdCount := uint8(0)
	traceHasCmdCount := false
	traceHasPendingRTC := false
	traceMeaningful := "NA"
	traceLastObserved := "NA"
	traceSeeded := "NA"
	traceInFlightCmdCount := "NA"
	traceResyncClientCmdCount := "NA"
	tracePendingCmdCode := "NA"
	traceAttempt := "NA"
	traceShouldResync := "NA"
	traceHasShouldResync := false
	traceRxHead := "NA"
	traceRxTail := "NA"
	traceRTCRelevant := requestManager.debug.Load() && (pendingRTCCount > 0 ||
		(currentReq != nil && !currentIsStatus && !currentIsMC))
	traceRTCOffset := "NA"
	traceRTCOffsetSource := "NA"
	traceHasRTCReplyBit := repBits&protocol_common.RespBitRTCReplyData != 0
	traceHasMonitoringBit := repBits&protocol_common.RespBitMonitoringChannel != 0
	traceCurrentCmdCode := "NA"

	if traceRTCRelevant && currentReq != nil && !currentIsStatus && !currentIsMC {
		traceCurrentCmdCode = fmt.Sprintf("0x%02X", extractRTCCommandCode(currentReq.request))
	}

	if traceWanted {
		headLen := len(packet)
		if headLen > 64 {
			headLen = 64
		}
		tailLen := len(packet)
		if tailLen > 16 {
			tailLen = 16
		}
		traceRxHead = hex.EncodeToString(packet[:headLen])
		traceRxTail = hex.EncodeToString(packet[len(packet)-tailLen:])
	}

	statusCandidate := repBits&protocol_common.RespBitStatusWord != 0 ||
		repBits&protocol_common.RespBitMonitoringChannel != 0
	status := (*protocol_common.Status)(nil)
	statusPlausible := false
	statusReason := "no_status_bits"
	if statusCandidate {
		status, statusPlausible, statusReason = requestManager.evaluateStatusPacket(reqBits, repBits, packet, rxFrom, rxAt)
	}

	if traceRTCRelevant {
		if rtcOffset, err := protocol_common.ResponseBlockOffsetFromPacket(packet, protocol_common.RespBitRTCReplyData); err == nil {
			traceRTCOffset = fmt.Sprintf("%d", rtcOffset)
			traceRTCOffsetSource = "resp_def"
		} else {
			traceRTCOffset = "invalid"
			traceRTCOffsetSource = "resp_def"
		}

		fmt.Printf("[RTC_PKT_DEF] reqBits=0x%08X repBits=0x%08X pendingRTC=%v currentCmd=%s hasRTCReplyBit=%v hasMonitoringBit=%v rtcOffset=%s rtcOffsetSource=%s\n",
			reqBits, repBits, pendingRTCCount > 0, traceCurrentCmdCode,
			traceHasRTCReplyBit, traceHasMonitoringBit, traceRTCOffset, traceRTCOffsetSource)
	}

	traceRouter := func(branch string) {
		if !traceWanted {
			return
		}
		cmdCount := "NA"
		existsRTC := "NA"
		shouldResync := "NA"
		if traceHasCmdCount {
			cmdCount = fmt.Sprintf("%d(0x%02X)", traceCmdCount, traceCmdCount)
			existsRTC = fmt.Sprintf("%v", traceHasPendingRTC)
		}
		if traceHasShouldResync {
			shouldResync = traceShouldResync
		}
		line := fmt.Sprintf("reqBits=0x%08X repBits=0x%08X branch=%s cmdCount=%s currentKind=%s isStatus=%v isMC=%v existsRTC=%s hasMeaningfulRTC=%s lastObservedDriveCmdCount=%s seeded=%s inFlightCmdCount=%s pendingCmdCode=%s resyncClientCmdCount=%s attempt=%s shouldResync=%s rx_head64=%s rx_tail16=%s",
			reqBits, repBits, branch, cmdCount, currentKind, currentIsStatus, currentIsMC, existsRTC, traceMeaningful, traceLastObserved, traceSeeded, traceInFlightCmdCount, tracePendingCmdCode, traceResyncClientCmdCount, traceAttempt, shouldResync, traceRxHead, traceRxTail)
		fmt.Printf("[ROUTER_RX] %s\n", line)
		if traceHasCmdCount {
			fmt.Printf("[RTC_RX_SUM] ts=%s branch=%s cmd=%s rx_cc=%s attempt=%s\n",
				time.Now().Format(time.RFC3339Nano), branch, tracePendingCmdCode, cmdCount, traceAttempt)
		}
		requestManager.recordRouterTrace(line)
	}

	// Routing logic based on what was REQUESTED (reqBits) and what was RETURNED (repBits):
	//
	// An RTC response is identified by:
	//   - Bit 8 set in RESPONSE definition (bytes 4-7): RTC reply data is present
	//   This handles both explicit RTC responses and status telegrams that include RTC reply data.
	//   We detect RTC replies based solely on response bit 8, not request bit 2, to catch
	//   RTC replies embedded in status telegrams where the request echo may not have bit 2 set.

	isRTCReplyPresent := repBits&protocol_common.RespBitRTCReplyData != 0 // Bit 8 in response = RTC reply data present

	if isRTCReplyPresent {
		requestManager.tryDeliverStatusLikeFromPacket(reqBits, repBits, packet, status, statusPlausible, statusReason, rxFrom, rxAt)
		// RTC reply data present in response => treat as RTC response.
		// This handles both explicit RTC responses and status telegrams that include RTC reply data.
		// For special commands, we need to look up the original command code
		// because byte29 in response contains status, not command code.
		cmdCount, err := protocol_rtc.ExtractRTCCommandCount(packet)
		if err != nil {
			return nil, err
		}
		if isValidRTCCounter(cmdCount) {
			requestManager.mu.Lock()
			requestManager.lastObservedDriveCmdCount = cmdCount
			requestManager.mu.Unlock()
		}
		if requestManager.debug.Load() {
			byte42 := "NA"
			if len(packet) > 42 {
				byte42 = fmt.Sprintf("0x%02X", packet[42])
			}
			last8 := packet
			if len(packet) > 8 {
				last8 = packet[len(packet)-8:]
			}
			rtcOffset := "NA"
			if offset, err := protocol_common.ResponseBlockOffsetFromPacket(packet, protocol_common.RespBitRTCReplyData); err == nil {
				rtcOffset = fmt.Sprintf("%d", offset)
			}
			fmt.Printf("[RTC_CMDCOUNT] len=%d byte42=%s last8=%s cmdCount=%d(0x%02X) offset=rtc_def(%s)\n",
				len(packet), byte42, hex.EncodeToString(last8), cmdCount, cmdCount, rtcOffset)
		}
		traceCmdCount = cmdCount
		traceHasCmdCount = true

		// Under lock, check if pending request exists and determine resync eligibility
		requestManager.mu.Lock()
		pendingReq, existsRTC := requestManager.pendingRequests[cmdCount]
		currentReq := requestManager.currentRequest
		hasInFlightRTC := currentReq != nil && !currentReq.isStatusRequest && !currentReq.isMCRequest
		resyncAttempted := requestManager.resyncAttempted
		oldCmdCount := uint8(0)
		attemptCount := uint(0)
		if currentReq != nil {
			oldCmdCount = currentReq.commandCount
			attemptCount = currentReq.sendCounter
		}
		var originalCmdCode uint8
		var originalUPID uint16
		if existsRTC {
			originalCmdCode = extractRTCCommandCode(pendingReq.request)
			originalUPID = pendingReq.originalUPID
			tracePendingCmdCode = fmt.Sprintf("0x%02X", originalCmdCode)
		}
		if currentReq != nil {
			traceInFlightCmdCount = fmt.Sprintf("%d(0x%02X)", oldCmdCount, oldCmdCount)
			traceResyncClientCmdCount = traceInFlightCmdCount
			traceAttempt = fmt.Sprintf("%d", attemptCount)
		}
		traceSeeded = fmt.Sprintf("%v", requestManager.rtcSeeded)
		if isValidRTCCounter(requestManager.lastObservedDriveCmdCount) {
			traceLastObserved = fmt.Sprintf("%d(0x%02X)", requestManager.lastObservedDriveCmdCount, requestManager.lastObservedDriveCmdCount)
		}
		requestManager.mu.Unlock()
		traceHasPendingRTC = existsRTC

		hasMeaningful := hasMeaningfulRTCReply(packet, pendingReq)
		traceMeaningful = fmt.Sprintf("%v", hasMeaningful)

		if requestManager.debug.Load() && currentReq != nil && !currentIsStatus && !currentIsMC && oldCmdCount != 0 && cmdCount != oldCmdCount {
			rtcOffset := -1
			rtcBytes := ""
			if offset, err := protocol_common.ResponseBlockOffsetFromPacket(packet, protocol_common.RespBitRTCReplyData); err == nil {
				rtcOffset = offset
				if rtcOffset+8 <= len(packet) {
					rtcBytes = hex.EncodeToString(packet[rtcOffset : rtcOffset+8])
				}
			}
			fmt.Printf("[RTC_MISMATCH_DUMP] tx_cc=%d rx_cc=%d reqBits=0x%08X repBits=0x%08X rtcOffset=%d rtcBytes=%s packet=%s\n",
				oldCmdCount, cmdCount, reqBits, repBits, rtcOffset, rtcBytes, hex.EncodeToString(packet))
		}

		if hasMeaningful {
			requestManager.lastDriveCmdCount = cmdCount // Track last cmdCount from drive
		}

		// Debug logging for ALL RTC RX packets when debug enabled
		if requestManager.debug.Load() {
			// Extract last 8 bytes (RTC data block)
			last8 := packet
			if len(packet) > 8 {
				last8 = last8[len(last8)-8:]
			}

			// Compute request echo bit 2 (RTC command sent)
			echoRTCBit2Present := reqBits&protocol_common.RequestFlags.RTCCommand != 0

			pendingCmdCode := "NA"
			if existsRTC {
				pendingCmdCode = fmt.Sprintf("0x%02X", originalCmdCode)
			}

			fmt.Printf("[RTC_RX] reqBits=0x%08X repBits=0x%08X isRTCReplyPresent=%v echoRTCBit2Present=%v cmdCount=%d(0x%02X) hasPendingRTC=%v pendingCmdCode=%s hasMeaningfulRTC=%v rxLast8=%s\n",
				reqBits, repBits, isRTCReplyPresent, echoRTCBit2Present,
				cmdCount, cmdCount, existsRTC, pendingCmdCode,
				hasMeaningful, hex.EncodeToString(last8))
		}

		if requestManager.debug.Load() && currentReq != nil && !currentIsStatus && !currentIsMC {
			currentCmd := extractRTCCommandCode(currentReq.request)
			if currentCmd == 0x70 && currentReq.commandCount != cmdCount {
				if rtcOffset, err := protocol_common.ResponseBlockOffsetFromPacket(packet, protocol_common.RespBitRTCReplyData); err == nil {
					if rtcOffset+8 <= len(packet) {
						fmt.Printf("[RTC_BLOCK] reason=cmd70_mismatch offset=%d bytes=%s\n",
							rtcOffset, hex.EncodeToString(packet[rtcOffset:rtcOffset+8]))
					}
				}
			}
		}

		if requestManager.debug.Load() {
			prevCount := prevRTCCounter(oldCmdCount)
			traceShouldResync = fmt.Sprintf("%v", shouldResync(cmdCount, oldCmdCount))
			traceHasShouldResync = true
			fmt.Printf("[RTC_ROUTE] driveCmdCount=%d clientCmdCount=%d prevClient=%d shouldResync=%v hasInFlightRTC=%v resyncAttempted=%v\n",
				cmdCount, oldCmdCount, prevCount, shouldResync(cmdCount, oldCmdCount), hasInFlightRTC, resyncAttempted)
		}

		if !hasMeaningful && hasInFlightRTC && attemptCount >= 2 && !resyncAttempted {
			requestManager.mu.Lock()
			beaconCount := requestManager.lastObservedDriveCmdCount
			if requestManager.currentRequest == currentReq && isValidRTCCounter(beaconCount) && shouldResync(beaconCount, oldCmdCount) {
				delete(requestManager.pendingRequests, oldCmdCount)
				desired := nextRTCCounter(beaconCount)
				newCmdCount := desired
				attempts := 0
				for _, exists := requestManager.pendingRequests[newCmdCount]; exists && attempts < 14; {
					newCmdCount = nextRTCCounter(newCmdCount)
					attempts++
				}
				if attempts >= 14 {
					newCmdCount = requestManager.rtcCount.Next()
					for newCmdCount == beaconCount || newCmdCount == oldCmdCount {
						newCmdCount = requestManager.rtcCount.Next()
					}
				}
				currentReq.commandCount = newCmdCount
				requestManager.pendingRequests[newCmdCount] = currentReq
				requestManager.resyncAttempted = true
				requestManager.rtcResyncCount.Add(1)
				requestManager.mu.Unlock()

				if requestManager.debug.Load() {
					fmt.Printf("[RTC_BEACON_RESYNC] beacon=%d inFlight=%d attempts=%d -> cmdCount=%d\n",
						beaconCount, oldCmdCount, attemptCount, newCmdCount)
				}

				if err := requestManager.sendRequestPacket(currentReq, currentReq.sendCounter+1); err == nil {
					requestManager.markRequestResent(currentReq)
				}
				traceRouter("rtc-beacon-resync")
				return nil, nil
			}
			requestManager.mu.Unlock()
		}

		// Decision tree: Only parse as RTC response if matching request exists OR resync occurred
		// Otherwise, fall through to normal status parsing to avoid swallowing status responses
		if hasMeaningful {
			if existsRTC {
				// Case 1: Matching RTC request exists => parse and return RTC response
				// Clear resync flag for next request
				requestManager.mu.Lock()
				requestManager.resyncAttempted = false
				requestManager.mu.Unlock()
				// Continue with RTC parsing below
			} else if hasInFlightRTC && shouldResync(cmdCount, oldCmdCount) && !resyncAttempted {
				// Case 2: In-flight RTC request, should resync => perform resync, then parse RTC
				// Drive is echoing a different cmdCount - we're desynced
				// Resync by re-keying the in-flight request with a new cmdCount
				requestManager.mu.Lock()
				if requestManager.currentRequest == currentReq && !requestManager.resyncAttempted {
					// Remove old entry
					delete(requestManager.pendingRequests, oldCmdCount)
					// Choose deterministic new cmdCount: next(driveCmdCount)
					// This ensures predictable resync behavior
					desired := nextRTCCounter(cmdCount)
					newCmdCount := desired
					// Check for collision with existing pending requests (paranoia guard)
					// In practice, with single-flight gating, there should only be one in-flight RTC
					// But guard against edge cases by iterating through the ring if needed
					attempts := 0
					for _, exists := requestManager.pendingRequests[newCmdCount]; exists && attempts < 14; {
						newCmdCount = nextRTCCounter(newCmdCount)
						attempts++
					}
					// If all 14 slots are occupied (shouldn't happen), fallback to rtcCount.Next()
					if attempts >= 14 {
						newCmdCount = requestManager.rtcCount.Next()
						for newCmdCount == cmdCount || newCmdCount == oldCmdCount {
							newCmdCount = requestManager.rtcCount.Next()
						}
					}
					// Update request with new cmdCount
					currentReq.commandCount = newCmdCount
					// Re-add to map with new key
					requestManager.pendingRequests[newCmdCount] = currentReq
					requestManager.resyncAttempted = true
					// Increment resync counter (pure instrumentation, always increments)
					requestManager.rtcResyncCount.Add(1)
					requestManager.mu.Unlock()

					if requestManager.debug.Load() {
						fmt.Printf("[RTC_DESYNC] detected: driveCmdCount=%d clientCmdCount=%d, resyncing to cmdCount=%d\n",
							cmdCount, oldCmdCount, newCmdCount)
					}

					// Resend with new cmdCount
					if err := requestManager.sendRequestPacket(currentReq, currentReq.sendCounter+1); err == nil {
						requestManager.mu.Lock()
						currentReq.lastSendTime = time.Now()
						requestManager.mu.Unlock()
					}
					// After resync, don't parse this unmatched response - it's from a previous cmdCount
					traceRouter("rtc-resync")
					return nil, nil
				} else {
					requestManager.mu.Unlock()
					// Resync didn't occur (request changed or already attempted) => fall through to status parsing
				}
			} else {
				// Case 3: No matching RTC request and no resync => fall through to status parsing
				// DO NOT treat as RTC response; let normal Status/Monitoring/MC parsing handle it
				// This ensures status requests still complete even when RTC reply data is present
			}
		}

		// Only continue with RTC parsing if we have a matching request (existsRTC is true)
		// If we reach here without existsRTC, we should have fallen through to status parsing
		if !existsRTC {
			// Fall through to normal status parsing (don't return early, don't parse as RTC)
			// This ensures status requests still complete even when RTC reply data is present
		} else {
			// Matched response - continue with RTC parsing
			// Debug logging for StopMC/StartMC (0x35/0x36) when debug enabled
			if requestManager.debug.Load() && (originalCmdCode == 0x35 || originalCmdCode == 0x36) {
				hexLen := len(packet)
				if hexLen > 64 {
					hexLen = 64
				}
				rxHex := hex.EncodeToString(packet[:hexLen])
				cmdName := "StopMC"
				if originalCmdCode == 0x36 {
					cmdName = "StartMC"
				}
				mapType := "RTC"
				fmt.Printf("[MC_DEBUG_RX] %s: packet=%s (len=%d), rtcCounter=%d, cmdCode=0x%02X, pendingKey=%d, mapType=%s, exists=%v\n",
					cmdName, rxHex, len(packet), cmdCount, originalCmdCode, cmdCount, mapType, existsRTC)
			}

			traceRouter("rtc")

			// Parse response with command code and UPID (needed for registry lookup).
			// Following C# library pattern: trust counter match, use request UPID for
			// response type determination (ignore response UPID which may be stale).
			resp, err := protocol_rtc.ReadRTCResponse(packet, originalCmdCode, originalUPID)
			if err != nil {
				return nil, err
			}
			return resp.(protocol_common.Response), nil
		}
		// If !existsRTC, fall through to normal status parsing below
	}

	if repBits&protocol_motion_control.StateVarResponse == protocol_motion_control.StateVarResponse && reqBits&protocol_common.RequestFlags.MotionControl != 0 {
		// Motion Control response - requires StatusWord (bit 0) and StateVar (bit 1)
		// May include additional bits like DemandPosition (bit 3)
		resp, err := protocol_motion_control.ReadMCResponse(packet)
		if err != nil {
			return nil, err
		}
		// Update last-seen StateVarLow for MC counter calculation (matches linudp.cs behavior)
		if status := resp.Status(); status != nil {
			requestManager.updateLastStateVarLow(status.StateVar)
		}
		traceRouter("mc")
		return resp, nil
	}

	if repBits != 0 {
		if statusCandidate && requestManager.shouldRejectStatusPacket(reqBits, repBits, packet, status, statusPlausible, statusReason, rxFrom, rxAt) {
			traceRouter("status_reject")
			return nil, nil
		}
		// Harden validateConnectivity: if there's a pending status request, only accept strict status telegrams
		requestManager.mu.Lock()
		hasPendingStatusRequest := requestManager.pendingStatusRequest != nil
		requestManager.mu.Unlock()
		if hasPendingStatusRequest {
			ok, statusReason := requestManager.acceptsPendingStatusPacket(reqBits, repBits, len(packet))
			if !ok {
				// Not a status telegram, ignore it (don't deliver to pending status request)
				if requestManager.debug.Load() {
					fromAddr := rxFrom
					if fromAddr == "" {
						fromAddr = "unknown"
					}
					fmt.Printf("[PENDING_STATUS_IGNORED] reqBits=0x%08X repBits=0x%08X len=%d from=%s reason=%s\n",
						reqBits, repBits, len(packet), fromAddr, statusReason)
				}
				traceRouter("nonstatus_ignored")
				return nil, nil
			}
		}
		// Check if this is a monitoring status response (bit 7 set)
		if repBits&protocol_common.RespBitMonitoringChannel != 0 {
			resp, err := protocol_common.ReadMonitoringStatusResponse(packet)
			if err != nil {
				return nil, err
			}
			// Update last-seen StateVarLow for MC counter calculation (matches linudp.cs behavior)
			if status := resp.Status(); status != nil && (repBits&protocol_common.RespBitStateVar != 0) {
				requestManager.updateLastStateVarLow(status.StateVar)
			}
			traceRouter("monitoring")
			return resp, nil
		} else {
			// Regular status response
			resp, err := protocol_common.ReadStatusResponse(packet)
			if err != nil {
				return nil, err
			}
			// Update last-seen StateVarLow for MC counter calculation (matches linudp.cs behavior)
			if status := resp.Status(); status != nil && (repBits&protocol_common.RespBitStateVar != 0) {
				requestManager.updateLastStateVarLow(status.StateVar)
			}
			traceRouter("status")
			return resp, nil
		}
	}

	// Invalid response flags (all zeros) – not a usable response.
	traceRouter("none")
	return nil, nil
}

// routeResponse dispatches a successfully parsed response to the appropriate
// pending request (status, monitoring, RTC, or motion control).
func (requestManager *RequestManager) routeResponse(response protocol_common.Response) {
	switch typed := response.(type) {
	case *protocol_common.StatusResponse:
		requestManager.deliverStatusLikeResponse(response)

	case *protocol_common.MonitoringStatusResponse:
		requestManager.deliverStatusLikeResponse(response)

	case protocol_rtc.RTCResponse:
		requestManager.deliverRTCResponse(typed, response)

	case protocol_motion_control.MCResponse:
		requestManager.deliverMCResponse(typed, response)
	}
}

// deliverStatusLikeResponse delivers either a StatusResponse or a
// MonitoringStatusResponse to the single pending status request, if any.
func (requestManager *RequestManager) deliverStatusLikeResponse(response protocol_common.Response) {
	requestManager.mu.Lock()
	request := requestManager.pendingStatusRequest
	if request != nil {
		requestManager.pendingStatusRequest = nil
	}
	requestManager.mu.Unlock()

	if request == nil {
		return
	}

	// Use idempotent completion (only first completion wins)
	request.tryCompleteWithResponse(response)
}

// deliverRTCResponse delivers an RTC response (read or write) to the pending
// request keyed by RTCCounter.
func (requestManager *RequestManager) deliverRTCResponse(rtcResp protocol_rtc.RTCResponse, response protocol_common.Response) {
	cmdCount := rtcResp.RTCCounter()

	requestManager.mu.Lock()
	request, exists := requestManager.pendingRequests[cmdCount]
	if exists {
		delete(requestManager.pendingRequests, cmdCount)
	}
	requestManager.mu.Unlock()

	if !exists || request == nil {
		return
	}

	// Use idempotent completion (only first completion wins)
	request.tryCompleteWithResponse(response)
}

// deliverMCResponse delivers a motion-control response to the pending request
// keyed by MCCounter.
//
// Some LinMot firmware versions do not echo the MC counter in StateVarLow bits 0-3
// (the nibble remains 0 regardless of the sent counter). In that case, the counter-based
// lookup will fail. Since the request gate ensures only one MC request is in-flight at a
// time, we fall back to delivering to the single pending MC request.
func (requestManager *RequestManager) deliverMCResponse(mcResp protocol_motion_control.MCResponse, response protocol_common.Response) {
	mcCounter := mcResp.MCCounter()

	requestManager.mu.Lock()
	request, exists := requestManager.pendingMCRequests[mcCounter]
	if exists {
		delete(requestManager.pendingMCRequests, mcCounter)
	}

	// Fallback: if counter-based lookup failed (e.g., drive doesn't echo MC counter
	// in StateVarLow), deliver to the single pending MC request. The request gate
	// ensures at most one MC request is in-flight.
	if !exists {
		for key, req := range requestManager.pendingMCRequests {
			request = req
			exists = true
			delete(requestManager.pendingMCRequests, key)
			break
		}
	}
	requestManager.mu.Unlock()

	if !exists || request == nil {
		return
	}

	// Use idempotent completion (only first completion wins)
	request.tryCompleteWithResponse(response)
}

// handleRTCParseError notifies any waiting request when parsing fails.
// For MissingResponseRegistryError, the counter is included in the error and is routed to the pending request.
func (requestManager *RequestManager) handleRTCParseError(err error) {
	if mfErr, ok := err.(*protocol_rtc.MissingResponseRegistryError); ok {
		requestManager.deliverRTCError(mfErr.Counter, err)
	}
}

// deliverRTCError sends an error to the pending request (if any) and cleans it up.
func (requestManager *RequestManager) deliverRTCError(cmdCount uint8, err error) bool {
	requestManager.mu.Lock()
	request, exists := requestManager.pendingRequests[cmdCount]
	if exists {
		delete(requestManager.pendingRequests, cmdCount)
	}
	requestManager.mu.Unlock()

	if !exists || request == nil {
		return false
	}

	// Use idempotent completion (only first completion wins)
	request.tryCompleteWithError(err)
	return true
}

// SetRTCCounterForTesting sets the RTC counter value for testing purposes only.
// This should only be used in test code.
func (requestManager *RequestManager) SetRTCCounterForTesting(value uint8) {
	requestManager.rtcCount.SetForTesting(value)
}

// LastDriveCmdCountForTesting returns the last cmdCount seen from the drive.
// This is a test-only function that should not be used in production code.
func (requestManager *RequestManager) LastDriveCmdCountForTesting() uint8 {
	requestManager.mu.RLock()
	defer requestManager.mu.RUnlock()
	return requestManager.lastDriveCmdCount
}

// PendingMCRequestCountForTesting returns the number of pending MC requests.
// This is a test-only function that should not be used in production code.
func (requestManager *RequestManager) PendingMCRequestCountForTesting() int {
	requestManager.mu.RLock()
	defer requestManager.mu.RUnlock()
	return len(requestManager.pendingMCRequests)
}
