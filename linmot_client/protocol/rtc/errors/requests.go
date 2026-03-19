package protocol_errors

import (
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/staged_robot/protocol/rtc"
)

// Compile-time interface checks
var (
	_ protocol_common.Request = (*GetErrorLogEntryCounterRequest)(nil)
	_ protocol_common.Request = (*GetErrorLogEntryCodeRequest)(nil)
	_ protocol_common.Request = (*GetErrorLogEntryTimeLowRequest)(nil)
	_ protocol_common.Request = (*GetErrorLogEntryTimeHighRequest)(nil)
	_ protocol_common.Request = (*GetErrorCodeTextStringletRequest)(nil)
)

// GetErrorLogEntryCounterRequest gets the error log entry counter.
type GetErrorLogEntryCounterRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewGetErrorLogEntryCounterRequest creates a new request to get error log entry counters.
func NewGetErrorLogEntryCounterRequest() *GetErrorLogEntryCounterRequest {
	return &GetErrorLogEntryCounterRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(0, 0, protocol_rtc.CommandCode.GetErrorLogEntryCounter),
	}
}

// GetErrorLogEntryCodeRequest gets the error code for a specific error log entry.
type GetErrorLogEntryCodeRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewGetErrorLogEntryCodeRequest creates a new request to get the error code for a specific entry.
// entryNumber is the error log entry number (0-20).
func NewGetErrorLogEntryCodeRequest(entryNumber uint16) *GetErrorLogEntryCodeRequest {
	return &GetErrorLogEntryCodeRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(entryNumber, 0, protocol_rtc.CommandCode.GetErrorLogEntryCode),
	}
}

// GetErrorLogEntryTimeLowRequest gets the low 32 bits of the error log entry time.
type GetErrorLogEntryTimeLowRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewGetErrorLogEntryTimeLowRequest creates a new request to get the error log entry time (low 32 bits).
// entryNumber is the error log entry number (0-20).
func NewGetErrorLogEntryTimeLowRequest(entryNumber uint16) *GetErrorLogEntryTimeLowRequest {
	return &GetErrorLogEntryTimeLowRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(entryNumber, 0, protocol_rtc.CommandCode.GetErrorLogEntryTimeLow),
	}
}

// GetErrorLogEntryTimeHighRequest gets the high 32 bits of the error log entry time.
type GetErrorLogEntryTimeHighRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewGetErrorLogEntryTimeHighRequest creates a new request to get the error log entry time (high 32 bits).
// entryNumber is the error log entry number (0-20).
func NewGetErrorLogEntryTimeHighRequest(entryNumber uint16) *GetErrorLogEntryTimeHighRequest {
	return &GetErrorLogEntryTimeHighRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(entryNumber, 0, protocol_rtc.CommandCode.GetErrorLogEntryTimeHigh),
	}
}

// GetErrorCodeTextStringletRequest gets a text stringlet for an error code.
type GetErrorCodeTextStringletRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewGetErrorCodeTextStringletRequest creates a new request to get error code text stringlet.
// errorCode is the error code to get text for.
// stringletNumber is the stringlet number (0-7).
func NewGetErrorCodeTextStringletRequest(errorCode, stringletNumber uint16) *GetErrorCodeTextStringletRequest {
	value := int32(uint32(stringletNumber))
	return &GetErrorCodeTextStringletRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(errorCode, value, protocol_rtc.CommandCode.GetErrorCodeTextStringlet),
	}
}
