package client_vai

import (
	"context"
	"fmt"

	client_common "github.com/Smart-Vision-Works/staged_robot/client/common"
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_motion_control "github.com/Smart-Vision-Works/staged_robot/protocol/motion_control"
	protocol_vai "github.com/Smart-Vision-Works/staged_robot/protocol/motion_control/vai"
)

// VAIManager handles VAI (Variable Interpolator) motion commands.
type VAIManager struct {
	requestManager *client_common.RequestManager
}

// NewVAIManager creates a new VAI manager.
func NewVAIManager(requestManager *client_common.RequestManager) *VAIManager {
	return &VAIManager{
		requestManager: requestManager,
	}
}

// GoToPosition sends a VAI Go To Position command.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.9
func (m *VAIManager) GoToPosition(ctx context.Context, positionMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	// Convert units
	position := protocol_common.ToProtocolPosition(positionMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIGoToPosCommand(position, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Go To Position command: %w", err)
	}

	return response.Status(), nil
}

// IncrementDemandPosition sends a VAI Increment Demand Position command.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.10
func (m *VAIManager) IncrementDemandPosition(ctx context.Context, incrementMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	// Convert units
	increment := protocol_common.ToProtocolPosition(incrementMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIIncrementDemPosCommand(increment, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Increment Demand Position command: %w", err)
	}

	return response.Status(), nil
}

// IncrementTargetPosition sends a VAI Increment Target Position command.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.11
func (m *VAIManager) IncrementTargetPosition(ctx context.Context, incrementMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	// Convert units
	increment := protocol_common.ToProtocolPosition(incrementMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIIncrementTargetPosCommand(increment, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Increment Target Position command: %w", err)
	}

	return response.Status(), nil
}

// Stop sends a VAI Stop command.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.16
func (m *VAIManager) Stop(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_vai.NewVAIStopCommand()

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Stop command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionFromActual sends a VAI Go To Position From Actual Position command.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.12
func (m *VAIManager) GoToPositionFromActual(ctx context.Context, positionMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	// Convert units
	position := protocol_common.ToProtocolPosition(positionMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIGoToPosFromActPosCommand(position, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Go To Position From Actual command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionFromActualDemVelZero sends a VAI Go To Position From Actual with start velocity = 0 command.
// Start velocity is forced to zero.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.13
func (m *VAIManager) GoToPositionFromActualDemVelZero(ctx context.Context, positionMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIGoToPosFromActPosDemVelZeroCommand(position, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Go To Position From Actual Dem Vel Zero command: %w", err)
	}

	return response.Status(), nil
}

// IncrementActualPosition sends a VAI Increment Actual Position command.
// New target position is calculated by adding increment to actual position.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.14
func (m *VAIManager) IncrementActualPosition(ctx context.Context, incrementMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	increment := protocol_common.ToProtocolPosition(incrementMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIIncrementActPosCommand(increment, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Increment Actual Position command: %w", err)
	}

	return response.Status(), nil
}

// IncrementActualPositionDemVelZero sends a VAI Increment Actual Position with start velocity = 0 command.
// Start velocity is forced to zero.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.15
func (m *VAIManager) IncrementActualPositionDemVelZero(ctx context.Context, incrementMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	increment := protocol_common.ToProtocolPosition(incrementMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIIncrementActPosDemVelZeroCommand(increment, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Increment Actual Position Dem Vel Zero command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionAfterActualCommand sends a VAI Go To Position After Actual Command.
// Waits until current motion finishes, then starts new VAI motion.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.17
func (m *VAIManager) GoToPositionAfterActualCommand(ctx context.Context, positionMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIGoToPosAfterActualCommandCommand(position, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Go To Position After Actual Command: %w", err)
	}

	return response.Status(), nil
}

// GoToAnalogPosition sends a VAI Go To Analog Position command.
// Target position comes from analog input; only motion parameters are specified.
//
// Parameters:
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.18
func (m *VAIManager) GoToAnalogPosition(ctx context.Context, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIGoToAnalogPosCommand(velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Go To Analog Position command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionOnRisingTrigger sends a VAI Go To Position On Rising Trigger Event command.
// Command waits for a rising trigger event before executing.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.19
func (m *VAIManager) GoToPositionOnRisingTrigger(ctx context.Context, positionMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIGoToPosOnRisingTriggerCommand(position, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Go To Position On Rising Trigger command: %w", err)
	}

	return response.Status(), nil
}

// IncrementTargetPositionOnRisingTrigger sends a VAI Increment Target Position On Rising Trigger Event command.
// Command waits for a rising trigger event before executing.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.20
func (m *VAIManager) IncrementTargetPositionOnRisingTrigger(ctx context.Context, incrementMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	increment := protocol_common.ToProtocolPosition(incrementMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIIncrementTargetPosOnRisingCommand(increment, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Increment Target Position On Rising Trigger command: %w", err)
	}

	return response.Status(), nil
}

// GoToPositionOnFallingTrigger sends a VAI Go To Position On Falling Trigger Event command.
// Command waits for a falling trigger event before executing.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.21
func (m *VAIManager) GoToPositionOnFallingTrigger(ctx context.Context, positionMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIGoToPosOnFallingTriggerCommand(position, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Go To Position On Falling Trigger command: %w", err)
	}

	return response.Status(), nil
}

// IncrementTargetPositionOnFallingTrigger sends a VAI Increment Target Position On Falling Trigger Event command.
// Command waits for a falling trigger event before executing.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.22
func (m *VAIManager) IncrementTargetPositionOnFallingTrigger(ctx context.Context, incrementMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	increment := protocol_common.ToProtocolPosition(incrementMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIIncrementTargetPosOnFallingCommand(increment, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Increment Target Position On Falling Trigger command: %w", err)
	}

	return response.Status(), nil
}

// ChangeMotionParamsOnPositiveTransition sends a VAI Change Motion Parameters On Positive Position Transition command.
// Changes motion parameters when demand position crosses the transition position in positive direction.
//
// Parameters:
//   - transitionPosMM: Position where parameters change in millimeters
//   - velocityMS: Maximal velocity after event in meters per second
//   - accelMS2: Acceleration after event in meters per second squared
//   - decelMS2: Deceleration after event in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.23
func (m *VAIManager) ChangeMotionParamsOnPositiveTransition(ctx context.Context, transitionPosMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	transitionPos := protocol_common.ToProtocolPosition(transitionPosMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIChangeParamsOnPosTransCommand(transitionPos, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Change Motion Parameters On Positive Transition command: %w", err)
	}

	return response.Status(), nil
}

// ChangeMotionParamsOnNegativeTransition sends a VAI Change Motion Parameters On Negative Position Transition command.
// Changes motion parameters when demand position crosses the transition position in negative direction.
//
// Parameters:
//   - transitionPosMM: Position where parameters change in millimeters
//   - velocityMS: Maximal velocity after event in meters per second
//   - accelMS2: Acceleration after event in meters per second squared
//   - decelMS2: Deceleration after event in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.24
func (m *VAIManager) ChangeMotionParamsOnNegativeTransition(ctx context.Context, transitionPosMM float64, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	transitionPos := protocol_common.ToProtocolPosition(transitionPosMM)
	velocity := protocol_common.ToProtocolVelocity(velocityMS)
	accel := protocol_common.ToProtocolAcceleration(accelMS2)
	decel := protocol_common.ToProtocolAcceleration(decelMS2)

	request := protocol_vai.NewVAIChangeParamsOnNegTransCommand(transitionPos, velocity, accel, decel)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send VAI Change Motion Parameters On Negative Transition command: %w", err)
	}

	return response.Status(), nil
}
