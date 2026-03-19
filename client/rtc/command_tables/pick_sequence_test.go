package client_command_tables

import (
	protocol_command_tables "gsail-go/linmot/protocol/rtc/command_tables"
	"testing"
)

// TestPickSequence_Load verifies that the pick_sequence.yaml file loads correctly
// and contains all expected entries with correct parameters.
func TestPickSequence_Load(t *testing.T) {
	ct, _ := loadExampleManifest(t, "pick_sequence.yaml")

	// Bind template variables used in pick_sequence.yaml
	ct.SetVar("POSITION_DOWN", 500000)
	ct.SetVar("POSITION_UP", 100000)
	ct.SetVar("MAX_VELOCITY", 100000)
	ct.SetVar("ACCELERATION", 100000)
	ct.SetVar("DECELERATION", 100000)
	ct.SetVar("DELAY_AT_BOTTOM", 1000)
	ct.SetVar("DELAY_PURGE", 5000)

	// Verify metadata
	if ct.Version != "1" {
		t.Errorf("Version = %q, want %q", ct.Version, "1")
	}
	if ct.DriveModel != "C1250-MI" {
		t.Errorf("DriveModel = %q, want %q", ct.DriveModel, "C1250-MI")
	}

	// Verify entry count
	if len(ct.Entries) != 11 {
		t.Fatalf("Entries count = %d, want 11", len(ct.Entries))
	}

	// Verify each entry in detail
	entries := []struct {
		id                     int
		name                   string
		typ                    string
		par1, par2, par3, par4 *Param
		sequencedEntry         *int
	}{
		{
			id:             1,
			name:           "von/poff",
			typ:            "SetDO",
			par1:           int64Ptr(0x0003),
			par2:           int64Ptr(0x0002),
			sequencedEntry: intPtrInt(2),
		},
		{
			id:             2,
			name:           "move down",
			typ:            "VAI_GoToPos",
			par1:           varPtr("POSITION_DOWN"),
			par2:           varPtr("MAX_VELOCITY"),
			par3:           varPtr("ACCELERATION"),
			par4:           varPtr("DECELERATION"),
			sequencedEntry: intPtrInt(3),
		},
		{
			id:             3,
			name:           "wait for motion",
			typ:            "WaitDemandVelLT",
			par1:           int64Ptr(0),
			sequencedEntry: intPtrInt(4),
		},
		{
			id:             4,
			name:           "wait at bottom",
			typ:            "Delay",
			par1:           varPtr("DELAY_AT_BOTTOM"),
			sequencedEntry: intPtrInt(5),
		},
		{
			id:             5,
			name:           "move up",
			typ:            "VAI_GoToPos",
			par1:           varPtr("POSITION_UP"),
			par2:           varPtr("MAX_VELOCITY"),
			par3:           varPtr("ACCELERATION"),
			par4:           varPtr("DECELERATION"),
			sequencedEntry: intPtrInt(6),
		},
		{
			id:             6,
			name:           "wait for motion",
			typ:            "WaitDemandVelLT",
			par1:           int64Ptr(0),
			sequencedEntry: intPtrInt(7),
		},
		{
			id:   7,
			name: "trigger check",
			typ:  "IfMaskedX4Eq",
			par1: int64Ptr(0x0002),
			par2: int64Ptr(0x0002),
			par3: int64Ptr(8),
			par4: int64Ptr(9),
		},
		{
			id:             8,
			name:           "falling edge",
			typ:            "WaitFalling",
			sequencedEntry: intPtrInt(9),
		},
		{
			id:             9,
			name:           "voff/pon",
			typ:            "ClearDO",
			par1:           int64Ptr(0x0003),
			par2:           int64Ptr(0x0001),
			sequencedEntry: intPtrInt(10),
		},
		{
			id:             10,
			name:           "wait",
			typ:            "Delay",
			par1:           varPtr("DELAY_PURGE"),
			sequencedEntry: intPtrInt(11),
		},
		{
			id:   11,
			name: "off",
			typ:  "ClearDO",
			par1: int64Ptr(0x0003),
			par2: int64Ptr(0x0000),
		},
	}

	for i, expected := range entries {
		entry := ct.Entries[i]

		// Verify basic fields
		if entry.ID != expected.id {
			t.Errorf("Entry %d: ID = %d, want %d", i, entry.ID, expected.id)
		}
		if entry.Name != expected.name {
			t.Errorf("Entry %d: Name = %q, want %q", i, entry.Name, expected.name)
		}
		if entry.Type != expected.typ {
			t.Errorf("Entry %d: Type = %q, want %q", i, entry.Type, expected.typ)
		}

		// Verify parameters
		verifyParameter(t, i, "Par1", entry.Par1, expected.par1)
		verifyParameter(t, i, "Par2", entry.Par2, expected.par2)
		verifyParameter(t, i, "Par3", entry.Par3, expected.par3)
		verifyParameter(t, i, "Par4", entry.Par4, expected.par4)

		// Verify sequenced entry
		if (entry.SequencedEntry == nil) != (expected.sequencedEntry == nil) {
			t.Errorf("Entry %d: SequencedEntry nil mismatch: got %v, want %v",
				i, entry.SequencedEntry, expected.sequencedEntry)
		} else if entry.SequencedEntry != nil && expected.sequencedEntry != nil {
			if *entry.SequencedEntry != *expected.sequencedEntry {
				t.Errorf("Entry %d: SequencedEntry = %d, want %d",
					i, *entry.SequencedEntry, *expected.sequencedEntry)
			}
		}
	}
}

// TestPickSequence_Encoding verifies that all entries can be encoded to 64-byte binary format
func TestPickSequence_Encoding(t *testing.T) {
	ct, _ := loadExampleManifest(t, "pick_sequence.yaml")

	// Bind template variables used in pick_sequence.yaml
	ct.SetVar("POSITION_DOWN", 500000)
	ct.SetVar("POSITION_UP", 100000)
	ct.SetVar("MAX_VELOCITY", 100000)
	ct.SetVar("ACCELERATION", 100000)
	ct.SetVar("DECELERATION", 100000)
	ct.SetVar("DELAY_AT_BOTTOM", 1000)
	ct.SetVar("DELAY_PURGE", 5000)

	for i, entry := range ct.Entries {
		encoded, err := entry.Encode(ct)
		if err != nil {
			t.Fatalf("Entry %d (ID=%d, Type=%s) Encode() failed: %v",
				i, entry.ID, entry.Type, err)
		}

		// Verify length
		if len(encoded) != 64 {
			t.Errorf("Entry %d encoded length = %d, want 64", i, len(encoded))
		}

		// Verify A701h version header (0x01 0xA7 in little-endian)
		if encoded[0] != 0x01 || encoded[1] != 0xA7 {
			t.Errorf("Entry %d missing A701h header: got %02x %02x", i, encoded[0], encoded[1])
		}

		// Verify name encoding (should be NUL-terminated within 16 bytes)
		nameBytes := encoded[38:54]
		hasNul := false
		for _, b := range nameBytes {
			if b == 0 {
				hasNul = true
				break
			}
		}
		if !hasNul {
			t.Errorf("Entry %d name not NUL-terminated", i)
		}
	}
}

// TestPickSequence_RoundTrip verifies that entries can be encoded and decoded
func TestPickSequence_RoundTrip(t *testing.T) {
	ct, _ := loadExampleManifest(t, "pick_sequence.yaml")

	// Bind template variables used in pick_sequence.yaml
	ct.SetVar("POSITION_DOWN", 500000)
	ct.SetVar("POSITION_UP", 100000)
	ct.SetVar("MAX_VELOCITY", 100000)
	ct.SetVar("ACCELERATION", 100000)
	ct.SetVar("DECELERATION", 100000)
	ct.SetVar("DELAY_AT_BOTTOM", 1000)
	ct.SetVar("DELAY_PURGE", 5000)

	for i, entry := range ct.Entries {
		// Encode
		encoded, err := entry.Encode(ct)
		if err != nil {
			t.Fatalf("Entry %d Encode() failed: %v", i, err)
		}

		// Decode
		decoded, err := protocol_command_tables.DecodeEntry(encoded)
		if err != nil {
			t.Fatalf("Entry %d DecodeEntry() failed: %v", i, err)
		}

		// Verify type matches
		if decoded.Type != entry.Type {
			t.Errorf("Entry %d: decoded Type = %q, want %q", i, decoded.Type, entry.Type)
		}

		// Verify name matches
		if decoded.Name != entry.Name {
			t.Errorf("Entry %d: decoded Name = %q, want %q", i, decoded.Name, entry.Name)
		}

		// Verify sequenced entry matches
		if (decoded.SequencedEntry == nil) != (entry.SequencedEntry == nil) {
			t.Errorf("Entry %d: SequencedEntry nil mismatch", i)
		} else if decoded.SequencedEntry != nil && entry.SequencedEntry != nil {
			if int(*decoded.SequencedEntry) != *entry.SequencedEntry {
				t.Errorf("Entry %d: decoded SequencedEntry = %d, want %d",
					i, *decoded.SequencedEntry, *entry.SequencedEntry)
			}
		}
	}
}

// TestPickSequence_SequencedEntryChain verifies the sequenced entry chain is correct
func TestPickSequence_SequencedEntryChain(t *testing.T) {
	ct, _ := loadExampleManifest(t, "pick_sequence.yaml")

	// Expected chain: 1→2→3→4→5→6→7 (no seq), 8→9→10→11 (no seq)
	expectedChains := [][]int{
		{1, 2, 3, 4, 5, 6, 7}, // Main sequence
		{8, 9, 10, 11},        // Suction/purge sequence
	}

	// Build ID to entry map
	entryMap := make(map[int]Entry)
	for _, entry := range ct.Entries {
		entryMap[entry.ID] = entry
	}

	// Verify each chain
	for chainIdx, chain := range expectedChains {
		for i := 0; i < len(chain)-1; i++ {
			currentID := chain[i]
			expectedNextID := chain[i+1]

			entry, ok := entryMap[currentID]
			if !ok {
				t.Fatalf("Chain %d: Entry ID %d not found", chainIdx, currentID)
			}

			if entry.SequencedEntry == nil {
				t.Errorf("Chain %d: Entry %d missing SequencedEntry", chainIdx, currentID)
				continue
			}

			if *entry.SequencedEntry != expectedNextID {
				t.Errorf("Chain %d: Entry %d SequencedEntry = %d, want %d",
					chainIdx, currentID, *entry.SequencedEntry, expectedNextID)
			}
		}

		// Verify last entry in chain has no sequenced entry (except for entry 7 which branches)
		lastID := chain[len(chain)-1]
		if lastID != 7 { // Entry 7 has no sequenced entry (it's an IF statement)
			entry := entryMap[lastID]
			if entry.SequencedEntry != nil {
				t.Errorf("Chain %d: Last entry %d should have no SequencedEntry, got %d",
					chainIdx, lastID, *entry.SequencedEntry)
			}
		}
	}
}

// TestPickSequence_IfBranchIntegrity verifies IF statement branches reference valid entries
func TestPickSequence_IfBranchIntegrity(t *testing.T) {
	ct, _ := loadExampleManifest(t, "pick_sequence.yaml")

	// Entry 7 is an IfMaskedX4Eq with par3=8 (true) and par4=9 (false)
	entry7 := ct.Entries[6] // Index 6 = ID 7
	if entry7.ID != 7 {
		t.Fatalf("Expected entry 7 at index 6, got ID %d", entry7.ID)
	}

	if entry7.Type != "IfMaskedX4Eq" {
		t.Fatalf("Entry 7 Type = %q, want IfMaskedX4Eq", entry7.Type)
	}

	// Verify true branch (par3) points to entry 8
	if entry7.Par3 == nil || entry7.Par3.Literal == nil || *entry7.Par3.Literal != 8 {
		t.Errorf("Entry 7 Par3 (true branch) = %v, want 8", entry7.Par3)
	}

	// Verify false branch (par4) points to entry 9
	if entry7.Par4 == nil || entry7.Par4.Literal == nil || *entry7.Par4.Literal != 9 {
		t.Errorf("Entry 7 Par4 (false branch) = %v, want 9", entry7.Par4)
	}

	// Verify both entries 8 and 9 exist
	entry8Exists := false
	entry9Exists := false
	for _, entry := range ct.Entries {
		if entry.ID == 8 {
			entry8Exists = true
		}
		if entry.ID == 9 {
			entry9Exists = true
		}
	}

	if !entry8Exists {
		t.Error("Entry 8 (Begin suction) not found - IF true branch broken")
	}
	if !entry9Exists {
		t.Error("Entry 9 (Stop suction) not found - IF false branch broken")
	}
}

// TestPickSequence_MotionParameters verifies motion command parameters are in valid ranges
func TestPickSequence_MotionParameters(t *testing.T) {
	ct, _ := loadExampleManifest(t, "pick_sequence.yaml")

	// Bind template variables to check resolved values
	ct.SetVar("POSITION_DOWN", 500000)
	ct.SetVar("POSITION_UP", 100000)
	ct.SetVar("MAX_VELOCITY", 100000)
	ct.SetVar("ACCELERATION", 100000)
	ct.SetVar("DECELERATION", 100000)
	ct.SetVar("DELAY_AT_BOTTOM", 1000)
	ct.SetVar("DELAY_PURGE", 5000)

	for i, entry := range ct.Entries {
		if entry.Type == "VAI_GoToPos" {
			// Verify position is reasonable (-2^31 to 2^31-1 in 0.1µm units)
			if entry.Par1 == nil {
				t.Errorf("Entry %d (VAI_GoToPos): Par1 (position) is nil", i)
			} else {
				pos := ct.resolveParam(entry.Par1)
				if pos != nil {
					if *pos < -2147483648 || *pos > 2147483647 {
						t.Errorf("Entry %d (VAI_GoToPos): Par1 (position) out of range: %d", i, *pos)
					}
				}
			}

			// Verify velocity is positive
			if entry.Par2 == nil {
				t.Errorf("Entry %d (VAI_GoToPos): Par2 (velocity) is nil", i)
			} else {
				vel := ct.resolveParam(entry.Par2)
				if vel != nil && *vel <= 0 {
					t.Errorf("Entry %d (VAI_GoToPos): Par2 (velocity) must be positive: %d", i, *vel)
				}
			}

			// Verify acceleration is positive
			if entry.Par3 == nil {
				t.Errorf("Entry %d (VAI_GoToPos): Par3 (acceleration) is nil", i)
			} else {
				acc := ct.resolveParam(entry.Par3)
				if acc != nil && *acc <= 0 {
					t.Errorf("Entry %d (VAI_GoToPos): Par3 (acceleration) must be positive: %d", i, *acc)
				}
			}

			// Verify deceleration is positive
			if entry.Par4 == nil {
				t.Errorf("Entry %d (VAI_GoToPos): Par4 (deceleration) is nil", i)
			} else {
				dec := ct.resolveParam(entry.Par4)
				if dec != nil && *dec <= 0 {
					t.Errorf("Entry %d (VAI_GoToPos): Par4 (deceleration) must be positive: %d", i, *dec)
				}
			}
		}
	}
}

// TestPickSequence_DigitalOutputs verifies digital output commands use valid bit masks
func TestPickSequence_DigitalOutputs(t *testing.T) {
	ct, _ := loadExampleManifest(t, "pick_sequence.yaml")

	for i, entry := range ct.Entries {
		if entry.Type == "SetDO" {
			// Verify mask is set
			if entry.Par1 == nil || entry.Par1.Literal == nil {
				t.Errorf("Entry %d (SetDO): Par1 (mask) is nil", i)
				continue
			}

			// Verify mask is non-zero
			if *entry.Par1.Literal == 0 {
				t.Errorf("Entry %d (SetDO): Par1 (mask) is zero", i)
			}

			// Verify value is within mask
			if entry.Par2 != nil && entry.Par2.Literal != nil {
				mask := uint16(*entry.Par1.Literal)
				value := uint16(*entry.Par2.Literal)
				if (value & ^mask) != 0 {
					t.Errorf("Entry %d (SetDO): Par2 (value) has bits outside mask: mask=0x%04x, value=0x%04x",
						i, mask, value)
				}
			}
		}
	}
}

// Helper function to verify parameter values
func verifyParameter(t *testing.T, entryIdx int, name string, got, want *Param) {
	t.Helper()
	if (got == nil) != (want == nil) {
		t.Errorf("Entry %d: %s nil mismatch: got %v, want %v", entryIdx, name, got, want)
		return
	}
	if got != nil && want != nil {
		// Compare VarName if both have them
		if got.VarName != "" || want.VarName != "" {
			if got.VarName != want.VarName {
				t.Errorf("Entry %d: %s VarName = %q, want %q", entryIdx, name, got.VarName, want.VarName)
			}
			return
		}
		// Compare Literal if both have them
		if got.Literal == nil || want.Literal == nil {
			if got.Literal != want.Literal {
				t.Errorf("Entry %d: %s Literal nil mismatch: got %v, want %v", entryIdx, name, got.Literal, want.Literal)
			}
		} else if *got.Literal != *want.Literal {
			t.Errorf("Entry %d: %s = %d, want %d", entryIdx, name, *got.Literal, *want.Literal)
		}
	}
}

// Helper to create int64 pointer
func int64Ptr(i int64) *Param {
	return &Param{Literal: &i}
}

// Helper to create variable reference
func varPtr(name string) *Param {
	return &Param{VarName: name}
}

// Note: intPtrInt is defined in command_table_test.go
