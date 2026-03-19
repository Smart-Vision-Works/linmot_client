package client_common

import (
	"context"
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	"sync"
	"time"
)

type pendingRequest struct {
	request            protocol_common.Request
	commandCount       uint8  // RTC counter (1-14 range) or 0 for non-RTC
	mcCounter          uint8  // MC counter (1-4 range) or 0 for non-MC
	originalUPID       uint16 // UPID from request for response matching (0 for non-parameter requests)
	responseChannel    chan protocol_common.Response
	errorChannel       chan error
	requestCtx         context.Context // Caller's context (checked in startRequest before sending)
	totalTimeout       time.Duration   // Total timeout duration (set when request is created)
	perAttemptTimeout  time.Duration   // Timeout per retry attempt (set when request is created)
	deadline           time.Time       // Deadline calculated in startRequest() when request is actually sent
	sendCounter        uint            // Tracks the number of times this request has been sent (starts at 1)
	lastSendTime       time.Time       // Time of last send attempt (set in startRequest())
	sentCh             chan struct{}   // Closed when first packet is successfully sent (signals timeout can start)
	sentOnce           sync.Once       // Ensures sentCh is closed exactly once
	completeOnce       sync.Once       // Ensures only one completion (response/error/timeout) can win
	isStatusRequest    bool            // True for Status/MonitoringStatus/ControlWord (non-RTC, non-MC)
	isMCRequest        bool            // True for Motion Control requests
	invalidStatusCount uint8           // Count of rejected status packets for this request
	statusTxHex        string          // Last status request packet hex (debug-only)
}

func newPendingRequest(
	request protocol_common.Request,
	commandCount uint8,
	mcCounter uint8,
	originalUPID uint16,
	responseChannel chan protocol_common.Response,
	errorChannel chan error,
	requestCtx context.Context,
	totalTimeout time.Duration,
	perAttemptTimeout time.Duration,
	isStatusRequest bool,
	isMCRequest bool,
) *pendingRequest {
	return &pendingRequest{
		request:           request,
		commandCount:      commandCount,
		mcCounter:         mcCounter,
		originalUPID:      originalUPID,
		responseChannel:   responseChannel,
		errorChannel:      errorChannel,
		requestCtx:        requestCtx,
		totalTimeout:      totalTimeout,
		perAttemptTimeout: perAttemptTimeout,
		deadline:          time.Time{},         // Will be set in startRequest() when request is actually sent
		sendCounter:       0,                   // Will be set to 1 in startRequest()
		lastSendTime:      time.Time{},         // Will be set in startRequest()
		sentCh:            make(chan struct{}), // Will be closed in startRequest() after first successful send
		isStatusRequest:   isStatusRequest,
		isMCRequest:       isMCRequest,
		statusTxHex:       "",
	}
}

// signalFirstSendComplete safely closes sentCh exactly once.
// This prevents deadlocks if multiple paths try to close sentCh.
func (req *pendingRequest) signalFirstSendComplete() {
	req.sentOnce.Do(func() {
		close(req.sentCh)
	})
}

// tryCompleteWithResponse attempts to complete the request with a response.
// Returns true if this was the winning completion, false if another completion already won.
func (req *pendingRequest) tryCompleteWithResponse(resp protocol_common.Response) bool {
	var sent bool
	req.completeOnce.Do(func() {
		select {
		case req.responseChannel <- resp:
			sent = true
		default:
			// Channel full or closed - another completion already won
		}
	})
	return sent
}

// tryCompleteWithError attempts to complete the request with an error.
// Returns true if this was the winning completion, false if another completion already won.
func (req *pendingRequest) tryCompleteWithError(err error) bool {
	var sent bool
	req.completeOnce.Do(func() {
		select {
		case req.errorChannel <- err:
			sent = true
		default:
			// Channel full or closed - another completion already won
		}
	})
	return sent
}
