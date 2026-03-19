package protocol_vai_dec_acc

import (
	protocol_motion_control "gsail-go/linmot/protocol/motion_control"
)

// ============================================================================
// VAI Dec=Acc Command Constructors
// ============================================================================
// These commands use symmetric motion where deceleration equals acceleration.

// NewVAIDecAccGoToPosCommand creates a VAI Dec=Acc Go To Position command (0C0xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//   - accel: Acceleration (and deceleration) in 1E-5 m/s² units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.109
func NewVAIDecAccGoToPosCommand(position int32, velocity, accel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.GoToPos),
		parameters,
	)
}

// NewVAIDecAccIncrementDemPosCommand creates a VAI Dec=Acc Increment Demand Position command (0C1xh).
func NewVAIDecAccIncrementDemPosCommand(increment int32, velocity, accel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.IncrementDemPos),
		parameters,
	)
}

// NewVAIDecAccIncrementTargetPosCommand creates a VAI Dec=Acc Increment Target Position command (0C2xh).
func NewVAIDecAccIncrementTargetPosCommand(increment int32, velocity, accel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.IncrementTargetPos),
		parameters,
	)
}

// NewVAIDecAccGoToPosFromActPosCommand creates a VAI Dec=Acc Go To Pos From Act Pos And Act Vel command (0C3xh).
func NewVAIDecAccGoToPosFromActPosCommand(position int32, velocity, accel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.GoToPosFromActPosAndActVel),
		parameters,
	)
}

// NewVAIDecAccGoToPosFromActPosDemVelZeroCommand creates a VAI Dec=Acc Go To Pos Starting With Dem Vel = 0 command (0C4xh).
func NewVAIDecAccGoToPosFromActPosDemVelZeroCommand(position int32, velocity, accel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.GoToPosFromActPosDemVelZero),
		parameters,
	)
}

// NewVAIDecAccGoToPosWithMaxCurrCommand creates a VAI Dec=Acc Go To Pos With Max Curr command (0C5xh).
// Includes maximum current parameter.
func NewVAIDecAccGoToPosWithMaxCurrCommand(position int32, velocity, accel, maxCurrent uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)
	parameters[6] = uint16(maxCurrent & 0xFFFF)
	parameters[7] = uint16((maxCurrent >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.GoToPosWithMaxCurr),
		parameters,
	)
}

// NewVAIDecAccGoToPosFromActPosMaxCurrCommand creates a VAI Dec=Acc Go To Pos From Act Pos With Max Curr command (0C6xh).
func NewVAIDecAccGoToPosFromActPosMaxCurrCommand(position int32, velocity, accel, maxCurrent uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)
	parameters[6] = uint16(maxCurrent & 0xFFFF)
	parameters[7] = uint16((maxCurrent >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.GoToPosFromActPosMaxCurr),
		parameters,
	)
}

// NewVAIDecAccGoToPosFromActPosDemVelZeroMaxCurrCommand creates a VAI Dec=Acc Go To Pos From Act Pos, Dem Vel = 0 and With Max Curr command (0C7xh).
func NewVAIDecAccGoToPosFromActPosDemVelZeroMaxCurrCommand(position int32, velocity, accel, maxCurrent uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)
	parameters[6] = uint16(maxCurrent & 0xFFFF)
	parameters[7] = uint16((maxCurrent >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.GoToPosFromActPosDemVelZeroMaxCurr),
		parameters,
	)
}

// NewVAIDecAccGoToPosAfterActualCommandCommand creates a VAI Dec=Acc Go To Pos After Actual Command (0C8xh).
func NewVAIDecAccGoToPosAfterActualCommandCommand(position int32, velocity, accel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.GoToPosAfterActualCommand),
		parameters,
	)
}

// NewVAIDecAccGoToPosOnRisingTriggerCommand creates a VAI Dec=Acc Go To Pos On Rising Trigger Event command (0CAxh).
func NewVAIDecAccGoToPosOnRisingTriggerCommand(position int32, velocity, accel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.GoToPosOnRisingTrigger),
		parameters,
	)
}

// NewVAIDecAccIncrementTargetPosOnRisingCommand creates a VAI Dec=Acc Increment Target Pos On Rising Trigger Event command (0CBxh).
func NewVAIDecAccIncrementTargetPosOnRisingCommand(increment int32, velocity, accel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.IncrementTargetPosOnRising),
		parameters,
	)
}

// NewVAIDecAccGoToPosOnFallingTriggerCommand creates a VAI Dec=Acc Go To Pos On Falling Trigger Event command (0CCxh).
func NewVAIDecAccGoToPosOnFallingTriggerCommand(position int32, velocity, accel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.GoToPosOnFallingTrigger),
		parameters,
	)
}

// NewVAIDecAccIncrementTargetPosOnFallingCommand creates a VAI Dec=Acc Increment Target Pos On Falling Trigger Event command (0CDxh).
func NewVAIDecAccIncrementTargetPosOnFallingCommand(increment int32, velocity, accel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)
	parameters[4] = uint16(accel & 0xFFFF)
	parameters[5] = uint16((accel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.IncrementTargetPosOnFalling),
		parameters,
	)
}

// NewVAIDecAccInfiniteMotionPositiveCommand creates a VAI Dec=Acc Infinite Motion Positive Direction command (0CExh).
func NewVAIDecAccInfiniteMotionPositiveCommand(velocity, accel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(velocity & 0xFFFF)
	parameters[1] = uint16((velocity >> 16) & 0xFFFF)
	parameters[2] = uint16(accel & 0xFFFF)
	parameters[3] = uint16((accel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.InfiniteMotionPositive),
		parameters,
	)
}

// NewVAIDecAccInfiniteMotionNegativeCommand creates a VAI Dec=Acc Infinite Motion Negative Direction command (0CFxh).
func NewVAIDecAccInfiniteMotionNegativeCommand(velocity, accel uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(velocity & 0xFFFF)
	parameters[1] = uint16((velocity >> 16) & 0xFFFF)
	parameters[2] = uint16(accel & 0xFFFF)
	parameters[3] = uint16((accel >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIDecEqualsAcc,
		uint8(SubIDs.InfiniteMotionNegative),
		parameters,
	)
}
