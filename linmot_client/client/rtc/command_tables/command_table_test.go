package client_command_tables

import (
	protocol_command_tables "github.com/Smart-Vision-Works/staged_robot/protocol/rtc/command_tables"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "bm_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	return f.Name()
}

// loadExampleManifest loads an example YAML file from the testdata directory.
// It returns the parsed CommandTable and the full path to the file.
func loadExampleManifest(t *testing.T, yamlFilename string) (*CommandTable, string) {
	t.Helper()
	// Get the path to the testdata directory
	_, filename, _, _ := runtime.Caller(0)
	testdataDir := filepath.Join(filepath.Dir(filename), "testdata")
	yamlPath := filepath.Join(testdataDir, yamlFilename)

	ct, err := Load(yamlPath)
	if err != nil {
		t.Fatalf("Load(%q) failed: %v", yamlFilename, err)
	}

	return ct, yamlPath
}

// TestCommandTable_BuildPickPlace tests loading the example test_command_table.yaml file
func TestCommandTable_BuildPickPlace(t *testing.T) {
	ct, _ := loadExampleManifest(t, "test_command_table.yaml")

	// Verify basic structure
	if ct.Version != "1" {
		t.Errorf("Version = %q, want %q", ct.Version, "1")
	}
	if ct.DriveModel != "C1250-MI" {
		t.Errorf("DriveModel = %q, want %q", ct.DriveModel, "C1250-MI")
	}
	if len(ct.Entries) != 6 {
		t.Errorf("Entries count = %d, want %d", len(ct.Entries), 6)
	}

	// Verify first entry (TrigSel - ParamWrite)
	if ct.Entries[0].ID != 1 {
		t.Errorf("Entry 0 ID = %d, want %d", ct.Entries[0].ID, 1)
	}
	if ct.Entries[0].Name != "TrigSel" {
		t.Errorf("Entry 0 Name = %q, want %q", ct.Entries[0].Name, "TrigSel")
	}
	if ct.Entries[0].Type != "ParamWrite" {
		t.Errorf("Entry 0 Type = %q, want %q", ct.Entries[0].Type, "ParamWrite")
	}
	if ct.Entries[0].Par1 == nil || ct.Entries[0].Par1.Literal == nil || *ct.Entries[0].Par1.Literal != 4153 {
		t.Errorf("Entry 0 Par1 = %v, want 4153", ct.Entries[0].Par1)
	}

	// Test encoding
	for i, entry := range ct.Entries {
		encoded, err := entry.Encode(ct)
		if err != nil {
			t.Fatalf("Entry %d Encode() failed: %v", i, err)
		}
		if len(encoded) != 64 {
			t.Errorf("Entry %d encoded length = %d, want %d", i, len(encoded), 64)
		}
		// Verify A701h version header (0x01 0xA7 in little-endian)
		if encoded[0] != 0x01 || encoded[1] != 0xA7 {
			t.Errorf("Entry %d missing A701h header: got %02x %02x", i, encoded[0], encoded[1])
		}
	}
}

func TestCommandTable_UnsupportedType(t *testing.T) {
	path := writeTemp(t, "version: 1\nentries:\n  - id: 5\n    name: X\n    type: UnknownType\n")
	if _, err := Load(path); err == nil {
		t.Fatalf("expected error for unsupported type")
	}
}

func TestCommandTable_DuplicateIDs(t *testing.T) {
	path := writeTemp(t, "version: 1\nentries:\n  - id: 5\n    name: A\n    type: VAI_GoToPos\n    par1: 1\n    par2: 2\n    par3: 3\n    par4: 4\n  - id: 5\n    name: B\n    type: VAI_GoToPos\n    par1: 1\n    par2: 2\n    par3: 3\n    par4: 4\n")
	if _, err := Load(path); err == nil {
		t.Fatalf("expected error for duplicate ids")
	}
}

func TestDecodeEntry_RoundTrip(t *testing.T) {
	// Test round-trip encoding/decoding for various entry types
	tests := []struct {
		name  string
		entry Entry
	}{
		{
			name: "VAI_GoToPos",
			entry: Entry{
				ID:   1,
				Name: "TestMove",
				Type: "VAI_GoToPos",
				Par1: intPtr(1000),
				Par2: intPtr(5000),
				Par3: intPtr(10000),
				Par4: intPtr(10000),
			},
		},
		{
			name: "MoveRel",
			entry: Entry{
				ID:   2,
				Name: "TestRel",
				Type: "MoveRel",
				Par1: intPtr(-500),
				Par2: intPtr(3000),
				Par3: intPtr(8000),
			},
		},
		{
			name: "Home",
			entry: Entry{
				ID:   3,
				Name: "HomePos",
				Type: "Home",
				Par1: intPtr(0),
			},
		},
		{
			name: "NoOp",
			entry: Entry{
				ID:   4,
				Name: "NoOpEntry",
				Type: "NoOp",
			},
		},
		{
			name: "ParamWrite",
			entry: Entry{
				ID:   5,
				Name: "Param",
				Type: "ParamWrite",
				Par1: intPtr(4153),
				Par2: intPtr(1),
			},
		},
		{
			name: "SetDO",
			entry: Entry{
				ID:   6,
				Name: "SetDO",
				Type: "SetDO",
				Par1: intPtr(0x0001),
				Par2: intPtr(0x0001),
			},
		},
		{
			name: "ClearDO",
			entry: Entry{
				ID:   7,
				Name: "ClearDO",
				Type: "ClearDO",
				Par1: intPtr(0x0001),
			},
		},
		{
			name: "WithSequencedEntry",
			entry: Entry{
				ID:             8,
				Name:           "Sequenced",
				Type:           "Delay",
				Par1:           intPtr(1000),
				SequencedEntry: intPtrInt(9),
			},
		},
		{
			name: "IfDemandPosLT",
			entry: Entry{
				ID:   10,
				Name: "IfTest",
				Type: "IfDemandPosLT",
				Par1: intPtr(5000),
				Par2: intPtr(11),
				Par3: intPtr(12),
			},
		},
		{
			name: "IfMaskedX4Eq",
			entry: Entry{
				ID:   20,
				Name: "IfMasked",
				Type: "IfMaskedX4Eq",
				Par1: intPtr(0x00FF),
				Par2: intPtr(0x0055),
				Par3: intPtr(21),
				Par4: intPtr(22),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal CommandTable for encoding (no vars needed for literal params)
			ct := &CommandTable{}
			// Encode
			encoded, err := tt.entry.Encode(ct)
			if err != nil {
				t.Fatalf("Encode() failed: %v", err)
			}

			// Decode
			decoded, err := protocol_command_tables.DecodeEntry(encoded)
			if err != nil {
				t.Fatalf("DecodeEntry() failed: %v", err)
			}

			// Compare fields
			if decoded.Type != tt.entry.Type {
				t.Errorf("Type = %q, want %q", decoded.Type, tt.entry.Type)
			}
			if decoded.Name != tt.entry.Name {
				t.Errorf("Name = %q, want %q", decoded.Name, tt.entry.Name)
			}

			// Compare SequencedEntry
			if (decoded.SequencedEntry == nil) != (tt.entry.SequencedEntry == nil) {
				t.Errorf("SequencedEntry nil mismatch: got %v, want %v", decoded.SequencedEntry, tt.entry.SequencedEntry)
			} else if decoded.SequencedEntry != nil && tt.entry.SequencedEntry != nil {
				if int(*decoded.SequencedEntry) != *tt.entry.SequencedEntry {
					t.Errorf("SequencedEntry = %d, want %d", *decoded.SequencedEntry, *tt.entry.SequencedEntry)
				}
			}

			// Compare parameters (decoded has *int64, entry has *Param)
			compareParamToInt64(t, "Par1", decoded.Par1, tt.entry.Par1)
			compareParamToInt64(t, "Par2", decoded.Par2, tt.entry.Par2)
			compareParamToInt64(t, "Par3", decoded.Par3, tt.entry.Par3)
			compareParamToInt64(t, "Par4", decoded.Par4, tt.entry.Par4)
		})
	}
}

func compareParamToInt64(t *testing.T, name string, got *int64, want *Param) {
	t.Helper()
	gotNil := got == nil
	wantNil := want == nil || want.Literal == nil
	if gotNil != wantNil {
		t.Errorf("%s nil mismatch: got %v, want %v", name, gotNil, wantNil)
		return
	}
	if got != nil && want != nil && want.Literal != nil {
		if *got != *want.Literal {
			t.Errorf("%s = %d, want %d", name, *got, *want.Literal)
		}
	}
}

func intPtr(i int64) *Param {
	return &Param{Literal: &i}
}

func intPtrInt(i int) *int {
	return &i
}

// Helper to convert *int64 to *Param for test compatibility
func paramFromInt64(p *int64) *Param {
	if p == nil {
		return nil
	}
	return &Param{Literal: p}
}

func TestDecodeEntry_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr string
	}{
		{
			name:    "TooShort",
			data:    make([]byte, 32),
			wantErr: "must be at least 64 bytes",
		},
		{
			name:    "InvalidVersion",
			data:    make([]byte, 64),
			wantErr: "invalid version header",
		},
		{
			name: "UnknownHeader",
			data: func() []byte {
				b := make([]byte, 64)
				b[0] = 0x01
				b[1] = 0xA7 // Valid version
				b[4] = 0xFF
				b[5] = 0xFF // Unknown header
				return b
			}(),
			wantErr: "unknown motion header",
		},
		{
			name: "NoNulTerminator",
			data: func() []byte {
				b := make([]byte, 64)
				b[0] = 0x01
				b[1] = 0xA7 // Valid version
				b[4] = 0x00
				b[5] = 0x00 // NoOp header
				// Fill name field with non-zero bytes
				for i := 38; i < 54; i++ {
					b[i] = 0xFF
				}
				return b
			}(),
			wantErr: "name field not NUL-terminated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := protocol_command_tables.DecodeEntry(tt.data)
			if err == nil {
				t.Fatalf("DecodeEntry() expected error containing %q, got nil", tt.wantErr)
			}
			if err.Error() == "" || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("DecodeEntry() error = %q, want error containing %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestCommandTable_ParamPresence(t *testing.T) {
	// missing par4
	path := writeTemp(t, "version: 1\nentries:\n  - id: 7\n    name: P\n    type: VAI_GoToPos\n    par1: 1\n    par2: 2\n    par3: 3\n")
	if _, err := Load(path); err == nil {
		t.Fatalf("expected error for missing required params")
	}
}

func TestCommandTable_SequencedEntryRange(t *testing.T) {
	path := writeTemp(t, "version: 1\nentries:\n  - id: 8\n    name: P\n    type: VAI_GoToPos\n    sequenced_entry: 0\n    par1: 1\n    par2: 2\n    par3: 3\n    par4: 4\n")
	if _, err := Load(path); err == nil {
		t.Fatalf("expected error for invalid sequenced_entry range")
	}
}

func TestCommandTable_ProfileGating(t *testing.T) {
	// type allowed by default; still OK when model provided
	path := writeTemp(t, "version: 1\ndrive_model: C1250-MI\nentries:\n  - id: 9\n    name: A\n    type: VAI_GoToPos\n    par1: 1\n    par2: 2\n    par3: 3\n    par4: 4\n")
	if _, err := Load(path); err != nil {
		t.Fatalf("unexpected gating error: %v", err)
	}
	// Unknown type should be rejected with gating before type lookup
	path2 := writeTemp(t, "version: 1\ndrive_model: C1250-MI\nentries:\n  - id: 10\n    name: B\n    type: NotAllowed\n")
	if _, err := Load(path2); err == nil {
		t.Fatalf("expected gating error for NotAllowed type")
	}
}

func TestCommandTable_SequencedEntryReferentialIntegrity(t *testing.T) {
	// sequenced_entry points to non-existent id 99
	yaml := "version: 1\nentries:\n  - id: 1\n    name: A\n    type: MoveAbs\n    par1: 0\n    par2: 1\n    par3: 1\n    par4: 1\n    sequenced_entry: 99\n"
	if _, err := Load(writeTemp(t, yaml)); err == nil {
		t.Fatalf("expected error for sequenced_entry referencing missing id")
	}
}

func TestCommandTable_WaitWithoutTriggerParam(t *testing.T) {
	// WaitRising/WaitFalling don't require ParamWrite in the same table
	// The trigger source can be configured via parameters outside the command table
	yaml := "version: 1\nentries:\n  - id: 1\n    name: W\n    type: WaitRising\n"
	if _, err := Load(writeTemp(t, yaml)); err != nil {
		t.Fatalf("WaitRising should work without ParamWrite in table: %v", err)
	}
	// WaitFalling should also work without ParamWrite
	yaml2 := "version: 1\nentries:\n  - id: 1\n    name: W\n    type: WaitFalling\n"
	if _, err := Load(writeTemp(t, yaml2)); err != nil {
		t.Fatalf("WaitFalling should work without ParamWrite in table: %v", err)
	}
}

func TestMoveRel_Encoding(t *testing.T) {
	d, v, a := int64(-250), int64(500000), int64(800000)
	be := Entry{ID: 11, Name: "Rel", Type: "MoveRel", Par1: paramFromInt64(&d), Par2: paramFromInt64(&v), Par3: paramFromInt64(&a)}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build MoveRel: %v", err)
	}
	// header 0x0110 (little-endian -> 0x10, 0x01)
	if b[4] != 0x10 || b[5] != 0x01 {
		t.Fatalf("MoveRel header: %02x %02x", b[4], b[5])
	}
}

func TestDelay_Encoding(t *testing.T) {
	ms := int64(250)
	be := Entry{ID: 12, Name: "Dwell", Type: "Delay", Par1: paramFromInt64(&ms)}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build Delay: %v", err)
	}
	// header 0x2100 (little-endian -> 0x00, 0x21)
	if b[4] != 0x00 || b[5] != 0x21 {
		t.Fatalf("Delay header: %02x %02x", b[4], b[5])
	}
}

func TestSetDO_Encoding(t *testing.T) {
	mask, value := int64(0x0001), int64(0x0001)
	be := Entry{ID: 13, Name: "DO1On", Type: "SetDO", Par1: paramFromInt64(&mask), Par2: paramFromInt64(&value)}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build SetDO: %v", err)
	}
	if b[4] != 0x30 || b[5] != 0x00 {
		t.Fatalf("SetDO header: %02x %02x", b[4], b[5])
	}
}

func TestSetX6_ClearX6_Encoding(t *testing.T) {
	mask, value := int64(0x0003), int64(0x0001)
	be := Entry{ID: 40, Name: "X6on", Type: "SetX6", Par1: paramFromInt64(&mask), Par2: paramFromInt64(&value)}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build SetX6: %v", err)
	}
	if b[4] != 0x40 || b[5] != 0x00 {
		t.Fatalf("SetX6 header: %02x %02x", b[4], b[5])
	}
	be2 := Entry{ID: 41, Name: "X6off", Type: "ClearX6", Par1: paramFromInt64(&mask)}
	b2, err := protocol_command_tables.BuildCTEntry(be2.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build ClearX6: %v", err)
	}
	if b2[4] != 0x40 || b2[5] != 0x00 {
		t.Fatalf("ClearX6 header: %02x %02x", b2[4], b2[5])
	}
	if b2[8] != 0x00 || b2[9] != 0x00 {
		t.Fatalf("ClearX6 value not zero: %02x %02x", b2[8], b2[9])
	}
}

func TestClearDO_Encoding(t *testing.T) {
	mask := int64(0x0003)
	be := Entry{ID: 18, Name: "Clear", Type: "ClearDO", Par1: paramFromInt64(&mask)}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build ClearDO: %v", err)
	}
	if b[4] != 0x30 || b[5] != 0x00 {
		t.Fatalf("ClearDO header: %02x %02x", b[4], b[5])
	}
	// value @ 8..9 should be 0
	if b[8] != 0x00 || b[9] != 0x00 {
		t.Fatalf("ClearDO value not zero: %02x %02x", b[8], b[9])
	}
}

func TestParamWrite_Encoding(t *testing.T) {
	upid, val := int64(0x1039), int64(2)
	be := Entry{ID: 14, Name: "Trig", Type: "ParamWrite", Par1: paramFromInt64(&upid), Par2: paramFromInt64(&val)}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build ParamWrite: %v", err)
	}
	if b[4] != 0x20 || b[5] != 0x00 {
		t.Fatalf("ParamWrite header: %02x %02x", b[4], b[5])
	}
}

func TestHome_Encoding(t *testing.T) {
	pos := int64(0)
	be := Entry{ID: 15, Name: "Home", Type: "Home", Par1: paramFromInt64(&pos)}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build Home: %v", err)
	}
	// header 0x0090 (little-endian -> 0x90, 0x00)
	if b[4] != 0x90 || b[5] != 0x00 {
		t.Fatalf("Home header: %02x %02x", b[4], b[5])
	}
}

func TestNoOp_Encoding(t *testing.T) {
	be := Entry{ID: 21, Name: "End", Type: "NoOp"}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build NoOp: %v", err)
	}
	if b[4] != 0x00 || b[5] != 0x00 {
		t.Fatalf("NoOp header: %02x %02x", b[4], b[5])
	}
	for i := 6; i < 38; i++ {
		if b[i] != 0 {
			t.Fatalf("NoOp params not zero at %d", i)
		}
	}
}

func TestWaits_Encoding(t *testing.T) {
	wr := Entry{ID: 16, Name: "WR", Type: "WaitRising"}
	wf := Entry{ID: 17, Name: "WF", Type: "WaitFalling"}
	ct := &CommandTable{}
	br, err := protocol_command_tables.BuildCTEntry(wr.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build WaitRising: %v", err)
	}
	bf, err := protocol_command_tables.BuildCTEntry(wf.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build WaitFalling: %v", err)
	}
	if br[4] != 0x30 || br[5] != 0x21 {
		t.Fatalf("WaitRising header: %02x %02x", br[4], br[5])
	}
	if bf[4] != 0x40 || bf[5] != 0x21 {
		t.Fatalf("WaitFalling header: %02x %02x", bf[4], bf[5])
	}
	// params should be zeroed for waits
	for i := 6; i < 38; i++ {
		if br[i] != 0 {
			t.Fatalf("WaitRising param not zero at %d", i)
		}
	}
	for i := 6; i < 38; i++ {
		if bf[i] != 0 {
			t.Fatalf("WaitFalling param not zero at %d", i)
		}
	}
}

func TestStop_Encoding(t *testing.T) {
	dec := int64(123456)
	be := Entry{ID: 22, Name: "Stop", Type: "Stop", Par1: paramFromInt64(&dec)}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build Stop: %v", err)
	}
	if b[4] != 0x70 || b[5] != 0x01 {
		t.Fatalf("Stop header: %02x %02x", b[4], b[5])
	}
}

func TestWaitDemandVelLT_Encoding(t *testing.T) {
	thr := int64(0)
	be := Entry{ID: 23, Name: "WaitV0", Type: "WaitDemandVelLT", Par1: paramFromInt64(&thr)}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build WaitDemandVelLT: %v", err)
	}
	if b[4] != 0x90 || b[5] != 0x22 {
		t.Fatalf("WaitDemandVelLT header: %02x %02x", b[4], b[5])
	}
}

func TestInfiniteMotion_NoParams(t *testing.T) {
	types := []string{"InfiniteMotionPos", "InfiniteMotionNeg", "InfiniteMotionPos_DecEqAcc", "InfiniteMotionNeg_DecEqAcc"}
	want := [][2]byte{{0xE0, 0x02}, {0xF0, 0x02}, {0xE0, 0x0C}, {0xF0, 0x0C}}
	ct := &CommandTable{}
	for i, typ := range types {
		be := Entry{ID: 30 + i, Name: "Jog", Type: typ}
		b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
		if err != nil {
			t.Fatalf("build %s: %v", typ, err)
		}
		if b[4] != want[i][0] || b[5] != want[i][1] {
			t.Fatalf("%s header: %02x %02x", typ, b[4], b[5])
		}
		for j := 6; j < 38; j++ {
			if b[j] != 0 {
				t.Fatalf("%s params not zero at %d", typ, j)
			}
		}
	}
}

func TestAliases_New(t *testing.T) {
	// VAI_Stop maps to Stop (0x0170)
	p := int64(100)
	be := Entry{ID: 50, Name: "S", Type: "VAI_Stop", Par1: paramFromInt64(&p)}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build VAI_Stop: %v", err)
	}
	if b[4] != 0x70 || b[5] != 0x01 {
		t.Fatalf("VAI_Stop header: %02x %02x", b[4], b[5])
	}
	// VAI_InfiniteMotionPlus maps to 0x02E0
	be2 := Entry{ID: 51, Name: "Jog+", Type: "VAI_InfiniteMotionPlus"}
	b2, err := protocol_command_tables.BuildCTEntry(be2.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build VAI_InfiniteMotionPlus: %v", err)
	}
	if b2[4] != 0xE0 || b2[5] != 0x02 {
		t.Fatalf("VAI_InfiniteMotionPlus header: %02x %02x", b2[4], b2[5])
	}
}

func TestSequencedEntry_LinkID(t *testing.T) {
	next := 42
	pos := int64(100)
	be := Entry{ID: 10, Name: "A", Type: "MoveAbs", Par1: paramFromInt64(&pos), Par2: paramFromInt64(&pos), Par3: paramFromInt64(&pos), Par4: paramFromInt64(&pos), SequencedEntry: &next}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build with link: %v", err)
	}
	// Link ID at CT bytes 2..3, little-endian
	if b[2] != byte(next&0xFF) || b[3] != byte((next>>8)&0xFF) {
		t.Fatalf("link id: got %02x %02x", b[2], b[3])
	}
}

func TestIfThreshold_Encoding(t *testing.T) {
	thr, tID, fID := int64(123456), int64(5), int64(6)
	be := Entry{ID: 19, Name: "IfDP<", Type: "IfDemandPosLT", Par1: paramFromInt64(&thr), Par2: paramFromInt64(&tID), Par3: paramFromInt64(&fID)}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build IfDemandPosLT: %v", err)
	}
	// header 0x2580 (little-endian -> 0x80, 0x25)
	if b[4] != 0x80 || b[5] != 0x25 {
		t.Fatalf("IfDemandPosLT header: %02x %02x", b[4], b[5])
	}
	// true/false ids
	if b[12] != byte(tID&0xFF) || b[13] != byte((tID>>8)&0xFF) {
		t.Fatalf("true id mismatch")
	}
	if b[14] != byte(fID&0xFF) || b[15] != byte((fID>>8)&0xFF) {
		t.Fatalf("false id mismatch")
	}
}

func TestIfMasked_Encoding(t *testing.T) {
	mask, val, tID, fID := int64(3), int64(1), int64(7), int64(8)
	be := Entry{ID: 20, Name: "IfX4==", Type: "IfMaskedX4Eq", Par1: paramFromInt64(&mask), Par2: paramFromInt64(&val), Par3: paramFromInt64(&tID), Par4: paramFromInt64(&fID)}
	ct := &CommandTable{}
	b, err := protocol_command_tables.BuildCTEntry(be.toWireEntry(ct))
	if err != nil {
		t.Fatalf("build IfMaskedX4Eq: %v", err)
	}
	if b[4] != 0x20 || b[5] != 0x26 {
		t.Fatalf("IfMaskedX4Eq header: %02x %02x", b[4], b[5])
	}
	if b[10] != byte(tID&0xFF) || b[11] != byte((tID>>8)&0xFF) {
		t.Fatalf("true id mismatch")
	}
	if b[12] != byte(fID&0xFF) || b[13] != byte((fID>>8)&0xFF) {
		t.Fatalf("false id mismatch")
	}
}

func TestIfBranchRefIntegrity(t *testing.T) {
	// Build manifest with IF referencing missing IDs
	yaml := "version: 1\nentries:\n  - id: 1\n    name: Cond\n    type: IfDemandPosLT\n    par1: 0\n    par2: 250\n    par3: 251\n"
	if _, err := Load(writeTemp(t, yaml)); err == nil {
		t.Fatalf("expected referential integrity error for IF targets")
	}
}
