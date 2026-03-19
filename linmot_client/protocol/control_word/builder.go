package protocol_control_word

// ControlWordBuilder provides a fluent API for building control word values
type ControlWordBuilder struct {
	word uint16
}

// NewControlWordBuilder creates a new builder with all bits cleared
func NewControlWordBuilder() *ControlWordBuilder {
	return &ControlWordBuilder{word: 0}
}

// NewControlWordBuilderFrom creates a new builder initialized with an existing word
func NewControlWordBuilderFrom(word uint16) *ControlWordBuilder {
	return &ControlWordBuilder{word: word}
}

// SwitchOn sets the Switch On bit (bit 0)
func (b *ControlWordBuilder) SwitchOn() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_SwitchOn)
	return b
}

// EnableVoltage sets the Enable Voltage bit (bit 1)
func (b *ControlWordBuilder) EnableVoltage() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_EnableVoltage)
	return b
}

// ReleaseQuickStop sets the Quick Stop bit (bit 2)
// Note: This bit has inverted logic - 1 = normal operation, 0 = quick stop
func (b *ControlWordBuilder) ReleaseQuickStop() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_QuickStop)
	return b
}

// QuickStop clears the Quick Stop bit (bit 2) to trigger quick stop
// Note: This bit has inverted logic - 1 = normal operation, 0 = quick stop
func (b *ControlWordBuilder) QuickStop() *ControlWordBuilder {
	b.word = ClearBit(b.word, ControlWordBit_QuickStop)
	return b
}

// EnableOperation sets the Enable Operation bit (bit 3)
func (b *ControlWordBuilder) EnableOperation() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_EnableOperation)
	return b
}

// ReleaseAbort sets the Abort bit (bit 4)
// Note: This bit has inverted logic - 1 = normal operation, 0 = abort
func (b *ControlWordBuilder) ReleaseAbort() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_Abort)
	return b
}

// Abort clears the Abort bit (bit 4) to trigger abort
// Note: This bit has inverted logic - 1 = normal operation, 0 = abort
func (b *ControlWordBuilder) Abort() *ControlWordBuilder {
	b.word = ClearBit(b.word, ControlWordBit_Abort)
	return b
}

// ReleaseFreeze sets the Freeze bit (bit 5)
// Note: This bit has inverted logic - 1 = normal operation, 0 = freeze
func (b *ControlWordBuilder) ReleaseFreeze() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_Freeze)
	return b
}

// Freeze clears the Freeze bit (bit 5) to trigger freeze
// Note: This bit has inverted logic - 1 = normal operation, 0 = freeze
func (b *ControlWordBuilder) Freeze() *ControlWordBuilder {
	b.word = ClearBit(b.word, ControlWordBit_Freeze)
	return b
}

// Home sets the Home bit (bit 6)
func (b *ControlWordBuilder) Home() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_Home)
	return b
}

// AcknowledgeError sets the Error Acknowledge bit (bit 7)
func (b *ControlWordBuilder) AcknowledgeError() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_ErrorAcknowledge)
	return b
}

// ClearErrorAcknowledge clears the Error Acknowledge bit (bit 7)
// Used for the falling edge after error acknowledge
func (b *ControlWordBuilder) ClearErrorAcknowledge() *ControlWordBuilder {
	b.word = ClearBit(b.word, ControlWordBit_ErrorAcknowledge)
	return b
}

// JogPlus sets the Jog Move + bit (bit 8)
func (b *ControlWordBuilder) JogPlus() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_JogPlus)
	return b
}

// JogMinus sets the Jog Move - bit (bit 9)
func (b *ControlWordBuilder) JogMinus() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_JogMinus)
	return b
}

// StopJog clears both jog bits (bits 8 and 9)
func (b *ControlWordBuilder) StopJog() *ControlWordBuilder {
	b.word = ClearBit(b.word, ControlWordBit_JogPlus)
	b.word = ClearBit(b.word, ControlWordBit_JogMinus)
	return b
}

// SpecialMode sets the Special Mode bit (bit 10)
func (b *ControlWordBuilder) SpecialMode() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_SpecialMode)
	return b
}

// ClearanceCheck sets the Clearance Check bit (bit 12)
func (b *ControlWordBuilder) ClearanceCheck() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_ClearanceCheck)
	return b
}

// GoToInitialPosition sets the Go To Initial Position bit (bit 13)
func (b *ControlWordBuilder) GoToInitialPosition() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_GoToInitialPos)
	return b
}

// PhaseSearch sets the Phase Search bit (bit 15)
func (b *ControlWordBuilder) PhaseSearch() *ControlWordBuilder {
	b.word = SetBit(b.word, ControlWordBit_PhaseSearch)
	return b
}

// WithBit sets an arbitrary bit (fluent interface)
func (b *ControlWordBuilder) WithBit(bit uint) *ControlWordBuilder {
	b.word = SetBit(b.word, bit)
	return b
}

// SetBit sets an arbitrary bit (fluent interface).
// Deprecated: use WithBit.
func (b *ControlWordBuilder) SetBit(bit uint) *ControlWordBuilder {
	return b.WithBit(bit)
}

// WithoutBit clears an arbitrary bit (fluent interface)
func (b *ControlWordBuilder) WithoutBit(bit uint) *ControlWordBuilder {
	b.word = ClearBit(b.word, bit)
	return b
}

// ClearBit clears an arbitrary bit (fluent interface).
// Deprecated: use WithoutBit.
func (b *ControlWordBuilder) ClearBit(bit uint) *ControlWordBuilder {
	return b.WithoutBit(bit)
}

// WithPattern replaces the current word with a pattern
func (b *ControlWordBuilder) WithPattern(pattern uint16) *ControlWordBuilder {
	b.word = pattern
	return b
}

// Build returns the final control word value
func (b *ControlWordBuilder) Build() uint16 {
	return b.word
}

// BuildRequest creates a ControlWordRequest from the builder
func (b *ControlWordBuilder) BuildRequest() *ControlWordRequest {
	return NewControlWordRequest(b.word)
}
