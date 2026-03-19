package transport

import (
	"context"
	"time"
)

// mockTransportClient implements Client using in-memory Go channels.
// This is a pure I/O implementation with NO protocol knowledge.
// It simply sends and receives raw bytes over channels.
// Thread-safe: can be used concurrently from multiple goroutines.
type mockTransportClient struct {
	server *mockTransportServer // The transport server (other end of the channel pair)
	config *MockTransportConfig // Optional config for packet loss simulation
}

// NewMockTransportClientWithServer creates a connected pair: a mockTransportClient
// and a mockTransportServer. This is purely transport - no protocol logic here.
// The protocol logic lives in LinUDPV2Protocol, and drive logic lives in mockLinMot.
func NewMockTransportClientWithServer() (*mockTransportClient, *mockTransportServer) {
	return NewMockTransportClientWithServerAndConfig(DefaultMockTransportConfig())
}

// NewMockTransportClientWithServerAndConfig creates a connected pair with packet loss simulation config.
func NewMockTransportClientWithServerAndConfig(config *MockTransportConfig) (*mockTransportClient, *mockTransportServer) {
	if config == nil {
		config = DefaultMockTransportConfig()
	}
	server := NewMockTransportServerWithConfig(config)
	client := &mockTransportClient{
		server: server,
		config: config,
	}
	return client, server
}

// SendPacket Send transmits raw bytes to the server over the channel.
// This is a pure I/O operation with no knowledge of packet contents.
func (c *mockTransportClient) SendPacket(ctx context.Context, data []byte) error {
	c.server.mu.RLock()
	defer c.server.mu.RUnlock()

	// Check if closed
	if c.server.closed.Load() {
		return context.Canceled
	}

	// Check context before starting
	if err := ctx.Err(); err != nil {
		return err
	}

	// Apply packet loss simulation
	if c.config != nil && c.config.EnablePacketLoss {
		// Check if packet should be dropped
		if c.config.shouldDrop() {
			// Packet dropped - simulate packet loss
			return nil // Return success to caller (packet "sent" but lost)
		}

		// Apply delay if configured
		if delay := c.config.randomDelay(); delay > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
				// Delay completed
			}
		}

		// Check if packet should be duplicated
		if c.config.shouldDuplicate() {
			// Send original packet
			select {
			case <-ctx.Done():
				return ctx.Err()
			case c.server.requestChan <- data:
				// Original sent
			}

			// Send duplicate (non-blocking, best effort)
			select {
			case c.server.requestChan <- data:
				// Duplicate sent
			default:
				// Channel full, duplicate dropped (acceptable for simulation)
			}
			return nil
		}
	}

	// Send packet to server (may block if server is slow)
	// The RLock ensures this won't race with Close()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.server.requestChan <- data:
		return nil
	}
}

// RecvPacket Receive reads raw bytes from the server over the channel.
// This is a pure I/O operation with no knowledge of packet contents.
// Returns the received bytes or an error.
func (c *mockTransportClient) RecvPacket(ctx context.Context) ([]byte, error) {
	data, _, err := c.RecvPacketWithAddr(ctx)
	return data, err
}

// RecvPacketWithAddr Receive reads raw bytes from the server over the channel and returns a mock sender address.
func (c *mockTransportClient) RecvPacketWithAddr(ctx context.Context) ([]byte, string, error) {
	// Check context before starting
	if err := ctx.Err(); err != nil {
		return nil, "", err
	}

	// Receive response from server (blocks until available)
	select {
	case <-ctx.Done():
		return nil, "", ctx.Err()
	case response := <-c.server.responseChan:
		// Apply packet loss simulation on receive
		if c.config != nil && c.config.EnablePacketLoss {
			// Check if response should be dropped
			if c.config.shouldDrop() {
				// Response dropped - continue waiting (simulate packet loss)
				// Recursively call to wait for next response
				return c.RecvPacketWithAddr(ctx)
			}

			// Apply delay if configured
			if delay := c.config.randomDelay(); delay > 0 {
				select {
				case <-ctx.Done():
					return nil, "", ctx.Err()
				case <-time.After(delay):
					// Delay completed
				}
			}
		}

		return response, "mock", nil
	}
}

// Close closes the mock transport (closes server channels).
func (c *mockTransportClient) Close() error {
	if c.server != nil {
		c.server.Close()
	}
	return nil
}
