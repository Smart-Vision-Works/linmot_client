package protocol_rtc

import (
	"encoding/binary"
	"encoding/hex"
	"testing"

	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
)

func expectedOffset(repBits uint32, targetBit uint32) int {
	offset := protocol_common.PacketHeaderSize
	for bit := uint32(1); bit != 0; bit <<= 1 {
		if repBits&bit == 0 {
			continue
		}
		if bit == targetBit {
			return offset
		}
		offset += protocol_common.BlockSizes[bit]
	}
	return -1
}

func buildResponsePacket(repBits uint32, cmdCount uint8) []byte {
	size := protocol_common.PacketHeaderSize + protocol_common.CalculateExpectedDataSize(repBits)
	packet := make([]byte, size)
	binary.LittleEndian.PutUint32(packet[0:4], 0x00000000)
	binary.LittleEndian.PutUint32(packet[4:8], repBits)

	offset := expectedOffset(repBits, protocol_common.RespBitRTCReplyData)
	if offset < 0 {
		return packet
	}
	packet[offset+RTCDataOffsetCounter] = cmdCount & RTCCounterMask
	packet[offset+RTCDataOffsetCmdOrStatus] = CommandCode.ReadRAM
	binary.LittleEndian.PutUint16(packet[offset+RTCDataOffsetUPID:offset+RTCDataOffsetUPID+2], 0x1234)
	binary.LittleEndian.PutUint32(packet[offset+RTCDataOffsetValue:offset+RTCDataOffsetValue+4], 0xDEADBEEF)
	return packet
}

func TestResponseBlockOffsetAndRTCCmdCount(t *testing.T) {
	testCases := []struct {
		name    string
		repBits uint32
	}{
		{
			name:    "rtc_only",
			repBits: protocol_common.RespBitRTCReplyData,
		},
		{
			name:    "monitoring_and_rtc",
			repBits: protocol_common.RespBitMonitoringChannel | protocol_common.RespBitRTCReplyData,
		},
		{
			name: "status_monitoring_rtc",
			repBits: protocol_common.RespBitStatusWord |
				protocol_common.RespBitStateVar |
				protocol_common.RespBitMonitoringChannel |
				protocol_common.RespBitRTCReplyData,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			packet := buildResponsePacket(tc.repBits, 14)
			wantOffset := expectedOffset(tc.repBits, protocol_common.RespBitRTCReplyData)
			gotOffset, err := protocol_common.ResponseBlockOffsetFromPacket(packet, protocol_common.RespBitRTCReplyData)
			if err != nil {
				t.Fatalf("ResponseBlockOffsetFromPacket error: %v", err)
			}
			if gotOffset != wantOffset {
				t.Fatalf("expected offset %d, got %d", wantOffset, gotOffset)
			}

			cmdCount, err := ExtractRTCCommandCount(packet)
			if err != nil {
				t.Fatalf("ExtractRTCCommandCount error: %v", err)
			}
			if cmdCount != 14 {
				t.Fatalf("expected cmdCount 14, got %d", cmdCount)
			}
		})
	}
}

func TestExtractRTCCommandCount_FromDumpedPacket(t *testing.T) {
	repBits := protocol_common.RespBitMonitoringChannel | protocol_common.RespBitRTCReplyData
	packet := buildResponsePacket(repBits, 14)
	hexDump := hex.EncodeToString(packet)

	decoded, err := hex.DecodeString(hexDump)
	if err != nil {
		t.Fatalf("hex decode failed: %v", err)
	}
	cmdCount, err := ExtractRTCCommandCount(decoded)
	if err != nil {
		t.Fatalf("ExtractRTCCommandCount error: %v", err)
	}
	if cmdCount != 14 {
		t.Fatalf("expected cmdCount 14, got %d", cmdCount)
	}
}
