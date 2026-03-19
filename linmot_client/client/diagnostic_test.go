package client

import (
	"context"
	"errors"
	"testing"
	"time"

	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	transport "github.com/Smart-Vision-Works/staged_robot/transport"
)

// newUDPClientNoValidation creates a UDP client WITHOUT connectivity validation.
// This is for diagnostic purposes only - to query drives in error state.
func newUDPClientNoValidation(driveIP string, drivePort, masterPort int, timeout time.Duration, debug bool) (*Client, error) {
	transportClient, err := transport.NewUDPTransportClient(driveIP, drivePort, masterPort, "", timeout)
	if err != nil {
		return nil, err
	}
	client := newClient(transportClient)
	client.driveIP = driveIP
	if debug {
		client.SetDebug(true)
	}
	return client, nil
}

// TestDiagnostic_RawStatus gets raw status from a LinMot without operational checks.
// Use: go test -v -run TestDiagnostic_RawStatus -linmot_mode=udp -linmot_ip=10.8.7.234 -linmot_timeout=30s
func TestDiagnostic_RawStatus(t *testing.T) {
	if *linmotMode != "udp" {
		t.Skip("Skipping diagnostic in mock mode - requires real hardware")
	}

	t.Logf("Connecting to LinMot at %s (no validation)...", *linmotIP)

	// Create client WITHOUT operational validation - for diagnostic use
	client, err := newUDPClientNoValidation(*linmotIP, *linmotDrivePort, *linmotMasterPort, *linmotTimeout, *linmotDebug)
	if err != nil {
		t.Fatalf("Failed to create UDP client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get raw status
	status, err := client.GetStatus(ctx)
	if err != nil {
		// Check if it's a timeout error and print diagnostic hint
		var timeoutErr *protocol_common.RequestTimeoutError
		if errors.As(err, &timeoutErr) {
			t.Logf("Timeout error details: %v", err)
			t.Logf("\n%s", timeoutErr.DiagnosticHint())
		}
		t.Fatalf("GetStatus failed: %v", err)
	}

	mainState := (status.StateVar >> 8) & 0xFF
	subState := status.StateVar & 0xFF

	t.Logf("\n=== LinMot Status at %s ===", *linmotIP)
	t.Logf("StatusWord:     0x%04X", status.StatusWord)
	t.Logf("StateVar:       0x%04X (MainState=%d, SubState=%d)", status.StateVar, mainState, subState)
	t.Logf("ErrorCode:      0x%04X", status.ErrorCode)
	t.Logf("WarnWord:       0x%04X", status.WarnWord)
	t.Logf("ActualPosition: %d (%.4f mm)", status.ActualPosition, status.ActualPositionMM())
	t.Logf("DemandPosition: %d (%.4f mm)", status.DemandPosition, status.DemandPositionMM())
	t.Logf("Current:        %d mA", status.Current)

	// State machine interpretation
	stateNames := map[uint16]string{
		0: "NotReadyToSwitchOn",
		1: "SwitchOnDisabled",
		2: "ReadyToSwitchOn",
		3: "SwitchedOn",
		4: "GeneralError",
		5: "SwitchOnInhibit",
		6: "ReadyToOperate",
		7: "OperationNotEnabled",
		8: "OperationEnabled",
		9: "Homing",
	}
	if name, ok := stateNames[mainState]; ok {
		t.Logf("State Machine: %s (%d)", name, mainState)
	} else {
		t.Logf("State Machine: Unknown (%d)", mainState)
	}

	// Get error text if there's an error
	if status.ErrorCode != 0 {
		errText, err := client.GetErrorText(ctx, status.ErrorCode)
		if err != nil {
			t.Logf("GetErrorText failed: %v", err)
		} else {
			t.Logf("\n!!! ERROR ACTIVE: Code=0x%04X Text=%q", status.ErrorCode, errText)
		}
	}

	// Get error log counts
	logged, occurred, err := client.GetErrorLogCounts(ctx)
	if err != nil {
		t.Logf("GetErrorLogCounts failed: %v", err)
	} else {
		t.Logf("\nError Log: %d logged, %d occurred", logged, occurred)

		// Get recent error log entries
		if logged > 0 {
			t.Logf("\n=== Recent Error Log Entries ===")
			maxEntries := logged
			if maxEntries > 10 {
				maxEntries = 10
			}
			for i := uint16(0); i < maxEntries; i++ {
				entry, err := client.GetErrorLogEntry(ctx, i)
				if err != nil {
					t.Logf("[%d] GetErrorLogEntry failed: %v", i, err)
					continue
				}
				errText, _ := client.GetErrorText(ctx, entry.ErrorCode)
				t.Logf("[%d] Code=0x%04X Time=%s Text=%q", i, entry.ErrorCode, entry.Timestamp.Format(time.RFC3339), errText)
			}
		}
	}
}

// TestDiagnostic_TryAcknowledgeError attempts to acknowledge any active error.
// Use: go test -v -run TestDiagnostic_TryAcknowledgeError -linmot_mode=udp -linmot_ip=10.8.7.234
func TestDiagnostic_TryAcknowledgeError(t *testing.T) {
	if *linmotMode != "udp" {
		t.Skip("Skipping diagnostic in mock mode - requires real hardware")
	}

	t.Logf("Connecting to LinMot at %s (no validation)...", *linmotIP)

	client, err := newUDPClientNoValidation(*linmotIP, *linmotDrivePort, *linmotMasterPort, *linmotTimeout, *linmotDebug)
	if err != nil {
		t.Fatalf("Failed to create UDP client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get initial status
	status, err := client.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	t.Logf("Before acknowledge: ErrorCode=0x%04X StateVar=0x%04X", status.ErrorCode, status.StateVar)

	if status.ErrorCode == 0 {
		t.Log("No error active - nothing to acknowledge")
		return
	}

	// Try to acknowledge the error
	t.Log("Attempting to acknowledge error...")
	status, err = client.AcknowledgeError(ctx)
	if err != nil {
		t.Logf("AcknowledgeError failed: %v", err)
	} else {
		t.Logf("After acknowledge: ErrorCode=0x%04X StateVar=0x%04X", status.ErrorCode, status.StateVar)

		mainState := (status.StateVar >> 8) & 0xFF
		if status.ErrorCode == 0 && mainState != 4 {
			t.Log("SUCCESS: Error acknowledged and cleared!")
		} else {
			t.Log("Error still active after acknowledge attempt")
		}
	}
}

// TestDiagnostic_ReadCommandTable reads the current command table from the LinMot.
// Use: go test -v -run TestDiagnostic_ReadCommandTable -linmot_mode=udp -linmot_ip=10.8.7.234 -linmot_timeout=60s
func TestDiagnostic_ReadCommandTable(t *testing.T) {
	if *linmotMode != "udp" {
		t.Skip("Skipping diagnostic in mock mode - requires real hardware")
	}

	t.Logf("Connecting to LinMot at %s (no validation)...", *linmotIP)

	client, err := newUDPClientNoValidation(*linmotIP, *linmotDrivePort, *linmotMasterPort, *linmotTimeout, *linmotDebug)
	if err != nil {
		t.Fatalf("Failed to create UDP client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// First get status to understand drive state
	status, err := client.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	mainState := (status.StateVar >> 8) & 0xFF
	t.Logf("Drive state: MainState=%d ErrorCode=0x%04X", mainState, status.ErrorCode)

	// Note: Reading command table may require MC to be stopped
	// If MC is running, we may need to stop it first

	t.Log("Attempting to read command table...")
	t.Log("(Note: This requires MC to be stopped - may fail if MC is running)")

	// Stop MC first if needed
	if mainState == 8 { // OperationEnabled
		t.Log("MC appears to be running, stopping it first...")
		err = client.StopMotionController(ctx)
		if err != nil {
			t.Logf("StopMotionController failed: %v", err)
		} else {
			t.Log("MC stopped successfully")
		}
	}

	// Read command table
	ct, err := client.GetCommandTable(ctx)
	if err != nil {
		t.Fatalf("GetCommandTable failed: %v", err)
	}

	t.Logf("\n=== Command Table (%d entries) ===", len(ct.Entries))
	for _, entry := range ct.Entries {
		par1Str := "nil"
		if entry.Par1 != nil && entry.Par1.Literal != nil {
			par1Str = string(rune(*entry.Par1.Literal))
		}
		t.Logf("ID=%d Name=%q Type=%q Par1=%s", entry.ID, entry.Name, entry.Type, par1Str)
	}
}

// TestDiagnostic_RecoverySequence attempts a full recovery sequence on a LinMot in error state.
// Use: go test -v -run TestDiagnostic_RecoverySequence -linmot_mode=udp -linmot_ip=10.8.7.234 -linmot_timeout=60s
func TestDiagnostic_RecoverySequence(t *testing.T) {
	if *linmotMode != "udp" {
		t.Skip("Skipping diagnostic in mock mode - requires real hardware")
	}

	t.Logf("=== LinMot Recovery Sequence for %s ===", *linmotIP)

	client, err := newUDPClientNoValidation(*linmotIP, *linmotDrivePort, *linmotMasterPort, *linmotTimeout, true)
	if err != nil {
		t.Fatalf("Failed to create UDP client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Step 1: Get initial status
	t.Log("\n--- Step 1: Initial Status ---")
	status, err := client.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	mainState := (status.StateVar >> 8) & 0xFF
	t.Logf("State: MainState=%d ErrorCode=0x%04X WarnWord=0x%04X", mainState, status.ErrorCode, status.WarnWord)
	t.Logf("Position: Actual=%.2fmm Demand=%.2fmm", status.ActualPositionMM(), status.DemandPositionMM())

	if status.ErrorCode == 0 {
		t.Log("No error active - drive appears healthy")
		return
	}

	// Get error text
	errText, _ := client.GetErrorText(ctx, status.ErrorCode)
	t.Logf("Active Error: 0x%04X = %q", status.ErrorCode, errText)

	// Step 2: Try to enable the drive first (sometimes needed before acknowledge)
	t.Log("\n--- Step 2: Attempting EnableDrive ---")
	ctxEnable, cancelEnable := context.WithTimeout(ctx, 5*time.Second)
	status, err = client.EnableDrive(ctxEnable)
	cancelEnable()
	if err != nil {
		t.Logf("EnableDrive failed (expected in error state): %v", err)
	} else {
		mainState = (status.StateVar >> 8) & 0xFF
		t.Logf("After EnableDrive: MainState=%d ErrorCode=0x%04X", mainState, status.ErrorCode)
	}

	// Step 3: Re-check status
	t.Log("\n--- Step 3: Re-check Status ---")
	status, err = client.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	mainState = (status.StateVar >> 8) & 0xFF
	t.Logf("Current State: MainState=%d ErrorCode=0x%04X", mainState, status.ErrorCode)

	// Step 4: Try AcknowledgeError with explicit longer timeout
	if status.ErrorCode != 0 {
		t.Log("\n--- Step 4: Attempting AcknowledgeError ---")
		ctxAck, cancelAck := context.WithTimeout(ctx, 15*time.Second)
		status, err = client.AcknowledgeError(ctxAck)
		cancelAck()
		if err != nil {
			t.Logf("AcknowledgeError failed: %v", err)
		} else {
			mainState = (status.StateVar >> 8) & 0xFF
			t.Logf("After AcknowledgeError: MainState=%d ErrorCode=0x%04X", mainState, status.ErrorCode)
		}
	}

	// Step 5: Final status check
	t.Log("\n--- Step 5: Final Status ---")
	status, err = client.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	mainState = (status.StateVar >> 8) & 0xFF
	t.Logf("Final State: MainState=%d ErrorCode=0x%04X WarnWord=0x%04X", mainState, status.ErrorCode, status.WarnWord)

	// Summary
	t.Log("\n--- Summary ---")
	if status.ErrorCode == 0 && mainState != 4 {
		t.Log("SUCCESS: Error cleared, drive recovered!")
	} else if status.ErrorCode == 0 {
		t.Log("PARTIAL: Error code cleared but still in error state")
	} else {
		t.Log("FAILED: Error still active")
		t.Log("This error may require manual intervention:")
		t.Log("  1. Check if motor is physically blocked")
		t.Log("  2. Power cycle the LinMot")
		t.Log("  3. Check motor sizing/load")
	}
}

// TestDiagnostic_CheckBothLinMots checks status of both known LinMot IPs.
// Use: go test -v -run TestDiagnostic_CheckBothLinMots
func TestDiagnostic_CheckBothLinMots(t *testing.T) {
	if *linmotMode != "udp" {
		t.Skip("Skipping diagnostic in mock mode - requires real hardware")
	}

	ips := []string{*linmotIP}

	for _, ip := range ips {
		t.Run(ip, func(t *testing.T) {
			t.Logf("Checking LinMot at %s (no validation)...", ip)

			client, err := newUDPClientNoValidation(ip, 49360, 41136, 5*time.Second, false)
			if err != nil {
				t.Logf("UNREACHABLE: %v", err)
				return
			}
			defer client.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			status, err := client.GetStatus(ctx)
			if err != nil {
				t.Logf("GetStatus FAILED: %v", err)
				return
			}

			mainState := (status.StateVar >> 8) & 0xFF
			stateNames := map[uint16]string{
				0: "NotReadyToSwitchOn", 1: "SwitchOnDisabled", 2: "ReadyToSwitchOn",
				3: "SwitchedOn", 4: "GeneralError", 5: "SwitchOnInhibit",
				6: "ReadyToOperate", 7: "OperationNotEnabled", 8: "OperationEnabled", 9: "Homing",
			}
			stateName := stateNames[mainState]

			t.Logf("ONLINE: State=%s(%d) Error=0x%04X Pos=%.2fmm",
				stateName, mainState, status.ErrorCode, status.ActualPositionMM())

			if status.ErrorCode != 0 {
				errText, _ := client.GetErrorText(ctx, status.ErrorCode)
				t.Logf("  ERROR: %q", errText)
			}
		})
	}
}
