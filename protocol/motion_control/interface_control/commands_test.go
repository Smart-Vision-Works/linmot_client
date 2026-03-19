package protocol_interface_control

import (
	"testing"

	protocol_motion_control "github.com/Smart-Vision-Works/linmot_client/protocol/motion_control"
)

// TestInterfaceControlCommands tests interface control commands
func TestInterfaceControlCommands(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *protocol_motion_control.MCCommandRequest
		expected uint8
	}{
		{
			name:     "WriteInterfaceControlWord",
			cmd:      NewWriteInterfaceControlWordCommand(0x1234),
			expected: uint8(SubIDs.WriteInterfaceControlWord),
		},
		{
			name:     "WriteLiveParameter",
			cmd:      NewWriteLiveParameterCommand(0x1000, 0x12345678),
			expected: uint8(SubIDs.WriteLiveParameter),
		},
		{
			name:     "WriteX4IntfOutputsWithMask",
			cmd:      NewWriteX4IntfOutputsWithMaskCommand(0xFF, 0x0F),
			expected: uint8(SubIDs.WriteX4IntfOutputsWithMask),
		},
		{
			name:     "Reset",
			cmd:      NewResetCommand(),
			expected: uint8(SubIDs.Reset),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cmd.SetCounter(1)

			_, err := tt.cmd.WritePacket()
			if err != nil {
				t.Fatalf("WritePacket() error = %v", err)
			}

			header := tt.cmd.Header()
			if header.MasterID != protocol_motion_control.MasterIDs.InterfaceControl {
				t.Errorf("expected MasterID 0x%02X, got 0x%02X",
					protocol_motion_control.MasterIDs.InterfaceControl, header.MasterID)
			}
			if header.SubID != tt.expected {
				t.Errorf("expected SubID 0x%02X, got 0x%02X",
					tt.expected, header.SubID)
			}
		})
	}
}
