package protocol_command_tables

import (
	"encoding/hex"
	"testing"
)

// TestAllocateEntryRequest_WireFormat tests that AllocateEntry requests encode correctly.
// The size parameter must be in Word3 (low 16 bits of value field), not Word4 (high 16 bits).
// This test prevents regression of the word-order bug that caused 0xD1 errors.
func TestAllocateEntryRequest_WireFormat(t *testing.T) {
	tests := []struct {
		name        string
		entryID     uint16
		size        uint16
		wantTail    string // Expected tail (bytes 12-15) in hex, little-endian
		description string
	}{
		{
			name:        "AllocateEntry_1_64",
			entryID:     1,
			size:        64,
			wantTail:    "40000000", // size=64 (0x40) in low 16 bits, 0 in high 16 bits
			description: "size=64 must appear as 0x40 in bytes 12-13 (little-endian: 40 00 00 00)",
		},
		{
			name:        "AllocateEntry_1_96",
			entryID:     1,
			size:        96,
			wantTail:    "60000000", // size=96 (0x60) in low 16 bits
			description: "size=96 must appear as 0x60 in bytes 12-13",
		},
		{
			name:        "AllocateEntry_2_128",
			entryID:     2,
			size:        128,
			wantTail:    "80000000", // size=128 (0x80) in low 16 bits
			description: "size=128 must appear as 0x80 in bytes 12-13",
		},
		{
			name:        "AllocateEntry_10_32",
			entryID:     10,
			size:        32,
			wantTail:    "20000000", // size=32 (0x20) in low 16 bits
			description: "size=32 must appear as 0x20 in bytes 12-13",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := NewAllocateEntryRequest(tt.entryID, tt.size)
			if err != nil {
				t.Fatalf("NewAllocateEntryRequest() error = %v", err)
			}

			// Generate wire packet (using counter=1 for test)
			packet, err := request.WriteRtcPacket(1)
			if err != nil {
				t.Fatalf("WriteRtcPacket() error = %v", err)
			}

			// Extract tail (bytes 12-15, the value field)
			if len(packet) < 16 {
				t.Fatalf("packet too short: got %d bytes, want at least 16", len(packet))
			}
			tail := hex.EncodeToString(packet[12:16])

			if tail != tt.wantTail {
				t.Errorf("Wire format mismatch for %s", tt.description)
				t.Errorf("  entryID=%d, size=%d", tt.entryID, tt.size)
				t.Errorf("  got tail:  %s", tail)
				t.Errorf("  want tail: %s", tt.wantTail)
				t.Errorf("  full packet: %s", hex.EncodeToString(packet))
				t.Errorf("  This indicates size is in the wrong word position (high vs low 16 bits)")
			} else {
				t.Logf("✓ Wire format correct: tail=%s (size=%d in low 16 bits)", tail, tt.size)
			}

			// Verify entryID is correct (bytes 10-11)
			entryIDBytes := hex.EncodeToString(packet[10:12])
			expectedEntryID := hex.EncodeToString([]byte{byte(tt.entryID), byte(tt.entryID >> 8)})
			if entryIDBytes != expectedEntryID {
				t.Errorf("entryID encoding mismatch: got %s, want %s", entryIDBytes, expectedEntryID)
			}
		})
	}
}

// TestAllocateEntryRequest_WordOrderRegression tests that the old buggy encoding is detected.
// This test ensures we don't regress to the bug where size was in Word4 (high 16 bits).
func TestAllocateEntryRequest_WordOrderRegression(t *testing.T) {
	request, err := NewAllocateEntryRequest(1, 64)
	if err != nil {
		t.Fatalf("NewAllocateEntryRequest() error = %v", err)
	}

	packet, err := request.WriteRtcPacket(1)
	if err != nil {
		t.Fatalf("WriteRtcPacket() error = %v", err)
	}

	tail := hex.EncodeToString(packet[12:16])

	// The buggy encoding would have been: 00004000 (size in high 16 bits)
	buggyTail := "00004000"
	if tail == buggyTail {
		t.Errorf("REGRESSION DETECTED: Wire format matches old buggy encoding!")
		t.Errorf("  got:  %s (BUGGY: size in high 16 bits)", tail)
		t.Errorf("  want: 40000000 (CORRECT: size in low 16 bits)")
		t.Errorf("  This would cause 0xD1 errors because drive receives size=0")
	}

	// Correct encoding: 40000000 (size in low 16 bits)
	correctTail := "40000000"
	if tail != correctTail {
		t.Errorf("Wire format incorrect: got %s, want %s", tail, correctTail)
	}
}
