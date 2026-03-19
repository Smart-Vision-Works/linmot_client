package linmot

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"

	linmot_command_tables "github.com/Smart-Vision-Works/staged_robot/client/rtc/command_tables"

	"stage_primer_config"
)

// Unit conversion constants for LinMot command tables.
// LinMot uses specific units that differ from our UI-friendly units.
const (
	// PositionUnit converts millimeters to 0.1 micrometer units (LinMot position unit)
	// 1 mm = 10,000 units of 0.1µm
	PositionUnit = 10000

	// SliderHomePositionMM is the resting position after homing (from LinMot-Talk:
	// Parameters → Motion Control SW → Homing → Homing Position Config → Slider Home Position).
	// The command table's "move up" entry returns to this position.
	SliderHomePositionMM = 10.0

	// VelocityUnit converts percent (0-100) to protocol velocity units (1µm/s).
	// 100% = 1 m/s = 1,000,000 µm/s, so each percent = 10,000 µm/s.
	// Reference: Factor.Speed = 1,000,000 (gsail-go/linmot/protocol/common/constants.go)
	VelocityUnit = 10000

	// AccelerationUnit converts percent (0-200) to 1e-5 m/s² units
	// Assumes 100% = 10 m/s² = 1,000,000 in 1e-5 m/s² units
	AccelerationUnit = 10000

	// TimeUnit converts seconds to 100 microsecond units
	// 1 second = 10,000 units of 100µs
	TimeUnit = 10000

	// DefaultPurgeDelayMs is the default delay for purge operations in milliseconds
	DefaultPurgeDelayMs = 500

	// DeploymentTimeout is the maximum time for the full deploy cycle:
	// write entries (~1s) + flash save (~75s observed) + recovery wait (~60s) + homing (~30s) = ~166s
	// Use 300s to accommodate flash save variability (observed 75s, documented 39s).
	DeploymentTimeout = 300 * time.Second

	// Hardware limits for validation
	maxZDistance    = 500.0 // Maximum Z travel in mm
	maxSpeed        = 100.0 // Maximum speed in percent
	maxAcceleration = 200.0 // Maximum acceleration in percent
	minPickTime     = 0.001 // Minimum pick time in seconds (1ms)
	maxPickTime     = 10.0  // Maximum pick time in seconds
)

var (
	// Template cache - loaded once and reused
	cachedTemplate      *linmot_command_tables.CommandTable
	cachedTemplateMutex sync.RWMutex

	// Inspect template cache - loaded once and reused (no vacuum version)
	cachedInspectTemplate      *linmot_command_tables.CommandTable
	cachedInspectTemplateMutex sync.RWMutex
)

// DeployConfig contains configuration for command table deployment
type DeployConfig struct {
	RobotIndex int
	StageIndex int
	Config     config.Config

	// LinMotIP is the direct IP address of the LinMot drive to deploy to.
	// When non-empty, this is used instead of looking up the IP from Config
	// via findLinMotIP(). This is the preferred path because it eliminates
	// the dependency on the ConfigStore being populated by SyncStagePrimers.
	LinMotIP string

	// Motion parameters from robot settings
	ZDistance            float64
	DefaultSpeed         float64
	DefaultAcceleration  float64
	PickTime             float64
}

// resolveCommandTablePath determines the correct path to the LinMot command table template.
// It checks production deployment path first, then development relative path, returning
// the first that exists. If neither exists, returns the production path (will error on
// first use with a clear message).
func resolveCommandTablePath() string {
	// Production path: /opt/svw/stage_primer/primer/linmot/linmot_command_table.yaml
	productionPath := "/opt/svw/stage_primer/primer/linmot/linmot_command_table.yaml"

	// Development path: primer/linmot/linmot_command_table.yaml (relative to working directory)
	devPath := "primer/linmot/linmot_command_table.yaml"

	// Try production path first
	if _, err := os.Stat(productionPath); err == nil {
		return productionPath
	}

	// Try development path
	if _, err := os.Stat(devPath); err == nil {
		return devPath
	}

	// Try relative to current working directory (for tests)
	wd, _ := os.Getwd()
	relativePath := filepath.Join(wd, "primer", "linmot", "linmot_command_table.yaml")
	if _, err := os.Stat(relativePath); err == nil {
		return relativePath
	}

	// Default to production path (will generate clear error on first use)
	return productionPath
}

// resolveInspectCommandTablePath determines the correct path to the LinMot inspect mode
// command table template (no vacuum version). Follows same resolution logic as standard template.
func resolveInspectCommandTablePath() string {
	// Production path
	productionPath := "/opt/svw/stage_primer/primer/linmot/linmot_command_table_inspect.yaml"

	// Development path (relative to working directory)
	devPath := "primer/linmot/linmot_command_table_inspect.yaml"

	// Try production path first
	if _, err := os.Stat(productionPath); err == nil {
		return productionPath
	}

	// Try development path
	if _, err := os.Stat(devPath); err == nil {
		return devPath
	}

	// Try relative to current working directory (for tests)
	wd, _ := os.Getwd()
	relativePath := filepath.Join(wd, "primer", "linmot", "linmot_command_table_inspect.yaml")
	if _, err := os.Stat(relativePath); err == nil {
		return relativePath
	}

	// Default to production path (will generate clear error on first use)
	return productionPath
}

// loadCommandTable loads the command table template from disk with caching.
// Uses double-checked locking for thread-safe lazy initialization.
// Returns a cached template for performance - callers must copy before modifying.
func loadCommandTable() (*linmot_command_tables.CommandTable, error) {
	// Fast path: check if cached (read lock only)
	cachedTemplateMutex.RLock()
	if cachedTemplate != nil {
		cachedTemplateMutex.RUnlock()
		return cachedTemplate, nil
	}
	cachedTemplateMutex.RUnlock()

	// Slow path: load from file (write lock)
	cachedTemplateMutex.Lock()
	defer cachedTemplateMutex.Unlock()

	// Double-check after acquiring write lock (another goroutine may have loaded)
	if cachedTemplate != nil {
		return cachedTemplate, nil
	}

	// Load template
	templatePath := resolveCommandTablePath()
	template, err := linmot_command_tables.Load(templatePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load command table template from %s", templatePath)
	}

	cachedTemplate = template
	return template, nil
}

// loadInspectCommandTable loads the inspect mode command table template from disk with caching.
// The inspect template provides motion-only commands without vacuum/purge control.
// Uses double-checked locking for thread-safe lazy initialization.
// Returns a cached template for performance - callers must copy before modifying.
func loadInspectCommandTable() (*linmot_command_tables.CommandTable, error) {
	// Fast path: check if cached (read lock only)
	cachedInspectTemplateMutex.RLock()
	if cachedInspectTemplate != nil {
		cachedInspectTemplateMutex.RUnlock()
		return cachedInspectTemplate, nil
	}
	cachedInspectTemplateMutex.RUnlock()

	// Slow path: load from file (write lock)
	cachedInspectTemplateMutex.Lock()
	defer cachedInspectTemplateMutex.Unlock()

	// Double-check after acquiring write lock (another goroutine may have loaded)
	if cachedInspectTemplate != nil {
		return cachedInspectTemplate, nil
	}

	// Load template
	templatePath := resolveInspectCommandTablePath()
	template, err := linmot_command_tables.Load(templatePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load inspect command table template from %s", templatePath)
	}

	cachedInspectTemplate = template
	return template, nil
}

// validateConfig validates the deployment configuration parameters.
// Returns error if any parameter is out of acceptable range.
func validateConfig(cfg DeployConfig) error {
	// Check for NaN or Inf values first (before range checks)
	if math.IsNaN(cfg.ZDistance) || math.IsInf(cfg.ZDistance, 0) {
		return fmt.Errorf("z_distance is not a valid number: %v", cfg.ZDistance)
	}
	if math.IsNaN(cfg.DefaultSpeed) || math.IsInf(cfg.DefaultSpeed, 0) {
		return fmt.Errorf("default_speed is not a valid number: %v", cfg.DefaultSpeed)
	}
	if math.IsNaN(cfg.DefaultAcceleration) || math.IsInf(cfg.DefaultAcceleration, 0) {
		return fmt.Errorf("default_acceleration is not a valid number: %v", cfg.DefaultAcceleration)
	}
	if math.IsNaN(cfg.PickTime) || math.IsInf(cfg.PickTime, 0) {
		return fmt.Errorf("pick_time is not a valid number: %v", cfg.PickTime)
	}

	// Validate Z distance
	if cfg.ZDistance < 0 {
		return fmt.Errorf("z_distance cannot be negative: %.2f mm", cfg.ZDistance)
	}
	if cfg.ZDistance > maxZDistance {
		return fmt.Errorf("z_distance exceeds maximum travel: %.2f mm > %.2f mm", cfg.ZDistance, maxZDistance)
	}

	// Validate speed
	if cfg.DefaultSpeed < 0 || cfg.DefaultSpeed > maxSpeed {
		return fmt.Errorf("default_speed out of range: %.1f%% (valid: 0-%.1f%%)", cfg.DefaultSpeed, maxSpeed)
	}

	// Validate acceleration
	if cfg.DefaultAcceleration < 0 || cfg.DefaultAcceleration > maxAcceleration {
		return fmt.Errorf("default_acceleration out of range: %.1f%% (valid: 0-%.1f%%)", cfg.DefaultAcceleration, maxAcceleration)
	}

	// Validate pick time
	if cfg.PickTime < minPickTime || cfg.PickTime > maxPickTime {
		return fmt.Errorf("pick_time out of range: %.3f s (valid: %.3f-%.1f s)", cfg.PickTime, minPickTime, maxPickTime)
	}

	return nil
}

// saveFlashAndRecover performs the full post-write recovery lifecycle:
//
//  1. Send SaveCommandTable command (fire-and-forget — the drive does not send a
//     standard RTC response for this command, so we expect a timeout).
//  2. Poll GetStatus until the drive's UDP interface recovers.
//  3. Restart the Motion Controller (MC was stopped for the write).
//  4. Send Home command and wait for homing to complete (establishes position
//     reference for VAI_GoToPos commands in the command table).
//
// IMPORTANT: The caller must manage recovery state (EnterRecoveryState /
// ExitRecoveryState). This function assumes the fault monitor is already
// paused for this IP and will remain paused until the caller exits recovery.
func saveFlashAndRecover(ctx context.Context, linmotIP string) error {
	const (
		homingPollInterval = 2 * time.Second
		maxHomingAttempts  = 15 // 15 × 2s = 30s max wait for homing
	)

	// Phase 1: Send SaveCommandTable command. The drive will save to flash but
	// does not send a standard RTC response — the command will time out. This is
	// expected. The save is confirmed later by reading back the command table.
	flashClient, err := globalClientFactory.CreateClient(linmotIP)
	if err != nil {
		return fmt.Errorf("failed to create client for flash save: %w", err)
	}

	fmt.Printf("[LinMot %s] Sending SaveCommandTable (fire-and-forget, timeout expected)...\n", linmotIP)
	flashCtx, flashCancel := context.WithTimeout(ctx, 30*time.Second)
	saveErr := flashClient.SaveCommandTableToFlash(flashCtx)
	flashCancel()
	if saveErr != nil {
		// Timeout is expected — the command was sent, drive is processing it.
		fmt.Printf("[LinMot %s] SaveCommandTable returned (expected timeout): %v\n", linmotIP, saveErr)
	} else {
		// Unexpected success — drive actually responded. Great!
		fmt.Printf("[LinMot %s] SaveCommandTable responded successfully (unexpected but welcome)\n", linmotIP)
	}

	// Phase 2: Wait for drive's UDP interface to recover after flash save.
	if err := waitForDriveRecovery(ctx, linmotIP); err != nil {
		return fmt.Errorf("drive recovery failed: %w", err)
	}

	// Phase 3: Restart MC. It was stopped for the write+flash cycle.
	freshClient, err := globalClientFactory.CreateClient(linmotIP)
	if err != nil {
		return fmt.Errorf("failed to create client for MC restart: %w", err)
	}

	fmt.Printf("[LinMot %s] Restarting motion controller...\n", linmotIP)
	restartCtx, restartCancel := context.WithTimeout(ctx, 10*time.Second)
	err = freshClient.StartMotionController(restartCtx)
	restartCancel()
	if err != nil {
		return fmt.Errorf("MC restart failed: %w", err)
	}
	fmt.Printf("[LinMot %s] Motion controller restarted successfully\n", linmotIP)

	// Phase 3: Send Home command to trigger the configured homing sequence.
	fmt.Printf("[LinMot %s] Sending Home command...\n", linmotIP)
	homeCtx, homeCancel := context.WithTimeout(ctx, 10*time.Second)
	_, err = freshClient.Home(homeCtx)
	homeCancel()
	if err != nil {
		return fmt.Errorf("Home command failed: %w", err)
	}

	// Phase 4: Wait for homing to complete (status word bit 11 = Homed)
	fmt.Printf("[LinMot %s] Waiting for homing to complete...\n", linmotIP)
	for i := range maxHomingAttempts {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during homing wait")
		case <-time.After(homingPollInterval):
		}

		pollCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		status, err := freshClient.GetStatus(pollCtx)
		cancel()
		if err != nil {
			fmt.Printf("[LinMot %s] Homing poll %d/%d: status read failed (%v)\n",
				linmotIP, i+1, maxHomingAttempts, err)
			continue
		}

		homedBit := (status.StatusWord >> 11) & 1
		mainState := (status.StateVar >> 8) & 0xFF
		fmt.Printf("[LinMot %s] Homing poll %d/%d: homed=%d, state=%d, warn=0x%04X\n",
			linmotIP, i+1, maxHomingAttempts, homedBit, mainState, status.WarnWord)

		if homedBit == 1 {
			fmt.Printf("[LinMot %s] Homing completed successfully!\n", linmotIP)
			return nil
		}
	}

	return fmt.Errorf("homing did not complete within %d attempts (motor may need physical intervention)", maxHomingAttempts)
}

// resolveLinMotIP determines the LinMot IP address for a deployment.
// It prefers the direct LinMotIP from the DeployConfig (passed by the caller
// from Redis robot_settings). If not provided, it falls back to looking up
// the IP from the ConfigStore via findLinMotIP.
func resolveLinMotIP(cfg DeployConfig) (string, error) {
	if cfg.LinMotIP != "" {
		fmt.Printf("[LinMot] Deploy target: %s (direct from caller, robot=%d, stage=%d)\n",
			cfg.LinMotIP, cfg.RobotIndex, cfg.StageIndex)
		return cfg.LinMotIP, nil
	}

	// Fallback: look up from ConfigStore (legacy path)
	linmotIP, err := findLinMotIP(cfg.RobotIndex, cfg.StageIndex, cfg.Config)
	if err != nil {
		return "", errors.Wrap(err, "no direct linmot_ip provided and ConfigStore lookup failed")
	}
	fmt.Printf("[LinMot] Deploy target: %s (ConfigStore fallback, robot=%d, stage=%d)\n",
		linmotIP, cfg.RobotIndex, cfg.StageIndex)
	return linmotIP, nil
}

// DeployCommandTable deploys a command table to a LinMot drive.
// It loads the template (with caching), validates input parameters, substitutes variables,
// and uploads the command table to the hardware.
//
// Thread-safety: Safe for concurrent calls. Template caching uses double-checked locking.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - cfg: Deployment configuration with robot/stage indices and motion parameters
//
// Returns error if validation fails, template loading fails, or deployment fails.
func DeployCommandTable(ctx context.Context, cfg DeployConfig) error {
	// 1. Validate input parameters
	if err := validateConfig(cfg); err != nil {
		return errors.Wrap(err, "invalid deployment configuration")
	}

	// 2. Resolve LinMot IP — prefer direct IP from caller, fall back to ConfigStore
	linmotIP, err := resolveLinMotIP(cfg)
	if err != nil {
		return err
	}

	// 3. Load template from yaml file (cached after first load)
	template, err := loadCommandTable()
	if err != nil {
		return err // Error already wrapped by loadCommandTable
	}

	// 4. Copy template to avoid modifying cached version
	templateCopy := *template

	// 5. Map settings to variables
	// Position variables: z_distance in mm → 0.1µm units (absolute positions)
	// POSITION_DOWN: the pick position (Z distance from the UI)
	// POSITION_UP: the Slider Home Position (10mm) — resting position after homing
	positionDown := int64(cfg.ZDistance * PositionUnit)
	positionUp := int64(SliderHomePositionMM * PositionUnit)

	templateCopy.SetVar("POSITION_DOWN", positionDown)
	templateCopy.SetVar("POSITION_UP", positionUp)

	// Velocity: default_speed is in percent (0-100), convert to µm/s
	// Assumes 100% = 10 mm/s = 10,000 µm/s, so percent * VelocityUnit
	maxVelocity := int64(cfg.DefaultSpeed * VelocityUnit)
	templateCopy.SetVar("MAX_VELOCITY", maxVelocity)

	// Acceleration/Deceleration: default_acceleration is in percent (0-200)
	// Convert to 1e-5 m/s² units. Assumes 100% = 10 m/s² = 1,000,000 units
	acceleration := int64(cfg.DefaultAcceleration * AccelerationUnit)
	templateCopy.SetVar("ACCELERATION", acceleration)
	templateCopy.SetVar("DECELERATION", acceleration) // Use same value for both

	// Delay at bottom: pick_time is in seconds, convert to 100µs units
	// 1 second = TimeUnit (10,000 units of 100µs each)
	delayAtBottom := int64(cfg.PickTime * TimeUnit)
	templateCopy.SetVar("DELAY_AT_BOTTOM", delayAtBottom)

	// Delay for purge: fixed value of 500ms
	delayPurge := int64(DefaultPurgeDelayMs * TimeUnit / 1000)
	templateCopy.SetVar("DELAY_PURGE", delayPurge)

	// 6. Validate after variable binding
	if err := templateCopy.Validate(); err != nil {
		return errors.Wrap(err, "command table validation failed after variable binding")
	}

	// 7. Create client and deploy with timeout
	deployCtx, cancel := context.WithTimeout(ctx, DeploymentTimeout)
	defer cancel()

	// Get client from pool (managed by factory, do NOT close)
	linmotClient, err := globalClientFactory.CreateClient(linmotIP)
	if err != nil {
		return errors.Wrapf(err, "failed to get LinMot client for %s", linmotIP)
	}

	fmt.Printf("[LinMot %s] Deploying command table: Z=%.2fmm (%d raw), speed=%.1f%% (%d raw), accel=%.1f%% (%d raw), pickTime=%.3fs (%d raw), entries=%d\n",
		linmotIP, cfg.ZDistance, positionDown, cfg.DefaultSpeed, maxVelocity,
		cfg.DefaultAcceleration, acceleration, cfg.PickTime, delayAtBottom, len(templateCopy.Entries))

	// Enter recovery state BEFORE the deploy. This pauses the fault monitor
	// for this IP during the entire flash save + recovery + MC restart + homing
	// cycle. The fault monitor only resumes once the drive is verified healthy.
	// Without this, the 750ms fault poll timeout fires every second during the
	// 75+ second flash save, flooding the UI with spurious "context deadline exceeded".
	EnterRecoveryState(linmotIP)

	deployStart := time.Now()

	// Write entries to RAM only — we manage flash save and MC restart ourselves.
	// SkipFlashSave: the LinMot drive does not send a standard RTC response for
	// SaveCommandTable, causing our request/response layer to always time out.
	// Instead, we send the save command fire-and-forget, wait for the drive to
	// recover, then restart MC and home.
	err = linmotClient.SetCommandTableWithOptions(deployCtx, &templateCopy, linmot_command_tables.SetCommandTableOptions{
		RestartMC:     false,
		SkipFlashSave: true,
	})
	if err != nil && !errors.Is(err, linmot_command_tables.ErrCommandTableUnchanged) {
		ExitRecoveryState(linmotIP)
		return errors.Wrapf(err, "failed to deploy command table to %s", linmotIP)
	}

	if errors.Is(err, linmot_command_tables.ErrCommandTableUnchanged) {
		// Entries already in RAM and flash — no MC stop/restart occurred, drive is healthy.
		ExitRecoveryState(linmotIP)
		fmt.Printf("[LinMot %s] Command table unchanged — skipping recovery (verified in %dms)\n",
			linmotIP, time.Since(deployStart).Milliseconds())
	} else {
		// Entries written to RAM. MC is stopped. Now persist to flash and recover.
		// The full sequence: save to flash → wait for drive → restart MC → home.
		fmt.Printf("[LinMot %s] Entries written to RAM in %dms — saving to flash and recovering\n",
			linmotIP, time.Since(deployStart).Milliseconds())
		if err := saveFlashAndRecover(deployCtx, linmotIP); err != nil {
			fmt.Printf("[LinMot %s] WARNING: post-deploy recovery incomplete: %v\n", linmotIP, err)
		}
		// Drive verified healthy (or recovery timed out) — resume fault monitoring.
		ExitRecoveryState(linmotIP)
	}

	fmt.Printf("[LinMot %s] Command table deployed successfully (robot=%d, stage=%d, Z=%.2fmm, speed=%.1f%%, accel=%.1f%%, pickTime=%.3fs, entries=%d, total=%dms)\n",
		linmotIP, cfg.RobotIndex, cfg.StageIndex, cfg.ZDistance, cfg.DefaultSpeed, cfg.DefaultAcceleration, cfg.PickTime, len(templateCopy.Entries), time.Since(deployStart).Milliseconds())
	return nil
}

// DeployInspectCommandTable deploys a command table for encoder inspection mode to a LinMot drive.
// This version provides motion-only commands WITHOUT vacuum/purge control, suitable for
// testing and encoder inspection without activating pneumatics.
//
// Thread-safety: Safe for concurrent calls. Template caching uses double-checked locking.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - cfg: Deployment configuration with robot/stage indices and motion parameters
//
// Returns error if validation fails, template loading fails, or deployment fails.
func DeployInspectCommandTable(ctx context.Context, cfg DeployConfig) error {
	// 1. Validate input parameters
	if err := validateConfig(cfg); err != nil {
		return errors.Wrap(err, "invalid deployment configuration")
	}

	// 2. Resolve LinMot IP — prefer direct IP from caller, fall back to ConfigStore
	linmotIP, err := resolveLinMotIP(cfg)
	if err != nil {
		return err
	}

	// 3. Load inspect template from yaml file (cached after first load)
	template, err := loadInspectCommandTable()
	if err != nil {
		return err // Error already wrapped by loadInspectCommandTable
	}

	// 4. Copy template to avoid modifying cached version
	templateCopy := *template

	// 5. Map settings to variables (same as normal deploy)
	// Position variables: absolute positions in 0.1µm units
	positionDown := int64(cfg.ZDistance * PositionUnit)
	positionUp := int64(SliderHomePositionMM * PositionUnit)

	templateCopy.SetVar("POSITION_DOWN", positionDown)
	templateCopy.SetVar("POSITION_UP", positionUp)

	// Velocity: default_speed is in percent (0-100), convert to µm/s
	maxVelocity := int64(cfg.DefaultSpeed * VelocityUnit)
	templateCopy.SetVar("MAX_VELOCITY", maxVelocity)

	// Acceleration/Deceleration: default_acceleration is in percent (0-200)
	acceleration := int64(cfg.DefaultAcceleration * AccelerationUnit)
	templateCopy.SetVar("ACCELERATION", acceleration)
	templateCopy.SetVar("DECELERATION", acceleration)

	// Delay at bottom: pick_time is in seconds, convert to 100µs units
	delayAtBottom := int64(cfg.PickTime * TimeUnit)
	templateCopy.SetVar("DELAY_AT_BOTTOM", delayAtBottom)

	// 6. Validate after variable binding
	if err := templateCopy.Validate(); err != nil {
		return errors.Wrap(err, "inspect command table validation failed after variable binding")
	}

	// 7. Create client and deploy with timeout
	deployCtx, cancel := context.WithTimeout(ctx, DeploymentTimeout)
	defer cancel()

	// Create client using factory (allows injection of mocks for testing)
	linmotClient, err := globalClientFactory.CreateClient(linmotIP)
	if err != nil {
		return errors.Wrapf(err, "failed to create LinMot client for %s", linmotIP)
	}

	fmt.Printf("[LinMot %s] Deploying inspect command table: Z=%.2fmm (%d raw), speed=%.1f%%, accel=%.1f%%, pickTime=%.3fs, entries=%d\n",
		linmotIP, cfg.ZDistance, positionDown, cfg.DefaultSpeed, cfg.DefaultAcceleration, cfg.PickTime, len(templateCopy.Entries))

	// Pause fault monitoring during deploy (flash save can take 75+ seconds).
	EnterRecoveryState(linmotIP)

	deployStart := time.Now()
	if err := linmotClient.SetCommandTable(deployCtx, &templateCopy); err != nil {
		ExitRecoveryState(linmotIP)
		return errors.Wrapf(err, "failed to deploy inspect command table to %s", linmotIP)
	}
	deployDuration := time.Since(deployStart)

	// Drive is healthy after deploy — resume fault monitoring.
	ExitRecoveryState(linmotIP)

	// Homing is handled by Entry 12 (Master Homing) at the start of every command
	// table execution. No separate Home() call needed.

	fmt.Printf("[LinMot %s] Inspect command table deployed successfully in %dms (robot=%d, stage=%d, Z=%.2fmm, speed=%.1f%%, accel=%.1f%%, pickTime=%.3fs, entries=%d)\n",
		linmotIP, deployDuration.Milliseconds(), cfg.RobotIndex, cfg.StageIndex, cfg.ZDistance, cfg.DefaultSpeed, cfg.DefaultAcceleration, cfg.PickTime, len(templateCopy.Entries))
	return nil
}
