package linmot

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockClientFactory_CreateClient(t *testing.T) {
	factory := NewMockClientFactory()
	defer factory.Close()

	// Create a mock client
	client, err := factory.CreateClient("192.168.1.100")
	require.NoError(t, err)
	require.NotNil(t, client)

	// Verify we can call methods on it
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pos, err := client.GetPosition(ctx)
	require.NoError(t, err)
	// MockLinMot starts at position 100000 (in 1/10 um), which is 10mm
	assert.InDelta(t, 10.0, pos, 0.1)
}

func TestMockClientFactory_MultipleMocks(t *testing.T) {
	factory := NewMockClientFactory()
	defer factory.Close()

	// Create clients for different IPs
	client1, err := factory.CreateClient("192.168.1.100")
	require.NoError(t, err)
	require.NotNil(t, client1)

	client2, err := factory.CreateClient("192.168.1.101")
	require.NoError(t, err)
	require.NotNil(t, client2)

	// Verify we have two mocks
	mocks := factory.GetAllMocks()
	assert.Len(t, mocks, 2)
	assert.NotNil(t, mocks["192.168.1.100"])
	assert.NotNil(t, mocks["192.168.1.101"])
}

func TestMockClientFactory_GetMock(t *testing.T) {
	factory := NewMockClientFactory()
	defer factory.Close()

	// No mock should exist initially
	mock := factory.GetMock("192.168.1.100")
	assert.Nil(t, mock)

	// Create a client
	_, err := factory.CreateClient("192.168.1.100")
	require.NoError(t, err)

	// Now we should be able to get the mock
	mock = factory.GetMock("192.168.1.100")
	assert.NotNil(t, mock)
}

func TestMockClientFactory_SetClientFactory(t *testing.T) {
	// Test that we can set the mock factory as the global factory
	factory := NewMockClientFactory()
	defer factory.Close()
	defer ResetClientFactory()

	SetClientFactory(factory)

	// Now globalClientFactory should use our mock factory
	// This is used by the actual stage_primer code
	client, err := globalClientFactory.CreateClient("192.168.1.100")
	require.NoError(t, err)
	require.NotNil(t, client)

	// Verify we can call methods on it
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pos, err := client.GetPosition(ctx)
	require.NoError(t, err)
	assert.InDelta(t, 10.0, pos, 0.1)
}
