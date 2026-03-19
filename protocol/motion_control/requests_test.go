package protocol_motion_control

import (
	"encoding/binary"
	"testing"

	protocol_common "gsail-go/linmot/protocol/common"
)

// TestMCCommandRequest_WritePacket tests the basic MC command request packet structure
func TestMCCommandRequest_WritePacket(t *testing.T) {
	// Create a simple MC request with known parameters
	params := [15]uint16{
		0x1234, 0x5678, 0x9ABC, 0xDEF0,
		0x1111, 0x2222, 0x3333, 0x4444,
		0x5555, 0x6666, 0x7777, 0x8888,
		0x9999, 0xAAAA, 0xBBBB,
	}

	header := MCHeader{
		MasterID: 0x01,
		SubID:    0x00,
		Counter:  3,
	}

	request := &MCCommandRequest{
		header:     header,
		parameters: params,
	}

	packet, err := request.WritePacket()
	if err != nil {
		t.Fatalf("WritePacket() error = %v", err)
	}

	// Verify packet length (8-byte header + 32-byte MC data = 40 bytes)
	if len(packet) != 40 {
		t.Errorf("expected packet length 40, got %d", len(packet))
	}

	// Verify request definition (bit 1: Motion Control)
	reqBits := binary.LittleEndian.Uint32(packet[0:4])
	if reqBits != protocol_common.RequestFlags.MotionControl {
		t.Errorf("expected request bits 0x%08X, got 0x%08X",
			protocol_common.RequestFlags.MotionControl, reqBits)
	}

	// Verify response definition (bits 0-1: StatusWord + StateVar, bit 3: DemandPosition)
	repBits := binary.LittleEndian.Uint32(packet[4:8])
	expectedRepBits := StateVarResponse | protocol_common.RespBitDemandPosition
	if repBits != expectedRepBits {
		t.Errorf("expected response bits 0x%08X, got 0x%08X",
			expectedRepBits, repBits)
	}

	// Verify MC header (SubID in high nibble, counter in low nibble)
	subIDCounterByte := packet[8]
	masterID := MasterID(packet[9])

	extractedSubID := (subIDCounterByte >> 4) & 0x0F
	extractedCounter := subIDCounterByte & 0x0F

	if extractedSubID != header.SubID {
		t.Errorf("expected SubID 0x%02X, got 0x%02X", header.SubID, extractedSubID)
	}
	if extractedCounter != header.Counter {
		t.Errorf("expected Counter %d, got %d", header.Counter, extractedCounter)
	}
	if masterID != header.MasterID {
		t.Errorf("expected MasterID 0x%02X, got 0x%02X", header.MasterID, masterID)
	}

	// Verify first parameter (little-endian)
	param1 := binary.LittleEndian.Uint16(packet[10:12])
	if param1 != params[0] {
		t.Errorf("expected param1 0x%04X, got 0x%04X", params[0], param1)
	}
}

// TestMCCommandRequest_Counter tests MC counter embedding in SubID
func TestMCCommandRequest_Counter(t *testing.T) {
	tests := []struct {
		name    string
		counter uint8
		subID   uint8
	}{
		{"Counter 1", 1, 0x00},
		{"Counter 2", 2, 0x05},
		{"Counter 3", 3, 0x0A},
		{"Counter 4", 4, 0x0F},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := MCHeader{
				MasterID: 0x01,
				SubID:    tt.subID,
				Counter:  tt.counter,
			}

			request := &MCCommandRequest{
				header:     header,
				parameters: [15]uint16{},
			}

			packet, err := request.WritePacket()
			if err != nil {
				t.Fatalf("WritePacket() error = %v", err)
			}

			// Verify counter is embedded in low nibble, SubID in high nibble
			subIDByte := packet[8]
			extractedSubID := (subIDByte >> 4) & 0x0F
			extractedCounter := subIDByte & 0x0F

			if extractedCounter != tt.counter {
				t.Errorf("expected counter %d, got %d", tt.counter, extractedCounter)
			}
			if extractedSubID != tt.subID {
				t.Errorf("expected subID 0x%X, got 0x%X", tt.subID, extractedSubID)
			}
		})
	}
}

// TestMCCommandRequest_SetCounter tests the SetCounter method
func TestMCCommandRequest_SetCounter(t *testing.T) {
	header := MCHeader{
		MasterID: 0x01,
		SubID:    0x00,
		Counter:  1,
	}

	request := &MCCommandRequest{
		header:     header,
		parameters: [15]uint16{},
	}

	// Set counter to 3
	request.SetCounter(3)

	if request.header.Counter != 3 {
		t.Errorf("expected counter 3, got %d", request.header.Counter)
	}

	// Verify it appears in packet
	packet, err := request.WritePacket()
	if err != nil {
		t.Fatalf("WritePacket() error = %v", err)
	}

	subIDByte := packet[8]
	extractedCounter := subIDByte & 0x0F
	if extractedCounter != 3 {
		t.Errorf("expected counter 3 in packet, got %d", extractedCounter)
	}
}

// TestMCCommandRequest_Header tests the Header() method
func TestMCCommandRequest_Header(t *testing.T) {
	expectedHeader := MCHeader{
		MasterID: 0x05,
		SubID:    0x0A,
		Counter:  2,
	}

	request := &MCCommandRequest{
		header:     expectedHeader,
		parameters: [15]uint16{},
	}

	gotHeader := request.Header()

	if gotHeader.MasterID != expectedHeader.MasterID {
		t.Errorf("expected MasterID 0x%02X, got 0x%02X",
			expectedHeader.MasterID, gotHeader.MasterID)
	}
	if gotHeader.SubID != expectedHeader.SubID {
		t.Errorf("expected SubID 0x%02X, got 0x%02X",
			expectedHeader.SubID, gotHeader.SubID)
	}
	if gotHeader.Counter != expectedHeader.Counter {
		t.Errorf("expected Counter %d, got %d",
			expectedHeader.Counter, gotHeader.Counter)
	}
}

// TestReadMCRequest tests parsing an MC request packet
func TestReadMCRequest(t *testing.T) {
	// Create a packet manually
	packet := make([]byte, 40)

	// Set request definition (Motion Control bit)
	binary.LittleEndian.PutUint32(packet[0:4], protocol_common.RequestFlags.MotionControl)

	// Set response definition (StateVar bit)
	binary.LittleEndian.PutUint32(packet[4:8], StateVarResponse)

	// Set MC header (SubID 0x05, counter 2)
	packet[8] = (0x05 << 4) | 2 // (SubID << 4) | Counter
	packet[9] = 0x01            // MasterID

	// Set first parameter
	binary.LittleEndian.PutUint16(packet[10:12], 0xABCD)

	// Parse packet
	request, err := ReadMCRequest(packet)
	if err != nil {
		t.Fatalf("ReadMCRequest() error = %v", err)
	}

	// Verify header
	header := request.Header()
	if header.MasterID != 0x01 {
		t.Errorf("expected MasterID 0x01, got 0x%02X", header.MasterID)
	}
	if header.SubID != 0x05 {
		t.Errorf("expected SubID 0x05, got 0x%02X", header.SubID)
	}
	if header.Counter != 2 {
		t.Errorf("expected Counter 2, got %d", header.Counter)
	}
}

// TestReadMCRequest_InvalidLength tests error handling for invalid packet length
func TestReadMCRequest_InvalidLength(t *testing.T) {
	// Too short packet
	packet := make([]byte, 30)

	_, err := ReadMCRequest(packet)
	if err == nil {
		t.Error("expected error for short packet, got nil")
	}
}
