package transport

import (
	"context"

	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

// udpTransportClient implements Client using UDP transport.
// This is a pure I/O implementation with NO protocol knowledge.
// It simply sends and receives raw bytes over UDP.
// Thread-safe: can be used concurrently from multiple goroutines.
type udpTransportClient struct {
	driveIP    string
	drivePort  int
	masterPort int
	connection *net.UDPConn
	timeout    time.Duration
	debug      atomic.Bool
}

// NewUDPTransportClient creates a new UDP transport client for the specified drive IP and ports.
// bindAddr specifies the local address to bind to (e.g., "0.0.0.0", "127.0.0.1", or "" for any).
// If bindAddr is empty, it defaults to "0.0.0.0" (all interfaces).
//
// This creates a pure I/O transport with no protocol awareness.
func NewUDPTransportClient(driveIP string, drivePort, masterPort int, bindAddr string, timeout time.Duration) (*udpTransportClient, error) {
	if driveIP == "" {
		return nil, errors.New("drive IP address cannot be empty")
	}

	channel := &udpTransportClient{
		driveIP:    driveIP,
		drivePort:  drivePort,
		masterPort: masterPort,
		timeout:    timeout,
	}

	localIP, _, err := resolveLocalIPv4(driveIP, drivePort, bindAddr)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to resolve local IPv4")
	}

	connection, err := createUDPConnection(localIP, masterPort)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create UDP connection")
	}
	channel.connection = connection

	return channel, nil
}

// SetDebug enables or disables transport-level debug logging.
func (c *udpTransportClient) SetDebug(enabled bool) {
	c.debug.Store(enabled)
}

// UDPInfo returns the socket details for debug logging.
func (c *udpTransportClient) UDPInfo() UDPInfo {
	localAddr := ""
	if c.connection != nil {
		localAddr = c.connection.LocalAddr().String()
	}
	remoteAddr := fmt.Sprintf("%s:%d", c.driveIP, c.drivePort)
	return UDPInfo{
		LocalAddr:  localAddr,
		RemoteAddr: remoteAddr,
		MasterPort: c.masterPort,
		DrivePort:  c.drivePort,
	}
}

// LocalAddr implements transport.ConnectionInfo for diagnostic error messages.
func (c *udpTransportClient) LocalAddr() string {
	if c.connection != nil {
		return c.connection.LocalAddr().String()
	}
	return ""
}

// RemoteAddr implements transport.ConnectionInfo for diagnostic error messages.
func (c *udpTransportClient) RemoteAddr() string {
	return fmt.Sprintf("%s:%d", c.driveIP, c.drivePort)
}

// resolveLocalIPv4 selects the local IPv4 address used to reach the drive.
// If bindAddr is provided and valid, it wins over route-based selection.
func resolveLocalIPv4(driveIP string, drivePort int, bindAddr string) (net.IP, *net.UDPAddr, error) {
	remoteAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", driveIP, drivePort))
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to resolve drive address")
	}

	if bindAddr != "" && bindAddr != "0.0.0.0" {
		localIP := net.ParseIP(bindAddr).To4()
		if localIP == nil {
			return nil, nil, errors.Errorf("bind address %q is not a valid IPv4 address", bindAddr)
		}
		return localIP, remoteAddr, nil
	}

	tmp, err := net.DialUDP("udp4", nil, remoteAddr)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to determine local route")
	}
	localAddr, ok := tmp.LocalAddr().(*net.UDPAddr)
	tmp.Close()
	if !ok || localAddr.IP == nil {
		return nil, nil, errors.New("failed to determine local IPv4 address")
	}
	localIP := localAddr.IP.To4()
	if localIP == nil {
		return nil, nil, errors.New("local route address is not IPv4")
	}
	return localIP, remoteAddr, nil
}

// createUDPConnection creates a UDP connection bound to the specified master port.
// CRITICAL: The LinMot drive ONLY responds to packets from the designated master port (default 41136).
// Sending from any other port will cause the LinMot to ignore responses and may wedge the LinUDP interface.
// If the master port is already in use, this function returns an error rather than falling back
// to a random port, which would break communication with the drive.
func createUDPConnection(localIP net.IP, masterPort int) (*net.UDPConn, error) {
	address := &net.UDPAddr{IP: localIP, Port: masterPort}
	connection, err := net.ListenUDP("udp4", address)
	if err != nil {
		return nil, errors.Wrapf(err,
			"failed to bind to master port %d - LinMot requires packets from this exact port. "+
				"Check if another process is using port %d (e.g., another process or LinMot client)",
			masterPort, masterPort)
	}

	return connection, nil
}

// SendPacket Send transmits raw bytes to the drive over UDP.
// This is a pure I/O operation with no knowledge of packet contents.
func (c *udpTransportClient) SendPacket(ctx context.Context, data []byte) error {
	// Check context before starting
	if err := ctx.Err(); err != nil {
		return err
	}

	driveAddress, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", c.driveIP, c.drivePort))
	if err != nil {
		return errors.WithMessage(err, "failed to resolve drive address")
	}

	// Set deadline from context if present, otherwise use timeout
	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		deadline = time.Now().Add(c.timeout)
	}
	if err := c.connection.SetWriteDeadline(deadline); err != nil {
		return errors.WithMessage(err, "failed to set write deadline")
	}

	// Run send in goroutine to support context cancellation
	type errorOpt struct {
		err error
	}
	resultChannel := make(chan errorOpt, 1)

	go func() {
		// Trace TX packet before sending
		localAddr := c.connection.LocalAddr().String()
		remoteAddr := driveAddress.String()
		TracePacket("TX", localAddr, remoteAddr, data)

		_, err := c.connection.WriteToUDP(data, driveAddress)
		if err != nil {
			resultChannel <- errorOpt{err: errors.WithMessage(err, "failed to send packet")}
		} else {
			resultChannel <- errorOpt{}
		}
	}()

	// Wait for either result or context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	case response := <-resultChannel:
		return response.err
	}
}

// RecvPacket Receive reads raw bytes from the drive over UDP.
// This is a pure I/O operation with no knowledge of packet contents.
// Returns the received bytes or an error.
func (c *udpTransportClient) RecvPacket(ctx context.Context) ([]byte, error) {
	data, _, err := c.RecvPacketWithAddr(ctx)
	return data, err
}

// RecvPacketWithAddr Receive reads raw bytes from the drive over UDP and returns the sender address.
// This is a pure I/O operation with no knowledge of packet contents.
// Returns the received bytes, sender address string, or an error.
func (c *udpTransportClient) RecvPacketWithAddr(ctx context.Context) ([]byte, string, error) {
	// Check context before starting
	if err := ctx.Err(); err != nil {
		return nil, "", err
	}

	// Set deadline from context if present, otherwise use timeout
	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		deadline = time.Now().Add(c.timeout)
	}
	if err := c.connection.SetReadDeadline(deadline); err != nil {
		return nil, "", errors.WithMessage(err, "failed to set read deadline")
	}

	// Run receive in goroutine to support context cancellation
	type readResult struct {
		data []byte
		addr string
		err  error
	}
	resultChannel := make(chan readResult, 1)

	go func() {
		buffer := make([]byte, 1500)
		n, remoteAddr, err := c.connection.ReadFromUDP(buffer)
		if err != nil {
			resultChannel <- readResult{err: errors.WithMessage(err, "failed to receive response")}
		} else {
			// Return a copy of the data, not the buffer
			data := make([]byte, n)
			copy(data, buffer[:n])

			// Trace RX packet after receiving
			localAddr := c.connection.LocalAddr().String()
			remoteAddrStr := ""
			if remoteAddr != nil {
				remoteAddrStr = remoteAddr.String()
			}
			TracePacket("RX", localAddr, remoteAddrStr, data)

			resultChannel <- readResult{data: data, addr: remoteAddrStr}
		}
	}()

	// Wait for either result or context cancellation
	select {
	case <-ctx.Done():
		return nil, "", ctx.Err()
	case result := <-resultChannel:
		return result.data, result.addr, result.err
	}
}

// Close closes the UDP connection.
func (c *udpTransportClient) Close() error {
	if c.connection != nil {
		return c.connection.Close()
	}
	return nil
}
