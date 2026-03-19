package linmot

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"

	clearcore "stage_primer_config"
)

// VacuumAction represents the desired vacuum state
type VacuumAction string

const (
	VacuumActionOn    VacuumAction = "on"
	VacuumActionOff   VacuumAction = "off"
	VacuumActionPurge VacuumAction = "purge"
)

// VacuumConfig contains configuration for vacuum control operations
type VacuumConfig struct {
	RobotIndex int
	StageIndex int
	Config     clearcore.Config
	Action     VacuumAction
}

// VacuumState represents the current state of vacuum outputs
type VacuumState struct {
	StageIndex int  `json:"stageIndex"`
	VacuumOn   bool `json:"vacuumOn"`
	PurgeOn    bool `json:"purgeOn"`
}

// SetVacuum controls the vacuum state on a LinMot drive by toggling the
// output pin function between AlwaysOn and None (same as LinMot-Talk's
// "Enable Manual Override" approach).
//
// Pin mapping: X4.3 = purge, X4.4 = vacuum
func SetVacuum(ctx context.Context, cfg VacuumConfig) error {
	linmotIP, err := findLinMotIP(cfg.RobotIndex, cfg.StageIndex, cfg.Config)
	if err != nil {
		return errors.Wrap(err, "failed to find LinMot IP")
	}

	linmotClient, err := globalClientFactory.CreateClient(linmotIP)
	if err != nil {
		return errors.Wrapf(err, "failed to get LinMot client for %s", linmotIP)
	}

	alwaysOn := protocol_common.OutputConfig.AlwaysOn
	none := protocol_common.OutputConfig.None
	ram := protocol_common.ParameterStorage.RAM

	switch cfg.Action {
	case VacuumActionOn:
		fmt.Printf("[LinMot %s] Stage %d: Turning vacuum ON\n", linmotIP, cfg.StageIndex)
		// X4.3 (purge) OFF, X4.4 (vacuum) ON
		if err := linmotClient.SetIODefOutputFunction(ctx, protocol_common.IOPin.Output43, none, ram); err != nil {
			fmt.Printf("[LinMot %s] Stage %d: ERROR setting Output43 to None: %v\n", linmotIP, cfg.StageIndex, err)
			return errors.Wrapf(err, "failed to turn off purge on LinMot %s", linmotIP)
		}
		fmt.Printf("[LinMot %s] Stage %d: Output43 set to None OK\n", linmotIP, cfg.StageIndex)
		if err := linmotClient.SetIODefOutputFunction(ctx, protocol_common.IOPin.Output44, alwaysOn, ram); err != nil {
			fmt.Printf("[LinMot %s] Stage %d: ERROR setting Output44 to AlwaysOn: %v\n", linmotIP, cfg.StageIndex, err)
			return errors.Wrapf(err, "failed to turn on vacuum on LinMot %s", linmotIP)
		}
		fmt.Printf("[LinMot %s] Stage %d: Output44 set to AlwaysOn OK\n", linmotIP, cfg.StageIndex)

	case VacuumActionPurge:
		fmt.Printf("[LinMot %s] Stage %d: Turning vacuum OFF with purge\n", linmotIP, cfg.StageIndex)
		// X4.3 (purge) ON, X4.4 (vacuum) OFF
		if err := linmotClient.SetIODefOutputFunction(ctx, protocol_common.IOPin.Output44, none, ram); err != nil {
			fmt.Printf("[LinMot %s] Stage %d: ERROR setting Output44 to None: %v\n", linmotIP, cfg.StageIndex, err)
			return errors.Wrapf(err, "failed to turn off vacuum on LinMot %s", linmotIP)
		}
		if err := linmotClient.SetIODefOutputFunction(ctx, protocol_common.IOPin.Output43, alwaysOn, ram); err != nil {
			fmt.Printf("[LinMot %s] Stage %d: ERROR setting Output43 to AlwaysOn: %v\n", linmotIP, cfg.StageIndex, err)
			return errors.Wrapf(err, "failed to turn on purge on LinMot %s", linmotIP)
		}

	case VacuumActionOff:
		fmt.Printf("[LinMot %s] Stage %d: Turning vacuum OFF (all off)\n", linmotIP, cfg.StageIndex)
		// Both OFF
		if err := linmotClient.SetIODefOutputFunction(ctx, protocol_common.IOPin.Output43, none, ram); err != nil {
			fmt.Printf("[LinMot %s] Stage %d: ERROR setting Output43 to None: %v\n", linmotIP, cfg.StageIndex, err)
			return errors.Wrapf(err, "failed to turn off purge on LinMot %s", linmotIP)
		}
		if err := linmotClient.SetIODefOutputFunction(ctx, protocol_common.IOPin.Output44, none, ram); err != nil {
			fmt.Printf("[LinMot %s] Stage %d: ERROR setting Output44 to None: %v\n", linmotIP, cfg.StageIndex, err)
			return errors.Wrapf(err, "failed to turn off vacuum on LinMot %s", linmotIP)
		}

	default:
		return errors.Errorf("invalid vacuum action: %s", cfg.Action)
	}

	fmt.Printf("[LinMot %s] Stage %d: Vacuum command sent successfully\n", linmotIP, cfg.StageIndex)
	return nil
}
