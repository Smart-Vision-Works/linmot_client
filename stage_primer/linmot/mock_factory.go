package linmot

import (
	"log"
	"sync"

	"github.com/Smart-Vision-Works/staged_robot/client"
	"github.com/Smart-Vision-Works/staged_robot/test"
	"github.com/Smart-Vision-Works/staged_robot/transport"
)

// MockClientFactory creates mock LinMot clients for testing.
// Each unique IP gets its own MockLinMot instance that simulates
// the full LinMot protocol including command tables, jog, and position.
//
// This factory is used when stage_primer is started with --mock-linmot flag,
// allowing E2E tests to run without real LinMot hardware.
type MockClientFactory struct {
	mu    sync.Mutex
	mocks map[string]*mockClientEntry
}

type mockClientEntry struct {
	client    *client.Client
	mock      *test.MockLinMot
	transport *transport.Server
}

// NewMockClientFactory creates a new factory for mock LinMot clients.
func NewMockClientFactory() *MockClientFactory {
	return &MockClientFactory{
		mocks: make(map[string]*mockClientEntry),
	}
}

// CreateClient returns a mock LinMot client for the given IP.
// Multiple calls with the same IP return new clients connected to the same mock drive.
// The mock drive simulates LinMot protocol including command tables, jog, and position.
func (f *MockClientFactory) CreateClient(linmotIP string) (LinMotClient, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Check if we already have a mock for this IP
	if entry, ok := f.mocks[linmotIP]; ok {
		// Create a new client connected to the existing mock
		newClient, transportServer := client.NewMockClient()
		// Wire up a new MockLinMot to handle this client's requests
		newMock := test.NewMockLinMot(transportServer)
		newMock.Start()
		log.Printf("[MockClientFactory] Created additional mock client for IP %s", linmotIP)
		// Note: We keep the original entry for reference but each client gets its own mock
		_ = entry
		return newClient, nil
	}

	// Create new mock transport pair
	mockClient, transportServer := client.NewMockClient()

	// Create MockLinMot to simulate the drive
	mock := test.NewMockLinMot(transportServer)
	mock.Start()

	// Store the entry
	f.mocks[linmotIP] = &mockClientEntry{
		client:    mockClient,
		mock:      mock,
		transport: &transportServer,
	}

	log.Printf("[MockClientFactory] Created mock LinMot for IP %s", linmotIP)
	return mockClient, nil
}

// GetMock returns the MockLinMot instance for a given IP, for test inspection.
// Returns nil if no mock exists for that IP.
func (f *MockClientFactory) GetMock(linmotIP string) *test.MockLinMot {
	f.mu.Lock()
	defer f.mu.Unlock()

	if entry, ok := f.mocks[linmotIP]; ok {
		return entry.mock
	}
	return nil
}

// GetAllMocks returns all mock instances for test inspection.
func (f *MockClientFactory) GetAllMocks() map[string]*test.MockLinMot {
	f.mu.Lock()
	defer f.mu.Unlock()

	result := make(map[string]*test.MockLinMot)
	for ip, entry := range f.mocks {
		result[ip] = entry.mock
	}
	return result
}

// Close stops all mock LinMot instances.
func (f *MockClientFactory) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()

	for ip, entry := range f.mocks {
		log.Printf("[MockClientFactory] Stopping mock LinMot for IP %s", ip)
		entry.mock.Stop()
		entry.client.Close()
	}
	f.mocks = make(map[string]*mockClientEntry)
}
