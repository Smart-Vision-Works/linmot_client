package protocol_interface_control

// ============================================================================
// Interface Control Sub ID Constants
// ============================================================================

// SubID represents an Interface Control command identifier.
type SubID uint8

// SubIDs groups all Interface Control command Sub ID constants.
// Reference: LinMot_MotionCtrl.txt, Section 4.3.1-4.3.8
var SubIDs = struct {
	NoOperation                 SubID // 0x0 - No Operation (000xh)
	WriteInterfaceControlWord   SubID // 0x1 - Write Interface Control Word (001xh)
	WriteLiveParameter          SubID // 0x2 - Write Live Parameter (002xh)
	WriteX4IntfOutputsWithMask  SubID // 0x3 - Write X4/X14 Intf Outputs with Mask (003xh)
	SelectPositionControllerSet SubID // 0x5 - Select Position Controller Set (005xh)
	ClearEventEvaluation        SubID // 0x8 - Clear Event Evaluation (008xh)
	MasterHoming                SubID // 0x9 - Master Homing (009xh)
	Reset                       SubID // 0xF - Reset (00Fxh)
}{
	NoOperation:                 0x0,
	WriteInterfaceControlWord:   0x1,
	WriteLiveParameter:          0x2,
	WriteX4IntfOutputsWithMask:  0x3,
	SelectPositionControllerSet: 0x5,
	ClearEventEvaluation:        0x8,
	MasterHoming:                0x9,
	Reset:                       0xF,
}
