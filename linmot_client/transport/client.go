package transport

import (
	"context"
)

// Client sends and receives raw bytes to/from a LinMot drive.
// Handles only byte I/O; protocol details are managed by the protocol package.
//
// All operations accept context.Context for cancellation and timeout control.
// Implementations must be thread-safe.
type Client interface {
	// SendPacket Send transmits raw bytes to the remote endpoint.
	// Returns immediately after sending (does not wait for response).
	// Returns context.Canceled if ctx is cancelled, context.DeadlineExceeded if ctx times out.
	SendPacket(ctx context.Context, data []byte) error

	// RecvPacket Receive reads raw bytes from the remote endpoint.
	// Blocks until data is received or context is cancelled/times out.
	// Returns context.Canceled if ctx is cancelled, context.DeadlineExceeded if ctx times out.
	RecvPacket(ctx context.Context) ([]byte, error)

	// Close closes the transport and releases resources (e.g., sockets, file descriptors).
	// Multiple calls to Close should be idempotent and safe.
	Close() error
}

// RecvPacketWithAddr allows callers to capture the remote address for a received packet.
// Implementations may return an empty string when the remote address is unavailable.
type RecvPacketWithAddr interface {
	RecvPacketWithAddr(ctx context.Context) ([]byte, string, error)
}

// ConnectionInfo provides diagnostic information about the connection.
// Optional interface that transports can implement for better error messages.
type ConnectionInfo interface {
	// LocalAddr returns the local bind address (e.g., "10.100.0.16:41136")
	LocalAddr() string
	// RemoteAddr returns the remote drive address (e.g., "10.8.7.234:49360")
	RemoteAddr() string
}
