package protocol_vai

import (
	"encoding/binary"
	"testing"

	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_motion_control "github.com/Smart-Vision-Works/staged_robot/protocol/motion_control"
)

// TestVAIGoToPosCommand_RoundTrip tests the Go To Position command
func TestVAIGoToPosCommand_RoundTrip(t *testing.T) {
	// Create command with known parameters
	position := int32(500000)  // 50.0mm
	velocity := uint32(500000) // 0.5 m/s
	accel := uint32(200000)    // 2.0 m/s²
	decel := uint32(200000)    // 2.0 m/s²

	cmd := NewVAIGoToPosCommand(position, velocity, accel, decel)
	cmd.SetCounter(3)

	packet, err := cmd.WritePacket()
	if err != nil {
		t.Fatalf("WritePacket() error = %v", err)
	}

	// Verify packet length (40 bytes total)
	if len(packet) != 40 {
		t.Errorf("expected packet length 40, got %d", len(packet))
	}

	// Verify request bit (Motion Control)
	reqBits := binary.LittleEndian.Uint32(packet[0:4])
	if reqBits != protocol_common.RequestFlags.MotionControl {
		t.Errorf("expected request bits 0x%08X, got 0x%08X",
			protocol_common.RequestFlags.MotionControl, reqBits)
	}

	// Verify header
	header := cmd.Header()
	if header.MasterID != protocol_motion_control.MasterIDs.VAI {
		t.Errorf("expected MasterID 0x%02X, got 0x%02X",
			protocol_motion_control.MasterIDs.VAI, header.MasterID)
	}
	if header.SubID != uint8(SubIDs.GoToPos) {
		t.Errorf("expected SubID 0x%02X, got 0x%02X",
			SubIDs.GoToPos, header.SubID)
	}
	if header.Counter != 3 {
		t.Errorf("expected Counter 3, got %d", header.Counter)
	}

	// Verify parameters in packet
	// Position (words 1-2)
	pos := int32(binary.LittleEndian.Uint32(packet[10:14]))
	if pos != position {
		t.Errorf("expected position %d, got %d", position, pos)
	}

	// Velocity (words 3-4)
	vel := binary.LittleEndian.Uint32(packet[14:18])
	if vel != velocity {
		t.Errorf("expected velocity %d, got %d", velocity, vel)
	}
}

// TestVAIStopCommand tests the Stop command (simpler, no parameters)
func TestVAIStopCommand(t *testing.T) {
	cmd := NewVAIStopCommand()
	cmd.SetCounter(2)

	_, err := cmd.WritePacket()
	if err != nil {
		t.Fatalf("WritePacket() error = %v", err)
	}

	// Verify header
	header := cmd.Header()
	if header.MasterID != protocol_motion_control.MasterIDs.VAI {
		t.Errorf("expected MasterID 0x%02X, got 0x%02X",
			protocol_motion_control.MasterIDs.VAI, header.MasterID)
	}
	if header.SubID != uint8(SubIDs.Stop) {
		t.Errorf("expected SubID 0x%02X, got 0x%02X",
			SubIDs.Stop, header.SubID)
	}
	if header.Counter != 2 {
		t.Errorf("expected Counter 2, got %d", header.Counter)
	}
}

// TestVAIIncrementCommands tests increment commands
func TestVAIIncrementDemPosCommand(t *testing.T) {
	increment := int32(-100000) // -10.0mm (negative increment)
	velocity := uint32(300000)
	accel := uint32(150000)
	decel := uint32(150000)

	cmd := NewVAIIncrementDemPosCommand(increment, velocity, accel, decel)
	cmd.SetCounter(1)

	packet, err := cmd.WritePacket()
	if err != nil {
		t.Fatalf("WritePacket() error = %v", err)
	}

	// Verify header
	header := cmd.Header()
	if header.SubID != uint8(SubIDs.IncrementDemPos) {
		t.Errorf("expected SubID 0x%02X, got 0x%02X",
			SubIDs.IncrementDemPos, header.SubID)
	}

	// Verify negative increment is preserved
	inc := int32(binary.LittleEndian.Uint32(packet[10:14]))
	if inc != increment {
		t.Errorf("expected increment %d, got %d", increment, inc)
	}
}

// TestMultipleCommands tests that different commands have different SubIDs
func TestMultipleCommands(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *protocol_motion_control.MCCommandRequest
		expected uint8
	}{
		{
			name:     "GoToPos",
			cmd:      NewVAIGoToPosCommand(0, 100000, 50000, 50000),
			expected: uint8(SubIDs.GoToPos),
		},
		{
			name:     "IncrementDemPos",
			cmd:      NewVAIIncrementDemPosCommand(1000, 100000, 50000, 50000),
			expected: uint8(SubIDs.IncrementDemPos),
		},
		{
			name:     "IncrementTargetPos",
			cmd:      NewVAIIncrementTargetPosCommand(1000, 100000, 50000, 50000),
			expected: uint8(SubIDs.IncrementTargetPos),
		},
		{
			name:     "Stop",
			cmd:      NewVAIStopCommand(),
			expected: uint8(SubIDs.Stop),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := tt.cmd.Header()
			if header.SubID != tt.expected {
				t.Errorf("expected SubID 0x%02X, got 0x%02X",
					tt.expected, header.SubID)
			}
			if header.MasterID != protocol_motion_control.MasterIDs.VAI {
				t.Errorf("expected MasterID 0x%02X, got 0x%02X",
					protocol_motion_control.MasterIDs.VAI, header.MasterID)
			}
		})
	}
}
