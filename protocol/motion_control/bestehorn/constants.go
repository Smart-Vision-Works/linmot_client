package protocol_bestehorn

// ============================================================================
// Bestehorn VAJ Sub ID Constants
// ============================================================================

// SubID represents a Bestehorn VAJ command identifier.
type SubID uint8

// SubIDs groups all Bestehorn VAJ command Sub ID constants.
// Bestehorn VAJ commands provide jerk-limited motion for smoother trajectories.
// Reference: LinMot_MotionCtrl.txt, Section 4.3.140-4.3.149
var SubIDs = struct {
	GoToPos                         SubID // 0x0 - Bestehorn VAJ Go To Pos (0F0xh)
	IncrementDemPos                 SubID // 0x1 - Bestehorn VAJ Increment Demand Pos (0F1xh)
	GoToPosFromActualPos            SubID // 0x4 - Bestehorn VAJ Go To Pos From Actual Pos (0F4xh)
	IncrementActualPos              SubID // 0x6 - Bestehorn VAJ Increment Actual Pos (0F6xh)
	GoToPosAfterActualCommand       SubID // 0x8 - Bestehorn VAJ Go To Pos After Actual Command (0F8xh)
	GoToAnalogPos                   SubID // 0x9 - Bestehorn VAJ Go To Analog Pos (0F9xh)
	GoToPosOnRisingTrigger          SubID // 0xA - Bestehorn VAJ Go To Pos On Rising Trigger Event (0FAxh)
	IncrementDemPosOnRisingTrigger  SubID // 0xB - Bestehorn VAJ Increment Demand Pos On Rising Trigger Event (0FBxh)
	GoToPosOnFallingTrigger         SubID // 0xC - Bestehorn VAJ Go To Pos On Falling Trigger Event (0FCxh)
	IncrementDemPosOnFallingTrigger SubID // 0xD - Bestehorn VAJ Increment Demand Pos On Falling Trigger Event (0FDxh)
}{
	GoToPos:                         0x0,
	IncrementDemPos:                 0x1,
	GoToPosFromActualPos:            0x4,
	IncrementActualPos:              0x6,
	GoToPosAfterActualCommand:       0x8,
	GoToAnalogPos:                   0x9,
	GoToPosOnRisingTrigger:          0xA,
	IncrementDemPosOnRisingTrigger:  0xB,
	GoToPosOnFallingTrigger:         0xC,
	IncrementDemPosOnFallingTrigger: 0xD,
}
