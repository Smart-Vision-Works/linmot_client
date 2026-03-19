package protocol_common

import (
	"fmt"
	"time"
)

// Protocol error codes
const (
	ErrCodePacketTooShort     = "PACKET_TOO_SHORT"
	ErrCodeInvalidFlags       = "INVALID_FLAGS"
	ErrCodeInvalidUPID        = "INVALID_UPID"
	ErrCodeInvalidCmdCode     = "INVALID_CMD_CODE"
	ErrCodeInvalidValue       = "INVALID_VALUE"
	ErrCodeUnknownRequestType = "UNKNOWN_REQUEST_TYPE"
)

// ProtocolError represents an error that occurred during protocol operations.
// It provides structured information about what went wrong, including the
// operation, error code, message, and the raw packet data for debugging.
type ProtocolError struct {
	Op      string // Operation that failed (e.g., "ParseStatus", "BuildRTCRead")
	Code    string // Error code (e.g., "PACKET_TOO_SHORT", "INVALID_FLAGS")
	Message string // Human-readable error message
	Data    []byte // The problematic packet data (for debugging)
}

// Error implements the error interface.
func (e *ProtocolError) Error() string {
	if len(e.Data) > 0 {
		// Limit data display to first 32 bytes
		dataLen := len(e.Data)
		if dataLen > 32 {
			return fmt.Sprintf("linudp protocol %s: %s (%s) [data: %d bytes: %x...]",
				e.Op, e.Message, e.Code, dataLen, e.Data[:32])
		}
		return fmt.Sprintf("linudp protocol %s: %s (%s) [data: %x]",
			e.Op, e.Message, e.Code, e.Data)
	}
	return fmt.Sprintf("linudp protocol %s: %s (%s)", e.Op, e.Message, e.Code)
}

// NewProtocolError creates a new ProtocolError.
func NewProtocolError(op, code, message string, data []byte) *ProtocolError {
	return &ProtocolError{
		Op:      op,
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewPacketTooShortError creates a ProtocolError for packets that are too short.
func NewPacketTooShortError(op string, actual, required int, data []byte) *ProtocolError {
	return &ProtocolError{
		Op:      op,
		Code:    ErrCodePacketTooShort,
		Message: fmt.Sprintf("packet too short: %d bytes, need %d", actual, required),
		Data:    data,
	}
}

// NewInvalidFlagsError creates a ProtocolError for invalid packet flags.
func NewInvalidFlagsError(op string, flags, expected uint32, data []byte) *ProtocolError {
	return &ProtocolError{
		Op:      op,
		Code:    ErrCodeInvalidFlags,
		Message: fmt.Sprintf("invalid flags: 0x%08X (expected at least 0x%08X)", flags, expected),
		Data:    data,
	}
}

// NewInvalidUPIDError creates a ProtocolError for invalid UPIDs.
func NewInvalidUPIDError(op string, upid uint16, reason string) *ProtocolError {
	return &ProtocolError{
		Op:      op,
		Code:    ErrCodeInvalidUPID,
		Message: fmt.Sprintf("invalid upid 0x%04X: %s", upid, reason),
		Data:    nil,
	}
}

// NewInvalidCmdCodeError creates a ProtocolError for invalid command codes.
func NewInvalidCmdCodeError(op string, cmdCode uint8, validCodes []uint8) *ProtocolError {
	return &ProtocolError{
		Op:      op,
		Code:    ErrCodeInvalidCmdCode,
		Message: fmt.Sprintf("invalid command code 0x%02X (valid: %v)", cmdCode, validCodes),
		Data:    nil,
	}
}

// InvalidEventError is returned when SendPacket() receives an unknown event type.
type InvalidEventError struct {
	Operation string
	Event     interface{} // Can be PacketWriteable, PacketReadable, or Request
	Message   string      // Optional additional message
}

func (e *InvalidEventError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: invalid event type %T: %s", e.Operation, e.Event, e.Message)
	}
	return fmt.Sprintf("%s: invalid event type %T", e.Operation, e.Event)
}

// RequestTimeoutError is returned when a request times out after all retry attempts.
type RequestTimeoutError struct {
	Attempts   int           // Number of retry attempts made
	Timeout    time.Duration // Total timeout duration
	LocalAddr  string        // Local UDP bind address (for diagnostics)
	RemoteAddr string        // Remote drive address (for diagnostics)
}

func (e *RequestTimeoutError) Error() string {
	base := fmt.Sprintf("request timeout after %d attempts (%v)", e.Attempts, e.Timeout)
	if e.LocalAddr != "" || e.RemoteAddr != "" {
		base += fmt.Sprintf(" [local=%s remote=%s]", e.LocalAddr, e.RemoteAddr)
	}
	return base
}

// DiagnosticHint returns a helpful message suggesting possible causes for the timeout.
func (e *RequestTimeoutError) DiagnosticHint() string {
	return `Possible causes:
  1. Another master is connected (LinMot-Talk, primer, or another client)
     - LinMot only responds to ONE master at a time
     - Close LinMot-Talk or stop the primer before running tests
  2. Wrong network interface - check local bind address above
  3. Firewall blocking UDP port 49360
  4. LinMot powered off or not on this network segment`
}
