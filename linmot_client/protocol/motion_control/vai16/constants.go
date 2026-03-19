package protocol_vai16

// ============================================================================
// 16-Bit VAI Sub ID Constants
// ============================================================================

// SubID represents a 16-bit VAI command identifier.
type SubID uint8

// SubIDs groups all 16-bit VAI command Sub ID constants.
// These commands use 16-bit position values (±3.28mm range) for reduced precision but faster execution.
// Reference: LinMot_MotionCtrl.txt, Section 4.3.73-4.3.87
var SubIDs = struct {
	GoToPos                     SubID // 0x0 - VAI 16 Bit Go To Pos (090xh)
	IncrementDemPos             SubID // 0x1 - VAI 16 Bit Increment Dem Pos (091xh)
	IncrementTargetPos          SubID // 0x2 - VAI 16 Bit Increment Target Pos (092xh)
	GoToPosFromActPosAndActVel  SubID // 0x3 - VAI 16 Bit Go To Pos From Act Pos And Act Vel (093xh)
	GoToPosFromActPosDemVelZero SubID // 0x4 - VAI 16 Bit Go To Pos From Act Pos Starting With Dem Vel = 0 (094xh)
	IncrementActPos             SubID // 0x5 - VAI 16 Bit Increment Act Pos (095xh)
	IncrementActPosDemVelZero   SubID // 0x6 - VAI 16 Bit Increment Act Pos Starting With Dem Vel = 0 (096xh)
	Stop                        SubID // 0x7 - VAI 16 Bit Stop (097xh)
	GoToPosAfterActualCommand   SubID // 0x8 - VAI 16 Bit Go To Pos After Actual Command (098xh)
	GoToPosOnRisingTrigger      SubID // 0xA - VAI 16 Bit Go To Pos On Rising Trigger Event (09Axh)
	IncrementTargetPosOnRising  SubID // 0xB - VAI 16 Bit Increment Target Pos On Rising Trigger Event (09Bxh)
	GoToPosOnFallingTrigger     SubID // 0xC - VAI 16 Bit Go To Pos On Falling Trigger Event (09Cxh)
	IncrementTargetPosOnFalling SubID // 0xD - VAI 16 Bit Increment Target Pos On Falling Trigger Event (09Dxh)
	ChangeParamsOnPosTrans      SubID // 0xE - VAI 16 Bit Change Motion Parameters On Positive Position Transition (09Exh)
	ChangeParamsOnNegTrans      SubID // 0xF - VAI 16 Bit Change Motion Parameters On Negative Position Transition (09Fxh)
}{
	GoToPos:                     0x0,
	IncrementDemPos:             0x1,
	IncrementTargetPos:          0x2,
	GoToPosFromActPosAndActVel:  0x3,
	GoToPosFromActPosDemVelZero: 0x4,
	IncrementActPos:             0x5,
	IncrementActPosDemVelZero:   0x6,
	Stop:                        0x7,
	GoToPosAfterActualCommand:   0x8,
	GoToPosOnRisingTrigger:      0xA,
	IncrementTargetPosOnRising:  0xB,
	GoToPosOnFallingTrigger:     0xC,
	IncrementTargetPosOnFalling: 0xD,
	ChangeParamsOnPosTrans:      0xE,
	ChangeParamsOnNegTrans:      0xF,
}
