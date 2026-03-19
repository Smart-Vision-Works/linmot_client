package client_errors

import (
	"context"
	"fmt"
	"time"

	client_common "gsail-go/linmot/client/common"
	protocol_errors "gsail-go/linmot/protocol/rtc/errors"
)

type ErrorManager struct {
	requestManager *client_common.RequestManager
}

func NewErrorManager(requestManager *client_common.RequestManager) *ErrorManager {
	return &ErrorManager{
		requestManager: requestManager,
	}
}

// ErrorLogEntry represents a single error log entry.
type ErrorLogEntry struct {
	Index       uint16
	ErrorCode   uint16
	Timestamp   time.Time
	Description string // Human-readable error text (empty if not retrieved)
}

// GetErrorLogCounts returns the count of logged and occurred errors.
func (m *ErrorManager) GetErrorLogCounts(ctx context.Context) (logged, occurred uint16, err error) {
	request := protocol_errors.NewGetErrorLogEntryCounterRequest()
	response, err := client_common.SendRequestAndReceive[*protocol_errors.GetErrorLogEntryCounterResponse](m.requestManager, ctx, request)
	if err != nil {
		return 0, 0, err
	}
	// Extract logged and occurred using the typed response methods
	logged = response.LoggedErrorsCount()
	occurred = response.OccurredErrorsCount()
	return logged, occurred, nil
}

// GetErrorLogEntry retrieves a single error log entry by index.
func (m *ErrorManager) GetErrorLogEntry(ctx context.Context, entryNum uint16) (*ErrorLogEntry, error) {
	// Get error code
	codeReq := protocol_errors.NewGetErrorLogEntryCodeRequest(entryNum)
	codeResp, err := client_common.SendRequestAndReceive[*protocol_errors.GetErrorLogEntryCodeResponse](m.requestManager, ctx, codeReq)
	if err != nil {
		return nil, err
	}
	errorCode := codeResp.ErrorCode()

	// Get timestamp (low 32 bits)
	timeLowReq := protocol_errors.NewGetErrorLogEntryTimeLowRequest(entryNum)
	timeLowResp, err := client_common.SendRequestAndReceive[*protocol_errors.GetErrorLogEntryTimeLowResponse](m.requestManager, ctx, timeLowReq)
	if err != nil {
		return nil, err
	}
	timeLowWord := timeLowResp.TimeLowWord()
	timeMidLowWord := timeLowResp.TimeMidLowWord()

	// Get timestamp (high 32 bits)
	timeHighReq := protocol_errors.NewGetErrorLogEntryTimeHighRequest(entryNum)
	timeHighResp, err := client_common.SendRequestAndReceive[*protocol_errors.GetErrorLogEntryTimeHighResponse](m.requestManager, ctx, timeHighReq)
	if err != nil {
		return nil, err
	}
	timeMidHighWord := timeHighResp.TimeMidHighWord()
	timeHighWord := timeHighResp.TimeHighWord()

	// Combine timestamp words into a 64-bit millisecond counter.
	timestampMillisecond := uint64(timeLowWord) |
		uint64(timeMidLowWord)<<16 |
		uint64(timeMidHighWord)<<32 |
		uint64(timeHighWord)<<48

	// Convert to time.Time assuming the counter represents milliseconds since epoch.
	timestamp := time.Unix(0, int64(timestampMillisecond)*int64(time.Millisecond))

	return &ErrorLogEntry{
		Index:     entryNum,
		ErrorCode: errorCode,
		Timestamp: timestamp,
	}, nil
}

// GetErrorLog retrieves all error log entries from the drive.
func (m *ErrorManager) GetErrorLog(ctx context.Context) ([]ErrorLogEntry, error) {
	// Get count of logged errors
	logged, _, err := m.GetErrorLogCounts(ctx)
	if err != nil {
		return nil, err
	}

	// Retrieve each error log entry
	entries := make([]ErrorLogEntry, 0, logged)
	for i := uint16(0); i < logged; i++ {
		entry, err := m.GetErrorLogEntry(ctx, i)
		if err != nil {
			return nil, err
		}
		entries = append(entries, *entry)
	}

	return entries, nil
}

// GetErrorText retrieves the human-readable description for an error code.
// Error text is stored as up to 8 stringlets of 4 bytes each (max 32 bytes).
// The function iterates through stringlets until it finds a null terminator or reaches the maximum.
func (m *ErrorManager) GetErrorText(ctx context.Context, errorCode uint16) (string, error) {
	var textBytes []byte
	const maxStringlets = 8

	for stringletNum := uint16(0); stringletNum < maxStringlets; stringletNum++ {
		// Request this stringlet
		request := protocol_errors.NewGetErrorCodeTextStringletRequest(errorCode, stringletNum)
		response, err := client_common.SendRequestAndReceive[*protocol_errors.GetErrorCodeTextStringletResponse](m.requestManager, ctx, request)
		if err != nil {
			return "", err
		}

		// Extract 4 bytes from the response
		stringletBytes := response.StringletBytes()

		// Append bytes to text buffer, stopping at null terminator
		for _, b := range stringletBytes {
			if b == 0 {
				// Found null terminator - end of string
				return string(textBytes), nil
			}
			textBytes = append(textBytes, b)
		}
	}

	// Reached max stringlets without finding null terminator
	return string(textBytes), nil
}

// GetErrorLogWithText retrieves all error log entries from the drive with human-readable descriptions.
func (m *ErrorManager) GetErrorLogWithText(ctx context.Context) ([]ErrorLogEntry, error) {
	entries, err := m.GetErrorLog(ctx)
	if err != nil {
		return nil, err
	}

	for i := range entries {
		text, textErr := m.GetErrorText(ctx, entries[i].ErrorCode)
		if textErr != nil {
			entries[i].Description = fmt.Sprintf("description unavailable: %v", textErr)
			continue
		}
		entries[i].Description = text
	}

	return entries, nil
}
