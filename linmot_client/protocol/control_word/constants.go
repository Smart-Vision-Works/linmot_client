package protocol_control_word

// Control Word Bit Positions
// From LinMot_MotionCtrl.txt section 3.22
const (
	ControlWordBit_SwitchOn         = 0
	ControlWordBit_EnableVoltage    = 1
	ControlWordBit_QuickStop        = 2 // Inverted logic: 1 = normal, 0 = quick stop
	ControlWordBit_EnableOperation  = 3
	ControlWordBit_Abort            = 4 // Inverted logic: 1 = normal, 0 = abort
	ControlWordBit_Freeze           = 5 // Inverted logic: 1 = normal, 0 = freeze
	ControlWordBit_Home             = 6
	ControlWordBit_ErrorAcknowledge = 7
	ControlWordBit_JogPlus          = 8
	ControlWordBit_JogMinus         = 9
	ControlWordBit_SpecialMode      = 10
	ControlWordBit_ClearanceCheck   = 12
	ControlWordBit_GoToInitialPos   = 13
	ControlWordBit_Reserved         = 14
	ControlWordBit_PhaseSearch      = 15
)

// Status Word Bit Positions
// From LinMot_MotionCtrl.txt section 3.23
const (
	StatusWordBit_OperationEnabled = 0
	StatusWordBit_SwitchOnActive   = 1
	StatusWordBit_EnableOperation  = 2
	StatusWordBit_Error            = 3
	StatusWordBit_VoltageEnable    = 4
	StatusWordBit_QuickStop        = 5
	StatusWordBit_SwitchOnLocked   = 6
	StatusWordBit_Warning          = 7
	StatusWordBit_EventHandler     = 8
	StatusWordBit_SpecialMotion    = 9
	StatusWordBit_InTargetPosition = 10
	StatusWordBit_Homed            = 11
	StatusWordBit_FatalError       = 12
	StatusWordBit_MotionActive     = 13
	StatusWordBit_RangeIndicator1  = 14
	StatusWordBit_RangeIndicator2  = 15
)

// State Machine States
// From LinMot_MotionCtrl.txt state diagram
type MainState uint8

const (
	State_NotReadyToSwitchOn MainState = 0
	State_SwitchOnDisabled   MainState = 1
	State_ReadyToSwitchOn    MainState = 2
	State_SetupError         MainState = 3
	State_Error              MainState = 4
	State_SwitchOn           MainState = 6
	State_OperationEnabled   MainState = 8
	State_Homing             MainState = 9
	State_ClearanceCheck     MainState = 10
	State_GoingToInitialPos  MainState = 11
	State_Aborting           MainState = 12
	State_Freezing           MainState = 13
	State_QuickStopActive    MainState = 14
	State_GoingToPosition    MainState = 15
	State_JoggingPlus        MainState = 16
	State_JoggingMinus       MainState = 17
	State_Linearizing        MainState = 18
	State_PhaseSearching     MainState = 19
	State_SpecialMode        MainState = 20
)

// String returns human-readable state name
func (s MainState) String() string {
	switch s {
	case State_NotReadyToSwitchOn:
		return "Not Ready To Switch On"
	case State_SwitchOnDisabled:
		return "Switch On Disabled"
	case State_ReadyToSwitchOn:
		return "Ready To Switch On"
	case State_SetupError:
		return "Setup Error"
	case State_Error:
		return "Error"
	case State_SwitchOn:
		return "Switch On"
	case State_OperationEnabled:
		return "Operation Enabled"
	case State_Homing:
		return "Homing"
	case State_ClearanceCheck:
		return "Clearance Check"
	case State_GoingToInitialPos:
		return "Going To Initial Position"
	case State_Aborting:
		return "Aborting"
	case State_Freezing:
		return "Freezing"
	case State_QuickStopActive:
		return "Quick Stop Active"
	case State_GoingToPosition:
		return "Going To Position"
	case State_JoggingPlus:
		return "Jogging Plus"
	case State_JoggingMinus:
		return "Jogging Minus"
	case State_Linearizing:
		return "Linearizing"
	case State_PhaseSearching:
		return "Phase Searching"
	case State_SpecialMode:
		return "Special Mode"
	default:
		return "Unknown State"
	}
}

// Common Control Word Patterns
// These are commonly used bit combinations for state transitions

// EnableDrivePattern returns control word to enable drive
// Sets: Switch On, Enable Voltage, Quick Stop (released), Enable Operation
func EnableDrivePattern() uint16 {
	return (1 << ControlWordBit_SwitchOn) |
		(1 << ControlWordBit_EnableVoltage) |
		(1 << ControlWordBit_QuickStop) |
		(1 << ControlWordBit_EnableOperation)
}

// DisableDrivePattern returns control word to disable drive
// Clears: Switch On (transitions to Switch On Disabled state)
func DisableDrivePattern() uint16 {
	return 0x0000
}

// ErrorAcknowledgePattern returns control word with error acknowledge bit set
// Used for rising edge of error acknowledge
func ErrorAcknowledgePattern() uint16 {
	return (1 << ControlWordBit_ErrorAcknowledge)
}

// QuickStopPattern returns control word to trigger quick stop
// Clears the Quick Stop bit (inverted logic)
func QuickStopPattern() uint16 {
	return 0x0000 // All bits clear, including Quick Stop
}

// SetBit sets a specific bit in the control word
func SetBit(word uint16, bit uint) uint16 {
	return word | (1 << bit)
}

// ClearBit clears a specific bit in the control word
func ClearBit(word uint16, bit uint) uint16 {
	return word &^ (1 << bit)
}

// IsBitSet checks if a specific bit is set in a word
func IsBitSet(word uint16, bit uint) bool {
	return (word & (1 << bit)) != 0
}
