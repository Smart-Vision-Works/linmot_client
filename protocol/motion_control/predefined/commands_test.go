package protocol_predefined

import (
	"testing"

	protocol_motion_control "github.com/Smart-Vision-Works/linmot_client/protocol/motion_control"
)

// TestPredefinedCommands tests multiple predefined VAI commands
func TestPredefinedCommands(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *protocol_motion_control.MCCommandRequest
		expected uint8
	}{
		{
			name:     "GoToPos",
			cmd:      NewPredefVAIGoToPosCommand(500000),
			expected: uint8(SubIDs.GoToPos),
		},
		{
			name:     "IncrementDemPos",
			cmd:      NewPredefVAIIncrementDemPosCommand(10000),
			expected: uint8(SubIDs.IncrementDemPos),
		},
		{
			name:     "GoToPosFromActPos",
			cmd:      NewPredefVAIGoToPosFromActPosCommand(500000),
			expected: uint8(SubIDs.GoToPosFromActPosAndActVel),
		},
		{
			name:     "Stop",
			cmd:      NewPredefVAIStopCommand(),
			expected: uint8(SubIDs.Stop),
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
			if header.MasterID != protocol_motion_control.MasterIDs.PredefVAI {
				t.Errorf("expected MasterID 0x%02X, got 0x%02X",
					protocol_motion_control.MasterIDs.PredefVAI, header.MasterID)
			}
			if header.SubID != tt.expected {
				t.Errorf("expected SubID 0x%02X, got 0x%02X",
					tt.expected, header.SubID)
			}
			if header.Counter != 2 {
				t.Errorf("expected Counter 2, got %d", header.Counter)
			}
		})
	}
}
