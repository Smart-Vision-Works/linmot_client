package protocol_bestehorn

import (
	protocol_motion_control "gsail-go/linmot/protocol/motion_control"
)

// packBestehornMotionParams packs position/increment + motion parameters into the parameter array.
func packBestehornMotionParams(posOrIncrement int32, velocity, accel, decel, jerk uint32) []uint16 {
	parameters := make([]uint16, 15)

	// Position or increment (int32 → 2 words, little-endian)
	parameters[0] = uint16(uint32(posOrIncrement) & 0xFFFF)
	parameters[1] = uint16((uint32(posOrIncrement) >> 16) & 0xFFFF)

	// Velocity (uint32 → 2 words)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	// Acceleration (uint32 → 2 words)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	// Deceleration (uint32 → 2 words)
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	// Jerk (uint32 → 2 words)
	parameters[8] = uint16(jerk & 0xFFFF)
	parameters[9] = uint16((jerk >> 16) & 0xFFFF)

	return parameters
}

// ============================================================================
// Bestehorn VAJ Command Constructors
// ============================================================================

// NewBestehornGoToPosCommand creates a Bestehorn VAJ Go To Position command (0F0xh).
//
// Bestehorn VAJ commands provide jerk-limited motion for smoother trajectories and
// reduced mechanical stress compared to standard VAI commands.
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//   - jerk: Maximum jerk in 1E-4 m/s³ units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.140
func NewBestehornGoToPosCommand(position int32, velocity, accel, decel, jerk uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packBestehornMotionParams(position, velocity, accel, decel, jerk)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Bestehorn,
		uint8(SubIDs.GoToPos),
		parameters,
	)
}

// NewBestehornIncrementDemPosCommand creates a Bestehorn VAJ Increment Demand Position command (0F1xh).
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//   - jerk: Maximum jerk in 1E-4 m/s³ units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.141
func NewBestehornIncrementDemPosCommand(increment int32, velocity, accel, decel, jerk uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packBestehornMotionParams(increment, velocity, accel, decel, jerk)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Bestehorn,
		uint8(SubIDs.IncrementDemPos),
		parameters,
	)
}

// NewBestehornGoToPosFromActualPosCommand creates a Bestehorn VAJ Go To Pos From Actual Pos command (0F4xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//   - jerk: Maximum jerk in 1E-4 m/s³ units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.142
func NewBestehornGoToPosFromActualPosCommand(position int32, velocity, accel, decel, jerk uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packBestehornMotionParams(position, velocity, accel, decel, jerk)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Bestehorn,
		uint8(SubIDs.GoToPosFromActualPos),
		parameters,
	)
}

// NewBestehornIncrementActualPosCommand creates a Bestehorn VAJ Increment Actual Pos command (0F6xh).
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//   - jerk: Maximum jerk in 1E-4 m/s³ units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.143
func NewBestehornIncrementActualPosCommand(increment int32, velocity, accel, decel, jerk uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packBestehornMotionParams(increment, velocity, accel, decel, jerk)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Bestehorn,
		uint8(SubIDs.IncrementActualPos),
		parameters,
	)
}

// NewBestehornGoToPosAfterActualCommandCommand creates a Bestehorn VAJ Go To Pos After Actual Command (0F8xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//   - jerk: Maximum jerk in 1E-4 m/s³ units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.144
func NewBestehornGoToPosAfterActualCommandCommand(position int32, velocity, accel, decel, jerk uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packBestehornMotionParams(position, velocity, accel, decel, jerk)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Bestehorn,
		uint8(SubIDs.GoToPosAfterActualCommand),
		parameters,
	)
}

// NewBestehornGoToAnalogPosCommand creates a Bestehorn VAJ Go To Analog Pos command (0F9xh).
//
// Parameters:
//   - analogPos: Analog position setpoint in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//   - jerk: Maximum jerk in 1E-4 m/s³ units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.145
func NewBestehornGoToAnalogPosCommand(analogPos int32, velocity, accel, decel, jerk uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packBestehornMotionParams(analogPos, velocity, accel, decel, jerk)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Bestehorn,
		uint8(SubIDs.GoToAnalogPos),
		parameters,
	)
}

// NewBestehornGoToPosOnRisingTriggerCommand creates a Bestehorn VAJ Go To Pos On Rising Trigger Event command (0FAxh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//   - jerk: Maximum jerk in 1E-4 m/s³ units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.146
func NewBestehornGoToPosOnRisingTriggerCommand(position int32, velocity, accel, decel, jerk uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packBestehornMotionParams(position, velocity, accel, decel, jerk)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Bestehorn,
		uint8(SubIDs.GoToPosOnRisingTrigger),
		parameters,
	)
}

// NewBestehornIncrementDemPosOnRisingTriggerCommand creates a Bestehorn VAJ Increment Demand Pos On Rising Trigger Event command (0FBxh).
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//   - jerk: Maximum jerk in 1E-4 m/s³ units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.147
func NewBestehornIncrementDemPosOnRisingTriggerCommand(increment int32, velocity, accel, decel, jerk uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packBestehornMotionParams(increment, velocity, accel, decel, jerk)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Bestehorn,
		uint8(SubIDs.IncrementDemPosOnRisingTrigger),
		parameters,
	)
}

// NewBestehornGoToPosOnFallingTriggerCommand creates a Bestehorn VAJ Go To Pos On Falling Trigger Event command (0FCxh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//   - jerk: Maximum jerk in 1E-4 m/s³ units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.148
func NewBestehornGoToPosOnFallingTriggerCommand(position int32, velocity, accel, decel, jerk uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packBestehornMotionParams(position, velocity, accel, decel, jerk)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Bestehorn,
		uint8(SubIDs.GoToPosOnFallingTrigger),
		parameters,
	)
}

// NewBestehornIncrementDemPosOnFallingTriggerCommand creates a Bestehorn VAJ Increment Demand Pos On Falling Trigger Event command (0FDxh).
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//   - jerk: Maximum jerk in 1E-4 m/s³ units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.149
func NewBestehornIncrementDemPosOnFallingTriggerCommand(increment int32, velocity, accel, decel, jerk uint32) *protocol_motion_control.MCCommandRequest {
	parameters := packBestehornMotionParams(increment, velocity, accel, decel, jerk)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Bestehorn,
		uint8(SubIDs.IncrementDemPosOnFallingTrigger),
		parameters,
	)
}
