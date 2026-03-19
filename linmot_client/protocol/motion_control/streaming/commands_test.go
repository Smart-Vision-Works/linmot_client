package protocol_streaming

import (
	"testing"

	protocol_motion_control "github.com/Smart-Vision-Works/staged_robot/protocol/motion_control"
)

// TestStreamingCommands tests streaming command construction
func TestStreamingCommands(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *protocol_motion_control.MCCommandRequest
		expected uint8
	}{
		{
			name:     "PStreamSlaveTimestamp",
			cmd:      NewPStreamCommand(100000),
			expected: uint8(SubIDs.PStreamSlaveTimestamp),
		},
		{
			name:     "PVStreamSlaveTimestamp",
			cmd:      NewPVStreamCommand(100000, 50000),
			expected: uint8(SubIDs.PVStreamSlaveTimestamp),
		},
		{
			name:     "PVAStreamSlaveTimestamp",
			cmd:      NewPVAStreamCommand(100000, 50000, 10000),
			expected: uint8(SubIDs.PVAStreamSlaveTimestamp),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cmd.SetCounter(2)

			_, err := tt.cmd.WritePacket()
			if err != nil {
				t.Fatalf("WritePacket() error = %v", err)
			}

			header := tt.cmd.Header()
			if header.MasterID != protocol_motion_control.MasterIDs.Streaming {
				t.Errorf("expected MasterID 0x%02X, got 0x%02X",
					protocol_motion_control.MasterIDs.Streaming, header.MasterID)
			}
			if header.SubID != tt.expected {
				t.Errorf("expected SubID 0x%02X, got 0x%02X",
					tt.expected, header.SubID)
			}
		})
	}
}
