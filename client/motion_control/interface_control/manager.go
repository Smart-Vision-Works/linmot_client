package client_interface_control

import (
	"context"
	"fmt"

	client_common "gsail-go/linmot/client/common"
	protocol_common "gsail-go/linmot/protocol/common"
	protocol_motion_control "gsail-go/linmot/protocol/motion_control"
	protocol_interface_control "gsail-go/linmot/protocol/motion_control/interface_control"
)

// InterfaceControlManager handles Interface Control commands.
type InterfaceControlManager struct {
	requestManager *client_common.RequestManager
}

// NewInterfaceControlManager creates a new Interface Control manager.
func NewInterfaceControlManager(requestManager *client_common.RequestManager) *InterfaceControlManager {
	return &InterfaceControlManager{
		requestManager: requestManager,
	}
}

// NoOperation sends a No Operation command.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.1
func (m *InterfaceControlManager) NoOperation(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_interface_control.NewNoOperationCommand()

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send No Operation command: %w", err)
	}

	return response.Status(), nil
}

// WriteInterfaceControlWord sends a Write Interface Control Word command.
//
// Parameters:
//   - controlWord: Interface control word value
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.2
func (m *InterfaceControlManager) WriteInterfaceControlWord(ctx context.Context, controlWord uint16) (*protocol_common.Status, error) {
	request := protocol_interface_control.NewWriteInterfaceControlWordCommand(controlWord)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Write Interface Control Word command: %w", err)
	}

	return response.Status(), nil
}

// WriteLiveParameter sends a Write Live Parameter command.
//
// Parameters:
//   - upid: Unique Parameter ID
//   - value: Parameter value
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.3
func (m *InterfaceControlManager) WriteLiveParameter(ctx context.Context, upid uint16, value uint32) (*protocol_common.Status, error) {
	request := protocol_interface_control.NewWriteLiveParameterCommand(upid, value)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Write Live Parameter command: %w", err)
	}

	return response.Status(), nil
}

// WriteOutputsWithMask sends a Write X4/X14 Interface Outputs with Mask command.
//
// Parameters:
//   - bitMask: Bit mask for which outputs to write
//   - bitValue: Bit values for the outputs
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.4
func (m *InterfaceControlManager) WriteOutputsWithMask(ctx context.Context, bitMask, bitValue uint16) (*protocol_common.Status, error) {
	request := protocol_interface_control.NewWriteX4IntfOutputsWithMaskCommand(bitMask, bitValue)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Write Outputs With Mask command: %w", err)
	}

	return response.Status(), nil
}

// SelectPositionControllerSet sends a Select Position Controller Set command.
//
// Parameters:
//   - controllerSet: Controller set selection (0 = Set A, 1 = Set B)
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.5
func (m *InterfaceControlManager) SelectPositionControllerSet(ctx context.Context, controllerSet uint16) (*protocol_common.Status, error) {
	request := protocol_interface_control.NewSelectPositionControllerSetCommand(controllerSet)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Select Position Controller Set command: %w", err)
	}

	return response.Status(), nil
}

// ClearEventEvaluation sends a Clear Event Evaluation command.
// Resets the event handler used for trigger-based commands.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.6
func (m *InterfaceControlManager) ClearEventEvaluation(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_interface_control.NewClearEventEvaluationCommand()

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Clear Event Evaluation command: %w", err)
	}

	return response.Status(), nil
}

// MasterHoming sends a Master Homing command.
//
// Parameters:
//   - homePositionMM: Home position in millimeters
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.7
func (m *InterfaceControlManager) MasterHoming(ctx context.Context, homePositionMM float64) (*protocol_common.Status, error) {
	homePosition := protocol_common.ToProtocolPosition(homePositionMM)

	request := protocol_interface_control.NewMasterHomingCommand(homePosition)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Master Homing command: %w", err)
	}

	return response.Status(), nil
}

// Reset sends a Reset command.
// Resets all firmware instances of the drive.
// WARNING: This will reboot the drive!
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.8
func (m *InterfaceControlManager) Reset(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_interface_control.NewResetCommand()

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Reset command: %w", err)
	}

	return response.Status(), nil
}
