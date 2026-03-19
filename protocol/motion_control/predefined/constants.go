package protocol_predefined

// ============================================================================
// Predefined VAI Sub ID Constants
// ============================================================================

// SubID represents a Predefined VAI command identifier.
type SubID uint8

// SubIDs groups all Predefined VAI command Sub ID constants.
// These commands use drive-configured motion parameters (velocity, accel, decel).
// Reference: LinMot_MotionCtrl.txt, Section 4.3.25-4.3.37
var SubIDs = struct {
	GoToPos                     SubID // 0x0 - Predef VAI Go To Pos (020xh)
	IncrementDemPos             SubID // 0x1 - Predef VAI Increment Dem Pos (021xh)
	IncrementTargetPos          SubID // 0x2 - Predef VAI Increment Target Pos (022xh)
	GoToPosFromActPosAndActVel  SubID // 0x3 - Predef VAI Go To Pos From Act Pos And Act Vel (023xh)
	GoToPosFromActPosDemVelZero SubID // 0x4 - Predef VAI Go To Pos From Act Pos Starting With Dem Vel = 0 (024xh)
	Stop                        SubID // 0x7 - Predef VAI Stop (027xh)
	GoToPosAfterActualCommand   SubID // 0x8 - Predef VAI Go To Pos After Actual Command (028xh)
	GoToPosOnRisingTrigger      SubID // 0xA - Predef VAI Go To Pos On Rising Trigger Event (02Axh)
	IncrementTargetPosOnRising  SubID // 0xB - Predef VAI Increment Target Pos On Rising Trigger Event (02Bxh)
	GoToPosOnFallingTrigger     SubID // 0xC - Predef VAI Go To Pos On Falling Trigger Event (02Cxh)
	IncrementTargetPosOnFalling SubID // 0xD - Predef VAI Increment Target Pos On Falling Trigger Event (02Dxh)
	InfiniteMotionPositive      SubID // 0xE - Predef VAI Infinite Motion Positive Direction (02Exh)
	InfiniteMotionNegative      SubID // 0xF - Predef VAI Infinite Motion Negative Direction (02Fxh)
}{
	GoToPos:                     0x0,
	IncrementDemPos:             0x1,
	IncrementTargetPos:          0x2,
	GoToPosFromActPosAndActVel:  0x3,
	GoToPosFromActPosDemVelZero: 0x4,
	Stop:                        0x7,
	GoToPosAfterActualCommand:   0x8,
	GoToPosOnRisingTrigger:      0xA,
	IncrementTargetPosOnRising:  0xB,
	GoToPosOnFallingTrigger:     0xC,
	IncrementTargetPosOnFalling: 0xD,
	InfiniteMotionPositive:      0xE,
	InfiniteMotionNegative:      0xF,
}
