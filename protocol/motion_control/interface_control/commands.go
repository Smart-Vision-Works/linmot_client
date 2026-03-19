package protocol_interface_control

import (
	protocol_motion_control "gsail-go/linmot/protocol/motion_control"
)

// ============================================================================
// Interface Control Command Constructors
// ============================================================================

// NewNoOperationCommand creates a No Operation command (000xh).
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.1
func NewNoOperationCommand() *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.InterfaceControl,
		uint8(SubIDs.NoOperation),
		parameters,
	)
}

// NewWriteInterfaceControlWordCommand creates a Write Interface Control Word command (001xh).
//
// Parameters:
//   - controlWord: Interface control word value (uint16)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.2
func NewWriteInterfaceControlWordCommand(controlWord uint16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)
	parameters[0] = controlWord

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.InterfaceControl,
		uint8(SubIDs.WriteInterfaceControlWord),
		parameters,
	)
}

// NewWriteLiveParameterCommand creates a Write Live Parameter command (002xh).
//
// Parameters:
//   - upid: Unique Parameter ID (uint16)
//   - value: Parameter value (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.3
func NewWriteLiveParameterCommand(upid uint16, value uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	// UPID (1 word)
	parameters[0] = upid

	// Value (uint32 → 2 words)
	parameters[1] = uint16(value & 0xFFFF)
	parameters[2] = uint16((value >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.InterfaceControl,
		uint8(SubIDs.WriteLiveParameter),
		parameters,
	)
}

// NewWriteX4IntfOutputsWithMaskCommand creates a Write X4/X14 Interface Outputs with Mask command (003xh).
//
// Parameters:
//   - bitMask: Bit mask for which outputs to write (uint16)
//   - bitValue: Bit values for the outputs (uint16)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.4
func NewWriteX4IntfOutputsWithMaskCommand(bitMask, bitValue uint16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = bitMask
	parameters[1] = bitValue

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.InterfaceControl,
		uint8(SubIDs.WriteX4IntfOutputsWithMask),
		parameters,
	)
}

// NewSelectPositionControllerSetCommand creates a Select Position Controller Set command (005xh).
//
// Parameters:
//   - controllerSet: Controller set selection (0 = Set A, 1 = Set B)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.5
func NewSelectPositionControllerSetCommand(controllerSet uint16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)
	parameters[0] = controllerSet

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.InterfaceControl,
		uint8(SubIDs.SelectPositionControllerSet),
		parameters,
	)
}

// NewClearEventEvaluationCommand creates a Clear Event Evaluation command (008xh).
// Resets the event handler used for trigger-based commands.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.6
func NewClearEventEvaluationCommand() *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.InterfaceControl,
		uint8(SubIDs.ClearEventEvaluation),
		parameters,
	)
}

// NewMasterHomingCommand creates a Master Homing command (009xh).
//
// Parameters:
//   - homePosition: Home position in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.7
func NewMasterHomingCommand(homePosition int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	// Home Position (int32 → 2 words)
	parameters[0] = uint16(uint32(homePosition) & 0xFFFF)
	parameters[1] = uint16((uint32(homePosition) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.InterfaceControl,
		uint8(SubIDs.MasterHoming),
		parameters,
	)
}

// NewResetCommand creates a Reset command (00Fxh).
// Resets all firmware instances of the drive.
// IMPORTANT: Use with counter = 0, otherwise the drive reboots cyclically.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.8
func NewResetCommand() *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.InterfaceControl,
		uint8(SubIDs.Reset),
		parameters,
	)
}
