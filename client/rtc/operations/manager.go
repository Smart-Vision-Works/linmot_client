package client_operations

import (
	"context"

	client_common "github.com/Smart-Vision-Works/linmot_client/client/common"
	protocol_operations "github.com/Smart-Vision-Works/linmot_client/protocol/rtc/operations"
)

type OperationManager struct {
	requestManager *client_common.RequestManager
}

func NewOperationManager(requestManager *client_common.RequestManager) *OperationManager {
	return &OperationManager{
		requestManager: requestManager,
	}
}

// RestartDrive restarts the drive.
// WARNING: This is a dangerous operation that will restart the entire drive.
func (m *OperationManager) RestartDrive(ctx context.Context) error {
	request := protocol_operations.NewRestartDriveRequest()
	_, err := client_common.SendRequestAndReceive[*protocol_operations.RestartDriveResponse](m.requestManager, ctx, request)
	return err
}

// SetOSROMToDefault resets OS SW parameter ROM values to default.
func (m *OperationManager) SetOSROMToDefault(ctx context.Context) error {
	request := protocol_operations.NewSetOSROMToDefaultRequest()
	_, err := client_common.SendRequestAndReceive[*protocol_operations.SetOSROMToDefaultResponse](m.requestManager, ctx, request)
	return err
}

// SetMCROMToDefault resets MC SW parameter ROM values to default.
func (m *OperationManager) SetMCROMToDefault(ctx context.Context) error {
	request := protocol_operations.NewSetMCROMToDefaultRequest()
	_, err := client_common.SendRequestAndReceive[*protocol_operations.SetMCROMToDefaultResponse](m.requestManager, ctx, request)
	return err
}

// SetInterfaceROMToDefault resets Interface SW parameter ROM values to default.
func (m *OperationManager) SetInterfaceROMToDefault(ctx context.Context) error {
	request := protocol_operations.NewSetInterfaceROMToDefaultRequest()
	_, err := client_common.SendRequestAndReceive[*protocol_operations.SetInterfaceROMToDefaultResponse](m.requestManager, ctx, request)
	return err
}

// SetApplicationROMToDefault resets Application SW parameter ROM values to default.
func (m *OperationManager) SetApplicationROMToDefault(ctx context.Context) error {
	request := protocol_operations.NewSetApplicationROMToDefaultRequest()
	_, err := client_common.SendRequestAndReceive[*protocol_operations.SetApplicationROMToDefaultResponse](m.requestManager, ctx, request)
	return err
}
