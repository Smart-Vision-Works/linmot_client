package linmot

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	gsail_client "github.com/Smart-Vision-Works/staged_robot/client"
	client_command_tables "github.com/Smart-Vision-Works/staged_robot/client/rtc/command_tables"
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
)

var (
	linmotMode = flag.String("linmot_mode", "mock", "client mode: mock|udp")
	linmotIP   = flag.String("linmot_ip", "10.8.7.232", "LinMot drive IP for UDP mode tests")
)

// panicFactory is installed as the global factory during all tests in this package.
// Any test that forgets to install its own mock factory will panic immediately with a
// clear message instead of silently attempting a real network connection.
type panicFactory struct{}

func (f *panicFactory) CreateClient(ip string) (LinMotClient, error) {
	panic(fmt.Sprintf(
		"BUG: test attempted to create a real LinMot connection to %q — "+
			"call linmot.SetClientFactory(mockFactory) before exercising code that creates clients",
		ip,
	))
}

func (f *panicFactory) Close() {}

// panicClient satisfies LinMotClient but panics on every call.
// It is never returned by panicFactory (which panics in CreateClient), but
// exists so the compiler can verify panicFactory satisfies ClientFactory.
type panicClient struct{ ip string }

func (c *panicClient) Close() error { panic("panicClient: Close called") }
func (c *panicClient) GetPosition(_ context.Context) (float64, error) {
	panic("panicClient: GetPosition called")
}
func (c *panicClient) GetStatus(_ context.Context) (*protocol_common.Status, error) {
	panic("panicClient: GetStatus called")
}
func (c *panicClient) CheckDriveFault(_ context.Context) error {
	panic("panicClient: CheckDriveFault called")
}
func (c *panicClient) EnableDrive(_ context.Context) (*protocol_common.Status, error) {
	panic("panicClient: EnableDrive called")
}
func (c *panicClient) AcknowledgeError(_ context.Context) (*protocol_common.Status, error) {
	panic("panicClient: AcknowledgeError called")
}
func (c *panicClient) SendControlWord(_ context.Context, _ uint16) (*protocol_common.Status, error) {
	panic("panicClient: SendControlWord called")
}
func (c *panicClient) WriteOutputsWithMask(_ context.Context, _, _ uint16) (*protocol_common.Status, error) {
	panic("panicClient: WriteOutputsWithMask called")
}
func (c *panicClient) VAIGoToPosition(_ context.Context, _, _, _, _ float64) (*protocol_common.Status, error) {
	panic("panicClient: VAIGoToPosition called")
}
func (c *panicClient) VAIGoToPositionFromActual(_ context.Context, _, _, _, _ float64) (*protocol_common.Status, error) {
	panic("panicClient: VAIGoToPositionFromActual called")
}
func (c *panicClient) SetPosition1(_ context.Context, _ float64, _ protocol_common.ParameterStorageType) error {
	panic("panicClient: SetPosition1 called")
}
func (c *panicClient) SetPosition2(_ context.Context, _ float64, _ protocol_common.ParameterStorageType) error {
	panic("panicClient: SetPosition2 called")
}
func (c *panicClient) SetRunMode(_ context.Context, _ protocol_common.RunMode, _ protocol_common.ParameterStorageType) error {
	panic("panicClient: SetRunMode called")
}
func (c *panicClient) SetCommandTable(_ context.Context, _ *client_command_tables.CommandTable) error {
	panic("panicClient: SetCommandTable called")
}
func (c *panicClient) SetEasyStepsAutoStart(_ context.Context, _ int32, _ protocol_common.ParameterStorageType) error {
	panic("panicClient: SetEasyStepsAutoStart called")
}
func (c *panicClient) SetEasyStepsAutoHome(_ context.Context, _ int32, _ protocol_common.ParameterStorageType) error {
	panic("panicClient: SetEasyStepsAutoHome called")
}
func (c *panicClient) SetEasyStepsInputRisingEdgeFunction(_ context.Context, _ protocol_common.IOPinNumber, _ int32, _ protocol_common.ParameterStorageType) error {
	panic("panicClient: SetEasyStepsInputRisingEdgeFunction called")
}
func (c *panicClient) SetEasyStepsInputCurveCmdID(_ context.Context, _ protocol_common.IOPinNumber, _ int32, _ protocol_common.ParameterStorageType) error {
	panic("panicClient: SetEasyStepsInputCurveCmdID called")
}
func (c *panicClient) SetIODefOutputFunction(_ context.Context, _ protocol_common.IOPinNumber, _ int32, _ protocol_common.ParameterStorageType) error {
	panic("panicClient: SetIODefOutputFunction called")
}
func (c *panicClient) SetIODefInputFunction(_ context.Context, _ protocol_common.IOPinNumber, _ int32, _ protocol_common.ParameterStorageType) error {
	panic("panicClient: SetIODefInputFunction called")
}
func (c *panicClient) SetTriggerMode(_ context.Context, _ int32, _ protocol_common.ParameterStorageType) error {
	panic("panicClient: SetTriggerMode called")
}
func (c *panicClient) Home(_ context.Context) (*protocol_common.Status, error) {
	panic("panicClient: Home called")
}
func (c *panicClient) SetCommandTableWithOptions(_ context.Context, _ *client_command_tables.CommandTable, _ client_command_tables.SetCommandTableOptions) error {
	panic("panicClient: SetCommandTableWithOptions called")
}
func (c *panicClient) GetCommandTable(_ context.Context) (*client_command_tables.CommandTable, error) {
	panic("panicClient: GetCommandTable called")
}
func (c *panicClient) StopMotionController(_ context.Context) error {
	panic("panicClient: StopMotionController called")
}
func (c *panicClient) StartMotionController(_ context.Context) error {
	panic("panicClient: StartMotionController called")
}
func (c *panicClient) SaveCommandTableToFlash(_ context.Context) error {
	panic("panicClient: SaveCommandTableToFlash called")
}
func (c *panicClient) ReadRAM(_ context.Context, _ uint16) (int32, error) {
	panic("panicClient: ReadRAM called")
}
func (c *panicClient) WriteRAMAndROM(_ context.Context, _ uint16, _ int32) error {
	panic("panicClient: WriteRAMAndROM called")
}

// Compile-time check: panicClient satisfies LinMotClient.
var _ LinMotClient = (*panicClient)(nil)

// mockSingleClientFactory returns the same pre-created mock client for every
// CreateClient call. Use this when a single test manages its own mock client
// and wants to inject it into the global factory.
type mockSingleClientFactory struct {
	client *gsail_client.Client
}

func (f *mockSingleClientFactory) CreateClient(_ string) (LinMotClient, error) {
	return &noopCloseClient{Client: f.client}, nil
}

func (f *mockSingleClientFactory) Close() {}

// noopCloseClient wraps a *client.Client and makes Close a no-op so tests can
// reuse the same underlying transport without stopping it between calls.
type noopCloseClient struct {
	*gsail_client.Client
}

func (c *noopCloseClient) Close() error { return nil }

// TestMain installs the panic factory before any test in this package runs.
// Both the global factory and the reset target are set to panicFactory, so
// any test that forgets to set up its own mock — or any cleanup that calls
// ResetClientFactory — will never fall back to creating real UDP connections.
func TestMain(m *testing.M) {
	flag.Parse()
	SetClientFactory(&panicFactory{})
	SetDefaultFactory(&panicFactory{})
	os.Exit(m.Run())
}
