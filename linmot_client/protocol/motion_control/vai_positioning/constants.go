package protocol_vai_positioning

// ============================================================================
// VAI Positioning Sub ID Constants
// ============================================================================

// SubID represents a VAI Positioning command identifier.
type SubID uint8

// SubIDs groups all VAI Positioning command Sub ID constants.
// These are advanced positioning variants including captured position, command table variables, and trigger config.
// Reference: LinMot_MotionCtrl.txt, Section 4.3.124-4.3.129
var SubIDs = struct {
	GoRelativeToCapturedPos       SubID // 0x0 - VAI Go Relative To Captured Pos (0D0xh)
	DecAcc16BitGoToPos            SubID // 0x1 - VAI Dec=Acc 16 Bit Go To Pos (0D1xh)
	GoToCmdTableVar1Pos           SubID // 0x4 - VAI Go To Cmd Table Var 1 Pos (0D4xh)
	GoToCmdTableVar2Pos           SubID // 0x5 - VAI Go To Cmd Table Var 2 Pos (0D5xh)
	GoToCmdTableVar1PosFromActPos SubID // 0x6 - VAI Go To Cmd Table Var 1 Pos From Act Pos And Act Vel (0D6xh)
	GoToCmdTableVar2PosFromActPos SubID // 0x7 - VAI Go To Cmd Table Var 2 Pos From Act Pos And Act Vel (0D7xh)
	StartTrigRiseConfigVAICommand SubID // 0xE - VAI Start Trig Rise Config VAI Command (0DExh)
	StartTrigFallConfigVAICommand SubID // 0xF - VAI Start Trig Fall Config VAI Command (0DFxh)
}{
	GoRelativeToCapturedPos:       0x0,
	DecAcc16BitGoToPos:            0x1,
	GoToCmdTableVar1Pos:           0x4,
	GoToCmdTableVar2Pos:           0x5,
	GoToCmdTableVar1PosFromActPos: 0x6,
	GoToCmdTableVar2PosFromActPos: 0x7,
	StartTrigRiseConfigVAICommand: 0xE,
	StartTrigFallConfigVAICommand: 0xF,
}
