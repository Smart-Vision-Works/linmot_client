package protocol_command_tables

import (
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/staged_robot/protocol/rtc"
)

// Compile-time interface checks
var (
	_ protocol_common.Response = (*SaveCommandTableResponse)(nil)
	_ protocol_common.Response = (*StopMotionControllerResponse)(nil)
	_ protocol_common.Response = (*StartMotionControllerResponse)(nil)
	_ protocol_common.Response = (*DeleteAllEntriesResponse)(nil)
	_ protocol_common.Response = (*DeleteEntryResponse)(nil)
	_ protocol_common.Response = (*AllocateEntryResponse)(nil)
	_ protocol_common.Response = (*WriteEntryDataResponse)(nil)
	_ protocol_common.Response = (*GetEntrySizeResponse)(nil)
	_ protocol_common.Response = (*ReadEntryDataResponse)(nil)
	_ protocol_common.Response = (*PresenceMaskResponse)(nil)
)

func init() {
	registerCommandTableResponseRegistries()
}

// SaveCommandTableResponse is the response to SaveCommandTableRequest.
type SaveCommandTableResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewSaveCommandTableResponse wraps an RTCSetParamResponse as a SaveCommandTableResponse.
func NewSaveCommandTableResponse(base *protocol_rtc.RTCSetParamResponse) *SaveCommandTableResponse {
	return &SaveCommandTableResponse{RTCSetParamResponse: *base}
}

// StopMotionControllerResponse is the response to StopMotionControllerRequest.
type StopMotionControllerResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewStopMotionControllerResponse wraps an RTCSetParamResponse as a StopMotionControllerResponse.
func NewStopMotionControllerResponse(base *protocol_rtc.RTCSetParamResponse) *StopMotionControllerResponse {
	return &StopMotionControllerResponse{RTCSetParamResponse: *base}
}

// StartMotionControllerResponse is the response to StartMotionControllerRequest.
type StartMotionControllerResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewStartMotionControllerResponse wraps an RTCSetParamResponse as a StartMotionControllerResponse.
func NewStartMotionControllerResponse(base *protocol_rtc.RTCSetParamResponse) *StartMotionControllerResponse {
	return &StartMotionControllerResponse{RTCSetParamResponse: *base}
}

// DeleteAllEntriesResponse is the response to DeleteAllEntriesRequest.
type DeleteAllEntriesResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewDeleteAllEntriesResponse wraps an RTCSetParamResponse as a DeleteAllEntriesResponse.
func NewDeleteAllEntriesResponse(base *protocol_rtc.RTCSetParamResponse) *DeleteAllEntriesResponse {
	return &DeleteAllEntriesResponse{RTCSetParamResponse: *base}
}

// DeleteEntryResponse is the response to DeleteEntryRequest.
type DeleteEntryResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewDeleteEntryResponse wraps an RTCSetParamResponse as a DeleteEntryResponse.
func NewDeleteEntryResponse(base *protocol_rtc.RTCSetParamResponse) *DeleteEntryResponse {
	return &DeleteEntryResponse{RTCSetParamResponse: *base}
}

// AllocateEntryResponse is the response to AllocateEntryRequest.
type AllocateEntryResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewAllocateEntryResponse wraps an RTCSetParamResponse as an AllocateEntryResponse.
func NewAllocateEntryResponse(base *protocol_rtc.RTCSetParamResponse) *AllocateEntryResponse {
	return &AllocateEntryResponse{RTCSetParamResponse: *base}
}

// WriteEntryDataResponse is the response to WriteEntryDataRequest.
type WriteEntryDataResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewWriteEntryDataResponse wraps an RTCSetParamResponse as a WriteEntryDataResponse.
func NewWriteEntryDataResponse(base *protocol_rtc.RTCSetParamResponse) *WriteEntryDataResponse {
	return &WriteEntryDataResponse{RTCSetParamResponse: *base}
}

// GetEntrySizeResponse is the response to GetEntrySizeRequest.
type GetEntrySizeResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewGetEntrySizeResponse wraps an RTCSetParamResponse as a GetEntrySizeResponse.
func NewGetEntrySizeResponse(base *protocol_rtc.RTCSetParamResponse) *GetEntrySizeResponse {
	return &GetEntrySizeResponse{RTCSetParamResponse: *base}
}

// ReadEntryDataResponse is the response to ReadEntryDataRequest.
type ReadEntryDataResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewReadEntryDataResponse wraps an RTCSetParamResponse as a ReadEntryDataResponse.
func NewReadEntryDataResponse(base *protocol_rtc.RTCSetParamResponse) *ReadEntryDataResponse {
	return &ReadEntryDataResponse{RTCSetParamResponse: *base}
}

// PresenceMaskResponse is the response to PresenceMaskRequest.
type PresenceMaskResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewPresenceMaskResponse wraps an RTCSetParamResponse as a PresenceMaskResponse.
func NewPresenceMaskResponse(base *protocol_rtc.RTCSetParamResponse) *PresenceMaskResponse {
	return &PresenceMaskResponse{RTCSetParamResponse: *base}
}

func registerCommandTableResponseRegistries() {
	registerCTSetRegistry(CommandCode.SaveCommandTable, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewSaveCommandTableResponse(base)
	})
	registerCTSetRegistry(CommandCode.StopMotionController, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewStopMotionControllerResponse(base)
	})
	registerCTSetRegistry(CommandCode.StartMotionController, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewStartMotionControllerResponse(base)
	})
	registerCTSetRegistry(CommandCode.DeleteAllEntries, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewDeleteAllEntriesResponse(base)
	})
	registerCTSetRegistry(CommandCode.DeleteEntry, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewDeleteEntryResponse(base)
	})
	registerCTSetRegistry(CommandCode.AllocateEntry, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewAllocateEntryResponse(base)
	})
	registerCTSetRegistry(CommandCode.WriteEntryData, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewWriteEntryDataResponse(base)
	})
	registerCTSetRegistry(CommandCode.GetEntrySize, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewGetEntrySizeResponse(base)
	})
	registerCTSetRegistry(CommandCode.ReadEntryData, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewReadEntryDataResponse(base)
	})

	// PresenceMask0-7 share the same response shape but differ by cmdCode.
	registerCTSetRegistry(CommandCode.PresenceMask0, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewPresenceMaskResponse(base)
	})
	registerCTSetRegistry(CommandCode.PresenceMask1, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewPresenceMaskResponse(base)
	})
	registerCTSetRegistry(CommandCode.PresenceMask2, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewPresenceMaskResponse(base)
	})
	registerCTSetRegistry(CommandCode.PresenceMask3, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewPresenceMaskResponse(base)
	})
	registerCTSetRegistry(CommandCode.PresenceMask4, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewPresenceMaskResponse(base)
	})
	registerCTSetRegistry(CommandCode.PresenceMask5, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewPresenceMaskResponse(base)
	})
	registerCTSetRegistry(CommandCode.PresenceMask6, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewPresenceMaskResponse(base)
	})
	registerCTSetRegistry(CommandCode.PresenceMask7, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewPresenceMaskResponse(base)
	})
}

func registerCTSetRegistry(cmdCode uint8, wrapper func(*protocol_rtc.RTCSetParamResponse) protocol_common.Response) {
	protocol_rtc.RegisterResponseRegistryByCmd(cmdCode,
		func(status *protocol_common.Status, value int32, upid uint16, rtcCounter, rtcStatus, cmdCode uint8) protocol_common.Response {
			base := protocol_rtc.NewRTCSetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
			return wrapper(base)
		})
}
