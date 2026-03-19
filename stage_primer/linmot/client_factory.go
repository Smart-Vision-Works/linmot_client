package linmot

import (
	"context"

	"github.com/Smart-Vision-Works/staged_robot/client"
	client_command_tables "github.com/Smart-Vision-Works/staged_robot/client/rtc/command_tables"
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
)

// LinMotClient defines the client operations consumed by this package.
// Using an interface lets the package work with pooled real clients in
// production and mock/test doubles in tests.
type LinMotClient interface {
	// Diagnostics
	GetPosition(ctx context.Context) (float64, error)
	GetStatus(ctx context.Context) (*protocol_common.Status, error)
	CheckDriveFault(ctx context.Context) error
	EnableDrive(ctx context.Context) (*protocol_common.Status, error)
	AcknowledgeError(ctx context.Context) (*protocol_common.Status, error)
	SendControlWord(ctx context.Context, word uint16) (*protocol_common.Status, error)
	SetPosition1(ctx context.Context, position float64, storage protocol_common.ParameterStorageType) error
	SetPosition2(ctx context.Context, position float64, storage protocol_common.ParameterStorageType) error
	SetRunMode(ctx context.Context, mode protocol_common.RunMode, storageType protocol_common.ParameterStorageType) error
	SetCommandTable(ctx context.Context, table *client_command_tables.CommandTable) error
	// SetCommandTableWithOptions deploys a command table with configurable MC restart behavior.
	SetCommandTableWithOptions(ctx context.Context, table *client_command_tables.CommandTable, opts client_command_tables.SetCommandTableOptions) error
	// GetCommandTable reads back the current command table from drive RAM. MC must be stopped first.
	GetCommandTable(ctx context.Context) (*client_command_tables.CommandTable, error)
	// StopMotionController stops the MC (required before GetCommandTable).
	StopMotionController(ctx context.Context) error
	// StartMotionController starts the MC.
	StartMotionController(ctx context.Context) error
	// SaveCommandTableToFlash sends the SaveCommandTable RTC command to persist RAM to flash.
	// MC must be stopped. The drive may not respond — callers should expect a timeout.
	SaveCommandTableToFlash(ctx context.Context) error
	// WriteOutputsWithMask sets digital outputs on the X4 interface with a bitmask.
	// Used for vacuum control: bit 0 = purge, bit 1 = vacuum.
	WriteOutputsWithMask(ctx context.Context, bitMask, bitValue uint16) (*protocol_common.Status, error)

	// Motion Control (VAI)
	VAIGoToPosition(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error)
	VAIGoToPositionFromActual(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error)

	// Homing — initiates the homing sequence.
	Home(ctx context.Context) (*protocol_common.Status, error)

	// ReadRAM reads the current RAM value of a parameter by UPID.
	ReadRAM(ctx context.Context, upid uint16) (int32, error)
	// WriteRAMAndROM writes a parameter value to both RAM and ROM by UPID.
	WriteRAMAndROM(ctx context.Context, upid uint16, value int32) error

	// Setup operations — used by Setup() to configure drive hardware parameters.
	SetEasyStepsAutoStart(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error
	SetEasyStepsAutoHome(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error
	SetEasyStepsInputRisingEdgeFunction(ctx context.Context, inputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error
	SetEasyStepsInputCurveCmdID(ctx context.Context, inputNumber protocol_common.IOPinNumber, curveCmdID int32, storageType protocol_common.ParameterStorageType) error
	SetIODefOutputFunction(ctx context.Context, outputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error
	SetIODefInputFunction(ctx context.Context, inputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error
	SetTriggerMode(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error
}

// ClientFactory creates LinMot clients. This interface allows injection of mock clients for testing.
type ClientFactory interface {
	// CreateClient gets a LinMot client for the given IP address from the pool.
	// The client is managed by the pool and should NOT be closed by the caller.
	CreateClient(linmotIP string) (LinMotClient, error)
	// Close closes the factory and all pooled clients.
	Close()
}

// globalClientPool is the shared pool used by all linmot package functions.
// This prevents port 41136 binding conflicts when multiple concurrent commands are sent.
// Both the globalClientFactory and Setup() use this same pool.
var globalClientPool = client.NewClientPool()

// pooledClientFactory uses the global ClientPool to maintain persistent connections.
type pooledClientFactory struct{}

// CreateClient gets a pooled client for the given IP address.
// The client is persistent and should NOT be closed by the caller.
func (f *pooledClientFactory) CreateClient(linmotIP string) (LinMotClient, error) {
	return globalClientPool.GetClient(linmotIP)
}

// Close closes all pooled clients.
func (f *pooledClientFactory) Close() {
	globalClientPool.Close()
}

// globalClientFactory is the factory used by all linmot package functions.
// It defaults to a pooled factory to prevent port binding conflicts.
var globalClientFactory ClientFactory = &pooledClientFactory{}

// resetFactory is the factory restored by ResetClientFactory.
// It defaults to the real UDP factory. Tests override it via SetDefaultFactory
// so that ResetClientFactory never restores to a factory that creates real connections.
var resetFactory ClientFactory = &pooledClientFactory{}

// SetClientFactory sets the global client factory.
// Call this in tests before exercising any code that creates LinMot clients.
// Always pair with a deferred ResetClientFactory call for cleanup.
func SetClientFactory(factory ClientFactory) {
	globalClientFactory = factory
}

// SetDefaultFactory sets the factory that ResetClientFactory will restore to.
// Tests should call this in TestMain to ensure cleanup never reverts to a real
// UDP factory between tests.
func SetDefaultFactory(factory ClientFactory) {
	resetFactory = factory
}

// ResetClientFactory restores the global client factory to the reset target
// (set by SetDefaultFactory). In production this is the real UDP factory;
// in tests TestMain overrides it to a panic factory.
func ResetClientFactory() {
	globalClientFactory = resetFactory
}

// CloseGlobalClientFactory closes the global client factory and all pooled clients.
// This should be called during application shutdown.
func CloseGlobalClientFactory() {
	globalClientFactory.Close()
}

// EvictPooledClient removes a specific client from the global pool, closing its
// socket. The next CreateClient call for that IP creates a fresh connection.
// Used after flash save to reset a poisoned UDP connection.
func EvictPooledClient(linmotIP string) {
	globalClientPool.EvictClient(linmotIP)
}
