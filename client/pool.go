package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"gsail-go/linmot/transport"
)

// ClientPool manages LinMot drive connections with automatic client pooling.
// It maintains a pool of clients indexed by drive IP address to avoid
// creating new UDP connections for every request.
//
// This is particularly useful in HTTP server contexts where multiple requests
// may target the same LinMot drive. The pool is thread-safe and can be used
// concurrently across goroutines.
//
// Example usage:
//
//	pool := client.NewClientPool()
//	defer pool.Close()
//
//	// Get or create a client for a specific drive
//	client, err := pool.GetClient("192.168.1.100")
//	if err != nil {
//	    return err
//	}
//
//	// Use the client (no need to defer client.Close())
//	position, err := client.GetPosition(ctx)
const failedEntryRetryInterval = 5 * time.Second

type clientEntry struct {
	ready    chan struct{}
	client   *Client
	err      error
	failedAt time.Time
}

type ClientPool struct {
	mu              sync.RWMutex
	clients         map[string]*clientEntry
	sharedTransport *transport.SharedUDPTransport
	// factory allows creating clients; tests may inject a mock factory
	factory func(ip string) (*Client, error)
}

// NewClientPool creates a new LinMot client pool.
// The pool uses a shared UDP transport to avoid port binding conflicts
// when multiple drives communicate from the same master port (41136).
func NewClientPool() *ClientPool {
	pool := &ClientPool{
		clients: make(map[string]*clientEntry),
	}

	// Factory creates clients using the shared transport
	pool.factory = func(ip string) (*Client, error) {
		// Lazily create shared transport on first client
		pool.mu.Lock()
		if pool.sharedTransport == nil {
			st, err := transport.NewSharedUDPTransport(
				transport.DefaultMasterPort,
				"", // Bind to any local address
				transport.DefaultTimeout,
			)
			if err != nil {
				pool.mu.Unlock()
				return nil, errors.WithMessage(err, "failed to create shared UDP transport")
			}
			pool.sharedTransport = st
		}
		st := pool.sharedTransport
		pool.mu.Unlock()

		// Register client transport for this drive IP
		clientTransport := st.RegisterClient(ip, transport.DefaultDrivePort)

		// Create client with the transport
		return NewClientWithTransport(ip, clientTransport)
	}

	return pool
}

// newClientPoolWithFactory creates a ClientPool that uses the provided factory for client creation.
// This is useful for tests to inject a mock or fast client implementation.
func newClientPoolWithFactory(factory func(ip string) (*Client, error)) *ClientPool {
	return &ClientPool{
		clients: make(map[string]*clientEntry),
		factory: factory,
	}
}

// GetClient returns an existing client for the given drive IP or creates a new one.
// The client is cached for future requests to the same IP address.
//
// This method is thread-safe and can be called concurrently. It uses a
// future-entry pattern so that only one goroutine creates the client for a
// given IP while others wait on the ready channel without holding the pool lock.
func (cp *ClientPool) GetClient(driveIP string) (*Client, error) {
	// Fast path: check if an entry exists
	cp.mu.RLock()
	entry, ok := cp.clients[driveIP]
	cp.mu.RUnlock()
	if ok {
		<-entry.ready
		if entry.err != nil {
			if time.Since(entry.failedAt) < failedEntryRetryInterval {
				return nil, entry.err
			}
			cp.mu.Lock()
			if current, exists := cp.clients[driveIP]; exists && current == entry {
				delete(cp.clients, driveIP)
			}
			cp.mu.Unlock()
			return cp.GetClient(driveIP)
		}
		return entry.client, nil
	}

	// Create a new entry to indicate we are creating the client
	cp.mu.Lock()
	// Check again in write lock
	if entry, ok = cp.clients[driveIP]; ok {
		cp.mu.Unlock()
		<-entry.ready
		return entry.client, entry.err
	}

	entry = &clientEntry{ready: make(chan struct{})}
	cp.clients[driveIP] = entry
	cp.mu.Unlock()

	// Ensure that ready channel is closed and entry.err is set even if a panic occurs
	defer func() {
		if r := recover(); r != nil {
			entry.err = fmt.Errorf("panic during client creation: %v", r)
			// remove entry on panic
			cp.mu.Lock()
			delete(cp.clients, driveIP)
			cp.mu.Unlock()
			close(entry.ready)
			panic(r)
		}
	}()

	// Create the client outside the lock using the factory
	client, err := cp.factory(driveIP)

	// Acquire lock to store result
	cp.mu.Lock()
	if err != nil {
		entry.err = err
		entry.failedAt = time.Now()
		close(entry.ready)
		cp.mu.Unlock()
		return nil, err
	}

	entry.client = client
	entry.err = nil
	close(entry.ready)
	cp.mu.Unlock()

	return client, nil
}

// Close closes all cached LinMot clients, the shared transport, and clears the pool.
// This should be called when the pool is no longer needed, typically
// during application shutdown.
//
// After Close() is called, the pool is reset; future GetClient calls behave
// as if the pool were new and will create a fresh shared transport on demand.
func (cp *ClientPool) Close() {
	cp.mu.Lock()
	entries := make([]*clientEntry, 0, len(cp.clients))
	for _, entry := range cp.clients {
		entries = append(entries, entry)
	}
	cp.clients = make(map[string]*clientEntry)
	st := cp.sharedTransport
	cp.sharedTransport = nil
	cp.mu.Unlock()

	// Close all clients
	for _, entry := range entries {
		<-entry.ready
		if entry.client != nil {
			_ = entry.client.Close()
		}
	}

	// Close shared transport
	if st != nil {
		_ = st.Close()
	}
}

// EvictClient removes and closes a single client from the pool by IP.
// The next GetClient call for that IP will create a fresh client with a new socket.
// This is used after flash save to reset a poisoned UDP connection.
func (cp *ClientPool) EvictClient(driveIP string) {
	cp.mu.Lock()
	entry, ok := cp.clients[driveIP]
	if ok {
		delete(cp.clients, driveIP)
	}
	cp.mu.Unlock()

	if ok && entry != nil {
		<-entry.ready
		if entry.client != nil {
			_ = entry.client.Close()
		}
	}
}

// Size returns the number of clients currently in the pool.
// This is primarily useful for debugging and testing.
func (cp *ClientPool) Size() int {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return len(cp.clients)
}
