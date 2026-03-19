package protocol_vai16

import (
	protocol_motion_control "gsail-go/linmot/protocol/motion_control"
)

// ============================================================================
// 16-Bit VAI Command Constructors
// ============================================================================
// These commands use 16-bit position values for reduced precision (±3.28mm range).

// NewVAI16GoToPosCommand creates a 16-bit VAI Go To Position command (090xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int16, ±3.28mm range)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.73
func NewVAI16GoToPosCommand(position int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	// Position (int16 → 1 word)
	parameters[0] = uint16(position)

	// Velocity (uint32 → 2 words, little-endian)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)

	// Acceleration (uint32 → 2 words, little-endian)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)

	// Deceleration (uint32 → 2 words, little-endian)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.GoToPos),
		parameters,
	)
}

// NewVAI16IncrementDemPosCommand creates a 16-bit VAI Increment Demand Position command (091xh).
func NewVAI16IncrementDemPosCommand(increment int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(increment)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.IncrementDemPos),
		parameters,
	)
}

// NewVAI16IncrementTargetPosCommand creates a 16-bit VAI Increment Target Position command (092xh).
func NewVAI16IncrementTargetPosCommand(increment int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(increment)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.IncrementTargetPos),
		parameters,
	)
}

// NewVAI16GoToPosFromActPosCommand creates a 16-bit VAI Go To Pos From Act Pos And Act Vel command (093xh).
func NewVAI16GoToPosFromActPosCommand(position int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(position)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.GoToPosFromActPosAndActVel),
		parameters,
	)
}

// NewVAI16GoToPosFromActPosDemVelZeroCommand creates a 16-bit VAI Go To Pos From Act Pos Starting With Dem Vel = 0 command (094xh).
func NewVAI16GoToPosFromActPosDemVelZeroCommand(position int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(position)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.GoToPosFromActPosDemVelZero),
		parameters,
	)
}

// NewVAI16IncrementActPosCommand creates a 16-bit VAI Increment Act Pos command (095xh).
func NewVAI16IncrementActPosCommand(increment int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(increment)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.IncrementActPos),
		parameters,
	)
}

// NewVAI16IncrementActPosDemVelZeroCommand creates a 16-bit VAI Increment Act Pos Starting With Dem Vel = 0 command (096xh).
func NewVAI16IncrementActPosDemVelZeroCommand(increment int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(increment)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.IncrementActPosDemVelZero),
		parameters,
	)
}

// NewVAI16StopCommand creates a 16-bit VAI Stop command (097xh).
func NewVAI16StopCommand() *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.Stop),
		parameters,
	)
}

// NewVAI16GoToPosAfterActualCommandCommand creates a 16-bit VAI Go To Pos After Actual Command (098xh).
func NewVAI16GoToPosAfterActualCommandCommand(position int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(position)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.GoToPosAfterActualCommand),
		parameters,
	)
}

// NewVAI16GoToPosOnRisingTriggerCommand creates a 16-bit VAI Go To Pos On Rising Trigger Event command (09Axh).
func NewVAI16GoToPosOnRisingTriggerCommand(position int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(position)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.GoToPosOnRisingTrigger),
		parameters,
	)
}

// NewVAI16IncrementTargetPosOnRisingCommand creates a 16-bit VAI Increment Target Pos On Rising Trigger Event command (09Bxh).
func NewVAI16IncrementTargetPosOnRisingCommand(increment int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(increment)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.IncrementTargetPosOnRising),
		parameters,
	)
}

// NewVAI16GoToPosOnFallingTriggerCommand creates a 16-bit VAI Go To Pos On Falling Trigger Event command (09Cxh).
func NewVAI16GoToPosOnFallingTriggerCommand(position int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(position)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.GoToPosOnFallingTrigger),
		parameters,
	)
}

// NewVAI16IncrementTargetPosOnFallingCommand creates a 16-bit VAI Increment Target Pos On Falling Trigger Event command (09Dxh).
func NewVAI16IncrementTargetPosOnFallingCommand(increment int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(increment)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.IncrementTargetPosOnFalling),
		parameters,
	)
}

// NewVAI16ChangeParamsOnPosTransCommand creates a 16-bit VAI Change Motion Parameters On Positive Position Transition command (09Exh).
func NewVAI16ChangeParamsOnPosTransCommand(transitionPos int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(transitionPos)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.ChangeParamsOnPosTrans),
		parameters,
	)
}

// NewVAI16ChangeParamsOnNegTransCommand creates a 16-bit VAI Change Motion Parameters On Negative Position Transition command (09Fxh).
func NewVAI16ChangeParamsOnNegTransCommand(transitionPos int16, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(transitionPos)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)
	parameters[5] = uint16(decel & 0xFFFF)
	parameters[6] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAI16Bit,
		uint8(SubIDs.ChangeParamsOnNegTrans),
		parameters,
	)
}
