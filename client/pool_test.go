package client

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockFactory(ip string) (*Client, error) {
	// ip is unused in the mock implementation, but keep it to match the factory signature
	_ = ip
	client, _ := NewMockClient()
	return client, nil
}

func TestClientPool_NewClientPool(t *testing.T) {
	pool := newClientPoolWithFactory(mockFactory)
	require.NotNil(t, pool)
	require.NotNil(t, pool.clients)
	assert.Equal(t, 0, pool.Size())
}

func TestClientPool_GetClient_SingleIP(t *testing.T) {
	pool := newClientPoolWithFactory(mockFactory)
	defer pool.Close()

	// First call should create new client
	client1, err := pool.GetClient("192.168.1.100")
	require.NoError(t, err)
	require.NotNil(t, client1)
	assert.Equal(t, 1, pool.Size())

	// Second call should return cached client
	client2, err := pool.GetClient("192.168.1.100")
	require.NoError(t, err)
	require.NotNil(t, client2)
	assert.Equal(t, 1, pool.Size())

	// Should be same client
	assert.Same(t, client1, client2)
}

func TestClientPool_GetClient_MultipleIPs(t *testing.T) {
	pool := newClientPoolWithFactory(mockFactory)
	defer pool.Close()

	// Create clients for different IPs
	client1, err := pool.GetClient("192.168.1.100")
	require.NoError(t, err)
	require.NotNil(t, client1)

	client2, err := pool.GetClient("192.168.1.101")
	require.NoError(t, err)
	require.NotNil(t, client2)

	client3, err := pool.GetClient("192.168.1.102")
	require.NoError(t, err)
	require.NotNil(t, client3)

	// Should have 3 separate clients
	assert.Equal(t, 3, pool.Size())
	assert.NotSame(t, client1, client2)
	assert.NotSame(t, client2, client3)
	assert.NotSame(t, client1, client3)

	// Retrieve cached clients
	cachedClient1, err := pool.GetClient("192.168.1.100")
	require.NoError(t, err)
	assert.Same(t, client1, cachedClient1)

	cachedClient2, err := pool.GetClient("192.168.1.101")
	require.NoError(t, err)
	assert.Same(t, client2, cachedClient2)
}

func TestClientPool_Close(t *testing.T) {
	pool := newClientPoolWithFactory(mockFactory)

	// Create some clients
	_, err := pool.GetClient("192.168.1.100")
	require.NoError(t, err)
	_, err = pool.GetClient("192.168.1.101")
	require.NoError(t, err)

	assert.Equal(t, 2, pool.Size())

	// Close pool
	pool.Close()

	// Pool should be empty
	assert.Equal(t, 0, pool.Size())
}

func TestClientPool_ConcurrentAccess(t *testing.T) {
	pool := newClientPoolWithFactory(mockFactory)
	defer pool.Close()

	const numGoroutines = 10
	const numRequests = 100
	var wg sync.WaitGroup

	// Multiple goroutines requesting same IP concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numRequests; j++ {
				client, err := pool.GetClient("192.168.1.100")
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		}()
	}

	wg.Wait()

	// Should only have created one client despite many concurrent requests
	assert.Equal(t, 1, pool.Size())
}

func TestClientPool_ConcurrentMultipleIPs(t *testing.T) {
	pool := newClientPoolWithFactory(mockFactory)
	defer pool.Close()

	const numGoroutines = 5
	const numIPs = 3
	var wg sync.WaitGroup

	ips := []string{
		"192.168.1.100",
		"192.168.1.101",
		"192.168.1.102",
	}

	// Multiple goroutines requesting different IPs
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			ip := ips[goroutineID%numIPs]
			for j := 0; j < 50; j++ {
				client, err := pool.GetClient(ip)
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		}(i)
	}

	wg.Wait()

	// Should have created exactly 3 clients (one per IP)
	assert.Equal(t, numIPs, pool.Size())
}

func TestClientPool_Size(t *testing.T) {
	pool := newClientPoolWithFactory(mockFactory)
	defer pool.Close()

	assert.Equal(t, 0, pool.Size())

	_, _ = pool.GetClient("192.168.1.100")
	assert.Equal(t, 1, pool.Size())

	_, _ = pool.GetClient("192.168.1.101")
	assert.Equal(t, 2, pool.Size())

	_, _ = pool.GetClient("192.168.1.100") // Cached
	assert.Equal(t, 2, pool.Size())

	pool.Close()
	assert.Equal(t, 0, pool.Size())
}
