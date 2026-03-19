package protocol_rtc

import (
	"encoding/binary"
	"testing"

	protocol_common "gsail-go/linmot/protocol/common"
)

func testRTCStatus() *protocol_common.Status {
	return &protocol_common.Status{
		StatusWord: 0x1234,
		StateVar:   0x5678,
	}
}

func responseUPIDFromPacket(t *testing.T, packet []byte) uint16 {
	t.Helper()

	rtcOffset, err := protocol_common.ResponseBlockOffsetFromPacket(packet, protocol_common.RespBitRTCReplyData)
	if err != nil {
		t.Fatalf("ResponseBlockOffsetFromPacket() error: %v", err)
	}

	return binary.LittleEndian.Uint16(packet[rtcOffset+RTCDataOffsetUPID : rtcOffset+RTCDataOffsetUPID+2])
}

func responseCmdOrStatusFromPacket(t *testing.T, packet []byte) uint8 {
	t.Helper()

	rtcOffset, err := protocol_common.ResponseBlockOffsetFromPacket(packet, protocol_common.RespBitRTCReplyData)
	if err != nil {
		t.Fatalf("ResponseBlockOffsetFromPacket() error: %v", err)
	}

	return packet[rtcOffset+RTCDataOffsetCmdOrStatus]
}

func setResponseCmdOrStatus(t *testing.T, packet []byte, value uint8) {
	t.Helper()

	rtcOffset, err := protocol_common.ResponseBlockOffsetFromPacket(packet, protocol_common.RespBitRTCReplyData)
	if err != nil {
		t.Fatalf("ResponseBlockOffsetFromPacket() error: %v", err)
	}

	packet[rtcOffset+RTCDataOffsetCmdOrStatus] = value
}

func withTemporaryResponseRegistryByCmdAndUPID(t *testing.T, cmdCode uint8, upid protocol_common.ParameterID, registry ResponseRegistry) {
	t.Helper()

	key := makeRegistryKey(cmdCode, true, upid, true)
	previous, hadPrevious := responseRegistryMap[key]
	responseRegistryMap[key] = registry

	t.Cleanup(func() {
		if hadPrevious {
			responseRegistryMap[key] = previous
			return
		}
		delete(responseRegistryMap, key)
	})
}

func withTemporaryResponseRegistryByCmd(t *testing.T, cmdCode uint8, registry ResponseRegistry) {
	t.Helper()

	key := makeRegistryKey(cmdCode, true, 0, false)
	previous, hadPrevious := responseRegistryMap[key]
	responseRegistryMap[key] = registry

	t.Cleanup(func() {
		if hadPrevious {
			responseRegistryMap[key] = previous
			return
		}
		delete(responseRegistryMap, key)
	})
}

func TestReadRTCResponse_StandardReadFallbackUsesTrustedRequestUPID(t *testing.T) {
	requestUPID := uint16(0x145A)
	responseUPID := uint16(0x9999)

	packet := writeRTCResponse(testRTCStatus(), 123, 3, responseUPID, CommandCode.ReadRAM, 0)
	resp, err := ReadRTCResponse(packet, CommandCode.ReadRAM, requestUPID)
	if err != nil {
		t.Fatalf("ReadRTCResponse() error: %v", err)
	}

	got, ok := resp.(*RTCGetParamResponse)
	if !ok {
		t.Fatalf("expected *RTCGetParamResponse, got %T", resp)
	}
	if got.upid != requestUPID {
		t.Fatalf("expected trusted UPID 0x%04X, got 0x%04X", requestUPID, got.upid)
	}

	roundTrip, err := got.WritePacket()
	if err != nil {
		t.Fatalf("WritePacket() error: %v", err)
	}
	if gotUPID := responseUPIDFromPacket(t, roundTrip); gotUPID != requestUPID {
		t.Fatalf("expected serialized UPID 0x%04X, got 0x%04X", requestUPID, gotUPID)
	}
}

func TestReadRTCResponse_StandardWriteFallbackUsesTrustedRequestUPID(t *testing.T) {
	requestUPID := uint16(0x145B)
	responseUPID := uint16(0x8888)

	packet := writeRTCResponse(testRTCStatus(), 456, 4, responseUPID, CommandCode.WriteRAM, 0)
	resp, err := ReadRTCResponse(packet, CommandCode.WriteRAM, requestUPID)
	if err != nil {
		t.Fatalf("ReadRTCResponse() error: %v", err)
	}

	got, ok := resp.(*RTCSetParamResponse)
	if !ok {
		t.Fatalf("expected *RTCSetParamResponse, got %T", resp)
	}
	if got.UPID() != requestUPID {
		t.Fatalf("expected trusted UPID 0x%04X, got 0x%04X", requestUPID, got.UPID())
	}
}

func TestReadRTCResponse_GetMinValueUsesTrustedRequestUPID(t *testing.T) {
	requestUPID := uint16(0x145C)
	responseUPID := uint16(0x7777)

	packet := writeRTCResponse(testRTCStatus(), 789, 5, responseUPID, CommandCode.GetMinValue, 0)
	resp, err := ReadRTCResponse(packet, CommandCode.GetMinValue, requestUPID)
	if err != nil {
		t.Fatalf("ReadRTCResponse() error: %v", err)
	}

	got, ok := resp.(*RTCGetParamResponse)
	if !ok {
		t.Fatalf("expected *RTCGetParamResponse, got %T", resp)
	}
	if got.upid != requestUPID {
		t.Fatalf("expected trusted UPID 0x%04X, got 0x%04X", requestUPID, got.upid)
	}

	roundTrip, err := got.WritePacket()
	if err != nil {
		t.Fatalf("WritePacket() error: %v", err)
	}
	if gotCmd := responseCmdOrStatusFromPacket(t, roundTrip); gotCmd != CommandCode.GetMinValue {
		t.Fatalf("expected serialized cmdCode 0x%02X, got 0x%02X", CommandCode.GetMinValue, gotCmd)
	}
}

func TestReadRTCResponse_WriteRAMAndROMUsesTrustedRequestUPID(t *testing.T) {
	requestUPID := uint16(0x145D)
	responseUPID := uint16(0x6666)

	packet := writeRTCResponse(testRTCStatus(), 321, 6, responseUPID, CommandCode.WriteRAMAndROM, 0)
	setResponseCmdOrStatus(t, packet, 0x00)
	resp, err := ReadRTCResponse(packet, CommandCode.WriteRAMAndROM, requestUPID)
	if err != nil {
		t.Fatalf("ReadRTCResponse() error: %v", err)
	}

	got, ok := resp.(*RTCSetParamResponse)
	if !ok {
		t.Fatalf("expected *RTCSetParamResponse, got %T", resp)
	}
	if got.UPID() != requestUPID {
		t.Fatalf("expected trusted UPID 0x%04X, got 0x%04X", requestUPID, got.UPID())
	}

	roundTrip, err := got.WritePacket()
	if err != nil {
		t.Fatalf("WritePacket() error: %v", err)
	}
	if gotCmd := responseCmdOrStatusFromPacket(t, roundTrip); gotCmd != CommandCode.WriteRAMAndROM {
		t.Fatalf("expected serialized cmdCode 0x%02X, got 0x%02X", CommandCode.WriteRAMAndROM, gotCmd)
	}
}

func TestReadRTCResponse_StandardReadRegistryReceivesTrustedRequestUPID(t *testing.T) {
	requestUPID := uint16(0x7F10)
	responseUPID := uint16(0x5555)
	var registryUPID uint16

	withTemporaryResponseRegistryByCmdAndUPID(t, CommandCode.ReadRAM, protocol_common.ParameterID(requestUPID),
		func(status *protocol_common.Status, value int32, upid uint16, rtcCounter, rtcStatus, cmdCode uint8) protocol_common.Response {
			registryUPID = upid
			return NewRTCGetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
		})

	packet := writeRTCResponse(testRTCStatus(), 654, 7, responseUPID, CommandCode.ReadRAM, 0)
	resp, err := ReadRTCResponse(packet, CommandCode.ReadRAM, requestUPID)
	if err != nil {
		t.Fatalf("ReadRTCResponse() error: %v", err)
	}

	if registryUPID != requestUPID {
		t.Fatalf("expected registry UPID 0x%04X, got 0x%04X", requestUPID, registryUPID)
	}

	got, ok := resp.(*RTCGetParamResponse)
	if !ok {
		t.Fatalf("expected *RTCGetParamResponse, got %T", resp)
	}
	if got.upid != requestUPID {
		t.Fatalf("expected typed response to store trusted UPID 0x%04X, got 0x%04X", requestUPID, got.upid)
	}
}

func TestReadRTCResponse_SpecialCommandRegistryKeepsResponseUPID(t *testing.T) {
	requestUPID := uint16(0x7F20)
	responseUPID := uint16(0x4444)
	var registryUPID uint16

	withTemporaryResponseRegistryByCmd(t, CommandCode.RestartDrive,
		func(status *protocol_common.Status, value int32, upid uint16, rtcCounter, rtcStatus, cmdCode uint8) protocol_common.Response {
			registryUPID = upid
			return NewRTCSetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
		})

	packet := writeRTCResponse(testRTCStatus(), 987, 8, responseUPID, CommandCode.RestartDrive, Status.OK)
	resp, err := ReadRTCResponse(packet, CommandCode.RestartDrive, requestUPID)
	if err != nil {
		t.Fatalf("ReadRTCResponse() error: %v", err)
	}

	if registryUPID != responseUPID {
		t.Fatalf("expected special-command registry to receive response UPID 0x%04X, got 0x%04X", responseUPID, registryUPID)
	}

	got, ok := resp.(*RTCSetParamResponse)
	if !ok {
		t.Fatalf("expected *RTCSetParamResponse, got %T", resp)
	}
	if got.UPID() != responseUPID {
		t.Fatalf("expected special-command response to keep response UPID 0x%04X, got 0x%04X", responseUPID, got.UPID())
	}
}
