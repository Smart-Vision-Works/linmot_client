package protocol_operations

import (
	"testing"

	protocol_rtc "gsail-go/linmot/protocol/rtc"
)

func TestRoundTrip_RestartDriveRequest(t *testing.T) {
	// Create request
	request := NewRestartDriveRequest()

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(2)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	// Verify counter is in packet
	if packet[8] != 2 {
		t.Errorf("expected rtcCounter 2 in packet[8], got %d", packet[8])
	}

	// Read from packet
	readReqAny, counter, err := protocol_rtc.ReadRTCSetParamRequest(packet)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	readReq, ok := readReqAny.(*protocol_rtc.RTCSetParamRequest)
	if !ok {
		t.Fatalf("expected *RTCSetParamRequest, got %T", readReqAny)
	}

	if readReq.CmdCode() != protocol_rtc.CommandCode.RestartDrive {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.RestartDrive, readReq.CmdCode())
	}
	if counter != 2 {
		t.Errorf("expected rtcCounter 2, got %d", counter)
	}
}

func TestRoundTrip_SetOSROMToDefaultRequest(t *testing.T) {
	// Create request
	request := NewSetOSROMToDefaultRequest()

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(4)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	// Read from packet
	readReqAny, counter, err := protocol_rtc.ReadRTCSetParamRequest(packet)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	readReq, ok := readReqAny.(*protocol_rtc.RTCSetParamRequest)
	if !ok {
		t.Fatalf("expected *RTCSetParamRequest, got %T", readReqAny)
	}

	if readReq.CmdCode() != protocol_rtc.CommandCode.SetOSROMToDefault {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.SetOSROMToDefault, readReq.CmdCode())
	}
	if counter != 4 {
		t.Errorf("expected rtcCounter 4, got %d", counter)
	}
}

func TestRoundTrip_SetMCROMToDefaultRequest(t *testing.T) {
	// Create request
	request := NewSetMCROMToDefaultRequest()

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(6)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	// Read from packet
	readReqAny, counter, err := protocol_rtc.ReadRTCSetParamRequest(packet)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	readReq, ok := readReqAny.(*protocol_rtc.RTCSetParamRequest)
	if !ok {
		t.Fatalf("expected *RTCSetParamRequest, got %T", readReqAny)
	}

	if readReq.CmdCode() != protocol_rtc.CommandCode.SetMCROMToDefault {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.SetMCROMToDefault, readReq.CmdCode())
	}
	if counter != 6 {
		t.Errorf("expected rtcCounter 6, got %d", counter)
	}
}

func TestRoundTrip_SetInterfaceROMToDefaultRequest(t *testing.T) {
	// Create request
	request := NewSetInterfaceROMToDefaultRequest()

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(8)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	// Read from packet
	readReqAny, counter, err := protocol_rtc.ReadRTCSetParamRequest(packet)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	readReq, ok := readReqAny.(*protocol_rtc.RTCSetParamRequest)
	if !ok {
		t.Fatalf("expected *RTCSetParamRequest, got %T", readReqAny)
	}

	if readReq.CmdCode() != protocol_rtc.CommandCode.SetInterfaceROMToDefault {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.SetInterfaceROMToDefault, readReq.CmdCode())
	}
	if counter != 8 {
		t.Errorf("expected rtcCounter 8, got %d", counter)
	}
}

func TestRoundTrip_SetApplicationROMToDefaultRequest(t *testing.T) {
	// Create request
	request := NewSetApplicationROMToDefaultRequest()

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(10)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	// Read from packet
	readReqAny, counter, err := protocol_rtc.ReadRTCSetParamRequest(packet)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	readReq, ok := readReqAny.(*protocol_rtc.RTCSetParamRequest)
	if !ok {
		t.Fatalf("expected *RTCSetParamRequest, got %T", readReqAny)
	}

	if readReq.CmdCode() != protocol_rtc.CommandCode.SetApplicationROMToDefault {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.SetApplicationROMToDefault, readReq.CmdCode())
	}
	if counter != 10 {
		t.Errorf("expected rtcCounter 10, got %d", counter)
	}
}
