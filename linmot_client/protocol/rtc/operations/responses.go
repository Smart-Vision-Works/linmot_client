package protocol_operations

import (
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/staged_robot/protocol/rtc"
)

// Compile-time interface checks
var (
	_ protocol_common.Response = (*RestartDriveResponse)(nil)
	_ protocol_common.Response = (*SetOSROMToDefaultResponse)(nil)
	_ protocol_common.Response = (*SetMCROMToDefaultResponse)(nil)
	_ protocol_common.Response = (*SetInterfaceROMToDefaultResponse)(nil)
	_ protocol_common.Response = (*SetApplicationROMToDefaultResponse)(nil)
)

func init() {
	registerOperationResponseRegistries()
}

// RestartDriveResponse is the response to RestartDriveRequest.
type RestartDriveResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewRestartDriveResponse wraps an RTCSetParamResponse as a RestartDriveResponse.
func NewRestartDriveResponse(base *protocol_rtc.RTCSetParamResponse) *RestartDriveResponse {
	return &RestartDriveResponse{RTCSetParamResponse: *base}
}

// SetOSROMToDefaultResponse is the response to SetOSROMToDefaultRequest.
type SetOSROMToDefaultResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewSetOSROMToDefaultResponse wraps an RTCSetParamResponse as a SetOSROMToDefaultResponse.
func NewSetOSROMToDefaultResponse(base *protocol_rtc.RTCSetParamResponse) *SetOSROMToDefaultResponse {
	return &SetOSROMToDefaultResponse{RTCSetParamResponse: *base}
}

// SetMCROMToDefaultResponse is the response to SetMCROMToDefaultRequest.
type SetMCROMToDefaultResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewSetMCROMToDefaultResponse wraps an RTCSetParamResponse as a SetMCROMToDefaultResponse.
func NewSetMCROMToDefaultResponse(base *protocol_rtc.RTCSetParamResponse) *SetMCROMToDefaultResponse {
	return &SetMCROMToDefaultResponse{RTCSetParamResponse: *base}
}

// SetInterfaceROMToDefaultResponse is the response to SetInterfaceROMToDefaultRequest.
type SetInterfaceROMToDefaultResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewSetInterfaceROMToDefaultResponse wraps an RTCSetParamResponse as a SetInterfaceROMToDefaultResponse.
func NewSetInterfaceROMToDefaultResponse(base *protocol_rtc.RTCSetParamResponse) *SetInterfaceROMToDefaultResponse {
	return &SetInterfaceROMToDefaultResponse{RTCSetParamResponse: *base}
}

// SetApplicationROMToDefaultResponse is the response to SetApplicationROMToDefaultRequest.
type SetApplicationROMToDefaultResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewSetApplicationROMToDefaultResponse wraps an RTCSetParamResponse as a SetApplicationROMToDefaultResponse.
func NewSetApplicationROMToDefaultResponse(base *protocol_rtc.RTCSetParamResponse) *SetApplicationROMToDefaultResponse {
	return &SetApplicationROMToDefaultResponse{RTCSetParamResponse: *base}
}

func registerOperationResponseRegistries() {
	registerOperationSetRegistry(protocol_rtc.CommandCode.RestartDrive, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewRestartDriveResponse(base)
	})
	registerOperationSetRegistry(protocol_rtc.CommandCode.SetOSROMToDefault, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewSetOSROMToDefaultResponse(base)
	})
	registerOperationSetRegistry(protocol_rtc.CommandCode.SetMCROMToDefault, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewSetMCROMToDefaultResponse(base)
	})
	registerOperationSetRegistry(protocol_rtc.CommandCode.SetInterfaceROMToDefault, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewSetInterfaceROMToDefaultResponse(base)
	})
	registerOperationSetRegistry(protocol_rtc.CommandCode.SetApplicationROMToDefault, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewSetApplicationROMToDefaultResponse(base)
	})
}

func registerOperationSetRegistry(cmdCode uint8, wrapper func(*protocol_rtc.RTCSetParamResponse) protocol_common.Response) {
	protocol_rtc.RegisterResponseRegistryByCmd(cmdCode,
		func(status *protocol_common.Status, value int32, upid uint16, rtcCounter, rtcStatus, cmdCode uint8) protocol_common.Response {
			base := protocol_rtc.NewRTCSetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
			return wrapper(base)
		})
}
