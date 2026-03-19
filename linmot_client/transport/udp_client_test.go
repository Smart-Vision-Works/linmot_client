package transport

import (
	"net"
	"strings"
	"testing"
)

func TestNewUDPTransportClient(t *testing.T) {
	tests := []struct {
		name          string
		driveIP       string
		expectError   bool
		errorContains string
	}{
		{
			name:          "Empty IP address",
			driveIP:       "",
			expectError:   true,
			errorContains: "cannot be empty",
		},
		{
			name:        "Valid IP address",
			driveIP:     "127.0.0.1",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport, err := NewUDPTransportClient(tt.driveIP, DefaultDrivePort, DefaultMasterPort, "", DefaultTimeout)

			if tt.expectError {
				if err == nil {
					t.Fatal("NewUDPTransportClient() expected error but got none")
				}
				if transport != nil {
					t.Fatal("NewUDPTransportClient() expected nil transport on error")
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("NewUDPTransportClient() error = %q, want error containing %q", err.Error(), tt.errorContains)
				}
			} else {
				if err != nil {
					t.Fatalf("NewUDPTransportClient() unexpected error: %v", err)
				}
				if transport == nil {
					t.Fatal("NewUDPTransportClient() returned nil transport")
				}
				defer transport.Close()

				// Verify transport is configured correctly
				if transport.driveIP != tt.driveIP {
					t.Errorf("Transport driveIP = %q, want %q", transport.driveIP, tt.driveIP)
				}
				if transport.drivePort != DefaultDrivePort {
					t.Errorf("Transport drivePort = %d, want %d", transport.drivePort, DefaultDrivePort)
				}
				if transport.masterPort != DefaultMasterPort {
					t.Errorf("Transport masterPort = %d, want %d", transport.masterPort, DefaultMasterPort)
				}
				if transport.connection == nil {
					t.Error("Transport connection is nil")
				}
			}
		})
	}
}

func TestNewUDPTransportClient_WithCustomPorts(t *testing.T) {
	tests := []struct {
		name        string
		driveIP     string
		drivePort   int
		masterPort  int
		expectError bool
	}{
		{
			name:        "Custom drive port",
			driveIP:     "192.168.1.100",
			drivePort:   50000,
			masterPort:  DefaultMasterPort,
			expectError: false,
		},
		{
			name:        "Custom master port",
			driveIP:     "192.168.1.100",
			drivePort:   DefaultDrivePort,
			masterPort:  41150,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport, err := NewUDPTransportClient(tt.driveIP, tt.drivePort, tt.masterPort, "", DefaultTimeout)

			if tt.expectError {
				if err == nil {
					t.Fatal("NewUDPTransportClient() expected error but got none")
				}
				if transport != nil {
					t.Fatal("NewUDPTransportClient() expected nil transport on error")
				}
			} else {
				if err != nil {
					t.Fatalf("NewUDPTransportClient() unexpected error: %v", err)
				}
				if transport == nil {
					t.Fatal("NewUDPTransportClient() returned nil transport")
				}
				defer transport.Close()

				// Verify transport is configured correctly
				if transport.driveIP != tt.driveIP {
					t.Errorf("Transport driveIP = %q, want %q", transport.driveIP, tt.driveIP)
				}
				if transport.drivePort != tt.drivePort {
					t.Errorf("Transport drivePort = %d, want %d", transport.drivePort, tt.drivePort)
				}
				if transport.masterPort != tt.masterPort {
					t.Errorf("Transport masterPort = %d, want %d", transport.masterPort, tt.masterPort)
				}
				if transport.connection == nil {
					t.Error("Transport connection is nil")
				}

				// Test that we can create a client with this transport
				// Note: NewUDPClient convenience wrapper is in root package linmot
				// This test only verifies the transport layer creation
			}
		})
	}
}

func TestCreateUDPConnection(t *testing.T) {
	tests := []struct {
		name        string
		masterPort  int
		expectError bool
	}{
		{
			name:        "Any available port",
			masterPort:  0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := createUDPConnection(net.ParseIP("0.0.0.0"), tt.masterPort)

			if tt.expectError {
				if err == nil {
					t.Fatal("createUDPConnection() expected error but got none")
				}
				if conn != nil {
					t.Fatal("createUDPConnection() expected nil connection on error")
					_ = conn.Close()
				}
			} else {
				if err != nil {
					t.Fatalf("createUDPConnection() unexpected error: %v", err)
				}
				if conn == nil {
					t.Fatal("createUDPConnection() returned nil connection")
				}
				defer conn.Close()

				// Verify connection is valid
				localAddr := conn.LocalAddr()
				if localAddr == nil {
					t.Error("Connection LocalAddr() is nil")
				}
			}
		})
	}
}

func TestCreateUDPConnection_PortInUse(t *testing.T) {
	// Test that createUDPConnection properly fails when the requested port is in use.
	// LinMot requires packets from the exact master port, so falling back to a random port
	// would break communication. The correct solution is to use SharedUDPTransport.

	// First, bind to a specific port
	firstConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 50002})
	if err != nil {
		t.Fatalf("Failed to create first connection: %v", err)
	}
	defer firstConn.Close()

	// Now try to create a connection on the same port - should fail
	secondConn, err := createUDPConnection(net.ParseIP("0.0.0.0"), 50002)
	if err == nil {
		secondConn.Close()
		t.Fatal("createUDPConnection() should fail when port is in use, but succeeded")
	}

	// Verify error message mentions the port conflict
	if !strings.Contains(err.Error(), "bind") && !strings.Contains(err.Error(), "address already in use") {
		t.Errorf("Error should mention port binding issue, got: %v", err)
	}
}

func TestUDPTransportClient_BuildRTCWritePacket(t *testing.T) {
	t.Skip("Packet building is now handled by LinUDPV2Protocol layer, tested in linudp_v2_protocol_test.go")
}

func TestUDPTransportClient_ConcurrentAccess(t *testing.T) {
	t.Skip("RTC counter thread safety is now tested at the protocol layer in linudp_v2_protocol_test.go")
}
