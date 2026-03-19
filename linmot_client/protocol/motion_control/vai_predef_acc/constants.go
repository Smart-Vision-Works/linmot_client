package protocol_vai_predef_acc

// ============================================================================
// VAI Predefined Acceleration Sub ID Constants
// ============================================================================

// SubID represents a VAI Predefined Acceleration command identifier.
type SubID uint8

// SubIDs groups all VAI Predefined Acceleration command Sub ID constants.
// These commands use drive-configured acceleration, but allow variable velocity and deceleration.
// Reference: LinMot_MotionCtrl.txt, Section 4.3.99-4.3.108
var SubIDs = struct {
	GoToPos                     SubID // 0x0 - VAI Predef Acc Go To Pos (0B0xh)
	IncrementDemPos             SubID // 0x1 - VAI Predef Acc Increment Dem Pos (0B1xh)
	IncrementTargetPos          SubID // 0x2 - VAI Predef Acc Increment Target Pos (0B2xh)
	GoToPosFromActPosAndActVel  SubID // 0x3 - VAI Predef Acc Go To Pos From Act Pos And Act Vel (0B3xh)
	GoToPosFromActPosDemVelZero SubID // 0x4 - VAI Predef Acc Go To Pos From Act Pos Starting With Dem Vel = 0 (0B4xh)
	GoToPosAfterActualCommand   SubID // 0x8 - VAI Predef Acc Go To Pos After Actual Command (0B8xh)
	GoToPosOnRisingTrigger      SubID // 0xA - VAI Predef Acc Go To Pos On Rising Trigger Event (0BAxh)
	IncrementTargetPosOnRising  SubID // 0xB - VAI Predef Acc Increment Target Pos On Rising Trigger Event (0BBxh)
	GoToPosOnFallingTrigger     SubID // 0xC - VAI Predef Acc Go To Pos On Falling Trigger Event (0BCxh)
	IncrementTargetPosOnFalling SubID // 0xD - VAI Predef Acc Increment Target Pos On Falling Trigger Event (0BDxh)
}{
	GoToPos:                     0x0,
	IncrementDemPos:             0x1,
	IncrementTargetPos:          0x2,
	GoToPosFromActPosAndActVel:  0x3,
	GoToPosFromActPosDemVelZero: 0x4,
	GoToPosAfterActualCommand:   0x8,
	GoToPosOnRisingTrigger:      0xA,
	IncrementTargetPosOnRising:  0xB,
	GoToPosOnFallingTrigger:     0xC,
	IncrementTargetPosOnFalling: 0xD,
}
