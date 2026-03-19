package protocol_advanced

import (
	protocol_motion_control "github.com/Smart-Vision-Works/linmot_client/protocol/motion_control"
)

// packMotionParams packs position/increment plus motion parameters into the command payload.
func packMotionParams(posOrIncrement int32, velocity, accel, decel uint32) []uint16 {
	parameters := make([]uint16, 15)

	// Position or increment (int32 → 2 words, little-endian)
	parameters[0] = uint16(uint32(posOrIncrement) & 0xFFFF)
	parameters[1] = uint16((uint32(posOrIncrement) >> 16) & 0xFFFF)

	// Velocity
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	// Acceleration
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	// Deceleration
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	return parameters
}

// ============================================================================
// Advanced Features Command Constructors
// ============================================================================
// Sin VA commands use sinusoidal acceleration profiles for smoother motion.

// NewSinVAGoToPosCommand creates a Sin VA Go To Pos command (0E0xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.130
func NewSinVAGoToPosCommand(position int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packMotionParams(position, velocity, accel, decel)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Advanced,
		uint8(SubIDs.SinVAGoToPos),
		parameters,
	)
}

// NewSinVAIncrementDemandPosCommand creates a Sin VA Increment Demand Pos command (0E1xh).
func NewSinVAIncrementDemandPosCommand(increment int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packMotionParams(increment, velocity, accel, decel)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Advanced,
		uint8(SubIDs.SinVAIncrementDemandPos),
		parameters,
	)
}

// NewSinVAGoToPosFromActualPosCommand creates a Sin VA Go To Pos From Actual Pos command (0E4xh).
func NewSinVAGoToPosFromActualPosCommand(position int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packMotionParams(position, velocity, accel, decel)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Advanced,
		uint8(SubIDs.SinVAGoToPosFromActualPos),
		parameters,
	)
}

// NewSinVAIncrementActualPosCommand creates a Sin VA Increment Actual Pos command (0E6xh).
func NewSinVAIncrementActualPosCommand(increment int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packMotionParams(increment, velocity, accel, decel)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Advanced,
		uint8(SubIDs.SinVAIncrementActualPos),
		parameters,
	)
}

// NewSinVAGoToPosAfterActualCommandCommand creates a Sin VA Go To Pos After Actual Command (0E8xh).
func NewSinVAGoToPosAfterActualCommandCommand(position int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packMotionParams(position, velocity, accel, decel)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Advanced,
		uint8(SubIDs.SinVAGoToPosAfterActualCommand),
		parameters,
	)
}

// NewSinVAGoToAnalogPosCommand creates a Sin VA Go To Analog Pos command (0E9xh).
func NewSinVAGoToAnalogPosCommand(velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(velocity & 0xFFFF)
	parameters[1] = uint16((velocity >> 16) & 0xFFFF)
	parameters[2] = uint16(accel & 0xFFFF)
	parameters[3] = uint16((accel >> 16) & 0xFFFF)
	parameters[4] = uint16(decel & 0xFFFF)
	parameters[5] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Advanced,
		uint8(SubIDs.SinVAGoToAnalogPos),
		parameters,
	)
}

// NewSinVAGoToPosOnRisingTriggerCommand creates a Sin VA Go To Pos On Rising Trigger Event command (0EAxh).
func NewSinVAGoToPosOnRisingTriggerCommand(position int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packMotionParams(position, velocity, accel, decel)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Advanced,
		uint8(SubIDs.SinVAGoToPosOnRisingTrigger),
		parameters,
	)
}

// NewSinVAIncrementDemPosOnRisingCommand creates a Sin VA Increment Demand Pos On Rising Trigger Event command (0EBxh).
func NewSinVAIncrementDemPosOnRisingCommand(increment int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packMotionParams(increment, velocity, accel, decel)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Advanced,
		uint8(SubIDs.SinVAIncrementDemPosOnRising),
		parameters,
	)
}

// NewSinVAGoToPosOnFallingTriggerCommand creates a Sin VA Go To Pos On Falling Trigger Event command (0ECxh).
func NewSinVAGoToPosOnFallingTriggerCommand(position int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packMotionParams(position, velocity, accel, decel)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Advanced,
		uint8(SubIDs.SinVAGoToPosOnFallingTrigger),
		parameters,
	)
}

// NewSinVAIncrementDemPosOnFallingCommand creates a Sin VA Increment Demand Pos On Falling Trigger Event command (0EDxh).
func NewSinVAIncrementDemPosOnFallingCommand(increment int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packMotionParams(increment, velocity, accel, decel)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Advanced,
		uint8(SubIDs.SinVAIncrementDemPosOnFalling),
		parameters,
	)
}
