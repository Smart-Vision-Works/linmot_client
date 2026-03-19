package client_rtc

import (
	"context"

	client_common "github.com/Smart-Vision-Works/linmot_client/client/common"
	client_command_tables "github.com/Smart-Vision-Works/linmot_client/client/rtc/command_tables"
	client_curves "github.com/Smart-Vision-Works/linmot_client/client/rtc/curves"
	client_errors "github.com/Smart-Vision-Works/linmot_client/client/rtc/errors"
	client_operations "github.com/Smart-Vision-Works/linmot_client/client/rtc/operations"
	client_parameters "github.com/Smart-Vision-Works/linmot_client/client/rtc/parameters"
	protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"
)

type RtcManager struct {
	requestManager *client_common.RequestManager

	commandTableManager *client_command_tables.CommandTableManager
	curveManager        *client_curves.CurveManager
	errorManager        *client_errors.ErrorManager
	operationManager    *client_operations.OperationManager
	parameterManager    *client_parameters.ParameterManager
}

func NewRtcManager(requestManager *client_common.RequestManager) *RtcManager {
	return &RtcManager{
		requestManager:      requestManager,
		commandTableManager: client_command_tables.NewCommandTableManager(requestManager),
		curveManager:        client_curves.NewCurveManager(requestManager),
		errorManager:        client_errors.NewErrorManager(requestManager),
		operationManager:    client_operations.NewOperationManager(requestManager),
		parameterManager:    client_parameters.NewParameterManager(requestManager),
	}
}

// SetDebug enables or disables debug logging for RTC operations.
func (manager *RtcManager) SetDebug(enabled bool) {
	manager.commandTableManager.SetDebug(enabled)
}

// GetCommandTable retrieves the current command table from the drive.
func (manager *RtcManager) GetCommandTable(ctx context.Context) (*client_command_tables.CommandTable, error) {
	return manager.commandTableManager.GetCommandTable(ctx)
}

// GetPresenceMasks retrieves the presence mask values from the drive.
// Returns 8 uint32 masks indicating which entries exist in each range.
func (manager *RtcManager) GetPresenceMasks(ctx context.Context) ([8]uint32, error) {
	return manager.commandTableManager.GetPresenceMasks(ctx)
}

// SetCommandTable sets the command table on the drive from a CommandTable struct.
func (manager *RtcManager) SetCommandTable(ctx context.Context, commandTable *client_command_tables.CommandTable) error {
	return manager.commandTableManager.SetCommandTable(ctx, commandTable)
}

// SetCommandTableWithOptions sets the command table on the drive from a CommandTable struct with configurable options.
func (manager *RtcManager) SetCommandTableWithOptions(ctx context.Context, commandTable *client_command_tables.CommandTable, opts client_command_tables.SetCommandTableOptions) error {
	return manager.commandTableManager.SetCommandTableWithOptions(ctx, commandTable, opts)
}

// StopMotionController stops the Motion Controller software on the drive.
func (manager *RtcManager) StopMotionController(ctx context.Context) error {
	return manager.commandTableManager.StopMotionController(ctx)
}

// StartMotionController restarts the Motion Controller software on the drive.
func (manager *RtcManager) StartMotionController(ctx context.Context) error {
	return manager.commandTableManager.StartMotionController(ctx)
}

// SaveCommandTableToFlash sends the SaveCommandTable (0x80) RTC command.
func (manager *RtcManager) SaveCommandTableToFlash(ctx context.Context) error {
	return manager.commandTableManager.SaveCommandTableToFlash(ctx)
}

// ReadRAM reads the current RAM value of a parameter by UPID.
func (manager *RtcManager) ReadRAM(ctx context.Context, upid uint16) (int32, error) {
	return manager.parameterManager.ReadRAM(ctx, upid)
}

// WriteRAMAndROM writes both RAM and ROM value of a parameter.
func (manager *RtcManager) WriteRAMAndROM(ctx context.Context, upid uint16, value int32) error {
	return manager.parameterManager.WriteRAMAndROM(ctx, upid, value)
}

// GetParameterMinValue gets the minimum allowed value for a parameter.
func (manager *RtcManager) GetParameterMinValue(ctx context.Context, upid uint16) (int32, error) {
	return manager.parameterManager.GetMinValue(ctx, upid)
}

// GetParameterMaxValue gets the maximum allowed value for a parameter.
func (manager *RtcManager) GetParameterMaxValue(ctx context.Context, upid uint16) (int32, error) {
	return manager.parameterManager.GetMaxValue(ctx, upid)
}

// GetParameterDefaultValue gets the default value for a parameter.
func (manager *RtcManager) GetParameterDefaultValue(ctx context.Context, upid uint16) (int32, error) {
	return manager.parameterManager.GetDefaultValue(ctx, upid)
}

// GetAllParameterIDs returns a list of all available parameter UPIDs.
func (manager *RtcManager) GetAllParameterIDs(ctx context.Context) ([]uint16, error) {
	return manager.parameterManager.GetAllUPIDs(ctx)
}

// GetModifiedParameterIDs returns a list of all modified parameter UPIDs.
func (manager *RtcManager) GetModifiedParameterIDs(ctx context.Context) ([]uint16, error) {
	return manager.parameterManager.GetModifiedUPIDs(ctx)
}

// GetAllParameters returns all available parameters with address usage information.
func (manager *RtcManager) GetAllParameters(ctx context.Context) ([]client_parameters.ParameterInfo, error) {
	return manager.parameterManager.GetAllParameters(ctx)
}

// GetModifiedParameters returns all modified parameters with address usage information.
func (manager *RtcManager) GetModifiedParameters(ctx context.Context) ([]client_parameters.ParameterInfo, error) {
	return manager.parameterManager.GetModifiedParameters(ctx)
}

// RestartDrive restarts the drive.
func (manager *RtcManager) RestartDrive(ctx context.Context) error {
	return manager.operationManager.RestartDrive(ctx)
}

// ResetOSParametersToDefault resets OS parameters to factory defaults.
func (manager *RtcManager) ResetOSParametersToDefault(ctx context.Context) error {
	return manager.operationManager.SetOSROMToDefault(ctx)
}

// ResetMCParametersToDefault resets MC parameters to factory defaults.
func (manager *RtcManager) ResetMCParametersToDefault(ctx context.Context) error {
	return manager.operationManager.SetMCROMToDefault(ctx)
}

// ResetInterfaceParametersToDefault resets interface parameters to factory defaults.
func (manager *RtcManager) ResetInterfaceParametersToDefault(ctx context.Context) error {
	return manager.operationManager.SetInterfaceROMToDefault(ctx)
}

// ResetApplicationParametersToDefault resets application parameters to factory defaults.
func (manager *RtcManager) ResetApplicationParametersToDefault(ctx context.Context) error {
	return manager.operationManager.SetApplicationROMToDefault(ctx)
}

// SaveAllCurves saves all curves from RAM to Flash.
func (manager *RtcManager) SaveAllCurves(ctx context.Context) error {
	return manager.curveManager.SaveAllCurves(ctx)
}

// DeleteAllCurves deletes all curves from RAM.
func (manager *RtcManager) DeleteAllCurves(ctx context.Context) error {
	return manager.curveManager.DeleteAllCurves(ctx)
}

// UploadCurve uploads a complete curve to the drive.
func (manager *RtcManager) UploadCurve(ctx context.Context, curveID uint16, infoBlock, dataBlock []byte) error {
	return manager.curveManager.UploadCurve(ctx, curveID, infoBlock, dataBlock)
}

// DownloadCurve downloads a complete curve from the drive.
func (manager *RtcManager) DownloadCurve(ctx context.Context, curveID uint16) ([]byte, []byte, error) {
	return manager.curveManager.DownloadCurve(ctx, curveID)
}

// ModifyCurve modifies an existing curve on the drive.
func (manager *RtcManager) ModifyCurve(ctx context.Context, curveID uint16, infoBlock, dataBlock []byte) error {
	return manager.curveManager.ModifyCurve(ctx, curveID, infoBlock, dataBlock)
}

// GetErrorLog retrieves the error log from the drive.
func (manager *RtcManager) GetErrorLog(ctx context.Context) ([]client_errors.ErrorLogEntry, error) {
	return manager.errorManager.GetErrorLog(ctx)
}

// GetErrorLogCounts returns the count of logged and occurred errors.
func (manager *RtcManager) GetErrorLogCounts(ctx context.Context) (logged, occurred uint16, err error) {
	return manager.errorManager.GetErrorLogCounts(ctx)
}

// GetErrorLogEntry retrieves a single error log entry by index.
func (manager *RtcManager) GetErrorLogEntry(ctx context.Context, entryNum uint16) (*client_errors.ErrorLogEntry, error) {
	return manager.errorManager.GetErrorLogEntry(ctx, entryNum)
}

// GetErrorText retrieves the human-readable description for an error code.
func (manager *RtcManager) GetErrorText(ctx context.Context, errorCode uint16) (string, error) {
	return manager.errorManager.GetErrorText(ctx, errorCode)
}

// GetErrorLogWithText retrieves the error log with human-readable descriptions.
func (manager *RtcManager) GetErrorLogWithText(ctx context.Context) ([]client_errors.ErrorLogEntry, error) {
	return manager.errorManager.GetErrorLogWithText(ctx)
}

// ============================================================================
// Motion Parameter Methods (delegated to ParameterManager)
// ============================================================================
// These methods provide convenient access to motion control parameters via RTC.

// SetPosition1 sets the target position for Position 1 in millimeters.
func (manager *RtcManager) SetPosition1(ctx context.Context, positionMM float64, storageType protocol_common.ParameterStorageType) error {
	return manager.parameterManager.SetPosition1(ctx, positionMM, storageType)
}

// SetPosition2 sets the target position for Position 2 in millimeters.
func (manager *RtcManager) SetPosition2(ctx context.Context, positionMM float64, storageType protocol_common.ParameterStorageType) error {
	return manager.parameterManager.SetPosition2(ctx, positionMM, storageType)
}

// GetVelocity returns the maximum velocity in meters per second.
func (manager *RtcManager) GetVelocity(ctx context.Context) (float64, error) {
	return manager.parameterManager.GetVelocity(ctx)
}

// SetVelocity sets the maximum velocity in meters per second.
func (manager *RtcManager) SetVelocity(ctx context.Context, velocityMS float64, storageType protocol_common.ParameterStorageType) error {
	return manager.parameterManager.SetVelocity(ctx, velocityMS, storageType)
}

// GetAcceleration returns the acceleration in meters per second squared.
func (manager *RtcManager) GetAcceleration(ctx context.Context) (float64, error) {
	return manager.parameterManager.GetAcceleration(ctx)
}

// SetAcceleration sets the acceleration in meters per second squared.
func (manager *RtcManager) SetAcceleration(ctx context.Context, accelMS2 float64, storageType protocol_common.ParameterStorageType) error {
	return manager.parameterManager.SetAcceleration(ctx, accelMS2, storageType)
}

// GetDeceleration returns the deceleration in meters per second squared.
func (manager *RtcManager) GetDeceleration(ctx context.Context) (float64, error) {
	return manager.parameterManager.GetDeceleration(ctx)
}

// SetDeceleration sets the deceleration in meters per second squared.
func (manager *RtcManager) SetDeceleration(ctx context.Context, decelMS2 float64, storageType protocol_common.ParameterStorageType) error {
	return manager.parameterManager.SetDeceleration(ctx, decelMS2, storageType)
}

// SetEasyStepsAutoStart sets the Easy Steps auto start configuration.
func (manager *RtcManager) SetEasyStepsAutoStart(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error {
	return manager.parameterManager.SetEasyStepsAutoStart(ctx, value, storageType)
}

// SetEasyStepsAutoHome sets the Easy Steps auto home configuration.
func (manager *RtcManager) SetEasyStepsAutoHome(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error {
	return manager.parameterManager.SetEasyStepsAutoHome(ctx, value, storageType)
}

// SetEasyStepsRisingEdge sets the Easy Steps rising edge action for an input pin.
func (manager *RtcManager) SetEasyStepsRisingEdge(ctx context.Context, inputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error {
	return manager.parameterManager.SetEasyStepsRisingEdge(ctx, inputNumber, value, storageType)
}

// SetEasyStepsIOMotionConfigCmd sets the Easy Steps IO motion config curve/CMD ID.
func (manager *RtcManager) SetEasyStepsIOMotionConfigCmd(ctx context.Context, inputNumber protocol_common.IOPinNumber, curveCmdID int32, storageType protocol_common.ParameterStorageType) error {
	return manager.parameterManager.SetEasyStepsIOMotionConfigCmd(ctx, inputNumber, curveCmdID, storageType)
}

// SetOutputFunction sets the output pin function configuration.
func (manager *RtcManager) SetOutputFunction(ctx context.Context, outputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error {
	return manager.parameterManager.SetOutputFunction(ctx, outputNumber, value, storageType)
}

// SetInputFunction sets the input pin function configuration.
func (manager *RtcManager) SetInputFunction(ctx context.Context, inputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error {
	return manager.parameterManager.SetInputFunction(ctx, inputNumber, value, storageType)
}

// SetRunMode sets the run mode configuration.
func (manager *RtcManager) SetRunMode(ctx context.Context, mode protocol_common.RunMode, storageType protocol_common.ParameterStorageType) error {
	return manager.parameterManager.SetRunMode(ctx, mode, storageType)
}

// SetTriggerMode sets the trigger mode configuration.
func (manager *RtcManager) SetTriggerMode(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error {
	return manager.parameterManager.SetTriggerMode(ctx, value, storageType)
}
