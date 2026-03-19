package protocol_encoder_cams

import (
	protocol_motion_control "github.com/Smart-Vision-Works/linmot_client/protocol/motion_control"
)

// ============================================================================
// Encoder CAM Command Constructors
// ============================================================================

// NewSetupOnRisingWithDelayCommand creates a Setup Encoder Cam On Rising Trigger Event With Delay Counts command (069xh).
//
// Parameters:
//   - delayCounts: Delay in encoder counts (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.63
func NewSetupOnRisingWithDelayCommand(delayCounts uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(delayCounts & 0xFFFF)
	parameters[1] = uint16((delayCounts >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.EncoderCam,
		uint8(SubIDs.SetupOnRisingWithDelay),
		parameters,
	)
}

// NewSetupOnRisingWithDelayTargetPosLenCommand creates a Setup Encoder Cam On Rising Trigger Event With Delay Counts, Target Pos and Length command (06Axh).
//
// Parameters:
//   - delayCounts: Delay in encoder counts (uint32)
//   - targetPosition: Target position in 0.1µm units (int32)
//   - length: Cam length in encoder counts (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.64
func NewSetupOnRisingWithDelayTargetPosLenCommand(delayCounts uint32, targetPosition int32, length uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(delayCounts & 0xFFFF)
	parameters[1] = uint16((delayCounts >> 16) & 0xFFFF)
	parameters[2] = uint16(uint32(targetPosition) & 0xFFFF)
	parameters[3] = uint16((uint32(targetPosition) >> 16) & 0xFFFF)
	parameters[4] = uint16(length & 0xFFFF)
	parameters[5] = uint16((length >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.EncoderCam,
		uint8(SubIDs.SetupOnRisingWithDelayTargetPosLen),
		parameters,
	)
}

// NewSetupOnFallingWithDelayCommand creates a Setup Encoder Cam On Falling Trigger Event With Delay Counts command (06Bxh).
//
// Parameters:
//   - delayCounts: Delay in encoder counts (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.65
func NewSetupOnFallingWithDelayCommand(delayCounts uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(delayCounts & 0xFFFF)
	parameters[1] = uint16((delayCounts >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.EncoderCam,
		uint8(SubIDs.SetupOnFallingWithDelay),
		parameters,
	)
}

// NewSetupOnFallingWithDelayTargetPosLenCommand creates a Setup Encoder Cam On Falling Trigger Event With Delay Counts, Target Pos and Length command (06Cxh).
//
// Parameters:
//   - delayCounts: Delay in encoder counts (uint32)
//   - targetPosition: Target position in 0.1µm units (int32)
//   - length: Cam length in encoder counts (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.66
func NewSetupOnFallingWithDelayTargetPosLenCommand(delayCounts uint32, targetPosition int32, length uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(delayCounts & 0xFFFF)
	parameters[1] = uint16((delayCounts >> 16) & 0xFFFF)
	parameters[2] = uint16(uint32(targetPosition) & 0xFFFF)
	parameters[3] = uint16((uint32(targetPosition) >> 16) & 0xFFFF)
	parameters[4] = uint16(length & 0xFFFF)
	parameters[5] = uint16((length >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.EncoderCam,
		uint8(SubIDs.SetupOnFallingWithDelayTargetPosLen),
		parameters,
	)
}

// NewSetupOnRisingWithDelayAmpScaleLenCommand creates a Setup Encoder Cam On Rising Trigger Event With Delay Counts, Amplitude scale and Length command (06Dxh).
//
// Parameters:
//   - delayCounts: Delay in encoder counts (uint32)
//   - amplitudeScale: Amplitude scale factor (uint16)
//   - length: Cam length in encoder counts (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.67
func NewSetupOnRisingWithDelayAmpScaleLenCommand(delayCounts uint32, amplitudeScale uint16, length uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(delayCounts & 0xFFFF)
	parameters[1] = uint16((delayCounts >> 16) & 0xFFFF)
	parameters[2] = amplitudeScale
	parameters[3] = uint16(length & 0xFFFF)
	parameters[4] = uint16((length >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.EncoderCam,
		uint8(SubIDs.SetupOnRisingWithDelayAmpScaleLen),
		parameters,
	)
}

// NewSetupOnFallingWithDelayAmpScaleLenCommand creates a Setup Encoder Cam On Falling Trigger Event With Delay Counts, Amplitude scale and Length command (06Exh).
//
// Parameters:
//   - delayCounts: Delay in encoder counts (uint32)
//   - amplitudeScale: Amplitude scale factor (uint16)
//   - length: Cam length in encoder counts (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.68
func NewSetupOnFallingWithDelayAmpScaleLenCommand(delayCounts uint32, amplitudeScale uint16, length uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(delayCounts & 0xFFFF)
	parameters[1] = uint16((delayCounts >> 16) & 0xFFFF)
	parameters[2] = amplitudeScale
	parameters[3] = uint16(length & 0xFFFF)
	parameters[4] = uint16((length >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.EncoderCam,
		uint8(SubIDs.SetupOnFallingWithDelayAmpScaleLen),
		parameters,
	)
}
