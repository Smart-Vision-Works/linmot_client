package client_common

import (
	"context"
	"encoding/binary"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	protocol_common "gsail-go/linmot/protocol/common"
	protocol_rtc "gsail-go/linmot/protocol/rtc"
	protocol_command_tables "gsail-go/linmot/protocol/rtc/command_tables"
	transport "gsail-go/linmot/transport"
)

// blockingTransport wraps a transport.Client and can block SendPacket calls
// for testing single-flight gating behavior.
type blockingTransport struct {
	transport.Client
	mu             sync.Mutex
	sendCount      atomic.Int32
	blockSends     atomic.Bool
	unblockChannel chan struct{}
	blockedSends   sync.WaitGroup
}

func newBlockingTransport(base transport.Client) *blockingTransport {
	return &blockingTransport{
		Client:         base,
		unblockChannel: make(chan struct{}),
	}
}

func (bt *blockingTransport) SendPacket(ctx context.Context, data []byte) error {
	bt.sendCount.Add(1)

	if bt.blockSends.Load() {
		bt.blockedSends.Add(1)
		// Block until unblocked - use mutex to safely read unblockChannel
		bt.mu.Lock()
		ch := bt.unblockChannel
		bt.mu.Unlock()

		select {
		case <-ch:
			// Unblocked
		case <-ctx.Done():
			bt.blockedSends.Done()
			return ctx.Err()
		}
		bt.blockedSends.Done()
	}

	return bt.Client.SendPacket(ctx, data)
}

func (bt *blockingTransport) setBlockSends(block bool) {
	bt.blockSends.Store(block)
}

func (bt *blockingTransport) unblock() {
	bt.mu.Lock()
	oldCh := bt.unblockChannel
	bt.unblockChannel = make(chan struct{})
	bt.mu.Unlock()

	close(oldCh)
	bt.blockedSends.Wait()
}

func (bt *blockingTransport) getSendCount() int32 {
	return bt.sendCount.Load()
}

func TestRequestManager_SingleFlight_Gating(t *testing.T) {
	// Create mock transport with send signal
	baseClient, server := transport.NewMockTransportClientWithServer()
	signalClient := newSendSignalTransport(baseClient)

	// Create RequestManager
	rm := NewRequestManager(signalClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Setup beforeFirstSend hook to pause request A BEFORE it sends
	blockCh := make(chan struct{})
	hookHit := make(chan struct{})
	requestCount := atomic.Int32{}

	rm.beforeFirstSend = func(req *pendingRequest) {
		count := requestCount.Add(1)
		if count == 1 {
			// First request (A) - pause it
			close(hookHit) // Signal that hook was hit
			<-blockCh      // Block until unblocked
		}
		// Second request (B) - let it proceed (but it should be blocked at gate)
	}
	defer func() {
		rm.beforeFirstSend = nil
	}()

	// Start request A in background - it will block in beforeFirstSend hook
	ctxA, cancelA := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelA()
	reqA := protocol_common.NewStatusRequest()

	var wg sync.WaitGroup
	var errA error
	var responseA *protocol_common.StatusResponse

	wg.Add(1)
	go func() {
		defer wg.Done()
		var resp *protocol_common.StatusResponse
		resp, errA = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctxA, reqA)
		responseA = resp
	}()

	// Wait for hook to be hit (request A is paused in beforeFirstSend)
	select {
	case <-hookHit:
		// Hook hit, request A is paused
	case <-time.After(200 * time.Millisecond):
		t.Fatal("beforeFirstSend hook was not hit within 200ms")
	}

	// Verify request A has NOT sent yet (it's paused in hook)
	if signalClient.sendCount.Load() != 0 {
		t.Errorf("Expected 0 SendPacket calls while request A is paused, got %d", signalClient.sendCount.Load())
	}

	// Start request B concurrently - it should block at gate acquisition in roundTripBase
	ctxB, cancelB := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelB()
	reqB := protocol_common.NewStatusRequest()

	var errB error
	var responseB *protocol_common.StatusResponse

	wg.Add(1)
	go func() {
		defer wg.Done()
		responseB, errB = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctxB, reqB)
	}()

	// Give request B time to reach gate acquisition (should block there)
	time.Sleep(50 * time.Millisecond)

	// Verify still only 0 sends (A is paused in hook, B is blocked at gate)
	if signalClient.sendCount.Load() != 0 {
		t.Errorf("Expected 0 SendPacket calls (A paused in hook, B at gate), got %d", signalClient.sendCount.Load())
	}

	// Prepare responses
	statusA := &protocol_common.Status{
		ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar,
		StatusWord:   0x0001,
		StateVar:     0x0002,
	}
	statusRespA := protocol_common.NewStatusResponse(statusA)
	packetA, _ := statusRespA.WritePacket()

	statusB := &protocol_common.Status{
		ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar,
		StatusWord:   0x0001,
		StateVar:     0x0002,
	}
	statusRespB := protocol_common.NewStatusResponse(statusB)
	packetB, _ := statusRespB.WritePacket()

	// Unblock request A by closing blockCh
	close(blockCh)

	// Wait for request A's SendPacket to be called (deterministic)
	if !signalClient.waitForSend(100 * time.Millisecond) {
		t.Fatal("Request A SendPacket was not called within 100ms after unblock")
	}

	// Verify exactly 1 send (request A only)
	sendCountA := signalClient.sendCount.Load()
	if sendCountA != 1 {
		t.Errorf("Expected exactly 1 SendPacket call (request A), got %d (request B may have sent too early)", sendCountA)
	}

	// Send response for request A immediately (before timeout)
	server.SendPacket(packetA)

	// Wait briefly for request A to complete and release gate
	time.Sleep(10 * time.Millisecond)

	// Now request B should have acquired gate and sent
	// Wait for second send (deterministic)
	if !signalClient.waitForSendCount(2, 100*time.Millisecond) {
		t.Fatalf("Request B did not send within 100ms, sendCount=%d", signalClient.sendCount.Load())
	}

	// Verify both sends occurred
	if signalClient.sendCount.Load() != 2 {
		t.Errorf("Expected 2 SendPacket calls (A and B), got %d", signalClient.sendCount.Load())
	}

	// Send response for request B immediately
	server.SendPacket(packetB)

	// Wait for both requests to complete
	wg.Wait()

	// Verify both requests succeeded
	if errA != nil {
		t.Errorf("Request A failed: %v", errA)
	}
	if responseA == nil {
		t.Error("Request A returned nil response")
	}
	if errB != nil {
		t.Errorf("Request B failed: %v", errB)
	}
	if responseB == nil {
		t.Error("Request B returned nil response")
	}

	// Verify both requests sent (total should be 2)
	finalSends := signalClient.sendCount.Load()
	if finalSends != 2 {
		t.Errorf("Expected 2 SendPacket calls (A and B), got %d", finalSends)
	}
}

func TestRequestManager_StartStopLifecycle(t *testing.T) {
	baseClient, _ := transport.NewMockTransportClientWithServer()
	rm := NewRequestManager(baseClient, 10*time.Millisecond)

	if rm.ctx != nil {
		t.Fatalf("expected ctx to be nil before Start, got %v", rm.ctx)
	}

	rm.Start()
	if rm.cancel == nil {
		t.Fatal("Start() should set cancel function")
	}
	if rm.ctx == nil {
		t.Fatal("Start() should set context")
	}
	if err := rm.ctx.Err(); err != nil {
		t.Fatalf("context should not be canceled immediately after Start, got %v", err)
	}

	startCtx := rm.ctx
	if err := rm.Stop(); err != nil {
		t.Fatalf("Stop() returned error: %v", err)
	}
	if startCtx.Err() != context.Canceled {
		t.Fatalf("expected context canceled after Stop(), got %v", startCtx.Err())
	}
	if rm.ctx != nil {
		t.Fatalf("expected ctx to be nil after Stop(), got %v", rm.ctx)
	}
	// Stop should be safe to call again
	if err := rm.Stop(); err != nil {
		t.Fatalf("Stop() should be idempotent, got error: %v", err)
	}
}

func TestRequestManager_CancelWhileWaitingForGate(t *testing.T) {
	// Create mock transport with send signal for deterministic coordination
	baseClient, server := transport.NewMockTransportClientWithServer()
	signalClient := newSendSignalTransport(baseClient)
	blockingClient := newBlockingTransport(signalClient)

	// Create RequestManager with large timeouts for race mode
	rm := NewRequestManagerWithConfig(blockingClient, 10*time.Millisecond, RequestManagerConfig{
		MaxRetries:     10,
		DefaultTimeout: 500 * time.Millisecond, // Large timeout for race mode
	})
	rm.Start()
	defer rm.Stop()

	// Block sends so request A will block
	blockingClient.setBlockSends(true)

	// Start request A in background (will block on SendPacket)
	ctxA := context.Background()
	reqA := protocol_common.NewStatusRequest()

	var wg sync.WaitGroup
	var wgB sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctxA, reqA)
	}()

	// Wait for request A's SendPacket to be called (deterministic)
	// blockingTransport increments count before blocking, so use blockingClient.getSendCount()
	deadline := time.Now().Add(200 * time.Millisecond)
	for time.Now().Before(deadline) {
		if blockingClient.getSendCount() >= 1 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if blockingClient.getSendCount() < 1 {
		t.Fatal("Request A SendPacket was not called within 200ms")
	}

	// Start request B with cancelable context
	ctxB, cancelB := context.WithCancel(context.Background())
	reqB := protocol_common.NewStatusRequest()

	var errB error

	wg.Add(1)
	wgB.Add(1)
	go func() {
		defer wg.Done()
		defer wgB.Done()
		_, errB = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctxB, reqB)
	}()

	// Wait for request B to reach gate (deterministic: it should block at gate, not send)
	// Give it a brief moment to reach gate acquisition
	time.Sleep(20 * time.Millisecond)

	// Verify request B hasn't sent yet (waiting for gate)
	if blockingClient.getSendCount() != 1 {
		t.Errorf("Expected 1 SendPacket call while request B waits, got %d", blockingClient.getSendCount())
	}

	// Cancel request B's context
	cancelB()

	// Wait for request B to complete with cancellation
	wgB.Wait()

	// Verify request B returned context.Canceled
	if errB != context.Canceled {
		t.Errorf("Expected context.Canceled for request B, got %v", errB)
	}

	// Verify request B never sent a packet
	if blockingClient.getSendCount() != 1 {
		t.Errorf("Expected 1 SendPacket call (only request A), got %d", blockingClient.getSendCount())
	}

	// Unblock and complete request A
	blockingClient.setBlockSends(false)
	blockingClient.unblock()

	statusA := &protocol_common.Status{
		ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar,
		StatusWord:   0x0001,
		StateVar:     0x0002,
	}
	statusRespA := protocol_common.NewStatusResponse(statusA)
	packetA, _ := statusRespA.WritePacket()
	server.SendPacket(packetA)

	// Wait for request A to complete
	wg.Wait()

	// Verify gate is released (can start new request C)
	// Now that blocking is disabled, use blockingClient.getSendCount() to detect when request C sends
	ctxC, cancelC := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelC()
	reqC := protocol_common.NewStatusRequest()

	var errC error
	var responseC *protocol_common.StatusResponse

	// Prepare response for request C
	statusC := &protocol_common.Status{
		ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar,
		StatusWord:   0x0001,
		StateVar:     0x0002,
	}
	statusRespC := protocol_common.NewStatusResponse(statusC)
	packetC, _ := statusRespC.WritePacket()

	// Start request C and wait for it to send (deterministic)
	var wgC sync.WaitGroup
	wgC.Add(1)
	go func() {
		defer wgC.Done()
		responseC, errC = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctxC, reqC)
	}()

	// Wait for request C's SendPacket to be called (blockingClient will see it)
	deadlineC := time.Now().Add(200 * time.Millisecond)
	for time.Now().Before(deadlineC) {
		if blockingClient.getSendCount() >= 2 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if blockingClient.getSendCount() < 2 {
		t.Fatal("Request C SendPacket was not called within 200ms")
	}

	// Immediately send response for request C
	server.SendPacket(packetC)

	// Wait for request C to complete
	wgC.Wait()

	if errC != nil {
		t.Errorf("Expected new request to succeed after gate release, got %v", errC)
	}
	if responseC == nil {
		t.Error("Expected non-nil response for request C")
	}
}

// blockingSendTransport wraps a transport.Client and blocks SendPacket until ctx.Done()
type blockingSendTransport struct {
	transport.Client
	blockUntilDone context.Context
}

func newBlockingSendTransport(base transport.Client, blockUntilDone context.Context) *blockingSendTransport {
	return &blockingSendTransport{
		Client:         base,
		blockUntilDone: blockUntilDone,
	}
}

func (bt *blockingSendTransport) SendPacket(ctx context.Context, data []byte) error {
	// Block until blockUntilDone is canceled
	// Don't check ctx.Done() here - we want to test that timeout runs even if SendPacket is blocked
	// The RequestManager timeout should trigger before the context timeout
	<-bt.blockUntilDone.Done()
	return bt.Client.SendPacket(ctx, data)
}

// sendSignalTransport wraps a transport.Client and signals when SendPacket is called
type sendSignalTransport struct {
	transport.Client
	mu        sync.Mutex
	sendChs   []chan struct{} // Channels to close on each SendPacket call
	sendCount atomic.Int32
}

func newSendSignalTransport(base transport.Client) *sendSignalTransport {
	return &sendSignalTransport{
		Client:  base,
		sendChs: make([]chan struct{}, 0),
	}
}

func (st *sendSignalTransport) SendPacket(ctx context.Context, data []byte) error {
	st.sendCount.Add(1)
	st.mu.Lock()
	// Close all pending channels
	for _, ch := range st.sendChs {
		close(ch)
	}
	st.sendChs = st.sendChs[:0] // Clear
	st.mu.Unlock()
	return st.Client.SendPacket(ctx, data)
}

func (st *sendSignalTransport) waitForSend(timeout time.Duration) bool {
	ch := make(chan struct{})
	st.mu.Lock()
	// If already sent, return immediately
	if st.sendCount.Load() > 0 {
		st.mu.Unlock()
		return true
	}
	st.sendChs = append(st.sendChs, ch)
	st.mu.Unlock()
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (st *sendSignalTransport) waitForSendCount(target int32, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if st.sendCount.Load() >= target {
			return true
		}
		time.Sleep(5 * time.Millisecond)
	}
	return false
}

// TestRequestManager_SaveCommandTableAllowed verifies that the Save to Flash
// command (0x80) is allowed through the request manager (tripwire removed).
func TestRequestManager_SaveCommandTableAllowed(t *testing.T) {
	baseClient, _ := transport.NewMockTransportClientWithServer()
	signalClient := newSendSignalTransport(baseClient)

	rm := NewRequestManager(signalClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	request := protocol_command_tables.NewSaveCommandTableRequest()
	_, err := SendRequestAndReceive[*protocol_command_tables.SaveCommandTableResponse](rm, ctx, request)
	// The command should be sent (not blocked). It will timeout because there's
	// no real drive, but the error should NOT be "FLASH TRIPWIRE".
	if err != nil && strings.Contains(err.Error(), "FLASH TRIPWIRE") {
		t.Fatal("Save to Flash command should not be blocked by tripwire")
	}
	if signalClient.sendCount.Load() == 0 {
		t.Fatal("expected at least one packet sent for Save to Flash command")
	}
}

// TestRequestManager_GateWaitDoesNotBurnTimeoutBudget verifies that waiting for
// requestGate does NOT consume timeout budget. Timeout only starts when deadline
// is stamped in startRequest(), after gate is acquired.
func TestRequestManager_GateWaitDoesNotBurnTimeoutBudget(t *testing.T) {
	// Create mock transport with send signal
	baseClient, server := transport.NewMockTransportClientWithServer()
	signalClient := newSendSignalTransport(baseClient)

	// Create RequestManager with large timeout for race mode
	// Total timeout = DefaultTimeout * MaxRetries = 500ms * 10 = 5s
	// This ensures timeout is much larger than gate hold time (700ms)
	rm := NewRequestManagerWithConfig(signalClient, 10*time.Millisecond, RequestManagerConfig{
		MaxRetries:     10,
		DefaultTimeout: 500 * time.Millisecond, // Large timeout for race mode
	})
	rm.Start()
	defer rm.Stop()

	// Manually acquire and hold the gate token
	token := <-rm.requestGate

	// Start a request in a goroutine - it will block waiting for the gate
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := protocol_common.NewStatusRequest()

	var wg sync.WaitGroup
	var err error
	var response *protocol_common.StatusResponse

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctx, req)
	}()

	// Hold gate for 700ms (much longer than per-attempt timeout, but less than total timeout)
	// If timeout budget were consumed during gate wait, request would timeout
	gateHoldDone := make(chan struct{})
	go func() {
		time.Sleep(700 * time.Millisecond)
		close(gateHoldDone)
	}()

	// Wait for gate hold to complete (request is still blocked waiting for gate)
	<-gateHoldDone

	// Release the gate - now the request can proceed
	rm.requestGate <- token

	// Wait for SendPacket to be called (deterministic, no arbitrary sleep)
	if !signalClient.waitForSend(200 * time.Millisecond) {
		t.Fatal("SendPacket was not called within 200ms after gate release")
	}

	// Immediately inject a valid response packet so it completes quickly
	// If timeout budget was consumed during gate wait, this would timeout
	status := &protocol_common.Status{
		ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar,
		StatusWord:   0x0001,
		StateVar:     0x0002,
	}
	statusResp := protocol_common.NewStatusResponse(status)
	packet, _ := statusResp.WritePacket()
	server.SendPacket(packet)

	// Wait for request to complete
	wg.Wait()

	// Verify request succeeded (proves gate wait didn't consume timeout budget)
	if err != nil {
		if timeoutErr, ok := err.(*protocol_common.RequestTimeoutError); ok {
			t.Errorf("Request timed out with %v - gate wait consumed timeout budget (should not happen)", timeoutErr)
		} else {
			t.Errorf("Request failed with unexpected error: %v", err)
		}
	}
	if response == nil {
		t.Error("Request returned nil response")
	}
}

// trackingTransport wraps a transport.Client to track SendPacket calls
type trackingTransport struct {
	transport.Client
	sendPacketCalled atomic.Bool
}

func (tt *trackingTransport) SendPacket(ctx context.Context, data []byte) error {
	tt.sendPacketCalled.Store(true)
	return tt.Client.SendPacket(ctx, data)
}

func TestRequestManager_CancelBeforeFirstSend(t *testing.T) {
	// Create mock transport that tracks SendPacket calls
	baseClient, _ := transport.NewMockTransportClientWithServer()
	trackingClient := &trackingTransport{
		Client: baseClient,
	}

	// Create RequestManager
	rm := NewRequestManager(trackingClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Create a latch to pause startRequest before first send
	var latch sync.WaitGroup
	var releaseLatch sync.WaitGroup
	latch.Add(1)
	releaseLatch.Add(1)

	// Set test hook to pause before first send
	rm.beforeFirstSend = func(req *pendingRequest) {
		// Signal that we're paused
		latch.Done()
		// Wait for release
		releaseLatch.Wait()
	}

	// Create cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	req := protocol_common.NewStatusRequest()

	var wg sync.WaitGroup
	var err error

	// Submit request in background
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctx, req)
	}()

	// Wait for startRequest to pause (via hook)
	latch.Wait()

	// Verify SendPacket was NOT called yet
	if trackingClient.sendPacketCalled.Load() {
		t.Error("SendPacket was called before hook released")
	}

	// Cancel context while paused
	cancel()

	// Give cancellation time to propagate
	time.Sleep(10 * time.Millisecond)

	// Verify SendPacket was still NOT called
	if trackingClient.sendPacketCalled.Load() {
		t.Error("SendPacket was called before context cancellation was checked")
	}

	// Release latch - startRequest should now check context and abort
	releaseLatch.Done()

	// Wait for request to complete
	wg.Wait()

	// Verify request returned context.Canceled
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}

	// Verify SendPacket was NEVER called
	if trackingClient.sendPacketCalled.Load() {
		t.Error("SendPacket was called even though context was canceled before first send")
	}

	// Verify no pending requests leaked
	rm.mu.RLock()
	pendingCount := len(rm.pendingRequests) + len(rm.pendingMCRequests)
	if rm.pendingStatusRequest != nil {
		pendingCount++
	}
	rm.mu.RUnlock()

	if pendingCount != 0 {
		t.Errorf("Expected 0 pending requests after cancellation, got %d", pendingCount)
	}

	// Verify currentRequest is nil
	rm.mu.RLock()
	currentReq := rm.currentRequest
	rm.mu.RUnlock()

	if currentReq != nil {
		t.Error("Expected currentRequest to be nil after cancellation, but it was not")
	}
}

// TestRequestManager_TimeoutStartsAfterDeadlineStamped_NotAfterSendCompletes verifies
// that timeout starts when deadline is stamped (request becomes "live"), NOT when
// SendPacket completes. Even if SendPacket is slow/blocked, timeout should still run.
func TestRequestManager_TimeoutStartsAfterDeadlineStamped_NotAfterSendCompletes(t *testing.T) {
	// Create mock transport
	baseClient, _ := transport.NewMockTransportClientWithServer()

	// Create a context that we can cancel to unblock SendPacket
	blockCtx, blockCancel := context.WithCancel(context.Background())
	defer blockCancel()

	// Wrap transport to block SendPacket until blockCtx is canceled
	blockingClient := newBlockingSendTransport(baseClient, blockCtx)

	// Create RequestManager with very short timeout (20ms) and MaxRetries=1
	rm := NewRequestManagerWithConfig(blockingClient, 10*time.Millisecond, RequestManagerConfig{
		MaxRetries:     1,
		DefaultTimeout: 20 * time.Millisecond, // Short timeout
	})
	rm.Start()
	defer rm.Stop()

	// Start request with large context timeout (2s) so RM timeout triggers first
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req := protocol_common.NewStatusRequest()

	var wg sync.WaitGroup
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctx, req)
	}()

	// Wait for request to reach SendPacket (which will block).
	// After deadline is stamped in startRequest(), sentCh is closed, so roundTripBase starts timeout timer.
	// We wait > 20ms timeout to ensure timeout triggers even though SendPacket is blocked.
	// This proves timeout starts when deadline is stamped, NOT when SendPacket completes.
	time.Sleep(25 * time.Millisecond) // Wait > 20ms timeout

	// Unblock SendPacket (but request should have already timed out)
	blockCancel()

	// Wait for request to complete
	wg.Wait()

	// Assert: request should have timed out with RequestTimeoutError
	// This proves timeout started when deadline was stamped, NOT when SendPacket completed
	if err == nil {
		t.Error("Expected RequestTimeoutError, got nil error (timeout didn't trigger)")
		return
	}

	timeoutErr, ok := err.(*protocol_common.RequestTimeoutError)
	if !ok {
		t.Errorf("Expected RequestTimeoutError, got %T: %v", err, err)
		return
	}

	// Verify it's a timeout error (proves timeout ran even though SendPacket was blocked)
	if timeoutErr.Timeout <= 0 {
		t.Errorf("Expected positive timeout duration, got %v", timeoutErr.Timeout)
	}
}

func TestRequestManager_CancelBeforeSend_NoCounterConsumption(t *testing.T) {
	// Create mock transport
	baseClient, _ := transport.NewMockTransportClientWithServer()
	trackingClient := &trackingTransport{
		Client: baseClient,
	}

	// Create RequestManager
	rm := NewRequestManager(trackingClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// We'll verify counter wasn't consumed by checking that no cmdCount was allocated

	// Create latch to pause before first send
	var latch sync.WaitGroup
	var releaseLatch sync.WaitGroup
	latch.Add(1)
	releaseLatch.Add(1)

	// Set test hook to pause before first send
	rm.beforeFirstSend = func(req *pendingRequest) {
		// Signal that we're paused
		latch.Done()
		// Wait for release
		releaseLatch.Wait()
	}

	// Create cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	req := protocol_common.NewStatusRequest()

	var wg sync.WaitGroup
	var err error

	// Submit request in background
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctx, req)
	}()

	// Wait for startRequest to pause (via hook)
	latch.Wait()

	// Verify SendPacket was NOT called yet
	if trackingClient.sendPacketCalled.Load() {
		t.Error("SendPacket was called before hook released")
	}

	// Cancel context while paused
	cancel()

	// Give cancellation time to propagate
	time.Sleep(10 * time.Millisecond)

	// Verify SendPacket was still NOT called
	if trackingClient.sendPacketCalled.Load() {
		t.Error("SendPacket was called before context cancellation was checked")
	}

	// Release latch - startRequest should now check context and abort
	releaseLatch.Done()

	// Wait for request to complete
	wg.Wait()

	// Verify request returned context.Canceled
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}

	// Verify SendPacket was NEVER called
	if trackingClient.sendPacketCalled.Load() {
		t.Error("SendPacket was called even though context was canceled before first send")
	}

	// Verify RTC counter was NOT consumed by checking that request was canceled
	// before cmdCount allocation. Since the request was canceled in startRequest
	// before allocating cmdCount, and we verified SendPacket was never called,
	// we know the counter wasn't consumed.

	// Verify no pending requests leaked
	rm.mu.RLock()
	pendingCount := len(rm.pendingRequests) + len(rm.pendingMCRequests)
	if rm.pendingStatusRequest != nil {
		pendingCount++
	}
	rm.mu.RUnlock()

	if pendingCount != 0 {
		t.Errorf("Expected 0 pending requests after cancellation, got %d", pendingCount)
	}

	// Verify currentRequest is nil
	rm.mu.RLock()
	currentReq := rm.currentRequest
	rm.mu.RUnlock()

	if currentReq != nil {
		t.Error("Expected currentRequest to be nil after cancellation, but it was not")
	}

	// Verify sentCh was closed (no deadlock)
	// This is implicit - if sentCh wasn't closed, roundTripBase would deadlock
	// and the test would hang. Since we got here, sentCh was closed.
}

func TestRequestManager_CompletionIdempotency(t *testing.T) {
	// Create mock transport
	baseClient, server := transport.NewMockTransportClientWithServer()

	// Create RequestManager
	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Create request with context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := protocol_common.NewStatusRequest()

	// Submit request and send response immediately
	status := &protocol_common.Status{
		StatusWord: 0x0001,
		StateVar:   0x0002,
	}
	statusResp := protocol_common.NewStatusResponse(status)
	packet, _ := statusResp.WritePacket()

	// Start request in background
	var wg1 sync.WaitGroup
	var response1 *protocol_common.StatusResponse
	var err1 error

	wg1.Add(1)
	go func() {
		defer wg1.Done()
		response1, err1 = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctx, req)
	}()

	// Give request time to send
	time.Sleep(50 * time.Millisecond)

	// Send response
	server.SendPacket(packet)

	// Wait for completion
	wg1.Wait()

	// Verify first completion succeeded
	if err1 != nil {
		t.Fatalf("First completion failed: %v", err1)
	}
	if response1 == nil {
		t.Fatal("First completion returned nil response")
	}

	// Try to complete again with error (should be ignored due to idempotency)
	// We need to get the pending request to test this, but it's already completed.
	// Instead, test that late responses don't cause issues.

	// Submit another request
	req2 := protocol_common.NewStatusRequest()

	// Send response for second request
	status2 := &protocol_common.Status{
		StatusWord: 0x0001,
		StateVar:   0x0002,
	}
	statusResp2 := protocol_common.NewStatusResponse(status2)
	packet2, _ := statusResp2.WritePacket()

	// Start second request in background
	var wg2 sync.WaitGroup
	var response2 *protocol_common.StatusResponse
	var err2 error

	wg2.Add(1)
	go func() {
		defer wg2.Done()
		response2, err2 = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctx, req2)
	}()

	// Give request time to send
	time.Sleep(50 * time.Millisecond)

	// Send response
	server.SendPacket(packet2)

	// Wait for completion
	wg2.Wait()

	// Verify second request also succeeds
	if err2 != nil {
		t.Fatalf("Second request failed: %v", err2)
	}
	if response2 == nil {
		t.Fatal("Second request returned nil response")
	}

	// Test timeout idempotency: submit a request, let it timeout, then try to send response
	// The response should be ignored (idempotent completion already fired)
	req3 := protocol_common.NewStatusRequest()

	// Use very short timeout
	ctxShort, cancelShort := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelShort()

	// Submit request with short timeout
	_, err3 := SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctxShort, req3)

	// Should timeout
	if err3 == nil {
		t.Error("Expected timeout error, got nil")
	}

	// Now try to send a late response (simulate race condition)
	// This should not panic or deadlock
	lateStatus := &protocol_common.Status{
		StatusWord: 0x0001,
		StateVar:   0x0002,
	}
	lateStatusResp := protocol_common.NewStatusResponse(lateStatus)
	latePacket, _ := lateStatusResp.WritePacket()

	// Send response after timeout - should be safely ignored
	server.SendPacket(latePacket)

	// Give it time to process
	time.Sleep(50 * time.Millisecond)

	// Verify no panic occurred and system is still functional
	req4 := protocol_common.NewStatusRequest()

	// Send response for 4th request
	status4 := &protocol_common.Status{
		StatusWord: 0x0001,
		StateVar:   0x0002,
	}
	statusResp4 := protocol_common.NewStatusResponse(status4)
	packet4, _ := statusResp4.WritePacket()

	// Start 4th request in background
	var wg4 sync.WaitGroup
	var response4 *protocol_common.StatusResponse
	var err4 error

	wg4.Add(1)
	go func() {
		defer wg4.Done()
		response4, err4 = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctx, req4)
	}()

	// Give request time to send
	time.Sleep(50 * time.Millisecond)

	// Send response
	server.SendPacket(packet4)

	// Wait for completion
	wg4.Wait()

	if err4 != nil {
		t.Errorf("System should still be functional after late response, got error: %v", err4)
	}
	if response4 == nil {
		t.Error("System should still be functional after late response")
	}

	t.Log("Completion idempotency verified: late responses after timeout are safely ignored")
}

// constructStatusTelegramWithRTCReply constructs a packet that simulates a status telegram
// that includes RTC reply data. The packet has:
// - reqBits without RTC bit (bit 2 = 0, simulating a status telegram)
// - repBits with RTC reply bit (bit 8 = 1)
// - Status data blocks in response-bit order
// - RTC reply block at the response-bit offset
func constructStatusTelegramWithRTCReply(reqBits, repBits uint32, status *protocol_common.Status, rtcCounter uint8, upid uint16, value int32, cmdCode uint8) []byte {
	size := protocol_common.PacketHeaderSize + protocol_common.CalculateExpectedDataSize(repBits)
	packet := make([]byte, size)

	// Header (8 bytes)
	binary.LittleEndian.PutUint32(packet[0:4], reqBits) // Request flags (no bit 2)
	binary.LittleEndian.PutUint32(packet[4:8], repBits) // Response flags (with bit 8)

	offset := protocol_common.PacketHeaderSize
	if repBits&protocol_common.RespBitStatusWord != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], status.StatusWord)
		offset += 2
	}
	if repBits&protocol_common.RespBitStateVar != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], status.StateVar)
		offset += 2
	}
	if repBits&protocol_common.RespBitActualPosition != 0 {
		binary.LittleEndian.PutUint32(packet[offset:offset+4], uint32(status.ActualPosition))
		offset += 4
	}
	if repBits&protocol_common.RespBitDemandPosition != 0 {
		binary.LittleEndian.PutUint32(packet[offset:offset+4], uint32(status.DemandPosition))
		offset += 4
	}
	if repBits&protocol_common.RespBitCurrent != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], uint16(status.Current))
		offset += 2
	}
	if repBits&protocol_common.RespBitWarnWord != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], status.WarnWord)
		offset += 2
	}
	if repBits&protocol_common.RespBitErrorCode != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], status.ErrorCode)
		offset += 2
	}
	if repBits&protocol_common.RespBitMonitoringChannel != 0 {
		for i := 0; i < 4; i++ {
			binary.LittleEndian.PutUint32(packet[offset:offset+4], uint32(status.MonitoringChannel[i]))
			offset += 4
		}
	}

	rtcOffset, err := protocol_common.ResponseBlockOffset(repBits, protocol_common.RespBitRTCReplyData)
	if err != nil {
		rtcOffset = len(packet) - protocol_common.BlockSizes[protocol_common.RespBitRTCReplyData]
	}
	// Byte 0: low nibble = Command Count
	packet[rtcOffset+protocol_rtc.RTCDataOffsetCounter] = rtcCounter & protocol_rtc.RTCCounterMask
	// Byte 1: cmdCode (for standard parameter access)
	packet[rtcOffset+protocol_rtc.RTCDataOffsetCmdOrStatus] = cmdCode
	// Bytes 2-3: UPID
	binary.LittleEndian.PutUint16(packet[rtcOffset+protocol_rtc.RTCDataOffsetUPID:rtcOffset+protocol_rtc.RTCDataOffsetUPID+2], upid)
	// Bytes 4-7: Value
	binary.LittleEndian.PutUint32(packet[rtcOffset+protocol_rtc.RTCDataOffsetValue:rtcOffset+protocol_rtc.RTCDataOffsetValue+4], uint32(value))

	return packet
}

func constructMonitoringStatusTelegramWithRTCReply(reqBits, repBits uint32, status *protocol_common.Status, rtcCounter uint8, upid uint16, value int32, cmdCode uint8) []byte {
	size := protocol_common.PacketHeaderSize
	statusBits := []uint32{
		protocol_common.RespBitStatusWord,
		protocol_common.RespBitStateVar,
		protocol_common.RespBitActualPosition,
		protocol_common.RespBitDemandPosition,
		protocol_common.RespBitCurrent,
		protocol_common.RespBitWarnWord,
		protocol_common.RespBitErrorCode,
		protocol_common.RespBitMonitoringChannel,
	}
	for _, bit := range statusBits {
		if repBits&bit != 0 {
			size += protocol_common.BlockSizes[bit]
		}
	}
	rtcOffset := size
	size += protocol_common.BlockSizes[protocol_common.RespBitRTCReplyData]

	packet := make([]byte, size)
	binary.LittleEndian.PutUint32(packet[0:4], reqBits)
	binary.LittleEndian.PutUint32(packet[4:8], repBits)

	offset := protocol_common.PacketHeaderSize
	if repBits&protocol_common.RespBitStatusWord != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], status.StatusWord)
		offset += 2
	}
	if repBits&protocol_common.RespBitStateVar != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], status.StateVar)
		offset += 2
	}
	if repBits&protocol_common.RespBitActualPosition != 0 {
		binary.LittleEndian.PutUint32(packet[offset:offset+4], uint32(status.ActualPosition))
		offset += 4
	}
	if repBits&protocol_common.RespBitDemandPosition != 0 {
		binary.LittleEndian.PutUint32(packet[offset:offset+4], uint32(status.DemandPosition))
		offset += 4
	}
	if repBits&protocol_common.RespBitCurrent != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], uint16(status.Current))
		offset += 2
	}
	if repBits&protocol_common.RespBitWarnWord != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], status.WarnWord)
		offset += 2
	}
	if repBits&protocol_common.RespBitErrorCode != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], status.ErrorCode)
		offset += 2
	}
	if repBits&protocol_common.RespBitMonitoringChannel != 0 {
		for i := 0; i < 4; i++ {
			binary.LittleEndian.PutUint32(packet[offset:offset+4], uint32(status.MonitoringChannel[i]))
			offset += 4
		}
	}

	packet[rtcOffset+protocol_rtc.RTCDataOffsetCounter] = rtcCounter & protocol_rtc.RTCCounterMask
	packet[rtcOffset+protocol_rtc.RTCDataOffsetCmdOrStatus] = cmdCode
	binary.LittleEndian.PutUint16(packet[rtcOffset+protocol_rtc.RTCDataOffsetUPID:rtcOffset+protocol_rtc.RTCDataOffsetUPID+2], upid)
	binary.LittleEndian.PutUint32(packet[rtcOffset+protocol_rtc.RTCDataOffsetValue:rtcOffset+protocol_rtc.RTCDataOffsetValue+4], uint32(value))

	return packet
}

// TestRequestManager_RTCReplyDetectedFromStatusTelegram verifies that RTC replies are detected
// from status telegrams that include RTC reply data (bit 8 set) even when the request echo
// doesn't have bit 2 set. This ensures we don't miss RTC replies embedded in status telegrams.
func TestRequestManager_RTCReplyDetectedFromStatusTelegram(t *testing.T) {
	// Create mock transport with server
	baseClient, server := transport.NewMockTransportClientWithServer()

	// Create RequestManager
	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Register a simple response registry for the test UPID
	// This is needed because ReadRTCResponse requires a registry to create typed responses
	testUPID := protocol_common.ParameterID(0x1234)
	protocol_rtc.RegisterResponseRegistry(testUPID, func(
		status *protocol_common.Status,
		value int32,
		upid uint16,
		rtcCounter uint8,
		rtcStatus uint8,
		cmdCode uint8,
	) protocol_common.Response {
		return protocol_rtc.NewRTCGetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
	})

	// Submit an RTC request to get a pending request with cmdCount
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upid := uint16(0x1234)
	cmdCode := uint8(0x10) // ReadRAM
	req := protocol_rtc.NewRTCGetParamRequest(upid, cmdCode)

	// Start request in background
	var wg sync.WaitGroup
	var response *protocol_rtc.RTCGetParamResponse
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](rm, ctx, req)
	}()

	// Wait a moment for the request to be submitted and cmdCount allocated
	time.Sleep(50 * time.Millisecond)

	// Get the cmdCount that was allocated (peek at pending requests)
	rm.mu.RLock()
	var cmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			cmdCount = count
			break
		}
	}
	rm.mu.RUnlock()

	if cmdCount == 0 {
		t.Fatal("Failed to find pending RTC request with cmdCount")
	}

	// Construct a status telegram with RTC reply data:
	// - reqBits without bit 2 (simulating status telegram, not explicit RTC response)
	// - repBits with bit 8 (RTC reply data present)
	reqBits := uint32(0x00000000) // No RTC command bit (bit 2 = 0)
	repBits := protocol_common.RespBitStatusWord |
		protocol_common.RespBitStateVar |
		protocol_common.RespBitActualPosition |
		protocol_common.RespBitDemandPosition |
		protocol_common.RespBitCurrent |
		protocol_common.RespBitWarnWord |
		protocol_common.RespBitErrorCode |
		protocol_common.RespBitRTCReplyData // Bit 8 set

	status := &protocol_common.Status{
		ResponseBits:   repBits,
		StatusWord:     0x0000,
		StateVar:       0x0000,
		ActualPosition: 0,
		DemandPosition: 0,
		Current:        0,
		WarnWord:       0x0000,
		ErrorCode:      0x0000,
	}

	value := int32(0x5678) // Test value

	// Construct the packet
	packet := constructStatusTelegramWithRTCReply(reqBits, repBits, status, cmdCount, upid, value, cmdCode)

	// Send the packet through the mock server
	server.SendPacket(packet)

	// Wait for request to complete
	wg.Wait()

	// Verify the request completed successfully
	if err != nil {
		t.Fatalf("Expected RTC request to complete successfully, got error: %v", err)
	}

	if response == nil {
		t.Fatal("Expected RTC response, got nil")
	}

	// Verify it's an RTC response (not StatusResponse)
	if response.Value() != value {
		t.Errorf("Expected value %d, got %d", value, response.Value())
	}

	if response.RTCCounter() != cmdCount {
		t.Errorf("Expected cmdCount %d, got %d", cmdCount, response.RTCCounter())
	}

	// Verify RTCStatus is 0x00 (success)
	if response.RTCStatus() != 0x00 {
		t.Errorf("Expected RTCStatus 0x00, got 0x%02X", response.RTCStatus())
	}

	t.Logf("RTC reply successfully detected from status telegram: cmdCount=%d, UPID=0x%04X, value=%d", cmdCount, upid, value)
}

// TestRequestManager_RTCDesyncResync_AlwaysOn verifies that RTC desync/resync behavior
// is always-on (production behavior), not gated by debug flag.
// When an unmatched RTC reply is received while an RTC request is in-flight,
// the request should be re-keyed with a new cmdCount and resent, even when debug is disabled.
func TestRequestManager_RTCDesyncResync_AlwaysOn(t *testing.T) {
	// Create mock transport with server
	baseClient, server := transport.NewMockTransportClientWithServer()

	// Create RequestManager with debug DISABLED
	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.SetDebug(false) // Explicitly disable debug
	rm.Start()
	defer rm.Stop()

	// Register a simple response registry for the test UPID
	testUPID := protocol_common.ParameterID(0x1234)
	protocol_rtc.RegisterResponseRegistry(testUPID, func(
		status *protocol_common.Status,
		value int32,
		upid uint16,
		rtcCounter uint8,
		rtcStatus uint8,
		cmdCode uint8,
	) protocol_common.Response {
		return protocol_rtc.NewRTCGetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
	})

	// Submit an RTC request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upid := uint16(0x1234)
	cmdCode := uint8(0x10) // ReadRAM
	req := protocol_rtc.NewRTCGetParamRequest(upid, cmdCode)

	// Start request in background
	var wg sync.WaitGroup
	var response *protocol_rtc.RTCGetParamResponse
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](rm, ctx, req)
	}()

	// Wait for request to be submitted and cmdCount allocated
	time.Sleep(50 * time.Millisecond)

	// Get the cmdCount that was allocated
	rm.mu.RLock()
	var clientCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			clientCmdCount = count
			break
		}
	}
	rm.mu.RUnlock()

	if clientCmdCount == 0 {
		t.Fatal("Failed to find pending RTC request with cmdCount")
	}

	// Inject an UNMATCHED RTC reply (different cmdCount than what client sent)
	// This simulates a desync scenario where the drive echoes a different cmdCount
	// Use delta=1 (plausible) to trigger resync
	unmatchedCmdCount := clientCmdCount + 1
	if unmatchedCmdCount > 14 {
		unmatchedCmdCount = 1
	}

	// Construct an unmatched RTC reply packet
	reqBits := uint32(0x00000000)
	repBits := protocol_common.RespBitStatusWord |
		protocol_common.RespBitStateVar |
		protocol_common.RespBitRTCReplyData // Bit 8 set

	status := &protocol_common.Status{
		ResponseBits:   repBits,
		StatusWord:     0x0000,
		StateVar:       0x0000,
		ActualPosition: 0,
		DemandPosition: 0,
		Current:        0,
		WarnWord:       0x0000,
		ErrorCode:      0x0000,
	}

	value := int32(0x5678)

	// Send unmatched reply (different cmdCount)
	unmatchedPacket := constructStatusTelegramWithRTCReply(reqBits, repBits, status, unmatchedCmdCount, upid, value, cmdCode)
	server.SendPacket(unmatchedPacket)

	// Wait a moment for resync to occur
	time.Sleep(100 * time.Millisecond)

	// Verify resync occurred: check that request was re-keyed with a new cmdCount
	rm.mu.RLock()
	var newCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			newCmdCount = count
			break
		}
	}
	resyncAttempted := rm.resyncAttempted
	rm.mu.RUnlock()

	if newCmdCount == 0 {
		t.Fatal("Request should still be pending after resync")
	}

	if newCmdCount == clientCmdCount {
		t.Errorf("Expected request to be re-keyed with new cmdCount, but still using old cmdCount %d", clientCmdCount)
	}

	if !resyncAttempted {
		t.Error("Expected resyncAttempted flag to be set, but it was false")
	}

	// Now send a matching reply with the NEW cmdCount
	matchingPacket := constructStatusTelegramWithRTCReply(reqBits, repBits, status, newCmdCount, upid, value, cmdCode)
	server.SendPacket(matchingPacket)

	// Wait for request to complete
	wg.Wait()

	// Verify the request completed successfully after resync
	if err != nil {
		t.Fatalf("Expected RTC request to complete successfully after resync, got error: %v", err)
	}

	if response == nil {
		t.Fatal("Expected RTC response, got nil")
	}

	if response.Value() != value {
		t.Errorf("Expected value %d, got %d", value, response.Value())
	}

	if response.RTCCounter() != newCmdCount {
		t.Errorf("Expected cmdCount %d, got %d", newCmdCount, response.RTCCounter())
	}

	t.Logf("RTC desync/resync worked correctly: originalCmdCount=%d, unmatchedCmdCount=%d, newCmdCount=%d (debug disabled)", clientCmdCount, unmatchedCmdCount, newCmdCount)
}

// TestRequestManager_RTCDesyncResync_NoInFlightRequest verifies that resync does NOT occur
// when there is no in-flight RTC request. Unmatched RTC replies should be ignored in this case.
func TestRequestManager_RTCDesyncResync_NoInFlightRequest(t *testing.T) {
	// Create mock transport with server
	baseClient, server := transport.NewMockTransportClientWithServer()

	// Create RequestManager
	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Register a simple response registry for the test UPID
	testUPID := protocol_common.ParameterID(0x1234)
	protocol_rtc.RegisterResponseRegistry(testUPID, func(
		status *protocol_common.Status,
		value int32,
		upid uint16,
		rtcCounter uint8,
		rtcStatus uint8,
		cmdCode uint8,
	) protocol_common.Response {
		return protocol_rtc.NewRTCGetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
	})

	// Verify no in-flight RTC request
	rm.mu.RLock()
	hasInFlightRTC := rm.currentRequest != nil && !rm.currentRequest.isStatusRequest && !rm.currentRequest.isMCRequest
	pendingCount := len(rm.pendingRequests)
	rm.mu.RUnlock()

	if hasInFlightRTC {
		t.Fatal("Expected no in-flight RTC request, but found one")
	}

	// Inject an unmatched RTC reply (cmdCount = 5, but no pending request)
	unmatchedCmdCount := uint8(5)
	reqBits := uint32(0x00000000)
	repBits := protocol_common.RespBitStatusWord |
		protocol_common.RespBitStateVar |
		protocol_common.RespBitRTCReplyData

	status := &protocol_common.Status{
		ResponseBits:   repBits,
		StatusWord:     0x0000,
		StateVar:       0x0000,
		ActualPosition: 0,
		DemandPosition: 0,
		Current:        0,
		WarnWord:       0x0000,
		ErrorCode:      0x0000,
	}

	value := int32(0x5678)
	upid := uint16(0x1234)
	cmdCode := uint8(0x10)

	unmatchedPacket := constructStatusTelegramWithRTCReply(reqBits, repBits, status, unmatchedCmdCount, upid, value, cmdCode)
	server.SendPacket(unmatchedPacket)

	// Wait a moment
	time.Sleep(50 * time.Millisecond)

	// Verify resync did NOT occur (no in-flight request, so no resync)
	rm.mu.RLock()
	resyncAttempted := rm.resyncAttempted
	pendingCountAfter := len(rm.pendingRequests)
	rm.mu.RUnlock()

	if resyncAttempted {
		t.Error("Expected resyncAttempted to remain false (no in-flight request), but it was true")
	}

	if pendingCountAfter != pendingCount {
		t.Errorf("Expected pending request count to remain %d, got %d", pendingCount, pendingCountAfter)
	}

	t.Logf("Correctly ignored unmatched RTC reply when no in-flight request: cmdCount=%d, resyncAttempted=%v", unmatchedCmdCount, resyncAttempted)
}

// TestRequestManager_MonitoringOnlyPacketDoesNotResync verifies that monitoring-only packets
// with RTC reply bits set but empty RTC payload do not trigger resync or consume pending RTC.
func TestRequestManager_MonitoringOnlyPacketDoesNotResync(t *testing.T) {
	baseClient, server := transport.NewMockTransportClientWithServer()

	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	testUPID := protocol_common.ParameterID(0x1234)
	protocol_rtc.RegisterResponseRegistry(testUPID, func(
		status *protocol_common.Status,
		value int32,
		upid uint16,
		rtcCounter uint8,
		rtcStatus uint8,
		cmdCode uint8,
	) protocol_common.Response {
		return protocol_rtc.NewRTCGetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upid := uint16(0x1234)
	cmdCode := uint8(0x10) // ReadRAM
	req := protocol_rtc.NewRTCGetParamRequest(upid, cmdCode)

	var wg sync.WaitGroup
	var response *protocol_rtc.RTCGetParamResponse
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](rm, ctx, req)
	}()

	time.Sleep(50 * time.Millisecond)

	rm.mu.RLock()
	var clientCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			clientCmdCount = count
			break
		}
	}
	rm.mu.RUnlock()

	if clientCmdCount == 0 {
		t.Fatal("Failed to find pending RTC request with cmdCount")
	}

	fakeCmdCount := nextRTCCounter(clientCmdCount) // mismatch that would trigger resync if meaningful

	reqBits := uint32(0x00000000)
	repBits := uint32(0x000001FF)
	status := &protocol_common.Status{
		ResponseBits: repBits,
	}

	monitoringPacket := constructMonitoringStatusTelegramWithRTCReply(reqBits, repBits, status, fakeCmdCount, upid, 0, 0)
	server.SendPacket(monitoringPacket)

	time.Sleep(100 * time.Millisecond)

	rm.mu.RLock()
	resyncAttempted := rm.resyncAttempted
	_, stillPending := rm.pendingRequests[clientCmdCount]
	rm.mu.RUnlock()

	if resyncAttempted {
		t.Fatal("Expected resyncAttempted to remain false for monitoring-only packet")
	}
	if !stillPending {
		t.Fatal("Expected RTC request to remain pending after monitoring-only packet")
	}

	value := int32(0x5678)
	matchingPacket := constructStatusTelegramWithRTCReply(reqBits, repBits, status, clientCmdCount, upid, value, cmdCode)
	server.SendPacket(matchingPacket)

	wg.Wait()

	if err != nil {
		t.Fatalf("Expected RTC request to complete successfully, got error: %v", err)
	}
	if response == nil {
		t.Fatal("Expected RTC response, got nil")
	}
	if response.RTCCounter() != clientCmdCount {
		t.Errorf("Expected cmdCount %d, got %d", clientCmdCount, response.RTCCounter())
	}
	if response.Value() != value {
		t.Errorf("Expected value %d, got %d", value, response.Value())
	}
}

func TestRequestManager_RTCSeed_FromMonitoringBeacon(t *testing.T) {
	baseClient, server := transport.NewMockTransportClientWithServer()

	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	testUPID := protocol_common.ParameterID(0x1234)
	protocol_rtc.RegisterResponseRegistry(testUPID, func(
		status *protocol_common.Status,
		value int32,
		upid uint16,
		rtcCounter uint8,
		rtcStatus uint8,
		cmdCode uint8,
	) protocol_common.Response {
		return protocol_rtc.NewRTCGetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
	})

	reqBits := uint32(0x00000000)
	repBits := uint32(0x000001FF)
	status := &protocol_common.Status{
		ResponseBits: repBits,
	}

	monitoringPacket := constructMonitoringStatusTelegramWithRTCReply(reqBits, repBits, status, 12, 0, 0, 0)
	server.SendPacket(monitoringPacket)

	time.Sleep(50 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upid := uint16(0x1234)
	cmdCode := uint8(0x10)
	req := protocol_rtc.NewRTCGetParamRequest(upid, cmdCode)

	var wg sync.WaitGroup
	var response *protocol_rtc.RTCGetParamResponse
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](rm, ctx, req)
	}()

	var seededCmdCount uint8
	for i := 0; i < 50; i++ {
		rm.mu.RLock()
		for count, pendingReq := range rm.pendingRequests {
			if pendingReq.request == req {
				seededCmdCount = count
				break
			}
		}
		rm.mu.RUnlock()
		if seededCmdCount != 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if seededCmdCount == 0 {
		t.Fatal("Failed to observe pending RTC request after seeding")
	}
	if seededCmdCount != 13 {
		t.Fatalf("Expected seeded cmdCount=13, got %d", seededCmdCount)
	}

	matchingPacket := constructStatusTelegramWithRTCReply(reqBits, repBits, status, seededCmdCount, upid, 0x5678, cmdCode)
	server.SendPacket(matchingPacket)

	wg.Wait()

	if err != nil {
		t.Fatalf("Expected RTC request to complete successfully, got error: %v", err)
	}
	if response == nil {
		t.Fatal("Expected RTC response, got nil")
	}
	if response.RTCCounter() != seededCmdCount {
		t.Errorf("Expected cmdCount %d, got %d", seededCmdCount, response.RTCCounter())
	}
}

func TestRequestManager_BeaconResync_AfterRetry(t *testing.T) {
	baseClient, server := transport.NewMockTransportClientWithServer()

	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	testUPID := protocol_common.ParameterID(0x1234)
	protocol_rtc.RegisterResponseRegistry(testUPID, func(
		status *protocol_common.Status,
		value int32,
		upid uint16,
		rtcCounter uint8,
		rtcStatus uint8,
		cmdCode uint8,
	) protocol_common.Response {
		return protocol_rtc.NewRTCGetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upid := uint16(0x1234)
	cmdCode := uint8(0x10)
	req := protocol_rtc.NewRTCGetParamRequest(upid, cmdCode)

	var wg sync.WaitGroup
	var response *protocol_rtc.RTCGetParamResponse
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](rm, ctx, req)
	}()

	var originalCmdCount uint8
	for i := 0; i < 50; i++ {
		rm.mu.RLock()
		for count, pendingReq := range rm.pendingRequests {
			if pendingReq.request == req {
				originalCmdCount = count
				break
			}
		}
		rm.mu.RUnlock()
		if originalCmdCount != 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if originalCmdCount == 0 {
		t.Fatal("Failed to observe pending RTC request")
	}

	rm.mu.Lock()
	if rm.currentRequest != nil {
		rm.currentRequest.sendCounter = 2
	}
	rm.mu.Unlock()

	reqBits := uint32(0x00000000)
	repBits := uint32(0x000001FF)
	status := &protocol_common.Status{
		ResponseBits: repBits,
	}

	beaconPacket := constructMonitoringStatusTelegramWithRTCReply(reqBits, repBits, status, 12, upid, 0, 0)
	server.SendPacket(beaconPacket)

	time.Sleep(100 * time.Millisecond)

	rm.mu.RLock()
	resyncAttempted := rm.resyncAttempted
	var newCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			newCmdCount = count
			break
		}
	}
	rm.mu.RUnlock()

	if !resyncAttempted {
		t.Fatal("Expected resyncAttempted to be true after beacon resync")
	}
	if newCmdCount != 13 {
		t.Fatalf("Expected beacon resync cmdCount=13, got %d (original=%d)", newCmdCount, originalCmdCount)
	}

	matchingPacket := constructStatusTelegramWithRTCReply(reqBits, repBits, status, newCmdCount, upid, 0x5678, cmdCode)
	server.SendPacket(matchingPacket)

	wg.Wait()

	if err != nil {
		t.Fatalf("Expected RTC request to complete successfully, got error: %v", err)
	}
	if response == nil {
		t.Fatal("Expected RTC response, got nil")
	}
	if response.RTCCounter() != newCmdCount {
		t.Errorf("Expected cmdCount %d, got %d", newCmdCount, response.RTCCounter())
	}
}

// TestRequestManager_RTCDesyncResync_OnlyOncePerRequest verifies that resync happens
// only once per request, even if multiple unmatched RTC replies are received.
// The resyncAttempted flag should prevent multiple resyncs for the same request.
func TestRequestManager_RTCDesyncResync_OnlyOncePerRequest(t *testing.T) {
	// Create mock transport with server
	baseClient, server := transport.NewMockTransportClientWithServer()

	// Create RequestManager
	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Register a simple response registry for the test UPID
	testUPID := protocol_common.ParameterID(0x1234)
	protocol_rtc.RegisterResponseRegistry(testUPID, func(
		status *protocol_common.Status,
		value int32,
		upid uint16,
		rtcCounter uint8,
		rtcStatus uint8,
		cmdCode uint8,
	) protocol_common.Response {
		return protocol_rtc.NewRTCGetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
	})

	// Submit an RTC request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upid := uint16(0x1234)
	cmdCode := uint8(0x10) // ReadRAM
	req := protocol_rtc.NewRTCGetParamRequest(upid, cmdCode)

	// Start request in background
	var wg sync.WaitGroup
	var response *protocol_rtc.RTCGetParamResponse
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](rm, ctx, req)
	}()

	// Wait for request to be submitted and cmdCount allocated
	time.Sleep(50 * time.Millisecond)

	// Get the cmdCount that was allocated
	rm.mu.RLock()
	var clientCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			clientCmdCount = count
			break
		}
	}
	rm.mu.RUnlock()

	if clientCmdCount == 0 {
		t.Fatal("Failed to find pending RTC request with cmdCount")
	}

	// Construct unmatched reply packets
	reqBits := uint32(0x00000000)
	repBits := protocol_common.RespBitStatusWord |
		protocol_common.RespBitStateVar |
		protocol_common.RespBitRTCReplyData

	status := &protocol_common.Status{
		ResponseBits:   repBits,
		StatusWord:     0x0000,
		StateVar:       0x0000,
		ActualPosition: 0,
		DemandPosition: 0,
		Current:        0,
		WarnWord:       0x0000,
		ErrorCode:      0x0000,
	}

	value := int32(0x5678)

	// Send FIRST unmatched reply (should trigger resync)
	// Use delta=1 (plausible) to trigger resync
	unmatchedCmdCount1 := clientCmdCount + 1
	if unmatchedCmdCount1 > 14 {
		unmatchedCmdCount1 = 1
	}
	unmatchedPacket1 := constructStatusTelegramWithRTCReply(reqBits, repBits, status, unmatchedCmdCount1, upid, value, cmdCode)
	server.SendPacket(unmatchedPacket1)

	// Wait for first resync
	time.Sleep(100 * time.Millisecond)

	// Get new cmdCount after first resync
	rm.mu.RLock()
	var newCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			newCmdCount = count
			break
		}
	}
	resyncAttempted1 := rm.resyncAttempted
	rm.mu.RUnlock()

	if newCmdCount == clientCmdCount {
		t.Fatal("Expected request to be re-keyed after first unmatched reply")
	}

	if !resyncAttempted1 {
		t.Fatal("Expected resyncAttempted to be true after first unmatched reply")
	}

	// Send SECOND unmatched reply (should NOT trigger another resync)
	// Use delta=1 (plausible) but resyncAttempted flag should prevent it
	unmatchedCmdCount2 := newCmdCount + 1
	if unmatchedCmdCount2 > 14 {
		unmatchedCmdCount2 = 1
	}
	unmatchedPacket2 := constructStatusTelegramWithRTCReply(reqBits, repBits, status, unmatchedCmdCount2, upid, value, cmdCode)
	server.SendPacket(unmatchedPacket2)

	// Wait a moment
	time.Sleep(50 * time.Millisecond)

	// Verify cmdCount did NOT change again (resync only happens once)
	rm.mu.RLock()
	var finalCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			finalCmdCount = count
			break
		}
	}
	resyncAttempted2 := rm.resyncAttempted
	rm.mu.RUnlock()

	if finalCmdCount != newCmdCount {
		t.Errorf("Expected cmdCount to remain %d after second unmatched reply, but changed to %d (resync should only happen once)", newCmdCount, finalCmdCount)
	}

	if !resyncAttempted2 {
		t.Error("Expected resyncAttempted to remain true after second unmatched reply")
	}

	// Now send a matching reply with the NEW cmdCount (from first resync)
	matchingPacket := constructStatusTelegramWithRTCReply(reqBits, repBits, status, newCmdCount, upid, value, cmdCode)
	server.SendPacket(matchingPacket)

	// Wait for request to complete
	wg.Wait()

	// Verify the request completed successfully
	if err != nil {
		t.Fatalf("Expected RTC request to complete successfully, got error: %v", err)
	}

	if response == nil {
		t.Fatal("Expected RTC response, got nil")
	}

	if response.RTCCounter() != newCmdCount {
		t.Errorf("Expected cmdCount %d, got %d", newCmdCount, response.RTCCounter())
	}

	t.Logf("Resync happened only once: originalCmdCount=%d, firstUnmatched=%d, newCmdCount=%d, secondUnmatched=%d, finalCmdCount=%d", clientCmdCount, unmatchedCmdCount1, newCmdCount, unmatchedCmdCount2, finalCmdCount)
}

// TestRequestManager_StatusResponse_NotSwallowedByRTCReplyData verifies that status responses
// are not swallowed when RespBitRTCReplyData is present but no matching RTC request exists.
// This ensures status requests can complete even when the drive includes RTC reply data in status telegrams.
func TestRequestManager_StatusResponse_NotSwallowedByRTCReplyData(t *testing.T) {
	// Create mock transport with server
	baseClient, server := transport.NewMockTransportClientWithServer()

	// Create RequestManager
	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Submit a status request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := protocol_common.NewStatusRequest()

	// Start request in background
	var wg sync.WaitGroup
	var response *protocol_common.StatusResponse
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctx, req)
	}()

	// Wait for request to be submitted
	time.Sleep(50 * time.Millisecond)

	// Verify no pending RTC requests exist
	rm.mu.RLock()
	pendingRTCCount := len(rm.pendingRequests)
	rm.mu.RUnlock()

	if pendingRTCCount != 0 {
		t.Fatalf("Expected no pending RTC requests, but found %d", pendingRTCCount)
	}

	// Construct a status telegram with RTC reply data but NO matching RTC request:
	// - Status bits set (StatusWord, StateVar, etc.)
	// - RespBitRTCReplyData set (bit 8)
	// - reqBits without RTC bit (bit 2 = 0)
	// - RTC reply data block with cmdCount that doesn't match any pending RTC request
	reqBits := uint32(0x00000000) // No RTC command bit (bit 2 = 0)
	repBits := protocol_common.RespBitStatusWord |
		protocol_common.RespBitStateVar |
		protocol_common.RespBitActualPosition |
		protocol_common.RespBitDemandPosition |
		protocol_common.RespBitCurrent |
		protocol_common.RespBitWarnWord |
		protocol_common.RespBitErrorCode |
		protocol_common.RespBitRTCReplyData // Bit 8 set

	status := &protocol_common.Status{
		ResponseBits:   repBits,
		StatusWord:     0x1234,
		StateVar:       0x5678,
		ActualPosition: 1000000,
		DemandPosition: 2000000,
		Current:        5000,
		WarnWord:       0x0000,
		ErrorCode:      0x0000,
	}

	// Use a cmdCount that doesn't match any pending RTC request
	unmatchedCmdCount := uint8(5)
	upid := uint16(0x1234)
	cmdCode := uint8(0x10)
	value := int32(0x5678)

	// Construct and send the packet
	packet := constructStatusTelegramWithRTCReply(reqBits, repBits, status, unmatchedCmdCount, upid, value, cmdCode)
	server.SendPacket(packet)

	// Wait for request to complete
	wg.Wait()

	// Verify the status request completed successfully (not swallowed)
	if err != nil {
		t.Fatalf("Expected status request to complete successfully, got error: %v", err)
	}

	if response == nil {
		t.Fatal("Expected StatusResponse, got nil")
	}

	// Verify it's a StatusResponse (not RTC response)
	if response.Status() == nil {
		t.Fatal("Expected StatusResponse with status data, got nil status")
	}

	// Verify status data is correct
	if response.Status().StatusWord != status.StatusWord {
		t.Errorf("Expected StatusWord 0x%04X, got 0x%04X", status.StatusWord, response.Status().StatusWord)
	}

	if response.Status().StateVar != status.StateVar {
		t.Errorf("Expected StateVar 0x%04X, got 0x%04X", status.StateVar, response.Status().StateVar)
	}

	// Verify lastStateVarLow was updated (MC counter calculation path)
	rm.mu.RLock()
	lastStateVarLowNibble := rm.lastStateVarLowNibble
	rm.mu.RUnlock()

	expectedNibble := uint8(status.StateVar & 0x0F)
	if lastStateVarLowNibble != expectedNibble {
		t.Errorf("Expected lastStateVarLowNibble 0x%X, got 0x%X", expectedNibble, lastStateVarLowNibble)
	}

	t.Logf("Status response correctly delivered despite RTC reply data: StatusWord=0x%04X, StateVar=0x%04X", status.StatusWord, status.StateVar)
}

// TestRequestManager_MonitoringAndRTCReply_CompletesBoth ensures a packet with
// monitoring data + RTC reply data satisfies both a monitoring request and an
// in-flight RTC request.
func TestRequestManager_MonitoringAndRTCReply_CompletesBoth(t *testing.T) {
	baseClient, server := transport.NewMockTransportClientWithServer()

	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	testUPID := protocol_common.ParameterID(0x4321)
	protocol_rtc.RegisterResponseRegistry(testUPID, func(
		status *protocol_common.Status,
		value int32,
		upid uint16,
		rtcCounter uint8,
		rtcStatus uint8,
		cmdCode uint8,
	) protocol_common.Response {
		return protocol_rtc.NewRTCGetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
	})

	entered := make(chan struct{})
	release := make(chan struct{})
	var rtcPending *pendingRequest
	var once sync.Once
	rm.beforeFirstSend = func(req *pendingRequest) {
		if req.kind() != "rtc" {
			return
		}
		once.Do(func() {
			rtcPending = req
			close(entered)
			<-release
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	monitorCtx, monitorCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer monitorCancel()
	monitorReq := protocol_common.NewMonitoringStatusRequest()
	perAttemptTimeout := monitorReq.OperationTimeout()
	if perAttemptTimeout <= 0 {
		perAttemptTimeout = rm.config.DefaultTimeout
	}
	totalTimeout := perAttemptTimeout * time.Duration(rm.config.MaxRetries)
	monitorPending, err := rm.submitRequest(monitorCtx, monitorReq, totalTimeout, perAttemptTimeout)
	if err != nil {
		t.Fatalf("submit monitoring request failed: %v", err)
	}
	<-monitorPending.sentCh

	rtcReq := protocol_rtc.NewRTCGetParamRequest(uint16(testUPID), protocol_rtc.CommandCode.ReadRAM)
	var rtcResp *protocol_rtc.RTCGetParamResponse
	var rtcErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		rtcResp, rtcErr = SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](rm, ctx, rtcReq)
	}()

	<-entered

	close(release)
	if rtcPending == nil {
		t.Fatal("expected rtcPending to be set")
	}
	<-rtcPending.sentCh

	reqBits := uint32(0x00000000)
	repBits := protocol_common.ResponseFlags.StandardWithMonitoring | protocol_common.RespBitRTCReplyData
	status := &protocol_common.Status{
		ResponseBits:   repBits,
		StatusWord:     0x1111,
		StateVar:       0x2222,
		ActualPosition: 1000,
		DemandPosition: 2000,
		Current:        300,
		WarnWord:       0x0000,
		ErrorCode:      0x0000,
		MonitoringChannel: [4]int32{
			10, 20, 30, 40,
		},
	}
	value := int32(0x5678)
	cmdCode := protocol_rtc.CommandCode.ReadRAM
	packet := constructMonitoringStatusTelegramWithRTCReply(reqBits, repBits, status, rtcPending.commandCount, uint16(testUPID), value, cmdCode)
	server.SendPacket(packet)

	var monitorResp protocol_common.Response
	select {
	case monitorResp = <-monitorPending.responseChannel:
	case err := <-monitorPending.errorChannel:
		t.Fatalf("monitoring request failed: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for monitoring response")
	}

	wg.Wait()
	if rtcErr != nil {
		t.Fatalf("RTC request failed: %v", rtcErr)
	}
	if rtcResp == nil {
		t.Fatal("expected RTC response, got nil")
	}

	monitoring, ok := monitorResp.(*protocol_common.MonitoringStatusResponse)
	if !ok || monitoring.Status() == nil {
		t.Fatalf("expected MonitoringStatusResponse, got %T", monitorResp)
	}
	if monitoring.Status().MonitoringChannel[0] != status.MonitoringChannel[0] {
		t.Fatalf("unexpected monitoring channel value: %d", monitoring.Status().MonitoringChannel[0])
	}
	if rtcResp.RTCCounter() != rtcPending.commandCount {
		t.Fatalf("expected rtc counter %d, got %d", rtcPending.commandCount, rtcResp.RTCCounter())
	}
}

// TestRequestManager_StatusResponse_AcceptsMonitoringResponse verifies that a monitoring status
// response can satisfy a StatusRequest by downcasting to StatusResponse.
func TestRequestManager_StatusResponse_AcceptsMonitoringResponse(t *testing.T) {
	// Create mock transport with server
	baseClient, server := transport.NewMockTransportClientWithServer()

	// Create RequestManager
	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Submit a status request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := protocol_common.NewStatusRequest()

	// Start request in background
	var wg sync.WaitGroup
	var response *protocol_common.StatusResponse
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_common.StatusResponse](rm, ctx, req)
	}()

	// Give request time to send
	time.Sleep(50 * time.Millisecond)

	status := &protocol_common.Status{
		ResponseBits:   protocol_common.ResponseFlags.StandardWithMonitoring,
		StatusWord:     0x1234,
		StateVar:       0x5678,
		ActualPosition: 100,
		DemandPosition: 200,
		Current:        300,
		WarnWord:       0x0000,
		ErrorCode:      0x0000,
		MonitoringChannel: [4]int32{
			0x11111111,
			0x22222222,
			0x33333333,
			0x44444444,
		},
	}
	monitoringResp := protocol_common.NewMonitoringStatusResponse(status)
	packet, err := monitoringResp.WritePacket()
	if err != nil {
		t.Fatalf("failed to write monitoring response packet: %v", err)
	}

	server.SendPacket(packet)

	// Wait for request to complete
	wg.Wait()

	if err != nil {
		t.Fatalf("Expected status request to complete successfully, got error: %v", err)
	}
	if response == nil {
		t.Fatal("Expected StatusResponse, got nil")
	}
	if response.Status() == nil {
		t.Fatal("Expected StatusResponse with status data, got nil status")
	}
	if response.Status().StatusWord != status.StatusWord {
		t.Errorf("Expected StatusWord 0x%04X, got 0x%04X", status.StatusWord, response.Status().StatusWord)
	}
	if response.Status().StateVar != status.StateVar {
		t.Errorf("Expected StateVar 0x%04X, got 0x%04X", status.StateVar, response.Status().StateVar)
	}
}

func TestRequestManager_StopMC_RoutesWithMonitoringRTCReply(t *testing.T) {
	baseClient, server := transport.NewMockTransportClientWithServer()

	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := protocol_command_tables.NewStopMotionControllerRequest()

	var wg sync.WaitGroup
	var response *protocol_command_tables.StopMotionControllerResponse
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_command_tables.StopMotionControllerResponse](rm, ctx, req)
	}()

	time.Sleep(50 * time.Millisecond)

	rm.mu.RLock()
	var cmdCount uint8
	for count, pending := range rm.pendingRequests {
		if pending.request == req {
			cmdCount = count
			break
		}
	}
	rm.mu.RUnlock()

	if cmdCount == 0 {
		t.Fatal("Expected pending StopMC request with RTC command count")
	}

	repBits := protocol_common.RespBitStatusWord |
		protocol_common.RespBitStateVar |
		protocol_common.RespBitActualPosition |
		protocol_common.RespBitDemandPosition |
		protocol_common.RespBitCurrent |
		protocol_common.RespBitWarnWord |
		protocol_common.RespBitErrorCode |
		protocol_common.RespBitMonitoringChannel |
		protocol_common.RespBitRTCReplyData

	status := &protocol_common.Status{
		ResponseBits:   repBits,
		StatusWord:     0x1234,
		StateVar:       0x5678,
		ActualPosition: 100,
		DemandPosition: 200,
		Current:        300,
		WarnWord:       0x0000,
		ErrorCode:      0x0000,
		MonitoringChannel: [4]int32{
			0x11111111,
			0x22222222,
			0x33333333,
			0x44444444,
		},
	}

	packet := constructMonitoringStatusTelegramWithRTCReply(
		0x00000000,
		repBits,
		status,
		cmdCount,
		0,
		0,
		protocol_command_tables.CommandCode.StopMotionController,
	)
	server.SendPacket(packet)

	wg.Wait()

	if err != nil {
		t.Fatalf("Expected StopMC request to complete, got error: %v", err)
	}
	if response == nil {
		t.Fatal("Expected StopMotionControllerResponse, got nil")
	}
	if response.RTCCounter() != cmdCount {
		t.Errorf("Expected cmdCount %d, got %d", cmdCount, response.RTCCounter())
	}
}

// TestRequestManager_RTCResync_RejectsInvalidCounters verifies that resync does NOT occur
// when invalid cmdCount values (0 or 15) are received, and the packet falls back to status parsing.
func TestRequestManager_RTCResync_RejectsInvalidCounters(t *testing.T) {
	// Create mock transport with server
	baseClient, server := transport.NewMockTransportClientWithServer()

	// Create RequestManager
	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Register a simple response registry for the test UPID
	testUPID := protocol_common.ParameterID(0x1234)
	protocol_rtc.RegisterResponseRegistry(testUPID, func(
		status *protocol_common.Status,
		value int32,
		upid uint16,
		rtcCounter uint8,
		rtcStatus uint8,
		cmdCode uint8,
	) protocol_common.Response {
		return protocol_rtc.NewRTCGetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
	})

	// Submit an RTC request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upid := uint16(0x1234)
	cmdCode := uint8(0x10) // ReadRAM
	req := protocol_rtc.NewRTCGetParamRequest(upid, cmdCode)

	// Start request in background
	var wg sync.WaitGroup
	var response *protocol_rtc.RTCGetParamResponse
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](rm, ctx, req)
	}()

	// Wait for request to be submitted and cmdCount allocated
	time.Sleep(50 * time.Millisecond)

	// Get the cmdCount that was allocated
	rm.mu.RLock()
	var clientCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			clientCmdCount = count
			break
		}
	}
	rm.mu.RUnlock()

	if clientCmdCount == 0 {
		t.Fatal("Failed to find pending RTC request with cmdCount")
	}

	// Test invalid cmdCount = 0
	reqBits := uint32(0x00000000)
	repBits := protocol_common.RespBitStatusWord |
		protocol_common.RespBitStateVar |
		protocol_common.RespBitRTCReplyData

	status := &protocol_common.Status{
		ResponseBits:   repBits,
		StatusWord:     0x0000,
		StateVar:       0x0000,
		ActualPosition: 0,
		DemandPosition: 0,
		Current:        0,
		WarnWord:       0x0000,
		ErrorCode:      0x0000,
	}

	value := int32(0x5678)

	// Send invalid cmdCount = 0
	invalidPacket0 := constructStatusTelegramWithRTCReply(reqBits, repBits, status, 0, upid, value, cmdCode)
	server.SendPacket(invalidPacket0)

	// Wait a moment
	time.Sleep(50 * time.Millisecond)

	// Verify resync did NOT occur (invalid counter rejected)
	rm.mu.RLock()
	resyncAttempted0 := rm.resyncAttempted
	// Verify request still has original cmdCount
	var stillPending bool
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			if count != clientCmdCount {
				t.Errorf("Expected cmdCount to remain %d, but changed to %d (resync should not occur for invalid counter)", clientCmdCount, count)
			}
			stillPending = true
			break
		}
	}
	rm.mu.RUnlock()

	if resyncAttempted0 {
		t.Error("Expected resyncAttempted to remain false (invalid counter 0 rejected), but it was true")
	}

	if !stillPending {
		t.Error("Expected request to still be pending after invalid counter 0")
	}

	// Test invalid cmdCount = 15
	invalidPacket15 := constructStatusTelegramWithRTCReply(reqBits, repBits, status, 15, upid, value, cmdCode)
	server.SendPacket(invalidPacket15)

	// Wait a moment
	time.Sleep(50 * time.Millisecond)

	// Verify resync did NOT occur (invalid counter rejected)
	rm.mu.RLock()
	resyncAttempted15 := rm.resyncAttempted
	rm.mu.RUnlock()

	if resyncAttempted15 {
		t.Error("Expected resyncAttempted to remain false (invalid counter 15 rejected), but it was true")
	}

	// Now send a matching reply with the original cmdCount to complete the request
	matchingPacket := constructStatusTelegramWithRTCReply(reqBits, repBits, status, clientCmdCount, upid, value, cmdCode)
	server.SendPacket(matchingPacket)

	// Wait for request to complete
	wg.Wait()

	// Verify the request completed successfully
	if err != nil {
		t.Fatalf("Expected RTC request to complete successfully, got error: %v", err)
	}

	if response == nil {
		t.Fatal("Expected RTC response, got nil")
	}

	t.Logf("Correctly rejected invalid counters (0 and 15): originalCmdCount=%d, resyncAttempted=%v", clientCmdCount, resyncAttempted15)
}

// TestRequestManager_RTCResync_DoesNotTriggerOnPrevCounterEcho verifies that resync does NOT
// trigger when the drive echoes the previous cmdCount (stale echo situation).
// This prevents false positives from common "drive is one behind" scenarios.
func TestRequestManager_RTCResync_DoesNotTriggerOnPrevCounterEcho(t *testing.T) {
	// Create mock transport with server
	baseClient, server := transport.NewMockTransportClientWithServer()

	// Create RequestManager
	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Register a simple response registry for the test UPID
	testUPID := protocol_common.ParameterID(0x1234)
	protocol_rtc.RegisterResponseRegistry(testUPID, func(
		status *protocol_common.Status,
		value int32,
		upid uint16,
		rtcCounter uint8,
		rtcStatus uint8,
		cmdCode uint8,
	) protocol_common.Response {
		return protocol_rtc.NewRTCGetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
	})

	// Submit an RTC request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upid := uint16(0x1234)
	cmdCode := uint8(0x10) // ReadRAM
	req := protocol_rtc.NewRTCGetParamRequest(upid, cmdCode)

	// Start request in background
	var wg sync.WaitGroup
	var response *protocol_rtc.RTCGetParamResponse
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](rm, ctx, req)
	}()

	// Wait for request to be submitted and cmdCount allocated
	time.Sleep(50 * time.Millisecond)

	// Get the cmdCount that was allocated
	rm.mu.RLock()
	var clientCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			clientCmdCount = count
			break
		}
	}
	rm.mu.RUnlock()

	if clientCmdCount == 0 {
		t.Fatal("Failed to find pending RTC request with cmdCount")
	}

	// Calculate previous counter (stale echo scenario)
	var prevCounter uint8
	if clientCmdCount == 1 {
		prevCounter = 14
	} else {
		prevCounter = clientCmdCount - 1
	}

	reqBits := uint32(0x00000000)
	repBits := protocol_common.RespBitStatusWord |
		protocol_common.RespBitStateVar |
		protocol_common.RespBitRTCReplyData

	status := &protocol_common.Status{
		ResponseBits:   repBits,
		StatusWord:     0x0000,
		StateVar:       0x0000,
		ActualPosition: 0,
		DemandPosition: 0,
		Current:        0,
		WarnWord:       0x0000,
		ErrorCode:      0x0000,
	}

	value := int32(0x5678)

	// Send status telegram with RTC reply data containing prev(clientCmdCount) - stale echo
	prevCounterPacket := constructStatusTelegramWithRTCReply(reqBits, repBits, status, prevCounter, upid, value, cmdCode)
	server.SendPacket(prevCounterPacket)

	// Wait a moment for processing
	time.Sleep(50 * time.Millisecond)

	// Verify resync did NOT occur (prev counter echo should be ignored)
	rm.mu.RLock()
	resyncAttempted := rm.resyncAttempted
	// Verify request still has original cmdCount (not re-keyed)
	var stillPending bool
	var stillHasOriginalCount bool
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			stillPending = true
			if count == clientCmdCount {
				stillHasOriginalCount = true
			}
			break
		}
	}
	rm.mu.RUnlock()

	if resyncAttempted {
		t.Error("Expected resyncAttempted to remain false (prev counter echo should not trigger resync), but it was true")
	}

	if !stillPending {
		t.Error("Expected request to still be pending after prev counter echo")
	}

	if !stillHasOriginalCount {
		t.Errorf("Expected request to still have original cmdCount %d, but it was re-keyed (resync should not occur)", clientCmdCount)
	}

	// Now send the real matching RTC reply with clientCmdCount to complete the request
	matchingPacket := constructStatusTelegramWithRTCReply(reqBits, repBits, status, clientCmdCount, upid, value, cmdCode)
	server.SendPacket(matchingPacket)

	// Wait for request to complete
	wg.Wait()

	// Verify the request completed successfully
	if err != nil {
		t.Fatalf("Expected RTC request to complete successfully, got error: %v", err)
	}

	if response == nil {
		t.Fatal("Expected RTC response, got nil")
	}

	t.Logf("Correctly ignored prev counter echo: clientCmdCount=%d, prevCounter=%d, resyncAttempted=%v", clientCmdCount, prevCounter, resyncAttempted)
}

// TestRequestManager_RTCResync_OnlyOnPlausibleDelta verifies that resync occurs
// for valid mismatches that are not prev(clientCmdCount) echoes.
// With the new predicate, any valid mismatch (not prev) triggers resync,
// regardless of the "delta" size. This test verifies the behavior with different mismatches.
func TestRequestManager_RTCResync_OnlyOnPlausibleDelta(t *testing.T) {
	// Create mock transport with server
	baseClient, server := transport.NewMockTransportClientWithServer()

	// Create RequestManager
	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Register a simple response registry for the test UPID
	testUPID := protocol_common.ParameterID(0x1234)
	protocol_rtc.RegisterResponseRegistry(testUPID, func(
		status *protocol_common.Status,
		value int32,
		upid uint16,
		rtcCounter uint8,
		rtcStatus uint8,
		cmdCode uint8,
	) protocol_common.Response {
		return protocol_rtc.NewRTCGetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
	})

	// Submit an RTC request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upid := uint16(0x1234)
	cmdCode := uint8(0x10) // ReadRAM
	req := protocol_rtc.NewRTCGetParamRequest(upid, cmdCode)

	// Start request in background
	var wg sync.WaitGroup
	var response *protocol_rtc.RTCGetParamResponse
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](rm, ctx, req)
	}()

	// Wait for request to be submitted and cmdCount allocated
	time.Sleep(50 * time.Millisecond)

	// Get the cmdCount that was allocated
	rm.mu.RLock()
	var clientCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			clientCmdCount = count
			break
		}
	}
	rm.mu.RUnlock()

	if clientCmdCount == 0 {
		t.Fatal("Failed to find pending RTC request with cmdCount")
	}

	// Set client cmdCount to 5 for testing
	rm.SetRTCCounterForTesting(4) // Next will be 5
	rm.mu.Lock()
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			pendingReq.commandCount = 5
			delete(rm.pendingRequests, count)
			rm.pendingRequests[5] = pendingReq
			clientCmdCount = 5
			break
		}
	}
	rm.mu.Unlock()

	reqBits := uint32(0x00000000)
	repBits := protocol_common.RespBitStatusWord |
		protocol_common.RespBitStateVar |
		protocol_common.RespBitRTCReplyData

	status := &protocol_common.Status{
		ResponseBits:   repBits,
		StatusWord:     0x0000,
		StateVar:       0x0000,
		ActualPosition: 0,
		DemandPosition: 0,
		Current:        0,
		WarnWord:       0x0000,
		ErrorCode:      0x0000,
	}

	value := int32(0x5678)

	// Test case 1: drive=6, client=5 (should resync: 6 != 5 and 6 != 4)
	unmatchedCmdCount6 := uint8(6)
	unmatchedPacket6 := constructStatusTelegramWithRTCReply(reqBits, repBits, status, unmatchedCmdCount6, upid, value, cmdCode)
	server.SendPacket(unmatchedPacket6)

	// Wait for resync
	time.Sleep(100 * time.Millisecond)

	// Verify resync occurred (6 is valid, not equal to 5, and not prev(5)=4)
	rm.mu.RLock()
	var newCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			newCmdCount = count
			break
		}
	}
	resyncAttempted6 := rm.resyncAttempted
	rm.mu.RUnlock()

	if newCmdCount == clientCmdCount {
		t.Error("Expected request to be re-keyed after unmatched reply with drive=6, but cmdCount didn't change")
	}

	if !resyncAttempted6 {
		t.Error("Expected resyncAttempted to be true after unmatched reply with drive=6 (valid mismatch, not prev)")
	}

	// Verify new cmdCount is next(6) = 7 (deterministic)
	expectedNewCount1 := uint8(7)
	if newCmdCount != expectedNewCount1 {
		t.Errorf("Expected new cmdCount to be next(6)=%d, but got %d", expectedNewCount1, newCmdCount)
	}

	// Reset for next test
	rm.mu.Lock()
	rm.resyncAttempted = false
	rm.mu.Unlock()

	// Update client cmdCount back to 5
	rm.mu.Lock()
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			pendingReq.commandCount = 5
			delete(rm.pendingRequests, count)
			rm.pendingRequests[5] = pendingReq
			break
		}
	}
	rm.mu.Unlock()

	// Test case 2: drive=10, client=5 (should resync: 10 != 5 and 10 != 4)
	// With new predicate, any valid mismatch (not prev) triggers resync, regardless of "delta" size
	unmatchedCmdCount10 := uint8(10)
	unmatchedPacket10 := constructStatusTelegramWithRTCReply(reqBits, repBits, status, unmatchedCmdCount10, upid, value, cmdCode)
	server.SendPacket(unmatchedPacket10)

	// Wait for resync
	time.Sleep(100 * time.Millisecond)

	// Verify resync occurred (10 is valid, not equal to 5, and not prev(5)=4)
	rm.mu.RLock()
	var finalCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			finalCmdCount = count
			break
		}
	}
	resyncAttempted10 := rm.resyncAttempted
	rm.mu.RUnlock()

	if finalCmdCount == 5 {
		t.Error("Expected request to be re-keyed after unmatched reply with drive=10, but cmdCount didn't change")
	}

	if !resyncAttempted10 {
		t.Error("Expected resyncAttempted to be true after unmatched reply with drive=10 (valid mismatch, not prev)")
	}

	// Verify new cmdCount is next(10) = 11 (deterministic)
	expectedNewCount2 := uint8(11)
	if finalCmdCount != expectedNewCount2 {
		t.Errorf("Expected new cmdCount to be next(10)=%d, but got %d", expectedNewCount2, finalCmdCount)
	}

	// Now send a matching reply with the new cmdCount to complete the request
	matchingPacket := constructStatusTelegramWithRTCReply(reqBits, repBits, status, finalCmdCount, upid, value, cmdCode)
	server.SendPacket(matchingPacket)

	// Wait for request to complete
	wg.Wait()

	// Verify the request completed successfully
	if err != nil {
		t.Fatalf("Expected RTC request to complete successfully, got error: %v", err)
	}

	if response == nil {
		t.Fatal("Expected RTC response, got nil")
	}

	t.Logf("Resync occurred for valid mismatches: clientCmdCount=5, drive=6 (resync to 7), drive=10 (resync to 11)")
}

// TestRequestManager_RTCResync_ChoosesNextOfDriveCounter verifies that when resync occurs,
// the new cmdCount is deterministically chosen as next(driveCmdCount).
// This ensures predictable resync behavior.
func TestRequestManager_RTCResync_ChoosesNextOfDriveCounter(t *testing.T) {
	// Create mock transport with server
	baseClient, server := transport.NewMockTransportClientWithServer()

	// Create RequestManager
	rm := NewRequestManager(baseClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Register a simple response registry for the test UPID
	testUPID := protocol_common.ParameterID(0x1234)
	protocol_rtc.RegisterResponseRegistry(testUPID, func(
		status *protocol_common.Status,
		value int32,
		upid uint16,
		rtcCounter uint8,
		rtcStatus uint8,
		cmdCode uint8,
	) protocol_common.Response {
		return protocol_rtc.NewRTCGetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
	})

	// Submit an RTC request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upid := uint16(0x1234)
	cmdCode := uint8(0x10) // ReadRAM
	req := protocol_rtc.NewRTCGetParamRequest(upid, cmdCode)

	// Start request in background
	var wg sync.WaitGroup
	var response *protocol_rtc.RTCGetParamResponse
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err = SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](rm, ctx, req)
	}()

	// Wait for request to be submitted and cmdCount allocated
	time.Sleep(50 * time.Millisecond)

	// Get the cmdCount that was allocated
	rm.mu.RLock()
	var clientCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			clientCmdCount = count
			break
		}
	}
	rm.mu.RUnlock()

	if clientCmdCount == 0 {
		t.Fatal("Failed to find pending RTC request with cmdCount")
	}

	// Choose a drive cmdCount that will trigger resync:
	// - Must be valid (1..14)
	// - Must not equal clientCmdCount
	// - Must not equal prev(clientCmdCount) (to avoid stale echo rejection)
	var driveCmdCount uint8
	if clientCmdCount >= 3 {
		// Use a value that's not client and not prev(client)
		driveCmdCount = clientCmdCount - 2
	} else {
		// If client is 1 or 2, use a value ahead
		driveCmdCount = clientCmdCount + 3
		if driveCmdCount > 14 {
			driveCmdCount = 14
		}
	}

	// Calculate expected new cmdCount: next(driveCmdCount)
	var expectedNewCmdCount uint8
	if driveCmdCount == 14 {
		expectedNewCmdCount = 1
	} else {
		expectedNewCmdCount = driveCmdCount + 1
	}

	reqBits := uint32(0x00000000)
	repBits := protocol_common.RespBitStatusWord |
		protocol_common.RespBitStateVar |
		protocol_common.RespBitRTCReplyData

	status := &protocol_common.Status{
		ResponseBits:   repBits,
		StatusWord:     0x0000,
		StateVar:       0x0000,
		ActualPosition: 0,
		DemandPosition: 0,
		Current:        0,
		WarnWord:       0x0000,
		ErrorCode:      0x0000,
	}

	value := int32(0x5678)

	// Send unmatched RTC reply packet with driveCmdCount that should trigger resync
	unmatchedPacket := constructStatusTelegramWithRTCReply(reqBits, repBits, status, driveCmdCount, upid, value, cmdCode)
	server.SendPacket(unmatchedPacket)

	// Wait for resync to occur
	time.Sleep(100 * time.Millisecond)

	// Verify resync occurred and new cmdCount is next(driveCmdCount)
	rm.mu.RLock()
	resyncAttempted := rm.resyncAttempted
	var newCmdCount uint8
	for count, pendingReq := range rm.pendingRequests {
		if pendingReq.request == req {
			newCmdCount = count
			break
		}
	}
	rm.mu.RUnlock()

	if !resyncAttempted {
		t.Errorf("Expected resync to occur for driveCmdCount=%d (not equal to clientCmdCount=%d and not prev)", driveCmdCount, clientCmdCount)
	}

	if newCmdCount != expectedNewCmdCount {
		t.Errorf("Expected new cmdCount to be next(driveCmdCount)=%d, but got %d", expectedNewCmdCount, newCmdCount)
	}

	// Now send matching reply with the new cmdCount to complete the request
	matchingPacket := constructStatusTelegramWithRTCReply(reqBits, repBits, status, newCmdCount, upid, value, cmdCode)
	server.SendPacket(matchingPacket)

	// Wait for request to complete
	wg.Wait()

	// Verify the request completed successfully
	if err != nil {
		t.Fatalf("Expected RTC request to complete successfully, got error: %v", err)
	}

	if response == nil {
		t.Fatal("Expected RTC response, got nil")
	}

	t.Logf("Resync chose deterministic counter: clientCmdCount=%d, driveCmdCount=%d, newCmdCount=%d (next(drive))", clientCmdCount, driveCmdCount, newCmdCount)
}

// TestShouldResync_Invariants verifies the shouldResync predicate invariants across all edge cases.
// This table-driven test ensures future edits don't accidentally reintroduce "delta math" or allow stale echo resync.
func TestShouldResync_Invariants(t *testing.T) {
	tests := []struct {
		name           string
		driveCmdCount  uint8
		clientCmdCount uint8
		expected       bool
		reason         string
	}{
		// Invalid counters => false
		{"invalid_drive_0", 0, 1, false, "invalid drive counter 0"},
		{"invalid_drive_15", 15, 1, false, "invalid drive counter 15"},
		{"invalid_drive_16", 16, 1, false, "invalid drive counter >14"},
		{"invalid_client_0", 1, 0, false, "invalid client counter 0"},
		{"invalid_client_15", 1, 15, false, "invalid client counter 15"},
		{"invalid_client_16", 1, 16, false, "invalid client counter >14"},
		{"both_invalid", 0, 15, false, "both counters invalid"},

		// drive == client => false
		{"match_1", 1, 1, false, "counters match"},
		{"match_7", 7, 7, false, "counters match"},
		{"match_14", 14, 14, false, "counters match"},

		// drive == prev(client) => false (stale echo)
		{"prev_echo_1", 14, 1, false, "drive is prev(client) - stale echo"},
		{"prev_echo_2", 1, 2, false, "drive is prev(client) - stale echo"},
		{"prev_echo_7", 6, 7, false, "drive is prev(client) - stale echo"},
		{"prev_echo_14", 13, 14, false, "drive is prev(client) - stale echo"},

		// All other valid pairs => true
		// Note: drive=1, client=2 is NOT valid (drive == prev(2) = 1, so it's a stale echo)
		// Note: drive=14, client=1 is NOT valid (drive == prev(1) = 14, so it's a stale echo)
		{"valid_mismatch_1_3", 1, 3, true, "valid mismatch, not prev"},
		{"valid_mismatch_1_4", 1, 4, true, "valid mismatch, not prev"},
		{"valid_mismatch_2_1", 2, 1, true, "valid mismatch, not prev"},
		{"valid_mismatch_2_4", 2, 4, true, "valid mismatch, not prev"},
		{"valid_mismatch_2_5", 2, 5, true, "valid mismatch, not prev"},
		{"valid_mismatch_7_1", 7, 1, true, "valid mismatch, not prev"},
		{"valid_mismatch_7_5", 7, 5, true, "valid mismatch, not prev"},
		{"valid_mismatch_14_2", 14, 2, true, "valid mismatch, not prev"},
		{"valid_mismatch_14_3", 14, 3, true, "valid mismatch, not prev"},
		{"valid_mismatch_14_13", 14, 13, true, "valid mismatch, not prev"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldResync(tt.driveCmdCount, tt.clientCmdCount)
			if result != tt.expected {
				t.Errorf("shouldResync(drive=%d, client=%d) = %v, want %v (%s)",
					tt.driveCmdCount, tt.clientCmdCount, result, tt.expected, tt.reason)
			}
		})
	}
}
