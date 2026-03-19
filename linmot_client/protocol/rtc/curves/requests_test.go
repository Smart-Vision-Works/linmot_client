package protocol_curves

import (
	"testing"

	protocol_rtc "github.com/Smart-Vision-Works/staged_robot/protocol/rtc"
)

func TestRoundTrip_SaveAllCurvesRequest(t *testing.T) {
	// Create request
	request := NewSaveAllCurvesRequest()

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

	if readReq.CmdCode() != protocol_rtc.CommandCode.SaveAllCurvesToFlash {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.SaveAllCurvesToFlash, readReq.CmdCode())
	}
	if counter != 1 {
		t.Errorf("expected rtcCounter 1, got %d", counter)
	}
}

func TestRoundTrip_DeleteAllCurvesRequest(t *testing.T) {
	// Create request
	request := NewDeleteAllCurvesRequest()

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

	if readReq.CmdCode() != protocol_rtc.CommandCode.DeleteAllCurves {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.DeleteAllCurves, readReq.CmdCode())
	}
	if counter != 2 {
		t.Errorf("expected rtcCounter 2, got %d", counter)
	}
}

func TestRoundTrip_StartAddingCurveRequest(t *testing.T) {
	// Create request
	request := NewStartAddingCurveRequest(5)

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

	if readReq.UPID() != 5 {
		t.Errorf("expected curveID 5, got %d", readReq.UPID())
	}
	if readReq.CmdCode() != protocol_rtc.CommandCode.StartAddingCurve {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.StartAddingCurve, readReq.CmdCode())
	}
	if counter != 3 {
		t.Errorf("expected rtcCounter 3, got %d", counter)
	}
}

func TestRoundTrip_AddCurveInfoBlockRequest(t *testing.T) {
	// Create request
	request := NewAddCurveInfoBlockRequest(10, 0x1234, 0x5678)

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

	if readReq.UPID() != 10 {
		t.Errorf("expected curveID 10, got %d", readReq.UPID())
	}
	if readReq.CmdCode() != protocol_rtc.CommandCode.AddCurveInfoBlock {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.AddCurveInfoBlock, readReq.CmdCode())
	}
	// Verify packed data
	expectedValue := int32(0x1234 | (0x5678 << 16))
	if readReq.Value() != expectedValue {
		t.Errorf("expected value 0x%08X, got 0x%08X", expectedValue, readReq.Value())
	}
	if counter != 4 {
		t.Errorf("expected rtcCounter 4, got %d", counter)
	}
}

func TestRoundTrip_StartGettingCurveRequest(t *testing.T) {
	// Create request
	request := NewStartGettingCurveRequest(15)

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

	if readReq.UPID() != 15 {
		t.Errorf("expected curveID 15, got %d", readReq.UPID())
	}
	if readReq.CmdCode() != protocol_rtc.CommandCode.StartGettingCurve {
		t.Errorf("expected cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.StartGettingCurve, readReq.CmdCode())
	}
	if counter != 5 {
		t.Errorf("expected rtcCounter 5, got %d", counter)
	}
}
