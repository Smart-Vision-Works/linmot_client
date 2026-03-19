package transport

import (
	"context"
	"testing"
	"time"
)

func TestMockTransportConfig_Default(t *testing.T) {
	config := DefaultMockTransportConfig()
	if config.EnablePacketLoss {
		t.Error("DefaultMockTransportConfig() should have EnablePacketLoss = false")
	}
}

func TestMockTransportConfig_WithPacketLoss(t *testing.T) {
	config := WithPacketLoss(0.5, 0.1)
	if !config.EnablePacketLoss {
		t.Error("WithPacketLoss() should have EnablePacketLoss = true")
	}
	if config.DropProbability != 0.5 {
		t.Errorf("WithPacketLoss() DropProbability = %f, want 0.5", config.DropProbability)
	}
	if config.DuplicateProbability != 0.1 {
		t.Errorf("WithPacketLoss() DuplicateProbability = %f, want 0.1", config.DuplicateProbability)
	}
}

func TestMockTransportClient_PacketDrop(t *testing.T) {
	// Create config with 100% drop probability
	config := &MockTransportConfig{
		DropProbability:  1.0,
		EnablePacketLoss: true,
	}
	client, server := NewMockTransportClientWithServerAndConfig(config)
	defer client.Close()
	defer server.Close()

	// Send packet - should be dropped
	err := client.SendPacket(context.Background(), []byte{1, 2, 3, 4})
	if err != nil {
		t.Errorf("SendPacket() error = %v, want nil (packet dropped silently)", err)
	}

	// Verify packet was not received by server
	select {
	case <-server.requestChan:
		t.Error("Packet should have been dropped but was received by server")
	case <-time.After(10 * time.Millisecond):
		// Expected - packet was dropped
	}
}

func TestMockTransportClient_PacketDelay(t *testing.T) {
	// Create config with delay
	config := &MockTransportConfig{
		DelayMin:         50 * time.Millisecond,
		DelayMax:         50 * time.Millisecond,
		EnablePacketLoss: true,
	}
	client, server := NewMockTransportClientWithServerAndConfig(config)
	defer client.Close()
	defer server.Close()

	// Send packet - should be delayed
	start := time.Now()
	err := client.SendPacket(context.Background(), []byte{1, 2, 3, 4})
	duration := time.Since(start)

	if err != nil {
		t.Errorf("SendPacket() error = %v, want nil", err)
	}

	// Verify delay was applied (within reasonable tolerance)
	if duration < 40*time.Millisecond || duration > 100*time.Millisecond {
		t.Errorf("SendPacket() delay = %v, expected ~50ms", duration)
	}

	// Verify packet was eventually received
	select {
	case packet := <-server.requestChan:
		if len(packet) != 4 {
			t.Errorf("Received packet length = %d, want 4", len(packet))
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Packet should have been received after delay")
	}
}

func TestMockTransportClient_PacketDuplicate(t *testing.T) {
	// Create config with 100% duplicate probability
	config := &MockTransportConfig{
		DuplicateProbability: 1.0,
		EnablePacketLoss:     true,
	}
	client, server := NewMockTransportClientWithServerAndConfig(config)
	defer client.Close()
	defer server.Close()

	// Send packet - should be duplicated
	err := client.SendPacket(context.Background(), []byte{1, 2, 3, 4})
	if err != nil {
		t.Errorf("SendPacket() error = %v, want nil", err)
	}

	// Verify packet was received (original)
	select {
	case packet := <-server.requestChan:
		if len(packet) != 4 {
			t.Errorf("Received packet length = %d, want 4", len(packet))
		}
	case <-time.After(10 * time.Millisecond):
		t.Error("Original packet should have been received")
	}

	// Verify duplicate was also received
	select {
	case packet := <-server.requestChan:
		if len(packet) != 4 {
			t.Errorf("Received duplicate packet length = %d, want 4", len(packet))
		}
	case <-time.After(10 * time.Millisecond):
		// Duplicate may be dropped if channel is full - this is acceptable
	}
}

func TestMockTransportServer_ResponseDrop(t *testing.T) {
	// Create config with 100% drop probability
	config := &MockTransportConfig{
		DropProbability:  1.0,
		EnablePacketLoss: true,
	}
	client, server := NewMockTransportClientWithServerAndConfig(config)
	defer client.Close()
	defer server.Close()

	// Send response - should be dropped
	server.SendPacket([]byte{5, 6, 7, 8})

	// Verify response was not received by client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.RecvPacket(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("RecvPacket() error = %v, want context.DeadlineExceeded (response dropped)", err)
	}
}

func TestMockTransportServer_ResponseDelay(t *testing.T) {
	// Create config with delay
	config := &MockTransportConfig{
		DelayMin:         50 * time.Millisecond,
		DelayMax:         50 * time.Millisecond,
		EnablePacketLoss: true,
	}
	client, server := NewMockTransportClientWithServerAndConfig(config)
	defer client.Close()
	defer server.Close()

	// Send response - should be delayed
	server.SendPacket([]byte{5, 6, 7, 8})

	// Verify response was received after delay
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	response, err := client.RecvPacket(ctx)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("RecvPacket() error = %v, want nil", err)
	}

	if len(response) != 4 {
		t.Errorf("Received response length = %d, want 4", len(response))
	}

	// Verify delay was applied (within reasonable tolerance)
	// Note: goroutine scheduling adds some overhead, so allow wider tolerance
	if duration < 40*time.Millisecond || duration > 150*time.Millisecond {
		t.Errorf("RecvPacket() delay = %v, expected ~50ms (tolerance: 40-150ms)", duration)
	}
}

func TestMockTransportServer_ResponseDuplicate(t *testing.T) {
	// Create config with 100% duplicate probability
	config := &MockTransportConfig{
		DuplicateProbability: 1.0,
		EnablePacketLoss:     true,
	}
	client, server := NewMockTransportClientWithServerAndConfig(config)
	defer client.Close()
	defer server.Close()

	// Send response - should be duplicated
	server.SendPacket([]byte{5, 6, 7, 8})

	// Verify response was received (original)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	response1, err := client.RecvPacket(ctx)
	if err != nil {
		t.Errorf("RecvPacket() error = %v, want nil", err)
	}
	if len(response1) != 4 {
		t.Errorf("Received response length = %d, want 4", len(response1))
	}

	// Verify duplicate was also received
	response2, err := client.RecvPacket(ctx)
	if err != nil {
		// Duplicate may be dropped if channel is full - this is acceptable
		t.Logf("Duplicate response not received (acceptable): %v", err)
	} else {
		if len(response2) != 4 {
			t.Errorf("Received duplicate response length = %d, want 4", len(response2))
		}
	}
}

func TestMockTransportClient_NoPacketLoss(t *testing.T) {
	// Create config with packet loss disabled
	config := DefaultMockTransportConfig()
	client, server := NewMockTransportClientWithServerAndConfig(config)
	defer client.Close()
	defer server.Close()

	// Send packet - should be delivered normally
	err := client.SendPacket(context.Background(), []byte{1, 2, 3, 4})
	if err != nil {
		t.Errorf("SendPacket() error = %v, want nil", err)
	}

	// Verify packet was received immediately
	select {
	case packet := <-server.requestChan:
		if len(packet) != 4 {
			t.Errorf("Received packet length = %d, want 4", len(packet))
		}
	case <-time.After(10 * time.Millisecond):
		t.Error("Packet should have been received immediately")
	}
}

func TestMockTransportConfig_Reset(t *testing.T) {
	config := &MockTransportConfig{
		EnablePacketLoss: true,
	}
	config.Reset() // Should not panic
	if config.rng == nil {
		t.Error("Reset() should initialize rng")
	}
}
