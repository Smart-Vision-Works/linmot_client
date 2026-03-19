package client

import (
	"context"
	"fmt"
	"strings"
)

// DriveFaultError represents a LinMot drive fault with captured evidence from status/probe calls.
type DriveFaultError struct {
	StatusWord  uint16
	StateVar    uint16
	ErrorCode   uint16
	ErrorText   string
	WarningWord uint16 // LinMot Warn Word bitfield from status (warning flags)
	WarningText string
}

func (e *DriveFaultError) Error() string {
	parts := []string{
		fmt.Sprintf("linmot drive fault status_word=0x%04X state_var=0x%04X", e.StatusWord, e.StateVar),
	}
	if e.ErrorCode != 0 {
		errPart := fmt.Sprintf("error_code=0x%04X", e.ErrorCode)
		if e.ErrorText != "" {
			errPart = fmt.Sprintf("%s error_text=%q", errPart, e.ErrorText)
		}
		parts = append(parts, errPart)
	}
	if e.WarningWord != 0 {
		warnPart := fmt.Sprintf("warning_word=0x%04X", e.WarningWord)
		if e.WarningText != "" {
			warnPart = fmt.Sprintf("%s warning_text=%q", warnPart, e.WarningText)
		}
		parts = append(parts, warnPart)
	}
	return strings.Join(parts, " ")
}

// DriveFaultProbeError indicates a fault was detected but probing for details failed.
type DriveFaultProbeError struct {
	Fault    *DriveFaultError
	ProbeErr error
}

func (e *DriveFaultProbeError) Error() string {
	if e.Fault == nil {
		return fmt.Sprintf("linmot drive fault probe failed: %v", e.ProbeErr)
	}
	return fmt.Sprintf("%s (probe failed: %v)", e.Fault.Error(), e.ProbeErr)
}

// Unwrap returns the underlying DriveFaultError so callers can extract evidence with errors.As.
func (e *DriveFaultProbeError) Unwrap() error {
	if e.Fault != nil {
		return e.Fault
	}
	return e.ProbeErr
}

// CheckDriveFault polls drive status and, if a fault is present, probes for error details.
// Returns nil when no fault exists. When a fault is detected, returns DriveFaultError with
// evidence; if probing the details fails, returns DriveFaultProbeError to indicate the attempt.
func (client *Client) CheckDriveFault(ctx context.Context) error {
	status, err := client.GetStatus(ctx)
	if err != nil {
		return err
	}
	if status == nil {
		return fmt.Errorf("linmot status poll returned nil")
	}

	if status.ErrorCode == 0 && status.WarnWord == 0 {
		return nil
	}

	fault := &DriveFaultError{
		StatusWord:  status.StatusWord,
		StateVar:    status.StateVar,
		ErrorCode:   status.ErrorCode,
		WarningWord: status.WarnWord,
	}

	if status.ErrorCode != 0 {
		errorText, probeErr := client.rtcManager.GetErrorText(ctx, status.ErrorCode)
		if probeErr != nil {
			return &DriveFaultProbeError{Fault: fault, ProbeErr: probeErr}
		}
		fault.ErrorText = errorText
	}

	return fault
}
