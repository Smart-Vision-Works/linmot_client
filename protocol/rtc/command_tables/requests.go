package protocol_command_tables

import (
	"fmt"
	"time"

	protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/linmot_client/protocol/rtc"
)

// Compile-time interface checks
var (
	_ protocol_common.Request = (*SaveCommandTableRequest)(nil)
	_ protocol_common.Request = (*StopMotionControllerRequest)(nil)
	_ protocol_common.Request = (*StartMotionControllerRequest)(nil)
	_ protocol_common.Request = (*DeleteAllEntriesRequest)(nil)
	_ protocol_common.Request = (*DeleteEntryRequest)(nil)
	_ protocol_common.Request = (*AllocateEntryRequest)(nil)
	_ protocol_common.Request = (*WriteEntryDataRequest)(nil)
	_ protocol_common.Request = (*GetEntrySizeRequest)(nil)
	_ protocol_common.Request = (*ReadEntryDataRequest)(nil)
	_ protocol_common.Request = (*PresenceMaskRequest)(nil)
)

// SaveCommandTableRequest requests to save the Command Table to Flash.
type SaveCommandTableRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// StopMotionControllerRequest requests to stop the Motion Controller.
type StopMotionControllerRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// StartMotionControllerRequest requests to start the Motion Controller.
type StartMotionControllerRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// DeleteAllEntriesRequest requests to delete all command table entries.
type DeleteAllEntriesRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// DeleteEntryRequest requests to delete a specific command table entry.
type DeleteEntryRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// AllocateEntryRequest requests to allocate a command table entry.
type AllocateEntryRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// WriteEntryDataRequest requests to write data to a command table entry.
type WriteEntryDataRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// GetEntrySizeRequest requests to get the size of a command table entry.
type GetEntrySizeRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// ReadEntryDataRequest requests to read data from a command table entry.
type ReadEntryDataRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// PresenceMaskRequest requests to get a presence mask.
type PresenceMaskRequest struct {
	protocol_rtc.RTCSetParamRequest
}

// NewSaveCommandTableRequest creates a new RTC write request to save the Command Table to Flash.
// This command does not require a upid or value (both are 0).
// Reference: LinUDP V2 specification, Command Code 0x80 (Command Table: Save to Flash)
func NewSaveCommandTableRequest() *SaveCommandTableRequest {
	return &SaveCommandTableRequest{
		RTCSetParamRequest: *newCommandTableRequest(CommandCode.SaveCommandTable, 0, 0, 0),
	}
}

// NewStopMotionControllerRequest creates an RTC request to stop the Motion Controller.
func NewStopMotionControllerRequest() *StopMotionControllerRequest {
	return &StopMotionControllerRequest{
		RTCSetParamRequest: *newCommandTableRequest(CommandCode.StopMotionController, 0, 0, 0),
	}
}

// NewStartMotionControllerRequest creates an RTC request to start the Motion Controller.
func NewStartMotionControllerRequest() *StartMotionControllerRequest {
	return &StartMotionControllerRequest{
		RTCSetParamRequest: *newCommandTableRequest(CommandCode.StartMotionController, 0, 0, 0),
	}
}

// OperationTimeout returns the timeout duration for StopMotionController requests.
// StopMC operations may take longer as they need to stop the motion controller safely.
func (*StopMotionControllerRequest) OperationTimeout() time.Duration {
	return 1 * time.Second // 1s per attempt
}

// OperationTimeout returns the timeout duration for StartMotionController requests.
// StartMC operations may take longer as they need to start the motion controller.
func (*StartMotionControllerRequest) OperationTimeout() time.Duration {
	return 1 * time.Second // 1s per attempt
}

// NewDeleteAllEntriesRequest creates an RTC request to delete all command table entries.
func NewDeleteAllEntriesRequest() *DeleteAllEntriesRequest {
	return &DeleteAllEntriesRequest{
		RTCSetParamRequest: *newCommandTableRequest(CommandCode.DeleteAllEntries, 0, 0, 0),
	}
}

// NewDeleteEntryRequest creates an RTC request to delete a specific command table entry.
func NewDeleteEntryRequest(entryID uint16) *DeleteEntryRequest {
	return &DeleteEntryRequest{
		RTCSetParamRequest: *newCommandTableRequest(CommandCode.DeleteEntry, entryID, 0, 0),
	}
}

// NewAllocateEntryRequest creates an RTC request to allocate a command table entry.
// size must be even. Returns an error if size is odd.
// NOTE: For AllocateEntry, the size must be in Word3 (low 16 bits of value), not Word4 (high 16 bits).
// This is different from the standard newCommandTableRequest packing.
func NewAllocateEntryRequest(entryID, size uint16) (*AllocateEntryRequest, error) {
	if size%2 == 1 {
		return nil, fmt.Errorf("size must be even, got: %d", size)
	}
	// For AllocateEntry, size goes in Word3 (low 16 bits), not Word4 (high 16 bits)
	// So we swap w3 and w4: value = (w4 << 16) | w3 = (0 << 16) | size = size
	value := int32(uint32(0)<<16 | uint32(size))
	return &AllocateEntryRequest{
		RTCSetParamRequest: *protocol_rtc.NewRTCSetParamRequest(entryID, value, CommandCode.AllocateEntry),
	}, nil
}

// NewWriteEntryDataRequest creates an RTC request to write data to a command table entry.
// data must be at most 4 bytes long. Returns an error if data is too long.
func NewWriteEntryDataRequest(entryID uint16, data []byte) (*WriteEntryDataRequest, error) {
	if len(data) > 4 {
		return nil, fmt.Errorf("data too long: %d bytes (max 4)", len(data))
	}
	w3, w4 := packBytes(data)
	return &WriteEntryDataRequest{
		RTCSetParamRequest: *newCommandTableRequest(CommandCode.WriteEntryData, entryID, w3, w4),
	}, nil
}

// NewGetEntrySizeRequest creates an RTC request to get the size of a command table entry.
func NewGetEntrySizeRequest(entryID uint16) *GetEntrySizeRequest {
	return &GetEntrySizeRequest{
		RTCSetParamRequest: *newCommandTableRequest(CommandCode.GetEntrySize, entryID, 0, 0),
	}
}

// NewReadEntryDataRequest creates an RTC request to read data from a command table entry.
func NewReadEntryDataRequest(entryID uint16) *ReadEntryDataRequest {
	return &ReadEntryDataRequest{
		RTCSetParamRequest: *newCommandTableRequest(CommandCode.ReadEntryData, entryID, 0, 0),
	}
}

// NewPresenceMaskRequest creates an RTC request to get a presence mask.
// maskIndex must be 0-7, corresponding to commands 0x87-0x8E.
func NewPresenceMaskRequest(maskIndex uint8) (*PresenceMaskRequest, error) {
	if maskIndex > 7 {
		return nil, fmt.Errorf("maskIndex must be 0-7, got %d", maskIndex)
	}
	cmdCode := CommandCode.PresenceMask0 + maskIndex
	return &PresenceMaskRequest{
		RTCSetParamRequest: *newCommandTableRequest(cmdCode, 0, 0, 0),
	}, nil
}

// newCommandTableRequest creates an RTC request for command table operations.
// This is exported for use by the command_table package.
// w2, w3, w4 are the CT-specific parameters packed into upid and value fields.
func newCommandTableRequest(cmdCode uint8, w2, w3, w4 uint16) *protocol_rtc.RTCSetParamRequest {
	// Pack w3 and w4 into value field: value = (w3 << 16) | w4
	value := int32(uint32(w3)<<16 | uint32(w4))
	return protocol_rtc.NewRTCSetParamRequest(w2, value, cmdCode)
}

// packBytes converts up to 4 bytes into w3 and w4 uint16 values.
// Wire format must match C# library: bytes go on wire as [data[0], data[1], data[2], data[3]].
// With little-endian encoding of value=(w3<<16)|w4:
//   - bytes[12] = w4 low byte = data[0]
//   - bytes[13] = w4 high byte = data[1]
//   - bytes[14] = w3 low byte = data[2]
//   - bytes[15] = w3 high byte = data[3]
func packBytes(data []byte) (w3, w4 uint16) {
	var b0, b1, b2, b3 byte
	if len(data) > 0 {
		b0 = data[0]
	}
	if len(data) > 1 {
		b1 = data[1]
	}
	if len(data) > 2 {
		b2 = data[2]
	}
	if len(data) > 3 {
		b3 = data[3]
	}
	// Pack for little-endian wire order: [b0, b1, b2, b3]
	w4 = uint16(b1)<<8 | uint16(b0) // w4 low=b0, w4 high=b1
	w3 = uint16(b3)<<8 | uint16(b2) // w3 low=b2, w3 high=b3
	return w3, w4
}
