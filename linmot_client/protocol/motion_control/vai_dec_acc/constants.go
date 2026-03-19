package protocol_vai_dec_acc

// ============================================================================
// VAI Dec=Acc Sub ID Constants
// ============================================================================

// SubID represents a VAI Dec=Acc command identifier.
type SubID uint8

// SubIDs groups all VAI Dec=Acc command Sub ID constants.
// These commands use symmetric motion (deceleration = acceleration).
// Reference: LinMot_MotionCtrl.txt, Section 4.3.109-4.3.123
var SubIDs = struct {
	GoToPos                            SubID // 0x0 - VAI Dec=Acc Go To Pos (0C0xh)
	IncrementDemPos                    SubID // 0x1 - VAI Dec=Acc Increment Dem Pos (0C1xh)
	IncrementTargetPos                 SubID // 0x2 - VAI Dec=Acc Increment Target Pos (0C2xh)
	GoToPosFromActPosAndActVel         SubID // 0x3 - VAI Dec=Acc Go To Pos From Act Pos And Act Vel (0C3xh)
	GoToPosFromActPosDemVelZero        SubID // 0x4 - VAI Dec=Acc Go To Pos From Act Pos Starting With Dem Vel = 0 (0C4xh)
	GoToPosWithMaxCurr                 SubID // 0x5 - VAI Dec=Acc Go To Pos With Max Curr (0C5xh)
	GoToPosFromActPosMaxCurr           SubID // 0x6 - VAI Dec=Acc Go To Pos From Act Pos And Act Vel With Max Curr (0C6xh)
	GoToPosFromActPosDemVelZeroMaxCurr SubID // 0x7 - VAI Dec=Acc Go To Pos From Act Pos, Dem Vel = 0 and With Max Curr (0C7xh)
	GoToPosAfterActualCommand          SubID // 0x8 - VAI Dec=Acc Go To Pos After Actual Command (0C8xh)
	GoToPosOnRisingTrigger             SubID // 0xA - VAI Dec=Acc Go To Pos On Rising Trigger Event (0CAxh)
	IncrementTargetPosOnRising         SubID // 0xB - VAI Dec=Acc Increment Target Pos On Rising Trigger Event (0CBxh)
	GoToPosOnFallingTrigger            SubID // 0xC - VAI Dec=Acc Go To Pos On Falling Trigger Event (0CCxh)
	IncrementTargetPosOnFalling        SubID // 0xD - VAI Dec=Acc Increment Target Pos On Falling Trigger Event (0CDxh)
	InfiniteMotionPositive             SubID // 0xE - VAI Dec=Acc Infinite Motion Positive Direction (0CExh)
	InfiniteMotionNegative             SubID // 0xF - VAI Dec=Acc Infinite Motion Negative Direction (0CFxh)
}{
	GoToPos:                            0x0,
	IncrementDemPos:                    0x1,
	IncrementTargetPos:                 0x2,
	GoToPosFromActPosAndActVel:         0x3,
	GoToPosFromActPosDemVelZero:        0x4,
	GoToPosWithMaxCurr:                 0x5,
	GoToPosFromActPosMaxCurr:           0x6,
	GoToPosFromActPosDemVelZeroMaxCurr: 0x7,
	GoToPosAfterActualCommand:          0x8,
	GoToPosOnRisingTrigger:             0xA,
	IncrementTargetPosOnRising:         0xB,
	GoToPosOnFallingTrigger:            0xC,
	IncrementTargetPosOnFalling:        0xD,
	InfiniteMotionPositive:             0xE,
	InfiniteMotionNegative:             0xF,
}
