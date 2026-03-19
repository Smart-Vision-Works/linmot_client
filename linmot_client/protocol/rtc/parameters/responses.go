package protocol_parameters

import (
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/staged_robot/protocol/rtc"
)

// Compile-time interface checks
var (
	_ protocol_common.Response = (*WriteRAMAndROMResponse)(nil)
	_ protocol_common.Response = (*GetMinValueResponse)(nil)
	_ protocol_common.Response = (*GetMaxValueResponse)(nil)
	_ protocol_common.Response = (*GetDefaultValueResponse)(nil)
	_ protocol_common.Response = (*StartGettingUPIDListResponse)(nil)
	_ protocol_common.Response = (*GetNextUPIDListItemResponse)(nil)
	_ protocol_common.Response = (*StartGettingModifiedUPIDListResponse)(nil)
	_ protocol_common.Response = (*GetNextModifiedUPIDListItemResponse)(nil)
	// Motion parameter responses (RTC-based wrappers)
	_ protocol_common.Response = (*ReadVelocityResponse)(nil)
	_ protocol_common.Response = (*ReadPosition1Response)(nil)
	_ protocol_common.Response = (*ReadAccelerationResponse)(nil)
	_ protocol_common.Response = (*ReadDecelerationResponse)(nil)
	_ protocol_common.Response = (*WritePosition1Response)(nil)
	_ protocol_common.Response = (*WritePosition2Response)(nil)
	_ protocol_common.Response = (*WriteVelocityResponse)(nil)
	_ protocol_common.Response = (*WriteAccelerationResponse)(nil)
	_ protocol_common.Response = (*WriteDecelerationResponse)(nil)
	_ protocol_common.Response = (*WriteEasyStepsAutoStartResponse)(nil)
	_ protocol_common.Response = (*WriteEasyStepsAutoHomeResponse)(nil)
	_ protocol_common.Response = (*WriteEasyStepsRisingEdgeResponse)(nil)
	_ protocol_common.Response = (*WriteEasyStepsIOMotionConfigCmdResponse)(nil)
	_ protocol_common.Response = (*WriteOutputFunctionResponse)(nil)
	_ protocol_common.Response = (*WriteInputFunctionResponse)(nil)
	_ protocol_common.Response = (*WriteTriggerModeResponse)(nil)
	_ protocol_common.Response = (*WriteRunModeResponse)(nil)
)

func init() {
	// Register factories for extended parameter operations using cmdCode-based routing
	// These operations work on ANY UPID, so we can't use UPID-based factories
	// Instead, the parser needs to check cmdCode first

	// Register UPID list operation response factories
	registerUPIDListResponseRegistries()

	// Register motion parameter response factories
	registerMotionParameterFactories()
}

// registerUPIDListResponseRegistries registers response factories for UPID list operations.
func registerUPIDListResponseRegistries() {
	registerUPIDListSetRegistry(protocol_rtc.CommandCode.StartGettingUPIDList, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewStartGettingUPIDListResponse(base)
	})
	registerUPIDListSetRegistry(protocol_rtc.CommandCode.GetNextUPIDListItem, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewGetNextUPIDListItemResponse(base)
	})
	registerUPIDListSetRegistry(protocol_rtc.CommandCode.StartGettingModifiedUPIDList, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewStartGettingModifiedUPIDListResponse(base)
	})
	registerUPIDListSetRegistry(protocol_rtc.CommandCode.GetNextModifiedUPIDListItem, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewGetNextModifiedUPIDListItemResponse(base)
	})
}

func registerUPIDListSetRegistry(cmdCode uint8, wrapper func(*protocol_rtc.RTCSetParamResponse) protocol_common.Response) {
	protocol_rtc.RegisterResponseRegistryByCmd(cmdCode,
		func(status *protocol_common.Status, value int32, upid uint16, rtcCounter, rtcStatus, cmdCode uint8) protocol_common.Response {
			base := protocol_rtc.NewRTCSetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
			return wrapper(base)
		})
}

// WriteRAMAndROMResponse is the response to WriteRAMAndROMRequest.
type WriteRAMAndROMResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWriteRAMAndROMResponse wraps an RTCSetParamResponse as a WriteRAMAndROMResponse.
func NewWriteRAMAndROMResponse(base *protocol_rtc.RTCSetParamResponse) *WriteRAMAndROMResponse {
	return &WriteRAMAndROMResponse{RTCSetParamResponse: *base}
}

// GetMinValueResponse is the response to GetMinValueRequest.
type GetMinValueResponse struct {
	protocol_rtc.RTCGetParamResponse
}

// NewGetMinValueResponse wraps an RTCGetParamResponse as a GetMinValueResponse.
func NewGetMinValueResponse(base *protocol_rtc.RTCGetParamResponse) *GetMinValueResponse {
	return &GetMinValueResponse{RTCGetParamResponse: *base}
}

// MinValue returns the minimum value for the parameter.
func (r *GetMinValueResponse) MinValue() int32 {
	return r.Value()
}

// GetMaxValueResponse is the response to GetMaxValueRequest.
type GetMaxValueResponse struct {
	protocol_rtc.RTCGetParamResponse
}

// NewGetMaxValueResponse wraps an RTCGetParamResponse as a GetMaxValueResponse.
func NewGetMaxValueResponse(base *protocol_rtc.RTCGetParamResponse) *GetMaxValueResponse {
	return &GetMaxValueResponse{RTCGetParamResponse: *base}
}

// MaxValue returns the maximum value for the parameter.
func (r *GetMaxValueResponse) MaxValue() int32 {
	return r.Value()
}

// GetDefaultValueResponse is the response to GetDefaultValueRequest.
type GetDefaultValueResponse struct {
	protocol_rtc.RTCGetParamResponse
}

// NewGetDefaultValueResponse wraps an RTCGetParamResponse as a GetDefaultValueResponse.
func NewGetDefaultValueResponse(base *protocol_rtc.RTCGetParamResponse) *GetDefaultValueResponse {
	return &GetDefaultValueResponse{RTCGetParamResponse: *base}
}

// DefaultValue returns the default value for the parameter.
func (r *GetDefaultValueResponse) DefaultValue() int32 {
	return r.Value()
}

// StartGettingUPIDListResponse is the response to StartGettingUPIDListRequest.
type StartGettingUPIDListResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewStartGettingUPIDListResponse wraps an RTCSetParamResponse as a StartGettingUPIDListResponse.
func NewStartGettingUPIDListResponse(base *protocol_rtc.RTCSetParamResponse) *StartGettingUPIDListResponse {
	return &StartGettingUPIDListResponse{RTCSetParamResponse: *base}
}

// GetNextUPIDListItemResponse is the response to GetNextUPIDListItemRequest.
type GetNextUPIDListItemResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewGetNextUPIDListItemResponse wraps an RTCSetParamResponse as a GetNextUPIDListItemResponse.
func NewGetNextUPIDListItemResponse(base *protocol_rtc.RTCSetParamResponse) *GetNextUPIDListItemResponse {
	return &GetNextUPIDListItemResponse{RTCSetParamResponse: *base}
}

// UPID returns the UPID found in the list (from response UPID field).
func (r *GetNextUPIDListItemResponse) FoundUPID() uint16 {
	return r.UPID()
}

// AddressUsage returns the address usage from the response value (bits 0-15).
func (r *GetNextUPIDListItemResponse) AddressUsage() uint16 {
	return uint16(r.Value() & 0xFFFF)
}

// StartGettingModifiedUPIDListResponse is the response to StartGettingModifiedUPIDListRequest.
type StartGettingModifiedUPIDListResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewStartGettingModifiedUPIDListResponse wraps an RTCSetParamResponse as a StartGettingModifiedUPIDListResponse.
func NewStartGettingModifiedUPIDListResponse(base *protocol_rtc.RTCSetParamResponse) *StartGettingModifiedUPIDListResponse {
	return &StartGettingModifiedUPIDListResponse{RTCSetParamResponse: *base}
}

// GetNextModifiedUPIDListItemResponse is the response to GetNextModifiedUPIDListItemRequest.
type GetNextModifiedUPIDListItemResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewGetNextModifiedUPIDListItemResponse wraps an RTCSetParamResponse as a GetNextModifiedUPIDListItemResponse.
func NewGetNextModifiedUPIDListItemResponse(base *protocol_rtc.RTCSetParamResponse) *GetNextModifiedUPIDListItemResponse {
	return &GetNextModifiedUPIDListItemResponse{RTCSetParamResponse: *base}
}

// FoundUPID returns the UPID found in the modified list (from response UPID field).
func (r *GetNextModifiedUPIDListItemResponse) FoundUPID() uint16 {
	return r.UPID()
}

// ModifiedValue returns the modified parameter value.
func (r *GetNextModifiedUPIDListItemResponse) ModifiedValue() int32 {
	return r.Value()
}

// ============================================================================
// Motion Parameter Responses (RTC-based wrappers for motion-related UPIDs)
// ============================================================================

// ReadVelocityResponse is the response to a ReadVelocityRequest.
type ReadVelocityResponse struct {
	protocol_rtc.RTCGetParamResponse
}

// NewReadVelocityResponse wraps an RTCGetParamResponse as a ReadVelocityResponse.
func NewReadVelocityResponse(base *protocol_rtc.RTCGetParamResponse) *ReadVelocityResponse {
	return &ReadVelocityResponse{RTCGetParamResponse: *base}
}

// VelocityMS returns the velocity in meters per second.
func (r *ReadVelocityResponse) VelocityMS() float64 {
	return r.RTCGetParamResponse.VelocityMS()
}

// ReadPosition1Response is the response to a ReadPosition1Request.
type ReadPosition1Response struct {
	protocol_rtc.RTCGetParamResponse
}

// NewReadPosition1Response wraps an RTCGetParamResponse as a ReadPosition1Response.
func NewReadPosition1Response(base *protocol_rtc.RTCGetParamResponse) *ReadPosition1Response {
	return &ReadPosition1Response{RTCGetParamResponse: *base}
}

// PositionMM returns the position in millimeters.
func (r *ReadPosition1Response) PositionMM() float64 {
	return protocol_common.FromProtocolPosition(r.Value())
}

// ReadAccelerationResponse is the response to a ReadAccelerationRequest.
type ReadAccelerationResponse struct {
	protocol_rtc.RTCGetParamResponse
}

// NewReadAccelerationResponse wraps an RTCGetParamResponse as a ReadAccelerationResponse.
func NewReadAccelerationResponse(base *protocol_rtc.RTCGetParamResponse) *ReadAccelerationResponse {
	return &ReadAccelerationResponse{RTCGetParamResponse: *base}
}

// AccelerationMS2 returns the acceleration in meters per second squared.
func (r *ReadAccelerationResponse) AccelerationMS2() float64 {
	return r.RTCGetParamResponse.AccelerationMS2()
}

// ReadDecelerationResponse is the response to a ReadDecelerationRequest.
type ReadDecelerationResponse struct {
	protocol_rtc.RTCGetParamResponse
}

// NewReadDecelerationResponse wraps an RTCGetParamResponse as a ReadDecelerationResponse.
func NewReadDecelerationResponse(base *protocol_rtc.RTCGetParamResponse) *ReadDecelerationResponse {
	return &ReadDecelerationResponse{RTCGetParamResponse: *base}
}

// DecelerationMS2 returns the deceleration in meters per second squared.
func (r *ReadDecelerationResponse) DecelerationMS2() float64 {
	return r.RTCGetParamResponse.AccelerationMS2()
}

// WritePosition1Response is the response to a WritePosition1Request.
type WritePosition1Response struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWritePosition1Response wraps an RTCSetParamResponse as a WritePosition1Response.
func NewWritePosition1Response(base *protocol_rtc.RTCSetParamResponse) *WritePosition1Response {
	return &WritePosition1Response{RTCSetParamResponse: *base}
}

// WritePosition2Response is the response to a WritePosition2Request.
type WritePosition2Response struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWritePosition2Response wraps an RTCSetParamResponse as a WritePosition2Response.
func NewWritePosition2Response(base *protocol_rtc.RTCSetParamResponse) *WritePosition2Response {
	return &WritePosition2Response{RTCSetParamResponse: *base}
}

// WriteVelocityResponse is the response to a WriteVelocityRequest.
type WriteVelocityResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWriteVelocityResponse wraps an RTCSetParamResponse as a WriteVelocityResponse.
func NewWriteVelocityResponse(base *protocol_rtc.RTCSetParamResponse) *WriteVelocityResponse {
	return &WriteVelocityResponse{RTCSetParamResponse: *base}
}

// WriteAccelerationResponse is the response to a WriteAccelerationRequest.
type WriteAccelerationResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWriteAccelerationResponse wraps an RTCSetParamResponse as a WriteAccelerationResponse.
func NewWriteAccelerationResponse(base *protocol_rtc.RTCSetParamResponse) *WriteAccelerationResponse {
	return &WriteAccelerationResponse{RTCSetParamResponse: *base}
}

// WriteDecelerationResponse is the response to a WriteDecelerationRequest.
type WriteDecelerationResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWriteDecelerationResponse wraps an RTCSetParamResponse as a WriteDecelerationResponse.
func NewWriteDecelerationResponse(base *protocol_rtc.RTCSetParamResponse) *WriteDecelerationResponse {
	return &WriteDecelerationResponse{RTCSetParamResponse: *base}
}

// WriteEasyStepsAutoStartResponse is the response to a WriteEasyStepsAutoStartRequest.
type WriteEasyStepsAutoStartResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWriteEasyStepsAutoStartResponse wraps an RTCSetParamResponse as a WriteEasyStepsAutoStartResponse.
func NewWriteEasyStepsAutoStartResponse(base *protocol_rtc.RTCSetParamResponse) *WriteEasyStepsAutoStartResponse {
	return &WriteEasyStepsAutoStartResponse{RTCSetParamResponse: *base}
}

// WriteEasyStepsAutoHomeResponse is the response to a WriteEasyStepsAutoHomeRequest.
type WriteEasyStepsAutoHomeResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWriteEasyStepsAutoHomeResponse wraps an RTCSetParamResponse as a WriteEasyStepsAutoHomeResponse.
func NewWriteEasyStepsAutoHomeResponse(base *protocol_rtc.RTCSetParamResponse) *WriteEasyStepsAutoHomeResponse {
	return &WriteEasyStepsAutoHomeResponse{RTCSetParamResponse: *base}
}

// WriteEasyStepsRisingEdgeResponse is the response to a WriteEasyStepsRisingEdgeRequest.
type WriteEasyStepsRisingEdgeResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWriteEasyStepsRisingEdgeResponse wraps an RTCSetParamResponse as a WriteEasyStepsRisingEdgeResponse.
func NewWriteEasyStepsRisingEdgeResponse(base *protocol_rtc.RTCSetParamResponse) *WriteEasyStepsRisingEdgeResponse {
	return &WriteEasyStepsRisingEdgeResponse{RTCSetParamResponse: *base}
}

// WriteEasyStepsIOMotionConfigCmdResponse is the response to a WriteEasyStepsIOMotionConfigCmdRequest.
type WriteEasyStepsIOMotionConfigCmdResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWriteEasyStepsIOMotionConfigCmdResponse wraps an RTCSetParamResponse as a WriteEasyStepsIOMotionConfigCmdResponse.
func NewWriteEasyStepsIOMotionConfigCmdResponse(base *protocol_rtc.RTCSetParamResponse) *WriteEasyStepsIOMotionConfigCmdResponse {
	return &WriteEasyStepsIOMotionConfigCmdResponse{RTCSetParamResponse: *base}
}

// WriteOutputFunctionResponse is the response to a WriteOutputFunctionRequest.
type WriteOutputFunctionResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWriteOutputFunctionResponse wraps an RTCSetParamResponse as a WriteOutputFunctionResponse.
func NewWriteOutputFunctionResponse(base *protocol_rtc.RTCSetParamResponse) *WriteOutputFunctionResponse {
	return &WriteOutputFunctionResponse{RTCSetParamResponse: *base}
}

// WriteInputFunctionResponse is the response to a WriteInputFunctionRequest.
type WriteInputFunctionResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWriteInputFunctionResponse wraps an RTCSetParamResponse as a WriteInputFunctionResponse.
func NewWriteInputFunctionResponse(base *protocol_rtc.RTCSetParamResponse) *WriteInputFunctionResponse {
	return &WriteInputFunctionResponse{RTCSetParamResponse: *base}
}

// WriteTriggerModeResponse is the response to a WriteTriggerModeRequest.
type WriteTriggerModeResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWriteTriggerModeResponse wraps an RTCSetParamResponse as a WriteTriggerModeResponse.
func NewWriteTriggerModeResponse(base *protocol_rtc.RTCSetParamResponse) *WriteTriggerModeResponse {
	return &WriteTriggerModeResponse{RTCSetParamResponse: *base}
}

// WriteRunModeResponse is the response to a WriteRunModeRequest.
type WriteRunModeResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWriteRunModeResponse wraps an RTCSetParamResponse as a WriteRunModeResponse.
func NewWriteRunModeResponse(base *protocol_rtc.RTCSetParamResponse) *WriteRunModeResponse {
	return &WriteRunModeResponse{RTCSetParamResponse: *base}
}

// ============================================================================
// Response Factory Registration
// ============================================================================

// registerReadResponseFactory is a helper to reduce boilerplate when registering read response factories.
func registerReadResponseFactory(upid protocol_common.ParameterID, wrapper func(*protocol_rtc.RTCGetParamResponse) protocol_common.Response) {

	for _, cmdCode := range []uint8{
		protocol_rtc.CommandCode.ReadRAM,
		protocol_rtc.CommandCode.ReadROM,
	} {
		protocol_rtc.RegisterResponseRegistryByCmdAndUPID(cmdCode, upid,
			func(status *protocol_common.Status, value int32, respUpid uint16, rtcCounter, rtcStatus, cmdCode uint8) protocol_common.Response {
				base := protocol_rtc.NewRTCGetParamResponseWithCmdCode(status, value, respUpid, rtcCounter, rtcStatus, cmdCode)
				return wrapper(base)
			})
	}
}

// registerWriteResponseFactory is a helper to reduce boilerplate when registering write response factories.
// Registers for all parameter storage command codes: WriteRAM, WriteROM, and WriteRAMAndROM.
func registerWriteResponseFactory(upid protocol_common.ParameterID, wrapper func(*protocol_rtc.RTCSetParamResponse) protocol_common.Response) {
	// All command codes that ToRTCCommandCode can return for parameter writes
	writeCommandCodes := []uint8{
		protocol_rtc.CommandCode.WriteROM,       // 0x12 - ParameterStorage.ROM
		protocol_rtc.CommandCode.WriteRAM,       // 0x13 - ParameterStorage.RAM
		protocol_rtc.CommandCode.WriteRAMAndROM, // 0x14 - ParameterStorage.Both
	}

	for _, cmdCode := range writeCommandCodes {
		protocol_rtc.RegisterResponseRegistryByCmdAndUPID(cmdCode, upid,
			func(status *protocol_common.Status, value int32, respUpid uint16, rtcCounter, rtcStatus, cmdCode uint8) protocol_common.Response {
				base := protocol_rtc.NewRTCSetParamResponseWithCmdCode(status, value, respUpid, rtcCounter, rtcStatus, cmdCode)
				return wrapper(base)
			})
	}
}

func registerMotionParameterFactories() {
	// Register typed read response factories
	registerReadResponseFactory(protocol_common.Parameter.Speed1, func(base *protocol_rtc.RTCGetParamResponse) protocol_common.Response {
		return NewReadVelocityResponse(base)
	})

	// Register Monitoring Channel configuration read responses (0x20A8-0x20AB)
	// These are standard parameter reads that read which UPID each monitoring channel monitors
	// Use cmd+UPID registration to avoid conflicts with write registrations
	for _, upid := range []protocol_common.ParameterID{
		protocol_common.PUID.MonitoringChannel1UPID,
		protocol_common.PUID.MonitoringChannel2UPID,
		protocol_common.PUID.MonitoringChannel3UPID,
		protocol_common.PUID.MonitoringChannel4UPID,
	} {
		protocol_rtc.RegisterResponseRegistryByCmdAndUPID(protocol_rtc.CommandCode.ReadRAM, upid,
			func(status *protocol_common.Status, value int32, respUpid uint16, rtcCounter, rtcStatus, cmdCode uint8) protocol_common.Response {
				return protocol_rtc.NewRTCGetParamResponse(status, value, respUpid, rtcCounter, rtcStatus)
			})
	}

	registerReadResponseFactory(protocol_common.Parameter.Position1, func(base *protocol_rtc.RTCGetParamResponse) protocol_common.Response {
		return NewReadPosition1Response(base)
	})

	registerReadResponseFactory(protocol_common.Parameter.Acceleration1, func(base *protocol_rtc.RTCGetParamResponse) protocol_common.Response {
		return NewReadAccelerationResponse(base)
	})

	registerReadResponseFactory(protocol_common.Parameter.Deceleration1, func(base *protocol_rtc.RTCGetParamResponse) protocol_common.Response {
		return NewReadDecelerationResponse(base)
	})

	// Register typed write response factories
	registerWriteResponseFactory(protocol_common.Parameter.Position1, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewWritePosition1Response(base)
	})

	registerWriteResponseFactory(protocol_common.Parameter.Position2, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewWritePosition2Response(base)
	})

	registerWriteResponseFactory(protocol_common.Parameter.Speed1, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewWriteVelocityResponse(base)
	})

	registerWriteResponseFactory(protocol_common.Parameter.Acceleration1, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewWriteAccelerationResponse(base)
	})

	registerWriteResponseFactory(protocol_common.Parameter.Deceleration1, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewWriteDecelerationResponse(base)
	})

	// Register EasySteps configuration responses
	registerWriteResponseFactory(protocol_common.PUID.EasyStepsAutoStart, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewWriteEasyStepsAutoStartResponse(base)
	})

	registerWriteResponseFactory(protocol_common.PUID.EasyStepsAutoHome, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewWriteEasyStepsAutoHomeResponse(base)
	})

	// Register EasySteps I/O responses (4 UPIDs per type)
	for _, inputNum := range []protocol_common.IOPinNumber{
		protocol_common.IOPin.Input45,
		protocol_common.IOPin.Input46,
		protocol_common.IOPin.Input47,
		protocol_common.IOPin.Input48,
	} {
		// Rising Edge responses
		if upid, err := protocol_common.GetEasyStepsRisingEdgePUID(inputNum); err == nil {
			registerWriteResponseFactory(upid, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
				return NewWriteEasyStepsRisingEdgeResponse(base)
			})
		}

		// IO Motion Config Cmd responses
		if upid, err := protocol_common.GetEasyStepsIOMotionConfigCmdPUID(inputNum); err == nil {
			registerWriteResponseFactory(upid, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
				return NewWriteEasyStepsIOMotionConfigCmdResponse(base)
			})
		}

		// Input Function responses
		if upid, err := protocol_common.GetInputFunctionPUID(inputNum); err == nil {
			registerWriteResponseFactory(upid, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
				return NewWriteInputFunctionResponse(base)
			})
		}
	}

	// Register Output Function responses (3 UPIDs)
	for _, outputNum := range []protocol_common.IOPinNumber{
		protocol_common.IOPin.Output36,
		protocol_common.IOPin.Output43,
		protocol_common.IOPin.Output44,
	} {
		if upid, err := protocol_common.GetOutputFunctionPUID(outputNum); err == nil {
			registerWriteResponseFactory(upid, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
				return NewWriteOutputFunctionResponse(base)
			})
		}
	}

	// Register Trigger Mode response
	registerWriteResponseFactory(protocol_common.PUID.TriggerMode, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewWriteTriggerModeResponse(base)
	})

	// Register Run Mode response
	registerWriteResponseFactory(protocol_common.PUID.RunMode, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewWriteRunModeResponse(base)
	})

	// Register Monitoring Channel configuration write responses (0x20A8-0x20AB)
	// These are standard parameter writes that configure which UPID each monitoring channel monitors
	// Use cmd+UPID registration to avoid conflicts with read registrations
	for _, upid := range []protocol_common.ParameterID{
		protocol_common.PUID.MonitoringChannel1UPID,
		protocol_common.PUID.MonitoringChannel2UPID,
		protocol_common.PUID.MonitoringChannel3UPID,
		protocol_common.PUID.MonitoringChannel4UPID,
	} {
		protocol_rtc.RegisterResponseRegistryByCmdAndUPID(protocol_rtc.CommandCode.WriteRAM, upid,
			func(status *protocol_common.Status, value int32, respUpid uint16, rtcCounter, rtcStatus, cmdCode uint8) protocol_common.Response {
				return protocol_rtc.NewRTCSetParamResponse(status, value, respUpid, rtcCounter, rtcStatus)
			})
	}
}
