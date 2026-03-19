package transport

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// sharedClientTransport is a per-client transport that wraps a SharedUDPTransport.
// It implements the transport.Client interface and provides isolated send/receive
// for a specific drive IP/port while sharing the underlying UDP socket.
type sharedClientTransport struct {
	shared    *SharedUDPTransport
	driveIP   string
	drivePort int
	driveAddr *net.UDPAddr
	recvChan  chan []byte
	closeOnce sync.Once
}

func (c *sharedClientTransport) closeRecvChan() {
	c.closeOnce.Do(func() {
		close(c.recvChan)
	})
}

// SendPacket transmits raw bytes to the drive over the shared UDP socket.
func (c *sharedClientTransport) SendPacket(ctx context.Context, data []byte) error {
	// Check context before starting
	if err := ctx.Err(); err != nil {
		return err
	}

	// Check if shared transport is closed
	if c.shared.closed.Load() {
		return errors.New("shared transport is closed")
	}

	// Set deadline from context if present, otherwise use timeout.
	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		deadline = time.Now().Add(c.shared.timeout)
	}

	// Run send in goroutine to support context cancellation
	type errorOpt struct {
		err error
	}
	resultChannel := make(chan errorOpt, 1)

	go func() {
		// Serialize write deadline + write because they share mutable conn state.
		c.shared.writeMu.Lock()
		defer c.shared.writeMu.Unlock()

		if err := c.shared.conn.SetWriteDeadline(deadline); err != nil {
			resultChannel <- errorOpt{err: errors.WithMessage(err, "failed to set write deadline")}
			return
		}

		// Trace TX packet before sending
		localAddr := c.shared.conn.LocalAddr().String()
		remoteAddr := c.driveAddr.String()
		TracePacket("TX", localAddr, remoteAddr, data)

		_, err := c.shared.conn.WriteToUDP(data, c.driveAddr)
		if err != nil {
			resultChannel <- errorOpt{err: errors.WithMessage(err, "failed to send packet")}
		} else {
			resultChannel <- errorOpt{}
		}
	}()

	// Wait for either result or context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	case response := <-resultChannel:
		return response.err
	}
}

// RecvPacket reads raw bytes from the drive via the shared receive channel.
// Blocks until data is received or context is cancelled/times out.
func (c *sharedClientTransport) RecvPacket(ctx context.Context) ([]byte, error) {
	// Check context before starting
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Determine deadline
	var deadline time.Time
	var hasDeadline bool
	if deadline, hasDeadline = ctx.Deadline(); !hasDeadline {
		deadline = time.Now().Add(c.shared.timeout)
	}

	// Wait for data or timeout
	timer := time.NewTimer(time.Until(deadline))
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
		return nil, context.DeadlineExceeded
	case data, ok := <-c.recvChan:
		if !ok {
			return nil, errors.New("receive channel closed")
		}
		return data, nil
	}
}

// Close unregisters this client from the shared transport.
// Does NOT close the underlying shared UDP socket.
func (c *sharedClientTransport) Close() error {
	c.shared.UnregisterClient(c.driveIP)
	return nil
}

// LocalAddr implements transport.ConnectionInfo for diagnostic error messages.
func (c *sharedClientTransport) LocalAddr() string {
	if c.shared.conn != nil {
		return c.shared.conn.LocalAddr().String()
	}
	return ""
}

// RemoteAddr implements transport.ConnectionInfo for diagnostic error messages.
func (c *sharedClientTransport) RemoteAddr() string {
	return fmt.Sprintf("%s:%d", c.driveIP, c.drivePort)
}
