package protocol_common

import (
	"time"
)

// Request is the base interface for all request types.
type Request interface {
	// OperationTimeout returns the appropriate timeout duration for this request.
	OperationTimeout() time.Duration
}
