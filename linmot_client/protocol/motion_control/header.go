package protocol_motion_control

import "encoding/binary"

// MCHeader represents the Motion Control command header (Word 1).
//
// Structure (16 bits):
// - Bits 15-8: Master ID (command group)
// - Bits 7-4: Sub ID (specific command within group)
// - Bits 3-0: Command Count (1-4 range, for response matching)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.1.1
type MCHeader struct {
	MasterID MasterID // Command group identifier
	SubID    uint8    // Command identifier within group
	Counter  uint8    // Command count (1-4 range)
}

// NewMCHeader creates a new MC command header.
func NewMCHeader(masterID MasterID, subID uint8, counter uint8) MCHeader {
	return MCHeader{
		MasterID: masterID,
		SubID:    subID,
		Counter:  counter,
	}
}

// EncodeHeader packs the header into a single 16-bit word.
// The counter is embedded in the low nibble of the SubID byte.
func (h MCHeader) EncodeHeader() uint16 {
	// Bits 15-8: Master ID
	// Bits 7-4: Sub ID (high nibble)
	// Bits 3-0: Counter (low nibble)
	masterByte := uint8(h.MasterID)
	subIDByte := (h.SubID << 4) | (h.Counter & 0x0F)

	// Little-endian: low byte first
	return uint16(subIDByte) | (uint16(masterByte) << 8)
}

// DecodeHeader unpacks a 16-bit header word into MCHeader components.
func DecodeHeader(word uint16) MCHeader {
	// Extract bytes (little-endian)
	lowByte := uint8(word & 0xFF)         // SubID + Counter
	highByte := uint8((word >> 8) & 0xFF) // Master ID

	return MCHeader{
		MasterID: MasterID(highByte),
		SubID:    (lowByte >> 4) & 0x0F, // High nibble
		Counter:  lowByte & 0x0F,        // Low nibble
	}
}

// WithCounter returns a new header with the specified counter value.
// This is useful for adding a counter to a command template.
func (h MCHeader) WithCounter(counter uint8) MCHeader {
	return MCHeader{
		MasterID: h.MasterID,
		SubID:    h.SubID,
		Counter:  counter & 0x0F, // Ensure 4-bit value
	}
}

// EncodeHeaderBytes writes the header to a byte slice (little-endian).
func (h MCHeader) EncodeHeaderBytes(data []byte) {
	word := h.EncodeHeader()
	binary.LittleEndian.PutUint16(data[0:2], word)
}

// DecodeHeaderBytes reads a header from a byte slice (little-endian).
func DecodeHeaderBytes(data []byte) MCHeader {
	word := binary.LittleEndian.Uint16(data[0:2])
	return DecodeHeader(word)
}
