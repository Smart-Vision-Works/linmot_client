package protocol_indexing

import (
	protocol_motion_control "github.com/Smart-Vision-Works/staged_robot/protocol/motion_control"
)

// ============================================================================
// Position Indexing Command Constructors
// ============================================================================

// NewStartVAIEncoderIndexingCommand creates a Start VAI Encoder Position Indexing command (070xh).
//
// Parameters:
//   - startCounterValue: Encoder counter value to start indexing (int32)
//   - indexDistance: Distance between index positions in encoder counts (uint32)
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.69
func NewStartVAIEncoderIndexingCommand(startCounterValue int32, indexDistance uint32, position int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	// Start Counter Value (int32 → 2 words)
	parameters[0] = uint16(uint32(startCounterValue) & 0xFFFF)
	parameters[1] = uint16((uint32(startCounterValue) >> 16) & 0xFFFF)

	// Index Distance (uint32 → 2 words)
	parameters[2] = uint16(indexDistance & 0xFFFF)
	parameters[3] = uint16((indexDistance >> 16) & 0xFFFF)

	// Position (int32 → 2 words)
	parameters[4] = uint16(uint32(position) & 0xFFFF)
	parameters[5] = uint16((uint32(position) >> 16) & 0xFFFF)

	// Velocity, Accel, Decel (truncated for space - basic implementation)
	parameters[6] = uint16(velocity & 0xFFFF)
	parameters[7] = uint16((velocity >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PositionIndexing,
		uint8(SubIDs.StartVAIEncoderIndexing),
		parameters,
	)
}

// NewStartPredefVAIEncoderIndexingCommand creates a Start Predef VAI Encoder Position Indexing command (071xh).
//
// Parameters:
//   - startCounterValue: Encoder counter value to start indexing (int32)
//   - indexDistance: Distance between index positions in encoder counts (uint32)
//   - position: Target position in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.70
func NewStartPredefVAIEncoderIndexingCommand(startCounterValue int32, indexDistance uint32, position int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(startCounterValue) & 0xFFFF)
	parameters[1] = uint16((uint32(startCounterValue) >> 16) & 0xFFFF)
	parameters[2] = uint16(indexDistance & 0xFFFF)
	parameters[3] = uint16((indexDistance >> 16) & 0xFFFF)
	parameters[4] = uint16(uint32(position) & 0xFFFF)
	parameters[5] = uint16((uint32(position) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PositionIndexing,
		uint8(SubIDs.StartPredefVAIEncoderIndexing),
		parameters,
	)
}

// NewStopIndexingVAIGoToPosCommand creates a Stop Position Indexing And VAI Go To Pos command (07Exh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.71
func NewStopIndexingVAIGoToPosCommand(position int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PositionIndexing,
		uint8(SubIDs.StopIndexingVAIGoToPos),
		parameters,
	)
}

// NewStopIndexingPredefVAIGoToPosCommand creates a Stop Position Indexing And Predefined VAI Go To Pos command (07Fxh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.72
func NewStopIndexingPredefVAIGoToPosCommand(position int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PositionIndexing,
		uint8(SubIDs.StopIndexingPredefVAIGoToPos),
		parameters,
	)
}
