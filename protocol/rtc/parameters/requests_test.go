package protocol_parameters

import (
	"testing"

	protocol_rtc "github.com/Smart-Vision-Works/linmot_client/protocol/rtc"
)

func TestRoundTrip_WriteRAMAndROMRequest(t *testing.T) {
	// Create request
	request := NewWriteRAMAndROMRequest(0x145A, 12345)

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(3)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	// Verify counter is in packet
	if packet[8] != 3 {
		t.Errorf("expected rtcCounter 3 in packet[8], got %d", packet[8])
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

	if readReq.UPID() != 0x145A {
		t.Errorf("expected upid 0x145A, got 0x%04X", readReq.UPID())
	}
	if readReq.Value() != 12345 {
		t.Errorf("expected value 12345, got %d", readReq.Value())
	}
	if readReq.CmdCode() != protocol_rtc.CommandCode.WriteRAMAndROM {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.WriteRAMAndROM, readReq.CmdCode())
	}
	if counter != 3 {
		t.Errorf("expected rtcCounter 3, got %d", counter)
	}
}

func TestRoundTrip_GetMinValueRequest(t *testing.T) {
	// Create request
	request := NewGetMinValueRequest(0x145B)

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(5)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	// Verify counter is in packet
	if packet[8] != 5 {
		t.Errorf("expected rtcCounter 5 in packet[8], got %d", packet[8])
	}

	// Read from packet
	readReqAny, counter, err := protocol_rtc.ReadRTCGetParamRequest(packet)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	readReq, ok := readReqAny.(*protocol_rtc.RTCGetParamRequest)
	if !ok {
		t.Fatalf("expected *RTCGetParamRequest, got %T", readReqAny)
	}

	if readReq.UPID() != 0x145B {
		t.Errorf("expected upid 0x145B, got 0x%04X", readReq.UPID())
	}
	if readReq.CmdCode() != protocol_rtc.CommandCode.GetMinValue {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.GetMinValue, readReq.CmdCode())
	}
	if counter != 5 {
		t.Errorf("expected rtcCounter 5, got %d", counter)
	}
}

func TestRoundTrip_GetMaxValueRequest(t *testing.T) {
	// Create request
	request := NewGetMaxValueRequest(0x145C)

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(7)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	// Read from packet
	readReqAny, counter, err := protocol_rtc.ReadRTCGetParamRequest(packet)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	readReq, ok := readReqAny.(*protocol_rtc.RTCGetParamRequest)
	if !ok {
		t.Fatalf("expected *RTCGetParamRequest, got %T", readReqAny)
	}

	if readReq.UPID() != 0x145C {
		t.Errorf("expected upid 0x145C, got 0x%04X", readReq.UPID())
	}
	if readReq.CmdCode() != protocol_rtc.CommandCode.GetMaxValue {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.GetMaxValue, readReq.CmdCode())
	}
	if counter != 7 {
		t.Errorf("expected rtcCounter 7, got %d", counter)
	}
}

func TestRoundTrip_GetDefaultValueRequest(t *testing.T) {
	// Create request
	request := NewGetDefaultValueRequest(0x145D)

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(9)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	// Read from packet
	readReqAny, counter, err := protocol_rtc.ReadRTCGetParamRequest(packet)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	readReq, ok := readReqAny.(*protocol_rtc.RTCGetParamRequest)
	if !ok {
		t.Fatalf("expected *RTCGetParamRequest, got %T", readReqAny)
	}

	if readReq.UPID() != 0x145D {
		t.Errorf("expected upid 0x145D, got 0x%04X", readReq.UPID())
	}
	if readReq.CmdCode() != protocol_rtc.CommandCode.GetDefaultValue {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.GetDefaultValue, readReq.CmdCode())
	}
	if counter != 9 {
		t.Errorf("expected rtcCounter 9, got %d", counter)
	}
}

func TestRoundTrip_StartGettingUPIDListRequest(t *testing.T) {
	// Create request
	request := NewStartGettingUPIDListRequest(0x1000)

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(11)
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

	if readReq.UPID() != 0x1000 {
		t.Errorf("expected upid 0x1000, got 0x%04X", readReq.UPID())
	}
	if readReq.CmdCode() != protocol_rtc.CommandCode.StartGettingUPIDList {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.StartGettingUPIDList, readReq.CmdCode())
	}
	if counter != 11 {
		t.Errorf("expected rtcCounter 11, got %d", counter)
	}
}

func TestRoundTrip_GetNextUPIDListItemRequest(t *testing.T) {
	// Create request
	request := NewGetNextUPIDListItemRequest()

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(13)
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

	if readReq.CmdCode() != protocol_rtc.CommandCode.GetNextUPIDListItem {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.GetNextUPIDListItem, readReq.CmdCode())
	}
	if counter != 13 {
		t.Errorf("expected rtcCounter 13, got %d", counter)
	}
}
