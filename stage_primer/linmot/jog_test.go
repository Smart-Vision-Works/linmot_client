package linmot

import (
	"context"
	"errors"
	"testing"
	"time"

	client_command_tables "github.com/Smart-Vision-Works/staged_robot/client/rtc/command_tables"
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_control_word "github.com/Smart-Vision-Works/staged_robot/protocol/control_word"

	config "stage_primer_config"
)

type jogTestClient struct {
	initialStatus *protocol_common.Status
	ackStatus     *protocol_common.Status
	enableStatus  *protocol_common.Status

	ackCalls       int
	enableCalls    int
	setPosition1   int
	setPosition2   int
	setRunModeCall int
}

func (c *jogTestClient) GetPosition(context.Context) (float64, error) { return 0, nil }
func (c *jogTestClient) GetStatus(context.Context) (*protocol_common.Status, error) {
	return c.initialStatus, nil
}
func (c *jogTestClient) CheckDriveFault(context.Context) error { return nil }
func (c *jogTestClient) EnableDrive(context.Context) (*protocol_common.Status, error) {
	c.enableCalls++
	return c.enableStatus, nil
}
func (c *jogTestClient) AcknowledgeError(context.Context) (*protocol_common.Status, error) {
	c.ackCalls++
	return c.ackStatus, nil
}
func (c *jogTestClient) SendControlWord(context.Context, uint16) (*protocol_common.Status, error) {
	return &protocol_common.Status{}, nil
}
func (c *jogTestClient) SetPosition1(context.Context, float64, protocol_common.ParameterStorageType) error {
	c.setPosition1++
	return nil
}
func (c *jogTestClient) SetPosition2(context.Context, float64, protocol_common.ParameterStorageType) error {
	c.setPosition2++
	return nil
}
func (c *jogTestClient) SetRunMode(context.Context, protocol_common.RunMode, protocol_common.ParameterStorageType) error {
	c.setRunModeCall++
	return nil
}
func (c *jogTestClient) SetCommandTable(context.Context, *client_command_tables.CommandTable) error {
	return nil
}
func (c *jogTestClient) WriteOutputsWithMask(context.Context, uint16, uint16) (*protocol_common.Status, error) {
	return &protocol_common.Status{}, nil
}
func (c *jogTestClient) VAIGoToPosition(context.Context, float64, float64, float64, float64) (*protocol_common.Status, error) {
	return &protocol_common.Status{}, nil
}
func (c *jogTestClient) VAIGoToPositionFromActual(context.Context, float64, float64, float64, float64) (*protocol_common.Status, error) {
	return &protocol_common.Status{}, nil
}
func (c *jogTestClient) SetEasyStepsAutoStart(context.Context, int32, protocol_common.ParameterStorageType) error {
	return nil
}
func (c *jogTestClient) SetEasyStepsAutoHome(context.Context, int32, protocol_common.ParameterStorageType) error {
	return nil
}
func (c *jogTestClient) SetEasyStepsInputRisingEdgeFunction(context.Context, protocol_common.IOPinNumber, int32, protocol_common.ParameterStorageType) error {
	return nil
}
func (c *jogTestClient) SetEasyStepsInputCurveCmdID(context.Context, protocol_common.IOPinNumber, int32, protocol_common.ParameterStorageType) error {
	return nil
}
func (c *jogTestClient) SetIODefOutputFunction(context.Context, protocol_common.IOPinNumber, int32, protocol_common.ParameterStorageType) error {
	return nil
}
func (c *jogTestClient) SetIODefInputFunction(context.Context, protocol_common.IOPinNumber, int32, protocol_common.ParameterStorageType) error {
	return nil
}
func (c *jogTestClient) SetTriggerMode(context.Context, int32, protocol_common.ParameterStorageType) error {
	return nil
}
func (c *jogTestClient) Home(context.Context) (*protocol_common.Status, error) {
	return &protocol_common.Status{}, nil
}
func (c *jogTestClient) SetCommandTableWithOptions(context.Context, *client_command_tables.CommandTable, client_command_tables.SetCommandTableOptions) error {
	return nil
}
func (c *jogTestClient) GetCommandTable(context.Context) (*client_command_tables.CommandTable, error) {
	return &client_command_tables.CommandTable{}, nil
}
func (c *jogTestClient) StopMotionController(context.Context) error    { return nil }
func (c *jogTestClient) StartMotionController(context.Context) error   { return nil }
func (c *jogTestClient) SaveCommandTableToFlash(context.Context) error { return nil }
func (c *jogTestClient) ReadRAM(context.Context, uint16) (int32, error) {
	return 0, nil
}
func (c *jogTestClient) WriteRAMAndROM(context.Context, uint16, int32) error {
	return nil
}

type jogTestFactory struct {
	client LinMotClient
}

func (f *jogTestFactory) CreateClient(string) (LinMotClient, error) { return f.client, nil }
func (f *jogTestFactory) Close()                                    {}

func TestJog_EnablesDriveWhenAcknowledgeDoesNotReachOperationEnabled(t *testing.T) {
	client := &jogTestClient{
		initialStatus: &protocol_common.Status{
			StatusWord: protocol_control_word.SetBit(0, protocol_control_word.StatusWordBit_Error),
			StateVar:   uint16(protocol_control_word.State_Error) << 8,
			ErrorCode:  0x1234,
		},
		ackStatus: &protocol_common.Status{
			StatusWord: 0x0000,
			StateVar:   uint16(protocol_control_word.State_SwitchOnDisabled) << 8,
		},
		enableStatus: &protocol_common.Status{
			StatusWord: protocol_control_word.SetBit(0, protocol_control_word.StatusWordBit_OperationEnabled),
			StateVar:   uint16(protocol_control_word.State_OperationEnabled) << 8,
		},
	}

	SetClientFactory(&jogTestFactory{client: client})
	defer ResetClientFactory()

	err := Jog(context.Background(), JogConfig{
		Position: 12.5,
		Config: config.Config{
			ClearCores: []config.ClearCoreConfig{
				{LinMots: []config.LinMotConfig{{IP: "10.0.0.5"}}},
			},
		},
	})
	if err != nil {
		t.Fatalf("Jog() error: %v", err)
	}

	if client.ackCalls != 1 {
		t.Fatalf("expected one acknowledge call, got %d", client.ackCalls)
	}
	if client.enableCalls != 1 {
		t.Fatalf("expected one enable call after acknowledge, got %d", client.enableCalls)
	}
	if client.setPosition1 != 1 || client.setPosition2 != 1 {
		t.Fatalf("expected both target positions to be written once, got pos1=%d pos2=%d", client.setPosition1, client.setPosition2)
	}
	if client.setRunModeCall != 2 {
		t.Fatalf("expected two run mode writes, got %d", client.setRunModeCall)
	}
}

func TestJog_CancelledDuringPostEnableSettleWait(t *testing.T) {
	client := &jogTestClient{
		initialStatus: &protocol_common.Status{
			StatusWord: 0x0000,
			StateVar:   uint16(protocol_control_word.State_SwitchOnDisabled) << 8,
		},
		enableStatus: &protocol_common.Status{
			StatusWord: protocol_control_word.SetBit(0, protocol_control_word.StatusWordBit_OperationEnabled),
			StateVar:   uint16(protocol_control_word.State_OperationEnabled) << 8,
		},
	}

	SetClientFactory(&jogTestFactory{client: client})
	defer ResetClientFactory()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := Jog(ctx, JogConfig{
		Position: 12.5,
		Config: config.Config{
			ClearCores: []config.ClearCoreConfig{
				{LinMots: []config.LinMotConfig{{IP: "10.0.0.5"}}},
			},
		},
	})
	if err == nil {
		t.Fatal("Jog() expected cancellation error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Jog() error = %v, want context.Canceled", err)
	}
	if client.enableCalls != 1 {
		t.Fatalf("expected one enable call before cancellation, got %d", client.enableCalls)
	}
	if client.setPosition1 != 0 || client.setPosition2 != 0 {
		t.Fatalf("expected no position writes after cancellation, got pos1=%d pos2=%d", client.setPosition1, client.setPosition2)
	}
	if client.setRunModeCall != 0 {
		t.Fatalf("expected no run mode writes after cancellation, got %d", client.setRunModeCall)
	}
}
