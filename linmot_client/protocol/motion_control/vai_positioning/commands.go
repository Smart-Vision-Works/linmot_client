package protocol_vai_positioning

import (
	protocol_motion_control "github.com/Smart-Vision-Works/staged_robot/protocol/motion_control"
)

// ============================================================================
// VAI Positioning Command Constructors
// ============================================================================
// Advanced positioning variants including captured position and command table variables.

// NewVAIGoRelativeToCapturedPosCommand creates a VAI Go Relative To Captured Pos command (0D0xh).
// Goes to a position relative to the last captured position.
//
// Parameters:
//   - offset: Offset from captured position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.124
func NewVAIGoRelativeToCapturedPosCommand(offset int32, velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(offset) & 0xFFFF)
	parameters[1] = uint16((uint32(offset) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)
	parameters[6] = uint16(decel & 0xFFFF)
	parameters[7] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPositioning,
		uint8(SubIDs.GoRelativeToCapturedPos),
		parameters,
	)
}

// NewVAIDecAcc16BitGoToPosCommand creates a VAI Dec=Acc 16 Bit Go To Pos command (0D1xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int16)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration (and deceleration) in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt
func NewVAIDecAcc16BitGoToPosCommand(position int16, velocity, accel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(position)
	parameters[1] = uint16(velocity & 0xFFFF)
	parameters[2] = uint16((velocity >> 16) & 0xFFFF)
	parameters[3] = uint16(accel & 0xFFFF)
	parameters[4] = uint16((accel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPositioning,
		uint8(SubIDs.DecAcc16BitGoToPos),
		parameters,
	)
}

// NewVAIGoToCmdTableVar1PosCommand creates a VAI Go To Cmd Table Var 1 Pos command (0D4xh).
// Target position comes from command table variable 1.
//
// Parameters:
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration in 1E-5 m/s² units (uint32)
//   - decel: Deceleration in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.124
func NewVAIGoToCmdTableVar1PosCommand(velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(velocity & 0xFFFF)
	parameters[1] = uint16((velocity >> 16) & 0xFFFF)
	parameters[2] = uint16(accel & 0xFFFF)
	parameters[3] = uint16((accel >> 16) & 0xFFFF)
	parameters[4] = uint16(decel & 0xFFFF)
	parameters[5] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPositioning,
		uint8(SubIDs.GoToCmdTableVar1Pos),
		parameters,
	)
}

// NewVAIGoToCmdTableVar2PosCommand creates a VAI Go To Cmd Table Var 2 Pos command (0D5xh).
// Target position comes from command table variable 2.
func NewVAIGoToCmdTableVar2PosCommand(velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(velocity & 0xFFFF)
	parameters[1] = uint16((velocity >> 16) & 0xFFFF)
	parameters[2] = uint16(accel & 0xFFFF)
	parameters[3] = uint16((accel >> 16) & 0xFFFF)
	parameters[4] = uint16(decel & 0xFFFF)
	parameters[5] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPositioning,
		uint8(SubIDs.GoToCmdTableVar2Pos),
		parameters,
	)
}

// NewVAIGoToCmdTableVar1PosFromActPosCommand creates a VAI Go To Cmd Table Var 1 Pos From Act Pos command (0D6xh).
func NewVAIGoToCmdTableVar1PosFromActPosCommand(velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(velocity & 0xFFFF)
	parameters[1] = uint16((velocity >> 16) & 0xFFFF)
	parameters[2] = uint16(accel & 0xFFFF)
	parameters[3] = uint16((accel >> 16) & 0xFFFF)
	parameters[4] = uint16(decel & 0xFFFF)
	parameters[5] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPositioning,
		uint8(SubIDs.GoToCmdTableVar1PosFromActPos),
		parameters,
	)
}

// NewVAIGoToCmdTableVar2PosFromActPosCommand creates a VAI Go To Cmd Table Var 2 Pos From Act Pos command (0D7xh).
func NewVAIGoToCmdTableVar2PosFromActPosCommand(velocity, accel, decel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(velocity & 0xFFFF)
	parameters[1] = uint16((velocity >> 16) & 0xFFFF)
	parameters[2] = uint16(accel & 0xFFFF)
	parameters[3] = uint16((accel >> 16) & 0xFFFF)
	parameters[4] = uint16(decel & 0xFFFF)
	parameters[5] = uint16((decel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPositioning,
		uint8(SubIDs.GoToCmdTableVar2PosFromActPos),
		parameters,
	)
}

// NewVAIStartTrigRiseConfigVAICommandCommand creates a VAI Start Trig Rise Config VAI Command (0DExh).
// Starts trigger configuration on rising edge.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.128
func NewVAIStartTrigRiseConfigVAICommandCommand() *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPositioning,
		uint8(SubIDs.StartTrigRiseConfigVAICommand),
		parameters,
	)
}

// NewVAIStartTrigFallConfigVAICommandCommand creates a VAI Start Trig Fall Config VAI Command (0DFxh).
// Starts trigger configuration on falling edge.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.129
func NewVAIStartTrigFallConfigVAICommandCommand() *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPositioning,
		uint8(SubIDs.StartTrigFallConfigVAICommand),
		parameters,
	)
}
