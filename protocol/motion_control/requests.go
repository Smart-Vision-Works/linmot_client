package protocol_motion_control

import (
	"encoding/binary"
	"time"

	protocol_common "gsail-go/linmot/protocol/common"
)

// Compile-time interface checks
var (
	_ MCRequest                      = (*MCCommandRequest)(nil)
	_ protocol_common.Request        = (*MCCommandRequest)(nil)
	_ protocol_common.PacketWritable = (*MCCommandRequest)(nil)
)

// MCRequest is the interface for all Motion Control requests.
type MCRequest interface {
	protocol_common.Request
	Header() MCHeader
	SetCounter(counter uint8)
}

// MCCommandRequest represents an Motion Control command with header and parameters.
type MCCommandRequest struct {
	header     MCHeader
	parameters [15]uint16 // Up to 15 parameter words (Word 2-16)
	paramCount int        // Actual number of parameters used
}

// NewMCCommandRequest creates a new MC command request.
func NewMCCommandRequest(masterID MasterID, subID uint8, parameters []uint16) *MCCommandRequest {
	req := &MCCommandRequest{
		header: MCHeader{
			MasterID: masterID,
			SubID:    subID,
			Counter:  0, // Will be set by request manager
		},
		paramCount: len(parameters),
	}

	// Copy parameters (max 15)
	if req.paramCount > 15 {
		req.paramCount = 15
	}
	for i := 0; i < req.paramCount; i++ {
		req.parameters[i] = parameters[i]
	}

	return req
}

// Header returns the MC command header.
func (r *MCCommandRequest) Header() MCHeader {
	return r.header
}

// SetCounter sets the counter value in the header.
// This is called by the request manager before sending.
func (r *MCCommandRequest) SetCounter(counter uint8) {
	r.header.Counter = counter & 0x0F
}

// Parameters returns the parameter words array.
// Parameters are stored as 15 uint16 words, with position/increment typically in words 0-1 (int32, little-endian).
func (r *MCCommandRequest) Parameters() [15]uint16 {
	return r.parameters
}

// OperationTimeout returns the timeout duration for Motion Control requests.
func (*MCCommandRequest) OperationTimeout() time.Duration {
	return protocol_common.DefaultOperationTimeout
}

// WritePacket serializes an MC command request to a UDP packet.
// Packet structure (40 bytes):
// - Bytes 0-3: Request definition (bit 1 set for Motion Control)
// - Bytes 4-7: Response definition
// - Bytes 8-9: MC command header (Master ID, Sub ID, Counter)
// - Bytes 10-39: MC parameter data (15 words)
func (r *MCCommandRequest) WritePacket() ([]byte, error) {
	packet := make([]byte, MinMCPacketSize)

	// Header: Request bit 1 (Motion Control), Response bits 0-1,3 (StatusWord + StateVar + DemandPosition for MC)
	// DemandPosition is requested to allow verification of motion commands
	binary.LittleEndian.PutUint32(packet[0:4], protocol_common.RequestFlags.MotionControl)
	responseBits := StateVarResponse | protocol_common.RespBitDemandPosition
	binary.LittleEndian.PutUint32(packet[4:8], responseBits)

	// MC command header (word 1)
	r.header.EncodeHeaderBytes(packet[8:10])

	// MC parameters (words 2-16)
	for i := 0; i < 15; i++ {
		offset := 10 + (i * 2)
		binary.LittleEndian.PutUint16(packet[offset:offset+2], r.parameters[i])
	}

	return packet, nil
}

// ReadMCRequest parses an Motion Control request packet.
func ReadMCRequest(data []byte) (*MCCommandRequest, error) {
	if len(data) < MinMCPacketSize {
		return nil, protocol_common.NewPacketTooShortError("ReadMCRequest", len(data), MinMCPacketSize, data)
	}

	// Parse header
	header := DecodeHeaderBytes(data[8:10])

	// Parse parameters
	var parameters [15]uint16
	for i := 0; i < 15; i++ {
		offset := 10 + (i * 2)
		parameters[i] = binary.LittleEndian.Uint16(data[offset : offset+2])
	}

	return &MCCommandRequest{
		header:     header,
		parameters: parameters,
		paramCount: 15, // Assume all parameters present when parsing
	}, nil
}
