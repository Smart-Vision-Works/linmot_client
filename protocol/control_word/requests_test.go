package protocol_control_word

import (
	"testing"
)

func TestControlWordRequest_WritePacket(t *testing.T) {
	tests := []struct {
		name     string
		word     uint16
		expected []byte
	}{
		{
			name: "EnableDrive pattern",
			word: EnableDrivePattern(),
			expected: []byte{
				0x01, 0x00, 0x00, 0x00, // Request Definition: bit 0 set (Control Word)
				0x01, 0x00, 0x00, 0x00, // Response Definition: bit 0 set (Status Word)
				0x0F, 0x00, // Control Word: 0x000F (bits 0-3 set)
			},
		},
		{
			name: "DisableDrive pattern",
			word: DisableDrivePattern(),
			expected: []byte{
				0x01, 0x00, 0x00, 0x00, // Request Definition
				0x01, 0x00, 0x00, 0x00, // Response Definition
				0x00, 0x00, // Control Word: 0x0000
			},
		},
		{
			name: "ErrorAcknowledge pattern",
			word: ErrorAcknowledgePattern(),
			expected: []byte{
				0x01, 0x00, 0x00, 0x00, // Request Definition
				0x01, 0x00, 0x00, 0x00, // Response Definition
				0x80, 0x00, // Control Word: 0x0080 (bit 7 set)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := NewControlWordRequest(tt.word)
			packet, err := request.WritePacket()
			if err != nil {
				t.Fatalf("WritePacket() error = %v", err)
			}

			if len(packet) != 10 {
				t.Errorf("expected packet length 10, got %d", len(packet))
			}

			for i, expectedByte := range tt.expected {
				if packet[i] != expectedByte {
					t.Errorf("byte %d: expected 0x%02X, got 0x%02X", i, expectedByte, packet[i])
				}
			}

			if request.GetControlWord() != tt.word {
				t.Errorf("GetControlWord() = 0x%04X, want 0x%04X", request.GetControlWord(), tt.word)
			}
		})
	}
}

func TestControlWordRequest_SetBit(t *testing.T) {
	request := NewControlWordRequest(0x0000)

	request.SetBit(ControlWordBit_SwitchOn)
	if request.GetControlWord() != 0x0001 {
		t.Errorf("After SetBit(0), expected 0x0001, got 0x%04X", request.GetControlWord())
	}

	request.SetBit(ControlWordBit_EnableVoltage)
	if request.GetControlWord() != 0x0003 {
		t.Errorf("After SetBit(1), expected 0x0003, got 0x%04X", request.GetControlWord())
	}

	request.SetBit(ControlWordBit_ErrorAcknowledge)
	if request.GetControlWord() != 0x0083 {
		t.Errorf("After SetBit(7), expected 0x0083, got 0x%04X", request.GetControlWord())
	}
}

func TestControlWordRequest_ClearBit(t *testing.T) {
	request := NewControlWordRequest(0xFFFF)

	request.ClearBit(ControlWordBit_SwitchOn)
	if request.GetControlWord() != 0xFFFE {
		t.Errorf("After ClearBit(0), expected 0xFFFE, got 0x%04X", request.GetControlWord())
	}

	request.ClearBit(ControlWordBit_EnableVoltage)
	if request.GetControlWord() != 0xFFFC {
		t.Errorf("After ClearBit(1), expected 0xFFFC, got 0x%04X", request.GetControlWord())
	}

	request.ClearBit(ControlWordBit_ErrorAcknowledge)
	if request.GetControlWord() != 0xFF7C {
		t.Errorf("After ClearBit(7), expected 0xFF7C, got 0x%04X", request.GetControlWord())
	}
}

func TestControlWordRequest_WithPattern(t *testing.T) {
	request := NewControlWordRequest(0x0000)

	request.WithPattern(0x1234)
	if request.GetControlWord() != 0x1234 {
		t.Errorf("After WithPattern(0x1234), expected 0x1234, got 0x%04X", request.GetControlWord())
	}

	request.WithPattern(0xABCD)
	if request.GetControlWord() != 0xABCD {
		t.Errorf("After WithPattern(0xABCD), expected 0xABCD, got 0x%04X", request.GetControlWord())
	}
}

func TestRoundTrip_ControlWordRequest(t *testing.T) {
	tests := []struct {
		name string
		word uint16
	}{
		{"Zero", 0x0000},
		{"EnableDrive", EnableDrivePattern()},
		{"ErrorAcknowledge", ErrorAcknowledgePattern()},
		{"AllBitsSet", 0xFFFF},
		{"SomeBits", 0x5555},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write request to packet
			request := NewControlWordRequest(tt.word)
			packet, err := request.WritePacket()
			if err != nil {
				t.Fatalf("WritePacket() error = %v", err)
			}

			// Read request back from packet
			requestRead, err := ReadControlWordRequest(packet)
			if err != nil {
				t.Fatalf("ReadControlWordRequest() error = %v", err)
			}

			// Compare
			if requestRead.GetControlWord() != request.GetControlWord() {
				t.Errorf("Round trip failed: wrote 0x%04X, read 0x%04X",
					request.GetControlWord(), requestRead.GetControlWord())
			}
		})
	}
}

func TestControlWordBuilder_Patterns(t *testing.T) {
	t.Run("EnableDrive", func(t *testing.T) {
		word := NewControlWordBuilder().
			SwitchOn().
			EnableVoltage().
			ReleaseQuickStop().
			EnableOperation().
			Build()

		expected := EnableDrivePattern()
		if word != expected {
			t.Errorf("EnableDrive pattern: expected 0x%04X, got 0x%04X", expected, word)
		}
	})

	t.Run("ErrorAcknowledge", func(t *testing.T) {
		word := NewControlWordBuilder().
			AcknowledgeError().
			Build()

		expected := ErrorAcknowledgePattern()
		if word != expected {
			t.Errorf("ErrorAcknowledge pattern: expected 0x%04X, got 0x%04X", expected, word)
		}
	})

	t.Run("QuickStop", func(t *testing.T) {
		word := NewControlWordBuilder().
			QuickStop().
			Build()

		if IsBitSet(word, ControlWordBit_QuickStop) {
			t.Errorf("QuickStop should clear bit 2, but it is set")
		}
	})

	t.Run("JogPlus", func(t *testing.T) {
		word := NewControlWordBuilder().
			JogPlus().
			Build()

		if !IsBitSet(word, ControlWordBit_JogPlus) {
			t.Errorf("JogPlus should set bit 8")
		}
	})

	t.Run("JogMinus", func(t *testing.T) {
		word := NewControlWordBuilder().
			JogMinus().
			Build()

		if !IsBitSet(word, ControlWordBit_JogMinus) {
			t.Errorf("JogMinus should set bit 9")
		}
	})

	t.Run("StopJog", func(t *testing.T) {
		word := NewControlWordBuilder().
			JogPlus().
			JogMinus().
			StopJog().
			Build()

		if IsBitSet(word, ControlWordBit_JogPlus) || IsBitSet(word, ControlWordBit_JogMinus) {
			t.Errorf("StopJog should clear both jog bits")
		}
	})

	t.Run("ComplexPattern", func(t *testing.T) {
		word := NewControlWordBuilder().
			SwitchOn().
			EnableVoltage().
			ReleaseQuickStop().
			EnableOperation().
			Home().
			Build()

		// Should have bits 0, 1, 2, 3, 6 set
		expectedBits := []uint{
			ControlWordBit_SwitchOn,
			ControlWordBit_EnableVoltage,
			ControlWordBit_QuickStop,
			ControlWordBit_EnableOperation,
			ControlWordBit_Home,
		}

		for _, bit := range expectedBits {
			if !IsBitSet(word, bit) {
				t.Errorf("Bit %d should be set", bit)
			}
		}
	})
}

func TestBitHelpers(t *testing.T) {
	t.Run("SetBit", func(t *testing.T) {
		word := uint16(0x0000)
		word = SetBit(word, 0)
		if word != 0x0001 {
			t.Errorf("SetBit(0) failed: expected 0x0001, got 0x%04X", word)
		}

		word = SetBit(word, 15)
		if word != 0x8001 {
			t.Errorf("SetBit(15) failed: expected 0x8001, got 0x%04X", word)
		}
	})

	t.Run("ClearBit", func(t *testing.T) {
		word := uint16(0xFFFF)
		word = ClearBit(word, 0)
		if word != 0xFFFE {
			t.Errorf("ClearBit(0) failed: expected 0xFFFE, got 0x%04X", word)
		}

		word = ClearBit(word, 15)
		if word != 0x7FFE {
			t.Errorf("ClearBit(15) failed: expected 0x7FFE, got 0x%04X", word)
		}
	})

	t.Run("IsBitSet", func(t *testing.T) {
		word := uint16(0x0081) // bits 0 and 7 set

		if !IsBitSet(word, 0) {
			t.Error("Bit 0 should be set")
		}
		if IsBitSet(word, 1) {
			t.Error("Bit 1 should not be set")
		}
		if !IsBitSet(word, 7) {
			t.Error("Bit 7 should be set")
		}
		if IsBitSet(word, 15) {
			t.Error("Bit 15 should not be set")
		}
	})
}
