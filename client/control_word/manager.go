package client_control_word

import (
	"context"
	"fmt"
	"time"

	client_common "gsail-go/linmot/client/common"
	protocol_common "gsail-go/linmot/protocol/common"
	protocol_control_word "gsail-go/linmot/protocol/control_word"
)

type ControlWordManager struct {
	requestManager *client_common.RequestManager
}

func NewControlWordManager(requestManager *client_common.RequestManager) *ControlWordManager {
	return &ControlWordManager{
		requestManager: requestManager,
	}
}

// EnableDrive transitions the drive to Operation Enabled state
// Sets: Switch On, Enable Voltage, Quick Stop (released), Enable Operation
func (m *ControlWordManager) EnableDrive(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_control_word.NewControlWordRequest(protocol_control_word.EnableDrivePattern())

	response, err := client_common.SendRequestAndReceive[*protocol_common.StatusResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send enable drive command: %w", err)
	}

	status := response.Status()
	helper := protocol_control_word.NewStatusWordHelper(status.StatusWord, status.StateVar)

	// Check for fatal error
	if helper.IsFatalError() {
		return status, fmt.Errorf("cannot enable drive: fatal error present (error code: 0x%04X)", status.ErrorCode)
	}

	// Check for error
	if helper.HasError() {
		return status, fmt.Errorf("cannot enable drive: error present (error code: 0x%04X)", status.ErrorCode)
	}

	return status, nil
}

// DisableDrive transitions the drive to Switch On Disabled state
func (m *ControlWordManager) DisableDrive(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_control_word.NewControlWordRequest(protocol_control_word.DisableDrivePattern())

	response, err := client_common.SendRequestAndReceive[*protocol_common.StatusResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send disable drive command: %w", err)
	}

	return response.Status(), nil
}

// acknowledgeErrorPollInterval is how often we poll drive status after sending
// an acknowledge edge, waiting for the error to clear.
const acknowledgeErrorPollInterval = 5 * time.Millisecond

// acknowledgeErrorTimeout is the maximum time to wait for the drive to clear
// an error after sending the acknowledge edges.
const acknowledgeErrorTimeout = 1 * time.Second

// AcknowledgeError acknowledges and clears a drive error.
// Sends rising edge on bit 7, polls until error clears, then sends a control word
// with EnableDrivePattern() (bit 7 cleared) to avoid leaving acknowledge latched.
// If no error is active, this is a no-op and returns current status.
// Fails if a fatal error is present (cannot be acknowledged) or if the error
// does not clear within the timeout.
func (m *ControlWordManager) AcknowledgeError(ctx context.Context) (*protocol_common.Status, error) {
	// First, check current status
	currentStatus, err := m.GetStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current status: %w", err)
	}

	helper := protocol_control_word.NewStatusWordHelper(currentStatus.StatusWord, currentStatus.StateVar)
	if helper.IsFatalError() {
		return currentStatus, fmt.Errorf("cannot acknowledge fatal error (error code: 0x%04X)", currentStatus.ErrorCode)
	}
	// No active error: do not toggle acknowledge bit at all.
	if !helper.HasError() {
		return currentStatus, nil
	}

	// Send rising edge (bit 7 set)
	request := protocol_control_word.NewControlWordRequest(protocol_control_word.ErrorAcknowledgePattern())
	_, err = client_common.SendRequestAndReceive[*protocol_common.StatusResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send error acknowledge (rising edge): %w", err)
	}

	// Poll until error clears or timeout
	pollCtx, pollCancel := context.WithTimeout(ctx, acknowledgeErrorTimeout)
	defer pollCancel()

	clearedStatus, err := m.pollUntilErrorCleared(pollCtx, currentStatus)
	if err != nil {
		return clearedStatus, fmt.Errorf("error did not clear after rising edge: %w", err)
	}

	// Send falling edge (bit 7 clear) while restoring enable bits.
	request = protocol_control_word.NewControlWordRequest(protocol_control_word.EnableDrivePattern())
	_, err = client_common.SendRequestAndReceive[*protocol_common.StatusResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send error acknowledge (falling edge): %w", err)
	}

	// Poll again to confirm error stays cleared
	pollCtx2, pollCancel2 := context.WithTimeout(ctx, acknowledgeErrorTimeout)
	defer pollCancel2()

	finalStatus, err := m.pollUntilErrorCleared(pollCtx2, clearedStatus)
	if err != nil {
		return finalStatus, fmt.Errorf("error reappeared after falling edge: %w", err)
	}

	return finalStatus, nil
}

// pollUntilErrorCleared polls drive status at acknowledgeErrorPollInterval until
// HasError() returns false. Returns the cleared status or the last-known status
// alongside an error so callers retain diagnostics when polling fails or times out.
func (m *ControlWordManager) pollUntilErrorCleared(ctx context.Context, initialStatus *protocol_common.Status) (*protocol_common.Status, error) {
	ticker := time.NewTicker(acknowledgeErrorPollInterval)
	defer ticker.Stop()

	lastStatus := initialStatus

	for {
		status, err := m.GetStatus(ctx)
		if err != nil {
			return lastStatus, fmt.Errorf("failed to poll status: %w", err)
		}
		lastStatus = status

		helper := protocol_control_word.NewStatusWordHelper(status.StatusWord, status.StateVar)
		if !helper.HasError() {
			return status, nil
		}

		select {
		case <-ctx.Done():
			lastErrorCode := uint16(0)
			if lastStatus != nil {
				lastErrorCode = lastStatus.ErrorCode
			}
			return lastStatus, fmt.Errorf("timeout waiting for error to clear (last error code: 0x%04X): %w", lastErrorCode, ctx.Err())
		case <-ticker.C:
		}
	}
}

// QuickStop triggers emergency quick stop
// Clears the Quick Stop bit (inverted logic)
func (m *ControlWordManager) QuickStop(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_control_word.NewControlWordRequest(protocol_control_word.QuickStopPattern())

	response, err := client_common.SendRequestAndReceive[*protocol_common.StatusResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send quick stop command: %w", err)
	}

	return response.Status(), nil
}

// Home initiates homing sequence
// Sets bit 6 (Home) while maintaining operation enabled state
func (m *ControlWordManager) Home(ctx context.Context) (*protocol_common.Status, error) {
	// Create context with longer timeout for homing
	homeCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Maintain operation enabled bits and add home bit
	word := protocol_control_word.NewControlWordBuilder().
		SwitchOn().
		EnableVoltage().
		QuickStop(). // Note: QuickStop sets the bit (released = bit set)
		EnableOperation().
		Home().
		Build()

	request := protocol_control_word.NewControlWordRequest(word)

	response, err := client_common.SendRequestAndReceive[*protocol_common.StatusResponse](m.requestManager, homeCtx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send home command: %w", err)
	}

	return response.Status(), nil
}

// ClearanceCheck initiates clearance check
// Sets bit 12 (Clearance Check)
func (m *ControlWordManager) ClearanceCheck(ctx context.Context) (*protocol_common.Status, error) {
	word := protocol_control_word.NewControlWordBuilder().
		ClearanceCheck().
		Build()

	request := protocol_control_word.NewControlWordRequest(word)

	response, err := client_common.SendRequestAndReceive[*protocol_common.StatusResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send clearance check command: %w", err)
	}

	return response.Status(), nil
}

// GoToInitialPosition moves drive to initial position
// Sets bit 13 (Go To Initial Position)
func (m *ControlWordManager) GoToInitialPosition(ctx context.Context) (*protocol_common.Status, error) {
	word := protocol_control_word.NewControlWordBuilder().
		GoToInitialPosition().
		Build()

	request := protocol_control_word.NewControlWordRequest(word)

	response, err := client_common.SendRequestAndReceive[*protocol_common.StatusResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send go to initial position command: %w", err)
	}

	return response.Status(), nil
}

// JogPlus starts positive jog motion
// Sets bit 8 (Jog Move +)
// Motion continues until StopJog is called
func (m *ControlWordManager) JogPlus(ctx context.Context) (*protocol_common.Status, error) {
	word := protocol_control_word.NewControlWordBuilder().
		JogPlus().
		Build()

	request := protocol_control_word.NewControlWordRequest(word)

	response, err := client_common.SendRequestAndReceive[*protocol_common.StatusResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send jog plus command: %w", err)
	}

	return response.Status(), nil
}

// JogMinus starts negative jog motion
// Sets bit 9 (Jog Move -)
// Motion continues until StopJog is called
func (m *ControlWordManager) JogMinus(ctx context.Context) (*protocol_common.Status, error) {
	word := protocol_control_word.NewControlWordBuilder().
		JogMinus().
		Build()

	request := protocol_control_word.NewControlWordRequest(word)

	response, err := client_common.SendRequestAndReceive[*protocol_common.StatusResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send jog minus command: %w", err)
	}

	return response.Status(), nil
}

// StopJog stops any active jog motion
// Clears bits 8 and 9 (Jog Move + and -)
func (m *ControlWordManager) StopJog(ctx context.Context) (*protocol_common.Status, error) {
	word := protocol_control_word.NewControlWordBuilder().
		StopJog().
		Build()

	request := protocol_control_word.NewControlWordRequest(word)

	response, err := client_common.SendRequestAndReceive[*protocol_common.StatusResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send stop jog command: %w", err)
	}

	return response.Status(), nil
}

// SendControlWord sends a raw control word value
// For advanced use cases requiring direct control word manipulation
func (m *ControlWordManager) SendControlWord(ctx context.Context, word uint16) (*protocol_common.Status, error) {
	request := protocol_control_word.NewControlWordRequest(word)

	response, err := client_common.SendRequestAndReceive[*protocol_common.StatusResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send control word: %w", err)
	}

	return response.Status(), nil
}

// GetStatus queries current drive status without changing control word
// Sends a status request (no control word module)
func (m *ControlWordManager) GetStatus(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_common.NewStatusRequest()

	response, err := client_common.SendRequestAndReceive[*protocol_common.StatusResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	return response.Status(), nil
}

// WaitForState polls status until the main state matches the target state
// Returns when state is reached or context times out
func (m *ControlWordManager) WaitForState(ctx context.Context, targetState protocol_control_word.MainState, pollInterval time.Duration) (*protocol_common.Status, error) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		status, err := m.GetStatus(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get status while waiting for state: %w", err)
		}
		helper := protocol_control_word.NewStatusWordHelper(status.StatusWord, status.StateVar)
		if helper.GetMainState() == targetState {
			return status, nil
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for state %s: %w", targetState.String(), ctx.Err())
		case <-ticker.C:
		}
	}
}

// WaitForOperationEnabled polls status until Operation Enabled bit is set
// Used internally by EnableDrive and can be used standalone
func (m *ControlWordManager) WaitForOperationEnabled(ctx context.Context, pollInterval time.Duration) (*protocol_common.Status, error) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		status, err := m.GetStatus(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get status while waiting for operation enabled: %w", err)
		}

		helper := protocol_control_word.NewStatusWordHelper(status.StatusWord, status.StateVar)
		if helper.IsOperationEnabled() {
			return status, nil
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for operation enabled: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}
