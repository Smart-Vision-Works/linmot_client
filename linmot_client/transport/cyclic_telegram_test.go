package transport

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"testing"
	"time"
)

// TestCyclicTelegram_C3Style sends a C#-style cyclic telegram probe (reqDef=6, respDef=0x1FF)
// and captures one response to verify drive connectivity.
// This bypasses the request manager and tests the raw protocol.
func TestCyclicTelegram_C3Style(t *testing.T) {
	t.Skip("hardware-only test; requires a real LinMot drive")

	// Configuration (from LinMot hardware test flags)
	masterPort := 41136
	driveIP := "10.8.7.232"
	drivePort := 49360

	// Create UDP connection bound to master port
	localAddr := net.UDPAddr{
		Port: masterPort,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &localAddr)
	if err != nil {
		t.Fatalf("Failed to bind UDP to port %d: %v", masterPort, err)
	}
	defer conn.Close()

	driveAddr := &net.UDPAddr{
		IP:   net.ParseIP(driveIP),
		Port: drivePort,
	}

	responseCount := 0
	var responses []string

	// Run 10 iterations
	for iteration := 0; iteration < 10; iteration++ {
		// Build 48-byte probe packet
		probe := make([]byte, 48)
		binary.LittleEndian.PutUint32(probe[0:4], 0x00000006) // reqDef = 6
		binary.LittleEndian.PutUint32(probe[4:8], 0x000001FF) // respDef = 0x1FF
		// Rest remains zero

		// Send probe
		_, err := conn.WriteToUDP(probe, driveAddr)
		if err != nil {
			t.Logf("Iteration %d: failed to send probe: %v", iteration, err)
			continue
		}
		t.Logf("Iteration %d: sent 48-byte probe to %s:%d", iteration, driveIP, drivePort)

		// Set read deadline (500ms)
		conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

		// Read ONE packet
		response := make([]byte, 2048)
		n, remoteAddr, err := conn.ReadFromUDP(response)
		conn.SetReadDeadline(time.Time{}) // Clear deadline

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				t.Logf("Iteration %d: read timeout (no response)", iteration)
				continue
			}
			t.Logf("Iteration %d: read error: %v", iteration, err)
			continue
		}

		// Got a response
		responseCount++
		response = response[:n]

		// Parse status word (bytes[8:10]) and status value (bytes[10:12])
		var sw, sv uint16
		if n >= 12 {
			sw = binary.LittleEndian.Uint16(response[8:10])
			sv = binary.LittleEndian.Uint16(response[10:12])
		}

		// Format output
		respHex := hex.EncodeToString(response)
		msg := fmt.Sprintf("Iteration %d: from=%s len=%d sw=0x%04X sv=0x%04X hex=%s",
			iteration, remoteAddr.String(), n, sw, sv, respHex)
		responses = append(responses, msg)
		t.Log(msg)
	}

	// Print summary
	fmt.Printf("\n\n==== C#-STYLE CYCLIC TELEGRAM PROBE RESULTS ====\n")
	fmt.Printf("Responses received: %d / 10\n", responseCount)
	fmt.Printf("Response length expectation: 50 bytes (respDef=0x1FF)\n")
	fmt.Printf("\nDetailed responses:\n")
	for _, resp := range responses {
		fmt.Printf("%s\n", resp)
	}
}

// TestCyclicTelegram_CSharpSteadyState_Req7 sends a C# steady-state cyclic telegram (reqDef=7, 50 bytes, with ControlWord)
// This is the pattern used in the C# library's continuous operational mode.
func TestCyclicTelegram_CSharpSteadyState_Req7(t *testing.T) {
	t.Skip("hardware-only test; requires a real LinMot drive")

	// Configuration
	masterPort := 41136
	driveIP := "10.8.7.232"
	drivePort := 49360

	// Create UDP connection bound to master port
	localAddr := net.UDPAddr{
		Port: masterPort,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &localAddr)
	if err != nil {
		t.Fatalf("Failed to bind UDP to port %d: %v", masterPort, err)
	}
	defer conn.Close()

	driveAddr := &net.UDPAddr{
		IP:   net.ParseIP(driveIP),
		Port: drivePort,
	}

	responseCount := 0
	var responses []string

	// Run 10 iterations with reqDef=7 probe (50 bytes total)
	for iteration := 0; iteration < 10; iteration++ {
		// Build 50-byte probe packet (C# steady-state format)
		probe := make([]byte, 50)
		binary.LittleEndian.PutUint32(probe[0:4], 0x00000007) // reqDef = 7 (DS402 cyclical)
		binary.LittleEndian.PutUint32(probe[4:8], 0x000001FF) // respDef = 0x1FF
		// Bytes 8-9: ControlWord = 0x0006 (safe/no-move DS402-ish value)
		binary.LittleEndian.PutUint16(probe[8:10], 0x0006)
		// Rest remains zero (including reserved bytes 10-11, and all trailing bytes)

		// Send probe
		_, err := conn.WriteToUDP(probe, driveAddr)
		if err != nil {
			t.Logf("Iteration %d: failed to send probe: %v", iteration, err)
			continue
		}
		t.Logf("Iteration %d: sent 50-byte probe (reqDef=7, ControlWord=0x0006) to %s:%d", iteration, driveIP, drivePort)

		// Set read deadline (500ms)
		conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

		// Read ONE packet
		response := make([]byte, 2048)
		n, remoteAddr, err := conn.ReadFromUDP(response)
		conn.SetReadDeadline(time.Time{}) // Clear deadline

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				t.Logf("Iteration %d: read timeout (no response)", iteration)
				continue
			}
			t.Logf("Iteration %d: read error: %v", iteration, err)
			continue
		}

		// Got a response
		responseCount++
		response = response[:n]

		// Parse status word (bytes[8:10]) and status value (bytes[10:12])
		var sw, sv uint16
		if n >= 12 {
			sw = binary.LittleEndian.Uint16(response[8:10])
			sv = binary.LittleEndian.Uint16(response[10:12])
		}

		// Format output
		respHex := hex.EncodeToString(response)
		msg := fmt.Sprintf("Iteration %d: from=%s len=%d sw=0x%04X sv=0x%04X hex=%s",
			iteration, remoteAddr.String(), n, sw, sv, respHex)
		responses = append(responses, msg)
		t.Log(msg)
	}

	// Print summary
	fmt.Printf("\n\n==== C# STEADY-STATE CYCLIC TELEGRAM (reqDef=7) RESULTS ====\n")
	fmt.Printf("Responses received: %d / 10\n", responseCount)
	fmt.Printf("Response length expectation: 50 bytes (respDef=0x1FF)\n")
	fmt.Printf("\nDetailed responses:\n")
	for _, resp := range responses {
		fmt.Printf("%s\n", resp)
	}
}
