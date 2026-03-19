package test

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	protocol_common "gsail-go/linmot/protocol/common"
	protocol_control_word "gsail-go/linmot/protocol/control_word"
	protocol_motion_control "gsail-go/linmot/protocol/motion_control"
	protocol_rtc "gsail-go/linmot/protocol/rtc"
	protocol_command_tables "gsail-go/linmot/protocol/rtc/command_tables"
	"gsail-go/linmot/transport"
)

// MockLinMot simulates a LinMot drive for testing.
// Thread-safe: can be accessed concurrently from multiple goroutines.
type MockLinMot struct {
	transport transport.Server
	mu        sync.RWMutex
	ramParams map[uint16]int32
	romParams map[uint16]int32
	status    *protocol_common.Status

	// Test control knobs
	simulateError   bool
	persistentError bool // If true, error acknowledge edge will not clear the error
	errorCode       uint16
	simulateWarning bool
	warningCode     uint16

	// RTC counter deduplication
	lastProcessedRtcCounter uint8
	lastRtcResponse         protocol_common.Response

	// CommandCount semantics: track last command count per command code
	// Real hardware only executes commands when the count changes
	lastCommandCount map[uint8]uint8 // cmdCode -> last counter value

	// Command Table state
	mcRunning       bool              // Motion Controller running state
	commandTable    map[uint16][]byte // CT entry ID -> entry data
	ctWriteBuffers  map[uint16][]byte // Entry ID -> partial write buffer (for progressive writes)
	ctReadPositions map[uint16]int    // Entry ID -> read position (for progressive reads)

	// State Machine state
	mainState        protocol_control_word.MainState // Current state machine state
	errorAcknowledge bool                            // Track error acknowledge edge
	isHomed          bool                            // Track if homing has been completed

	// Monitoring Channel state
	monitoringChannelConfig [4]uint16        // UPIDs configured for each monitoring channel (0x20A8-0x20AB)
	monitoredVariables      map[uint16]int32 // Mock variable storage for monitoring

	// Motion Control state
	mcCounter        uint8  // Motion Control counter (1-4 range)
	targetPosition   int32  // Target position for VAI commands
	demandPosition   int32  // Demand position (intermediate setpoint)
	motionVelocity   uint32 // Current motion velocity
	motionInProgress bool   // Whether motion is active

	// Control
	stopChan chan struct{}
	stopped  atomic.Bool

	statusDelay    time.Duration
	errorTextDelay time.Duration
}

// mockErrorTexts contains mock error code descriptions for testing.
var mockErrorTexts = map[uint16]string{
	0x0020: "Position Lag Error",
	0x0030: "Motor Overtemperature",
	0x0040: "Drive Overload",
	0x0050: "Supply Voltage Too Low",
	0x0060: "Supply Voltage Too High",
}

var mockLinMotDebug = os.Getenv("LINMOT_TEST_DEBUG") == "1"

// NewMockLinMot creates a new mock LinMot drive for testing.
func NewMockLinMot(transportServer transport.Server) *MockLinMot {
	return &MockLinMot{
		transport:               transportServer,
		ramParams:               make(map[uint16]int32),
		romParams:               make(map[uint16]int32),
		lastProcessedRtcCounter: 0, // 0 indicates no counter processed yet
		lastRtcResponse:         nil,
		lastCommandCount:        make(map[uint8]uint8), // Track command counts per cmdCode
		mcRunning:               true,                  // MC starts running
		commandTable:            make(map[uint16][]byte),
		ctWriteBuffers:          make(map[uint16][]byte),
		ctReadPositions:         make(map[uint16]int),
		mainState:               protocol_control_word.State_SwitchOnDisabled,
		errorAcknowledge:        false,
		isHomed:                 false,
		monitoredVariables:      make(map[uint16]int32),
		mcCounter:               0, // MC counter starts at 0, first command will use 1
		targetPosition:          100000,
		demandPosition:          100000,
		motionVelocity:          0,
		motionInProgress:        false,
		status: &protocol_common.Status{
			StatusWord:     0x0000, // Start with all bits clear
			StateVar:       0x0100, // State 1: Switch On Disabled
			ActualPosition: 100000,
			DemandPosition: 100000,
			Current:        100,
			WarnWord:       0,
			ErrorCode:      0,
		},
		stopChan: make(chan struct{}),
	}
}

// Start starts the drive processing loop.
func (server *MockLinMot) Start() {
	go server.process()
}

// process handles incoming requests and sends responses.
func (server *MockLinMot) process() {
	for {
		select {
		case <-server.stopChan:
			return
		default:
			// 1. Transport: Receive (I/O)
			requestPacket := server.transport.RecvPacket()
			if requestPacket == nil {
				return
			}

			// 2. Parse packet to request - inline dispatcher logic
			reqBits, repBits, err := protocol_common.ReadPacketHeader(requestPacket)
			if err != nil {
				// Send error response
				errorStatus := &protocol_common.Status{ErrorCode: 0xFFFF}
				responsePacket, _ := protocol_common.NewStatusResponse(errorStatus).WritePacket()
				server.transport.SendPacket(responsePacket)
				continue
			}

			var request protocol_common.Request
			var rtcCounter uint8

			// Check if this is a Control Word request
			if reqBits&protocol_common.RequestFlags.ControlWord != 0 {
				// Control Word request (10 bytes)
				req, err := protocol_control_word.ReadControlWordRequest(requestPacket)
				if err != nil {
					errorStatus := &protocol_common.Status{ErrorCode: 0xFFFF}
					responsePacket, _ := protocol_common.NewStatusResponse(errorStatus).WritePacket()
					server.transport.SendPacket(responsePacket)
					continue
				}
				request = req
				rtcCounter = 0 // Control Word requests don't have RTC counters
			} else if reqBits&protocol_common.RequestFlags.MotionControl != 0 {
				// Motion Control request (40 bytes)
				req, err := protocol_motion_control.ReadMCRequest(requestPacket)
				if err != nil {
					errorStatus := &protocol_common.Status{ErrorCode: 0xFFFF}
					responsePacket, _ := protocol_common.NewStatusResponse(errorStatus).WritePacket()
					server.transport.SendPacket(responsePacket)
					continue
				}
				request = req
				rtcCounter = 0 // MC requests don't have RTC counters
			} else if reqBits&0x00000004 != 0 { // Bit 2: RTC Command (use literal to validate correct bit)
				// RTC command request - extract counter and parse
				if len(requestPacket) < 10 {
					errorStatus := &protocol_common.Status{ErrorCode: 0xFFFF}
					responsePacket, _ := protocol_common.NewStatusResponse(errorStatus).WritePacket()
					server.transport.SendPacket(responsePacket)
					continue
				}
				cmdCode := requestPacket[9]
				// Check if this is a "get/read" style command
				isGetCommand := cmdCode == protocol_rtc.CommandCode.ReadROM ||
					cmdCode == protocol_rtc.CommandCode.ReadRAM ||
					cmdCode == protocol_rtc.CommandCode.GetMinValue ||
					cmdCode == protocol_rtc.CommandCode.GetMaxValue ||
					cmdCode == protocol_rtc.CommandCode.GetDefaultValue

				if isGetCommand {
					event, counter, err := protocol_rtc.ReadRTCGetParamRequest(requestPacket)
					if err != nil {
						errorStatus := &protocol_common.Status{ErrorCode: 0xFFFF}
						responsePacket, _ := protocol_common.NewStatusResponse(errorStatus).WritePacket()
						server.transport.SendPacket(responsePacket)
						continue
					}
					request = event.(protocol_common.Request)
					rtcCounter = counter
				} else {
					event, counter, err := protocol_rtc.ReadRTCSetParamRequest(requestPacket)
					if err != nil {
						errorStatus := &protocol_common.Status{ErrorCode: 0xFFFF}
						responsePacket, _ := protocol_common.NewStatusResponse(errorStatus).WritePacket()
						server.transport.SendPacket(responsePacket)
						continue
					}
					request = event.(protocol_common.Request)
					rtcCounter = counter
				}
			} else if repBits&protocol_common.ResponseFlags.Standard != 0 {
				// Status or Monitoring Status request (8 bytes, header only)
				if len(requestPacket) == protocol_common.PacketHeaderSize {
					// Check if monitoring channel is requested (bit 7)
					if repBits&protocol_common.ResponseFlags.MonitoringChannel != 0 {
						req, err := protocol_common.ReadMonitoringStatusRequest(requestPacket)
						if err != nil {
							errorStatus := &protocol_common.Status{ErrorCode: 0xFFFF}
							responsePacket, _ := protocol_common.NewStatusResponse(errorStatus).WritePacket()
							server.transport.SendPacket(responsePacket)
							continue
						}
						request = req
					} else {
						req, err := protocol_common.ReadStatusRequest(requestPacket)
						if err != nil {
							errorStatus := &protocol_common.Status{ErrorCode: 0xFFFF}
							responsePacket, _ := protocol_common.NewStatusResponse(errorStatus).WritePacket()
							server.transport.SendPacket(responsePacket)
							continue
						}
						request = req
					}
					rtcCounter = 0 // Status requests don't have counters
				} else {
					// Invalid packet length
					errorStatus := &protocol_common.Status{ErrorCode: 0xFFFF}
					responsePacket, _ := protocol_common.NewStatusResponse(errorStatus).WritePacket()
					server.transport.SendPacket(responsePacket)
					continue
				}
			} else {
				// Invalid request flags
				errorStatus := &protocol_common.Status{ErrorCode: 0xFFFF}
				responsePacket, _ := protocol_common.NewStatusResponse(errorStatus).WritePacket()
				server.transport.SendPacket(responsePacket)
				continue
			}

			// 3. Route to appropriate handler based on request type
			var response protocol_common.Response
			if rtcCounter != 0 {
				// RTC request - always process normally (no duplicate detection for simplicity)
				// Real hardware duplicate detection would need to rebuild response with new counter
				response = server.handleRtcRequest(request, rtcCounter)
			} else {
				// Non-RTC request (StatusRequest)
				response = server.handleRequest(request)
			}

			// 4. Convert response request to packet
			responsePacket, _ := response.WritePacket()

			// 5. Transport: Send (I/O)
			server.transport.SendPacket(responsePacket)
		}
	}
}

// handleRequest processes a non-RTC request and returns a Response.
func (server *MockLinMot) handleRequest(request protocol_common.Request) protocol_common.Response {
	// Check if this is a Control Word request
	if cwRequest, ok := request.(*protocol_control_word.ControlWordRequest); ok {
		return server.handleControlWordRequest(cwRequest)
	}

	// Check if this is an Motion Control request
	if mcRequest, ok := request.(*protocol_motion_control.MCCommandRequest); ok {
		return server.handleMCRequest(mcRequest)
	}

	// Check if this is a Monitoring Status request
	if _, ok := request.(*protocol_common.MonitoringStatusRequest); ok {
		return server.handleMonitoringStatusRequest()
	}

	// Otherwise it's a StatusRequest
	status := server.getStatus()
	return protocol_common.NewStatusResponse(status)
}

// handleRtcRequest processes an RTC request and returns a Response.
func (server *MockLinMot) handleRtcRequest(request protocol_common.Request, rtcCounter uint8) protocol_common.Response {
	switch typedRequest := request.(type) {
	case *protocol_rtc.RTCGetParamRequest:
		cmdCode := typedRequest.CmdCode()
		value, err := server.processRTCRead(typedRequest.UPID(), cmdCode)
		if err != nil {
			return protocol_common.NewStatusResponse(
				protocol_common.NewErrorStatus(0xFFFF),
			)
		}
		status := server.getStatus()
		return protocol_rtc.NewRTCGetParamResponseWithCmdCode(
			status,
			value,
			typedRequest.UPID(),
			rtcCounter,
			0x00,
			cmdCode,
		)

	case *protocol_rtc.RTCSetParamRequest:
		cmdCode := typedRequest.CmdCode()

		// Check if this is a CT command (use IsCTCommand helper for consistency)
		if protocol_rtc.IsCTCommand(cmdCode) {
			// CT command - extract w2, w3, w4 from upid and value
			upid := typedRequest.UPID()
			value := typedRequest.Value()
			w2 := upid
			w3 := uint16(value >> 16)
			w4 := uint16(value & 0xFFFF)

			// Enforce CommandCount semantics: only execute if count changed
			server.mu.Lock()
			lastCount, hasLastCount := server.lastCommandCount[cmdCode]
			shouldExecute := !hasLastCount || rtcCounter != lastCount
			if shouldExecute {
				server.lastCommandCount[cmdCode] = rtcCounter
			}
			server.mu.Unlock()

			var w3Out, w4Out uint16
			var rtcStatus uint8

			if shouldExecute {
				// Count changed - execute the command
				w3Out, w4Out, rtcStatus = server.processCTCommand(cmdCode, w2, w3, w4)
			} else {
				// Count unchanged - return success without executing (simulating hardware behavior)
				// Return the same response as if command succeeded, but without side effects
				w3Out, w4Out, rtcStatus = 0, 0, 0x00
			}

			status := server.getStatus()
			responseValue := int32(uint32(w3Out)<<16 | uint32(w4Out))

			return protocol_rtc.NewRTCSetParamResponseWithCmdCode(
				status,
				responseValue,
				upid,
				rtcCounter,
				rtcStatus,
				cmdCode,
			)
		}

		// Handle special commands that need custom responses
		upid := typedRequest.UPID()
		value := typedRequest.Value()
		status := server.getStatus()

		switch cmdCode {
		case protocol_rtc.CommandCode.StartGettingUPIDList:
			// Acknowledge start - return OK
			return protocol_rtc.NewRTCSetParamResponseWithCmdCode(
				status,
				0,
				upid,
				rtcCounter,
				0x00,
				cmdCode,
			)

		case protocol_rtc.CommandCode.GetNextUPIDListItem:
			// Return end-of-list status (0xC6) immediately (mock has no UPIDs)
			return protocol_rtc.NewRTCSetParamResponseWithCmdCode(
				status,
				0,
				0,
				rtcCounter,
				0xC6, // End of list
				cmdCode,
			)

		case protocol_rtc.CommandCode.StartGettingModifiedUPIDList:
			// Acknowledge start - return OK
			return protocol_rtc.NewRTCSetParamResponseWithCmdCode(
				status,
				0,
				upid,
				rtcCounter,
				0x00,
				cmdCode,
			)

		case protocol_rtc.CommandCode.GetNextModifiedUPIDListItem:
			// Return end-of-list status (0xC6) immediately (mock has no modified UPIDs)
			return protocol_rtc.NewRTCSetParamResponseWithCmdCode(
				status,
				0,
				0,
				rtcCounter,
				0xC6, // End of list
				cmdCode,
			)

		case 0x70: // GetErrorLogEntryCounter
			// Return 0 logged errors, 0 occurred errors (no errors in mock)
			return protocol_rtc.NewRTCSetParamResponseWithCmdCode(
				status,
				0, // w3=0 (logged), w4=0 (occurred)
				upid,
				rtcCounter,
				0x00,
				cmdCode,
			)

		case 0x71, 0x72, 0x73: // Error log entry code/time commands
			// Return 0 values (no error data in mock)
			return protocol_rtc.NewRTCSetParamResponseWithCmdCode(
				status,
				0,
				upid,
				rtcCounter,
				0x00,
				cmdCode,
			)

		case 0x74: // GetErrorCodeTextStringlet
			if server.errorTextDelay > 0 {
				time.Sleep(server.errorTextDelay)
			}
			// upid = error code, value = stringlet number in low word (word 3)
			errorCode := upid
			stringletNum := uint16(value & 0xFFFF)

			// Get error text from mock database
			errorText, exists := mockErrorTexts[errorCode]
			if !exists {
				errorText = "Unknown Error"
			}

			// Convert to bytes
			textBytes := []byte(errorText)

			// Calculate which 4-byte chunk to return
			startIdx := int(stringletNum) * 4

			// Extract 4 bytes for this stringlet
			var stringletBytes [4]byte
			for i := 0; i < 4; i++ {
				idx := startIdx + i
				if idx < len(textBytes) {
					stringletBytes[i] = textBytes[idx]
				} else {
					stringletBytes[i] = 0 // Null padding
				}
			}

			// Pack 4 bytes into int32 value
			stringletValue := int32(stringletBytes[0]) |
				int32(stringletBytes[1])<<8 |
				int32(stringletBytes[2])<<16 |
				int32(stringletBytes[3])<<24

			return protocol_rtc.NewRTCSetParamResponseWithCmdCode(
				status,
				stringletValue,
				errorCode, // Return error code in UPID field
				rtcCounter,
				0x00,
				cmdCode,
			)
		}

		// Regular parameter write
		err := server.processRTCWrite(upid, value, cmdCode)
		if err != nil {
			return protocol_common.NewStatusResponse(
				protocol_common.NewErrorStatus(0xFFFF),
			)
		}
		return protocol_rtc.NewRTCSetParamResponseWithCmdCode(
			status,
			value,
			upid,
			rtcCounter,
			0x00,
			cmdCode,
		)

	default:
		// Should never reach here due to interface-based routing in process()
		return protocol_common.NewStatusResponse(
			protocol_common.NewErrorStatus(0xFFFF),
		)
	}
}

// getStatus generates status data for a status request.
func (server *MockLinMot) getStatus() *protocol_common.Status {
	server.mu.RLock()
	defer server.mu.RUnlock()

	if server.statusDelay > 0 {
		time.Sleep(server.statusDelay)
	}

	status := &protocol_common.Status{
		StatusWord:     server.status.StatusWord,
		StateVar:       server.status.StateVar,
		ActualPosition: server.status.ActualPosition,
		DemandPosition: server.status.DemandPosition,
		Current:        server.status.Current,
		WarnWord:       server.status.WarnWord,
		ErrorCode:      server.status.ErrorCode,
	}

	// Apply error/warning simulation
	if server.simulateError {
		status.ErrorCode = server.errorCode
	}
	if server.simulateWarning {
		status.WarnWord = server.warningCode
	}

	return status
}

// handleMonitoringStatusRequest processes a monitoring status request and returns a response with monitoring data.
func (server *MockLinMot) handleMonitoringStatusRequest() protocol_common.Response {
	server.mu.RLock()
	defer server.mu.RUnlock()

	// Get the base status
	status := server.getStatus()

	// Populate monitoring channel values based on configured UPIDs
	for i := 0; i < 4; i++ {
		configuredUPID := server.monitoringChannelConfig[i]
		if configuredUPID == 0 {
			// No UPID configured for this channel, set to 0
			status.MonitoringChannel[i] = 0
			continue
		}

		// Check if there's a mock value for this UPID
		if mockValue, ok := server.monitoredVariables[configuredUPID]; ok {
			status.MonitoringChannel[i] = mockValue
		} else {
			// Simulate realistic values based on common UPIDs
			// For position-related UPIDs, return actual position
			// For other UPIDs, check ramParams or return 0
			if value, ok := server.ramParams[configuredUPID]; ok {
				status.MonitoringChannel[i] = value
			} else {
				// Default to actual position for unknown UPIDs (common use case)
				status.MonitoringChannel[i] = status.ActualPosition
			}
		}
	}

	return protocol_common.NewMonitoringStatusResponse(status)
}

// handleMCRequest processes an Motion Control request and returns a response.
func (server *MockLinMot) handleMCRequest(request *protocol_motion_control.MCCommandRequest) protocol_common.Response {
	server.mu.Lock()
	defer server.mu.Unlock()

	header := request.Header()

	// Update MC counter (echo from request)
	server.mcCounter = header.Counter

	// Simulate VAI command execution based on Master ID and Sub ID
	if header.MasterID == protocol_motion_control.MasterIDs.VAI {
		server.handleVAICommand(request)
	}
	// Add handling for other Master IDs as needed (Predefined VAI, Streaming, etc.)

	// Update StateVar to include MC counter in low nibble (bits 0-3)
	stateVarHigh := (server.status.StateVar & 0xFF00) | (uint16(server.mainState) << 8)
	stateVarLow := server.mcCounter & 0x0F
	server.status.StateVar = stateVarHigh | uint16(stateVarLow)

	// Build status response (can't call getStatus() as it would deadlock with RLock)
	status := &protocol_common.Status{
		StatusWord:     server.status.StatusWord,
		StateVar:       server.status.StateVar,
		ActualPosition: server.status.ActualPosition,
		DemandPosition: server.status.DemandPosition,
		Current:        server.status.Current,
		WarnWord:       server.status.WarnWord,
		ErrorCode:      server.status.ErrorCode,
	}

	// Apply error/warning simulation
	if server.simulateError {
		status.ErrorCode = server.errorCode
	}
	if server.simulateWarning {
		status.WarnWord = server.warningCode
	}

	// Return MC response with echoed counter
	return protocol_motion_control.NewMCCommandResponse(status, server.mcCounter)
}

// handleVAICommand processes VAI commands and updates motion state.
func (server *MockLinMot) handleVAICommand(request *protocol_motion_control.MCCommandRequest) {
	header := request.Header()
	parameters := request.Parameters()

	// Extract position/increment from first 2 words (int32, little-endian)
	// Position is encoded as: parameters[0] = low word, parameters[1] = high word
	// Reconstruct int32 from two uint16 words (little-endian)
	positionValue := int32(uint32(parameters[0]) | (uint32(parameters[1]) << 16))

	switch header.SubID {
	case 0x0: // VAI Go To Pos
		// Extract target position from parameters[0-1] (int32, little-endian)
		server.motionInProgress = true
		server.targetPosition = positionValue
		server.demandPosition = server.targetPosition
		server.status.DemandPosition = server.demandPosition

	case 0x1: // VAI Increment Dem Pos
		// Extract increment from parameters[0-1] (int32, little-endian)
		server.motionInProgress = true
		server.demandPosition += positionValue
		server.targetPosition = server.demandPosition
		server.status.DemandPosition = server.demandPosition

	case 0x2: // VAI Increment Target Pos
		// Extract increment from parameters[0-1] (int32, little-endian)
		server.motionInProgress = true
		server.targetPosition += positionValue
		server.demandPosition = server.targetPosition
		server.status.DemandPosition = server.demandPosition

	case 0x7: // VAI Stop
		server.motionInProgress = false

	case 0x3: // VAI Go To Pos From Act Pos And Act Vel
		// Extract target position from parameters[0-1] (int32, little-endian)
		server.motionInProgress = true
		server.targetPosition = positionValue
		server.demandPosition = server.targetPosition
		server.status.DemandPosition = server.demandPosition

	case 0x4: // VAI Go To Pos From Act Pos Dem Vel Zero
		// Extract target position from parameters[0-1] (int32, little-endian)
		server.motionInProgress = true
		server.targetPosition = positionValue
		server.demandPosition = server.targetPosition
		server.status.DemandPosition = server.demandPosition

	case 0x5: // VAI Increment Act Pos
		// Extract increment from parameters[0-1] (int32, little-endian)
		server.motionInProgress = true
		server.status.ActualPosition += positionValue
		server.demandPosition = server.status.ActualPosition
		server.targetPosition = server.demandPosition
		server.status.DemandPosition = server.demandPosition

	case 0x8: // VAI Go To Analog Pos
		// No position parameter - uses analog input
		server.motionInProgress = true
		// For testing, just increment position slightly
		server.targetPosition = server.status.ActualPosition + 10000
		server.demandPosition = server.targetPosition
		server.status.DemandPosition = server.demandPosition

	case 0xA: // VAI Go To Pos On Rising Trigger
		// Extract target position from parameters[0-1] (int32, little-endian)
		server.motionInProgress = true
		server.targetPosition = positionValue
		server.demandPosition = server.targetPosition
		server.status.DemandPosition = server.demandPosition

	case 0x17: // VAI Change Motion Params On Positive Transition
		// Extract transition position from parameters[0-1] (int32, little-endian)
		// For testing, just update target to transition position
		server.motionInProgress = true
		server.targetPosition = positionValue
		server.demandPosition = server.targetPosition
		server.status.DemandPosition = server.demandPosition
	}

	// Simulate motion: move actual position toward demand position
	if server.motionInProgress {
		// Simplified motion: move actual toward demand (instant for testing)
		server.status.ActualPosition = server.demandPosition
	}
}

// processRTCRead processes an RTC read request and returns the value.
func (server *MockLinMot) processRTCRead(upid uint16, cmdCode uint8) (int32, error) {
	server.mu.RLock()
	defer server.mu.RUnlock()

	// Handle monitoring channel configuration reads (0x20A8-0x20AB)
	if upid >= uint16(protocol_common.PUID.MonitoringChannel1UPID) && upid <= uint16(protocol_common.PUID.MonitoringChannel4UPID) {
		channelIndex := upid - uint16(protocol_common.PUID.MonitoringChannel1UPID)
		return int32(server.monitoringChannelConfig[channelIndex]), nil
	}

	// Determine which parameter store to use based on command code
	var value int32
	var ok bool

	switch cmdCode {
	case protocol_rtc.CommandCode.ReadROM:
		value, ok = server.romParams[upid]
	case protocol_rtc.CommandCode.ReadRAM:
		value, ok = server.ramParams[upid]
	case protocol_rtc.CommandCode.GetMinValue:
		// Return simulated min value (e.g., -100000 for position parameters)
		value = -100000
		ok = true
	case protocol_rtc.CommandCode.GetMaxValue:
		// Return simulated max value (e.g., 100000 for position parameters)
		value = 100000
		ok = true
	case protocol_rtc.CommandCode.GetDefaultValue:
		// Return simulated default value (e.g., 0 for most parameters)
		value = 0
		ok = true
	default:
		// Unknown read command
		value = 0
		ok = false
	}

	if !ok {
		// Parameter not found, return 0
		value = 0
	}

	return value, nil
}

// processRTCWrite processes an RTC write request.
func (server *MockLinMot) processRTCWrite(upid uint16, value int32, cmdCode uint8) error {
	server.mu.Lock()
	defer server.mu.Unlock()

	// Handle monitoring channel configuration writes (0x20A8-0x20AB)
	if upid >= uint16(protocol_common.PUID.MonitoringChannel1UPID) && upid <= uint16(protocol_common.PUID.MonitoringChannel4UPID) {
		channelIndex := upid - uint16(protocol_common.PUID.MonitoringChannel1UPID)
		server.monitoringChannelConfig[channelIndex] = uint16(value)
	}

	switch cmdCode {
	case protocol_rtc.CommandCode.WriteROM:
		server.romParams[upid] = value
	case protocol_rtc.CommandCode.WriteRAM:
		server.ramParams[upid] = value
		server.updateStatusFromWrite(upid, value)
	case protocol_rtc.CommandCode.WriteRAMAndROM:
		// Write to both RAM and ROM
		server.ramParams[upid] = value
		server.romParams[upid] = value
		server.updateStatusFromWrite(upid, value)

	// UPID List Operations (0x20-0x23) - handled in handleRtcRequest special case
	// (removed from here to avoid fall-through)

	// Drive Operations (0x30-0x36)
	case protocol_rtc.CommandCode.RestartDrive:
		// Simulate drive restart (no-op in mock)
	case protocol_rtc.CommandCode.SetOSROMToDefault:
		// Simulate reset OS parameters (no-op in mock)
	case protocol_rtc.CommandCode.SetMCROMToDefault:
		// Simulate reset MC parameters (no-op in mock)
	case protocol_rtc.CommandCode.SetInterfaceROMToDefault:
		// Simulate reset interface parameters (no-op in mock)
	case protocol_rtc.CommandCode.SetApplicationROMToDefault:
		// Simulate reset application parameters (no-op in mock)

	// Curve Service (0x40-0x62)
	case protocol_rtc.CommandCode.SaveAllCurvesToFlash:
		// Simulate save curves (no-op in mock)
	case protocol_rtc.CommandCode.DeleteAllCurves:
		// Simulate delete curves (no-op in mock)
	case 0x50, 0x51, 0x52, 0x53, 0x54, 0x55: // Curve add/modify commands
		// Simulate curve operations (no-op in mock)
	case 0x60, 0x61, 0x62: // Curve get commands
		// Simulate curve get operations (no-op in mock)

	// Error Log (0x70-0x74)
	case 0x70, 0x71, 0x72, 0x73, 0x74:
		// Simulate error log operations (no-op in mock, return 0 errors)

	default:
		// Unknown command - return nil to avoid error
	}

	return nil
}

// updateStatusFromWrite updates drive status based on parameter writes.
func (server *MockLinMot) updateStatusFromWrite(upid uint16, value int32) {
	switch upid {
	case uint16(protocol_common.Parameter.Position1):
		// Position 1: update both demand and actual position (drive moves to Position 1)
		server.status.DemandPosition = value
		server.status.ActualPosition = value
	case uint16(protocol_common.Parameter.Position2):
		// Position 2: update both demand and actual position (drive moves to Position 2)
		// In sequential moves, this would be the second stage
		server.status.DemandPosition = value
		server.status.ActualPosition = value
	case uint16(protocol_common.PUID.RunMode):
		server.status.StatusWord = 0x0001
	}
}

// SetSimulateError enables or disables error simulation.
func (server *MockLinMot) SetSimulateError(simulate bool, code uint16) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.simulateError = simulate
	server.errorCode = code

	if simulate {
		server.status.ErrorCode = code
		server.mainState = protocol_control_word.State_Error
		server.status.StateVar = (uint16(server.mainState) << 8) | (server.status.StateVar & 0x00FF)
	} else {
		server.status.ErrorCode = 0
		if server.mainState == protocol_control_word.State_Error {
			server.mainState = protocol_control_word.State_SwitchOnDisabled
		}
		server.status.StateVar = (uint16(server.mainState) << 8) | (server.status.StateVar & 0x00FF)
	}
	server.updateStatusWord()
}

// SetSimulateWarning enables or disables warning simulation.
func (server *MockLinMot) SetSimulateWarning(simulate bool, code uint16) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.simulateWarning = simulate
	server.warningCode = code
}

// SetPersistentError enables or disables persistent error simulation (error won't clear on acknowledge).
func (server *MockLinMot) SetPersistentError(persistent bool) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.persistentError = persistent
}

// SetStatus sets the simulated drive status.
func (server *MockLinMot) SetStatus(status *protocol_common.Status) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.status = status
}

// SetStatusDelay sets an artificial delay for status responses (for testing timeouts).
func (server *MockLinMot) SetStatusDelay(delay time.Duration) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.statusDelay = delay
}

// SetErrorTextDelay sets an artificial delay for error text responses (for testing timeouts).
func (server *MockLinMot) SetErrorTextDelay(delay time.Duration) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.errorTextDelay = delay
}

// GetRAMParameter gets a RAM parameter value for testing.
func (server *MockLinMot) GetRAMParameter(upid uint16) int32 {
	server.mu.RLock()
	defer server.mu.RUnlock()
	return server.ramParams[upid]
}

// GetROMParameter gets a ROM parameter value for testing.
func (server *MockLinMot) GetROMParameter(upid uint16) int32 {
	server.mu.RLock()
	defer server.mu.RUnlock()
	return server.romParams[upid]
}

// Close stops the drive processing loop and closes the transport.
func (server *MockLinMot) Close() {
	if server.stopped.CompareAndSwap(false, true) {
		close(server.stopChan)
		server.transport.Close()
	}
}

// Stop is an alias for Close() for backwards compatibility.
func (server *MockLinMot) Stop() {
	server.Close()
}

// GetCommandTableEntry retrieves a CT entry for test verification.
func (server *MockLinMot) GetCommandTableEntry(id uint16) ([]byte, bool) {
	server.mu.RLock()
	defer server.mu.RUnlock()
	entry, exists := server.commandTable[id]
	if !exists {
		return nil, false
	}
	// Return copy to avoid race conditions
	entryCopy := make([]byte, len(entry))
	copy(entryCopy, entry)
	return entryCopy, true
}

// IsMotionControllerRunning returns the MC running state for test verification.
func (server *MockLinMot) IsMotionControllerRunning() bool {
	server.mu.RLock()
	defer server.mu.RUnlock()
	return server.mcRunning
}

// GetCommandTableSize returns the number of entries for test verification.
func (server *MockLinMot) GetCommandTableSize() int {
	server.mu.RLock()
	defer server.mu.RUnlock()
	return len(server.commandTable)
}

// processCTCommand processes command table operations.
// Returns (w3, w4, rtcStatus) where w3/w4 are packed into the response value.
// Enforces rule: CT commands must fail unless MainState is 1 or 2.
// StopMotionController and StartMotionController are exempt (they control MC state).
func (server *MockLinMot) processCTCommand(cmdCode uint8, w2, w3, w4 uint16) (uint16, uint16, uint8) {
	server.mu.Lock()
	defer server.mu.Unlock()

	switch cmdCode {
	case protocol_command_tables.CommandCode.StopMotionController:
		server.mcRunning = false
		// Transition to safe state (MainState 1 or 2) to allow CT operations
		// Set MainState to 2 (Ready To Switch On) - safe for CT operations
		server.status.StateVar = (server.status.StateVar & 0x00FF) | (0x02 << 8)
		return 0, 0, 0x00 // OK

	case protocol_command_tables.CommandCode.StartMotionController:
		server.mcRunning = true
		return 0, 0, 0x00 // OK
	}

	// For all other CT commands, check if MC is in a safe state (MainState 1, 2, 5, or 6)
	// MainState 1 = Switch On Disabled
	// MainState 2 = Ready To Switch On
	// MainState 5 = Switch On
	// MainState 6 = Ready to Operate
	mainState := (server.status.StateVar >> 8) & 0xFF
	if mainState != 1 && mainState != 2 && mainState != 5 && mainState != 6 {
		// CT commands fail if not in safe state (matches real hardware behavior)
		return 0, 0, 0xD1 // Size Error - actually state error, but using 0xD1 to match production
	}

	switch cmdCode {

	case protocol_command_tables.CommandCode.DeleteAllEntries:
		server.commandTable = make(map[uint16][]byte)
		server.ctWriteBuffers = make(map[uint16][]byte)
		server.ctReadPositions = make(map[uint16]int)
		return 0, 0, 0x00 // OK

	case protocol_command_tables.CommandCode.DeleteEntry:
		id := w2
		delete(server.commandTable, id)
		delete(server.ctWriteBuffers, id)
		delete(server.ctReadPositions, id)
		return 0, 0, 0x00 // OK

	case protocol_command_tables.CommandCode.AllocateEntry:
		id := w2
		// For AllocateEntry, size is in Word3 (low 16 bits), which is w4 in our extraction
		// value = (w3 << 16) | w4, but AllocateEntry puts size in low 16 bits (w4)
		size := w4
		// Validate size is even
		if size%2 == 1 {
			return 0, 0, 0xC0 // Error: odd size
		}
		// Real hardware behavior: If entry already exists and wasn't properly deleted,
		// attempting to allocate it again returns 0xD1 (Size Error).
		if existingEntry, exists := server.commandTable[id]; exists && existingEntry != nil {
			return 0, 0, 0xD1 // Size Error: entry already allocated
		}
		// Allocate entry buffer - store nil initially, will populate when write completes
		// Use write buffer capacity to track allocated size
		server.commandTable[id] = nil
		server.ctWriteBuffers[id] = make([]byte, 0, size)
		server.ctReadPositions[id] = 0
		if mockLinMotDebug {
			fmt.Printf("[MOCK] AllocateEntry id=%d size=%d\n", id, size)
		}
		return 0, 0, 0x00 // OK

	case protocol_command_tables.CommandCode.WriteEntryData:
		id := w2
		// Extract 4 bytes from w3 and w4 using little-endian wire order
		// Wire format: [b0, b1, b2, b3] where b0=w4_low, b1=w4_high, b2=w3_low, b3=w3_high
		b0 := byte(w4 & 0xFF) // w4 low
		b1 := byte(w4 >> 8)   // w4 high
		b2 := byte(w3 & 0xFF) // w3 low
		b3 := byte(w3 >> 8)   // w3 high
		data := []byte{b0, b1, b2, b3}

		// Append to write buffer
		buffer, exists := server.ctWriteBuffers[id]
		if !exists {
			return 0, 0, 0xC1 // Error: entry not allocated
		}
		buffer = append(buffer, data...)
		server.ctWriteBuffers[id] = buffer
		if mockLinMotDebug {
			fmt.Printf("[MOCK] WriteEntryData id=%d appended_bytes=%d cap=%d len=%d\n", id, len(data), cap(buffer), len(buffer))
		}

		// Check if complete (buffer length matches allocated capacity)
		allocatedSize := cap(buffer)
		if len(buffer) >= allocatedSize {
			// Copy complete buffer to entry
			fullEntry := make([]byte, allocatedSize)
			copy(fullEntry, buffer[:allocatedSize])
			server.commandTable[id] = fullEntry
			// Reset read position
			server.ctReadPositions[id] = 0
			if mockLinMotDebug {
				fmt.Printf("[MOCK] WriteEntryData id=%d COMPLETE allocatedSize=%d\n", id, allocatedSize)
			}
			return 0, 0, 0x00 // OK: complete
		}
		if mockLinMotDebug {
			fmt.Printf("[MOCK] WriteEntryData id=%d INCOMPLETE len=%d allocated=%d\n", id, len(buffer), allocatedSize)
		}
		return 0, 0, 0x04 // OK: incomplete

	case protocol_command_tables.CommandCode.GetEntrySize:
		id := w2
		entry, exists := server.commandTable[id]
		if !exists || entry == nil {
			return 0, 0, 0xC2 // Error: entry not found
		}
		size := uint16(len(entry))
		// Return size in low 16 bits (w4) so caller that extracts w4 sees the size
		return 0, size, 0x00 // Return size in w4

	case protocol_command_tables.CommandCode.ReadEntryData:
		id := w2
		entry, exists := server.commandTable[id]
		if !exists {
			return 0, 0, 0xC2 // Error: entry not found
		}

		// Get current read position
		pos := server.ctReadPositions[id]

		// Read next 4 bytes from entry
		if pos >= len(entry) {
			// Already read everything
			return 0, 0, 0x00 // Complete, no more data
		}

		// Read up to 4 bytes
		end := pos + 4
		if end > len(entry) {
			end = len(entry)
		}

		readData := entry[pos:end]
		// Pad to 4 bytes if needed
		for len(readData) < 4 {
			readData = append(readData, 0)
		}

		// Pack for little-endian wire order: [b0, b1, b2, b3]
		// w4 low=b0, w4 high=b1, w3 low=b2, w3 high=b3
		w4Out := uint16(readData[1])<<8 | uint16(readData[0])
		w3Out := uint16(readData[3])<<8 | uint16(readData[2])

		// Update read position
		server.ctReadPositions[id] = end

		if end >= len(entry) {
			// Reset for next read cycle
			server.ctReadPositions[id] = 0
			return w3Out, w4Out, 0x00 // Complete
		}
		return w3Out, w4Out, 0x04 // Incomplete

	case protocol_command_tables.CommandCode.PresenceMask0, protocol_command_tables.CommandCode.PresenceMask1,
		protocol_command_tables.CommandCode.PresenceMask2, protocol_command_tables.CommandCode.PresenceMask3,
		protocol_command_tables.CommandCode.PresenceMask4, protocol_command_tables.CommandCode.PresenceMask5,
		protocol_command_tables.CommandCode.PresenceMask6, protocol_command_tables.CommandCode.PresenceMask7:
		// Calculate mask index (0-7)
		maskIndex := cmdCode - protocol_command_tables.CommandCode.PresenceMask0
		baseID := uint16(maskIndex * 32)

		// Build 32-bit mask with inverted logic (0 = present, 1 = absent)
		// Initialize to all 1s (all entries absent)
		mask := uint32(0xFFFFFFFF)
		for bit := uint(0); bit < 32; bit++ {
			entryID := baseID + uint16(bit)
			entry, exists := server.commandTable[entryID]
			if exists && entry != nil {
				// Entry present - clear bit to 0
				mask &^= (1 << bit)
			}
			// Entry absent - bit remains 1
		}

		// Return mask in w3 (high 16 bits) and w4 (low 16 bits)
		w3 := uint16(mask >> 16)
		w4 := uint16(mask & 0xFFFF)
		return w3, w4, 0x00 // OK

	default:
		return 0, 0, 0xFF // Unknown command
	}
}

// handleControlWordRequest processes a control word request and updates state machine
func (server *MockLinMot) handleControlWordRequest(request *protocol_control_word.ControlWordRequest) protocol_common.Response {
	server.mu.Lock()
	defer server.mu.Unlock()

	controlWord := request.GetControlWord()

	// Extract control word bits
	switchOn := protocol_control_word.IsBitSet(controlWord, protocol_control_word.ControlWordBit_SwitchOn)
	enableVoltage := protocol_control_word.IsBitSet(controlWord, protocol_control_word.ControlWordBit_EnableVoltage)
	quickStopReleased := protocol_control_word.IsBitSet(controlWord, protocol_control_word.ControlWordBit_QuickStop)
	enableOperation := protocol_control_word.IsBitSet(controlWord, protocol_control_word.ControlWordBit_EnableOperation)
	errorAck := protocol_control_word.IsBitSet(controlWord, protocol_control_word.ControlWordBit_ErrorAcknowledge)
	home := protocol_control_word.IsBitSet(controlWord, protocol_control_word.ControlWordBit_Home)

	// Simple state machine simulation
	// This is a simplified version - real drives have more complex transitions

	// Handle error acknowledge (rising edge detection)
	if errorAck && !server.errorAcknowledge && server.mainState == protocol_control_word.State_Error {
		if !server.persistentError {
			// Rising edge on error acknowledge - clear error
			server.status.ErrorCode = 0
			server.simulateError = false
			server.mainState = protocol_control_word.State_SwitchOnDisabled
		}
	}
	server.errorAcknowledge = errorAck

	// Handle state transitions based on control word
	// Allow multi-step transitions in a single request for testing convenience
	transitionCount := 0
	maxTransitions := 10 // Prevent infinite loops

	for transitionCount < maxTransitions {
		oldState := server.mainState

		switch server.mainState {
		case protocol_control_word.State_SwitchOnDisabled:
			if switchOn && enableVoltage {
				server.mainState = protocol_control_word.State_ReadyToSwitchOn
			}

		case protocol_control_word.State_ReadyToSwitchOn:
			if !switchOn {
				server.mainState = protocol_control_word.State_SwitchOnDisabled
			} else if switchOn && enableVoltage && quickStopReleased {
				server.mainState = protocol_control_word.State_SwitchOn
			}

		case protocol_control_word.State_SwitchOn:
			if !switchOn {
				server.mainState = protocol_control_word.State_SwitchOnDisabled
			} else if enableOperation {
				server.mainState = protocol_control_word.State_OperationEnabled
			}

		case protocol_control_word.State_OperationEnabled:
			if !switchOn {
				server.mainState = protocol_control_word.State_SwitchOnDisabled
			} else if !enableOperation {
				server.mainState = protocol_control_word.State_SwitchOn
			} else if home {
				server.mainState = protocol_control_word.State_Homing
			}

		case protocol_control_word.State_Homing:
			// Simulate instant homing completion for testing
			server.mainState = protocol_control_word.State_OperationEnabled
			// Mark as homed
			server.isHomed = true

		case protocol_control_word.State_Error:
			// Stay in error state until acknowledged
			// (handled above with error acknowledge logic)
		}

		// If state didn't change, we're done
		if server.mainState == oldState {
			break
		}
		transitionCount++
	}

	// Update StatusWord based on current state
	server.updateStatusWord()

	// Update StateVar (MainState in high byte)
	server.status.StateVar = (uint16(server.mainState) << 8) | (server.status.StateVar & 0x00FF)

	// Return a custom raw status response and set reqBits bit 0 so RequestManager
	// recognizes it as a control-word response instead of a plain status packet.
	// The mock needs this explicit shape because RequestManager routes these paths
	// based on reqBits, not just on the response payload contents.
	statusPkt, _ := protocol_common.NewStatusResponse(server.status).WritePacket()
	if len(statusPkt) >= 8 {
		// Set ControlWord bit in reqBits
		statusPkt[0] |= 0x01
	}

	return &customRawResponse{data: statusPkt}
}

// customRawResponse wraps a raw byte slice to implement the Response interface.
type customRawResponse struct {
	data []byte
}

func (r *customRawResponse) WritePacket() ([]byte, error) {
	return r.data, nil
}

// updateStatusWord sets status word bits based on current state
func (server *MockLinMot) updateStatusWord() {
	// Clear all state-dependent bits first
	server.status.StatusWord = 0x0000

	// Set persistent bits
	if server.isHomed {
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_Homed)
	}

	// Set bits based on current state
	switch server.mainState {
	case protocol_control_word.State_SwitchOnDisabled:
		// No special bits

	case protocol_control_word.State_ReadyToSwitchOn:
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_SwitchOnActive)

	case protocol_control_word.State_SwitchOn:
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_SwitchOnActive)
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_VoltageEnable)
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_QuickStop)

	case protocol_control_word.State_OperationEnabled:
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_OperationEnabled)
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_SwitchOnActive)
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_EnableOperation)
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_VoltageEnable)
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_QuickStop)

	case protocol_control_word.State_Homing:
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_OperationEnabled)
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_SwitchOnActive)
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_MotionActive)

	case protocol_control_word.State_Error:
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_Error)
		if server.simulateError {
			// Check if it's a fatal error (just use error code > 0x8000 as criterion)
			if server.errorCode >= 0x8000 {
				server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_FatalError)
			}
		}
	}

	// Apply warnings if simulated
	if server.simulateWarning {
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_Warning)
	}

	// Apply errors if simulated
	if server.simulateError {
		server.status.StatusWord = protocol_control_word.SetBit(server.status.StatusWord, protocol_control_word.StatusWordBit_Error)
		// Transition to error state
		if server.mainState != protocol_control_word.State_Error {
			server.mainState = protocol_control_word.State_Error
		}
	}
}

// SimulateError sets the mock to simulate an error condition
func (server *MockLinMot) SimulateError(errorCode uint16) {
	server.mu.Lock()
	defer server.mu.Unlock()

	server.simulateError = true
	server.errorCode = errorCode
	server.status.ErrorCode = errorCode
	server.mainState = protocol_control_word.State_Error
	server.status.StateVar = (uint16(server.mainState) << 8) | (server.status.StateVar & 0x00FF)
	server.updateStatusWord()
}
