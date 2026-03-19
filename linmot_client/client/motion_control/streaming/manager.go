package client_streaming

import (
	"context"
	"fmt"

	client_common "github.com/Smart-Vision-Works/staged_robot/client/common"
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_motion_control "github.com/Smart-Vision-Works/staged_robot/protocol/motion_control"
	protocol_streaming "github.com/Smart-Vision-Works/staged_robot/protocol/motion_control/streaming"
)

// StreamingManager handles streaming motion commands.
type StreamingManager struct {
	requestManager *client_common.RequestManager
}

// NewStreamingManager creates a new Streaming manager.
func NewStreamingManager(requestManager *client_common.RequestManager) *StreamingManager {
	return &StreamingManager{
		requestManager: requestManager,
	}
}

// SendPStream sends a Position Stream command.
//
// Parameters:
//   - positionMM: Target position in millimeters
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.38
func (m *StreamingManager) SendPStream(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)

	request := protocol_streaming.NewPStreamCommand(position)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send P Stream command: %w", err)
	}

	return response.Status(), nil
}

// SendPVStream sends a Position-Velocity Stream command.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Velocity in meters per second
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.39
func (m *StreamingManager) SendPVStream(ctx context.Context, positionMM, velocityMS float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)
	velocity := protocol_common.ToProtocolVelocitySigned(velocityMS)

	request := protocol_streaming.NewPVStreamCommand(position, velocity)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send PV Stream command: %w", err)
	}

	return response.Status(), nil
}

// SendPStreamConfigPeriod sends a P Stream With Configured Period command.
//
// Parameters:
//   - positionMM: Target position in millimeters
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.40
func (m *StreamingManager) SendPStreamConfigPeriod(ctx context.Context, positionMM float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)

	request := protocol_streaming.NewPStreamConfigPeriodCommand(position)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send P Stream Config Period command: %w", err)
	}

	return response.Status(), nil
}

// SendPVStreamConfigPeriod sends a PV Stream With Configured Period command.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Velocity in meters per second
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.41
func (m *StreamingManager) SendPVStreamConfigPeriod(ctx context.Context, positionMM, velocityMS float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)
	velocity := protocol_common.ToProtocolVelocitySigned(velocityMS)

	request := protocol_streaming.NewPVStreamConfigPeriodCommand(position, velocity)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send PV Stream Config Period command: %w", err)
	}

	return response.Status(), nil
}

// SendPVAStream sends a Position-Velocity-Acceleration Stream command.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.42
func (m *StreamingManager) SendPVAStream(ctx context.Context, positionMM, velocityMS, accelMS2 float64) (*protocol_common.Status, error) {
	position := protocol_common.ToProtocolPosition(positionMM)
	velocity := protocol_common.ToProtocolVelocitySigned(velocityMS)
	acceleration := protocol_common.ToProtocolAccelerationSigned(accelMS2)

	request := protocol_streaming.NewPVAStreamCommand(position, velocity, acceleration)

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send PVA Stream command: %w", err)
	}

	return response.Status(), nil
}

// StopStreaming sends a Stop Streaming command.
//
// Reference: LinMot_MotionCtrl.txt, Section 4.3.43
func (m *StreamingManager) StopStreaming(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_streaming.NewStopStreamingCommand()

	response, err := client_common.SendRequestAndReceive[*protocol_motion_control.MCCommandResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send Stop Streaming command: %w", err)
	}

	return response.Status(), nil
}
