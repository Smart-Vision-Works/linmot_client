package transport

import (
	"sync"
	"sync/atomic"
	"time"
)

// mockTransportServer implements TransportServer using in-memory Go channels.
// It's purely a transport layer - just exposes channels for sending/receiving packets.
// It contains NO drive simulation logic - that belongs in mockLinMot.
type mockTransportServer struct {
	requestChan  chan []byte
	responseChan chan []byte
	mu           sync.RWMutex // Protects channel operations during close
	closed       atomic.Bool
	closeOnce    sync.Once
	config       *MockTransportConfig // Optional config for packet loss simulation
}

// NewMockTransportServer creates a new mock transport server (transport layer only).
func NewMockTransportServer() *mockTransportServer {
	return NewMockTransportServerWithConfig(DefaultMockTransportConfig())
}

// NewMockTransportServerWithConfig creates a new mock transport server with packet loss simulation config.
func NewMockTransportServerWithConfig(config *MockTransportConfig) *mockTransportServer {
	if config == nil {
		config = DefaultMockTransportConfig()
	}
	return &mockTransportServer{
		requestChan:  make(chan []byte, 10), // Larger buffer for duplicate packets
		responseChan: make(chan []byte, 10), // Larger buffer for duplicate packets
		config:       config,
	}
}

// ReceiveRequest receives a request packet from the client.
// This is called by mockLinMot to get incoming requests.
// Returns nil if the transport is closed.
func (s *mockTransportServer) RecvPacket() []byte {
	if s.closed.Load() {
		return nil
	}
	packet, ok := <-s.requestChan
	if !ok {
		return nil
	}
	return packet
}

// SendResponse sends a response packet to the client.
// This is called by mockLinMot to send responses.
// Blocks until the client receives the response.
// No-op if the transport is closed.
func (s *mockTransportServer) SendPacket(response []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed.Load() {
		return
	}

	// Apply packet loss simulation
	if s.config != nil && s.config.EnablePacketLoss {
		// Check if response should be dropped
		if s.config.shouldDrop() {
			// Response dropped - simulate packet loss
			return
		}

		// Apply delay if configured (in a goroutine to avoid blocking)
		if delay := s.config.randomDelay(); delay > 0 {
			go func() {
				time.Sleep(delay)
				s.mu.RLock()
				defer s.mu.RUnlock()
				if !s.closed.Load() {
					select {
					case s.responseChan <- response:
					default:
						// Channel full, response dropped (acceptable for simulation)
					}
				}
			}()
			return
		}

		// Check if response should be duplicated
		if s.config.shouldDuplicate() {
			// Send original response
			select {
			case s.responseChan <- response:
				// Original sent
			default:
				// Channel full, original dropped
			}

			// Send duplicate (non-blocking, best effort)
			select {
			case s.responseChan <- response:
				// Duplicate sent
			default:
				// Channel full, duplicate dropped (acceptable for simulation)
			}
			return
		}
	}

	// Block until client receives response - don't drop responses!
	// The RLock ensures this won't race with Close()
	s.responseChan <- response
}

// Close closes the transport channels.
// Safe to call multiple times (idempotent).
func (s *mockTransportServer) Close() {
	s.closeOnce.Do(func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.closed.Store(true)
		close(s.requestChan)
		close(s.responseChan)
	})
}
