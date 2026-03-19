package client_predefined

import (
	"context"
	"fmt"

	client_common "github.com/Smart-Vision-Works/staged_robot/client/common"
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_motion_control "github.com/Smart-Vision-Works/staged_robot/protocol/motion_control"
	protocol_predefined "github.com/Smart-Vision-Works/staged_robot/protocol/motion_control/predefined"
)

// PredefVAIManager handles Predefined VAI (Variable Interpolator) motion commands.
// These commands use drive-configured motion parameters (velocity, accel, decel).
type PredefVAIManager struct {
	requestManager *client_common.RequestManager
}

// NewPredefVAIManager creates a new Predefined VAI manager.
func NewPredefVAIManager(requestManager *client_common.RequestManager) *PredefVAIManager {
	return &PredefVAIManager{
		requestManager: requestManager,
	}
}

// GoToPosition sends a Predefined VAI Go To Position command.
// Uses drive-configured velocity, acceleration, and deceleration.
//
// Parameters:
//   - positionMM: Target position in millimeters
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.25
func (m *PredefVAIManager) GoToPosition(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)

	request := protocol_predefined.NewPredefVAIGoToPosCommand(position)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Predef VAI Go To Position command: %w", err)
	}

	return response.Status(), nil
}

// IncrementDemandPosition sends a Predefined VAI Increment Demand Position command.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.26
func (m *PredefVAIManager) IncrementDemandPosition(ctx context.Context, incrementMM float64) (*protocol_common.Status, error) {
	increment := protocol_common.ToProtocolPosition(incrementMM)

	request := protocol_predefined.NewPredefVAIIncrementDemPosCommand(increment)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Predef VAI Increment Demand Position command: %w", err)
	}

	return response.Status(), nil
}

// IncrementTargetPosition sends a Predefined VAI Increment Target Position command.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.27
func (m *PredefVAIManager) IncrementTargetPosition(ctx context.Context, incrementMM float64) (*protocol_common.Status, error) {
	increment := protocol_common.ToProtocolPosition(incrementMM)

	request := protocol_predefined.NewPredefVAIIncrementTargetPosCommand(increment)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Predef VAI Increment Target Position command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionFromActual sends a Predefined VAI Go To Pos From Actual Position command.
//
// Parameters:
//   - positionMM: Target position in millimeters
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.28
func (m *PredefVAIManager) GoToPositionFromActual(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)

	request := protocol_predefined.NewPredefVAIGoToPosFromActPosCommand(position)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Predef VAI Go To Position From Actual command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionFromActualDemVelZero sends a Predefined VAI Go To Pos Starting With Dem Vel = 0 command.
//
// Parameters:
//   - positionMM: Target position in millimeters
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.29
func (m *PredefVAIManager) GoToPositionFromActualDemVelZero(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)

	request := protocol_predefined.NewPredefVAIGoToPosFromActPosDemVelZeroCommand(position)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Predef VAI Go To Position From Actual Dem Vel Zero command: %w", err)
	}

	return response.Status(), nil
}

// Stop sends a Predefined VAI Stop command.
// Uses drive-configured quick stop deceleration.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.30
func (m *PredefVAIManager) Stop(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_predefined.NewPredefVAIStopCommand()

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Predef VAI Stop command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionAfterActualCommand sends a Predefined VAI Go To Pos After Actual Command.
//
// Parameters:
//   - positionMM: Target position in millimeters
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.31
func (m *PredefVAIManager) GoToPositionAfterActualCommand(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)

	request := protocol_predefined.NewPredefVAIGoToPosAfterActualCommandCommand(position)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Predef VAI Go To Position After Actual Command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionOnRisingTrigger sends a Predefined VAI Go To Pos On Rising Trigger Event command.
//
// Parameters:
//   - positionMM: Target position in millimeters
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.32
func (m *PredefVAIManager) GoToPositionOnRisingTrigger(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)

	request := protocol_predefined.NewPredefVAIGoToPosOnRisingTriggerCommand(position)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Predef VAI Go To Position On Rising Trigger command: %w", err)
	}

	return response.Status(), nil
}

// IncrementTargetPositionOnRisingTrigger sends a Predefined VAI Increment Target Pos On Rising Trigger Event command.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.33
func (m *PredefVAIManager) IncrementTargetPositionOnRisingTrigger(ctx context.Context, incrementMM float64) (*protocol_common.Status, error) {
	increment := protocol_common.ToProtocolPosition(incrementMM)

	request := protocol_predefined.NewPredefVAIIncrementTargetPosOnRisingCommand(increment)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Predef VAI Increment Target Position On Rising Trigger command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionOnFallingTrigger sends a Predefined VAI Go To Pos On Falling Trigger Event command.
//
// Parameters:
//   - positionMM: Target position in millimeters
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.34
func (m *PredefVAIManager) GoToPositionOnFallingTrigger(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)

	request := protocol_predefined.NewPredefVAIGoToPosOnFallingTriggerCommand(position)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Predef VAI Go To Position On Falling Trigger command: %w", err)
	}

	return response.Status(), nil
}

// IncrementTargetPositionOnFallingTrigger sends a Predefined VAI Increment Target Pos On Falling Trigger Event command.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.35
func (m *PredefVAIManager) IncrementTargetPositionOnFallingTrigger(ctx context.Context, incrementMM float64) (*protocol_common.Status, error) {
	increment := protocol_common.ToProtocolPosition(incrementMM)

	request := protocol_predefined.NewPredefVAIIncrementTargetPosOnFallingCommand(increment)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Predef VAI Increment Target Position On Falling Trigger command: %w", err)
	}

	return response.Status(), nil
}

// InfiniteMotionPositive sends a Predefined VAI Infinite Motion Positive Direction command.
// Starts continuous motion in positive direction.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.36
func (m *PredefVAIManager) InfiniteMotionPositive(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_predefined.NewPredefVAIInfiniteMotionPositiveCommand()

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Predef VAI Infinite Motion Positive command: %w", err)
	}

	return response.Status(), nil
}

// InfiniteMotionNegative sends a Predefined VAI Infinite Motion Negative Direction command.
// Starts continuous motion in negative direction.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.37
func (m *PredefVAIManager) InfiniteMotionNegative(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_predefined.NewPredefVAIInfiniteMotionNegativeCommand()

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Predef VAI Infinite Motion Negative command: %w", err)
	}

	return response.Status(), nil
}
