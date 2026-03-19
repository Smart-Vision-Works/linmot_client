package linmot

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_control_word "github.com/Smart-Vision-Works/staged_robot/protocol/control_word"

	"stage_primer_config"
)

const jogPostEnableSettleDelay = 200 * time.Millisecond

func waitForCtx(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// JogConfig contains configuration for jogging a LinMot drive (API usage)
type JogConfig struct {
	RobotIndex int
	StageIndex int
	Config     config.Config
	Position   float64 // Z-axis position in mm to set
}

// PositionConfig contains configuration for getting position from a LinMot drive (API usage)
type PositionConfig struct {
	RobotIndex int
	StageIndex int
	Config     config.Config
}

// Jog performs an absolute Z-axis movement on a LinMot drive
func Jog(ctx context.Context, cfg JogConfig) error {
	// Find LinMot IP from config
	linmotIP, err := findLinMotIP(cfg.RobotIndex, cfg.StageIndex, cfg.Config)
	if err != nil {
		return errors.Wrap(err, "failed to find LinMot IP")
	}

	// Get client from pool (managed by factory, do NOT close)
	linmotClient, err := globalClientFactory.CreateClient(linmotIP)
	if err != nil {
		return errors.Wrapf(err, "failed to get LinMot client for %s", linmotIP)
	}

	// Check current status before acting
	status, err := linmotClient.GetStatus(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to get status before jog on %s", linmotIP)
	}
	statusHelper := protocol_control_word.NewStatusWordHelper(status.StatusWord, status.StateVar)

	// Only act if the drive is not operation-enabled or currently reports a fault.
	if !statusHelper.IsOperationEnabled() || statusHelper.HasError() {
		if statusHelper.HasError() {
			// Drive is in fault — acknowledge and clear it.
			fmt.Printf("[LinMot %s] Drive in fault (Status: 0x%04X). Attempting Acknowledge...\n", linmotIP, status.StatusWord)
			ackStatus, err := linmotClient.AcknowledgeError(ctx)
			if err != nil {
				return errors.Wrapf(err, "failed to acknowledge error before jog on LinMot %s", linmotIP)
			}
			if ackStatus == nil {
				if _, err := linmotClient.EnableDrive(ctx); err != nil {
					return errors.Wrapf(err, "failed to enable drive after acknowledge on LinMot %s", linmotIP)
				}
			} else {
				ackHelper := protocol_control_word.NewStatusWordHelper(ackStatus.StatusWord, ackStatus.StateVar)
				if !ackHelper.IsOperationEnabled() {
					fmt.Printf("[LinMot %s] Drive not yet enabled after acknowledge (Status: 0x%04X). Enabling...\n", linmotIP, ackStatus.StatusWord)
					if _, err := linmotClient.EnableDrive(ctx); err != nil {
						return errors.Wrapf(err, "failed to enable drive after acknowledge on LinMot %s", linmotIP)
					}
				}
			}
		} else {
			// Drive is not enabled but has no fault (e.g. after power cycle) — enable it directly.
			fmt.Printf("[LinMot %s] Drive not enabled (Status: 0x%04X). Enabling...\n", linmotIP, status.StatusWord)
			if _, err := linmotClient.EnableDrive(ctx); err != nil {
				return errors.Wrapf(err, "failed to enable drive before jog on LinMot %s", linmotIP)
			}
		}
		// Give the drive a brief hardware-settle window after acknowledge/enable before sending motion setup writes.
		if err := waitForCtx(ctx, jogPostEnableSettleDelay); err != nil {
			return errors.Wrapf(err, "jog canceled while waiting for LinMot %s to settle after enable", linmotIP)
		}
	}

	// Use VAI2PosContinuous for point-to-point motion by forcing a mode transition.
	// 1. Set the target position in BOTH Position 1 and Position 2.
	// This ensures that regardless of whether the drive evaluates the Special Mode
	// bit as high or low, it seeks our desired position.
	if err := linmotClient.SetPosition1(ctx, cfg.Position, protocol_common.ParameterStorage.RAM); err != nil {
		return errors.Wrapf(err, "failed to set position 1 on LinMot %s", linmotIP)
	}
	if err := linmotClient.SetPosition2(ctx, cfg.Position, protocol_common.ParameterStorage.RAM); err != nil {
		return errors.Wrapf(err, "failed to set position 2 on LinMot %s", linmotIP)
	}

	// 2. Exit VAI2PosContinuous by switching to MotionCommandInterface first.
	// This rearms the subsequent transition back into VAI2PosContinuous.
	if err := linmotClient.SetRunMode(ctx, protocol_common.RunModes.MotionCommandInterface, protocol_common.ParameterStorage.RAM); err != nil {
		return errors.Wrapf(err, "failed to clear run mode on LinMot %s", linmotIP)
	}
	if err := waitForCtx(ctx, 50*time.Millisecond); err != nil {
		return errors.Wrapf(err, "jog canceled while waiting for run-mode transition on LinMot %s", linmotIP)
	}

	// 3. Re-enter VAI2PosContinuous. That mode transition triggers the drive to seek Position 1.
	if err := linmotClient.SetRunMode(ctx, protocol_common.RunModes.VAI2PosContinuous, protocol_common.ParameterStorage.RAM); err != nil {
		return errors.Wrapf(err, "failed to set positioning mode on LinMot %s", linmotIP)
	}

	fmt.Printf("[LinMot %s] Z-jog command sent successfully\n", linmotIP)
	return nil
}

// GetPosition retrieves the current position from a LinMot drive
func GetPosition(ctx context.Context, cfg PositionConfig) (float64, error) {
	// Find LinMot IP from config
	linmotIP, err := findLinMotIP(cfg.RobotIndex, cfg.StageIndex, cfg.Config)
	if err != nil {
		return 0, errors.Wrap(err, "failed to find LinMot IP")
	}

	// Get client from pool (managed by factory, do NOT close)
	linmotClient, err := globalClientFactory.CreateClient(linmotIP)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get LinMot client for %s", linmotIP)
	}

	// Get current position
	position, err := linmotClient.GetPosition(ctx)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get position from LinMot %s", linmotIP)
	}

	fmt.Printf("[LinMot %s] Current position: %.2f mm\n", linmotIP, position)
	return position, nil
}
