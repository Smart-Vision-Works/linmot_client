package client

import (
	"fmt"
	"strings"

	protocol_common "gsail-go/linmot/protocol/common"
)

type ErrDriveNotOperational struct {
	FirmwareRunning bool
	MotorSwitchedOn bool
	OperationEnable bool
	EnableOperation bool
	FatalError      bool
	ErrorActive     bool
	StatusWord      uint16
	StateVar        uint16
	ErrorCode       uint16
	WarnWord        uint16
	StateMachine    string
}

func (e ErrDriveNotOperational) Error() string {
	reasons := []string{}
	if !e.FirmwareRunning {
		reasons = append(reasons, "Firmware Stopped")
	}
	if !e.MotorSwitchedOn {
		reasons = append(reasons, "Motor Switched Off")
	}
	if e.FatalError {
		reasons = append(reasons, "Fatal Error")
	}
	if e.ErrorActive {
		reasons = append(reasons, "Error Active")
	}
	if len(reasons) == 0 {
		reasons = append(reasons, "Drive Not Operational")
	}

	return fmt.Sprintf("linmot drive not operational: %s (StateMachine=%s StatusWord=0x%04X StateVar=0x%04X ErrorCode=0x%04X WarnWord=0x%04X)",
		strings.Join(reasons, "; "), e.StateMachine, e.StatusWord, e.StateVar, e.ErrorCode, e.WarnWord)
}

func assertDriveOperational(status *protocol_common.Status) error {
	if status == nil {
		return ErrDriveNotOperational{
			FirmwareRunning: false,
			MotorSwitchedOn: false,
			StateMachine:    "Unknown",
		}
	}

	state := uint8(status.StateVar >> 8)
	statusLow := uint8(status.StatusWord & 0x00FF)
	statusHigh := uint8(status.StatusWord >> 8)

	firmwareRunning := state != 0
	motorSwitchedOn := statusLow&0x02 != 0
	operationEnable := statusLow&0x01 != 0
	enableOperation := statusLow&0x04 != 0
	fatalError := statusHigh&0x10 != 0
	errorActive := statusLow&0x08 != 0

	stateName := stateMachineStateName(state)
	stateNotOperational := state == 0 || state == 1 || state == 3 || state == 4 || state == 5

	if !firmwareRunning || !motorSwitchedOn || fatalError || errorActive || stateNotOperational {
		return ErrDriveNotOperational{
			FirmwareRunning: firmwareRunning,
			MotorSwitchedOn: motorSwitchedOn,
			OperationEnable: operationEnable,
			EnableOperation: enableOperation,
			FatalError:      fatalError,
			ErrorActive:     errorActive,
			StatusWord:      status.StatusWord,
			StateVar:        status.StateVar,
			ErrorCode:       status.ErrorCode,
			WarnWord:        status.WarnWord,
			StateMachine:    stateName,
		}
	}

	return nil
}

func stateMachineStateName(state uint8) string {
	switch state {
	case 0:
		return "NotReadyToSwitchOn(0)"
	case 1:
		return "SwitchOnDisabled(1)"
	case 2:
		return "ReadyToSwitchOn(2)"
	case 3:
		return "SetupError(3)"
	case 4:
		return "GeneralError(4)"
	case 5:
		return "HWTests(5)"
	case 6:
		return "ReadyToOperate(6)"
	case 8:
		return "OperationEnabled(8)"
	case 9:
		return "Homing(9)"
	default:
		return fmt.Sprintf("State(%d)", state)
	}
}
