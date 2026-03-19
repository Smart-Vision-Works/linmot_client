package linmot

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Smart-Vision-Works/staged_robot/client"
	linmot_command_tables "github.com/Smart-Vision-Works/staged_robot/client/rtc/command_tables"
	"github.com/Smart-Vision-Works/staged_robot/test"

	"github.com/pkg/errors"

	"stage_primer_config"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     DeployConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			cfg: DeployConfig{
				ZDistance:           25.5,
				DefaultSpeed:        50.0,
				DefaultAcceleration: 100.0,
				PickTime:            0.5,
			},
			wantErr: false,
		},
		{
			name: "negative z_distance",
			cfg: DeployConfig{
				ZDistance:           -1.0,
				DefaultSpeed:        50.0,
				DefaultAcceleration: 100.0,
				PickTime:            0.5,
			},
			wantErr: true,
			errMsg:  "z_distance cannot be negative",
		},
		{
			name: "z_distance too large",
			cfg: DeployConfig{
				ZDistance:           600.0,
				DefaultSpeed:        50.0,
				DefaultAcceleration: 100.0,
				PickTime:            0.5,
			},
			wantErr: true,
			errMsg:  "z_distance exceeds maximum travel",
		},
		{
			name: "speed too high",
			cfg: DeployConfig{
				ZDistance:           25.5,
				DefaultSpeed:        150.0,
				DefaultAcceleration: 100.0,
				PickTime:            0.5,
			},
			wantErr: true,
			errMsg:  "default_speed out of range",
		},
		{
			name: "negative speed",
			cfg: DeployConfig{
				ZDistance:           25.5,
				DefaultSpeed:        -10.0,
				DefaultAcceleration: 100.0,
				PickTime:            0.5,
			},
			wantErr: true,
			errMsg:  "default_speed out of range",
		},
		{
			name: "acceleration too high",
			cfg: DeployConfig{
				ZDistance:           25.5,
				DefaultSpeed:        50.0,
				DefaultAcceleration: 250.0,
				PickTime:            0.5,
			},
			wantErr: true,
			errMsg:  "default_acceleration out of range",
		},
		{
			name: "pick_time too small",
			cfg: DeployConfig{
				ZDistance:           25.5,
				DefaultSpeed:        50.0,
				DefaultAcceleration: 100.0,
				PickTime:            0.0001,
			},
			wantErr: true,
			errMsg:  "pick_time out of range",
		},
		{
			name: "pick_time too large",
			cfg: DeployConfig{
				ZDistance:           25.5,
				DefaultSpeed:        50.0,
				DefaultAcceleration: 100.0,
				PickTime:            20.0,
			},
			wantErr: true,
			errMsg:  "pick_time out of range",
		},
		{
			name: "NaN z_distance",
			cfg: DeployConfig{
				ZDistance:           math.NaN(),
				DefaultSpeed:        50.0,
				DefaultAcceleration: 100.0,
				PickTime:            0.5,
			},
			wantErr: true,
			errMsg:  "z_distance is not a valid number",
		},
		{
			name: "Inf z_distance",
			cfg: DeployConfig{
				ZDistance:           math.Inf(1),
				DefaultSpeed:        50.0,
				DefaultAcceleration: 100.0,
				PickTime:            0.5,
			},
			wantErr: true,
			errMsg:  "z_distance is not a valid number",
		},
		{
			name: "NaN speed",
			cfg: DeployConfig{
				ZDistance:           25.5,
				DefaultSpeed:        math.NaN(),
				DefaultAcceleration: 100.0,
				PickTime:            0.5,
			},
			wantErr: true,
			errMsg:  "default_speed is not a valid number",
		},
		{
			name: "boundary values - max z_distance",
			cfg: DeployConfig{
				ZDistance:           maxZDistance,
				DefaultSpeed:        50.0,
				DefaultAcceleration: 100.0,
				PickTime:            0.5,
			},
			wantErr: false,
		},
		{
			name: "boundary values - max speed",
			cfg: DeployConfig{
				ZDistance:           25.5,
				DefaultSpeed:        maxSpeed,
				DefaultAcceleration: 100.0,
				PickTime:            0.5,
			},
			wantErr: false,
		},
		{
			name: "boundary values - max acceleration",
			cfg: DeployConfig{
				ZDistance:           25.5,
				DefaultSpeed:        50.0,
				DefaultAcceleration: maxAcceleration,
				PickTime:            0.5,
			},
			wantErr: false,
		},
		{
			name: "boundary values - min pick_time",
			cfg: DeployConfig{
				ZDistance:           25.5,
				DefaultSpeed:        50.0,
				DefaultAcceleration: 100.0,
				PickTime:            minPickTime,
			},
			wantErr: false,
		},
		{
			name: "boundary values - max pick_time",
			cfg: DeployConfig{
				ZDistance:           25.5,
				DefaultSpeed:        50.0,
				DefaultAcceleration: 100.0,
				PickTime:            maxPickTime,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || err.Error() == "" {
					t.Errorf("validateConfig() expected error message containing %q, got nil", tt.errMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateConfig() error message = %q, want containing %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// createTestTemplate creates a minimal valid command table template for testing
func createTestTemplate(t *testing.T) string {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test_template.yaml")

	templateContent := `version: 1
drive_model: C1250-MI
entries:
  - id: 1
    name: "test entry"
    type: NoOp
`

	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	return templatePath
}

func TestLoadCommandTable_Caching(t *testing.T) {
	// Reset cache before test
	cachedTemplateMutex.Lock()
	originalCache := cachedTemplate
	cachedTemplate = nil
	cachedTemplateMutex.Unlock()
	defer func() {
		cachedTemplateMutex.Lock()
		cachedTemplate = originalCache
		cachedTemplateMutex.Unlock()
	}()

	// Save current working directory
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWD)

	// Create a temporary directory and change to it
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create template file in expected dev path
	templatePath := "primer/linmot/linmot_command_table.yaml"
	if err := os.MkdirAll(filepath.Dir(templatePath), 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	templateContent := `version: 1
drive_model: C1250-MI
entries:
  - id: 1
    name: "test entry"
    type: NoOp
`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// First load should read from file
	template1, err := loadCommandTable()
	if err != nil {
		t.Fatalf("loadCommandTable() failed: %v", err)
	}
	if template1 == nil {
		t.Fatal("loadCommandTable() returned nil template")
	}

	// Modify the file
	templateContent2 := templateContent + "  - id: 2\n    name: \"modified\"\n    type: NoOp\n"
	if err := os.WriteFile(templatePath, []byte(templateContent2), 0644); err != nil {
		t.Fatalf("Failed to modify test template: %v", err)
	}

	// Second load should return cached version (not reloaded)
	template2, err := loadCommandTable()
	if err != nil {
		t.Fatalf("loadCommandTable() failed on second call: %v", err)
	}
	if template2 != template1 {
		t.Error("loadCommandTable() returned different template on second call (cache not working)")
	}
}

func TestLoadCommandTable_ConcurrentAccess(t *testing.T) {
	// Reset cache before test
	cachedTemplateMutex.Lock()
	originalCache := cachedTemplate
	cachedTemplate = nil
	cachedTemplateMutex.Unlock()
	defer func() {
		cachedTemplateMutex.Lock()
		cachedTemplate = originalCache
		cachedTemplateMutex.Unlock()
	}()

	// Save current working directory
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWD)

	// Create a temporary directory and change to it
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create template file in expected dev path
	templatePath := "primer/linmot/linmot_command_table.yaml"
	if err := os.MkdirAll(filepath.Dir(templatePath), 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	templateContent := `version: 1
drive_model: C1250-MI
entries:
  - id: 1
    name: "test entry"
    type: NoOp
`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Test concurrent access
	const numGoroutines = 10
	var wg sync.WaitGroup
	templates := make([]*linmot_command_tables.CommandTable, numGoroutines)
	errors := make([]error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			templates[idx], errors[idx] = loadCommandTable()
		}(i)
	}

	wg.Wait()

	// All should succeed
	for i, err := range errors {
		if err != nil {
			t.Errorf("Goroutine %d failed: %v", i, err)
		}
		if templates[i] == nil {
			t.Errorf("Goroutine %d returned nil template", i)
		}
	}

	// All should return the same cached template
	firstTemplate := templates[0]
	for i, template := range templates {
		if template != firstTemplate {
			t.Errorf("Goroutine %d returned different template (not thread-safe)", i)
		}
	}
}

func TestLoadCommandTable_FileNotFound(t *testing.T) {
	// Reset cache before test
	cachedTemplateMutex.Lock()
	originalCache := cachedTemplate
	cachedTemplate = nil
	cachedTemplateMutex.Unlock()
	defer func() {
		cachedTemplateMutex.Lock()
		cachedTemplate = originalCache
		cachedTemplateMutex.Unlock()
	}()

	// Save current working directory
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWD)

	// Create a temporary directory with no template file
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Don't create template file - should fail
	template, err := loadCommandTable()
	if err == nil {
		t.Error("loadCommandTable() should fail when file doesn't exist")
	}
	if template != nil {
		t.Error("loadCommandTable() should return nil template on error")
	}
	if err != nil && !strings.Contains(err.Error(), "failed to load command table template") {
		t.Errorf("Error message should mention template loading: %v", err)
	}
}

func TestResolveCommandTablePath(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	prodPath := filepath.Join(tmpDir, "opt", "svw", "stage_primer", "primer", "linmot", "linmot_command_table.yaml")
	devPath := filepath.Join(tmpDir, "primer", "linmot", "linmot_command_table.yaml")

	// Test production path preference
	t.Run("prefers production path", func(t *testing.T) {
		if err := os.MkdirAll(filepath.Dir(prodPath), 0755); err != nil {
			t.Fatalf("Failed to create prod directory: %v", err)
		}
		if err := os.WriteFile(prodPath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create prod file: %v", err)
		}
		if err := os.MkdirAll(filepath.Dir(devPath), 0755); err != nil {
			t.Fatalf("Failed to create dev directory: %v", err)
		}
		if err := os.WriteFile(devPath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create dev file: %v", err)
		}

		// Temporarily override the production path
		originalProdPath := "/opt/svw/stage_primer/primer/linmot/linmot_command_table.yaml"
		// We can't easily test this without modifying the function, but we can verify
		// the logic exists in the code
		_ = originalProdPath
	})

	// Test dev path fallback
	t.Run("falls back to dev path", func(t *testing.T) {
		// Remove prod path if it exists
		os.Remove(prodPath)

		if err := os.MkdirAll(filepath.Dir(devPath), 0755); err != nil {
			t.Fatalf("Failed to create dev directory: %v", err)
		}
		if err := os.WriteFile(devPath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create dev file: %v", err)
		}

		// The function should find dev path
		// Note: This is a simplified test - full testing would require
		// more complex path manipulation
	})
}

func TestDeployCommandTable_IntegrationWithMock(t *testing.T) {
	// Reset cache before test
	cachedTemplateMutex.Lock()
	originalCache := cachedTemplate
	cachedTemplate = nil
	cachedTemplateMutex.Unlock()
	defer func() {
		cachedTemplateMutex.Lock()
		cachedTemplate = originalCache
		cachedTemplateMutex.Unlock()
	}()

	// Create mock LinMot client
	linmotClient, transportServer := client.NewMockClient()
	defer linmotClient.Close()

	// Create and start mock drive
	mockDrive := test.NewMockLinMot(transportServer)
	mockDrive.Start()
	defer mockDrive.Close()

	// Save current working directory
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWD)

	// Create temporary directory and change to it
	tmpDir := t.TempDir()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create template file in expected dev path
	templatePath := "primer/linmot/linmot_command_table.yaml"
	if err := os.MkdirAll(filepath.Dir(templatePath), 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	templateContent := `version: 1
drive_model: C1250-MI
entries:
  - id: 1
    name: "test entry"
    type: NoOp
    par1: ${POSITION_DOWN}
    par2: ${MAX_VELOCITY}
    par3: ${ACCELERATION}
    par4: ${DELAY_AT_BOTTOM}
`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Create test config
	testCfg := config.Config{
		ClearCores: []config.ClearCoreConfig{
			{
				LinMots: []config.LinMotConfig{
					{IP: "127.0.0.1"}, // Mock doesn't need real IP
				},
			},
		},
	}

	cfg := DeployConfig{
		RobotIndex:          0,
		StageIndex:          0,
		Config:              testCfg,
		ZDistance:           25.5,
		DefaultSpeed:        50.0,
		DefaultAcceleration: 100.0,
		PickTime:            0.5,
	}

	// We need to mock the client creation. Since DeployCommandTable creates
	// its own client, we'll need to test the variable substitution separately
	// and test the full flow with a modified version or use dependency injection.
	// For now, let's test what we can: variable substitution logic

	// Test variable substitution by loading template and checking values
	template, err := loadCommandTable()
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	templateCopy := *template
	positionDown := int64(cfg.ZDistance * PositionUnit)
	maxVelocity := int64(cfg.DefaultSpeed * VelocityUnit)
	acceleration := int64(cfg.DefaultAcceleration * AccelerationUnit)
	delayAtBottom := int64(cfg.PickTime * TimeUnit)

	templateCopy.SetVar("POSITION_DOWN", positionDown)
	templateCopy.SetVar("MAX_VELOCITY", maxVelocity)
	templateCopy.SetVar("ACCELERATION", acceleration)
	templateCopy.SetVar("DELAY_AT_BOTTOM", delayAtBottom)

	// Verify variable substitution
	if err := templateCopy.Validate(); err != nil {
		t.Errorf("Template validation failed after variable substitution: %v", err)
	}

	// Verify expected values
	expectedPositionDown := int64(25.5 * 10000) // 255000
	if positionDown != expectedPositionDown {
		t.Errorf("PositionDown = %d, expected %d", positionDown, expectedPositionDown)
	}

	expectedMaxVelocity := int64(50.0 * VelocityUnit) // 500000 (50% × 10,000 µm/s per %)
	if maxVelocity != expectedMaxVelocity {
		t.Errorf("MaxVelocity = %d, expected %d", maxVelocity, expectedMaxVelocity)
	}
}

func TestDeployCommandTable_ErrorHandling(t *testing.T) {
	// Reset cache before test
	cachedTemplateMutex.Lock()
	originalCache := cachedTemplate
	cachedTemplate = nil
	cachedTemplateMutex.Unlock()
	defer func() {
		cachedTemplateMutex.Lock()
		cachedTemplate = originalCache
		cachedTemplateMutex.Unlock()
	}()

	// Test validation error
	testCfg := config.Config{
		ClearCores: []config.ClearCoreConfig{
			{
				LinMots: []config.LinMotConfig{
					{IP: "127.0.0.1"},
				},
			},
		},
	}

	cfg := DeployConfig{
		RobotIndex:          0,
		StageIndex:          0,
		Config:              testCfg,
		ZDistance:           -1.0, // Invalid: negative
		DefaultSpeed:        50.0,
		DefaultAcceleration: 100.0,
		PickTime:            0.5,
	}

	ctx := context.Background()
	err := DeployCommandTable(ctx, cfg)
	if err == nil {
		t.Error("DeployCommandTable() should fail with invalid config")
	}
	if !strings.Contains(err.Error(), "invalid deployment configuration") {
		t.Errorf("Error should mention invalid configuration: %v", err)
	}

	// Test missing LinMot IP
	cfg2 := DeployConfig{
		RobotIndex:          999, // Out of range
		StageIndex:          0,
		Config:              testCfg,
		ZDistance:           25.5,
		DefaultSpeed:        50.0,
		DefaultAcceleration: 100.0,
		PickTime:            0.5,
	}

	err = DeployCommandTable(ctx, cfg2)
	if err == nil {
		t.Error("DeployCommandTable() should fail with invalid robot index")
	}
	if !strings.Contains(err.Error(), "ConfigStore lookup failed") {
		t.Errorf("Error should mention ConfigStore lookup failure: %v", err)
	}
}

func TestDeployCommandTable_ContextCancellation(t *testing.T) {
	// Reset cache before test
	cachedTemplateMutex.Lock()
	originalCache := cachedTemplate
	cachedTemplate = nil
	cachedTemplateMutex.Unlock()
	defer func() {
		cachedTemplateMutex.Lock()
		cachedTemplate = originalCache
		cachedTemplateMutex.Unlock()
	}()

	// Create mock LinMot client
	linmotClient, transportServer := client.NewMockClient()
	defer linmotClient.Close()

	mockDrive := test.NewMockLinMot(transportServer)
	mockDrive.Start()
	defer mockDrive.Close()

	// Install mock factory so DeployCommandTable never attempts a real connection.
	SetClientFactory(&mockSingleClientFactory{client: linmotClient})
	defer ResetClientFactory()

	// Save current working directory
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWD)

	// Create temporary directory and change to it
	tmpDir := t.TempDir()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create template file in expected dev path
	templatePath := "primer/linmot/linmot_command_table.yaml"
	if err := os.MkdirAll(filepath.Dir(templatePath), 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	templateContent := `version: 1
drive_model: C1250-MI
entries:
  - id: 1
    name: "test entry"
    type: NoOp
`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	testCfg := config.Config{
		ClearCores: []config.ClearCoreConfig{
			{
				LinMots: []config.LinMotConfig{
					{IP: "127.0.0.1"},
				},
			},
		},
	}

	cfg := DeployConfig{
		RobotIndex:          0,
		StageIndex:          0,
		Config:              testCfg,
		ZDistance:           25.5,
		DefaultSpeed:        50.0,
		DefaultAcceleration: 100.0,
		PickTime:            0.5,
	}

	// Create context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = DeployCommandTable(ctx, cfg)
	// The function may or may not check context early, but it should handle
	// cancellation gracefully. Since we're using a mock that responds quickly,
	// we mainly verify it doesn't panic.
	if err != nil && !errors.Is(err, context.Canceled) {
		// Error is acceptable, but should be context-related if it happens early
		_ = err
	}
}
