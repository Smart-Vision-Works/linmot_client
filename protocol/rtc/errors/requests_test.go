package protocol_errors

import (
	"testing"

	protocol_rtc "github.com/Smart-Vision-Works/linmot_client/protocol/rtc"
)

func TestRoundTrip_GetErrorLogEntryCounterRequest(t *testing.T) {
	// Create request
	request := NewGetErrorLogEntryCounterRequest()

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(1)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	// Verify counter is in packet
	if packet[8] != 1 {
		t.Errorf("expected rtcCounter 1 in packet[8], got %d", packet[8])
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

	if readReq.CmdCode() != protocol_rtc.CommandCode.GetErrorLogEntryCounter {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.GetErrorLogEntryCounter, readReq.CmdCode())
	}
	if counter != 1 {
		t.Errorf("expected rtcCounter 1, got %d", counter)
	}
}

func TestRoundTrip_GetErrorLogEntryCodeRequest(t *testing.T) {
	// Create request for entry 5
	request := NewGetErrorLogEntryCodeRequest(5)

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(2)
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

	if readReq.UPID() != 5 {
		t.Errorf("expected entryNumber 5, got %d", readReq.UPID())
	}
	if readReq.CmdCode() != protocol_rtc.CommandCode.GetErrorLogEntryCode {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.GetErrorLogEntryCode, readReq.CmdCode())
	}
	if counter != 2 {
		t.Errorf("expected rtcCounter 2, got %d", counter)
	}
}

func TestRoundTrip_GetErrorLogEntryTimeLowRequest(t *testing.T) {
	// Create request for entry 10
	request := NewGetErrorLogEntryTimeLowRequest(10)

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(3)
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

	if readReq.UPID() != 10 {
		t.Errorf("expected entryNumber 10, got %d", readReq.UPID())
	}
	if readReq.CmdCode() != protocol_rtc.CommandCode.GetErrorLogEntryTimeLow {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.GetErrorLogEntryTimeLow, readReq.CmdCode())
	}
	if counter != 3 {
		t.Errorf("expected rtcCounter 3, got %d", counter)
	}
}

func TestRoundTrip_GetErrorLogEntryTimeHighRequest(t *testing.T) {
	// Create request for entry 15
	request := NewGetErrorLogEntryTimeHighRequest(15)

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

	if readReq.UPID() != 15 {
		t.Errorf("expected entryNumber 15, got %d", readReq.UPID())
	}
	if readReq.CmdCode() != protocol_rtc.CommandCode.GetErrorLogEntryTimeHigh {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.GetErrorLogEntryTimeHigh, readReq.CmdCode())
	}
	if counter != 4 {
		t.Errorf("expected rtcCounter 4, got %d", counter)
	}
}

func TestRoundTrip_GetErrorCodeTextStringletRequest(t *testing.T) {
	// Create request for error code 0x0020, stringlet 3
	request := NewGetErrorCodeTextStringletRequest(0x0020, 3)

	// Write to packet with counter
	packet, err := request.WriteRtcPacket(5)
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

	if readReq.UPID() != 0x0020 {
		t.Errorf("expected errorCode 0x0020, got 0x%04X", readReq.UPID())
	}
	if readReq.Value() != 3 {
		t.Errorf("expected stringletNumber 3, got %d", readReq.Value())
	}
	if readReq.CmdCode() != protocol_rtc.CommandCode.GetErrorCodeTextStringlet {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.GetErrorCodeTextStringlet, readReq.CmdCode())
	}
	if counter != 5 {
		t.Errorf("expected rtcCounter 5, got %d", counter)
	}
}
