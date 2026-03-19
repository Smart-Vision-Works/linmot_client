package protocol_streaming

import (
	protocol_motion_control "gsail-go/linmot/protocol/motion_control"
)

// ============================================================================
// Streaming Command Constructors
// ============================================================================

// NewPStreamCommand creates a Position Stream With Slave Generated Time Stamp command (030xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.38
func NewPStreamCommand(position int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Streaming,
		uint8(SubIDs.PStreamSlaveTimestamp),
		parameters,
	)
}

// NewPVStreamCommand creates a PV Stream With Slave Generated Time Stamp command (031xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Velocity in 1E-6 m/s units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.39
func NewPVStreamCommand(position, velocity int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(uint32(velocity) & 0xFFFF)
	parameters[3] = uint16((uint32(velocity) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Streaming,
		uint8(SubIDs.PVStreamSlaveTimestamp),
		parameters,
	)
}

// NewPStreamConfigPeriodCommand creates a P Stream With Configured Period Time command (032xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.40
func NewPStreamConfigPeriodCommand(position int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Streaming,
		uint8(SubIDs.PStreamSlaveTimestampConfigPeriod),
		parameters,
	)
}

// NewPVStreamConfigPeriodCommand creates a PV Stream With Configured Period Time command (033xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Velocity in 1E-6 m/s units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.41
func NewPVStreamConfigPeriodCommand(position, velocity int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(uint32(velocity) & 0xFFFF)
	parameters[3] = uint16((uint32(velocity) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Streaming,
		uint8(SubIDs.PVStreamSlaveTimestampConfigPeriod),
		parameters,
	)
}

// NewPVAStreamCommand creates a PVA Stream With Slave Generated Time Stamp command (034xh).
//
// Parameters:
//   - position: Target position in 0.1µm units (int32)
//   - velocity: Velocity in 1E-6 m/s units (int32)
//   - acceleration: Acceleration in 1E-5 m/s² units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.42
func NewPVAStreamCommand(position, velocity, acceleration int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(position) & 0xFFFF)
	parameters[1] = uint16((uint32(position) >> 16) & 0xFFFF)
	parameters[2] = uint16(uint32(velocity) & 0xFFFF)
	parameters[3] = uint16((uint32(velocity) >> 16) & 0xFFFF)
	parameters[4] = uint16(uint32(acceleration) & 0xFFFF)
	parameters[5] = uint16((uint32(acceleration) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Streaming,
		uint8(SubIDs.PVAStreamSlaveTimestamp),
		parameters,
	)
}

// NewStopStreamingCommand creates a Stop Streaming command (03Fxh).
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.43
func NewStopStreamingCommand() *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Streaming,
		uint8(SubIDs.StopStreaming),
		parameters,
	)
}
