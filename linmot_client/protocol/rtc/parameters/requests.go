package protocol_parameters

import (
	"fmt"

	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/staged_robot/protocol/rtc"
)

// Compile-time interface checks
var (
	_ protocol_common.Request = (*WriteRAMAndROMRequest)(nil)
	_ protocol_common.Request = (*GetMinValueRequest)(nil)
	_ protocol_common.Request = (*GetMaxValueRequest)(nil)
	_ protocol_common.Request = (*GetDefaultValueRequest)(nil)
	_ protocol_common.Request = (*StartGettingUPIDListRequest)(nil)
	_ protocol_common.Request = (*GetNextUPIDListItemRequest)(nil)
	_ protocol_common.Request = (*StartGettingModifiedUPIDListRequest)(nil)
	_ protocol_common.Request = (*GetNextModifiedUPIDListItemRequest)(nil)
	// Motion parameter requests (RTC-based wrappers for motion-related UPIDs)
	_ protocol_common.Request = (*ReadVelocityRequest)(nil)
	_ protocol_common.Request = (*ReadPosition1Request)(nil)
	_ protocol_common.Request = (*ReadAccelerationRequest)(nil)
	_ protocol_common.Request = (*ReadDecelerationRequest)(nil)
	_ protocol_common.Request = (*WritePosition1Request)(nil)
	_ protocol_common.Request = (*WritePosition2Request)(nil)
	_ protocol_common.Request = (*WriteVelocityRequest)(nil)
	_ protocol_common.Request = (*WriteAccelerationRequest)(nil)
	_ protocol_common.Request = (*WriteDecelerationRequest)(nil)
	_ protocol_common.Request = (*WriteRunModeRequest)(nil)
	_ protocol_common.Request = (*WriteEasyStepsAutoStartRequest)(nil)
	_ protocol_common.Request = (*WriteEasyStepsAutoHomeRequest)(nil)
)

// WriteRAMAndROMRequest writes both RAM and ROM value of a parameter by UPID.
type WriteRAMAndROMRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWriteRAMAndROMRequest creates a new request to write both RAM and ROM parameter values.
func NewWriteRAMAndROMRequest(upid uint16, value int32) *WriteRAMAndROMRequest {
	return &WriteRAMAndROMRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(upid, value, protocol_rtc.CommandCode.WriteRAMAndROM),
	}
}

// GetMinValueRequest gets the minimal value of a parameter by UPID.
type GetMinValueRequest struct {
	protocol_rtc.RTCGetParamRequest
}

// NewGetMinValueRequest creates a new request to get the minimum value of a parameter.
func NewGetMinValueRequest(upid uint16) *GetMinValueRequest {
	return &GetMinValueRequest{
		RTCGetParamRequest: *protocol_rtc.NewRTCGetParamRequest(upid, protocol_rtc.CommandCode.GetMinValue),
	}
}

// GetMaxValueRequest gets the maximal value of a parameter by UPID.
type GetMaxValueRequest struct {
	protocol_rtc.RTCGetParamRequest
}

// NewGetMaxValueRequest creates a new request to get the maximum value of a parameter.
func NewGetMaxValueRequest(upid uint16) *GetMaxValueRequest {
	return &GetMaxValueRequest{
		RTCGetParamRequest: *protocol_rtc.NewRTCGetParamRequest(upid, protocol_rtc.CommandCode.GetMaxValue),
	}
}

// GetDefaultValueRequest gets the default value of a parameter by UPID.
type GetDefaultValueRequest struct {
	protocol_rtc.RTCGetParamRequest
}

// NewGetDefaultValueRequest creates a new request to get the default value of a parameter.
func NewGetDefaultValueRequest(upid uint16) *GetDefaultValueRequest {
	return &GetDefaultValueRequest{
		RTCGetParamRequest: *protocol_rtc.NewRTCGetParamRequest(upid, protocol_rtc.CommandCode.GetDefaultValue),
	}
}

// StartGettingUPIDListRequest starts UPID list iteration.
type StartGettingUPIDListRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewStartGettingUPIDListRequest creates a new request to start getting the UPID list.
// startUPID is the UPID to start searching from.
func NewStartGettingUPIDListRequest(startUPID uint16) *StartGettingUPIDListRequest {
	return &StartGettingUPIDListRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(startUPID, 0, protocol_rtc.CommandCode.StartGettingUPIDList),
	}
}

// GetNextUPIDListItemRequest gets the next UPID in the list.
type GetNextUPIDListItemRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewGetNextUPIDListItemRequest creates a new request to get the next UPID list item.
func NewGetNextUPIDListItemRequest() *GetNextUPIDListItemRequest {
	return &GetNextUPIDListItemRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(0, 0, protocol_rtc.CommandCode.GetNextUPIDListItem),
	}
}

// StartGettingModifiedUPIDListRequest starts modified UPID list iteration.
type StartGettingModifiedUPIDListRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewStartGettingModifiedUPIDListRequest creates a new request to start getting the modified UPID list.
// startUPID is the UPID to start searching from.
func NewStartGettingModifiedUPIDListRequest(startUPID uint16) *StartGettingModifiedUPIDListRequest {
	return &StartGettingModifiedUPIDListRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(startUPID, 0, protocol_rtc.CommandCode.StartGettingModifiedUPIDList),
	}
}

// GetNextModifiedUPIDListItemRequest gets the next modified UPID in the list.
type GetNextModifiedUPIDListItemRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewGetNextModifiedUPIDListItemRequest creates a new request to get the next modified UPID list item.
func NewGetNextModifiedUPIDListItemRequest() *GetNextModifiedUPIDListItemRequest {
	return &GetNextModifiedUPIDListItemRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(0, 0, protocol_rtc.CommandCode.GetNextModifiedUPIDListItem),
	}
}

// ============================================================================
// Motion Parameter Requests (RTC-based wrappers for motion-related UPIDs)
// ============================================================================
// These are convenience wrappers for commonly-used motion control parameters.
// They use RTC Parameter Access commands (0x10-0x13) under the hood.

// ReadVelocityRequest reads velocity (upid 0x145B).
type ReadVelocityRequest struct {
	protocol_rtc.RTCGetParamRequest
}

// NewReadVelocityRequest creates a new RTC read request for velocity (upid 0x145B).
func NewReadVelocityRequest() *ReadVelocityRequest {
	return &ReadVelocityRequest{
		RTCGetParamRequest: *protocol_rtc.NewRTCGetParamRequest(
			uint16(protocol_common.Parameter.Speed1),
			protocol_rtc.CommandCode.ReadRAM,
		),
	}
}

// ReadPosition1Request reads Position 1 (upid 0x145A).
type ReadPosition1Request struct {
	protocol_rtc.RTCGetParamRequest
}

// NewReadPosition1Request creates a new RTC read request for Position 1 (upid 0x145A).
func NewReadPosition1Request() *ReadPosition1Request {
	return &ReadPosition1Request{
		RTCGetParamRequest: *protocol_rtc.NewRTCGetParamRequest(
			uint16(protocol_common.Parameter.Position1),
			protocol_rtc.CommandCode.ReadRAM,
		),
	}
}

// ReadAccelerationRequest reads acceleration (upid 0x145C).
type ReadAccelerationRequest struct {
	protocol_rtc.RTCGetParamRequest
}

// NewReadAccelerationRequest creates a new RTC read request for acceleration (upid 0x145C).
func NewReadAccelerationRequest() *ReadAccelerationRequest {
	return &ReadAccelerationRequest{
		RTCGetParamRequest: *protocol_rtc.NewRTCGetParamRequest(
			uint16(protocol_common.Parameter.Acceleration1),
			protocol_rtc.CommandCode.ReadRAM,
		),
	}
}

// ReadDecelerationRequest reads deceleration (upid 0x145D).
type ReadDecelerationRequest struct {
	protocol_rtc.RTCGetParamRequest
}

// NewReadDecelerationRequest creates a new RTC read request for deceleration (upid 0x145D).
func NewReadDecelerationRequest() *ReadDecelerationRequest {
	return &ReadDecelerationRequest{
		RTCGetParamRequest: *protocol_rtc.NewRTCGetParamRequest(
			uint16(protocol_common.Parameter.Deceleration1),
			protocol_rtc.CommandCode.ReadRAM,
		),
	}
}

// WritePosition1Request writes Position 1 (upid 0x145A).
type WritePosition1Request struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWritePosition1Request creates a new RTC write request for Position 1 (upid 0x145A).
// positionMM is the position in millimeters.
func NewWritePosition1Request(positionMM float64, storageType protocol_common.ParameterStorageType) (*WritePosition1Request, error) {
	upid := protocol_common.Parameter.Position1
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	cmdCode, err := protocol_common.ToRTCCommandCode(storageType)
	if err != nil {
		return nil, err
	}
	value := protocol_common.ToProtocolPosition(positionMM)
	return &WritePosition1Request{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(
			uint16(upid),
			value,
			cmdCode,
		),
	}, nil
}

// WritePosition2Request writes Position 2 (upid 0x145F).
type WritePosition2Request struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWritePosition2Request creates a new RTC write request for Position 2 (upid 0x145F).
// positionMM is the position in millimeters.
func NewWritePosition2Request(positionMM float64, storageType protocol_common.ParameterStorageType) (*WritePosition2Request, error) {
	upid := protocol_common.Parameter.Position2
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	cmdCode, err := protocol_common.ToRTCCommandCode(storageType)
	if err != nil {
		return nil, err
	}
	value := protocol_common.ToProtocolPosition(positionMM)
	return &WritePosition2Request{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(
			uint16(upid),
			value,
			cmdCode,
		),
	}, nil
}

// WriteVelocityRequest writes velocity (upid 0x145B).
type WriteVelocityRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWriteVelocityRequest creates a new RTC write request for velocity (upid 0x145B).
// velocityMS is the velocity in meters per second.
func NewWriteVelocityRequest(velocityMS float64, storageType protocol_common.ParameterStorageType) (*WriteVelocityRequest, error) {
	upid := protocol_common.Parameter.Speed1
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	cmdCode, err := protocol_common.ToRTCCommandCode(storageType)
	if err != nil {
		return nil, err
	}
	value := protocol_common.ToProtocolVelocitySigned(velocityMS)
	return &WriteVelocityRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(
			uint16(upid),
			value,
			cmdCode,
		),
	}, nil
}

// WriteAccelerationRequest writes acceleration (upid 0x145C).
type WriteAccelerationRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWriteAccelerationRequest creates a new RTC write request for acceleration (upid 0x145C).
// accelMS2 is the acceleration in meters per second squared.
func NewWriteAccelerationRequest(accelMS2 float64, storageType protocol_common.ParameterStorageType) (*WriteAccelerationRequest, error) {
	upid := protocol_common.Parameter.Acceleration1
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	cmdCode, err := protocol_common.ToRTCCommandCode(storageType)
	if err != nil {
		return nil, err
	}
	value := protocol_common.ToProtocolAccelerationSigned(accelMS2)
	return &WriteAccelerationRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(
			uint16(upid),
			value,
			cmdCode,
		),
	}, nil
}

// WriteDecelerationRequest writes deceleration (upid 0x145D).
type WriteDecelerationRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWriteDecelerationRequest creates a new RTC write request for deceleration (upid 0x145D).
// decelMS2 is the deceleration in meters per second squared.
func NewWriteDecelerationRequest(decelMS2 float64, storageType protocol_common.ParameterStorageType) (*WriteDecelerationRequest, error) {
	upid := protocol_common.Parameter.Deceleration1
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	cmdCode, err := protocol_common.ToRTCCommandCode(storageType)
	if err != nil {
		return nil, err
	}
	value := protocol_common.ToProtocolAccelerationSigned(decelMS2)
	return &WriteDecelerationRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(
			uint16(upid),
			value,
			cmdCode,
		),
	}, nil
}

// WriteEasyStepsAutoStartRequest writes Easy Steps auto start configuration (upid 0x30D4).
type WriteEasyStepsAutoStartRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWriteEasyStepsAutoStartRequest creates a new RTC write request for Easy Steps auto start (upid 0x30D4).
// value should be from protocol_common.EasyStepsAutoStart (Disabled or Enabled).
func NewWriteEasyStepsAutoStartRequest(value int32, storageType protocol_common.ParameterStorageType) (*WriteEasyStepsAutoStartRequest, error) {
	upid := protocol_common.PUID.EasyStepsAutoStart
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	cmdCode, err := protocol_common.ToRTCCommandCode(storageType)
	if err != nil {
		return nil, err
	}
	return &WriteEasyStepsAutoStartRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(
			uint16(upid),
			value,
			cmdCode,
		),
	}, nil
}

// WriteEasyStepsAutoHomeRequest writes Easy Steps auto home configuration (upid 0x30D5).
type WriteEasyStepsAutoHomeRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWriteEasyStepsAutoHomeRequest creates a new RTC write request for Easy Steps auto home (upid 0x30D5).
// value should be from protocol_common.EasyStepsAutoHome (Disabled or Enabled).
func NewWriteEasyStepsAutoHomeRequest(value int32, storageType protocol_common.ParameterStorageType) (*WriteEasyStepsAutoHomeRequest, error) {
	upid := protocol_common.PUID.EasyStepsAutoHome
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	cmdCode, err := protocol_common.ToRTCCommandCode(storageType)
	if err != nil {
		return nil, err
	}
	return &WriteEasyStepsAutoHomeRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(
			uint16(upid),
			value,
			cmdCode,
		),
	}, nil
}

// WriteEasyStepsRisingEdgeRequest writes Easy Steps rising edge action configuration for any input pin.
type WriteEasyStepsRisingEdgeRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWriteEasyStepsRisingEdgeRequest creates a new RTC write request for Easy Steps rising edge action.
// inputNumber should be protocol_common.IOPin.Input45, .Input46, .Input47, or .Input48.
// value should be from protocol_common.EasyStepsIOMotion.
func NewWriteEasyStepsRisingEdgeRequest(inputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) (*WriteEasyStepsRisingEdgeRequest, error) {
	upid, err := protocol_common.GetEasyStepsRisingEdgePUID(inputNumber)
	if err != nil {
		return nil, err
	}
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	cmdCode, err := protocol_common.ToRTCCommandCode(storageType)
	if err != nil {
		return nil, err
	}
	return &WriteEasyStepsRisingEdgeRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(
			uint16(upid),
			value,
			cmdCode,
		),
	}, nil
}

// WriteEasyStepsIOMotionConfigCmdRequest writes Easy Steps IO motion config curve/CMD ID for any input pin.
type WriteEasyStepsIOMotionConfigCmdRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWriteEasyStepsIOMotionConfigCmdRequest creates a new RTC write request for Easy Steps IO motion config.
// inputNumber should be protocol_common.IOPin.Input45, .Input46, .Input47, or .Input48.
// curveCmdID is the curve or command table ID to execute.
func NewWriteEasyStepsIOMotionConfigCmdRequest(inputNumber protocol_common.IOPinNumber, curveCmdID int32, storageType protocol_common.ParameterStorageType) (*WriteEasyStepsIOMotionConfigCmdRequest, error) {
	upid, err := protocol_common.GetEasyStepsIOMotionConfigCmdPUID(inputNumber)
	if err != nil {
		return nil, err
	}
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	cmdCode, err := protocol_common.ToRTCCommandCode(storageType)
	if err != nil {
		return nil, err
	}
	return &WriteEasyStepsIOMotionConfigCmdRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(
			uint16(upid),
			curveCmdID,
			cmdCode,
		),
	}, nil
}

// WriteOutputFunctionRequest writes output pin function configuration for any output pin.
type WriteOutputFunctionRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWriteOutputFunctionRequest creates a new RTC write request for output function configuration.
// outputNumber should be protocol_common.IOPin.Output36, .Output43, or .Output44.
// value should be from protocol_common.OutputConfig.
func NewWriteOutputFunctionRequest(outputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) (*WriteOutputFunctionRequest, error) {
	upid, err := protocol_common.GetOutputFunctionPUID(outputNumber)
	if err != nil {
		return nil, err
	}
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	cmdCode, err := protocol_common.ToRTCCommandCode(storageType)
	if err != nil {
		return nil, err
	}
	return &WriteOutputFunctionRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(
			uint16(upid),
			value,
			cmdCode,
		),
	}, nil
}

// WriteInputFunctionRequest writes input pin function configuration for any input pin.
type WriteInputFunctionRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWriteInputFunctionRequest creates a new RTC write request for input function configuration.
// inputNumber should be protocol_common.IOPin.Input45, .Input46, .Input47, or .Input48.
// value should be from protocol_common.InputFunction.
func NewWriteInputFunctionRequest(inputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) (*WriteInputFunctionRequest, error) {
	upid, err := protocol_common.GetInputFunctionPUID(inputNumber)
	if err != nil {
		return nil, err
	}
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	cmdCode, err := protocol_common.ToRTCCommandCode(storageType)
	if err != nil {
		return nil, err
	}
	return &WriteInputFunctionRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(
			uint16(upid),
			value,
			cmdCode,
		),
	}, nil
}

// WriteRunModeRequest writes run mode configuration (upid 0x1450).
type WriteRunModeRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWriteRunModeRequest creates a new RTC write request for run mode (upid 0x1450).
// mode should be from protocol_common.RunModes.
func NewWriteRunModeRequest(mode protocol_common.RunMode, storageType protocol_common.ParameterStorageType) (*WriteRunModeRequest, error) {
	upid := protocol_common.PUID.RunMode
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	cmdCode, err := protocol_common.ToRTCCommandCode(storageType)
	if err != nil {
		return nil, err
	}
	return &WriteRunModeRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(
			uint16(upid),
			int32(mode),
			cmdCode,
		),
	}, nil
}

// WriteTriggerModeRequest writes trigger mode configuration (upid 0x170C).
type WriteTriggerModeRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewWriteTriggerModeRequest creates a new RTC write request for trigger mode (upid 0x170C).
// value should be from protocol_common.TriggerModeConfig.
func NewWriteTriggerModeRequest(value int32, storageType protocol_common.ParameterStorageType) (*WriteTriggerModeRequest, error) {
	upid := protocol_common.PUID.TriggerMode
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	cmdCode, err := protocol_common.ToRTCCommandCode(storageType)
	if err != nil {
		return nil, err
	}
	return &WriteTriggerModeRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(
			uint16(upid),
			value,
			cmdCode,
		),
	}, nil
}
