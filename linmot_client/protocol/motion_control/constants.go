package protocol_motion_control

// ============================================================================
// Motion Control Response Flags
// ============================================================================

// StateVarResponse is the response bit flag for Motion Control responses.
// Motion Control responses need StatusWord (bit 0) and StateVar (bit 1).
// The MC counter is echoed in bits 0-3 of StateVar for request/response correlation.
// Reference: LINUDP_LIBRARY_ANALYSIS.md - Motion Control Counter System
const StateVarResponse uint32 = 0x00000003 // Bits 0-1: StatusWord + StateVar

// Note: Motion Control REQUEST flag (bit 1) is now in protocol_common.RequestFlags.MotionControl

// ============================================================================
// Master ID Constants
// ============================================================================

// MasterID represents a command group identifier in Motion Control.
type MasterID uint8

// MasterIDs groups all Master ID constants for Motion Control command groups.
// Reference: LinMot_MotionCtrl.txt, Section 4.2
var MasterIDs = struct {
	InterfaceControl MasterID // 0x00 - Interface control word, outputs, reset
	VAI              MasterID // 0x01 - Variable Interpolator commands
	PredefVAI        MasterID // 0x02 - Predefined VAI commands
	Streaming        MasterID // 0x03 - P/PV/PVA streaming commands
	Curve            MasterID // 0x04 - Curve commands
	Cam              MasterID // 0x05 - Cam commands
	EncoderCam       MasterID // 0x06 - Encoder cam commands
	PositionIndexing MasterID // 0x07 - Position indexing
	VAI16Bit         MasterID // 0x09 - 16-bit VAI variants
	PredefVAI16Bit   MasterID // 0x0A - 16-bit predefined VAI
	VAIPredefAcc     MasterID // 0x0B - VAI with predefined acceleration
	VAIDecEqualsAcc  MasterID // 0x0C - VAI with deceleration = acceleration
	VAIPositioning   MasterID // 0x0D - Advanced VAI positioning
	Advanced         MasterID // 0x0E - MC-PV upload/download, capture
	Bestehorn        MasterID // 0x0F - Bestehorn VAJ (jerk-limited motion)
}{
	InterfaceControl: 0x00,
	VAI:              0x01,
	PredefVAI:        0x02,
	Streaming:        0x03,
	Curve:            0x04,
	Cam:              0x05,
	EncoderCam:       0x06,
	PositionIndexing: 0x07,
	VAI16Bit:         0x09,
	PredefVAI16Bit:   0x0A,
	VAIPredefAcc:     0x0B,
	VAIDecEqualsAcc:  0x0C,
	VAIPositioning:   0x0D,
	Advanced:         0x0E,
	Bestehorn:        0x0F,
}

// ============================================================================
// Motion Control Packet Constants
// ============================================================================

const (
	// MCDataSize is the size of Motion Control command data (32 bytes = 16 words).
	MCDataSize = 32

	// MCHeaderSize is the size of the MC command header (2 bytes = 1 word).
	MCHeaderSize = 2

	// MCParameterWords is the maximum number of parameter words (header + 15 parameters).
	MCParameterWords = 16

	// MinMCPacketSize is the minimum size for an Motion Control packet (header + MC data).
	MinMCPacketSize = 8 + MCDataSize // 8-byte LinUDP header + 32-byte MC data = 40 bytes
)

// ============================================================================
// MC Counter Constants
// ============================================================================

const (
	// MCCounterMin is the minimum value for Motion Control counter.
	MCCounterMin uint8 = 1

	// MCCounterMax is the maximum value for Motion Control counter.
	MCCounterMax uint8 = 4
)
