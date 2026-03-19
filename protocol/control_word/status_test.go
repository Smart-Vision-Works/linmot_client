package protocol_control_word

import (
	"testing"
)

func TestStatusWordHelper_BitChecks(t *testing.T) {
	tests := []struct {
		name       string
		statusWord uint16
		checkFunc  func(*StatusWordHelper) bool
		expected   bool
	}{
		{"OperationEnabled_Set", 0x0001, (*StatusWordHelper).IsOperationEnabled, true},
		{"OperationEnabled_Clear", 0x0000, (*StatusWordHelper).IsOperationEnabled, false},
		{"SwitchOnActive_Set", 0x0002, (*StatusWordHelper).IsSwitchOnActive, true},
		{"SwitchOnActive_Clear", 0x0000, (*StatusWordHelper).IsSwitchOnActive, false},
		{"EnableOperation_Set", 0x0004, (*StatusWordHelper).IsEnableOperation, true},
		{"EnableOperation_Clear", 0x0000, (*StatusWordHelper).IsEnableOperation, false},
		{"Error_Set", 0x0008, (*StatusWordHelper).HasError, true},
		{"Error_Clear", 0x0000, (*StatusWordHelper).HasError, false},
		{"VoltageEnabled_Set", 0x0010, (*StatusWordHelper).IsVoltageEnabled, true},
		{"VoltageEnabled_Clear", 0x0000, (*StatusWordHelper).IsVoltageEnabled, false},
		{"QuickStopActive_Set", 0x0020, (*StatusWordHelper).IsQuickStopActive, true},
		{"QuickStopActive_Clear", 0x0000, (*StatusWordHelper).IsQuickStopActive, false},
		{"SwitchOnLocked_Set", 0x0040, (*StatusWordHelper).IsSwitchOnLocked, true},
		{"SwitchOnLocked_Clear", 0x0000, (*StatusWordHelper).IsSwitchOnLocked, false},
		{"Warning_Set", 0x0080, (*StatusWordHelper).HasWarning, true},
		{"Warning_Clear", 0x0000, (*StatusWordHelper).HasWarning, false},
		{"EventHandler_Set", 0x0100, (*StatusWordHelper).IsEventHandlerActive, true},
		{"EventHandler_Clear", 0x0000, (*StatusWordHelper).IsEventHandlerActive, false},
		{"SpecialMotion_Set", 0x0200, (*StatusWordHelper).IsSpecialMotionActive, true},
		{"SpecialMotion_Clear", 0x0000, (*StatusWordHelper).IsSpecialMotionActive, false},
		{"InTargetPosition_Set", 0x0400, (*StatusWordHelper).IsInTargetPosition, true},
		{"InTargetPosition_Clear", 0x0000, (*StatusWordHelper).IsInTargetPosition, false},
		{"Homed_Set", 0x0800, (*StatusWordHelper).IsHomed, true},
		{"Homed_Clear", 0x0000, (*StatusWordHelper).IsHomed, false},
		{"FatalError_Set", 0x1000, (*StatusWordHelper).IsFatalError, true},
		{"FatalError_Clear", 0x0000, (*StatusWordHelper).IsFatalError, false},
		{"MotionActive_Set", 0x2000, (*StatusWordHelper).IsMotionActive, true},
		{"MotionActive_Clear", 0x0000, (*StatusWordHelper).IsMotionActive, false},
		{"RangeIndicator1_Set", 0x4000, (*StatusWordHelper).GetRangeIndicator1, true},
		{"RangeIndicator1_Clear", 0x0000, (*StatusWordHelper).GetRangeIndicator1, false},
		{"RangeIndicator2_Set", 0x8000, (*StatusWordHelper).GetRangeIndicator2, true},
		{"RangeIndicator2_Clear", 0x0000, (*StatusWordHelper).GetRangeIndicator2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := NewStatusWordHelper(tt.statusWord, 0x0000)
			result := tt.checkFunc(helper)
			if result != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, result)
			}
		})
	}
}

func TestStatusWordHelper_StateExtraction(t *testing.T) {
	tests := []struct {
		name         string
		stateVar     uint16
		expectedMain MainState
		expectedSub  uint8
	}{
		{"NotReadyToSwitchOn", 0x0000, State_NotReadyToSwitchOn, 0x00},
		{"SwitchOnDisabled", 0x0100, State_SwitchOnDisabled, 0x00},
		{"SwitchOnDisabled_SubState5", 0x0105, State_SwitchOnDisabled, 0x05},
		{"ReadyToSwitchOn", 0x0200, State_ReadyToSwitchOn, 0x00},
		{"SetupError", 0x0300, State_SetupError, 0x00},
		{"Error", 0x0400, State_Error, 0x00},
		{"SwitchOn", 0x0600, State_SwitchOn, 0x00},
		{"OperationEnabled", 0x0800, State_OperationEnabled, 0x00},
		{"OperationEnabled_SubState10", 0x080A, State_OperationEnabled, 0x0A},
		{"Homing", 0x0900, State_Homing, 0x00},
		{"ClearanceCheck", 0x0A00, State_ClearanceCheck, 0x00},
		{"GoingToInitialPos", 0x0B00, State_GoingToInitialPos, 0x00},
		{"Aborting", 0x0C00, State_Aborting, 0x00},
		{"Freezing", 0x0D00, State_Freezing, 0x00},
		{"QuickStopActive", 0x0E00, State_QuickStopActive, 0x00},
		{"GoingToPosition", 0x0F00, State_GoingToPosition, 0x00},
		{"JoggingPlus", 0x1000, State_JoggingPlus, 0x00},
		{"JoggingMinus", 0x1100, State_JoggingMinus, 0x00},
		{"Linearizing", 0x1200, State_Linearizing, 0x00},
		{"PhaseSearching", 0x1300, State_PhaseSearching, 0x00},
		{"SpecialMode", 0x1400, State_SpecialMode, 0x00},
		{"MaxState_MaxSubState", 0xFFFF, MainState(0xFF), 0xFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := NewStatusWordHelper(0x0000, tt.stateVar)

			mainState := helper.GetMainState()
			if mainState != tt.expectedMain {
				t.Errorf("GetMainState(): expected %d, got %d", tt.expectedMain, mainState)
			}

			subState := helper.GetSubState()
			if subState != tt.expectedSub {
				t.Errorf("GetSubState(): expected %d, got %d", tt.expectedSub, subState)
			}
		})
	}
}

func TestStatusWordHelper_StateName(t *testing.T) {
	tests := []struct {
		stateVar     uint16
		expectedName string
	}{
		{0x0000, "Not Ready To Switch On"},
		{0x0100, "Switch On Disabled"},
		{0x0200, "Ready To Switch On"},
		{0x0300, "Setup Error"},
		{0x0400, "Error"},
		{0x0600, "Switch On"},
		{0x0800, "Operation Enabled"},
		{0x0900, "Homing"},
		{0x0A00, "Clearance Check"},
		{0x0B00, "Going To Initial Position"},
		{0x0C00, "Aborting"},
		{0x0D00, "Freezing"},
		{0x0E00, "Quick Stop Active"},
		{0x0F00, "Going To Position"},
		{0x1000, "Jogging Plus"},
		{0x1100, "Jogging Minus"},
		{0x1200, "Linearizing"},
		{0x1300, "Phase Searching"},
		{0x1400, "Special Mode"},
		{0xFF00, "Unknown State"},
	}

	for _, tt := range tests {
		t.Run(tt.expectedName, func(t *testing.T) {
			helper := NewStatusWordHelper(0x0000, tt.stateVar)
			name := helper.GetStateName()
			if name != tt.expectedName {
				t.Errorf("GetStateName(): expected %q, got %q", tt.expectedName, name)
			}
		})
	}
}

func TestStatusWordHelper_Getters(t *testing.T) {
	statusWord := uint16(0x1234)
	stateVar := uint16(0x5678)

	helper := NewStatusWordHelper(statusWord, stateVar)

	if helper.GetStatusWord() != statusWord {
		t.Errorf("GetStatusWord(): expected 0x%04X, got 0x%04X", statusWord, helper.GetStatusWord())
	}

	if helper.GetStateVar() != stateVar {
		t.Errorf("GetStateVar(): expected 0x%04X, got 0x%04X", stateVar, helper.GetStateVar())
	}
}

func TestStatusWordHelper_CombinedBits(t *testing.T) {
	// Test realistic status word with multiple bits set
	// Operation Enabled (bit 0) + Error (bit 3) + Warning (bit 7) = 0x0089
	statusWord := uint16(0x0089)
	helper := NewStatusWordHelper(statusWord, 0x0400) // State 4: Error

	if !helper.IsOperationEnabled() {
		t.Error("Operation Enabled should be set")
	}
	if !helper.HasError() {
		t.Error("Error should be set")
	}
	if !helper.HasWarning() {
		t.Error("Warning should be set")
	}
	if helper.IsFatalError() {
		t.Error("Fatal Error should not be set")
	}

	if helper.GetMainState() != State_Error {
		t.Errorf("Main state should be Error (4), got %d", helper.GetMainState())
	}
}

func TestMainState_String(t *testing.T) {
	tests := []struct {
		state    MainState
		expected string
	}{
		{State_NotReadyToSwitchOn, "Not Ready To Switch On"},
		{State_SwitchOnDisabled, "Switch On Disabled"},
		{State_ReadyToSwitchOn, "Ready To Switch On"},
		{State_SetupError, "Setup Error"},
		{State_Error, "Error"},
		{State_SwitchOn, "Switch On"},
		{State_OperationEnabled, "Operation Enabled"},
		{State_Homing, "Homing"},
		{State_ClearanceCheck, "Clearance Check"},
		{State_GoingToInitialPos, "Going To Initial Position"},
		{State_Aborting, "Aborting"},
		{State_Freezing, "Freezing"},
		{State_QuickStopActive, "Quick Stop Active"},
		{State_GoingToPosition, "Going To Position"},
		{State_JoggingPlus, "Jogging Plus"},
		{State_JoggingMinus, "Jogging Minus"},
		{State_Linearizing, "Linearizing"},
		{State_PhaseSearching, "Phase Searching"},
		{State_SpecialMode, "Special Mode"},
		{MainState(255), "Unknown State"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.state.String()
			if result != tt.expected {
				t.Errorf("String(): expected %q, got %q", tt.expected, result)
			}
		})
	}
}
