package protocol_rtc

import (
	"encoding/binary"
	"testing"

	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
)

// TestRTCSetParamRequest_WritePacket tests RTC Set Parameter request serialization
func TestRTCSetParamRequest_WritePacket(t *testing.T) {
	// Create a set parameter request
	request := NewRTCSetParamRequest(0x1450, 12345, CommandCode.WriteRAM)

	packet, err := request.WriteRtcPacket(5)
	if err != nil {
		t.Fatalf("WriteRtcPacket() error = %v", err)
	}

	// Verify packet length (8-byte header + 8-byte RTC data = 16 bytes)
	if len(packet) != 16 {
		t.Errorf("expected packet length 16, got %d", len(packet))
	}

	// Verify request definition (bit 2: RTC Command)
	reqBits := binary.LittleEndian.Uint32(packet[0:4])
	if reqBits != protocol_common.RequestFlags.RTCCommand {
		t.Errorf("expected request bits 0x%08X, got 0x%08X",
			protocol_common.RequestFlags.RTCCommand, reqBits)
	}

	// Verify response definition (bits 0-8: Standard + RTC Reply)
	repBits := binary.LittleEndian.Uint32(packet[4:8])
	if repBits != protocol_common.ResponseFlags.RTCReply {
		t.Errorf("expected response bits 0x%08X, got 0x%08X",
			protocol_common.ResponseFlags.RTCReply, repBits)
	}

	// Verify counter is embedded
	counter := packet[8]
	if counter != 5 {
		t.Errorf("expected counter 5, got %d", counter)
	}

	// Verify command code
	cmdCode := packet[9]
	if cmdCode != CommandCode.WriteRAM {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", CommandCode.WriteRAM, cmdCode)
	}
}

// TestRTCGetParamRequest_WritePacket tests RTC Get Parameter request serialization
func TestRTCGetParamRequest_WritePacket(t *testing.T) {
	// Create a get parameter request
	request := NewRTCGetParamRequest(0x1450, CommandCode.ReadRAM)

	packet, err := request.WriteRtcPacket(3)
	if err != nil {
		t.Fatalf("WriteRtcPacket() error = %v", err)
	}

	// Verify packet length
	if len(packet) != 16 {
		t.Errorf("expected packet length 16, got %d", len(packet))
	}

	// Verify request bit (RTC Command)
	reqBits := binary.LittleEndian.Uint32(packet[0:4])
	if reqBits != protocol_common.RequestFlags.RTCCommand {
		t.Errorf("expected request bits 0x%08X, got 0x%08X",
			protocol_common.RequestFlags.RTCCommand, reqBits)
	}

	// Verify counter
	counter := packet[8]
	if counter != 3 {
		t.Errorf("expected counter 3, got %d", counter)
	}
}

// TestRTCRequest_CounterRange tests RTC counter range (1-14)
func TestRTCRequest_CounterRange(t *testing.T) {
	request := NewRTCSetParamRequest(0x1450, 0, CommandCode.WriteRAM)

	// Test valid counter values (1-14)
	for counter := uint8(1); counter <= 14; counter++ {
		packet, err := request.WriteRtcPacket(counter)
		if err != nil {
			t.Fatalf("WriteRtcPacket(%d) error = %v", counter, err)
		}

		gotCounter := packet[8]
		if gotCounter != counter {
			t.Errorf("counter %d: expected %d in packet, got %d",
				counter, counter, gotCounter)
		}
	}
}

// TestReadRTCSetParamRequest_RoundTrip tests parsing an RTC Set Parameter request
func TestReadRTCSetParamRequest_RoundTrip(t *testing.T) {
	original := NewRTCSetParamRequest(0x1450, 54321, CommandCode.WriteRAM)

	packet, err := original.WriteRtcPacket(7)
	if err != nil {
		t.Fatalf("WriteRtcPacket() error = %v", err)
	}

	// Parse back
	parsedAny, counter, err := ReadRTCSetParamRequest(packet)
	if err != nil {
		t.Fatalf("ReadRTCSetParamRequest() error = %v", err)
	}

	parsed, ok := parsedAny.(*RTCSetParamRequest)
	if !ok {
		t.Fatalf("expected *RTCSetParamRequest, got %T", parsedAny)
	}

	if counter != 7 {
		t.Errorf("expected counter 7, got %d", counter)
	}

	if parsed.CmdCode() != CommandCode.WriteRAM {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X",
			CommandCode.WriteRAM, parsed.CmdCode())
	}
}

// TestReadRTCGetParamRequest_RoundTrip tests parsing an RTC Get Parameter request
func TestReadRTCGetParamRequest_RoundTrip(t *testing.T) {
	original := NewRTCGetParamRequest(0x1450, CommandCode.ReadRAM)

	packet, err := original.WriteRtcPacket(10)
	if err != nil {
		t.Fatalf("WriteRtcPacket() error = %v", err)
	}

	// Parse back
	parsedAny, counter, err := ReadRTCGetParamRequest(packet)
	if err != nil {
		t.Fatalf("ReadRTCGetParamRequest() error = %v", err)
	}

	parsed, ok := parsedAny.(*RTCGetParamRequest)
	if !ok {
		t.Fatalf("expected *RTCGetParamRequest, got %T", parsedAny)
	}

	if counter != 10 {
		t.Errorf("expected counter 10, got %d", counter)
	}

	if parsed.CmdCode() != CommandCode.ReadRAM {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X",
			CommandCode.ReadRAM, parsed.CmdCode())
	}
}
