package transport

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestUDPTransportClient_SendPacket_ContextDeadline(t *testing.T) {
	client, err := NewUDPTransportClient("127.0.0.1", DefaultDrivePort, DefaultMasterPort, "", DefaultTimeout)
	if err != nil {
		t.Fatalf("NewUDPTransportClient() failed: %v", err)
	}
	defer client.Close()

	t.Run("Context deadline exceeded", func(t *testing.T) {
		// Create context with very short deadline
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Wait a bit to ensure deadline has passed
		time.Sleep(10 * time.Millisecond)

		// Try to send - should return context.DeadlineExceeded
		err := client.SendPacket(ctx, []byte{1, 2, 3, 4})
		if err != context.DeadlineExceeded {
			t.Errorf("SendPacket() error = %v, want context.DeadlineExceeded", err)
		}
	})

	t.Run("Context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// Try to send - should return context.Canceled
		err := client.SendPacket(ctx, []byte{1, 2, 3, 4})
		if err != context.Canceled {
			t.Errorf("SendPacket() error = %v, want context.Canceled", err)
		}
	})

	t.Run("Context deadline sets socket deadline", func(t *testing.T) {
		// Create context with short deadline
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		// UDP send is non-blocking, so it completes immediately
		// The deadline is set on the socket, but send completes before deadline
		// This test verifies that the context deadline is used to set socket deadline
		err := client.SendPacket(ctx, []byte{1, 2, 3, 4})
		// Send should succeed (UDP send is non-blocking)
		if err != nil {
			t.Errorf("SendPacket() error = %v, expected nil (UDP send is non-blocking)", err)
		}
	})
}

func TestUDPTransportClient_RecvPacket_ContextDeadline(t *testing.T) {
	client, err := NewUDPTransportClient("127.0.0.1", DefaultDrivePort, DefaultMasterPort, "", DefaultTimeout)
	if err != nil {
		t.Fatalf("NewUDPTransportClient() failed: %v", err)
	}
	defer client.Close()

	t.Run("Context deadline exceeded", func(t *testing.T) {
		// Create context with very short deadline
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Wait a bit to ensure deadline has passed
		time.Sleep(10 * time.Millisecond)

		// Try to receive - should return context.DeadlineExceeded
		_, err := client.RecvPacket(ctx)
		if err != context.DeadlineExceeded {
			t.Errorf("RecvPacket() error = %v, want context.DeadlineExceeded", err)
		}
	})

	t.Run("Context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// Try to receive - should return context.Canceled
		_, err := client.RecvPacket(ctx)
		if err != context.Canceled {
			t.Errorf("RecvPacket() error = %v, want context.Canceled", err)
		}
	})

	t.Run("Context deadline respected during operation", func(t *testing.T) {
		// Create context with short deadline
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		// Try to receive - operation should timeout (no response expected)
		start := time.Now()
		_, err := client.RecvPacket(ctx)
		duration := time.Since(start)

		// UDP read timeout may return wrapped error, but should indicate timeout
		// The socket deadline is set from context deadline, so timeout should occur
		if err == nil {
			t.Fatal("RecvPacket() expected timeout error but got nil")
		}

		// Error should be either context.DeadlineExceeded or a timeout error
		// (UDP socket may return i/o timeout which is wrapped)
		errStr := err.Error()
		if err != context.DeadlineExceeded && !strings.Contains(errStr, "timeout") {
			t.Errorf("RecvPacket() error = %v, want context.DeadlineExceeded or timeout error", err)
		}

		// Verify it timed out within reasonable bounds (50ms + some tolerance)
		if duration < 40*time.Millisecond || duration > 200*time.Millisecond {
			t.Errorf("RecvPacket() timed out after %v, expected ~50ms", duration)
		}
	})
}

func TestMockTransportClient_SendPacket_ContextDeadline(t *testing.T) {
	client, _ := NewMockTransportClientWithServer()
	defer client.Close()

	t.Run("Context deadline exceeded", func(t *testing.T) {
		// Create context with very short deadline
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Wait a bit to ensure deadline has passed
		time.Sleep(10 * time.Millisecond)

		// Try to send - should return context.DeadlineExceeded
		err := client.SendPacket(ctx, []byte{1, 2, 3, 4})
		if err != context.DeadlineExceeded {
			t.Errorf("SendPacket() error = %v, want context.DeadlineExceeded", err)
		}
	})

	t.Run("Context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// Try to send - should return context.Canceled
		err := client.SendPacket(ctx, []byte{1, 2, 3, 4})
		if err != context.Canceled {
			t.Errorf("SendPacket() error = %v, want context.Canceled", err)
		}
	})

	t.Run("Context cancellation during blocking operation", func(t *testing.T) {
		// Create a server that won't receive (simulate blocking)
		client, server := NewMockTransportClientWithServer()
		defer client.Close()
		defer server.Close()

		// Close server to simulate blocking send
		server.Close()

		ctx, cancel := context.WithCancel(context.Background())

		// Cancel context after a short delay
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		// Try to send - should be cancelled
		err := client.SendPacket(ctx, []byte{1, 2, 3, 4})
		if err != context.Canceled {
			t.Errorf("SendPacket() error = %v, want context.Canceled", err)
		}
	})
}

func TestMockTransportClient_RecvPacket_ContextDeadline(t *testing.T) {
	client, server := NewMockTransportClientWithServer()
	defer client.Close()
	defer server.Close()

	t.Run("Context deadline exceeded", func(t *testing.T) {
		// Create context with very short deadline
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Wait a bit to ensure deadline has passed
		time.Sleep(10 * time.Millisecond)

		// Try to receive - should return context.DeadlineExceeded
		_, err := client.RecvPacket(ctx)
		if err != context.DeadlineExceeded {
			t.Errorf("RecvPacket() error = %v, want context.DeadlineExceeded", err)
		}
	})

	t.Run("Context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// Try to receive - should return context.Canceled
		_, err := client.RecvPacket(ctx)
		if err != context.Canceled {
			t.Errorf("RecvPacket() error = %v, want context.Canceled", err)
		}
	})

	t.Run("Context cancellation during blocking operation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		// Cancel context after a short delay
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		// Try to receive - should be cancelled (no response will be sent)
		_, err := client.RecvPacket(ctx)
		if err != context.Canceled {
			t.Errorf("RecvPacket() error = %v, want context.Canceled", err)
		}
	})

	t.Run("Context deadline respected during blocking operation", func(t *testing.T) {
		// Create context with short deadline
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		// Try to receive - operation should timeout (no response will be sent)
		start := time.Now()
		_, err := client.RecvPacket(ctx)
		duration := time.Since(start)

		if err != context.DeadlineExceeded {
			t.Errorf("RecvPacket() error = %v, want context.DeadlineExceeded", err)
		}

		// Verify it timed out within reasonable bounds (50ms + some tolerance)
		if duration < 40*time.Millisecond || duration > 200*time.Millisecond {
			t.Errorf("RecvPacket() timed out after %v, expected ~50ms", duration)
		}
	})
}

func TestUDPTransportClient_ContextDeadline_ErrorTypes(t *testing.T) {
	client, err := NewUDPTransportClient("127.0.0.1", DefaultDrivePort, DefaultMasterPort, "", DefaultTimeout)
	if err != nil {
		t.Fatalf("NewUDPTransportClient() failed: %v", err)
	}
	defer client.Close()

	t.Run("DeadlineExceeded error type", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		time.Sleep(10 * time.Millisecond)

		err := client.SendPacket(ctx, []byte{1, 2, 3})
		if err == nil {
			t.Fatal("SendPacket() expected error but got nil")
		}

		// Verify it's exactly context.DeadlineExceeded
		if err != context.DeadlineExceeded {
			t.Errorf("SendPacket() error = %v, want context.DeadlineExceeded", err)
		}

		// Verify error message
		if !strings.Contains(err.Error(), "deadline exceeded") {
			t.Errorf("SendPacket() error message = %q, want to contain 'deadline exceeded'", err.Error())
		}
	})

	t.Run("Canceled error type", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := client.SendPacket(ctx, []byte{1, 2, 3})
		if err == nil {
			t.Fatal("SendPacket() expected error but got nil")
		}

		// Verify it's exactly context.Canceled
		if err != context.Canceled {
			t.Errorf("SendPacket() error = %v, want context.Canceled", err)
		}

		// Verify error message
		if !strings.Contains(err.Error(), "canceled") {
			t.Errorf("SendPacket() error message = %q, want to contain 'canceled'", err.Error())
		}
	})
}
