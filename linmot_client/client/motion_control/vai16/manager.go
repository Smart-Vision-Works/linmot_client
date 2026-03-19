package client_vai16

import (
	"context"
	"fmt"

	client_common "github.com/Smart-Vision-Works/staged_robot/client/common"
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_motion_control "github.com/Smart-Vision-Works/staged_robot/protocol/motion_control"
	protocol_vai16 "github.com/Smart-Vision-Works/staged_robot/protocol/motion_control/vai16"
)

// VAI16Manager handles 16-bit VAI (Variable Interpolator) motion commands.
// Position range limited to ±32767 units (±3.28mm).
type VAI16Manager struct {
	requestManager *client_common.RequestManager
}

const (
	maxVAI16Units = 32767
	minVAI16Units = -32768
)

func toVAI16Units(mm float64) (int16, error) {
	raw := protocol_common.ToProtocolPosition(mm)
	if raw < minVAI16Units || raw > maxVAI16Units {
		return 0, fmt.Errorf("%.4fmm exceeds the VAI16 range (±3.2767mm)", mm)
	}
	return int16(raw), nil
}

// NewVAI16Manager creates a new 16-bit VAI manager.
func NewVAI16Manager(requestManager *client_common.RequestManager) *VAI16Manager {
	return &VAI16Manager{
		requestManager: requestManager,
	}
}

// GoToPosition sends a 16-bit VAI Go To Position command.
// Position range: ±3.28mm
func (m *VAI16Manager) GoToPosition(ctx context.Context, positionMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	position, err := toVAI16Units(positionMM)
	if err != nil {
		return nil, fmt.Errorf("GoToPosition: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16GoToPosCommand(position, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Go To Position command: %w", err)
	}

	return response.Status(), nil
}

// IncrementDemandPosition sends a 16-bit VAI Increment Demand Position command.
func (m *VAI16Manager) IncrementDemandPosition(ctx context.Context, incrementMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	increment, err := toVAI16Units(incrementMM)
	if err != nil {
		return nil, fmt.Errorf("IncrementDemandPosition: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16IncrementDemPosCommand(increment, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Increment Demand Position command: %w", err)
	}

	return response.Status(), nil
}

// IncrementTargetPosition sends a 16-bit VAI Increment Target Position command.
func (m *VAI16Manager) IncrementTargetPosition(ctx context.Context, incrementMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	increment, err := toVAI16Units(incrementMM)
	if err != nil {
		return nil, fmt.Errorf("IncrementTargetPosition: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16IncrementTargetPosCommand(increment, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Increment Target Position command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionFromActual sends a 16-bit VAI Go To Position From Actual Position command.
func (m *VAI16Manager) GoToPositionFromActual(ctx context.Context, positionMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	position, err := toVAI16Units(positionMM)
	if err != nil {
		return nil, fmt.Errorf("GoToPositionFromActual: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16GoToPosFromActPosCommand(position, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Go To Position From Actual command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionFromActualDemVelZero sends a 16-bit VAI Go To Position From Actual with start velocity = 0 command.
func (m *VAI16Manager) GoToPositionFromActualDemVelZero(ctx context.Context, positionMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	position, err := toVAI16Units(positionMM)
	if err != nil {
		return nil, fmt.Errorf("GoToPositionFromActualDemVelZero: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16GoToPosFromActPosDemVelZeroCommand(position, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Go To Position From Actual Dem Vel Zero command: %w", err)
	}

	return response.Status(), nil
}

// IncrementActualPosition sends a 16-bit VAI Increment Actual Position command.
func (m *VAI16Manager) IncrementActualPosition(ctx context.Context, incrementMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	increment, err := toVAI16Units(incrementMM)
	if err != nil {
		return nil, fmt.Errorf("IncrementActualPosition: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16IncrementActPosCommand(increment, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Increment Actual Position command: %w", err)
	}

	return response.Status(), nil
}

// IncrementActualPositionDemVelZero sends a 16-bit VAI Increment Actual Position with start velocity = 0 command.
func (m *VAI16Manager) IncrementActualPositionDemVelZero(ctx context.Context, incrementMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	increment, err := toVAI16Units(incrementMM)
	if err != nil {
		return nil, fmt.Errorf("IncrementActualPositionDemVelZero: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16IncrementActPosDemVelZeroCommand(increment, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Increment Actual Position Dem Vel Zero command: %w", err)
	}

	return response.Status(), nil
}

// Stop sends a 16-bit VAI Stop command.
func (m *VAI16Manager) Stop(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_vai16.NewVAI16StopCommand()

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Stop command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionAfterActualCommand sends a 16-bit VAI Go To Position After Actual Command.
func (m *VAI16Manager) GoToPositionAfterActualCommand(ctx context.Context, positionMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	position, err := toVAI16Units(positionMM)
	if err != nil {
		return nil, fmt.Errorf("GoToPositionAfterActualCommand: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16GoToPosAfterActualCommandCommand(position, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Go To Position After Actual Command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionOnRisingTrigger sends a 16-bit VAI Go To Position On Rising Trigger Event command.
func (m *VAI16Manager) GoToPositionOnRisingTrigger(ctx context.Context, positionMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	position, err := toVAI16Units(positionMM)
	if err != nil {
		return nil, fmt.Errorf("GoToPositionOnRisingTrigger: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16GoToPosOnRisingTriggerCommand(position, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Go To Position On Rising Trigger command: %w", err)
	}

	return response.Status(), nil
}

// IncrementTargetPositionOnRisingTrigger sends a 16-bit VAI Increment Target Position On Rising Trigger Event command.
func (m *VAI16Manager) IncrementTargetPositionOnRisingTrigger(ctx context.Context, incrementMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	increment, err := toVAI16Units(incrementMM)
	if err != nil {
		return nil, fmt.Errorf("IncrementTargetPositionOnRisingTrigger: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16IncrementTargetPosOnRisingCommand(increment, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Increment Target Position On Rising Trigger command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionOnFallingTrigger sends a 16-bit VAI Go To Position On Falling Trigger Event command.
func (m *VAI16Manager) GoToPositionOnFallingTrigger(ctx context.Context, positionMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	position, err := toVAI16Units(positionMM)
	if err != nil {
		return nil, fmt.Errorf("GoToPositionOnFallingTrigger: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16GoToPosOnFallingTriggerCommand(position, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Go To Position On Falling Trigger command: %w", err)
	}

	return response.Status(), nil
}

// IncrementTargetPositionOnFallingTrigger sends a 16-bit VAI Increment Target Position On Falling Trigger Event command.
func (m *VAI16Manager) IncrementTargetPositionOnFallingTrigger(ctx context.Context, incrementMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	increment, err := toVAI16Units(incrementMM)
	if err != nil {
		return nil, fmt.Errorf("IncrementTargetPositionOnFallingTrigger: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16IncrementTargetPosOnFallingCommand(increment, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Increment Target Position On Falling Trigger command: %w", err)
	}

	return response.Status(), nil
}

// ChangeMotionParamsOnPositiveTransition sends a 16-bit VAI Change Motion Parameters On Positive Position Transition command.
func (m *VAI16Manager) ChangeMotionParamsOnPositiveTransition(ctx context.Context, transitionPosMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	transitionPos, err := toVAI16Units(transitionPosMM)
	if err != nil {
		return nil, fmt.Errorf("ChangeMotionParamsOnPositiveTransition: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16ChangeParamsOnPosTransCommand(transitionPos, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Change Motion Parameters On Positive Transition command: %w", err)
	}

	return response.Status(), nil
}

// ChangeMotionParamsOnNegativeTransition sends a 16-bit VAI Change Motion Parameters On Negative Position Transition command.
func (m *VAI16Manager) ChangeMotionParamsOnNegativeTransition(ctx context.Context, transitionPosMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	transitionPos, err := toVAI16Units(transitionPosMM)
	if err != nil {
		return nil, fmt.Errorf("ChangeMotionParamsOnNegativeTransition: %w", err)
	}
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai16.NewVAI16ChangeParamsOnNegTransCommand(transitionPos, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI16 Change Motion Parameters On Negative Transition command: %w", err)
	}

	return response.Status(), nil
}
