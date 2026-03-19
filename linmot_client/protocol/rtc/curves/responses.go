package protocol_curves

import (
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/staged_robot/protocol/rtc"
)

// Compile-time interface checks
var (
	_ protocol_common.Response = (*SaveAllCurvesResponse)(nil)
	_ protocol_common.Response = (*DeleteAllCurvesResponse)(nil)
	_ protocol_common.Response = (*StartAddingCurveResponse)(nil)
	_ protocol_common.Response = (*AddCurveInfoBlockResponse)(nil)
	_ protocol_common.Response = (*AddCurveDataResponse)(nil)
	_ protocol_common.Response = (*StartModifyingCurveResponse)(nil)
	_ protocol_common.Response = (*ModifyCurveInfoBlockResponse)(nil)
	_ protocol_common.Response = (*ModifyCurveDataResponse)(nil)
	_ protocol_common.Response = (*StartGettingCurveResponse)(nil)
	_ protocol_common.Response = (*GetCurveInfoBlockResponse)(nil)
	_ protocol_common.Response = (*GetCurveDataResponse)(nil)
)

func init() {
	registerCurveResponseRegistries()
}

// SaveAllCurvesResponse is the response to SaveAllCurvesRequest.
type SaveAllCurvesResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewSaveAllCurvesResponse wraps an RTCSetParamResponse as a SaveAllCurvesResponse.
func NewSaveAllCurvesResponse(base *protocol_rtc.RTCSetParamResponse) *SaveAllCurvesResponse {
	return &SaveAllCurvesResponse{RTCSetParamResponse: *base}
}

// DeleteAllCurvesResponse is the response to DeleteAllCurvesRequest.
type DeleteAllCurvesResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewDeleteAllCurvesResponse wraps an RTCSetParamResponse as a DeleteAllCurvesResponse.
func NewDeleteAllCurvesResponse(base *protocol_rtc.RTCSetParamResponse) *DeleteAllCurvesResponse {
	return &DeleteAllCurvesResponse{RTCSetParamResponse: *base}
}

// StartAddingCurveResponse is the response to StartAddingCurveRequest.
type StartAddingCurveResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewStartAddingCurveResponse wraps an RTCSetParamResponse as a StartAddingCurveResponse.
func NewStartAddingCurveResponse(base *protocol_rtc.RTCSetParamResponse) *StartAddingCurveResponse {
	return &StartAddingCurveResponse{RTCSetParamResponse: *base}
}

// AddCurveInfoBlockResponse is the response to AddCurveInfoBlockRequest.
type AddCurveInfoBlockResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewAddCurveInfoBlockResponse wraps an RTCSetParamResponse as an AddCurveInfoBlockResponse.
func NewAddCurveInfoBlockResponse(base *protocol_rtc.RTCSetParamResponse) *AddCurveInfoBlockResponse {
	return &AddCurveInfoBlockResponse{RTCSetParamResponse: *base}
}

// AddCurveDataResponse is the response to AddCurveDataRequest.
type AddCurveDataResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewAddCurveDataResponse wraps an RTCSetParamResponse as an AddCurveDataResponse.
func NewAddCurveDataResponse(base *protocol_rtc.RTCSetParamResponse) *AddCurveDataResponse {
	return &AddCurveDataResponse{RTCSetParamResponse: *base}
}

// StartModifyingCurveResponse is the response to StartModifyingCurveRequest.
type StartModifyingCurveResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewStartModifyingCurveResponse wraps an RTCSetParamResponse as a StartModifyingCurveResponse.
func NewStartModifyingCurveResponse(base *protocol_rtc.RTCSetParamResponse) *StartModifyingCurveResponse {
	return &StartModifyingCurveResponse{RTCSetParamResponse: *base}
}

// ModifyCurveInfoBlockResponse is the response to ModifyCurveInfoBlockRequest.
type ModifyCurveInfoBlockResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewModifyCurveInfoBlockResponse wraps an RTCSetParamResponse as a ModifyCurveInfoBlockResponse.
func NewModifyCurveInfoBlockResponse(base *protocol_rtc.RTCSetParamResponse) *ModifyCurveInfoBlockResponse {
	return &ModifyCurveInfoBlockResponse{RTCSetParamResponse: *base}
}

// ModifyCurveDataResponse is the response to ModifyCurveDataRequest.
type ModifyCurveDataResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewModifyCurveDataResponse wraps an RTCSetParamResponse as a ModifyCurveDataResponse.
func NewModifyCurveDataResponse(base *protocol_rtc.RTCSetParamResponse) *ModifyCurveDataResponse {
	return &ModifyCurveDataResponse{RTCSetParamResponse: *base}
}

// StartGettingCurveResponse is the response to StartGettingCurveRequest.
type StartGettingCurveResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewStartGettingCurveResponse wraps an RTCSetParamResponse as a StartGettingCurveResponse.
func NewStartGettingCurveResponse(base *protocol_rtc.RTCSetParamResponse) *StartGettingCurveResponse {
	return &StartGettingCurveResponse{RTCSetParamResponse: *base}
}

// GetCurveInfoBlockResponse is the response to GetCurveInfoBlockRequest.
type GetCurveInfoBlockResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewGetCurveInfoBlockResponse wraps an RTCSetParamResponse as a GetCurveInfoBlockResponse.
func NewGetCurveInfoBlockResponse(base *protocol_rtc.RTCSetParamResponse) *GetCurveInfoBlockResponse {
	return &GetCurveInfoBlockResponse{RTCSetParamResponse: *base}
}

// InfoBlockSize returns the info block size from the response value (bits 0-15).
func (r *GetCurveInfoBlockResponse) InfoBlockSize() uint16 {
	return uint16(r.Value() & 0xFFFF)
}

// DataBlockSize returns the data block size from the response value (bits 16-31).
func (r *GetCurveInfoBlockResponse) DataBlockSize() uint16 {
	return uint16(r.Value() >> 16)
}

// GetCurveDataResponse is the response to GetCurveDataRequest.
type GetCurveDataResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewGetCurveDataResponse wraps an RTCSetParamResponse as a GetCurveDataResponse.
func NewGetCurveDataResponse(base *protocol_rtc.RTCSetParamResponse) *GetCurveDataResponse {
	return &GetCurveDataResponse{RTCSetParamResponse: *base}
}

// DataLow returns the low word of curve data from the response value (bits 0-15).
func (r *GetCurveDataResponse) DataLow() uint16 {
	return uint16(r.Value() & 0xFFFF)
}

// DataHigh returns the high word of curve data from the response value (bits 16-31).
func (r *GetCurveDataResponse) DataHigh() uint16 {
	return uint16(r.Value() >> 16)
}

func registerCurveResponseRegistries() {
	registerCurveSetRegistry(protocol_rtc.CommandCode.SaveAllCurvesToFlash, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewSaveAllCurvesResponse(base)
	})
	registerCurveSetRegistry(protocol_rtc.CommandCode.DeleteAllCurves, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewDeleteAllCurvesResponse(base)
	})
	registerCurveSetRegistry(protocol_rtc.CommandCode.StartAddingCurve, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewStartAddingCurveResponse(base)
	})
	registerCurveSetRegistry(protocol_rtc.CommandCode.AddCurveInfoBlock, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewAddCurveInfoBlockResponse(base)
	})
	registerCurveSetRegistry(protocol_rtc.CommandCode.AddCurveData, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewAddCurveDataResponse(base)
	})
	registerCurveSetRegistry(protocol_rtc.CommandCode.StartModifyingCurve, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewStartModifyingCurveResponse(base)
	})
	registerCurveSetRegistry(protocol_rtc.CommandCode.ModifyCurveInfoBlock, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewModifyCurveInfoBlockResponse(base)
	})
	registerCurveSetRegistry(protocol_rtc.CommandCode.ModifyCurveData, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewModifyCurveDataResponse(base)
	})
	registerCurveSetRegistry(protocol_rtc.CommandCode.StartGettingCurve, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewStartGettingCurveResponse(base)
	})
	registerCurveSetRegistry(protocol_rtc.CommandCode.GetCurveInfoBlock, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewGetCurveInfoBlockResponse(base)
	})
	registerCurveSetRegistry(protocol_rtc.CommandCode.GetCurveData, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewGetCurveDataResponse(base)
	})
}

func registerCurveSetRegistry(cmdCode uint8, wrapper func(*protocol_rtc.RTCSetParamResponse) protocol_common.Response) {
	protocol_rtc.RegisterResponseRegistryByCmd(cmdCode,
		func(status *protocol_common.Status, value int32, upid uint16, rtcCounter, rtcStatus, cmdCode uint8) protocol_common.Response {
			base := protocol_rtc.NewRTCSetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
			return wrapper(base)
		})
}
