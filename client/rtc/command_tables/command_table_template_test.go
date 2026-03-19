package client_command_tables

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// Helper functions for tests
func paramLiteral(i int64) *Param {
	return &Param{Literal: &i}
}

func paramVar(name string) *Param {
	return &Param{VarName: name}
}

// TestParam_UnmarshalYAML_Literal tests unmarshaling literal integers
func TestParam_UnmarshalYAML_Literal(t *testing.T) {
	tests := []struct {
		name  string
		yaml  string
		want  *Param
		error bool
	}{
		{
			name: "simple integer",
			yaml: "123",
			want: paramLiteral(123),
		},
		{
			name: "zero",
			yaml: "0",
			want: paramLiteral(0),
		},
		{
			name: "negative",
			yaml: "-456",
			want: paramLiteral(-456),
		},
		{
			name: "large number",
			yaml: "100000",
			want: paramLiteral(100000),
		},
		{
			name: "hex literal",
			yaml: "0x0003",
			want: paramLiteral(3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Param
			node := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: tt.yaml,
				Tag:   "!!int",
			}
			err := p.UnmarshalYAML(node)
			if tt.error {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p.Literal == nil {
				t.Fatalf("expected literal value, got nil")
			}
			if *p.Literal != *tt.want.Literal {
				t.Errorf("Literal = %d, want %d", *p.Literal, *tt.want.Literal)
			}
			if p.VarName != "" {
				t.Errorf("VarName = %q, want empty", p.VarName)
			}
		})
	}
}

// TestParam_UnmarshalYAML_Variable tests unmarshaling variable references
func TestParam_UnmarshalYAML_Variable(t *testing.T) {
	tests := []struct {
		name  string
		yaml  string
		want  string // expected variable name
		error bool
	}{
		{
			name: "simple variable unquoted",
			yaml: "${MAX_VELOCITY}",
			want: "MAX_VELOCITY",
		},
		{
			name: "variable with underscores",
			yaml: "${DELAY_AT_BOTTOM}",
			want: "DELAY_AT_BOTTOM",
		},
		{
			name: "variable with whitespace trimmed",
			yaml: "${ MAX_VELOCITY }",
			want: "MAX_VELOCITY",
		},
		{
			name:  "empty variable name",
			yaml:  "${}",
			error: true,
		},
		{
			name:  "whitespace only variable name",
			yaml:  "${  }",
			error: true,
		},
		{
			name:  "invalid pattern - missing $",
			yaml:  "{MAX_VELOCITY}",
			error: true,
		},
		{
			name:  "invalid pattern - missing braces",
			yaml:  "$MAX_VELOCITY",
			error: true,
		},
		{
			name: "quoted variable",
			yaml: "\"${MAX_VELOCITY}\"",
			want: "MAX_VELOCITY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Param
			node := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: strings.Trim(tt.yaml, "\""),
				Tag:   "!!str",
			}
			err := p.UnmarshalYAML(node)
			if tt.error {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p.VarName != tt.want {
				t.Errorf("VarName = %q, want %q", p.VarName, tt.want)
			}
			if p.Literal != nil {
				t.Errorf("Literal = %v, want nil", p.Literal)
			}
		})
	}
}

// TestParam_UnmarshalYAML_Null tests null/omitted values
func TestParam_UnmarshalYAML_Null(t *testing.T) {
	var p *Param
	node := &yaml.Node{
		Kind: yaml.ScalarNode,
		Tag:  "!!null",
	}
	err := p.UnmarshalYAML(node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != nil {
		t.Errorf("expected nil Param for null value, got %v", p)
	}
}

// TestCommandTable_SetVar tests variable binding
func TestCommandTable_SetVar(t *testing.T) {
	ct := &CommandTable{}

	// Set a variable
	ct.SetVar("MAX_VELOCITY", 100000)

	if ct.vars == nil {
		t.Fatalf("vars map not initialized")
	}
	if val, ok := ct.vars["MAX_VELOCITY"]; !ok {
		t.Fatalf("variable not set")
	} else if val != 100000 {
		t.Errorf("variable value = %d, want 100000", val)
	}

	// Overwrite existing variable
	ct.SetVar("MAX_VELOCITY", 150000)
	if val := ct.vars["MAX_VELOCITY"]; val != 150000 {
		t.Errorf("overwritten variable value = %d, want 150000", val)
	}

	// Set another variable
	ct.SetVar("ACCELERATION", 50000)
	if len(ct.vars) != 2 {
		t.Errorf("vars map length = %d, want 2", len(ct.vars))
	}
}

// TestCommandTable_RequiredVars tests variable discovery
func TestCommandTable_RequiredVars(t *testing.T) {
	ct := &CommandTable{
		Entries: []Entry{
			{
				ID:   1,
				Name: "move down",
				Type: "VAI_GoToPos",
				Par1: paramLiteral(0),
				Par2: paramVar("MAX_VELOCITY"),
				Par3: paramVar("ACCELERATION"),
				Par4: paramVar("DECELERATION"),
			},
			{
				ID:   2,
				Name: "move up",
				Type: "VAI_GoToPos",
				Par1: paramLiteral(100000),
				Par2: paramVar("MAX_VELOCITY"), // Reused variable
				Par3: paramVar("ACCELERATION"),
				Par4: paramVar("DECELERATION"),
			},
			{
				ID:   3,
				Name: "wait",
				Type: "Delay",
				Par1: paramVar("DELAY_AT_BOTTOM"),
			},
		},
	}

	required := ct.RequiredVars()
	expectedCount := 4 // MAX_VELOCITY, ACCELERATION, DECELERATION, DELAY_AT_BOTTOM
	if len(required) != expectedCount {
		t.Errorf("RequiredVars() length = %d, want %d. Got: %v", len(required), expectedCount, required)
	}

	expected := map[string]bool{
		"MAX_VELOCITY":    true,
		"ACCELERATION":    true,
		"DECELERATION":    true,
		"DELAY_AT_BOTTOM": true,
	}
	for _, name := range required {
		if !expected[name] {
			t.Errorf("unexpected variable: %q", name)
		}
		delete(expected, name)
	}
	if len(expected) > 0 {
		t.Errorf("missing variables: %v", expected)
	}
}

// TestCommandTable_LoadTemplate tests loading templated YAML
func TestCommandTable_LoadTemplate(t *testing.T) {
	yaml := `version: 1
drive_model: C1250-MI
entries:
  - id: 1
    name: "move"
    type: VAI_GoToPos
    par1: 500000
    par2: ${MAX_VELOCITY}
    par3: ${ACCELERATION}
    par4: ${DECELERATION}
`
	path := writeTemp(t, yaml)
	defer os.Remove(path)

	ct, err := Load(path)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if len(ct.Entries) != 1 {
		t.Fatalf("Entries length = %d, want 1", len(ct.Entries))
	}

	entry := ct.Entries[0]
	if entry.Par1 == nil || entry.Par1.Literal == nil || *entry.Par1.Literal != 500000 {
		t.Errorf("Par1 (literal) = %v, want 500000", entry.Par1)
	}
	if entry.Par2 == nil || entry.Par2.VarName != "MAX_VELOCITY" {
		t.Errorf("Par2 (variable) = %v, want MAX_VELOCITY", entry.Par2)
	}
	if entry.Par3 == nil || entry.Par3.VarName != "ACCELERATION" {
		t.Errorf("Par3 (variable) = %v, want ACCELERATION", entry.Par3)
	}
	if entry.Par4 == nil || entry.Par4.VarName != "DECELERATION" {
		t.Errorf("Par4 (variable) = %v, want DECELERATION", entry.Par4)
	}

	// Should have variables but not validated yet
	required := ct.RequiredVars()
	if len(required) != 3 {
		t.Errorf("RequiredVars() = %v, want 3 variables", required)
	}
}

// TestCommandTable_Validate_UnboundVars tests validation with unbound variables
func TestCommandTable_Validate_UnboundVars(t *testing.T) {
	ct := &CommandTable{
		Entries: []Entry{
			{
				ID:   1,
				Name: "move",
				Type: "VAI_GoToPos",
				Par1: paramLiteral(0),
				Par2: paramVar("MAX_VELOCITY"),
				Par3: paramVar("ACCELERATION"),
				Par4: paramVar("DECELERATION"),
			},
		},
	}

	// Validate without binding variables - should fail
	err := ct.Validate()
	if err == nil {
		t.Fatalf("expected error for unbound variables, got nil")
	}
	if !strings.Contains(err.Error(), "unbound template variables") {
		t.Errorf("error message doesn't mention unbound variables: %v", err)
	}
	if !strings.Contains(err.Error(), "MAX_VELOCITY") {
		t.Errorf("error message doesn't list MAX_VELOCITY: %v", err)
	}
}

// TestCommandTable_Validate_BoundVars tests validation with all variables bound
func TestCommandTable_Validate_BoundVars(t *testing.T) {
	ct := &CommandTable{
		Entries: []Entry{
			{
				ID:   1,
				Name: "move",
				Type: "VAI_GoToPos",
				Par1: paramLiteral(0),
				Par2: paramVar("MAX_VELOCITY"),
				Par3: paramVar("ACCELERATION"),
				Par4: paramVar("DECELERATION"),
			},
		},
	}

	// Bind all variables
	ct.SetVar("MAX_VELOCITY", 100000)
	ct.SetVar("ACCELERATION", 50000)
	ct.SetVar("DECELERATION", 60000)

	// Should validate successfully
	err := ct.Validate()
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

// TestCommandTable_Instantiate_Encode tests encoding with variable resolution
func TestCommandTable_Instantiate_Encode(t *testing.T) {
	ct := &CommandTable{
		Entries: []Entry{
			{
				ID:   1,
				Name: "move",
				Type: "VAI_GoToPos",
				Par1: paramLiteral(0),
				Par2: paramVar("MAX_VELOCITY"),
				Par3: paramVar("ACCELERATION"),
				Par4: paramVar("DECELERATION"),
			},
		},
	}

	// Bind variables
	ct.SetVar("MAX_VELOCITY", 100000)
	ct.SetVar("ACCELERATION", 50000)
	ct.SetVar("DECELERATION", 60000)

	// Validate
	if err := ct.Validate(); err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	// Encode entry - should use resolved values
	entry := ct.Entries[0]
	encoded, err := entry.Encode(ct)
	if err != nil {
		t.Fatalf("Encode() failed: %v", err)
	}

	if len(encoded) != 64 {
		t.Errorf("encoded length = %d, want 64", len(encoded))
	}

	// Verify A701h header
	if encoded[0] != 0x01 || encoded[1] != 0xA7 {
		t.Errorf("missing A701h header: got %02x %02x", encoded[0], encoded[1])
	}
}

// TestCommandTable_Load_BackwardCompatible tests backward compatibility
func TestCommandTable_Load_BackwardCompatible(t *testing.T) {
	// Load existing non-templated YAML
	ct, _ := loadExampleManifest(t, "test_command_table.yaml")

	// Should have no required variables
	required := ct.RequiredVars()
	if len(required) != 0 {
		t.Errorf("RequiredVars() = %v, want empty (no template variables)", required)
	}

	// Should validate successfully (already validated at load time)
	if err := ct.Validate(); err != nil {
		t.Fatalf("validation failed: %v", err)
	}
}

// TestCommandTable_ResolveParam tests parameter resolution
func TestCommandTable_ResolveParam(t *testing.T) {
	ct := &CommandTable{}
	ct.SetVar("MAX_VELOCITY", 100000)

	tests := []struct {
		name   string
		param  *Param
		want   *int64
		panics bool
	}{
		{
			name:  "nil param",
			param: nil,
			want:  nil,
		},
		{
			name:  "literal param",
			param: paramLiteral(50000),
			want:  int64PtrForTest(50000),
		},
		{
			name:  "bound variable",
			param: paramVar("MAX_VELOCITY"),
			want:  int64PtrForTest(100000),
		},
		{
			name:   "unbound variable",
			param:  paramVar("UNKNOWN_VAR"),
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic, got none")
					}
				}()
			}

			got := ct.resolveParam(tt.param)
			if tt.want == nil {
				if got != nil {
					t.Errorf("resolveParam() = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Errorf("resolveParam() = nil, want %d", *tt.want)
				} else if *got != *tt.want {
					t.Errorf("resolveParam() = %d, want %d", *got, *tt.want)
				}
			}
		})
	}
}

// TestCommandTable_LoadTemplate_ConditionalValidation tests conditional validation
func TestCommandTable_LoadTemplate_ConditionalValidation(t *testing.T) {
	// Test with templated YAML - should not fully validate at load time
	yamlTemplate := `version: 1
drive_model: C1250-MI
entries:
  - id: 1
    name: "move"
    type: VAI_GoToPos
    par1: 0
    par2: ${MAX_VELOCITY}
    par3: 100000
    par4: 100000
`
	path := writeTemp(t, yamlTemplate)
	defer os.Remove(path)

	ct, err := Load(path)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Should have variables
	required := ct.RequiredVars()
	if len(required) == 0 {
		t.Fatalf("expected template variables, got none")
	}

	// Should fail validation until vars are bound
	err = ct.Validate()
	if err == nil {
		t.Fatalf("expected validation error for unbound vars, got nil")
	}

	// Bind variable and validate again
	ct.SetVar("MAX_VELOCITY", 100000)
	err = ct.Validate()
	if err != nil {
		t.Fatalf("validation failed after binding vars: %v", err)
	}
}

// Integration test with example template file
func TestCommandTable_TemplateIntegration(t *testing.T) {
	// Load template YAML
	_, filename, _, _ := runtime.Caller(0)
	testdataDir := filepath.Join(filepath.Dir(filename), "testdata")
	templatePath := filepath.Join(testdataDir, "test_command_table_template.yaml")

	ct, err := Load(templatePath)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Check required variables
	required := ct.RequiredVars()
	if len(required) == 0 {
		t.Fatalf("expected template variables, got none")
	}

	// Bind all variables
	ct.SetVar("MAX_VELOCITY", 100000)
	ct.SetVar("ACCELERATION", 100000)
	ct.SetVar("DECELERATION", 100000)
	ct.SetVar("DELAY_AT_BOTTOM", 10000)
	ct.SetVar("POSITION_UP", 100000)

	// Validate
	if err := ct.Validate(); err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	// Encode all entries
	for i, entry := range ct.Entries {
		encoded, err := entry.Encode(ct)
		if err != nil {
			t.Fatalf("Entry %d Encode() failed: %v", i, err)
		}
		if len(encoded) != 64 {
			t.Errorf("Entry %d encoded length = %d, want 64", i, len(encoded))
		}
	}
}

func int64PtrForTest(i int64) *int64 {
	return &i
}
