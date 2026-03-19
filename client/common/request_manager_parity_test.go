package client_common

import (
	"context"
	"testing"
	"time"

	protocol_rtc "gsail-go/linmot/protocol/rtc"
	transport "gsail-go/linmot/transport"
)

// TestRequestManager_CmdCountNotConsumedOnSendFailure verifies that RTC cmdCount is consumed
// even if send fails, matching C# linudp.cs behavior (counter incremented before send,
// not rolled back on send failure).
func TestRequestManager_CmdCountNotConsumedOnSendFailure(t *testing.T) {
	// Create a transport that fails on first send
	baseClient, _ := transport.NewMockTransportClientWithServer()
	failingClient := &failingTransport{
		Client:        baseClient,
		failCount:     1, // Fail first send, succeed subsequent
		failRemaining: 1,
	}

	// Create RequestManager
	rm := NewRequestManager(failingClient, 10*time.Millisecond)
	rm.Start()
	defer rm.Stop()

	// Get initial counter state by peeking (Allocate doesn't consume, just returns current)
	rm.mu.RLock()
	initialCounter := rm.rtcCount.Allocate() // Peek without consuming
	rm.mu.RUnlock()

	// Submit RTC request (will fail on send)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create an RTC request (e.g., GetParam) - UPID and cmdCode
	req := protocol_rtc.NewRTCGetParamRequest(0x1234, 0x10) // UPID=0x1234, cmdCode=0x10

	_, err := SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](rm, ctx, req)

	// Verify request failed with send error
	if err == nil {
		t.Fatal("Expected send error, got nil")
	}

	// Verify counter was consumed (not rolled back on send failure)
	// The counter should have advanced by 1 even though send failed
	// Check by peeking again - should be initialCounter + 1
	rm.mu.RLock()
	nextCounter := rm.rtcCount.Allocate() // Peek again
	rm.mu.RUnlock()

	expectedNext := initialCounter + 1
	if expectedNext > 14 {
		expectedNext = 1
	}

	if nextCounter != expectedNext {
		t.Errorf("Expected counter to advance to %d (consumed on send failure), got %d. Counter was NOT consumed on send failure (mismatch with C#)", expectedNext, nextCounter)
	}

	t.Logf("Counter consumed on send failure: %d -> %d (matches C#: counter not rolled back on send failure)", initialCounter, nextCounter)
}

// failingTransport wraps a transport.Client and fails the first N sends
type failingTransport struct {
	transport.Client
	failCount     int
	failRemaining int
}

func (ft *failingTransport) SendPacket(ctx context.Context, data []byte) error {
	if ft.failRemaining > 0 {
		ft.failRemaining--
		return context.DeadlineExceeded // Simulate send failure
	}
	return ft.Client.SendPacket(ctx, data)
}
