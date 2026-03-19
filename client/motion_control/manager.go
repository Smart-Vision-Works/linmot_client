package client_motion_control

import (
	"context"

	client_common "github.com/Smart-Vision-Works/linmot_client/client/common"
	client_interface_control "github.com/Smart-Vision-Works/linmot_client/client/motion_control/interface_control"
	client_predefined "github.com/Smart-Vision-Works/linmot_client/client/motion_control/predefined"
	client_streaming "github.com/Smart-Vision-Works/linmot_client/client/motion_control/streaming"
	client_vai "github.com/Smart-Vision-Works/linmot_client/client/motion_control/vai"
	client_vai16 "github.com/Smart-Vision-Works/linmot_client/client/motion_control/vai16"
	protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"
)

// MotionControlManager handles MC (Motion Command) Interface operations.
type MotionControlManager struct {
	requestManager          *client_common.RequestManager
	vaiManager              *client_vai.VAIManager
	predefVAIManager        *client_predefined.PredefVAIManager
	vai16Manager            *client_vai16.VAI16Manager
	interfaceControlManager *client_interface_control.InterfaceControlManager
	streamingManager        *client_streaming.StreamingManager
}

// NewMotionControlManager creates a new Motion Control manager.
func NewMotionControlManager(requestManager *client_common.RequestManager) *MotionControlManager {
	if requestManager == nil {
		panic("requestManager cannot be nil")
	}
	return &MotionControlManager{
		requestManager:          requestManager,
		vaiManager:              client_vai.NewVAIManager(requestManager),
		predefVAIManager:        client_predefined.NewPredefVAIManager(requestManager),
		vai16Manager:            client_vai16.NewVAI16Manager(requestManager),
		interfaceControlManager: client_interface_control.NewInterfaceControlManager(requestManager),
		streamingManager:        client_streaming.NewStreamingManager(requestManager),
	}
}

// ============================================================================
// VAI Methods (delegated to VAIManager)
// ============================================================================

// VAIGoToPosition sends a VAI Go To Position command.
func (m *MotionControlManager) VAIGoToPosition(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.GoToPosition(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAIIncrementDemandPosition sends a VAI Increment Demand Position command.
func (m *MotionControlManager) VAIIncrementDemandPosition(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.IncrementDemandPosition(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAIIncrementTargetPosition sends a VAI Increment Target Position command.
func (m *MotionControlManager) VAIIncrementTargetPosition(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.IncrementTargetPosition(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAIStop sends a VAI Stop command.
func (m *MotionControlManager) VAIStop(ctx context.Context) (*protocol_common.Status, error) {
	return m.vaiManager.Stop(ctx)
}

// VAIGoToPositionFromActual sends a VAI Go To Position From Actual Position command.
func (m *MotionControlManager) VAIGoToPositionFromActual(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.GoToPositionFromActual(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAIGoToPositionFromActualDemVelZero sends a VAI Go To Position From Actual with start velocity = 0 command.
func (m *MotionControlManager) VAIGoToPositionFromActualDemVelZero(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.GoToPositionFromActualDemVelZero(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAIIncrementActualPosition sends a VAI Increment Actual Position command.
func (m *MotionControlManager) VAIIncrementActualPosition(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.IncrementActualPosition(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAIIncrementActualPositionDemVelZero sends a VAI Increment Actual Position with start velocity = 0 command.
func (m *MotionControlManager) VAIIncrementActualPositionDemVelZero(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.IncrementActualPositionDemVelZero(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAIGoToPositionAfterActualCommand sends a VAI Go To Position After Actual Command.
func (m *MotionControlManager) VAIGoToPositionAfterActualCommand(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.GoToPositionAfterActualCommand(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAIGoToAnalogPosition sends a VAI Go To Analog Position command.
func (m *MotionControlManager) VAIGoToAnalogPosition(ctx context.Context, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.GoToAnalogPosition(ctx, velocityMS, accelMS2, decelMS2)
}

// VAIGoToPositionOnRisingTrigger sends a VAI Go To Position On Rising Trigger Event command.
func (m *MotionControlManager) VAIGoToPositionOnRisingTrigger(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.GoToPositionOnRisingTrigger(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAIIncrementTargetPositionOnRisingTrigger sends a VAI Increment Target Position On Rising Trigger Event command.
func (m *MotionControlManager) VAIIncrementTargetPositionOnRisingTrigger(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.IncrementTargetPositionOnRisingTrigger(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAIGoToPositionOnFallingTrigger sends a VAI Go To Position On Falling Trigger Event command.
func (m *MotionControlManager) VAIGoToPositionOnFallingTrigger(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.GoToPositionOnFallingTrigger(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAIIncrementTargetPositionOnFallingTrigger sends a VAI Increment Target Position On Falling Trigger Event command.
func (m *MotionControlManager) VAIIncrementTargetPositionOnFallingTrigger(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.IncrementTargetPositionOnFallingTrigger(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAIChangeMotionParamsOnPositiveTransition sends a VAI Change Motion Parameters On Positive Position Transition command.
func (m *MotionControlManager) VAIChangeMotionParamsOnPositiveTransition(ctx context.Context, transitionPosMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.ChangeMotionParamsOnPositiveTransition(ctx, transitionPosMM, velocityMS, accelMS2, decelMS2)
}

// VAIChangeMotionParamsOnNegativeTransition sends a VAI Change Motion Parameters On Negative Position Transition command.
func (m *MotionControlManager) VAIChangeMotionParamsOnNegativeTransition(ctx context.Context, transitionPosMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vaiManager.ChangeMotionParamsOnNegativeTransition(ctx, transitionPosMM, velocityMS, accelMS2, decelMS2)
}

// ============================================================================
// Predefined VAI Methods (delegated to PredefVAIManager)
// ============================================================================

// PredefVAIGoToPosition sends a Predefined VAI Go To Position command.
func (m *MotionControlManager) PredefVAIGoToPosition(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	return m.predefVAIManager.GoToPosition(ctx, positionMM)
}

// PredefVAIIncrementDemandPosition sends a Predefined VAI Increment Demand Position command.
func (m *MotionControlManager) PredefVAIIncrementDemandPosition(ctx context.Context, incrementMM float64) (*protocol_common.Status, error) {
	return m.predefVAIManager.IncrementDemandPosition(ctx, incrementMM)
}

// PredefVAIIncrementTargetPosition sends a Predefined VAI Increment Target Position command.
func (m *MotionControlManager) PredefVAIIncrementTargetPosition(ctx context.Context, incrementMM float64) (*protocol_common.Status, error) {
	return m.predefVAIManager.IncrementTargetPosition(ctx, incrementMM)
}

// PredefVAIGoToPositionFromActual sends a Predefined VAI Go To Position From Actual command.
func (m *MotionControlManager) PredefVAIGoToPositionFromActual(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	return m.predefVAIManager.GoToPositionFromActual(ctx, positionMM)
}

// PredefVAIGoToPositionFromActualDemVelZero sends a Predefined VAI Go To Position Starting With Dem Vel = 0 command.
func (m *MotionControlManager) PredefVAIGoToPositionFromActualDemVelZero(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	return m.predefVAIManager.GoToPositionFromActualDemVelZero(ctx, positionMM)
}

// PredefVAIStop sends a Predefined VAI Stop command.
func (m *MotionControlManager) PredefVAIStop(ctx context.Context) (*protocol_common.Status, error) {
	return m.predefVAIManager.Stop(ctx)
}

// PredefVAIGoToPositionAfterActualCommand sends a Predefined VAI Go To Position After Actual Command.
func (m *MotionControlManager) PredefVAIGoToPositionAfterActualCommand(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	return m.predefVAIManager.GoToPositionAfterActualCommand(ctx, positionMM)
}

// PredefVAIGoToPositionOnRisingTrigger sends a Predefined VAI Go To Position On Rising Trigger Event command.
func (m *MotionControlManager) PredefVAIGoToPositionOnRisingTrigger(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	return m.predefVAIManager.GoToPositionOnRisingTrigger(ctx, positionMM)
}

// PredefVAIIncrementTargetPositionOnRisingTrigger sends a Predefined VAI Increment Target Position On Rising Trigger Event command.
func (m *MotionControlManager) PredefVAIIncrementTargetPositionOnRisingTrigger(ctx context.Context, incrementMM float64) (*protocol_common.Status, error) {
	return m.predefVAIManager.IncrementTargetPositionOnRisingTrigger(ctx, incrementMM)
}

// PredefVAIGoToPositionOnFallingTrigger sends a Predefined VAI Go To Position On Falling Trigger Event command.
func (m *MotionControlManager) PredefVAIGoToPositionOnFallingTrigger(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	return m.predefVAIManager.GoToPositionOnFallingTrigger(ctx, positionMM)
}

// PredefVAIIncrementTargetPositionOnFallingTrigger sends a Predefined VAI Increment Target Position On Falling Trigger Event command.
func (m *MotionControlManager) PredefVAIIncrementTargetPositionOnFallingTrigger(ctx context.Context, incrementMM float64) (*protocol_common.Status, error) {
	return m.predefVAIManager.IncrementTargetPositionOnFallingTrigger(ctx, incrementMM)
}

// PredefVAIInfiniteMotionPositive sends a Predefined VAI Infinite Motion Positive Direction command.
func (m *MotionControlManager) PredefVAIInfiniteMotionPositive(ctx context.Context) (*protocol_common.Status, error) {
	return m.predefVAIManager.InfiniteMotionPositive(ctx)
}

// PredefVAIInfiniteMotionNegative sends a Predefined VAI Infinite Motion Negative Direction command.
func (m *MotionControlManager) PredefVAIInfiniteMotionNegative(ctx context.Context) (*protocol_common.Status, error) {
	return m.predefVAIManager.InfiniteMotionNegative(ctx)
}

// ============================================================================
// 16-Bit VAI Methods (delegated to VAI16Manager)
// ============================================================================

// VAI16GoToPosition sends a 16-bit VAI Go To Position command.
func (m *MotionControlManager) VAI16GoToPosition(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.GoToPosition(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAI16IncrementDemandPosition sends a 16-bit VAI Increment Demand Position command.
func (m *MotionControlManager) VAI16IncrementDemandPosition(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.IncrementDemandPosition(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAI16IncrementTargetPosition sends a 16-bit VAI Increment Target Position command.
func (m *MotionControlManager) VAI16IncrementTargetPosition(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.IncrementTargetPosition(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAI16GoToPositionFromActual sends a 16-bit VAI Go To Position From Actual Position command.
func (m *MotionControlManager) VAI16GoToPositionFromActual(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.GoToPositionFromActual(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAI16GoToPositionFromActualDemVelZero sends a 16-bit VAI Go To Position From Actual with start velocity = 0 command.
func (m *MotionControlManager) VAI16GoToPositionFromActualDemVelZero(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.GoToPositionFromActualDemVelZero(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAI16IncrementActualPosition sends a 16-bit VAI Increment Actual Position command.
func (m *MotionControlManager) VAI16IncrementActualPosition(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.IncrementActualPosition(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAI16IncrementActualPositionDemVelZero sends a 16-bit VAI Increment Actual Position with start velocity = 0 command.
func (m *MotionControlManager) VAI16IncrementActualPositionDemVelZero(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.IncrementActualPositionDemVelZero(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAI16Stop sends a 16-bit VAI Stop command.
func (m *MotionControlManager) VAI16Stop(ctx context.Context) (*protocol_common.Status, error) {
	return m.vai16Manager.Stop(ctx)
}

// VAI16GoToPositionAfterActualCommand sends a 16-bit VAI Go To Position After Actual Command.
func (m *MotionControlManager) VAI16GoToPositionAfterActualCommand(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.GoToPositionAfterActualCommand(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAI16GoToPositionOnRisingTrigger sends a 16-bit VAI Go To Position On Rising Trigger Event command.
func (m *MotionControlManager) VAI16GoToPositionOnRisingTrigger(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.GoToPositionOnRisingTrigger(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAI16IncrementTargetPositionOnRisingTrigger sends a 16-bit VAI Increment Target Position On Rising Trigger Event command.
func (m *MotionControlManager) VAI16IncrementTargetPositionOnRisingTrigger(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.IncrementTargetPositionOnRisingTrigger(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAI16GoToPositionOnFallingTrigger sends a 16-bit VAI Go To Position On Falling Trigger Event command.
func (m *MotionControlManager) VAI16GoToPositionOnFallingTrigger(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.GoToPositionOnFallingTrigger(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAI16IncrementTargetPositionOnFallingTrigger sends a 16-bit VAI Increment Target Position On Falling Trigger Event command.
func (m *MotionControlManager) VAI16IncrementTargetPositionOnFallingTrigger(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.IncrementTargetPositionOnFallingTrigger(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAI16ChangeMotionParamsOnPositiveTransition sends a 16-bit VAI Change Motion Parameters On Positive Position Transition command.
func (m *MotionControlManager) VAI16ChangeMotionParamsOnPositiveTransition(ctx context.Context, transitionPosMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.ChangeMotionParamsOnPositiveTransition(ctx, transitionPosMM, velocityMS, accelMS2, decelMS2)
}

// VAI16ChangeMotionParamsOnNegativeTransition sends a 16-bit VAI Change Motion Parameters On Negative Position Transition command.
func (m *MotionControlManager) VAI16ChangeMotionParamsOnNegativeTransition(ctx context.Context, transitionPosMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return m.vai16Manager.ChangeMotionParamsOnNegativeTransition(ctx, transitionPosMM, velocityMS, accelMS2, decelMS2)
}

// ============================================================================
// Interface Control Methods (delegated to InterfaceControlManager)
// ============================================================================

// NoOperation sends a No Operation command.
func (m *MotionControlManager) NoOperation(ctx context.Context) (*protocol_common.Status, error) {
	return m.interfaceControlManager.NoOperation(ctx)
}

// WriteInterfaceControlWord sends a Write Interface Control Word command.
func (m *MotionControlManager) WriteInterfaceControlWord(ctx context.Context, controlWord uint16) (*protocol_common.Status, error) {
	return m.interfaceControlManager.WriteInterfaceControlWord(ctx, controlWord)
}

// WriteLiveParameter sends a Write Live Parameter command.
func (m *MotionControlManager) WriteLiveParameter(ctx context.Context, upid uint16, value uint32) (*protocol_common.Status, error) {
	return m.interfaceControlManager.WriteLiveParameter(ctx, upid, value)
}

// WriteOutputsWithMask sends a Write X4/X14 Interface Outputs with Mask command.
func (m *MotionControlManager) WriteOutputsWithMask(ctx context.Context, bitMask, bitValue uint16) (*protocol_common.Status, error) {
	return m.interfaceControlManager.WriteOutputsWithMask(ctx, bitMask, bitValue)
}

// SelectPositionControllerSet sends a Select Position Controller Set command.
func (m *MotionControlManager) SelectPositionControllerSet(ctx context.Context, controllerSet uint16) (*protocol_common.Status, error) {
	return m.interfaceControlManager.SelectPositionControllerSet(ctx, controllerSet)
}

// ClearEventEvaluation sends a Clear Event Evaluation command.
func (m *MotionControlManager) ClearEventEvaluation(ctx context.Context) (*protocol_common.Status, error) {
	return m.interfaceControlManager.ClearEventEvaluation(ctx)
}

// MasterHoming sends a Master Homing command.
func (m *MotionControlManager) MasterHoming(ctx context.Context, homePositionMM float64) (*protocol_common.Status, error) {
	return m.interfaceControlManager.MasterHoming(ctx, homePositionMM)
}

// Reset sends a Reset command. WARNING: This will reboot the drive!
func (m *MotionControlManager) Reset(ctx context.Context) (*protocol_common.Status, error) {
	return m.interfaceControlManager.Reset(ctx)
}

// ============================================================================
// Streaming Methods (delegated to StreamingManager)
// ============================================================================

// SendPStream sends a Position Stream command.
func (m *MotionControlManager) SendPStream(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	return m.streamingManager.SendPStream(ctx, positionMM)
}

// SendPVStream sends a Position-Velocity Stream command.
func (m *MotionControlManager) SendPVStream(ctx context.Context, positionMM, velocityMS float64) (*protocol_common.Status, error) {
	return m.streamingManager.SendPVStream(ctx, positionMM, velocityMS)
}

// SendPStreamConfigPeriod sends a P Stream With Configured Period command.
func (m *MotionControlManager) SendPStreamConfigPeriod(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	return m.streamingManager.SendPStreamConfigPeriod(ctx, positionMM)
}

// SendPVStreamConfigPeriod sends a PV Stream With Configured Period command.
func (m *MotionControlManager) SendPVStreamConfigPeriod(ctx context.Context, positionMM, velocityMS float64) (*protocol_common.Status, error) {
	return m.streamingManager.SendPVStreamConfigPeriod(ctx, positionMM, velocityMS)
}

// SendPVAStream sends a Position-Velocity-Acceleration Stream command.
func (m *MotionControlManager) SendPVAStream(ctx context.Context, positionMM, velocityMS, accelMS2 float64) (*protocol_common.Status, error) {
	return m.streamingManager.SendPVAStream(ctx, positionMM, velocityMS, accelMS2)
}

// StopStreaming sends a Stop Streaming command.
func (m *MotionControlManager) StopStreaming(ctx context.Context) (*protocol_common.Status, error) {
	return m.streamingManager.StopStreaming(ctx)
}
