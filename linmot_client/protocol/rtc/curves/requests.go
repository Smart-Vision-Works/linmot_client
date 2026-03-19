package protocol_curves

import (
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/staged_robot/protocol/rtc"
)

// Compile-time interface checks
var (
	_ protocol_common.Request = (*SaveAllCurvesRequest)(nil)
	_ protocol_common.Request = (*DeleteAllCurvesRequest)(nil)
	_ protocol_common.Request = (*StartAddingCurveRequest)(nil)
	_ protocol_common.Request = (*AddCurveInfoBlockRequest)(nil)
	_ protocol_common.Request = (*AddCurveDataRequest)(nil)
	_ protocol_common.Request = (*StartModifyingCurveRequest)(nil)
	_ protocol_common.Request = (*ModifyCurveInfoBlockRequest)(nil)
	_ protocol_common.Request = (*ModifyCurveDataRequest)(nil)
	_ protocol_common.Request = (*StartGettingCurveRequest)(nil)
	_ protocol_common.Request = (*GetCurveInfoBlockRequest)(nil)
	_ protocol_common.Request = (*GetCurveDataRequest)(nil)
)

// SaveAllCurvesRequest saves all curves from RAM to Flash.
type SaveAllCurvesRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewSaveAllCurvesRequest creates a new request to save all curves to Flash.
func NewSaveAllCurvesRequest() *SaveAllCurvesRequest {
	return &SaveAllCurvesRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(0, 0, protocol_rtc.CommandCode.SaveAllCurvesToFlash),
	}
}

// DeleteAllCurvesRequest deletes all curves from RAM.
type DeleteAllCurvesRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewDeleteAllCurvesRequest creates a new request to delete all curves from RAM.
func NewDeleteAllCurvesRequest() *DeleteAllCurvesRequest {
	return &DeleteAllCurvesRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(0, 0, protocol_rtc.CommandCode.DeleteAllCurves),
	}
}

// StartAddingCurveRequest starts adding a curve to RAM.
type StartAddingCurveRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewStartAddingCurveRequest creates a new request to start adding a curve.
// curveID is the curve number (1-100).
func NewStartAddingCurveRequest(curveID uint16) *StartAddingCurveRequest {
	return &StartAddingCurveRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(curveID, 0, protocol_rtc.CommandCode.StartAddingCurve),
	}
}

// AddCurveInfoBlockRequest adds a curve info block.
type AddCurveInfoBlockRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewAddCurveInfoBlockRequest creates a new request to add curve info block data.
// curveID is the curve number, dataLow and dataHigh are the info block data words.
func NewAddCurveInfoBlockRequest(curveID uint16, dataLow, dataHigh uint16) *AddCurveInfoBlockRequest {
	value := int32(uint32(dataLow) | (uint32(dataHigh) << 16))
	return &AddCurveInfoBlockRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(curveID, value, protocol_rtc.CommandCode.AddCurveInfoBlock),
	}
}

// AddCurveDataRequest adds curve data.
type AddCurveDataRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewAddCurveDataRequest creates a new request to add curve data.
// curveID is the curve number, dataLow and dataHigh are the curve data words.
func NewAddCurveDataRequest(curveID uint16, dataLow, dataHigh uint16) *AddCurveDataRequest {
	value := int32(uint32(dataLow) | (uint32(dataHigh) << 16))
	return &AddCurveDataRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(curveID, value, protocol_rtc.CommandCode.AddCurveData),
	}
}

// StartModifyingCurveRequest starts modifying an existing curve.
type StartModifyingCurveRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewStartModifyingCurveRequest creates a new request to start modifying an existing curve.
// curveID is the curve number (1-100).
func NewStartModifyingCurveRequest(curveID uint16) *StartModifyingCurveRequest {
	return &StartModifyingCurveRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(curveID, 0, protocol_rtc.CommandCode.StartModifyingCurve),
	}
}

// ModifyCurveInfoBlockRequest modifies a curve info block.
type ModifyCurveInfoBlockRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewModifyCurveInfoBlockRequest creates a new request to modify curve info block data.
// curveID is the curve number, dataLow and dataHigh are the info block data words.
func NewModifyCurveInfoBlockRequest(curveID uint16, dataLow, dataHigh uint16) *ModifyCurveInfoBlockRequest {
	value := int32(uint32(dataLow) | (uint32(dataHigh) << 16))
	return &ModifyCurveInfoBlockRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(curveID, value, protocol_rtc.CommandCode.ModifyCurveInfoBlock),
	}
}

// ModifyCurveDataRequest modifies curve data.
type ModifyCurveDataRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewModifyCurveDataRequest creates a new request to modify curve data.
// curveID is the curve number, dataLow and dataHigh are the curve data words.
func NewModifyCurveDataRequest(curveID uint16, dataLow, dataHigh uint16) *ModifyCurveDataRequest {
	value := int32(uint32(dataLow) | (uint32(dataHigh) << 16))
	return &ModifyCurveDataRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(curveID, value, protocol_rtc.CommandCode.ModifyCurveData),
	}
}

// StartGettingCurveRequest starts getting a curve from RAM.
type StartGettingCurveRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewStartGettingCurveRequest creates a new request to start getting a curve.
// curveID is the curve number (1-100).
func NewStartGettingCurveRequest(curveID uint16) *StartGettingCurveRequest {
	return &StartGettingCurveRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(curveID, 0, protocol_rtc.CommandCode.StartGettingCurve),
	}
}

// GetCurveInfoBlockRequest gets a curve info block.
type GetCurveInfoBlockRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewGetCurveInfoBlockRequest creates a new request to get curve info block data.
// curveID is the curve number.
func NewGetCurveInfoBlockRequest(curveID uint16) *GetCurveInfoBlockRequest {
	return &GetCurveInfoBlockRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(curveID, 0, protocol_rtc.CommandCode.GetCurveInfoBlock),
	}
}

// GetCurveDataRequest gets curve data.
type GetCurveDataRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewGetCurveDataRequest creates a new request to get curve data.
// curveID is the curve number.
func NewGetCurveDataRequest(curveID uint16) *GetCurveDataRequest {
	return &GetCurveDataRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(curveID, 0, protocol_rtc.CommandCode.GetCurveData),
	}
}
