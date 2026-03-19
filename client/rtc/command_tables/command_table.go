package client_command_tables

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	protocol_command_tables "gsail-go/linmot/protocol/rtc/command_tables"
)

// Param represents a template-capable parameter value.
// It can hold either a literal int64 or a variable reference like ${VAR}.
type Param struct {
	Literal *int64 // Set when value is literal int (e.g., par2: 100000)
	VarName string // Set when value is variable reference (e.g., par2: ${MAX_VELOCITY})
}

// UnmarshalYAML implements custom YAML unmarshaling for Param.
// Supports literal integers and variable references like ${VAR}.
func (p *Param) UnmarshalYAML(node *yaml.Node) error {
	// Handle null/omitted values
	if node.Kind == yaml.ScalarNode && node.Tag == "!!null" {
		return nil // Param pointer will remain nil
	}

	// Try to decode as int64 literal first (most common case)
	var v int64
	if err := node.Decode(&v); err == nil {
		p.Literal = &v
		p.VarName = ""
		return nil
	}

	// Fallback: try to decode as string for variable reference
	var s string
	if err := node.Decode(&s); err != nil {
		return fmt.Errorf("Param: expected int64 or ${VAR} string, got %q", node.Value)
	}

	// Parse ${VAR} pattern
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "${") && strings.HasSuffix(s, "}") {
		varName := strings.TrimSpace(s[2 : len(s)-1])
		if varName == "" {
			return fmt.Errorf("Param: empty variable name in %q", s)
		}
		p.VarName = varName
		p.Literal = nil
		return nil
	}

	return fmt.Errorf("Param: invalid variable reference %q (expected ${VAR})", s)
}

// CommandTable represents a command table that can be loaded from YAML
// and validated. This is the primary domain type for command table operations.
type CommandTable struct {
	Version    string           `yaml:"version"`
	DriveModel string           `yaml:"drive_model"`
	Entries    []Entry          `yaml:"entries"`
	vars       map[string]int64 // internal, not YAML-mapped; variable bindings
}

// Entry represents a single command table entry with its parameters.
// This type can be encoded to a 64-byte binary representation for deployment to drives.
type Entry struct {
	ID             int    `yaml:"id"`
	Name           string `yaml:"name"`
	Type           string `yaml:"type"`
	SequencedEntry *int   `yaml:"sequenced_entry"`
	Par1           *Param `yaml:"par1"`
	Par2           *Param `yaml:"par2"`
	Par3           *Param `yaml:"par3"`
	Par4           *Param `yaml:"par4"`
	Par5           *Param `yaml:"par5"`
	Par6           *Param `yaml:"par6"`
	Par7           *Param `yaml:"par7"`
	Par8           *Param `yaml:"par8"`
}

// Encode converts this entry to a 64-byte binary representation.
// Requires CommandTable context for variable resolution.
func (e *Entry) Encode(ct *CommandTable) ([]byte, error) {
	wireEntry := e.toWireEntry(ct)
	return protocol_command_tables.BuildCTEntry(wireEntry)
}

// toWireEntry converts an Entry to protocol WireEntry, resolving template variables.
func (e *Entry) toWireEntry(ct *CommandTable) protocol_command_tables.WireEntry {
	var sequencedEntry *uint8
	if e.SequencedEntry != nil {
		id := uint8(*e.SequencedEntry)
		sequencedEntry = &id
	}
	return protocol_command_tables.WireEntry{
		ID:             uint8(e.ID),
		Name:           e.Name,
		Type:           e.Type,
		SequencedEntry: sequencedEntry,
		Par1:           ct.resolveParam(e.Par1),
		Par2:           ct.resolveParam(e.Par2),
		Par3:           ct.resolveParam(e.Par3),
		Par4:           ct.resolveParam(e.Par4),
		Par5:           ct.resolveParam(e.Par5),
		Par6:           ct.resolveParam(e.Par6),
		Par7:           ct.resolveParam(e.Par7),
		Par8:           ct.resolveParam(e.Par8),
	}
}

// SetVar binds a variable name to a value (CT units: int64).
// Variable names should match references in YAML like ${VAR_NAME}.
func (ct *CommandTable) SetVar(name string, value int64) {
	if ct.vars == nil {
		ct.vars = make(map[string]int64)
	}
	ct.vars[name] = value
}

// RequiredVars returns the list of variable names referenced in this table.
// Scans all entries for ${VAR} references in any parameter field.
func (ct *CommandTable) RequiredVars() []string {
	seen := make(map[string]struct{})
	for _, e := range ct.Entries {
		for _, p := range []*Param{e.Par1, e.Par2, e.Par3, e.Par4, e.Par5, e.Par6, e.Par7, e.Par8} {
			if p != nil && p.VarName != "" {
				seen[p.VarName] = struct{}{}
			}
		}
	}
	out := make([]string, 0, len(seen))
	for name := range seen {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// resolveParam converts a Param to *int64 using variable bindings.
// Returns nil if Param is nil.
// Returns Literal value if Param has a literal.
// Returns bound variable value if Param has VarName.
// Panics if variable is unbound (caller must ensure all vars are bound).
func (ct *CommandTable) resolveParam(p *Param) *int64 {
	if p == nil {
		return nil
	}
	if p.Literal != nil {
		return p.Literal
	}
	if p.VarName != "" {
		if ct.vars == nil {
			panic(fmt.Sprintf("unbound variable: %s", p.VarName))
		}
		val, ok := ct.vars[p.VarName]
		if !ok {
			panic(fmt.Sprintf("unbound variable: %s", p.VarName))
		}
		return &val
	}
	return nil
}

// Validate validates the command table, checking entry IDs, types, parameters,
// and referential integrity.
func (ct *CommandTable) Validate() error {
	return validateCommandTable(ct)
}

// Load reads and validates a command table YAML file from disk.
func Load(path string) (*CommandTable, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

// Parse parses a command table from YAML bytes.
func Parse(data []byte) (*CommandTable, error) {
	var ct CommandTable
	if err := yaml.Unmarshal(data, &ct); err != nil {
		return nil, err
	}
	// Run structural validation (template-safe)
	if err := validateStructure(&ct); err != nil {
		return nil, err
	}
	// If no template variables, run full validation immediately (backward compatible)
	if len(ct.RequiredVars()) == 0 {
		if err := ct.Validate(); err != nil {
			return nil, err
		}
	}
	return &ct, nil
}

// validateStructure performs template-safe structural validation.
// Can run even when variables are unbound.
func validateStructure(ct *CommandTable) error {
	ids := make(map[int]struct{}, len(ct.Entries))
	for i, e := range ct.Entries {
		if e.ID <= 0 || e.ID > 255 {
			return fmt.Errorf("entry %d: id out of range: %d", i, e.ID)
		}
		if _, ok := ids[e.ID]; ok {
			return fmt.Errorf("duplicate id: %d", e.ID)
		}
		ids[e.ID] = struct{}{}
		if e.Name == "" {
			return fmt.Errorf("entry %d: missing name", i)
		}
		if len(e.Name) > 15 { // enforce ≤15 so we can write trailing NUL in 16-byte field
			return fmt.Errorf("entry %d: name too long (>15)", i)
		}
		if e.Type == "" {
			return fmt.Errorf("entry %d: missing type", i)
		}
		// Sequenced entry optional; if set, must be 1..255 (FFFFh represented by nil)
		if e.SequencedEntry != nil {
			if *e.SequencedEntry <= 0 || *e.SequencedEntry > 255 {
				return fmt.Errorf("entry %d: sequenced_entry out of range: %d", i, *e.SequencedEntry)
			}
		}
	}

	// Referential integrity for sequenced_entry (non-templated field)
	for i, e := range ct.Entries {
		if e.SequencedEntry != nil {
			if _, ok := ids[*e.SequencedEntry]; !ok {
				return fmt.Errorf("entry %d: sequenced_entry references missing id: %d", i, *e.SequencedEntry)
			}
		}
	}
	return nil
}

func validateCommandTable(ct *CommandTable) error {
	// Check for unbound template variables
	required := ct.RequiredVars()
	missing := make([]string, 0)
	for _, name := range required {
		if ct.vars == nil {
			missing = append(missing, name)
			continue
		}
		if _, ok := ct.vars[name]; !ok {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("unbound template variables: %v", missing)
	}

	// Run structural validation
	if err := validateStructure(ct); err != nil {
		return err
	}

	// Build ID map for referential integrity checks
	ids := make(map[int]struct{}, len(ct.Entries))
	for _, e := range ct.Entries {
		ids[e.ID] = struct{}{}
	}

	// Validate each entry with resolved parameters
	for i, e := range ct.Entries {
		// Resolve all params to concrete values for validation
		wireEntry := e.toWireEntry(ct)
		if err := protocol_command_tables.ValidateWireEntry(wireEntry); err != nil {
			return fmt.Errorf("entry %d (type=%s): %v", i, e.Type, err)
		}

		// IF-branch referential integrity: check resolved parameter values
		switch e.Type {
		case "IfDemandPosLT", "IfDemandPosGT", "IfActualPosLT", "IfActualPosGT", "IfDiffPosLT", "IfDiffPosGT", "IfCurrentLT", "IfCurrentGT", "IfAnalogX44LT":
			if wireEntry.Par2 != nil {
				if _, ok := ids[int(*wireEntry.Par2)]; !ok {
					return fmt.Errorf("entry %d: IF trueID references missing id: %d", i, *wireEntry.Par2)
				}
			}
			if wireEntry.Par3 != nil {
				if _, ok := ids[int(*wireEntry.Par3)]; !ok {
					return fmt.Errorf("entry %d: IF falseID references missing id: %d", i, *wireEntry.Par3)
				}
			}
		case "IfMaskedX4Eq", "IfMaskedX6Eq", "IfMaskedStatusEq", "IfMaskedWarnEq":
			if wireEntry.Par3 != nil {
				if _, ok := ids[int(*wireEntry.Par3)]; !ok {
					return fmt.Errorf("entry %d: IF trueID references missing id: %d", i, *wireEntry.Par3)
				}
			}
			if wireEntry.Par4 != nil {
				if _, ok := ids[int(*wireEntry.Par4)]; !ok {
					return fmt.Errorf("entry %d: IF falseID references missing id: %d", i, *wireEntry.Par4)
				}
			}
		}
	}

	return nil
}
