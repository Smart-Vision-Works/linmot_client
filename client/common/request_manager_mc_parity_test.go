package client_common

import (
	"context"
	"encoding/binary"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"
	protocol_motion_control "github.com/Smart-Vision-Works/linmot_client/protocol/motion_control"
	transport "github.com/Smart-Vision-Works/linmot_client/transport"
)

// mcTrackingTransport wraps a transport.Client and records all SendPacket calls
// for testing MC counter selection and header encoding.
type mcTrackingTransport struct {
	transport.Client
	mu          sync.Mutex
	sentPackets [][]byte
	sendCount   atomic.Int32
}

func newMCTrackingTransport(base transport.Client) *mcTrackingTransport {
	return &mcTrackingTransport{
		Client:      base,
		sentPackets: make([][]byte, 0),
	}
}

func (tt *mcTrackingTransport) SendPacket(ctx context.Context, data []byte) error {
	tt.sendCount.Add(1)
	tt.mu.Lock()
	// Make a copy of the packet for later inspection
	packetCopy := make([]byte, len(data))
	copy(packetCopy, data)
	tt.sentPackets = append(tt.sentPackets, packetCopy)
	tt.mu.Unlock()
	return tt.Client.SendPacket(ctx, data)
}

func (tt *mcTrackingTransport) getLastSentPacket() []byte {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	if len(tt.sentPackets) == 0 {
		return nil
	}
	last := tt.sentPackets[len(tt.sentPackets)-1]
	result := make([]byte, len(last))
	copy(result, last)
	return result
}

// TestMCCountNibble_UsesLastSeenStateVarLowPlusOne verifies that MC counter
// selection uses (lastSeenNibble + 1) wrap 1-4, matching linudp.cs behavior.
func TestMCCountNibble_UsesLastSeenStateVarLowPlusOne(t *testing.T) {
	tests := []struct {
		name           string
		stateVarLowNib uint8 // Last seen StateVarLow nibble (0-4)
		expectedCount  uint8 // Expected MC counter (1-4)
	}{
		{"nibble_0_to_1", 0, 1}, // No StateVar seen yet -> default to 1
		{"nibble_1_to_2", 1, 2},
		{"nibble_2_to_3", 2, 3},
		{"nibble_3_to_4", 3, 4},
		{"nibble_4_wrap_to_1", 4, 1}, // Wrap: 4+1=5 -> 1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock transport
			baseClient, server := transport.NewMockTransportClientWithServer()
			trackingClient := newMCTrackingTransport(baseClient)

			// Create RequestManager
			rm := NewRequestManager(trackingClient, 10*time.Millisecond)
			rm.Start()
			defer rm.Stop()

			// Inject a StatusResponse packet with the specified StateVarLow nibble
			// This will update lastStateVarLowNibble
			stateVar := uint16(0x5670 | uint16(tt.stateVarLowNib)) // High byte 0x56, low byte 0x70|nibble
			status := &protocol_common.Status{
				ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar,
				StatusWord:   0x0001,
				StateVar:     stateVar,
			}
			statusResp := protocol_common.NewStatusResponse(status)
			statusPacket, err := statusResp.WritePacket()
			if err != nil {
				t.Fatalf("WritePacket() error = %v", err)
			}
			server.SendPacket(statusPacket)

			// Wait for packet to be processed (rxLoop should parse it and update lastStateVarLowNibble)
			time.Sleep(100 * time.Millisecond)

			// Send an MC request
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			mcReq := protocol_motion_control.NewMCCommandRequest(
				protocol_motion_control.MasterIDs.InterfaceControl,
				0x01,             // SubID
				[]uint16{0x1234}, // Parameters
			)

			// Send request (will use nextMCCountFromLastStateVar)
			var wg sync.WaitGroup
			var mcResp protocol_motion_control.MCResponse
			var mcErr error

			wg.Add(1)
			go func() {
				defer wg.Done()
				mcResp, mcErr = SendRequestAndReceive[protocol_motion_control.MCResponse](rm, ctx, mcReq)
			}()

			// Wait a bit for the request to be sent
			time.Sleep(50 * time.Millisecond)

			// Get the sent packet and verify MC header counter
			sentPacket := trackingClient.getLastSentPacket()
			if sentPacket == nil {
				t.Fatal("No packet was sent")
			}

			if len(sentPacket) < 10 {
				t.Fatalf("Packet too short: got %d bytes, need at least 10", len(sentPacket))
			}

			// Extract MC header from bytes 8-9 (little-endian)
			mcHeaderWord := binary.LittleEndian.Uint16(sentPacket[8:10])
			// Decode header: low byte = (SubID << 4) | Counter
			lowByte := uint8(mcHeaderWord & 0xFF)
			actualCounter := lowByte & 0x0F // Low nibble is the counter

			if actualCounter != tt.expectedCount {
				t.Errorf("MC counter mismatch: got %d, expected %d (lastSeenNibble=%d)",
					actualCounter, tt.expectedCount, tt.stateVarLowNib)
			}

			// Send a response to allow the request to complete
			mcResponseStatus := &protocol_common.Status{
				ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar | protocol_common.RespBitDemandPosition,
				StatusWord:   0x0001,
				StateVar:     uint16(0x5670 | uint16(tt.expectedCount)), // Echo the counter we sent
			}
			mcResponse := protocol_motion_control.NewMCCommandResponse(mcResponseStatus, tt.expectedCount)
			mcResponsePacket, err := mcResponse.WritePacket()
			if err != nil {
				t.Fatalf("MC response WritePacket() error = %v", err)
			}
			server.SendPacket(mcResponsePacket)

			// Wait for request to complete
			wg.Wait()

			if mcErr != nil {
				t.Errorf("MC request failed: %v", mcErr)
			}
			if mcResp == nil {
				t.Error("MC response is nil")
			}
		})
	}
}

// TestMCHeaderEncoding_SubByteORsNibble verifies that MC header encoding
// correctly ORs the counter nibble with SubID: (SubID << 4) | (counterNibble & 0x0F).
func TestMCHeaderEncoding_SubByteORsNibble(t *testing.T) {
	// Create mock transport
	baseClient, server := transport.NewMockTransportClientWithServer()
	trackingClient := newMCTrackingTransport(baseClient)

	// Create RequestManager
	rm := NewRequestManager(trackingClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Set lastStateVarLowNibble to 2 (so next counter will be 3)
	stateVar := uint16(0x5670 | uint16(2))
	status := &protocol_common.Status{
		ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar,
		StatusWord:   0x0001,
		StateVar:     stateVar,
	}
	statusResp := protocol_common.NewStatusResponse(status)
	statusPacket, _ := statusResp.WritePacket()
	server.SendPacket(statusPacket)

	// Wait for packet to be processed
	time.Sleep(100 * time.Millisecond)

	// Send MC request with SubID = 0x05
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mcReq := protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.InterfaceControl,
		0x05, // SubID
		[]uint16{0x1234},
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = SendRequestAndReceive[protocol_motion_control.MCResponse](rm, ctx, mcReq)
	}()

	// Wait for packet to be sent
	time.Sleep(50 * time.Millisecond)

	// Get sent packet
	sentPacket := trackingClient.getLastSentPacket()
	if sentPacket == nil {
		t.Fatal("No packet was sent")
	}

	if len(sentPacket) < 10 {
		t.Fatalf("Packet too short: got %d bytes, need at least 10", len(sentPacket))
	}

	// Extract MC header from bytes 8-9
	mcHeaderWord := binary.LittleEndian.Uint16(sentPacket[8:10])
	lowByte := uint8(mcHeaderWord & 0xFF)
	highByte := uint8((mcHeaderWord >> 8) & 0xFF)

	// Verify encoding: lowByte should be (SubID << 4) | Counter
	expectedSubID := uint8(0x05)
	expectedCounter := uint8(3) // (lastSeenNibble=2) + 1 = 3
	expectedLowByte := (expectedSubID << 4) | (expectedCounter & 0x0F)

	if lowByte != expectedLowByte {
		t.Errorf("MC header low byte mismatch: got 0x%02X, expected 0x%02X (SubID=0x%02X, Counter=%d)",
			lowByte, expectedLowByte, expectedSubID, expectedCounter)
	}

	// Verify high byte is MasterID
	expectedMasterID := uint8(protocol_motion_control.MasterIDs.InterfaceControl)
	if highByte != expectedMasterID {
		t.Errorf("MC header high byte (MasterID) mismatch: got 0x%02X, expected 0x%02X",
			highByte, expectedMasterID)
	}

	// Send response to allow request to complete
	mcResponseStatus := &protocol_common.Status{
		ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar | protocol_common.RespBitDemandPosition,
		StatusWord:   0x0001,
		StateVar:     uint16(0x5670 | uint16(expectedCounter)),
	}
	mcResponse := protocol_motion_control.NewMCCommandResponse(mcResponseStatus, expectedCounter)
	mcResponsePacket, _ := mcResponse.WritePacket()
	server.SendPacket(mcResponsePacket)

	wg.Wait()
}

// TestMCResponseRouting_UsesEchoedNibble verifies that MC response routing
// uses the echoed counter from StateVarLow for matching pending requests.
func TestMCResponseRouting_UsesEchoedNibble(t *testing.T) {
	// Create mock transport
	baseClient, server := transport.NewMockTransportClientWithServer()
	trackingClient := newMCTrackingTransport(baseClient)

	// Create RequestManager
	rm := NewRequestManager(trackingClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Set lastStateVarLowNibble to 1 (inject StatusResponse packet)
	stateVar := uint16(0x5670 | uint16(1))
	status := &protocol_common.Status{
		ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar,
		StatusWord:   0x0001,
		StateVar:     stateVar,
	}
	statusResp := protocol_common.NewStatusResponse(status)
	statusPacket, _ := statusResp.WritePacket()
	server.SendPacket(statusPacket)

	// Wait for packet to be processed
	time.Sleep(100 * time.Millisecond)

	// Send MC request (should use counter=2, since lastSeenNibble=1)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mcReq := protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.InterfaceControl,
		0x01,
		[]uint16{0x1234},
	)

	var wg sync.WaitGroup
	var mcResp protocol_motion_control.MCResponse
	var mcErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		mcResp, mcErr = SendRequestAndReceive[protocol_motion_control.MCResponse](rm, ctx, mcReq)
	}()

	// Wait for request to be sent
	time.Sleep(50 * time.Millisecond)

	// Verify request was sent with counter=2
	sentPacket := trackingClient.getLastSentPacket()
	if sentPacket == nil {
		t.Fatal("No packet was sent")
	}

	mcHeaderWord := binary.LittleEndian.Uint16(sentPacket[8:10])
	lowByte := uint8(mcHeaderWord & 0xFF)
	sentCounter := lowByte & 0x0F

	expectedCounter := uint8(2)
	if sentCounter != expectedCounter {
		t.Errorf("Sent counter mismatch: got %d, expected %d", sentCounter, expectedCounter)
	}

	// Inject MC response packet with StateVarLow nibble=2 (echoing the counter we sent)
	mcResponseStatus := &protocol_common.Status{
		ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar | protocol_common.RespBitDemandPosition,
		StatusWord:   0x0001,
		StateVar:     uint16(0x5670 | uint16(expectedCounter)), // Echo counter=2 in StateVarLow
	}
	mcResponse := protocol_motion_control.NewMCCommandResponse(mcResponseStatus, expectedCounter)
	mcResponsePacket, err := mcResponse.WritePacket()
	if err != nil {
		t.Fatalf("MC response WritePacket() error = %v", err)
	}
	server.SendPacket(mcResponsePacket)

	// Wait for request to complete
	wg.Wait()

	// Verify response was delivered
	if mcErr != nil {
		t.Errorf("MC request failed: %v", mcErr)
	}
	if mcResp == nil {
		t.Fatal("MC response is nil - routing failed")
	}

	// Verify the response has the correct counter
	if mcResp.MCCounter() != expectedCounter {
		t.Errorf("MC response counter mismatch: got %d, expected %d",
			mcResp.MCCounter(), expectedCounter)
	}
}

// TestMCCountSelectedAtSendTime_SeesLatestStateVar verifies that MC counter
// selection happens at send-time (startRequest), not queue-time (submitRequest),
// ensuring it uses the latest StateVarLow value even if it changes while queued.
func TestMCCountSelectedAtSendTime_SeesLatestStateVar(t *testing.T) {
	// Create mock transport
	baseClient, server := transport.NewMockTransportClientWithServer()
	trackingClient := newMCTrackingTransport(baseClient)

	// Create RequestManager
	rm := NewRequestManager(trackingClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Preload lastStateVarLowNibble to 1 (so next should be 2)
	stateVar1 := uint16(0x5670 | uint16(1))
	status1 := &protocol_common.Status{
		ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar,
		StatusWord:   0x0001,
		StateVar:     stateVar1,
	}
	statusResp1 := protocol_common.NewStatusResponse(status1)
	statusPacket1, _ := statusResp1.WritePacket()
	server.SendPacket(statusPacket1)

	// Wait for rxLoop to process the packet
	time.Sleep(100 * time.Millisecond)

	// Setup hook coordination channels
	hookHit := make(chan struct{})
	unblockSend := make(chan struct{})

	// Install beforeFirstSend hook that blocks until unblockSend is closed
	rm.beforeFirstSend = func(req *pendingRequest) {
		close(hookHit) // Signal that we've reached the hook
		<-unblockSend  // Block until unblocked
	}

	// Launch MC request in goroutine
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mcReq := protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.InterfaceControl,
		0x01,
		[]uint16{0x1234},
	)

	var wg sync.WaitGroup
	var mcResp protocol_motion_control.MCResponse
	var mcErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		mcResp, mcErr = SendRequestAndReceive[protocol_motion_control.MCResponse](rm, ctx, mcReq)
	}()

	// Wait for hook to be hit (request is queued and startRequest has reached the hook)
	<-hookHit

	// While still blocked in beforeFirstSend, inject a NEW StatusResponse
	// with StateVarLow nibble = 3 (so next should be 4)
	stateVar3 := uint16(0x5670 | uint16(3))
	status3 := &protocol_common.Status{
		ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar,
		StatusWord:   0x0001,
		StateVar:     stateVar3,
	}
	statusResp3 := protocol_common.NewStatusResponse(status3)
	statusPacket3, _ := statusResp3.WritePacket()
	server.SendPacket(statusPacket3)

	// Wait briefly for rxLoop to process the new packet and update lastStateVarLowNibble
	time.Sleep(100 * time.Millisecond)

	// Now unblock the hook so startRequest proceeds to assign mcCounter and send
	close(unblockSend)

	// Wait for packet to be sent
	time.Sleep(50 * time.Millisecond)

	// Clear the hook
	rm.beforeFirstSend = nil

	// Get the sent packet and verify MC header counter
	sentPacket := trackingClient.getLastSentPacket()
	if sentPacket == nil {
		t.Fatal("No packet was sent")
	}

	if len(sentPacket) < 10 {
		t.Fatalf("Packet too short: got %d bytes, need at least 10", len(sentPacket))
	}

	// Extract MC header from bytes 8-9 (little-endian)
	mcHeaderWord := binary.LittleEndian.Uint16(sentPacket[8:10])
	lowByte := uint8(mcHeaderWord & 0xFF)
	actualCounter := lowByte & 0x0F // Low nibble is the counter

	// Assert: counter should be 4 (derived from latest nibble 3, not initial nibble 1)
	expectedCounter := uint8(4)
	if actualCounter != expectedCounter {
		t.Errorf("MC counter mismatch: got %d, expected %d (should use latest StateVarLow nibble=3, not initial nibble=1)",
			actualCounter, expectedCounter)
	}

	// Send MC response packet that echoes counter=4 in StateVarLow so the request completes
	mcResponseStatus := &protocol_common.Status{
		ResponseBits: protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar | protocol_common.RespBitDemandPosition,
		StatusWord:   0x0001,
		StateVar:     uint16(0x5670 | uint16(expectedCounter)),
	}
	mcResponse := protocol_motion_control.NewMCCommandResponse(mcResponseStatus, expectedCounter)
	mcResponsePacket, err := mcResponse.WritePacket()
	if err != nil {
		t.Fatalf("MC response WritePacket() error = %v", err)
	}
	server.SendPacket(mcResponsePacket)

	// Wait for request to complete
	wg.Wait()

	// Verify no error
	if mcErr != nil {
		t.Errorf("MC request failed: %v", mcErr)
	}
	if mcResp == nil {
		t.Error("MC response is nil")
	}
}
