package protocol_predefined

import (
	protocol_motion_control "gsail-go/linmot/protocol/motion_control"
)

// ============================================================================
// Predefined VAI Command Constructors
// ============================================================================
// These commands use drive-configured motion parameters (velocity, accel, decel).

// NewPredefVAIGoToPosCommand creates a Predefined VAI Go To Position command (020xh).
// Uses drive-configured velocity, acceleration, and deceleration.
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.25
func NewPredefVAIGoToPosCommand(position int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	// Position (int32 → 2 words, little-endian)
	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI,
		uint8(SubIDs.GoToPos),
		parameters,
	)
}

// NewPredefVAIIncrementDemPosCommand creates a Predefined VAI Increment Demand Position command (021xh).
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.26
func NewPredefVAIIncrementDemPosCommand(increment int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI,
		uint8(SubIDs.IncrementDemPos),
		parameters,
	)
}

// NewPredefVAIIncrementTargetPosCommand creates a Predefined VAI Increment Target Position command (022xh).
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.27
func NewPredefVAIIncrementTargetPosCommand(increment int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI,
		uint8(SubIDs.IncrementTargetPos),
		parameters,
	)
}

// NewPredefVAIGoToPosFromActPosCommand creates a Predefined VAI Go To Pos From Act Pos And Act Vel command (023xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.28
func NewPredefVAIGoToPosFromActPosCommand(position int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI,
		uint8(SubIDs.GoToPosFromActPosAndActVel),
		parameters,
	)
}

// NewPredefVAIGoToPosFromActPosDemVelZeroCommand creates a Predefined VAI Go To Pos Starting With Dem Vel = 0 command (024xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.29
func NewPredefVAIGoToPosFromActPosDemVelZeroCommand(position int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI,
		uint8(SubIDs.GoToPosFromActPosDemVelZero),
		parameters,
	)
}

// NewPredefVAIStopCommand creates a Predefined VAI Stop command (027xh).
// Uses drive-configured quick stop deceleration.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.30
func NewPredefVAIStopCommand() *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI,
		uint8(SubIDs.Stop),
		parameters,
	)
}

// NewPredefVAIGoToPosAfterActualCommandCommand creates a Predefined VAI Go To Pos After Actual Command (028xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.31
func NewPredefVAIGoToPosAfterActualCommandCommand(position int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI,
		uint8(SubIDs.GoToPosAfterActualCommand),
		parameters,
	)
}

// NewPredefVAIGoToPosOnRisingTriggerCommand creates a Predefined VAI Go To Pos On Rising Trigger Event command (02Axh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.32
func NewPredefVAIGoToPosOnRisingTriggerCommand(position int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI,
		uint8(SubIDs.GoToPosOnRisingTrigger),
		parameters,
	)
}

// NewPredefVAIIncrementTargetPosOnRisingCommand creates a Predefined VAI Increment Target Pos On Rising Trigger Event command (02Bxh).
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.33
func NewPredefVAIIncrementTargetPosOnRisingCommand(increment int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI,
		uint8(SubIDs.IncrementTargetPosOnRising),
		parameters,
	)
}

// NewPredefVAIGoToPosOnFallingTriggerCommand creates a Predefined VAI Go To Pos On Falling Trigger Event command (02Cxh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.34
func NewPredefVAIGoToPosOnFallingTriggerCommand(position int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI,
		uint8(SubIDs.GoToPosOnFallingTrigger),
		parameters,
	)
}

// NewPredefVAIIncrementTargetPosOnFallingCommand creates a Predefined VAI Increment Target Pos On Falling Trigger Event command (02Dxh).
//
// Parameters:
//   - increment: Position increment in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.35
func NewPredefVAIIncrementTargetPosOnFallingCommand(increment int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI,
		uint8(SubIDs.IncrementTargetPosOnFalling),
		parameters,
	)
}

// NewPredefVAIInfiniteMotionPositiveCommand creates a Predefined VAI Infinite Motion Positive Direction command (02Exh).
// Starts continuous motion in positive direction using drive-configured parameters.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.36
func NewPredefVAIInfiniteMotionPositiveCommand() *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI,
		uint8(SubIDs.InfiniteMotionPositive),
		parameters,
	)
}

// NewPredefVAIInfiniteMotionNegativeCommand creates a Predefined VAI Infinite Motion Negative Direction command (02Fxh).
// Starts continuous motion in negative direction using drive-configured parameters.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.37
func NewPredefVAIInfiniteMotionNegativeCommand() *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI,
		uint8(SubIDs.InfiniteMotionNegative),
		parameters,
	)
}
