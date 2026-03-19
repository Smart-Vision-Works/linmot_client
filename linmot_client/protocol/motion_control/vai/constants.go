package protocol_vai

// ============================================================================
// VAI Sub ID Constants
// ============================================================================

// SubID represents a VAI command identifier.
type SubID uint8

// SubIDs groups all VAI command Sub ID constants.
// Reference: LinMot_MotionCtrl.txt, Section 4.3
var SubIDs = struct {
	GoToPos                     SubID // 0x0 - VAI Go To Pos (010xh)
	IncrementDemPos             SubID // 0x1 - VAI Increment Dem Pos (011xh)
	IncrementTargetPos          SubID // 0x2 - VAI Increment Target Pos (012xh)
	GoToPosFromActPosAndActVel  SubID // 0x3 - VAI Go To Pos From Act Pos And Act Vel (013xh)
	GoToPosFromActPosDemVelZero SubID // 0x4 - VAI Go To Pos From Act Pos Starting With Dem Vel = 0 (014xh)
	IncrementActPos             SubID // 0x5 - VAI Increment Act Pos (015xh)
	IncrementActPosDemVelZero   SubID // 0x6 - VAI Increment Act Pos Starting With Dem Vel = 0 (016xh)
	Stop                        SubID // 0x7 - VAI Stop (017xh)
	GoToPosAfterActualCommand   SubID // 0x8 - VAI Go To Pos After Actual Command (018xh)
	GoToAnalogPos               SubID // 0x9 - VAI Go To Analog Pos (019xh)
	GoToPosOnRisingTrigger      SubID // 0xA - VAI Go To Pos On Rising Trigger Event (01Axh)
	IncrementTargetPosOnRising  SubID // 0xB - VAI Increment Target Pos On Rising Trigger Event (01Bxh)
	GoToPosOnFallingTrigger     SubID // 0xC - VAI Go To Pos On Falling Trigger Event (01Cxh)
	IncrementTargetPosOnFalling SubID // 0xD - VAI Increment Target Pos On Falling Trigger Event (01Dxh)
	ChangeParamsOnPosTrans      SubID // 0xE - VAI Change Motion Parameters On Positive Position Transition (01Exh)
	ChangeParamsOnNegTrans      SubID // 0xF - VAI Change Motion Parameters On Negative Position Transition (01Fxh)
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
	GoToAnalogPos:               0x9,
	GoToPosOnRisingTrigger:      0xA,
	IncrementTargetPosOnRising:  0xB,
	GoToPosOnFallingTrigger:     0xC,
	IncrementTargetPosOnFalling: 0xD,
	ChangeParamsOnPosTrans:      0xE,
	ChangeParamsOnNegTrans:      0xF,
}
