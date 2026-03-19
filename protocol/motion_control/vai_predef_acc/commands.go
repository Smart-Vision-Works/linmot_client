package protocol_vai_predef_acc

import (
	protocol_motion_control "github.com/Smart-Vision-Works/linmot_client/protocol/motion_control"
)

// ============================================================================
// VAI Predefined Acceleration Command Constructors
// ============================================================================
// These commands use drive-configured acceleration, but allow variable velocity and deceleration.

// NewVAIPredefAccGoToPosCommand creates a VAI Predefined Acc Go To Position command (0B0xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Maximal velocity in 1E-6 m/s units (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.99
func NewVAIPredefAccGoToPosCommand(position int32, velocity uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	// Position (int32 → 2 words)
	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)

	// Velocity (uint32 → 2 words)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPredefAcc,
		uint8(SubIDs.GoToPos),
		parameters,
	)
}

// NewVAIPredefAccIncrementDemPosCommand creates a VAI Predefined Acc Increment Demand Position command (0B1xh).
func NewVAIPredefAccIncrementDemPosCommand(increment int32, velocity uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPredefAcc,
		uint8(SubIDs.IncrementDemPos),
		parameters,
	)
}

// NewVAIPredefAccIncrementTargetPosCommand creates a VAI Predefined Acc Increment Target Position command (0B2xh).
func NewVAIPredefAccIncrementTargetPosCommand(increment int32, velocity uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPredefAcc,
		uint8(SubIDs.IncrementTargetPos),
		parameters,
	)
}

// NewVAIPredefAccGoToPosFromActPosCommand creates a VAI Predefined Acc Go To Pos From Act Pos And Act Vel command (0B3xh).
func NewVAIPredefAccGoToPosFromActPosCommand(position int32, velocity uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPredefAcc,
		uint8(SubIDs.GoToPosFromActPosAndActVel),
		parameters,
	)
}

// NewVAIPredefAccGoToPosFromActPosDemVelZeroCommand creates a VAI Predefined Acc Go To Pos Starting With Dem Vel = 0 command (0B4xh).
func NewVAIPredefAccGoToPosFromActPosDemVelZeroCommand(position int32, velocity uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPredefAcc,
		uint8(SubIDs.GoToPosFromActPosDemVelZero),
		parameters,
	)
}

// NewVAIPredefAccGoToPosAfterActualCommandCommand creates a VAI Predefined Acc Go To Pos After Actual Command (0B8xh).
func NewVAIPredefAccGoToPosAfterActualCommandCommand(position int32, velocity uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPredefAcc,
		uint8(SubIDs.GoToPosAfterActualCommand),
		parameters,
	)
}

// NewVAIPredefAccGoToPosOnRisingTriggerCommand creates a VAI Predefined Acc Go To Pos On Rising Trigger Event command (0BAxh).
func NewVAIPredefAccGoToPosOnRisingTriggerCommand(position int32, velocity uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPredefAcc,
		uint8(SubIDs.GoToPosOnRisingTrigger),
		parameters,
	)
}

// NewVAIPredefAccIncrementTargetPosOnRisingCommand creates a VAI Predefined Acc Increment Target Pos On Rising Trigger Event command (0BBxh).
func NewVAIPredefAccIncrementTargetPosOnRisingCommand(increment int32, velocity uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPredefAcc,
		uint8(SubIDs.IncrementTargetPosOnRising),
		parameters,
	)
}

// NewVAIPredefAccGoToPosOnFallingTriggerCommand creates a VAI Predefined Acc Go To Pos On Falling Trigger Event command (0BCxh).
func NewVAIPredefAccGoToPosOnFallingTriggerCommand(position int32, velocity uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPredefAcc,
		uint8(SubIDs.GoToPosOnFallingTrigger),
		parameters,
	)
}

// NewVAIPredefAccIncrementTargetPosOnFallingCommand creates a VAI Predefined Acc Increment Target Pos On Falling Trigger Event command (0BDxh).
func NewVAIPredefAccIncrementTargetPosOnFallingCommand(increment int32, velocity uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(increment) & 0xFFFF)
	parameters[1] = uint16((uint32(increment) >> 16) & 0xFFFF)
	parameters[2] = uint16(velocity & 0xFFFF)
	parameters[3] = uint16((velocity >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.VAIPredefAcc,
		uint8(SubIDs.IncrementTargetPosOnFalling),
		parameters,
	)
}
