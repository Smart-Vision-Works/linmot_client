package protocol_vai

import (
	protocol_motion_control "gsail-go/linmot/protocol/motion_control"
)

// ============================================================================
// VAI Command Constructors
// ============================================================================

// NewVAIGoToPosCommand creates a VAI Go To Position command (010xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.9
func NewVAIGoToPosCommand(position int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	// Pack parameters into 15-word array
	// Word 1-2: Target Position (int32, little-endian)
	// Word 3-4: Maximal Velocity (uint32, little-endian)
	// Word 5-6: Acceleration (uint32, little-endian)
	// Word 7-8: Deceleration (uint32, little-endian)
	parameters := make([]uint16, 15)

	// Position (int32 → 2 words, little-endian)
	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)

	// Velocity (uint32 → 2 words, little-endian)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	// Acceleration (uint32 → 2 words, little-endian)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	// Deceleration (uint32 → 2 words, little-endian)
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	// Remaining words unused (already zero)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.GoToPos),
		parameters,
	)
}

// NewVAIIncrementDemPosCommand creates a VAI Increment Demand Position command (011xh).
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.10
func NewVAIIncrementDemPosCommand(increment int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	// Position Increment (int32 → 2 words, little-endian)
	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)

	// Velocity (uint32 → 2 words, little-endian)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	// Acceleration (uint32 → 2 words, little-endian)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	// Deceleration (uint32 → 2 words, little-endian)
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.IncrementDemPos),
		parameters,
	)
}

// NewVAIIncrementTargetPosCommand creates a VAI Increment Target Position command (012xh).
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.11
func NewVAIIncrementTargetPosCommand(increment int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	// Position Increment (int32 → 2 words, little-endian)
	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)

	// Velocity (uint32 → 2 words, little-endian)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	// Acceleration (uint32 → 2 words, little-endian)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	// Deceleration (uint32 → 2 words, little-endian)
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.IncrementTargetPos),
		parameters,
	)
}

// NewVAIStopCommand creates a VAI Stop command (017xh).
// This command stops the current motion.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.16
func NewVAIStopCommand() *protocol_motion_control.MCCommandRequest {
	// Stop command has no parameters (all zeros)
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.Stop),
		parameters,
	)
}

// NewVAIGoToPosFromActPosCommand creates a VAI Go To Pos From Actual Position command (013xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.12
func NewVAIGoToPosFromActPosCommand(position int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	// Target Position (int32 → 2 words, little-endian)
	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)

	// Velocity (uint32 → 2 words, little-endian)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	// Acceleration (uint32 → 2 words, little-endian)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	// Deceleration (uint32 → 2 words, little-endian)
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.GoToPosFromActPosAndActVel),
		parameters,
	)
}

// NewVAIGoToPosFromActPosDemVelZeroCommand creates a VAI Go To Pos From Act Pos Starting With Dem Vel = 0 command (014xh).
// This starts VAI setpoint generation from the actual position with start velocity forced to zero.
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.13
func NewVAIGoToPosFromActPosDemVelZeroCommand(position int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
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
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.GoToPosFromActPosDemVelZero),
		parameters,
	)
}

// NewVAIIncrementActPosCommand creates a VAI Increment Actual Position command (015xh).
// Target position is calculated by adding the increment to the actual position.
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.14
func NewVAIIncrementActPosCommand(increment int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.IncrementActPos),
		parameters,
	)
}

// NewVAIIncrementActPosDemVelZeroCommand creates a VAI Increment Act Pos Starting With Dem Vel = 0 command (016xh).
// Starts VAI from actual position with start velocity forced to zero.
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.15
func NewVAIIncrementActPosDemVelZeroCommand(increment int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.IncrementActPosDemVelZero),
		parameters,
	)
}

// NewVAIGoToPosAfterActualCommandCommand creates a VAI Go To Pos After Actual Command (018xh).
// Waits until current motion setpoint generation finishes, then starts the new VAI motion.
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.17
func NewVAIGoToPosAfterActualCommandCommand(position int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
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
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.GoToPosAfterActualCommand),
		parameters,
	)
}

// NewVAIGoToAnalogPosCommand creates a VAI Go To Analog Pos command (019xh).
// Target position comes from analog input; only motion parameters are specified.
//
// Parameters:
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.18
func NewVAIGoToAnalogPosCommand(velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	// No position parameter - target comes from analog input
	// Velocity starts at offset 0 (2 words)
	parameters[0] = uint16(velocity & 0xFFFF)
	parameters[1] = uint16((velocity >> 16) & 0xFFFF)
	parameters[2] = uint16(accel & 0xFFFF)
	parameters[3] = uint16((accel >> 16) & 0xFFFF)
	parameters[4] = uint16(decel & 0xFFFF)
	parameters[5] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.GoToAnalogPos),
		parameters,
	)
}

// NewVAIGoToPosOnRisingTriggerCommand creates a VAI Go To Pos On Rising Trigger Event command (01Axh).
// Command waits for a rising trigger event before executing.
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.19
func NewVAIGoToPosOnRisingTriggerCommand(position int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
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
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.GoToPosOnRisingTrigger),
		parameters,
	)
}

// NewVAIIncrementTargetPosOnRisingCommand creates a VAI Increment Target Pos On Rising Trigger Event command (01Bxh).
// Command waits for a rising trigger event before executing.
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.20
func NewVAIIncrementTargetPosOnRisingCommand(increment int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.IncrementTargetPosOnRising),
		parameters,
	)
}

// NewVAIGoToPosOnFallingTriggerCommand creates a VAI Go To Pos On Falling Trigger Event command (01Cxh).
// Command waits for a falling trigger event before executing.
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.21
func NewVAIGoToPosOnFallingTriggerCommand(position int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
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
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.GoToPosOnFallingTrigger),
		parameters,
	)
}

// NewVAIIncrementTargetPosOnFallingCommand creates a VAI Increment Target Pos On Falling Trigger Event command (01Dxh).
// Command waits for a falling trigger event before executing.
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.22
func NewVAIIncrementTargetPosOnFallingCommand(increment int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.IncrementTargetPosOnFalling),
		parameters,
	)
}

// NewVAIChangeParamsOnPosTransCommand creates a VAI Change Motion Parameters On Positive Position Transition command (01Exh).
// Changes motion parameters when demand position crosses the transition position in positive direction.
//
// Parameters:
//   - transitionPos: Position where parameters change in 0.1µm units (int32)
//   - velocity: Maximal velocity after event in 1E-6 m/s units (uint32)
//   - accel: Acceleration after event in 1E-5 m/s² units (uint32)
//   - decel: Deceleration after event in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.23
func NewVAIChangeParamsOnPosTransCommand(transitionPos int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(transitionPos) & 0xFFFF)
	parameters[1] = uint16((uint32(transitionPos) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.ChangeParamsOnPosTrans),
		parameters,
	)
}

// NewVAIChangeParamsOnNegTransCommand creates a VAI Change Motion Parameters On Negative Position Transition command (01Fxh).
// Changes motion parameters when demand position crosses the transition position in negative direction.
//
// Parameters:
//   - transitionPos: Position where parameters change in 0.1µm units (int32)
//   - velocity: Maximal velocity after event in 1E-6 m/s units (uint32)
//   - accel: Acceleration after event in 1E-5 m/s² units (uint32)
//   - decel: Deceleration after event in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.24
func NewVAIChangeParamsOnNegTransCommand(transitionPos int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(transitionPos) & 0xFFFF)
	parameters[1] = uint16((uint32(transitionPos) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI,
		uint8(SubIDs.ChangeParamsOnNegTrans),
		parameters,
	)
}
