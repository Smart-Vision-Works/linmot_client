package protocol_common

import (
	"encoding/binary"
)

// ReadPacketHeader extracts request and response flags from packet header.
// Returns reqBits, repBits, and an error if the packet is too short.
func ReadPacketHeader(data []byte) (uint32, uint32, error) {
	if len(data) < PacketHeaderSize {
		return 0, 0, NewPacketTooShortError("ReadPacketHeader", len(data), PacketHeaderSize, data)
	}
	reqBits := binary.LittleEndian.Uint32(data[0:4])
	repBits := binary.LittleEndian.Uint32(data[4:8])
	return reqBits, repBits, nil
}
