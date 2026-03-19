package protocol_operations

import (
	protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/linmot_client/protocol/rtc"
)

// Compile-time interface checks
var (
	_ protocol_common.Request = (*RestartDriveRequest)(nil)
	_ protocol_common.Request = (*SetOSROMToDefaultRequest)(nil)
	_ protocol_common.Request = (*SetMCROMToDefaultRequest)(nil)
	_ protocol_common.Request = (*SetInterfaceROMToDefaultRequest)(nil)
	_ protocol_common.Request = (*SetApplicationROMToDefaultRequest)(nil)
)

// RestartDriveRequest requests a drive restart.
type RestartDriveRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewRestartDriveRequest creates a new request to restart the drive.
func NewRestartDriveRequest() *RestartDriveRequest {
	return &RestartDriveRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(0, 0, protocol_rtc.CommandCode.RestartDrive),
	}
}

// SetOSROMToDefaultRequest sets OS SW parameter ROM values to default.
type SetOSROMToDefaultRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewSetOSROMToDefaultRequest creates a new request to reset OS SW parameters to default.
func NewSetOSROMToDefaultRequest() *SetOSROMToDefaultRequest {
	return &SetOSROMToDefaultRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(0, 0, protocol_rtc.CommandCode.SetOSROMToDefault),
	}
}

// SetMCROMToDefaultRequest sets MC SW parameter ROM values to default.
type SetMCROMToDefaultRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewSetMCROMToDefaultRequest creates a new request to reset MC SW parameters to default.
func NewSetMCROMToDefaultRequest() *SetMCROMToDefaultRequest {
	return &SetMCROMToDefaultRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(0, 0, protocol_rtc.CommandCode.SetMCROMToDefault),
	}
}

// SetInterfaceROMToDefaultRequest sets Interface SW parameter ROM values to default.
type SetInterfaceROMToDefaultRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewSetInterfaceROMToDefaultRequest creates a new request to reset Interface SW parameters to default.
func NewSetInterfaceROMToDefaultRequest() *SetInterfaceROMToDefaultRequest {
	return &SetInterfaceROMToDefaultRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(0, 0, protocol_rtc.CommandCode.SetInterfaceROMToDefault),
	}
}

// SetApplicationROMToDefaultRequest sets Application SW parameter ROM values to default.
type SetApplicationROMToDefaultRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewSetApplicationROMToDefaultRequest creates a new request to reset Application SW parameters to default.
func NewSetApplicationROMToDefaultRequest() *SetApplicationROMToDefaultRequest {
	return &SetApplicationROMToDefaultRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(0, 0, protocol_rtc.CommandCode.SetApplicationROMToDefault),
	}
}
