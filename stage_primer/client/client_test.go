package client

import (
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/Smart-Vision-Works/staged_robot/client"
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	"github.com/Smart-Vision-Works/staged_robot/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	clearcore "stage_primer_config"
	"primer/linmot"
	pb "primer/proto"
	server "primer/server"
)

// mockClientFactory creates mock LinMot clients for testing
type mockClientFactory struct {
	mockClient *client.Client
	mockDrive  *test.MockLinMot
	mu         sync.Mutex
}

func newMockClientFactory(t *testing.T) (*mockClientFactory, func()) {
	// Create mock LinMot client and drive
	mockClient, transportServer := client.NewMockClient()
	mockDrive := test.NewMockLinMot(transportServer)
	mockDrive.Start()

	// Initialize mock drive with a known position (10.0 mm)
	mockDrive.SetStatus(&protocol_common.Status{
		StatusWord:     0x0001,
		StateVar:       0x0002,
		ActualPosition: 100000, // 10.0 mm in LinMot units (100000 = 10.0 mm)
		DemandPosition: 100000,
		Current:        100,
		WarnWord:       0,
		ErrorCode:      0,
	})

	factory := &mockClientFactory{
		mockClient: mockClient,
		mockDrive:  mockDrive,
	}

	cleanup := func() {
		mockDrive.Close()
		mockClient.Close() // Close the real client at cleanup time
	}

	return factory, cleanup
}

// reusableClient wraps a client.Client to prevent Close() from stopping the request manager.
type reusableClient struct {
	*client.Client
}

// Close is a no-op to prevent the request manager from being stopped.
// This allows the client to be reused after "closing".
func (c *reusableClient) Close() error {
	// Don't actually close - allow reuse
	return nil
}

func (f *mockClientFactory) CreateClient(linmotIP string) (linmot.LinMotClient, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Return a wrapper that prevents Close() from stopping the request manager.
	// Since the factory now returns LinMotClient interface, we can return our wrapper.
	return &reusableClient{Client: f.mockClient}, nil
}

// Close closes the factory and all its resources.
func (f *mockClientFactory) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()
	// No-op for testing - cleanup is done via cleanup func
}

func resetClientConnections() {
	connectionsMu.Lock()
	defer connectionsMu.Unlock()
	for _, conn := range connections {
		_ = conn.Close()
	}
	connections = make(map[string]*grpc.ClientConn)
	stageClients = make(map[string]pb.StagePrimerClient)
}

// setupTestGRPCServer creates a real stage primer gRPC server with mock LinMot clients
func setupTestGRPCServer(t *testing.T) (string, func()) {
	// Create mock client factory
	mockFactory, mockCleanup := newMockClientFactory(t)

	// Set the global factory to use mocks
	linmot.SetClientFactory(mockFactory)

	// Create config with the mock LinMot IP
	config := server.CreateDefaultTestConfig()

	configPath, cleanup := server.CreateTempConfigFile(t, config)

	// Create gRPC server instance
	grpcSrv := server.NewGRPCServer(config)
	grpcSrv.SetMockMode(true)
	grpcSrv.SetConfigPath(configPath)
	grpcServer := grpc.NewServer()
	pb.RegisterStagePrimerServer(grpcServer, grpcSrv)
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "Failed to listen on test port")
	go func() {
		_ = grpcServer.Serve(lis)
	}()

	cleanupFunc := func() {
		grpcServer.Stop()
		_ = lis.Close()
		resetClientConnections()
		linmot.ResetClientFactory() // Reset to default factory
		mockCleanup()
		cleanup()
	}

	return lis.Addr().String(), cleanupFunc
}

func setupTestGRPCServerOnAddr(t *testing.T, addr string) (string, func()) {
	// Create mock client factory
	mockFactory, mockCleanup := newMockClientFactory(t)

	// Set the global factory to use mocks
	linmot.SetClientFactory(mockFactory)

	// Create config with the mock LinMot IP
	config := clearcore.Config{
		ClearCores: []clearcore.ClearCoreConfig{
			{
				LinMots: []clearcore.LinMotConfig{
					{IP: "127.0.0.1"}, // This will use the mock client
				},
			},
		},
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStagePrimerServer(grpcServer, server.NewGRPCServer(config))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		mockCleanup()
		linmot.ResetClientFactory()
		t.Skipf("unable to listen on %s: %v", addr, err)
	}
	go func() {
		_ = grpcServer.Serve(lis)
	}()

	cleanupFunc := func() {
		grpcServer.Stop()
		_ = lis.Close()
		resetClientConnections()
		linmot.ResetClientFactory()
		mockCleanup()
	}

	return lis.Addr().String(), cleanupFunc
}

func TestGetStagePosition_Success(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	// Now with mock LinMot, this should succeed!
	position, err := GetStagePosition(grpcAddr, 0, 0)
	require.NoError(t, err)
	// Mock drive is initialized to 10.0 mm (100000 units = 10.0 mm)
	assert.Equal(t, 10.0, position)
}

func TestGetStagePosition_InvalidRobotIndex(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	// Robot index out of range
	position, err := GetStagePosition(grpcAddr, 999, 0)
	require.Error(t, err)
	assert.Equal(t, 0.0, position)
	assert.Contains(t, err.Error(), "failed to find LinMot IP")
}

func TestGetStagePosition_InvalidStageIndex(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	// Stage index out of range
	position, err := GetStagePosition(grpcAddr, 0, 999)
	require.Error(t, err)
	assert.Equal(t, 0.0, position)
	assert.Contains(t, err.Error(), "failed to find LinMot IP")
}

func TestJogStage_Success(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	// With mock LinMot, this should succeed!
	err := JogStage(grpcAddr, 0, 0, 100.0)
	require.NoError(t, err)
}

func TestJogStage_InvalidRobotIndex(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	err := JogStage(grpcAddr, 999, 0, 100.0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find LinMot IP")
}

func TestJogStageOffset_Success(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	// With mock LinMot, this should succeed!
	err := JogStageOffset(grpcAddr, 0, 0, 10.5)
	require.NoError(t, err)
}

func TestJogStageOffset_InvalidRobotIndex(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	err := JogStageOffset(grpcAddr, 999, 0, 10.5)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find LinMot IP")
}

func TestDeployStageCommandTable_Success(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	// Save current working directory
	originalWD, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWD)

	// Create temporary directory and change to it
	tmpDir := t.TempDir()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create template file in expected dev path
	templatePath := "primer/linmot/linmot_command_table.yaml"
	err = os.MkdirAll(filepath.Dir(templatePath), 0755)
	require.NoError(t, err)

	templateContent := `version: 1
drive_model: C1250-MI
entries:
  - id: 1
    name: "test entry"
    type: NoOp
    par1: ${POSITION_DOWN}
    par2: ${MAX_VELOCITY}
    par3: ${ACCELERATION}
    par4: ${DELAY_AT_BOTTOM}
`
	err = os.WriteFile(templatePath, []byte(templateContent), 0644)
	require.NoError(t, err)

	params := CommandTableParams{
		DefaultSpeed:        50.0,
		DefaultAcceleration: 100.0,
		PickTime:            0.5,
		ZDistance:           25.5,
	}

	// With mock LinMot and template file, this should succeed!
	err = DeployStageCommandTable(grpcAddr, 0, 0, params)
	require.NoError(t, err)
}

func TestDeployStageCommandTable_InvalidParams(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	// Invalid: negative zDistance
	params := CommandTableParams{
		DefaultSpeed:        50.0,
		DefaultAcceleration: 100.0,
		PickTime:            0.5,
		ZDistance:           -1.0, // Invalid
	}

	err := DeployStageCommandTable(grpcAddr, 0, 0, params)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid deployment configuration")
}

func TestDeployStageCommandTable_InvalidRobotIndex(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	params := CommandTableParams{
		DefaultSpeed:        50.0,
		DefaultAcceleration: 100.0,
		PickTime:            0.5,
		ZDistance:           25.5,
	}

	err := DeployStageCommandTable(grpcAddr, 999, 0, params)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ConfigStore lookup failed")
}

func TestDeployStageCommandTable_InvalidSpeed(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	// Invalid: speed too high
	params := CommandTableParams{
		DefaultSpeed:        150.0, // Invalid: > 100%
		DefaultAcceleration: 100.0,
		PickTime:            0.5,
		ZDistance:           25.5,
	}

	err := DeployStageCommandTable(grpcAddr, 0, 0, params)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "default_speed out of range")
}

func TestContainsPort(t *testing.T) {
	tests := []struct {
		name     string
		addr     string
		expected bool
	}{
		{"IPv4 with port", "192.168.1.1:8080", true},
		{"IPv4 without port", "192.168.1.1", false},
		{"IPv6 with port", "[2001:db8::1]:8080", true},
		{"IPv6 without port", "2001:db8::1", true}, // IPv6 addresses contain colons, so containsPort returns true (known limitation)
		{"localhost with port", "localhost:8080", true},
		{"localhost without port", "localhost", false},
		{"empty string", "", false},
		{"just colon", ":", false},
		{"colon at end", "192.168.1.1:", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsPort(tt.addr)
			assert.Equal(t, tt.expected, result, "containsPort(%q) = %v, want %v", tt.addr, result, tt.expected)
		})
	}
}

func TestGetStagePosition_WithPortInIP(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	// Use the full address with port (as returned by net.Listen)
	position, err := GetStagePosition(grpcAddr, 0, 0)
	require.NoError(t, err)
	assert.Equal(t, 10.0, position)
}

func TestGetStagePosition_WithoutPortInIP(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServerOnAddr(t, "127.0.0.1:50051")
	defer cleanup()

	// Extract just the IP without port - simulate real-world usage where IP comes without port
	host, _, err := net.SplitHostPort(grpcAddr)
	require.NoError(t, err)
	// The client should add :50051 automatically when no port is supplied.
	position, err := GetStagePosition(host, 0, 0)
	require.NoError(t, err)
	assert.Equal(t, 10.0, position)
}

func TestJogStage_NetworkError(t *testing.T) {
	// Use an invalid IP to simulate network error
	primerIP := "192.0.2.0:50051" // Test IP that should not be reachable

	err := JogStage(primerIP, 0, 0, 100.0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestDeployStageCommandTable_AllZeroParams(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	params := CommandTableParams{
		DefaultSpeed:        0.0,
		DefaultAcceleration: 0.0,
		PickTime:            0.0,
		ZDistance:           0.0,
	}

	err := DeployStageCommandTable(grpcAddr, 0, 0, params)
	// Zero params will fail validation, which is expected
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid deployment configuration")
}

func TestSetVacuum_Success(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	err := SetVacuum(grpcAddr, 0, 0, "on")
	require.NoError(t, err)
}

func TestGetUSBDevices_Success(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	devices, err := GetUSBDevices(grpcAddr)
	require.NoError(t, err)
	assert.NotEmpty(t, devices)
	assert.Equal(t, "Teknic Inc.", devices[0].Manufacturer)
}

func TestConfig_Success(t *testing.T) {
	grpcAddr, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	initialConfig, err := GetConfig(grpcAddr)
	require.NoError(t, err)
	require.NotNil(t, initialConfig)

	// Update config
	newConfig := *initialConfig
	newConfig.ClearCores = append(newConfig.ClearCores, ClearCoreConfig{
		USBID:     "new-device",
		DHCP:      false,
		IPAddress: "10.0.0.60",
		Gateway:   "10.0.0.1",
		Subnet:    "255.255.255.0",
		DNS:       "8.8.8.8",
		LinMots:   []LinMotConfig{},
	})

	err = SetConfig(grpcAddr, newConfig)
	require.NoError(t, err)

	updatedConfig, err := GetConfig(grpcAddr)
	require.NoError(t, err)
	assert.Equal(t, newConfig, *updatedConfig)
}
