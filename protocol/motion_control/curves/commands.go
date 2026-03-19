package protocol_curves

import (
	protocol_motion_control "github.com/Smart-Vision-Works/linmot_client/protocol/motion_control"
)

// ============================================================================
// Time Curve Command Constructors
// ============================================================================

// NewTimeCurveDefaultParamsCommand creates a Time Curve With Default Parameters command (040xh).
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.45
func NewTimeCurveDefaultParamsCommand() *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		uint8(SubIDs.TimeCurveDefaultParams),
		parameters,
	)
}

// NewTimeCurveDefaultParamsFromActPosCommand creates a Time Curve With Default Parameters From Act Pos command (041xh).
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.46
func NewTimeCurveDefaultParamsFromActPosCommand() *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		uint8(SubIDs.TimeCurveDefaultParamsFromActPos),
		parameters,
	)
}

// NewTimeCurveToPosDefaultSpeedCommand creates a Time Curve To Pos With Default Speed command (042xh).
//
// Parameters:
//   - targetPosition: Target position in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.47
func NewTimeCurveToPosDefaultSpeedCommand(targetPosition int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(targetPosition) & 0xFFFF)
	parameters[1] = uint16((uint32(targetPosition) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		uint8(SubIDs.TimeCurveToPosDefaultSpeed),
		parameters,
	)
}

// NewTimeCurveToPosAdjustableTimeCommand creates a Time Curve To Pos With Adjustable Time command (043xh).
//
// Parameters:
//   - targetPosition: Target position in 0.1µm units (int32)
//   - time: Time to reach position in milliseconds (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.48
func NewTimeCurveToPosAdjustableTimeCommand(targetPosition int32, time uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(targetPosition) & 0xFFFF)
	parameters[1] = uint16((uint32(targetPosition) >> 16) & 0xFFFF)
	parameters[2] = uint16(time & 0xFFFF)
	parameters[3] = uint16((time >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		uint8(SubIDs.TimeCurveToPosAdjustableTime),
		parameters,
	)
}

// NewTimeCurveOffsetTimeAmpScaleCommand creates a Time Curve With Adjustable Offset, Time & Amplitude Scale command (045xh).
//
// Parameters:
//   - offset: Position offset in 0.1µm units (int32)
//   - timeScale: Time scale factor (uint16)
//   - amplitudeScale: Amplitude scale factor (uint16)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.50
func NewTimeCurveOffsetTimeAmpScaleCommand(offset int32, timeScale, amplitudeScale uint16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(offset) & 0xFFFF)
	parameters[1] = uint16((uint32(offset) >> 16) & 0xFFFF)
	parameters[2] = timeScale
	parameters[3] = amplitudeScale

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		uint8(SubIDs.TimeCurveOffsetTimeAmpScale),
		parameters,
	)
}

// NewTimeCurveOffsetTimeAmpScaleRisingCommand creates a Time Curve With Adjustable Offset, Time & Amplitude Scale On Rising Trigger Event command (046xh).
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.51
func NewTimeCurveOffsetTimeAmpScaleRisingCommand(offset int32, timeScale, amplitudeScale uint16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(offset) & 0xFFFF)
	parameters[1] = uint16((uint32(offset) >> 16) & 0xFFFF)
	parameters[2] = timeScale
	parameters[3] = amplitudeScale

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		uint8(SubIDs.TimeCurveOffsetTimeAmpScaleRising),
		parameters,
	)
}

// NewTimeCurveOffsetTimeAmpScaleFallingCommand creates a Time Curve With Adjustable Offset, Time & Amplitude Scale On Falling Trigger Event command (047xh).
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.52
func NewTimeCurveOffsetTimeAmpScaleFallingCommand(offset int32, timeScale, amplitudeScale uint16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(offset) & 0xFFFF)
	parameters[1] = uint16((uint32(offset) >> 16) & 0xFFFF)
	parameters[2] = timeScale
	parameters[3] = amplitudeScale

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		uint8(SubIDs.TimeCurveOffsetTimeAmpScaleFalling),
		parameters,
	)
}

// NewTimeCurveToPosDefaultSpeedRisingCommand creates a Time Curve To Pos With Default Speed On Rising Trigger Event command (04Axh).
//
// Parameters:
//   - targetPosition: Target position in 0.1µm units (int32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.53
func NewTimeCurveToPosDefaultSpeedRisingCommand(targetPosition int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(targetPosition) & 0xFFFF)
	parameters[1] = uint16((uint32(targetPosition) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		uint8(SubIDs.TimeCurveToPosDefaultSpeedRising),
		parameters,
	)
}

// NewTimeCurveToPosDefaultSpeedFallingCommand creates a Time Curve To Pos With Default Speed On Falling Trigger Event command (04Cxh).
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.54
func NewTimeCurveToPosDefaultSpeedFallingCommand(targetPosition int32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(targetPosition) & 0xFFFF)
	parameters[1] = uint16((uint32(targetPosition) >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		uint8(SubIDs.TimeCurveToPosDefaultSpeedFalling),
		parameters,
	)
}

// NewTimeCurveToPosAdjustableTimeRisingCommand creates a Time Curve To Pos With Adjustable Time On Rising Trigger Event command (04Exh).
//
// Parameters:
//   - targetPosition: Target position in 0.1µm units (int32)
//   - time: Time to reach position in milliseconds (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.55
func NewTimeCurveToPosAdjustableTimeRisingCommand(targetPosition int32, time uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(targetPosition) & 0xFFFF)
	parameters[1] = uint16((uint32(targetPosition) >> 16) & 0xFFFF)
	parameters[2] = uint16(time & 0xFFFF)
	parameters[3] = uint16((time >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		uint8(SubIDs.TimeCurveToPosAdjustableTimeRising),
		parameters,
	)
}

// NewTimeCurveToPosAdjustableTimeFallingCommand creates a Time Curve To Pos With Adjustable Time On Falling Trigger Event command (04Fxh).
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.56
func NewTimeCurveToPosAdjustableTimeFallingCommand(targetPosition int32, time uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = uint16(uint32(targetPosition) & 0xFFFF)
	parameters[1] = uint16((uint32(targetPosition) >> 16) & 0xFFFF)
	parameters[2] = uint16(time & 0xFFFF)
	parameters[3] = uint16((time >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		uint8(SubIDs.TimeCurveToPosAdjustableTimeFalling),
		parameters,
	)
}

// ============================================================================
// Curve Modification Commands (Sub-group 0x50-0x56)
// ============================================================================

// NewModifyCurveStartAddressCommand creates a Modify Curve Start Address in RAM command (050xh).
//
// Parameters:
//   - curveID: Curve identifier (uint16)
//   - startAddress: New start address (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.57
func NewModifyCurveStartAddressCommand(curveID uint16, startAddress uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = curveID
	parameters[1] = uint16(startAddress & 0xFFFF)
	parameters[2] = uint16((startAddress >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		0x0, // Sub ID 0x0 in the 0x50 sub-group
		parameters,
	)
}

// NewModifyCurveInfoBlock16BitCommand creates a Modify Curve Info Block 16 Bit Value in RAM command (051xh).
//
// Parameters:
//   - curveID: Curve identifier (uint16)
//   - offset: Byte offset in info block (uint16)
//   - value: 16-bit value to write (uint16)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.58
func NewModifyCurveInfoBlock16BitCommand(curveID, offset, value uint16) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = curveID
	parameters[1] = offset
	parameters[2] = value

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		0x1, // Sub ID 0x1 in the 0x50 sub-group
		parameters,
	)
}

// NewModifyCurveInfoBlock32BitCommand creates a Modify Curve Info Block 32 Bit Value in RAM command (052xh).
//
// Parameters:
//   - curveID: Curve identifier (uint16)
//   - offset: Byte offset in info block (uint16)
//   - value: 32-bit value to write (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.59
func NewModifyCurveInfoBlock32BitCommand(curveID, offset uint16, value uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = curveID
	parameters[1] = offset
	parameters[2] = uint16(value & 0xFFFF)
	parameters[3] = uint16((value >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		0x2, // Sub ID 0x2 in the 0x50 sub-group
		parameters,
	)
}

// NewModifyCurveDataBlock32BitCommand creates a Modify Curve Data Block 32 Bit Value in RAM command (054xh).
//
// Parameters:
//   - curveID: Curve identifier (uint16)
//   - offset: Byte offset in data block (uint16)
//   - value: 32-bit value to write (uint32)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.60
func NewModifyCurveDataBlock32BitCommand(curveID, offset uint16, value uint32) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = curveID
	parameters[1] = offset
	parameters[2] = uint16(value & 0xFFFF)
	parameters[3] = uint16((value >> 16) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		0x4, // Sub ID 0x4 in the 0x50 sub-group
		parameters,
	)
}

// NewModifyCurveDataBlock64BitCommand creates a Modify Curve Data Block 64 Bit Value in RAM command (055xh).
//
// Parameters:
//   - curveID: Curve identifier (uint16)
//   - offset: Byte offset in data block (uint16)
//   - value: 64-bit value to write (uint64)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.61
func NewModifyCurveDataBlock64BitCommand(curveID, offset uint16, value uint64) *protocol_motion_control.MCCommandRequest {
	parameters := make([]uint16, 15)

	parameters[0] = curveID
	parameters[1] = offset
	parameters[2] = uint16(value & 0xFFFF)
	parameters[3] = uint16((value >> 16) & 0xFFFF)
	parameters[4] = uint16((value >> 32) & 0xFFFF)
	parameters[5] = uint16((value >> 48) & 0xFFFF)

	return protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterIDs.Curve,
		0x5, // Sub ID 0x5 in the 0x50 sub-group
		parameters,
	)
}
