package transport

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSharedUDPTransport_Creation(t *testing.T) {
	st, err := NewSharedUDPTransport(0, "127.0.0.1", time.Second)
	require.NoError(t, err)
	require.NotNil(t, st)
	defer st.Close()

	assert.NotNil(t, st.conn)
	assert.Equal(t, "127.0.0.1", st.bindAddr)
}

func TestSharedUDPTransport_RegisterClient(t *testing.T) {
	st, err := NewSharedUDPTransport(0, "127.0.0.1", time.Second)
	require.NoError(t, err)
	defer st.Close()

	client1 := st.RegisterClient("10.8.7.232", 49360)
	require.NotNil(t, client1)
	assert.Equal(t, "10.8.7.232", client1.driveIP)
	assert.Equal(t, 49360, client1.drivePort)

	// Registering same IP again should return same client
	client2 := st.RegisterClient("10.8.7.232", 49360)
	assert.Same(t, client1, client2)

	// Different IP should get different client
	client3 := st.RegisterClient("10.8.7.234", 49360)
	assert.NotSame(t, client1, client3)
	assert.Equal(t, "10.8.7.234", client3.driveIP)
}

func TestSharedUDPTransport_UnregisterClient(t *testing.T) {
	st, err := NewSharedUDPTransport(0, "127.0.0.1", time.Second)
	require.NoError(t, err)
	defer st.Close()

	client := st.RegisterClient("10.8.7.232", 49360)
	require.NotNil(t, client)

	// Should have 1 client
	st.mu.RLock()
	assert.Len(t, st.clients, 1)
	st.mu.RUnlock()

	// Unregister
	st.UnregisterClient("10.8.7.232")

	// Should have 0 clients
	st.mu.RLock()
	assert.Len(t, st.clients, 0)
	st.mu.RUnlock()
}

func TestSharedUDPTransport_SendRecv(t *testing.T) {
	// Create shared transport
	st, err := NewSharedUDPTransport(0, "127.0.0.1", time.Second)
	require.NoError(t, err)
	defer st.Close()

	masterPort := st.conn.LocalAddr().(*net.UDPAddr).Port

	// Create a mock drive server
	driveAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:0")
	require.NoError(t, err)
	driveConn, err := net.ListenUDP("udp4", driveAddr)
	require.NoError(t, err)
	defer driveConn.Close()

	drivePort := driveConn.LocalAddr().(*net.UDPAddr).Port

	// Register client for the drive
	client := st.RegisterClient("127.0.0.1", drivePort)

	// Send packet from client
	ctx := context.Background()
	testData := []byte{0x01, 0x02, 0x03, 0x04}
	err = client.SendPacket(ctx, testData)
	require.NoError(t, err)

	// Receive on drive side
	buffer := make([]byte, 1500)
	driveConn.SetReadDeadline(time.Now().Add(time.Second))
	n, clientAddr, err := driveConn.ReadFromUDP(buffer)
	require.NoError(t, err)
	assert.Equal(t, testData, buffer[:n])
	assert.Equal(t, masterPort, clientAddr.Port) // Should come from master port

	// Send response from drive
	responseData := []byte{0x05, 0x06, 0x07, 0x08}
	_, err = driveConn.WriteToUDP(responseData, clientAddr)
	require.NoError(t, err)

	// Receive on client side
	receivedData, err := client.RecvPacket(ctx)
	require.NoError(t, err)
	assert.Equal(t, responseData, receivedData)
}

func TestSharedUDPTransport_ConcurrentSend(t *testing.T) {
	st, err := NewSharedUDPTransport(0, "127.0.0.1", time.Second)
	require.NoError(t, err)
	defer st.Close()

	driveConn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})
	require.NoError(t, err)
	defer driveConn.Close()

	client := st.RegisterClient("127.0.0.1", driveConn.LocalAddr().(*net.UDPAddr).Port)
	require.NotNil(t, client)

	const sends = 32
	received := make(chan struct{}, sends)
	recvErr := make(chan error, 1)

	go func() {
		buffer := make([]byte, 1500)
		for i := 0; i < sends; i++ {
			driveConn.SetReadDeadline(time.Now().Add(2 * time.Second))
			_, _, err := driveConn.ReadFromUDP(buffer)
			if err != nil {
				recvErr <- err
				return
			}
			received <- struct{}{}
		}
	}()

	var wg sync.WaitGroup
	errs := make(chan error, sends)
	for i := 0; i < sends; i++ {
		wg.Add(1)
		go func(id byte) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if err := client.SendPacket(ctx, []byte{0xAA, id}); err != nil {
				errs <- err
			}
		}(byte(i))
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}

	for i := 0; i < sends; i++ {
		select {
		case <-received:
		case err := <-recvErr:
			require.NoError(t, err)
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for drive receives")
		}
	}
}

func TestSharedUDPTransport_MultipleClients(t *testing.T) {
	// Create shared transport
	st, err := NewSharedUDPTransport(0, "127.0.0.1", time.Second)
	require.NoError(t, err)
	defer st.Close()

	masterPort := st.conn.LocalAddr().(*net.UDPAddr).Port

	// Create two mock drives on different ports (both on localhost)
	drive1Conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})
	require.NoError(t, err)
	defer drive1Conn.Close()
	drive1Port := drive1Conn.LocalAddr().(*net.UDPAddr).Port

	drive2Conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})
	require.NoError(t, err)
	defer drive2Conn.Close()
	drive2Port := drive2Conn.LocalAddr().(*net.UDPAddr).Port

	// Register two clients - both using localhost but different ports
	// This simulates multiple drives on the same network
	client1 := st.RegisterClient("127.0.0.1", drive1Port)
	client2 := st.RegisterClient("127.0.0.1", drive2Port)

	// First registration should succeed
	require.NotNil(t, client1)
	// Second registration with same IP should return the same client
	// (since routing is by IP, not port)
	assert.Same(t, client1, client2)

	ctx := context.Background()

	// Test that the client can send packets from the shared master port
	data1 := []byte{0x11, 0x22}
	err = client1.SendPacket(ctx, data1)
	require.NoError(t, err)

	// Receive on drive1
	buffer := make([]byte, 1500)
	drive1Conn.SetReadDeadline(time.Now().Add(time.Second))
	n, addr1, err := drive1Conn.ReadFromUDP(buffer)
	require.NoError(t, err)
	assert.Equal(t, data1, buffer[:n])
	assert.Equal(t, masterPort, addr1.Port, "packet should come from shared master port")

	// Verify only one client is registered (since both registrations were for same IP)
	st.mu.RLock()
	assert.Len(t, st.clients, 1)
	st.mu.RUnlock()

	// Now test registering truly different clients (using invalid IPs that we won't actually send to)
	client3 := st.RegisterClient("10.8.7.232", 49360)
	client4 := st.RegisterClient("10.8.7.234", 49360)

	assert.NotSame(t, client3, client4)
	assert.NotSame(t, client1, client3)

	// Verify we now have 3 clients registered
	st.mu.RLock()
	assert.Len(t, st.clients, 3)
	st.mu.RUnlock()
}

func TestSharedUDPTransport_Close(t *testing.T) {
	st, err := NewSharedUDPTransport(0, "127.0.0.1", time.Second)
	require.NoError(t, err)

	client := st.RegisterClient("10.8.7.232", 49360)
	require.NotNil(t, client)

	// Close transport
	err = st.Close()
	assert.NoError(t, err)

	// Should be marked as closed
	assert.True(t, st.closed.Load())

	// Client's receive channel should be closed
	_, ok := <-client.recvChan
	assert.False(t, ok)

	// Second close should be idempotent
	err = st.Close()
	assert.NoError(t, err)
}

func TestSharedClientTransport_CloseRecvChan_Idempotent(t *testing.T) {
	client := &sharedClientTransport{
		recvChan: make(chan []byte, 1),
	}

	client.closeRecvChan()
	client.closeRecvChan()

	_, ok := <-client.recvChan
	assert.False(t, ok)
}

func TestSharedClientTransport_CloseRecvChan_Concurrent(t *testing.T) {
	client := &sharedClientTransport{
		recvChan: make(chan []byte, 1),
	}

	const n = 32
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			client.closeRecvChan()
		}()
	}
	wg.Wait()

	_, ok := <-client.recvChan
	assert.False(t, ok)
}

func TestSharedUDPTransport_ContextCancellation(t *testing.T) {
	st, err := NewSharedUDPTransport(0, "127.0.0.1", time.Second)
	require.NoError(t, err)
	defer st.Close()

	client := st.RegisterClient("10.8.7.232", 49360)

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Send should fail with context cancelled
	err = client.SendPacket(ctx, []byte{0x01})
	assert.ErrorIs(t, err, context.Canceled)

	// Recv should fail with context cancelled
	_, err = client.RecvPacket(ctx)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestSharedUDPTransport_Timeout(t *testing.T) {
	st, err := NewSharedUDPTransport(0, "127.0.0.1", 50*time.Millisecond)
	require.NoError(t, err)
	defer st.Close()

	client := st.RegisterClient("10.8.7.232", 49360)

	ctx := context.Background()

	// Recv should timeout since no data is available
	start := time.Now()
	_, err = client.RecvPacket(ctx)
	duration := time.Since(start)

	assert.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.GreaterOrEqual(t, duration, 50*time.Millisecond)
	assert.Less(t, duration, 200*time.Millisecond)
}

func TestSharedUDPTransport_ConcurrentAccess(t *testing.T) {
	st, err := NewSharedUDPTransport(0, "127.0.0.1", time.Second)
	require.NoError(t, err)
	defer st.Close()

	const numClients = 10
	var wg sync.WaitGroup

	// Register multiple clients concurrently
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			driveIP := fmt.Sprintf("10.8.7.%d", 100+id)
			client := st.RegisterClient(driveIP, 49360)
			assert.NotNil(t, client)
		}(i)
	}

	wg.Wait()

	st.mu.RLock()
	count := len(st.clients)
	st.mu.RUnlock()

	assert.Equal(t, numClients, count)
}

func TestSharedClientTransport_ConnectionInfo(t *testing.T) {
	st, err := NewSharedUDPTransport(0, "127.0.0.1", time.Second)
	require.NoError(t, err)
	defer st.Close()

	client := st.RegisterClient("10.8.7.232", 49360)

	// Should implement ConnectionInfo interface
	localAddr := client.LocalAddr()
	assert.Contains(t, localAddr, "127.0.0.1")

	remoteAddr := client.RemoteAddr()
	assert.Equal(t, "10.8.7.232:49360", remoteAddr)
}
