package protocol_motion_control

import (
	"encoding/binary"
	"testing"

	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
)

// TestMCCommandResponse_WritePacket tests MC response packet structure
func TestMCCommandResponse_WritePacket(t *testing.T) {
	status := &protocol_common.Status{
		StatusWord:     0x1234,
		StateVar:       0x5678,
		ActualPosition: 100000,
		DemandPosition: 200000,
		Current:        500,
		WarnWord:       0x0000,
		ErrorCode:      0x0000,
	}
	mcCounter := uint8(3)

	response := NewMCCommandResponse(status, mcCounter)
	packet, err := response.WritePacket()
	if err != nil {
		t.Fatalf("WritePacket() error = %v", err)
	}

	// Verify packet length (16 bytes: header 8 + StatusWord 2 + StateVar 2 + DemandPosition 4)
	if len(packet) != 16 {
		t.Errorf("expected packet length 16, got %d", len(packet))
	}

	// Verify request definition (should echo MotionControl bit per LinUDP V2 spec 4.4.2)
	reqBits := binary.LittleEndian.Uint32(packet[0:4])
	expectedReqBits := protocol_common.RequestFlags.MotionControl
	if reqBits != expectedReqBits {
		t.Errorf("expected request bits 0x%08X, got 0x%08X", expectedReqBits, reqBits)
	}

	// Verify response definition (bits 0-1: StatusWord + StateVar, bit 3: DemandPosition)
	repBits := binary.LittleEndian.Uint32(packet[4:8])
	expectedRepBits := StateVarResponse | protocol_common.RespBitDemandPosition
	if repBits != expectedRepBits {
		t.Errorf("expected response bits 0x%08X, got 0x%08X", expectedRepBits, repBits)
	}

	// Verify StatusWord (bit 0)
	statusWord := binary.LittleEndian.Uint16(packet[8:10])
	if statusWord != status.StatusWord {
		t.Errorf("expected StatusWord 0x%04X, got 0x%04X", status.StatusWord, statusWord)
	}

	// Verify StateVar with MC counter (bit 1)
	stateVar := binary.LittleEndian.Uint16(packet[10:12])
	expectedStateVar := (status.StateVar & 0xFFF0) | uint16(mcCounter&0x0F)
	if stateVar != expectedStateVar {
		t.Errorf("expected StateVar 0x%04X (with counter %d), got 0x%04X", expectedStateVar, mcCounter, stateVar)
	}

	// Verify DemandPosition (bit 3)
	demandPos := int32(binary.LittleEndian.Uint32(packet[12:16]))
	if demandPos != status.DemandPosition {
		t.Errorf("expected DemandPosition %d, got %d", status.DemandPosition, demandPos)
	}

	// Verify MC counter is in low nibble
	extractedCounter := uint8(stateVar & 0x0F)
	if extractedCounter != mcCounter {
		t.Errorf("expected MC counter %d, got %d", mcCounter, extractedCounter)
	}
}

// TestReadMCResponse tests parsing MC response with only StatusWord + StateVar bits
func TestReadMCResponse(t *testing.T) {
	// Create minimal MC response packet (12 bytes)
	packet := make([]byte, 12)

	// Header: request definition = 0, response definition = StatusWord + StateVar
	binary.LittleEndian.PutUint32(packet[0:4], 0x00000000)
	binary.LittleEndian.PutUint32(packet[4:8], StateVarResponse)

	// StatusWord (bit 0)
	binary.LittleEndian.PutUint16(packet[8:10], 0xABCD)

	// StateVar (bit 1) with MC counter 2 in low nibble
	stateVar := uint16(0x5670 | 2) // High nibble 0x567, counter 2 in low nibble
	binary.LittleEndian.PutUint16(packet[10:12], stateVar)

	// Parse packet
	response, err := ReadMCResponse(packet)
	if err != nil {
		t.Fatalf("ReadMCResponse() error = %v", err)
	}

	// Verify StatusWord
	if response.Status().StatusWord != 0xABCD {
		t.Errorf("expected StatusWord 0xABCD, got 0x%04X", response.Status().StatusWord)
	}

	// Verify StateVar
	if response.Status().StateVar != stateVar {
		t.Errorf("expected StateVar 0x%04X, got 0x%04X", stateVar, response.Status().StateVar)
	}

	// Verify MC counter extraction
	expectedCounter := uint8(2)
	if response.MCCounter() != expectedCounter {
		t.Errorf("expected MC counter %d, got %d", expectedCounter, response.MCCounter())
	}
}

// TestReadMCResponse_MinimalPacket tests that parsing works with minimal 12-byte packet
func TestReadMCResponse_MinimalPacket(t *testing.T) {
	// Create minimal packet with only required fields
	packet := make([]byte, 12)
	binary.LittleEndian.PutUint32(packet[0:4], 0x00000000)
	binary.LittleEndian.PutUint32(packet[4:8], StateVarResponse)
	binary.LittleEndian.PutUint16(packet[8:10], 0x1234)
	// StateVar with counter 4 in low nibble: mask high nibble and set counter
	stateVar := uint16((0x5678 & 0xFFF0) | 4) // Preserve high nibble, set counter 4 in low nibble
	binary.LittleEndian.PutUint16(packet[10:12], stateVar)

	response, err := ReadMCResponse(packet)
	if err != nil {
		t.Fatalf("ReadMCResponse() error = %v", err)
	}

	if response == nil {
		t.Fatal("ReadMCResponse() returned nil response")
	}

	if response.MCCounter() != 4 {
		t.Errorf("expected MC counter 4, got %d", response.MCCounter())
	}
}

// TestReadMCResponse_InvalidLength tests error handling for short packets
func TestReadMCResponse_InvalidLength(t *testing.T) {
	// Too short packet (less than 12 bytes)
	packet := make([]byte, 10)

	_, err := ReadMCResponse(packet)
	if err == nil {
		t.Error("expected error for short packet, got nil")
	}
}

// TestMCCommandResponse_RoundTrip tests write-then-read round trip
func TestMCCommandResponse_RoundTrip(t *testing.T) {
	originalStatus := &protocol_common.Status{
		StatusWord:     0x9ABC,
		StateVar:       0xDEF0,
		ActualPosition: 123456,
		DemandPosition: 789012,
		Current:        1000,
		WarnWord:       0x1111,
		ErrorCode:      0x2222,
	}
	originalCounter := uint8(1)

	// Write packet
	response := NewMCCommandResponse(originalStatus, originalCounter)
	packet, err := response.WritePacket()
	if err != nil {
		t.Fatalf("WritePacket() error = %v", err)
	}

	// Read packet back
	parsedResponse, err := ReadMCResponse(packet)
	if err != nil {
		t.Fatalf("ReadMCResponse() error = %v", err)
	}

	// Verify StatusWord (only field we write/read)
	if parsedResponse.Status().StatusWord != originalStatus.StatusWord {
		t.Errorf("StatusWord mismatch: expected 0x%04X, got 0x%04X",
			originalStatus.StatusWord, parsedResponse.Status().StatusWord)
	}

	// Verify StateVar (only field we write/read)
	expectedStateVar := (originalStatus.StateVar & 0xFFF0) | uint16(originalCounter&0x0F)
	if parsedResponse.Status().StateVar != expectedStateVar {
		t.Errorf("StateVar mismatch: expected 0x%04X, got 0x%04X",
			expectedStateVar, parsedResponse.Status().StateVar)
	}

	// Verify MC counter
	if parsedResponse.MCCounter() != originalCounter {
		t.Errorf("MC counter mismatch: expected %d, got %d",
			originalCounter, parsedResponse.MCCounter())
	}
}

// TestMCCommandResponse_ResponseBitsOnly tests that only StatusWord + StateVar bits are set
func TestMCCommandResponse_ResponseBitsOnly(t *testing.T) {
	status := &protocol_common.Status{
		StatusWord: 0x1234,
		StateVar:   0x5678,
	}
	response := NewMCCommandResponse(status, 2)

	packet, err := response.WritePacket()
	if err != nil {
		t.Fatalf("WritePacket() error = %v", err)
	}

	repBits := binary.LittleEndian.Uint32(packet[4:8])

	// Verify bits 0-1 and 3 are set (StatusWord + StateVar + DemandPosition)
	expectedBits := uint32(protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar | protocol_common.RespBitDemandPosition)
	if repBits != expectedBits {
		t.Errorf("expected response bits 0x%08X (StatusWord + StateVar + DemandPosition), got 0x%08X", expectedBits, repBits)
	}

	// Verify required bits are set (allow additional bits for future extensibility)
	requiredBits := uint32(protocol_common.RespBitStatusWord | protocol_common.RespBitStateVar)
	if repBits&requiredBits != requiredBits {
		t.Errorf("required response bits 0x%08X not all set in 0x%08X", requiredBits, repBits)
	}
}
