package client_monitoring

import (
	"context"
	"fmt"

	client_common "github.com/Smart-Vision-Works/linmot_client/client/common"
	protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/linmot_client/protocol/rtc"
)

// MonitoringManager handles monitoring channel configuration and data retrieval.
type MonitoringManager struct {
	requestManager *client_common.RequestManager
}

// NewMonitoringManager creates a new MonitoringManager.
func NewMonitoringManager(requestManager *client_common.RequestManager) *MonitoringManager {
	return &MonitoringManager{
		requestManager: requestManager,
	}
}

// ConfigureChannel configures a single monitoring channel to monitor a specific UPID.
// channelNum must be 1-4.
// upid is the parameter to monitor (e.g., position, velocity, current).
func (m *MonitoringManager) ConfigureChannel(ctx context.Context, channelNum int, upid uint16) error {
	if channelNum < 1 || channelNum > 4 {
		return fmt.Errorf("invalid channel number %d: must be 1-4", channelNum)
	}

	// Map channel number (1-4) to configuration UPID (0x20A8-0x20AB)
	configUPID := uint16(protocol_common.PUID.MonitoringChannel1UPID) + uint16(channelNum-1)

	// Write the target UPID to the monitoring channel configuration parameter
	request := protocol_rtc.NewRTCSetParamRequest(configUPID, int32(upid), protocol_rtc.CommandCode.WriteRAM)
	_, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCSetParamResponse](m.requestManager, ctx, request)
	return err
}

// ConfigureChannels bulk configures all 4 monitoring channels.
// upids[0] configures channel 1, upids[1] configures channel 2, etc.
func (m *MonitoringManager) ConfigureChannels(ctx context.Context, upids [4]uint16) error {
	for i := 0; i < 4; i++ {
		if err := m.ConfigureChannel(ctx, i+1, upids[i]); err != nil {
			return fmt.Errorf("failed to configure channel %d: %w", i+1, err)
		}
	}
	return nil
}

// GetChannelConfiguration reads which UPID is assigned to a monitoring channel.
// channelNum must be 1-4.
func (m *MonitoringManager) GetChannelConfiguration(ctx context.Context, channelNum int) (uint16, error) {
	if channelNum < 1 || channelNum > 4 {
		return 0, fmt.Errorf("invalid channel number %d: must be 1-4", channelNum)
	}

	// Map channel number (1-4) to configuration UPID (0x20A8-0x20AB)
	configUPID := uint16(protocol_common.PUID.MonitoringChannel1UPID) + uint16(channelNum-1)

	// Read the target UPID from the monitoring channel configuration parameter
	request := protocol_rtc.NewRTCGetParamRequest(configUPID, protocol_rtc.CommandCode.ReadRAM)
	response, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](m.requestManager, ctx, request)
	if err != nil {
		return 0, fmt.Errorf("failed to read channel %d configuration: %w", channelNum, err)
	}

	return uint16(response.Value()), nil
}

// GetAllChannelConfigurations returns the UPIDs configured for all 4 monitoring channels.
func (m *MonitoringManager) GetAllChannelConfigurations(ctx context.Context) ([4]uint16, error) {
	var upids [4]uint16
	for i := 0; i < 4; i++ {
		upid, err := m.GetChannelConfiguration(ctx, i+1)
		if err != nil {
			return upids, err
		}
		upids[i] = upid
	}
	return upids, nil
}

// MonitoringSnapshot contains monitoring channel data with convenience accessors.
type MonitoringSnapshot struct {
	Status        *protocol_common.Status
	Channel1Value int32
	Channel2Value int32
	Channel3Value int32
	Channel4Value int32
}

// GetMonitoringData retrieves drive status with monitoring channel data.
// Returns a Status struct with MonitoringChannel field populated.
func (m *MonitoringManager) GetMonitoringData(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_common.NewMonitoringStatusRequest()
	response, err := client_common.SendRequestAndReceive[*protocol_common.MonitoringStatusResponse](m.requestManager, ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve monitoring data: %w", err)
	}

	return response.Status(), nil
}

// GetMonitoringSnapshot retrieves monitoring data as a convenient snapshot struct.
// This method retrieves status and monitoring channel values in one call.
func (m *MonitoringManager) GetMonitoringSnapshot(ctx context.Context) (*MonitoringSnapshot, error) {
	status, err := m.GetMonitoringData(ctx)
	if err != nil {
		return nil, err
	}

	return &MonitoringSnapshot{
		Status:        status,
		Channel1Value: status.MonitoringChannel[0],
		Channel2Value: status.MonitoringChannel[1],
		Channel3Value: status.MonitoringChannel[2],
		Channel4Value: status.MonitoringChannel[3],
	}, nil
}
