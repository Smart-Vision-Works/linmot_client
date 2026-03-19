package transport

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

// SharedUDPTransport provides a single UDP socket shared by multiple LinMot clients.
// This solves the port binding conflict when multiple drives need to communicate
// from the same master port (41136).
//
// The shared transport owns one UDP socket and routes incoming packets to registered
// clients based on source IP address. Each client gets its own receive channel.
//
// Thread-safe: can be used concurrently from multiple goroutines.
type SharedUDPTransport struct {
	masterPort int
	bindAddr   string
	timeout    time.Duration
	conn       *net.UDPConn

	mu      sync.RWMutex
	writeMu sync.Mutex                        // Guards conn write deadline + write as one atomic operation.
	clients map[string]*sharedClientTransport // driveIP -> client transport
	closed  atomic.Bool

	stopChan chan struct{}
	wg       sync.WaitGroup
}

// NewSharedUDPTransport creates a new shared UDP transport bound to the specified master port.
// bindAddr specifies the local address to bind to (e.g., "0.0.0.0", "127.0.0.1", or "" for any).
// If bindAddr is empty, it defaults to "0.0.0.0" (all interfaces).
//
// The transport immediately binds to the master port and starts a background receive goroutine
// that routes packets to registered clients.
func NewSharedUDPTransport(masterPort int, bindAddr string, timeout time.Duration) (*SharedUDPTransport, error) {
	if bindAddr == "" {
		bindAddr = "0.0.0.0"
	}

	localIP := net.ParseIP(bindAddr).To4()
	if localIP == nil {
		return nil, errors.Errorf("bind address %q is not a valid IPv4 address", bindAddr)
	}

	address := &net.UDPAddr{IP: localIP, Port: masterPort}
	conn, err := net.ListenUDP("udp4", address)
	if err != nil {
		return nil, errors.Wrapf(err,
			"failed to bind to master port %d - LinMot requires packets from this exact port. "+
				"Check if another process is using port %d (e.g., another LinMot client, stage_primer)",
			masterPort, masterPort)
	}

	st := &SharedUDPTransport{
		masterPort: masterPort,
		bindAddr:   bindAddr,
		timeout:    timeout,
		conn:       conn,
		clients:    make(map[string]*sharedClientTransport),
		stopChan:   make(chan struct{}),
	}

	// Start background receive goroutine
	st.wg.Add(1)
	go st.rxLoop()

	return st, nil
}

// RegisterClient registers a new client for a specific drive IP and port.
// Returns a transport.Client that can be used to send/receive packets to/from that drive.
// The returned client shares the underlying UDP socket with all other registered clients.
func (st *SharedUDPTransport) RegisterClient(driveIP string, drivePort int) *sharedClientTransport {
	st.mu.Lock()
	defer st.mu.Unlock()

	// Check if already registered
	if existing, ok := st.clients[driveIP]; ok {
		return existing
	}

	driveAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", driveIP, drivePort))
	if err != nil {
		// This shouldn't happen for valid IPs, but return a client that will fail on send
		driveAddr = &net.UDPAddr{IP: net.ParseIP(driveIP), Port: drivePort}
	}

	client := &sharedClientTransport{
		shared:    st,
		driveIP:   driveIP,
		drivePort: drivePort,
		driveAddr: driveAddr,
		recvChan:  make(chan []byte, 32), // Buffered to avoid blocking rxLoop
	}

	st.clients[driveIP] = client
	return client
}

// UnregisterClient removes a client from the routing table.
// After unregistration, packets from that drive IP will be dropped.
func (st *SharedUDPTransport) UnregisterClient(driveIP string) {
	st.mu.Lock()
	defer st.mu.Unlock()

	if client, ok := st.clients[driveIP]; ok {
		client.closeRecvChan()
		delete(st.clients, driveIP)
	}
}

// Close closes the shared UDP socket and stops the receive goroutine.
// All registered clients will stop receiving packets.
func (st *SharedUDPTransport) Close() error {
	if st.closed.Swap(true) {
		return nil // Already closed
	}

	// Signal stop
	close(st.stopChan)

	// Close socket (this will unblock ReadFromUDP)
	var closeErr error
	if st.conn != nil {
		closeErr = st.conn.Close()
	}

	// Wait for rxLoop to finish
	st.wg.Wait()

	// Close all client channels
	st.mu.Lock()
	for _, client := range st.clients {
		client.closeRecvChan()
	}
	st.clients = make(map[string]*sharedClientTransport)
	st.mu.Unlock()

	return closeErr
}

// rxLoop is the background goroutine that receives packets and routes them to clients.
func (st *SharedUDPTransport) rxLoop() {
	defer st.wg.Done()

	buffer := make([]byte, 1500)

	for {
		select {
		case <-st.stopChan:
			return
		default:
		}

		// Set read deadline to avoid blocking forever
		st.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

		n, remoteAddr, err := st.conn.ReadFromUDP(buffer)
		if err != nil {
			if st.closed.Load() {
				return
			}
			// Timeout or temporary error - continue
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			// Other errors - continue (connection might be closing)
			continue
		}

		// Make a copy of the data
		data := make([]byte, n)
		copy(data, buffer[:n])

		// Extract source IP (strip port)
		sourceIP := remoteAddr.IP.String()

		// Route to registered client
		st.mu.RLock()
		client, ok := st.clients[sourceIP]
		st.mu.RUnlock()

		if !ok {
			// No client registered for this IP - drop packet
			// This is normal when drives respond after client is unregistered
			continue
		}

		// Trace packet after routing
		localAddr := st.conn.LocalAddr().String()
		remoteAddrStr := remoteAddr.String()
		TracePacket("RX", localAddr, remoteAddrStr, data)

		// Non-blocking send to client channel
		select {
		case client.recvChan <- data:
			// Packet delivered
		default:
			// Channel full - drop packet to avoid blocking rxLoop
			// This should be rare with a 32-buffer channel
		}
	}
}
