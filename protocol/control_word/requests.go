package protocol_control_word

import (
	"encoding/binary"
	protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"
	"time"
)

// Compile-time interface checks
var (
	_ protocol_common.Request        = (*ControlWordRequest)(nil)
	_ protocol_common.PacketWritable = (*ControlWordRequest)(nil)
)

// ControlWordRequest sends a control word to the drive to control state machine
type ControlWordRequest struct {
	controlWord uint16
}

// NewControlWordRequest creates a control word request with the specified word value
func NewControlWordRequest(word uint16) *ControlWordRequest {
	return &ControlWordRequest{
		controlWord: word,
	}
}

// OperationTimeout returns the timeout for control word operations
func (*ControlWordRequest) OperationTimeout() time.Duration {
	return 5 * time.Second // State transitions can take a few seconds
}

// GetControlWord returns the current control word value
func (r *ControlWordRequest) GetControlWord() uint16 {
	return r.controlWord
}

// SetBit sets a specific bit in the control word (fluent interface)
func (r *ControlWordRequest) SetBit(bit uint) *ControlWordRequest {
	// bit should be 0-15 for the uint16 control word
	r.controlWord = SetBit(r.controlWord, bit)
	return r
}

// ClearBit clears a specific bit in the control word (fluent interface)
func (r *ControlWordRequest) ClearBit(bit uint) *ControlWordRequest {
	// bit should be 0-15 for the uint16 control word
	r.controlWord = ClearBit(r.controlWord, bit)
	return r
}

// WithPattern replaces the control word with a new pattern
func (r *ControlWordRequest) WithPattern(pattern uint16) *ControlWordRequest {
	r.controlWord = pattern
	return r
}

// WritePacket serializes the control word request to a UDP packet
// Packet structure (10 bytes):
//
//	Offset 0-3: Request Definition (0x00000001 - Control Word bit set)
//	Offset 4-7: Response Definition (0x00000001 - Status Word bit set)
//	Offset 8-9: Control Word (2 bytes, little endian)
func (r *ControlWordRequest) WritePacket() ([]byte, error) {
	packet := make([]byte, 10)

	// Request Definition: bit 0 set (Control Word module)
	binary.LittleEndian.PutUint32(packet[0:4], 0x00000001)

	// Response Definition: bit 0 set (Status Word module)
	binary.LittleEndian.PutUint32(packet[4:8], 0x00000001)

	// Control Word (2 bytes)
	binary.LittleEndian.PutUint16(packet[8:10], r.controlWord)

	return packet, nil
}

// ReadControlWordRequest parses a control word request from packet data
// This is primarily for testing and mock implementations
func ReadControlWordRequest(data []byte) (*ControlWordRequest, error) {
	if len(data) < 10 {
		return nil, protocol_common.NewPacketTooShortError("ReadControlWordRequest", len(data), 10, data)
	}

	controlWord := binary.LittleEndian.Uint16(data[8:10])
	return NewControlWordRequest(controlWord), nil
}
