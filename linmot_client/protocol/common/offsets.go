package protocol_common

import "fmt"

// ResponseBlockOffset returns the byte offset for the start of a response block.
// The offset is computed by iterating response-definition bits in ascending order
// and summing the sizes of all blocks that appear before targetBit.
func ResponseBlockOffset(repBits uint32, targetBit uint32) (int, error) {
	if targetBit == 0 || targetBit&(targetBit-1) != 0 {
		return 0, fmt.Errorf("ResponseBlockOffset: targetBit must be a single bit, got 0x%08X", targetBit)
	}
	if repBits&targetBit == 0 {
		return 0, fmt.Errorf("ResponseBlockOffset: target bit 0x%08X not set in repBits 0x%08X", targetBit, repBits)
	}

	offset := PacketHeaderSize
	for bit := uint32(1); bit != 0; bit <<= 1 {
		if repBits&bit == 0 {
			continue
		}
		if bit == targetBit {
			return offset, nil
		}
		blockSize, ok := BlockSizes[bit]
		if !ok {
			return 0, fmt.Errorf("ResponseBlockOffset: unknown block size for bit 0x%08X", bit)
		}
		offset += blockSize
	}

	return 0, fmt.Errorf("ResponseBlockOffset: target bit 0x%08X not found in repBits 0x%08X", targetBit, repBits)
}

// ResponseBlockOffsetFromPacket reads the response-definition bits from the packet
// header and returns the byte offset for the requested response block.
func ResponseBlockOffsetFromPacket(data []byte, targetBit uint32) (int, error) {
	_, repBits, err := ReadPacketHeader(data)
	if err != nil {
		return 0, err
	}
	return ResponseBlockOffset(repBits, targetBit)
}
