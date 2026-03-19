package protocol_motion_control

import (
	"encoding/binary"

	protocol_common "gsail-go/linmot/protocol/common"
)

// Compile-time interface checks
var (
	_ MCResponse                     = (*MCCommandResponse)(nil)
	_ protocol_common.Response       = (*MCCommandResponse)(nil)
	_ protocol_common.PacketWritable = (*MCCommandResponse)(nil)
)

// MCResponse is the interface for all Motion Control responses.
// The MC counter is echoed in StateVarLow bits 0-3 for response matching.
type MCResponse interface {
	protocol_common.Response
	MCCounter() uint8
}

// MCCommandResponse represents an Motion Control command response.
// The response includes Status data with the MC counter echoed in StateVarLow.
type MCCommandResponse struct {
	status    *protocol_common.Status
	mcCounter uint8 // Counter echoed from request (extracted from StateVarLow bits 0-3)
}

// NewMCCommandResponse creates a new MC command response.
// This is primarily for testing and mock implementations.
func NewMCCommandResponse(status *protocol_common.Status, mcCounter uint8) *MCCommandResponse {
	return &MCCommandResponse{
		status:    status,
		mcCounter: mcCounter,
	}
}

// ReadMCResponse parses an Motion Control response packet.
// The MC counter is extracted from StateVarLow bits 0-3.
// MC responses include StatusWord (bit 0), StateVar (bit 1), and DemandPosition (bit 3).
func ReadMCResponse(data []byte) (*MCCommandResponse, error) {
	// Minimum packet must have status data (StatusWord + StateVar = 4 bytes minimum)
	// Full packet includes DemandPosition (4 bytes) = 16 bytes total
	if len(data) < 8+4 { // Header (8) + StatusWord (2) + StateVar (2)
		return nil, protocol_common.NewPacketTooShortError("ReadMCResponse", len(data), 12, data)
	}

	// Parse status - MC responses include StatusWord and StateVar
	// We need to parse with the correct response bits set
	status, err := protocol_common.ReadStatus(data, false)
	if err != nil {
		return nil, err
	}

	// Extract MC counter from StateVarLow (bits 0-3)
	// The drive echoes the counter value directly in StateVarLow & 0xF (range 0-3, but drive uses 1-4).
	// linudp.cs behavior:
	//   - When sending: CountNibble = (StateVarLow & 0xF) + 1, then wraps to 1 if > 4
	//   - When receiving: compares CountNibble == (StateVarLow & 0xF) (direct comparison)
	// Our extraction matches the drive's echoed value (0-3 range, but we interpret as 1-4).
	// TODO: Verify with hardware that drive actually echoes 0-3 or 1-4, and ensure our
	//       MCCounter.Next() allocation (which returns 1-4) matches linudp.cs behavior.
	stateVarLow := uint8(status.StateVar & 0xFF)
	mcCounter := stateVarLow & 0x0F

	return &MCCommandResponse{
		status:    status,
		mcCounter: mcCounter,
	}, nil
}

// Status returns the drive status from the response.
func (r *MCCommandResponse) Status() *protocol_common.Status {
	return r.status
}

// MCCounter returns the MC counter value echoed in the response.
func (r *MCCommandResponse) MCCounter() uint8 {
	return r.mcCounter
}

// WritePacket serializes an MC command response to a UDP packet.
// Packet structure (16 bytes):
// - Bytes 0-3: Request definition (echoes MotionControl bit per LinUDP V2 spec 4.4.2)
// - Bytes 4-7: Response definition (bits 0-1: StatusWord + StateVar, bit 3: DemandPosition)
// - Bytes 8-9: StatusWord (2 bytes)
// - Bytes 10-11: StateVar (2 bytes) with MC counter in bits 0-3
// - Bytes 12-15: DemandPosition (4 bytes)
// Note: DemandPosition is included because MC requests ask for it to verify motion commands.
func (r *MCCommandResponse) WritePacket() ([]byte, error) {
	// Packet size: header (8) + StatusWord (2) + StateVar (2) + DemandPosition (4) = 16 bytes
	packet := make([]byte, 16)

	// Header
	// MC responses echo the request definition (bit 1 = MotionControl) per LinUDP V2 spec 4.4.2
	// Response definition includes StatusWord (bit 0), StateVar (bit 1), and DemandPosition (bit 3)
	// StateVar contains the MC counter in bits 0-3
	binary.LittleEndian.PutUint32(packet[0:4], protocol_common.RequestFlags.MotionControl)
	// Include StatusWord (bit 0), StateVar (bit 1), and DemandPosition (bit 3)
	responseBits := StateVarResponse | protocol_common.RespBitDemandPosition
	binary.LittleEndian.PutUint32(packet[4:8], responseBits)

	// StatusWord (bit 0) - 2 bytes at offset 8
	binary.LittleEndian.PutUint16(packet[8:10], r.status.StatusWord)

	// StateVar (bit 1) - 2 bytes at offset 10
	// Embed MC counter in low nibble (bits 0-3)
	stateVarWithCounter := (r.status.StateVar & 0xFFF0) | uint16(r.mcCounter&0x0F)
	binary.LittleEndian.PutUint16(packet[10:12], stateVarWithCounter)

	// DemandPosition (bit 3) - 4 bytes at offset 12
	binary.LittleEndian.PutUint32(packet[12:16], uint32(r.status.DemandPosition))

	return packet, nil
}
