package protocol_predef_vai16

// ============================================================================
// Predefined 16-Bit VAI Sub ID Constants
// ============================================================================

// SubID represents a Predefined 16-bit VAI command identifier.
type SubID uint8

// SubIDs groups all Predefined 16-bit VAI command Sub ID constants.
// These commands combine predefined motion parameters with 16-bit position values.
// Reference: LinMot_MotionCtrl.txt, Section 4.3.88-4.3.98
var SubIDs = struct {
	GoToPos                     SubID // 0x0 - Predef VAI 16 Bit Go To Pos (0A0xh)
	IncrementDemPos             SubID // 0x1 - Predef VAI 16 Bit Increment Dem Pos (0A1xh)
	IncrementTargetPos          SubID // 0x2 - Predef VAI 16 Bit Increment Target Pos (0A2xh)
	GoToPosFromActPosAndActVel  SubID // 0x3 - Predef VAI 16 Bit Go To Pos From Act Pos And Act Vel (0A3xh)
	GoToPosFromActPosDemVelZero SubID // 0x4 - Predef VAI 16 Bit Go To Pos From Act Pos Starting With Dem Vel = 0 (0A4xh)
	Stop                        SubID // 0x7 - Predef VAI 16 Bit Stop (0A7xh)
	GoToPosAfterActualCommand   SubID // 0x8 - Predef VAI 16 Bit Go To Pos After Actual Command (0A8xh)
	GoToPosOnRisingTrigger      SubID // 0xA - Predef VAI 16 Bit Go To Pos On Rising Trigger Event (0AAxh)
	IncrementTargetPosOnRising  SubID // 0xB - Predef VAI 16 Bit Increment Target Pos On Rising Trigger Event (0ABxh)
	GoToPosOnFallingTrigger     SubID // 0xC - Predef VAI 16 Bit Go To Pos On Falling Trigger Event (0ACxh)
	IncrementTargetPosOnFalling SubID // 0xD - Predef VAI 16 Bit Increment Target Pos On Falling Trigger Event (0ADxh)
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
}
