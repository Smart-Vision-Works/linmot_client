package protocol_command_tables

// CommandCode groups all Command Table command code constants.
// Usage: CommandCode.StopMotionController, CommandCode.AllocateEntry
// Reference: LinUDP V2 protocol specification
var CommandCode = struct {
	// StopMotionController stops the Motion Controller software on the drive (RTC 0x35).
	StopMotionController uint8

	// StartMotionController restarts the Motion Controller software on the drive (RTC 0x36).
	StartMotionController uint8

	// SaveCommandTable saves the command table to Flash.
	// Future: not yet implemented
	SaveCommandTable uint8

	// DeleteAllEntries clears all Command Table entries from RAM (RTC 0x81).
	DeleteAllEntries uint8

	// DeleteEntry deletes a single Command Table entry by ID from RAM (RTC 0x82).
	DeleteEntry uint8

	// AllocateEntry creates/allocates a CT entry with the given ID and size (RTC 0x83).
	AllocateEntry uint8

	// WriteEntryData appends up to 4 bytes per call to the CT entry data buffer (RTC 0x84).
	WriteEntryData uint8

	// GetEntrySize returns the block size of a CT entry (RTC 0x85).
	GetEntrySize uint8

	// ReadEntryData reads up to 4 bytes from a CT entry (RTC 0x86).
	ReadEntryData uint8

	// PresenceMask0-7 return the eight 32-bit masks of defined/free entries (RTC 0x87-0x8E).
	PresenceMask0 uint8
	PresenceMask1 uint8
	PresenceMask2 uint8
	PresenceMask3 uint8
	PresenceMask4 uint8
	PresenceMask5 uint8
	PresenceMask6 uint8
	PresenceMask7 uint8
}{
	StopMotionController:  0x35,
	StartMotionController: 0x36,
	SaveCommandTable:      0x80,
	DeleteAllEntries:      0x81,
	DeleteEntry:           0x82,
	AllocateEntry:         0x83,
	WriteEntryData:        0x84,
	GetEntrySize:          0x85,
	ReadEntryData:         0x86,
	PresenceMask0:         0x87,
	PresenceMask1:         0x88,
	PresenceMask2:         0x89,
	PresenceMask3:         0x8A,
	PresenceMask4:         0x8B,
	PresenceMask5:         0x8C,
	PresenceMask6:         0x8D,
	PresenceMask7:         0x8E,
}

// IsCTCommand checks if a command code is a Command Table operation.
// CT commands are in two ranges:
//   - 0x35-0x36: Motion Controller control (StopMC, StartMC)
//   - 0x80-0x8E: Command Table operations (SaveCommandTable, DeleteAllEntries, etc.)
func IsCTCommand(cmdCode uint8) bool {
	return (cmdCode >= 0x35 && cmdCode <= 0x36) ||
		(cmdCode >= 0x80 && cmdCode <= 0x8E)
}

const (
	// PresenceMaskCount is the number of presence masks (8 masks cover 256 entry IDs).
	// Each mask covers 32 entry IDs (mask 0 = IDs 0-31, mask 1 = IDs 32-63, etc.).
	// Reference: LinUDP V2 protocol specification
	PresenceMaskCount = 8

	// EntryMinSize is the minimum CT entry size in bytes (64 bytes for A701h format).
	// Reference: Command table entry header format (A701h version identifier)
	EntryMinSize = 64

	// EntryHeaderSize is the size of the CT entry header in bytes (first 2 bytes are version).
	// Reference: Command table entry header format
	EntryHeaderSize = 2
)
