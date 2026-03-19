package protocol_rtc

// ============================================================================
// RTC Packet Constants
// ============================================================================

const (
	// PacketSize is the size of an RTC packet (16 bytes: header 8 + RTC data 8).
	PacketSize = 16

	// MinReplySize is the minimum size of an RTC reply packet (36 bytes: header 8 + status data 20 + RTC reply 8).
	MinReplySize = 36

	// CounterMax is the maximum value for the RTC counter (wraps at 14, range 1-14).
	// Reference: C# library behavior (matches linudp.cs implementation)
	CounterMax = 14

	// RTCDataFieldOffsets define the byte offsets within the 8-byte RTC data block.
	// The RTC data block is located at the end of the packet (last 8 bytes).
	RTCDataOffsetCounter     = 0 // Byte 0: rtcCounter (low nibble, bits 3-0)
	RTCDataOffsetCmdOrStatus = 1 // Byte 1: cmdCode (standard) or rtcStatus (CT/special)
	RTCDataOffsetUPID        = 2 // Bytes 2-3: upid (uint16)
	RTCDataOffsetValue       = 4 // Bytes 4-7: value (uint32)

	// Bit masks for RTC data fields
	RTCCounterMask = 0x0F // Low nibble mask for rtcCounter (bits 3-0)
)

// ============================================================================
// RTC Command Codes
// ============================================================================
// CommandCode groups all RTC command code constants.
// Usage: CommandCode.ReadROM, CommandCode.WriteRAM
// Reference: LinUDP V2 protocol specification

var CommandCode = struct {
	// NoOperation does nothing (can be sent in any state).
	NoOperation uint8

	// Parameter Access (0x10-0x17)
	ReadROM         uint8 // Read ROM parameter by UPID
	ReadRAM         uint8 // Read RAM parameter by UPID
	WriteROM        uint8 // Write ROM parameter by UPID
	WriteRAM        uint8 // Write RAM parameter by UPID
	WriteRAMAndROM  uint8 // Write both RAM and ROM parameter by UPID
	GetMinValue     uint8 // Get minimal value of parameter by UPID
	GetMaxValue     uint8 // Get maximal value of parameter by UPID
	GetDefaultValue uint8 // Get default value of parameter by UPID

	// UPID List Operations (0x20-0x23)
	StartGettingUPIDList         uint8 // Start UPID list iteration
	GetNextUPIDListItem          uint8 // Get next UPID in list
	StartGettingModifiedUPIDList uint8 // Start modified UPID list iteration
	GetNextModifiedUPIDListItem  uint8 // Get next modified UPID

	// Drive Control Operations (0x30-0x36)
	RestartDrive               uint8 // Restart the drive
	SetOSROMToDefault          uint8 // Reset OS SW parameters to default
	SetMCROMToDefault          uint8 // Reset MC SW parameters to default
	SetInterfaceROMToDefault   uint8 // Reset Interface SW parameters to default
	SetApplicationROMToDefault uint8 // Reset Application SW parameters to default
	StopMC                     uint8 // Stop MC and Application Software
	StartMC                    uint8 // Start MC and Application Software

	// Curve Service (0x40-0x62)
	SaveAllCurvesToFlash uint8 // Save all curves from RAM to Flash
	DeleteAllCurves      uint8 // Delete all curves from RAM
	StartAddingCurve     uint8 // Start adding curve to RAM
	AddCurveInfoBlock    uint8 // Add curve info block
	AddCurveData         uint8 // Add curve data
	StartModifyingCurve  uint8 // Start modifying existing curve
	ModifyCurveInfoBlock uint8 // Modify curve info block
	ModifyCurveData      uint8 // Modify curve data
	StartGettingCurve    uint8 // Start getting curve from RAM
	GetCurveInfoBlock    uint8 // Get curve info block
	GetCurveData         uint8 // Get curve data

	// Error Log (0x70-0x74)
	GetErrorLogEntryCounter   uint8 // Get error log entry counter
	GetErrorLogEntryCode      uint8 // Get error code for specific entry
	GetErrorLogEntryTimeLow   uint8 // Get error log entry time (low 32 bits)
	GetErrorLogEntryTimeHigh  uint8 // Get error log entry time (high 32 bits)
	GetErrorCodeTextStringlet uint8 // Get error code text stringlet
}{
	NoOperation: 0x00,

	// Parameter Access
	ReadROM:         0x10,
	ReadRAM:         0x11,
	WriteROM:        0x12,
	WriteRAM:        0x13,
	WriteRAMAndROM:  0x14,
	GetMinValue:     0x15,
	GetMaxValue:     0x16,
	GetDefaultValue: 0x17,

	// UPID List
	StartGettingUPIDList:         0x20,
	GetNextUPIDListItem:          0x21,
	StartGettingModifiedUPIDList: 0x22,
	GetNextModifiedUPIDListItem:  0x23,

	// Drive Control
	RestartDrive:               0x30,
	SetOSROMToDefault:          0x31,
	SetMCROMToDefault:          0x32,
	SetInterfaceROMToDefault:   0x33,
	SetApplicationROMToDefault: 0x34,
	StopMC:                     0x35,
	StartMC:                    0x36,

	// Curve Service
	SaveAllCurvesToFlash: 0x40,
	DeleteAllCurves:      0x41,
	StartAddingCurve:     0x50,
	AddCurveInfoBlock:    0x51,
	AddCurveData:         0x52,
	StartModifyingCurve:  0x53,
	ModifyCurveInfoBlock: 0x54,
	ModifyCurveData:      0x55,
	StartGettingCurve:    0x60,
	GetCurveInfoBlock:    0x61,
	GetCurveData:         0x62,

	// Error Log
	GetErrorLogEntryCounter:   0x70,
	GetErrorLogEntryCode:      0x71,
	GetErrorLogEntryTimeLow:   0x72,
	GetErrorLogEntryTimeHigh:  0x73,
	GetErrorCodeTextStringlet: 0x74,
}

// IsCTCommand checks if a command code is a Command Table operation.
// This is a wrapper around command_table.IsCTCommand for use by general protocol code.
func IsCTCommand(cmdCode uint8) bool {
	// Import command_table to avoid circular dependency
	// We inline the logic here to avoid import cycle
	return (cmdCode >= 0x35 && cmdCode <= 0x36) ||
		(cmdCode >= 0x80 && cmdCode <= 0x8E)
}

// ============================================================================
// RTC Status Codes
// ============================================================================
// Status groups RTC command status code constants.
// Usage: Status.OK, Status.Busy, Status.Incomplete
// Reference: LinUDP V2 protocol specification

var Status = struct {
	// OK indicates successful operation (0x00).
	OK uint8

	// Busy indicates the drive is busy (0x02).
	// This is typically a transient state and the operation should be retried.
	Busy uint8

	// Incomplete indicates a progressive operation is incomplete (0x04).
	// Used for multi-part operations like WriteEntryData or ReadEntryData.
	// The caller should continue the operation until status becomes OK.
	Incomplete uint8

	// Error codes (0xC0-0xD4) indicate various error conditions.
	// Reference: LinUDP V2 specification error code ranges
	ErrorRangeStart uint8 // 0xC0
	ErrorRangeEnd   uint8 // 0xD4
}{
	OK:              0x00,
	Busy:            0x02,
	Incomplete:      0x04,
	ErrorRangeStart: 0xC0,
	ErrorRangeEnd:   0xD4,
}

// IsStatusOK checks if an RTC status code indicates success or acceptable intermediate state.
// Returns true for OK (0x00), Busy (0x02), or Incomplete (0x04).
func IsStatusOK(status uint8) bool {
	return status == Status.OK || status == Status.Busy || status == Status.Incomplete
}

// IsStatusError checks if an RTC status code indicates an error.
// Returns true for status codes in the error range (0xC0-0xD4).
func IsStatusError(status uint8) bool {
	return status >= Status.ErrorRangeStart && status <= Status.ErrorRangeEnd
}
