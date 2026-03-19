package protocol_command_tables

import "fmt"

// EntryNotFoundError is returned when a command table entry is not found.
type EntryNotFoundError struct {
	EntryID uint16
}

func (e *EntryNotFoundError) Error() string {
	return fmt.Sprintf("command table entry %d not found", e.EntryID)
}

// EntryAllocationError is returned when entry allocation fails.
type EntryAllocationError struct {
	EntryID uint16
	Size    uint16
	Status  uint8
}

func (e *EntryAllocationError) Error() string {
	return fmt.Sprintf("failed to allocate entry %d (size=%d): status 0x%02X", e.EntryID, e.Size, e.Status)
}

// EntryWriteError is returned when writing entry data fails.
type EntryWriteError struct {
	EntryID uint16
	Offset  int
	Status  uint8
}

func (e *EntryWriteError) Error() string {
	return fmt.Sprintf("failed to write entry %d data at offset %d: status 0x%02X", e.EntryID, e.Offset, e.Status)
}

// EntryReadError is returned when reading entry data fails.
type EntryReadError struct {
	EntryID uint16
	Status  uint8
}

func (e *EntryReadError) Error() string {
	return fmt.Sprintf("failed to read entry %d: status 0x%02X", e.EntryID, e.Status)
}

// MCControlError is returned when motion controller control operations fail.
type MCControlError struct {
	Operation string // "stop" or "start"
	Status    uint8
}

func (e *MCControlError) Error() string {
	return fmt.Sprintf("failed to %s motion controller: status 0x%02X", e.Operation, e.Status)
}
