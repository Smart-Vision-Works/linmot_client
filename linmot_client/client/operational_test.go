package client

import (
	"strings"
	"testing"
	"time"

	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
)

func TestAssertDriveOperational_FirmwareStopped(t *testing.T) {
	status := &protocol_common.Status{
		StatusWord: 0x0000,
		StateVar:   0x0000, // StateVarHigh=0 -> NotReadyToSwitchOn
		ErrorCode:  0x0000,
		WarnWord:   0x0000,
	}

	err := assertDriveOperational(status)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Firmware Stopped") {
		t.Fatalf("expected Firmware Stopped in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "NotReadyToSwitchOn(0)") {
		t.Fatalf("expected state machine name in error, got: %v", err)
	}
}

func TestAssertDriveOperational_MotorSwitchedOff(t *testing.T) {
	status := &protocol_common.Status{
		StatusWord: 0x0001, // OperationEnable only, switch on active bit not set
		StateVar:   0x0600, // StateVarHigh=6 -> ReadyToOperate
		ErrorCode:  0x0000,
		WarnWord:   0x0000,
	}

	err := assertDriveOperational(status)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Motor Switched Off") {
		t.Fatalf("expected Motor Switched Off in error, got: %v", err)
	}
}

func TestAssertDriveOperational_OK(t *testing.T) {
	status := &protocol_common.Status{
		StatusWord: 0x0003, // OperationEnable + SwitchOnActive
		StateVar:   0x0800, // StateVarHigh=8 -> OperationEnabled
		ErrorCode:  0x0000,
		WarnWord:   0x0000,
	}

	if err := assertDriveOperational(status); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateConnectivity_FailFast(t *testing.T) {
	client, transportServer := NewMockClient()
	defer client.Close()

	drive := NewMockLinMot(transportServer)
	drive.Start()
	defer drive.Close()

	drive.SetStatus(&protocol_common.Status{
		StatusWord: 0x0001,
		StateVar:   0x0600,
		ErrorCode:  0x0000,
		WarnWord:   0x0000,
	})

	err := client.validateConnectivity(500 * time.Millisecond)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Motor Switched Off") {
		t.Fatalf("expected Motor Switched Off in error, got: %v", err)
	}
}

func TestValidateConnectivity_OK(t *testing.T) {
	client, transportServer := NewMockClient()
	defer client.Close()

	drive := NewMockLinMot(transportServer)
	drive.Start()
	defer drive.Close()

	drive.SetStatus(&protocol_common.Status{
		StatusWord: 0x0003,
		StateVar:   0x0800,
		ErrorCode:  0x0000,
		WarnWord:   0x0000,
	})

	if err := client.validateConnectivity(500 * time.Millisecond); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
