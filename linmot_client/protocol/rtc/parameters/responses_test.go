package protocol_parameters

import (
	"encoding/binary"
	"testing"

	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/staged_robot/protocol/rtc"
)

func testRTCStatus() *protocol_common.Status {
	return &protocol_common.Status{
		StatusWord: 0x1234,
		StateVar:   0x5678,
	}
}

func writeRTCResponse(status *protocol_common.Status, value int32, rtcCounter uint8, upid uint16, cmdCode uint8, rtcStatus uint8) []byte {
	packet := make([]byte, 50)

	binary.LittleEndian.PutUint32(packet[0:4], protocol_common.RequestFlags.RTCCommand)
	binary.LittleEndian.PutUint32(packet[4:8], protocol_common.ResponseFlags.RTCReply)
	binary.LittleEndian.PutUint16(packet[8:10], status.StatusWord)
	binary.LittleEndian.PutUint16(packet[10:12], status.StateVar)
	binary.LittleEndian.PutUint32(packet[12:16], uint32(status.ActualPosition))
	binary.LittleEndian.PutUint32(packet[16:20], uint32(status.DemandPosition))
	binary.LittleEndian.PutUint16(packet[20:22], uint16(status.Current))
	binary.LittleEndian.PutUint16(packet[22:24], status.WarnWord)
	binary.LittleEndian.PutUint16(packet[24:26], status.ErrorCode)

	rtcOffset, err := protocol_common.ResponseBlockOffset(protocol_common.ResponseFlags.RTCReply, protocol_common.RespBitRTCReplyData)
	if err != nil {
		rtcOffset = len(packet) - protocol_common.BlockSizes[protocol_common.RespBitRTCReplyData]
	}

	packet[rtcOffset+protocol_rtc.RTCDataOffsetCounter] = rtcCounter & protocol_rtc.RTCCounterMask
	packet[rtcOffset+protocol_rtc.RTCDataOffsetCmdOrStatus] = cmdCode
	binary.LittleEndian.PutUint16(packet[rtcOffset+protocol_rtc.RTCDataOffsetUPID:rtcOffset+protocol_rtc.RTCDataOffsetUPID+2], upid)
	binary.LittleEndian.PutUint32(packet[rtcOffset+protocol_rtc.RTCDataOffsetValue:rtcOffset+protocol_rtc.RTCDataOffsetValue+4], uint32(value))

	_ = rtcStatus

	return packet
}

func responseCmdOrStatusFromPacket(t *testing.T, packet []byte) uint8 {
	t.Helper()

	rtcOffset, err := protocol_common.ResponseBlockOffsetFromPacket(packet, protocol_common.RespBitRTCReplyData)
	if err != nil {
		t.Fatalf("ResponseBlockOffsetFromPacket() error: %v", err)
	}

	return packet[rtcOffset+protocol_rtc.RTCDataOffsetCmdOrStatus]
}

func TestReadRTCResponse_TypedReadROMPreservesCommandCode(t *testing.T) {
	packet := writeRTCResponse(testRTCStatus(), 123, 3, uint16(protocol_common.Parameter.Speed1), protocol_rtc.CommandCode.ReadROM, 0)

	resp, err := protocol_rtc.ReadRTCResponse(packet, protocol_rtc.CommandCode.ReadROM, uint16(protocol_common.Parameter.Speed1))
	if err != nil {
		t.Fatalf("ReadRTCResponse() error: %v", err)
	}

	got, ok := resp.(*ReadVelocityResponse)
	if !ok {
		t.Fatalf("expected *ReadVelocityResponse, got %T", resp)
	}

	roundTrip, err := got.WritePacket()
	if err != nil {
		t.Fatalf("WritePacket() error: %v", err)
	}
	if gotCmd := responseCmdOrStatusFromPacket(t, roundTrip); gotCmd != protocol_rtc.CommandCode.ReadROM {
		t.Fatalf("expected serialized cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.ReadROM, gotCmd)
	}
}

func TestReadRTCResponse_TypedWriteROMPreservesCommandCode(t *testing.T) {
	packet := writeRTCResponse(testRTCStatus(), 456, 4, uint16(protocol_common.Parameter.Position1), protocol_rtc.CommandCode.WriteROM, 0)

	resp, err := protocol_rtc.ReadRTCResponse(packet, protocol_rtc.CommandCode.WriteROM, uint16(protocol_common.Parameter.Position1))
	if err != nil {
		t.Fatalf("ReadRTCResponse() error: %v", err)
	}

	got, ok := resp.(*WritePosition1Response)
	if !ok {
		t.Fatalf("expected *WritePosition1Response, got %T", resp)
	}

	roundTrip, err := got.WritePacket()
	if err != nil {
		t.Fatalf("WritePacket() error: %v", err)
	}
	if gotCmd := responseCmdOrStatusFromPacket(t, roundTrip); gotCmd != protocol_rtc.CommandCode.WriteROM {
		t.Fatalf("expected serialized cmdCode 0x%02X, got 0x%02X", protocol_rtc.CommandCode.WriteROM, gotCmd)
	}
}
