package protocol_curves

// ============================================================================
// Time Curve Sub ID Constants
// ============================================================================

// SubID represents a Time Curve command identifier.
type SubID uint8

// SubIDs groups all Time Curve command Sub ID constants.
// Reference: LinMot_MotionCtrl.txt, Section 4.3.45-4.3.62
var SubIDs = struct {
	TimeCurveDefaultParams              SubID // 0x0 - Time Curve With Default Parameters (040xh)
	TimeCurveDefaultParamsFromActPos    SubID // 0x1 - Time Curve With Default Parameters From Act Pos (041xh)
	TimeCurveToPosDefaultSpeed          SubID // 0x2 - Time Curve To Pos With Default Speed (042xh)
	TimeCurveToPosAdjustableTime        SubID // 0x3 - Time Curve To Pos With Adjustable Time (043xh)
	TimeCurveOffsetTimeScaleAmpScale    SubID // 0x4 - Time Curve With Adjustable Offset, Time Scale & Amplitude Scale (044xh)
	TimeCurveOffsetTimeAmpScale         SubID // 0x5 - Time Curve With Adjustable Offset, Time & Amplitude Scale (045xh)
	TimeCurveOffsetTimeAmpScaleRising   SubID // 0x6 - Time Curve With Adjustable Offset, Time & Amplitude Scale On Rising Trigger Event (046xh)
	TimeCurveOffsetTimeAmpScaleFalling  SubID // 0x7 - Time Curve With Adjustable Offset, Time & Amplitude Scale On Falling Trigger Event (047xh)
	TimeCurveToPosDefaultSpeedRising    SubID // 0xA - Time Curve To Pos With Default Speed On Rising Trigger Event (04Axh)
	TimeCurveToPosDefaultSpeedFalling   SubID // 0xC - Time Curve To Pos With Default Speed On Falling Trigger Event (04Cxh)
	TimeCurveToPosAdjustableTimeRising  SubID // 0xE - Time Curve To Pos With Adjustable Time On Rising Trigger Event (04Exh)
	TimeCurveToPosAdjustableTimeFalling SubID // 0xF - Time Curve To Pos With Adjustable Time On Falling Trigger Event (04Fxh)
	ModifyCurveStartAddress             SubID // 0x0 in sub-group - Modify Curve Start Address in RAM (050xh)
	ModifyCurveInfoBlock16Bit           SubID // 0x1 in sub-group - Modify Curve Info Block 16 Bit Value in RAM (051xh)
	ModifyCurveInfoBlock32Bit           SubID // 0x2 in sub-group - Modify Curve Info Block 32 Bit Value in RAM (052xh)
	ModifyCurveDataBlock32Bit           SubID // 0x4 in sub-group - Modify Curve Data Block 32 Bit Value in RAM (054xh)
	ModifyCurveDataBlock64Bit           SubID // 0x5 in sub-group - Modify Curve Data Block 64 Bit Value in RAM (055xh)
	ModifyCurveDataBlock96Bit           SubID // 0x6 in sub-group - Modify Curve Data Block 96 Bit Value in RAM (056xh)
}{
	TimeCurveDefaultParams:              0x0,
	TimeCurveDefaultParamsFromActPos:    0x1,
	TimeCurveToPosDefaultSpeed:          0x2,
	TimeCurveToPosAdjustableTime:        0x3,
	TimeCurveOffsetTimeScaleAmpScale:    0x4,
	TimeCurveOffsetTimeAmpScale:         0x5,
	TimeCurveOffsetTimeAmpScaleRising:   0x6,
	TimeCurveOffsetTimeAmpScaleFalling:  0x7,
	TimeCurveToPosDefaultSpeedRising:    0xA,
	TimeCurveToPosDefaultSpeedFalling:   0xC,
	TimeCurveToPosAdjustableTimeRising:  0xE,
	TimeCurveToPosAdjustableTimeFalling: 0xF,
	ModifyCurveStartAddress:             0x0, // Note: Different sub-group (0x50-0x56)
	ModifyCurveInfoBlock16Bit:           0x1,
	ModifyCurveInfoBlock32Bit:           0x2,
	ModifyCurveDataBlock32Bit:           0x4,
	ModifyCurveDataBlock64Bit:           0x5,
	ModifyCurveDataBlock96Bit:           0x6,
}
