package protocol_errors

import (
	protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/linmot_client/protocol/rtc"
)

// Compile-time interface checks
var (
	_ protocol_common.Response = (*GetErrorLogEntryCounterResponse)(nil)
	_ protocol_common.Response = (*GetErrorLogEntryCodeResponse)(nil)
	_ protocol_common.Response = (*GetErrorLogEntryTimeLowResponse)(nil)
	_ protocol_common.Response = (*GetErrorLogEntryTimeHighResponse)(nil)
	_ protocol_common.Response = (*GetErrorCodeTextStringletResponse)(nil)
)

func init() {
	registerErrorResponseRegistries()
}

// GetErrorLogEntryCounterResponse is the response to GetErrorLogEntryCounterRequest.
type GetErrorLogEntryCounterResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewGetErrorLogEntryCounterResponse wraps an RTCSetParamResponse as a GetErrorLogEntryCounterResponse.
func NewGetErrorLogEntryCounterResponse(base *protocol_rtc.RTCSetParamResponse) *GetErrorLogEntryCounterResponse {
	return &GetErrorLogEntryCounterResponse{RTCSetParamResponse: *base}
}

// LoggedErrorsCount returns the number of logged errors (low 16 bits).
func (r *GetErrorLogEntryCounterResponse) LoggedErrorsCount() uint16 {
	return uint16(r.Value() & 0xFFFF)
}

// OccurredErrorsCount returns the number of occurred errors (high 16 bits).
func (r *GetErrorLogEntryCounterResponse) OccurredErrorsCount() uint16 {
	return uint16(r.Value() >> 16)
}

// GetErrorLogEntryCodeResponse is the response to GetErrorLogEntryCodeRequest.
type GetErrorLogEntryCodeResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewGetErrorLogEntryCodeResponse wraps an RTCSetParamResponse as a GetErrorLogEntryCodeResponse.
func NewGetErrorLogEntryCodeResponse(base *protocol_rtc.RTCSetParamResponse) *GetErrorLogEntryCodeResponse {
	return &GetErrorLogEntryCodeResponse{RTCSetParamResponse: *base}
}

// EntryNumber returns the entry number from the response UPID field.
func (r *GetErrorLogEntryCodeResponse) EntryNumber() uint16 {
	return r.UPID()
}

// ErrorCode returns the logged error code (low 16 bits).
func (r *GetErrorLogEntryCodeResponse) ErrorCode() uint16 {
	return uint16(r.Value() & 0xFFFF)
}

// GetErrorLogEntryTimeLowResponse is the response to GetErrorLogEntryTimeLowRequest.
type GetErrorLogEntryTimeLowResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewGetErrorLogEntryTimeLowResponse wraps an RTCSetParamResponse as a GetErrorLogEntryTimeLowResponse.
func NewGetErrorLogEntryTimeLowResponse(base *protocol_rtc.RTCSetParamResponse) *GetErrorLogEntryTimeLowResponse {
	return &GetErrorLogEntryTimeLowResponse{RTCSetParamResponse: *base}
}

// EntryNumber returns the entry number from the response UPID field.
func (r *GetErrorLogEntryTimeLowResponse) EntryNumber() uint16 {
	return r.UPID()
}

// TimeLowWord returns the entry time low word (low 16 bits).
func (r *GetErrorLogEntryTimeLowResponse) TimeLowWord() uint16 {
	return uint16(r.Value() & 0xFFFF)
}

// TimeMidLowWord returns the entry time mid-low word (high 16 bits).
func (r *GetErrorLogEntryTimeLowResponse) TimeMidLowWord() uint16 {
	return uint16(r.Value() >> 16)
}

// GetErrorLogEntryTimeHighResponse is the response to GetErrorLogEntryTimeHighRequest.
type GetErrorLogEntryTimeHighResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewGetErrorLogEntryTimeHighResponse wraps an RTCSetParamResponse as a GetErrorLogEntryTimeHighResponse.
func NewGetErrorLogEntryTimeHighResponse(base *protocol_rtc.RTCSetParamResponse) *GetErrorLogEntryTimeHighResponse {
	return &GetErrorLogEntryTimeHighResponse{RTCSetParamResponse: *base}
}

// EntryNumber returns the entry number from the response UPID field.
func (r *GetErrorLogEntryTimeHighResponse) EntryNumber() uint16 {
	return r.UPID()
}

// TimeMidHighWord returns the entry time mid-high word (low 16 bits).
func (r *GetErrorLogEntryTimeHighResponse) TimeMidHighWord() uint16 {
	return uint16(r.Value() & 0xFFFF)
}

// TimeHighWord returns the entry time high word (high 16 bits).
func (r *GetErrorLogEntryTimeHighResponse) TimeHighWord() uint16 {
	return uint16(r.Value() >> 16)
}

// GetErrorCodeTextStringletResponse is the response to GetErrorCodeTextStringletRequest.
type GetErrorCodeTextStringletResponse struct {
	protocol_rtc.RTCSetParamResponse
}

// NewGetErrorCodeTextStringletResponse wraps an RTCSetParamResponse as a GetErrorCodeTextStringletResponse.
func NewGetErrorCodeTextStringletResponse(base *protocol_rtc.RTCSetParamResponse) *GetErrorCodeTextStringletResponse {
	return &GetErrorCodeTextStringletResponse{RTCSetParamResponse: *base}
}

// ErrorCode returns the error code from the response UPID field.
func (r *GetErrorCodeTextStringletResponse) ErrorCode() uint16 {
	return r.UPID()
}

// StringletBytes returns the 4 bytes of text from the stringlet.
func (r *GetErrorCodeTextStringletResponse) StringletBytes() [4]byte {
	val := uint32(r.Value())
	return [4]byte{
		byte(val & 0xFF),
		byte((val >> 8) & 0xFF),
		byte((val >> 16) & 0xFF),
		byte((val >> 24) & 0xFF),
	}
}

func registerErrorResponseRegistries() {
	registerErrorSetRegistry(protocol_rtc.CommandCode.GetErrorLogEntryCounter, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewGetErrorLogEntryCounterResponse(base)
	})
	registerErrorSetRegistry(protocol_rtc.CommandCode.GetErrorLogEntryCode, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewGetErrorLogEntryCodeResponse(base)
	})
	registerErrorSetRegistry(protocol_rtc.CommandCode.GetErrorLogEntryTimeLow, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewGetErrorLogEntryTimeLowResponse(base)
	})
	registerErrorSetRegistry(protocol_rtc.CommandCode.GetErrorLogEntryTimeHigh, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewGetErrorLogEntryTimeHighResponse(base)
	})
	registerErrorSetRegistry(protocol_rtc.CommandCode.GetErrorCodeTextStringlet, func(base *protocol_rtc.RTCSetParamResponse) protocol_common.Response {
		return NewGetErrorCodeTextStringletResponse(base)
	})
}

func registerErrorSetRegistry(cmdCode uint8, wrapper func(*protocol_rtc.RTCSetParamResponse) protocol_common.Response) {
	protocol_rtc.RegisterResponseRegistryByCmd(cmdCode,
		func(status *protocol_common.Status, value int32, upid uint16, rtcCounter, rtcStatus, cmdCode uint8) protocol_common.Response {
			base := protocol_rtc.NewRTCSetParamResponseWithCmdCode(status, value, upid, rtcCounter, rtcStatus, cmdCode)
			return wrapper(base)
		})
}
