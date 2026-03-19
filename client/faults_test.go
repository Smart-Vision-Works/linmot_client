package client

import (
	"context"
	"errors"
	"testing"
	"time"

	protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"
)

func newFaultTestClient(t *testing.T) (*Client, *MockLinMot, func()) {
	t.Helper()

	client, transportServer := NewMockClient()
	drive := NewMockLinMot(transportServer)
	drive.Start()

	cleanup := func() {
		drive.Close()
		client.Close()
	}

	return client, drive, cleanup
}

func TestCheckDriveFault_NoFault(t *testing.T) {
	client, drive, cleanup := newFaultTestClient(t)
	defer cleanup()

	drive.SetStatus(&protocol_common.Status{
		StatusWord: 0x0001,
		StateVar:   0x0200,
		ErrorCode:  0,
		WarnWord:   0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := client.CheckDriveFault(ctx); err != nil {
		t.Fatalf("CheckDriveFault() expected nil, got %v", err)
	}
}

func TestCheckDriveFault_WithFault(t *testing.T) {
	client, drive, cleanup := newFaultTestClient(t)
	defer cleanup()

	drive.SetSimulateError(true, 0x0020)
	drive.SetSimulateWarning(true, 0x0001)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := client.CheckDriveFault(ctx)
	if err == nil {
		t.Fatal("CheckDriveFault() expected fault error, got nil")
	}

	var fault *DriveFaultError
	if !errors.As(err, &fault) {
		t.Fatalf("expected DriveFaultError, got %T", err)
	}

	if fault.ErrorCode != 0x0020 {
		t.Fatalf("ErrorCode = 0x%04X, want 0x0020", fault.ErrorCode)
	}
	if fault.ErrorText != "Position Lag Error" {
		t.Fatalf("ErrorText = %q, want %q", fault.ErrorText, "Position Lag Error")
	}
	if fault.WarningWord != 0x0001 {
		t.Fatalf("WarningWord = 0x%04X, want 0x0001", fault.WarningWord)
	}
}

func TestCheckDriveFault_ProbeFailure(t *testing.T) {
	client, drive, cleanup := newFaultTestClient(t)
	defer cleanup()

	drive.SetSimulateError(true, 0x0020)
	drive.SetErrorTextDelay(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := client.CheckDriveFault(ctx)
	if err == nil {
		t.Fatal("CheckDriveFault() expected probe failure, got nil")
	}

	var probeErr *DriveFaultProbeError
	if !errors.As(err, &probeErr) {
		t.Fatalf("expected DriveFaultProbeError, got %T", err)
	}
	if probeErr.ProbeErr == nil {
		t.Fatalf("expected ProbeErr to be set, got nil")
	}

	var fault *DriveFaultError
	if !errors.As(err, &fault) {
		t.Fatalf("expected to unwrap DriveFaultError, got %T", err)
	}
	if fault.ErrorCode != 0x0020 {
		t.Fatalf("ErrorCode = 0x%04X, want 0x0020", fault.ErrorCode)
	}
}
