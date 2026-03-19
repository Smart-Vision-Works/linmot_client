package protocol_rtc

import (
	protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"
)

// ValidateRTCRead validates an RTC read operation.
// It checks that the upid is within valid range and the command code is valid.
func ValidateRTCRead(upid uint16, cmdCode uint8) error {
	// upid 0xFFFF is reserved and should not be used
	if upid == 0xFFFF {
		return protocol_common.NewInvalidUPIDError("ValidateRTCRead", upid, "reserved upid")
	}

	// Validate command code - allow all read-style commands (0x10-0x17)
	validCodes := []uint8{
		CommandCode.ReadROM,
		CommandCode.ReadRAM,
		CommandCode.GetMinValue,
		CommandCode.GetMaxValue,
		CommandCode.GetDefaultValue,
	}
	isValid := cmdCode == CommandCode.ReadROM ||
		cmdCode == CommandCode.ReadRAM ||
		cmdCode == CommandCode.GetMinValue ||
		cmdCode == CommandCode.GetMaxValue ||
		cmdCode == CommandCode.GetDefaultValue
	if !isValid {
		return protocol_common.NewInvalidCmdCodeError("ValidateRTCRead", cmdCode, validCodes)
	}

	return nil
}

// ValidateRTCWrite validates an RTC write operation.
// It checks that the upid is valid for writing and the command code is valid.
// For special operations (CT, curves, errors, etc.), validation is relaxed as they use custom data layouts.
func ValidateRTCWrite(upid uint16, value int32, cmdCode uint8) error {
	// For special operations, skip strict UPID validation
	if IsSpecialCommand(cmdCode) {
		// Only validate that upid is not reserved (0xFFFF is sometimes used for special commands)
		// Allow 0x0000 for many special commands
		return nil
	}

	// Standard and extended parameter write operations (0x12, 0x13, 0x14)
	// Validate upid is not reserved
	if upid == 0xFFFF {
		return protocol_common.NewInvalidUPIDError("ValidateRTCWrite", upid, "reserved upid")
	}

	// Validate command code
	validCodes := []uint8{CommandCode.WriteROM, CommandCode.WriteRAM, CommandCode.WriteRAMAndROM}
	isValid := cmdCode == CommandCode.WriteROM || cmdCode == CommandCode.WriteRAM || cmdCode == CommandCode.WriteRAMAndROM
	if !isValid {
		return protocol_common.NewInvalidCmdCodeError("ValidateRTCWrite", cmdCode, validCodes)
	}

	// Only validate UPID for standard write operations
	if cmdCode == CommandCode.WriteRAM || cmdCode == CommandCode.WriteROM {
		if !protocol_common.IsValidPUID(protocol_common.ParameterID(upid)) {
			return protocol_common.NewInvalidUPIDError("ValidateRTCWrite", upid, "not a valid PUID")
		}
	}

	return nil
}

// IsSpecialCommand returns true if the command code is for special operations
// (CT, UPID lists, curves, error logs, operations, etc.) that encode STATUS in byte29 instead of cmdCode.
// Extended parameter access (0x14-0x17) are NOT special - they encode cmdCode like standard operations.
func IsSpecialCommand(cmdCode uint8) bool {
	return cmdCode == CommandCode.NoOperation ||
		// UPID List operations (0x20-0x23)
		cmdCode == CommandCode.StartGettingUPIDList ||
		cmdCode == CommandCode.GetNextUPIDListItem ||
		cmdCode == CommandCode.StartGettingModifiedUPIDList ||
		cmdCode == CommandCode.GetNextModifiedUPIDListItem ||
		// Drive operations (0x30-0x36)
		cmdCode == CommandCode.RestartDrive ||
		cmdCode == CommandCode.SetOSROMToDefault ||
		cmdCode == CommandCode.SetMCROMToDefault ||
		cmdCode == CommandCode.SetInterfaceROMToDefault ||
		cmdCode == CommandCode.SetApplicationROMToDefault ||
		IsCTCommand(cmdCode) ||
		// Curve operations (0x40-0x42, 0x50-0x62)
		cmdCode == CommandCode.SaveAllCurvesToFlash ||
		cmdCode == CommandCode.DeleteAllCurves ||
		cmdCode >= 0x50 && cmdCode <= 0x62 ||
		// Error log operations (0x70-0x74)
		cmdCode >= 0x70 && cmdCode <= 0x74
}
