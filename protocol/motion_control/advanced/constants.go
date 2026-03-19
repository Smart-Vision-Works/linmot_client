package protocol_advanced

// ============================================================================
// Advanced Features Sub ID Constants
// ============================================================================

// SubID represents an Advanced Features command identifier.
type SubID uint8

// SubIDs groups all Advanced Features command Sub ID constants.
// These include Sin VA (Sinusoidal Variable Acceleration) commands.
// Reference: LinMot_MotionCtrl.txt, Section 4.3.130-4.3.139
var SubIDs = struct {
	SinVAGoToPos                   SubID // 0x0 - Sin VA Go To Pos (0E0xh)
	SinVAIncrementDemandPos        SubID // 0x1 - Sin VA Increment Demand Pos (0E1xh)
	SinVAGoToPosFromActualPos      SubID // 0x4 - Sin VA Go To Pos From Actual Pos (0E4xh)
	SinVAIncrementActualPos        SubID // 0x6 - Sin VA Increment Actual Pos (0E6xh)
	SinVAGoToPosAfterActualCommand SubID // 0x8 - Sin VA Go To Pos After Actual Command (0E8xh)
	SinVAGoToAnalogPos             SubID // 0x9 - Sin VA Go To Analog Pos (0E9xh)
	SinVAGoToPosOnRisingTrigger    SubID // 0xA - Sin VA Go To Pos On Rising Trigger Event (0EAxh)
	SinVAIncrementDemPosOnRising   SubID // 0xB - Sin VA Increment Demand Pos On Rising Trigger Event (0EBxh)
	SinVAGoToPosOnFallingTrigger   SubID // 0xC - Sin VA Go To Pos On Falling Trigger Event (0ECxh)
	SinVAIncrementDemPosOnFalling  SubID // 0xD - Sin VA Increment Demand Pos On Falling Trigger Event (0EDxh)
}{
	SinVAGoToPos:                   0x0,
	SinVAIncrementDemandPos:        0x1,
	SinVAGoToPosFromActualPos:      0x4,
	SinVAIncrementActualPos:        0x6,
	SinVAGoToPosAfterActualCommand: 0x8,
	SinVAGoToAnalogPos:             0x9,
	SinVAGoToPosOnRisingTrigger:    0xA,
	SinVAIncrementDemPosOnRising:   0xB,
	SinVAGoToPosOnFallingTrigger:   0xC,
	SinVAIncrementDemPosOnFalling:  0xD,
}
