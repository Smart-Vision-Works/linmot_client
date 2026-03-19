package protocol_control_word

// StatusWordHelper provides helper methods for interpreting StatusWord and StateVar
type StatusWordHelper struct {
	statusWord uint16
	stateVar   uint16
}

// NewStatusWordHelper creates a new helper for interpreting status information
func NewStatusWordHelper(statusWord, stateVar uint16) *StatusWordHelper {
	return &StatusWordHelper{
		statusWord: statusWord,
		stateVar:   stateVar,
	}
}

// IsOperationEnabled checks if the drive is in Operation Enabled state (bit 0)
func (h *StatusWordHelper) IsOperationEnabled() bool {
	return IsBitSet(h.statusWord, StatusWordBit_OperationEnabled)
}

// IsSwitchOnActive checks if Switch On is active (bit 1)
func (h *StatusWordHelper) IsSwitchOnActive() bool {
	return IsBitSet(h.statusWord, StatusWordBit_SwitchOnActive)
}

// IsEnableOperation checks if Enable Operation is set (bit 2)
func (h *StatusWordHelper) IsEnableOperation() bool {
	return IsBitSet(h.statusWord, StatusWordBit_EnableOperation)
}

// HasError checks if an error is present (bit 3)
func (h *StatusWordHelper) HasError() bool {
	return IsBitSet(h.statusWord, StatusWordBit_Error)
}

// IsVoltageEnabled checks if voltage is enabled (bit 4)
func (h *StatusWordHelper) IsVoltageEnabled() bool {
	return IsBitSet(h.statusWord, StatusWordBit_VoltageEnable)
}

// IsQuickStopActive checks if quick stop is active (bit 5)
// Note: This bit has inverted logic in control word but normal logic in status word
func (h *StatusWordHelper) IsQuickStopActive() bool {
	return IsBitSet(h.statusWord, StatusWordBit_QuickStop)
}

// IsSwitchOnLocked checks if switch on is locked (bit 6)
func (h *StatusWordHelper) IsSwitchOnLocked() bool {
	return IsBitSet(h.statusWord, StatusWordBit_SwitchOnLocked)
}

// HasWarning checks if a warning is present (bit 7)
func (h *StatusWordHelper) HasWarning() bool {
	return IsBitSet(h.statusWord, StatusWordBit_Warning)
}

// IsEventHandlerActive checks if the event handler is active (bit 8)
func (h *StatusWordHelper) IsEventHandlerActive() bool {
	return IsBitSet(h.statusWord, StatusWordBit_EventHandler)
}

// IsSpecialMotionActive checks if special motion is active (bit 9)
func (h *StatusWordHelper) IsSpecialMotionActive() bool {
	return IsBitSet(h.statusWord, StatusWordBit_SpecialMotion)
}

// IsInTargetPosition checks if the drive is in target position (bit 10)
func (h *StatusWordHelper) IsInTargetPosition() bool {
	return IsBitSet(h.statusWord, StatusWordBit_InTargetPosition)
}

// IsHomed checks if the drive has been homed (bit 11)
func (h *StatusWordHelper) IsHomed() bool {
	return IsBitSet(h.statusWord, StatusWordBit_Homed)
}

// IsFatalError checks if a fatal error is present (bit 12)
func (h *StatusWordHelper) IsFatalError() bool {
	return IsBitSet(h.statusWord, StatusWordBit_FatalError)
}

// IsMotionActive checks if motion is active (bit 13)
func (h *StatusWordHelper) IsMotionActive() bool {
	return IsBitSet(h.statusWord, StatusWordBit_MotionActive)
}

// GetRangeIndicator1 checks Range Indicator 1 (bit 14)
func (h *StatusWordHelper) GetRangeIndicator1() bool {
	return IsBitSet(h.statusWord, StatusWordBit_RangeIndicator1)
}

// GetRangeIndicator2 checks Range Indicator 2 (bit 15)
func (h *StatusWordHelper) GetRangeIndicator2() bool {
	return IsBitSet(h.statusWord, StatusWordBit_RangeIndicator2)
}

// GetMainState extracts the main state from StateVar (high byte)
func (h *StatusWordHelper) GetMainState() MainState {
	return MainState((h.stateVar >> 8) & 0xFF)
}

// GetSubState extracts the sub-state from StateVar (low byte)
func (h *StatusWordHelper) GetSubState() uint8 {
	return uint8(h.stateVar & 0xFF)
}

// GetStateName returns a human-readable name for the current main state
func (h *StatusWordHelper) GetStateName() string {
	return h.GetMainState().String()
}

// GetStatusWord returns the raw status word value
func (h *StatusWordHelper) GetStatusWord() uint16 {
	return h.statusWord
}

// GetStateVar returns the raw state var value
func (h *StatusWordHelper) GetStateVar() uint16 {
	return h.stateVar
}
