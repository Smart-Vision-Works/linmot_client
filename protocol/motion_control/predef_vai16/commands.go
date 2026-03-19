package protocol_predef_vai16

import (
	protocol_motion_control "github.com/Smart-Vision-Works/linmot_client/protocol/motion_control"
)

// ============================================================================
// Predefined 16-Bit VAI Command Constructors
// ============================================================================
// These commands combine drive-configured motion parameters with 16-bit position values.

// NewPredefVAI16GoToPosCommand creates a Predefined 16-bit VAI Go To Position command (0A0xh).
func NewPredefVAI16GoToPosCommand(position int16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)
	parameters[0] = uint16(position)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI16Bit,
		uint8(SubIDs.GoToPos),
		parameters,
	)
}

// NewPredefVAI16IncrementDemPosCommand creates a Predefined 16-bit VAI Increment Demand Position command (0A1xh).
func NewPredefVAI16IncrementDemPosCommand(increment int16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)
	parameters[0] = uint16(increment)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI16Bit,
		uint8(SubIDs.IncrementDemPos),
		parameters,
	)
}

// NewPredefVAI16IncrementTargetPosCommand creates a Predefined 16-bit VAI Increment Target Position command (0A2xh).
func NewPredefVAI16IncrementTargetPosCommand(increment int16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)
	parameters[0] = uint16(increment)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI16Bit,
		uint8(SubIDs.IncrementTargetPos),
		parameters,
	)
}

// NewPredefVAI16GoToPosFromActPosCommand creates a Predefined 16-bit VAI Go To Pos From Act Pos And Act Vel command (0A3xh).
func NewPredefVAI16GoToPosFromActPosCommand(position int16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)
	parameters[0] = uint16(position)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI16Bit,
		uint8(SubIDs.GoToPosFromActPosAndActVel),
		parameters,
	)
}

// NewPredefVAI16GoToPosFromActPosDemVelZeroCommand creates a Predefined 16-bit VAI Go To Pos Starting With Dem Vel = 0 command (0A4xh).
func NewPredefVAI16GoToPosFromActPosDemVelZeroCommand(position int16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)
	parameters[0] = uint16(position)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI16Bit,
		uint8(SubIDs.GoToPosFromActPosDemVelZero),
		parameters,
	)
}

// NewPredefVAI16StopCommand creates a Predefined 16-bit VAI Stop command (0A7xh).
func NewPredefVAI16StopCommand() *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI16Bit,
		uint8(SubIDs.Stop),
		parameters,
	)
}

// NewPredefVAI16GoToPosAfterActualCommandCommand creates a Predefined 16-bit VAI Go To Pos After Actual Command (0A8xh).
func NewPredefVAI16GoToPosAfterActualCommandCommand(position int16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)
	parameters[0] = uint16(position)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI16Bit,
		uint8(SubIDs.GoToPosAfterActualCommand),
		parameters,
	)
}

// NewPredefVAI16GoToPosOnRisingTriggerCommand creates a Predefined 16-bit VAI Go To Pos On Rising Trigger Event command (0AAxh).
func NewPredefVAI16GoToPosOnRisingTriggerCommand(position int16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)
	parameters[0] = uint16(position)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI16Bit,
		uint8(SubIDs.GoToPosOnRisingTrigger),
		parameters,
	)
}

// NewPredefVAI16IncrementTargetPosOnRisingCommand creates a Predefined 16-bit VAI Increment Target Pos On Rising Trigger Event command (0ABxh).
func NewPredefVAI16IncrementTargetPosOnRisingCommand(increment int16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)
	parameters[0] = uint16(increment)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI16Bit,
		uint8(SubIDs.IncrementTargetPosOnRising),
		parameters,
	)
}

// NewPredefVAI16GoToPosOnFallingTriggerCommand creates a Predefined 16-bit VAI Go To Pos On Falling Trigger Event command (0ACxh).
func NewPredefVAI16GoToPosOnFallingTriggerCommand(position int16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)
	parameters[0] = uint16(position)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI16Bit,
		uint8(SubIDs.GoToPosOnFallingTrigger),
		parameters,
	)
}

// NewPredefVAI16IncrementTargetPosOnFallingCommand creates a Predefined 16-bit VAI Increment Target Pos On Falling Trigger Event command (0ADxh).
func NewPredefVAI16IncrementTargetPosOnFallingCommand(increment int16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)
	parameters[0] = uint16(increment)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.PredefVAI16Bit,
		uint8(SubIDs.IncrementTargetPosOnFalling),
		parameters,
	)
}
