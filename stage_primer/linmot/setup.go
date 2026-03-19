package linmot

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
)

// SetupConfig contains configuration for setting up a LinMot drive.
type SetupConfig struct {
	IP string
}

// setupParameter defines a single drive parameter with its UPID and desired value.
type setupParameter struct {
	UPID  uint16
	Value int32
	Name  string
}

// desiredParameters returns the full list of parameters that Setup configures.
// Every LinMot drive in the system must have these values for the triggered
// command table architecture to work end-to-end.
func desiredParameters() []setupParameter {
	return []setupParameter{
		{uint16(protocol_common.PUID.RunMode), int32(protocol_common.RunModes.TriggeredCommandTable), "Run Mode"},
		{uint16(protocol_common.PUID.EasyStepsAutoStart), protocol_common.EasyStepsAutoStart.Enabled, "Easy Steps Auto Start"},
		{uint16(protocol_common.PUID.EasyStepsAutoHome), protocol_common.EasyStepsAutoHome.Enabled, "Easy Steps Auto Home"},
		{uint16(protocol_common.PUID.Output43Function), protocol_common.OutputConfig.InterfaceOutput, "Output 4.3 Function"},
		{uint16(protocol_common.PUID.Output44Function), protocol_common.OutputConfig.InterfaceOutput, "Output 4.4 Function"},
		{uint16(protocol_common.PUID.Input46Function), protocol_common.InputFunction.Trigger, "Input 4.6 Function"},
		{uint16(protocol_common.PUID.TriggerMode), protocol_common.TriggerModeConfig.Direct, "Trigger Mode"},
		{uint16(protocol_common.PUID.EasySteps46RisingEdge), protocol_common.EasyStepsIOMotion.EvalCommandTableCommand, "Easy Steps 4.6 Rising Edge"},
		{uint16(protocol_common.PUID.EasySteps46IOMotionConfigCmd), 1, "Easy Steps 4.6 Command Table Start"},
		// State machine positions — motor rests at 10mm after homing (away from
		// mechanical stop at -1mm to prevent overheating). These are critical for
		// the drive to reach Operation Enabled where triggers are evaluated.
		{0x1725, 100000, "Go To Position (10mm)"},          // State Machine → Go To Position
		{0x13D9, 100000, "Initial Position (10mm)"},         // Homing → Initial Position Config
		{0x30D7, 1, "Intf Go To Initial Pos Flag (Enter Op Enabled)"}, // EasySteps → Smart Control Word
		// Triggered Command Table entry IDs — the MC-level trigger mechanism
		// determines which command table entry to execute on rise/fall edges of X4.
		// Rise = Entry 1 (move down + vacuum), Fall = Entry 5 (move up + purge).
		{0x1486, 1, "Rise Command Table Entry ID"},  // Triggered CT Settings → Rise Entry
		{0x1487, 5, "Fall Command Table Entry ID"},   // Triggered CT Settings → Fall Entry
	}
}

const (
	// recoveryPollInterval is how often we poll GetStatus during recovery.
	recoveryPollInterval = 5 * time.Second

	// maxRecoveryAttempts limits recovery polling (60 × 5s = 5 minutes).
	// ROM writes and flash saves can leave drives unresponsive for 30-90+ seconds.
	maxRecoveryAttempts = 60
)

// Setup performs idempotent, compare-before-write configuration of a LinMot drive.
//
// For each of the 12 required parameters, the current RAM value is read and
// compared with the desired value. Only parameters that differ are written to
// RAM+ROM, minimising flash wear and avoiding unnecessary drive recovery time.
//
// If any ROM writes occurred, the pooled UDP client is evicted and Setup waits
// for the drive to become responsive again before returning.
func Setup(ctx context.Context, cfg SetupConfig) error {
	if cfg.IP == "" {
		return fmt.Errorf("LinMot IP address is required")
	}

	linmotClient, err := globalClientFactory.CreateClient(cfg.IP)
	if err != nil {
		return errors.Wrapf(err, "failed to get LinMot client for %s", cfg.IP)
	}

	params := desiredParameters()
	fmt.Printf("[LinMot %s] Starting setup (%d parameters)...\n", cfg.IP, len(params))

	var written int
	for _, p := range params {
		current, err := linmotClient.ReadRAM(ctx, p.UPID)
		if err != nil {
			return errors.Wrapf(err, "failed to read %s (UPID 0x%04X) on LinMot %s", p.Name, p.UPID, cfg.IP)
		}

		if current == p.Value {
			fmt.Printf("[LinMot %s] ✓ %s already correct (%d)\n", cfg.IP, p.Name, current)
			continue
		}

		fmt.Printf("[LinMot %s] Writing %s: %d → %d\n", cfg.IP, p.Name, current, p.Value)
		if err := linmotClient.WriteRAMAndROM(ctx, p.UPID, p.Value); err != nil {
			return errors.Wrapf(err, "failed to write %s (UPID 0x%04X) on LinMot %s", p.Name, p.UPID, cfg.IP)
		}
		fmt.Printf("[LinMot %s] ✓ %s updated\n", cfg.IP, p.Name)
		written++
	}

	if written == 0 {
		fmt.Printf("[LinMot %s] Setup complete — drive already configured correctly\n", cfg.IP)
		return nil
	}

	fmt.Printf("[LinMot %s] Setup wrote %d parameter(s) to ROM — waiting for drive recovery\n", cfg.IP, written)

	// ROM writes may disrupt the drive's LinUDP stack. Enter recovery state so
	// the fault monitor stops polling this IP entirely, then wait for the drive
	// to come back before resuming normal operation.
	EnterRecoveryState(cfg.IP)
	recoveryErr := waitForDriveRecovery(ctx, cfg.IP)
	ExitRecoveryState(cfg.IP)
	if recoveryErr != nil {
		return errors.Wrapf(recoveryErr, "drive recovery failed after setup on LinMot %s", cfg.IP)
	}

	fmt.Printf("[LinMot %s] Setup completed successfully (%d parameter(s) updated)\n", cfg.IP, written)
	return nil
}

// waitForDriveRecovery evicts the stale pooled client once, then polls
// GetStatus until the drive responds. Used after any operation that writes to
// ROM (Setup, flash save) which can disrupt the LinUDP interface.
func waitForDriveRecovery(ctx context.Context, ip string) error {
	// Evict once — the socket that was open during ROM writes is poisoned.
	// All subsequent CreateClient calls will get a fresh connection from the pool.
	EvictPooledClient(ip)

	fmt.Printf("[LinMot %s] Waiting for drive to recover (polling every %v, up to %v)...\n",
		ip, recoveryPollInterval, time.Duration(maxRecoveryAttempts)*recoveryPollInterval)

	for i := 0; i < maxRecoveryAttempts; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during recovery wait")
		case <-time.After(recoveryPollInterval):
		}

		client, err := globalClientFactory.CreateClient(ip)
		if err != nil {
			fmt.Printf("[LinMot %s] Recovery poll %d/%d: can't create client (%v)\n",
				ip, i+1, maxRecoveryAttempts, err)
			continue
		}

		pollCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		status, err := client.GetStatus(pollCtx)
		cancel()
		if err != nil {
			fmt.Printf("[LinMot %s] Recovery poll %d/%d: not yet (%v)\n",
				ip, i+1, maxRecoveryAttempts, err)
			continue
		}

		fmt.Printf("[LinMot %s] Recovery poll %d/%d: drive responding (status=0x%04X, state=0x%04X)\n",
			ip, i+1, maxRecoveryAttempts, status.StatusWord, status.StateVar)
		return nil
	}

	return fmt.Errorf("drive did not recover within %v", time.Duration(maxRecoveryAttempts)*recoveryPollInterval)
}

// SetupAll performs setup for multiple LinMot drives concurrently.
func SetupAll(ctx context.Context, configs []SetupConfig) error {
	if len(configs) == 0 {
		return nil
	}

	type result struct {
		ip  string
		err error
	}

	results := make(chan result, len(configs))

	for _, cfg := range configs {
		go func(cfg SetupConfig) {
			err := Setup(ctx, cfg)
			results <- result{ip: cfg.IP, err: err}
		}(cfg)
	}

	var errs []error
	for i := 0; i < len(configs); i++ {
		res := <-results
		if res.err != nil {
			errs = append(errs, res.err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to setup %d LinMot drive(s): %v", len(errs), errs)
	}

	return nil
}
