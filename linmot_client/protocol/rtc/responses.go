package protocol_rtc

import (
	"encoding/binary"
	"fmt"

	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
)

// Compile-time interface checks
var (
	_ protocol_common.PacketWritable = (*RTCGetParamResponse)(nil)
	_ protocol_common.PacketWritable = (*RTCSetParamResponse)(nil)
)

// RTCResponse is implemented by any RTC response that provides counter and status metadata.
// This interface allows RequestManager to route responses without depending on concrete types,
// enabling typed responses to work seamlessly.
type RTCResponse interface {
	protocol_common.Response
	RTCCounter() uint8
	RTCStatus() uint8
}

// MissingResponseRegistryError indicates no response registry was registered for the given UPID/cmd.
type MissingResponseRegistryError struct {
	Op      string
	UPID    uint16
	CmdCode uint8
	Counter uint8
}

func (e *MissingResponseRegistryError) Error() string {
	return fmt.Sprintf("%s: missing response registry for parameter 0x%04X (cmd 0x%02X counter %d)", e.Op, e.UPID, e.CmdCode, e.Counter)
}

// RTCGetParamResponse indicates an RTC parameter get completed.
type RTCGetParamResponse struct {
	status     *protocol_common.Status
	value      int32
	upid       uint16
	rtcCounter uint8 // Echo from request
	rtcStatus  uint8 // RTC command status code (0x00=OK, 0x02=Busy, 0xC0-0xD4=Errors)
	cmdCode    uint8 // Command code for serialization (0 = not set, will default to ReadRAM)
}

// NewRTCGetParamResponse creates a new RTCGetParamResponse with the given values.
// This is primarily for testing and mock implementations.
func NewRTCGetParamResponse(status *protocol_common.Status, value int32, upid uint16, rtcCounter uint8, rtcStatus uint8) *RTCGetParamResponse {
	return &RTCGetParamResponse{
		status:     status,
		value:      value,
		upid:       upid,
		rtcCounter: rtcCounter,
		rtcStatus:  rtcStatus,
		cmdCode:    0, // Default, will use ReadRAM
	}
}

// NewRTCGetParamResponseWithCmdCode creates a new RTCGetParamResponse with the given values including command code.
func NewRTCGetParamResponseWithCmdCode(status *protocol_common.Status, value int32, upid uint16, rtcCounter uint8, rtcStatus uint8, cmdCode uint8) *RTCGetParamResponse {
	return &RTCGetParamResponse{
		status:     status,
		value:      value,
		upid:       upid,
		rtcCounter: rtcCounter,
		rtcStatus:  rtcStatus,
		cmdCode:    cmdCode,
	}
}

func (response *RTCGetParamResponse) Status() *protocol_common.Status { return response.status }

// Value returns the raw protocol value (int32).
// For physical units, use VelocityMS() or AccelerationMS2() instead.
func (response *RTCGetParamResponse) Value() int32 {
	return response.value
}

// VelocityMS converts the value field from protocol units to meters per second.
// Use this when reading velocity parameters (e.g., Parameter.Speed1, Parameter.Speed2).
func (response *RTCGetParamResponse) VelocityMS() float64 {
	return protocol_common.FromProtocolVelocity(uint32(response.value))
}

// AccelerationMS2 converts the value field from protocol units to meters per second squared.
// Use this when reading acceleration or deceleration parameters (e.g., Parameter.Acceleration1, Parameter.Deceleration1).
func (response *RTCGetParamResponse) AccelerationMS2() float64 {
	return protocol_common.FromProtocolAcceleration(uint32(response.value))
}

func (response *RTCGetParamResponse) RTCCounter() uint8 {
	return response.rtcCounter
}

func (response *RTCGetParamResponse) RTCStatus() uint8 {
	return response.rtcStatus
}

// WritePacket serializes an RTCGetParamResponse to a UDP packet (48 bytes).
// The cmdCode is determined from the stored cmdCode field, or defaults to ReadRAM if not set.
func (response *RTCGetParamResponse) WritePacket() ([]byte, error) {
	cmdCode := response.cmdCode
	if cmdCode == 0 {
		// Default to ReadRAM for standard parameter reads
		cmdCode = CommandCode.ReadRAM
	}
	return writeRTCResponse(response.status, response.Value(), response.rtcCounter, response.upid, cmdCode, response.rtcStatus), nil
}

// RTCSetParamResponse indicates an RTC parameter set completed.
type RTCSetParamResponse struct {
	status     *protocol_common.Status
	upid       uint16
	value      int32 // Echo of written value
	rtcCounter uint8 // Echo from request
	rtcStatus  uint8 // RTC command status code (0x00=OK, 0x02=Busy, 0xC0-0xD4=Errors)
	cmdCode    uint8 // Command code for serialization (0 = not set, will default to WriteRAM)
}

// NewRTCSetParamResponse creates a new RTCSetParamResponse with the given values.
// This is primarily for testing and mock implementations.
func NewRTCSetParamResponse(status *protocol_common.Status, value int32, upid uint16, rtcCounter uint8, rtcStatus uint8) *RTCSetParamResponse {
	return &RTCSetParamResponse{
		status:     status,
		value:      value,
		upid:       upid,
		rtcCounter: rtcCounter,
		rtcStatus:  rtcStatus,
		cmdCode:    0, // Default, will use WriteRAM
	}
}

// NewRTCSetParamResponseWithCmdCode creates a new RTCSetParamResponse with the given values including command code.
// This is used for CT commands that need to encode their command code in the response.
func NewRTCSetParamResponseWithCmdCode(status *protocol_common.Status, value int32, upid uint16, rtcCounter uint8, rtcStatus uint8, cmdCode uint8) *RTCSetParamResponse {
	return &RTCSetParamResponse{
		status:     status,
		value:      value,
		upid:       upid,
		rtcCounter: rtcCounter,
		rtcStatus:  rtcStatus,
		cmdCode:    cmdCode,
	}
}

func (response *RTCSetParamResponse) Status() *protocol_common.Status { return response.status }

func (response *RTCSetParamResponse) UPID() uint16 { return response.upid }

func (response *RTCSetParamResponse) Value() int32 {
	return response.value
}

func (response *RTCSetParamResponse) RTCStatus() uint8 {
	return response.rtcStatus
}

func (response *RTCSetParamResponse) RTCCounter() uint8 {
	return response.rtcCounter
}

// WritePacket serializes an RTCSetParamResponse to a UDP packet (48 bytes).
// The cmdCode is determined from the stored cmdCode field, or defaults to WriteRAM if not set.
// For CT commands, rtcStatus contains the actual command status.
func (response *RTCSetParamResponse) WritePacket() ([]byte, error) {
	cmdCode := response.cmdCode
	if cmdCode == 0 {
		// Default to WriteRAM for standard parameter writes
		cmdCode = CommandCode.WriteRAM
	}
	return writeRTCResponse(response.status, response.Value(), response.rtcCounter, response.upid, cmdCode, response.rtcStatus), nil
}

// writeRTCResponse builds an RTC response packet with explicit RTC status code.
// For CT commands, rtcStatus contains the actual status; for standard commands, it's typically 0x00.
func writeRTCResponse(status *protocol_common.Status, value int32, rtcCounter uint8, upid uint16, cmdCode uint8, rtcStatus uint8) []byte {
	// Size calculation:
	// Header (8) + Standard Status (18) + Monitoring (16) + RTC Data (8) = 50 bytes
	// Note: Standard Status is bits 0-6. Monitoring is bit 7. RTC is bit 8.
	// ResponseFlags.RTCReply includes all these bits (0x1FF).
	packet := make([]byte, 50) // Full RTC reply size

	// Header - request flags must have bit 2 set (RTC command) so response is routed correctly
	binary.LittleEndian.PutUint32(packet[0:4], protocol_common.RequestFlags.RTCCommand)
	binary.LittleEndian.PutUint32(packet[4:8], protocol_common.ResponseFlags.RTCReply)

	// status data (20 bytes, offsets 8-27)
	binary.LittleEndian.PutUint16(packet[8:10], status.StatusWord)
	binary.LittleEndian.PutUint16(packet[10:12], status.StateVar)
	binary.LittleEndian.PutUint32(packet[12:16], uint32(status.ActualPosition))
	binary.LittleEndian.PutUint32(packet[16:20], uint32(status.DemandPosition))
	binary.LittleEndian.PutUint16(packet[20:22], uint16(status.Current))
	binary.LittleEndian.PutUint16(packet[22:24], status.WarnWord)
	binary.LittleEndian.PutUint16(packet[24:26], status.ErrorCode)

	// Padding (2 bytes, offsets 26-27)
	packet[26] = 0
	packet[27] = 0

	rtcOffset, err := protocol_common.ResponseBlockOffset(protocol_common.ResponseFlags.RTCReply, protocol_common.RespBitRTCReplyData)
	if err != nil {
		rtcOffset = len(packet) - protocol_common.BlockSizes[protocol_common.RespBitRTCReplyData]
	}

	// Byte 0: low nibble (bits 3-0) = Command Count, high nibble (bits 7-4) = 0 (reserved)
	packet[rtcOffset+RTCDataOffsetCounter] = rtcCounter & RTCCounterMask // Only low nibble
	// Byte 1: Encoding depends on command type
	// - Standard parameter access (0x10-0x17): cmdCode
	// - CT commands (0x35-0x36, 0x80-0x8E): rtcStatus
	// - Special commands (0x20-0x23, 0x30-0x34, 0x40-0x74): rtcStatus
	if IsCTCommand(cmdCode) || IsSpecialCommand(cmdCode) {
		// CT and special commands: encode status in byte29
		packet[rtcOffset+RTCDataOffsetCmdOrStatus] = rtcStatus
	} else {
		// Standard parameter access: encode cmdCode for read/write determination
		packet[rtcOffset+RTCDataOffsetCmdOrStatus] = cmdCode
	}
	binary.LittleEndian.PutUint16(packet[rtcOffset+RTCDataOffsetUPID:rtcOffset+RTCDataOffsetUPID+2], upid)
	binary.LittleEndian.PutUint32(packet[rtcOffset+RTCDataOffsetValue:rtcOffset+RTCDataOffsetValue+4], uint32(value))

	// Bytes between status data and RTC data (28-39) are zero-filled by make()

	return packet
}

// ExtractRTCCommandCount extracts just the command count from an RTC response packet.
// This is used to look up the pending request before full parsing.
func ExtractRTCCommandCount(data []byte) (uint8, error) {
	if len(data) < protocol_common.PacketHeaderSize {
		return 0, protocol_common.NewPacketTooShortError("ExtractRTCCommandCount", len(data), protocol_common.PacketHeaderSize, data)
	}

	_, repBits, err := protocol_common.ReadPacketHeader(data)
	if err != nil {
		return 0, err
	}
	if repBits&protocol_common.RespBitRTCReplyData == 0 {
		return 0, protocol_common.NewInvalidFlagsError("ExtractRTCCommandCount", repBits, protocol_common.RespBitRTCReplyData, data)
	}

	rtcOffset, err := protocol_common.ResponseBlockOffset(repBits, protocol_common.RespBitRTCReplyData)
	if err != nil {
		return 0, err
	}
	rtcSize := protocol_common.BlockSizes[protocol_common.RespBitRTCReplyData]
	if rtcOffset+rtcSize > len(data) {
		return 0, protocol_common.NewPacketTooShortError("ExtractRTCCommandCount (RTC data)", len(data), rtcOffset+rtcSize, data)
	}

	return data[rtcOffset+RTCDataOffsetCounter] & RTCCounterMask, nil
}

// ReadRTCResponse parses an RTC response packet (variable length, server → client).
// The RTC reply data offset is computed from response-definition bits (bit order).
//
// cmdCode is the original command code from the request (0 if unknown).
// For special commands (0x20-0x74, 0x80-0x8E), cmdCode must be provided for registry lookup
// because byte29 contains the status, not the command code.
//
// requestUPID is the UPID from the original request (0 if unknown).
// Following the C# library pattern, we trust the counter match and use the request UPID
// for response type determination, ignoring the response UPID which may be stale or incorrect.
func ReadRTCResponse(data []byte, cmdCode uint8, requestUPID uint16) (any, error) {
	minSize := protocol_common.PacketHeaderSize + protocol_common.BlockSizes[protocol_common.RespBitRTCReplyData]
	if len(data) < minSize {
		return nil, protocol_common.NewPacketTooShortError("ReadRTCResponse", len(data), minSize, data)
	}

	// Parse status dynamically
	status, err := protocol_common.ReadStatus(data, false)
	if err != nil {
		return nil, err
	}

	rtcOffset, err := protocol_common.ResponseBlockOffsetFromPacket(data, protocol_common.RespBitRTCReplyData)
	if err != nil {
		return nil, err
	}
	rtcSize := protocol_common.BlockSizes[protocol_common.RespBitRTCReplyData]
	if rtcOffset+rtcSize > len(data) {
		return nil, protocol_common.NewPacketTooShortError("ReadRTCResponse (RTC data)", len(data), rtcOffset+rtcSize, data)
	}

	// Parse RTC data at the end of the packet
	// Byte 0: low nibble (bits 3-0) = Command Count, high nibble (bits 7-4) = reserved
	rtcCounter := data[rtcOffset+RTCDataOffsetCounter] & RTCCounterMask
	// Byte 1: For standard commands, this is the cmdCode; for CT commands, this is the RTC status
	byte29 := data[rtcOffset+RTCDataOffsetCmdOrStatus]
	// Byte 2-3: UPID from response (may be stale or incorrect - we use requestUPID for registry lookup)
	responseUPID := binary.LittleEndian.Uint16(data[rtcOffset+RTCDataOffsetUPID : rtcOffset+RTCDataOffsetUPID+2])
	value := int32(binary.LittleEndian.Uint32(data[rtcOffset+RTCDataOffsetValue : rtcOffset+RTCDataOffsetValue+4]))

	// Use request UPID for registry lookup (C# library pattern).
	// If requestUPID is 0 (unknown), fall back to response UPID for backward compatibility.
	upidForRegistry := requestUPID
	if upidForRegistry == 0 {
		upidForRegistry = responseUPID
	}

	// Check if byte29 is a read command code
	// OR if byte29 is 0x00 (success status) and cmdCode indicates a read command
	// (Some real drives return status in byte29 instead of echoing the command code)
	isReadCommand := byte29 == CommandCode.ReadROM || byte29 == CommandCode.ReadRAM ||
		byte29 == CommandCode.GetMinValue || byte29 == CommandCode.GetMaxValue || byte29 == CommandCode.GetDefaultValue
	if !isReadCommand && byte29 == 0x00 && (cmdCode == CommandCode.ReadRAM || cmdCode == CommandCode.ReadROM ||
		cmdCode == CommandCode.GetMinValue || cmdCode == CommandCode.GetMaxValue || cmdCode == CommandCode.GetDefaultValue) {
		// Real drive returned status code instead of command code - use provided cmdCode
		isReadCommand = true
	}

	if isReadCommand {
		// Extended parameter read commands - don't use UPID registry
		readCmdCode := byte29
		if byte29 == 0x00 && (cmdCode == CommandCode.ReadRAM || cmdCode == CommandCode.ReadROM ||
			cmdCode == CommandCode.GetMinValue || cmdCode == CommandCode.GetMaxValue || cmdCode == CommandCode.GetDefaultValue) {
			readCmdCode = cmdCode
		}
		if readCmdCode == CommandCode.GetMinValue || readCmdCode == CommandCode.GetMaxValue || readCmdCode == CommandCode.GetDefaultValue {
			return NewRTCGetParamResponseWithCmdCode(
				status,
				value,
				upidForRegistry,
				rtcCounter,
				0x00,
				readCmdCode,
			), nil
		}

		// Standard parameter reads use the trusted request UPID for both lookup and
		// response identity because some drives echo a stale UPID in the reply.
		param := protocol_common.ParameterID(upidForRegistry)
		if registry, ok := LookupResponseRegistryCmdAware(readCmdCode, param); ok {
			return registry(status, value, upidForRegistry, rtcCounter, 0x00, readCmdCode), nil
		}
		// Fallback: no typed registry entry exists for this UPID/cmd pair.
		// Return a generic RTCGetParamResponse instead of failing parse.
		return NewRTCGetParamResponseWithCmdCode(
			status,
			value,
			upidForRegistry,
			rtcCounter,
			0x00,
			readCmdCode,
		), nil
	}

	// Check if byte29 is a standard write command code (0x12, 0x13)
	// OR if byte29 is 0x00 (success status) and cmdCode indicates a write command
	// (Some real drives return status in byte29 instead of echoing the command code)
	isWriteCommand := byte29 == CommandCode.WriteRAM || byte29 == CommandCode.WriteROM
	if !isWriteCommand && byte29 == 0x00 && (cmdCode == CommandCode.WriteRAM || cmdCode == CommandCode.WriteROM) {
		// Real drive returned status code instead of command code - use provided cmdCode
		isWriteCommand = true
	}

	if isWriteCommand {
		// Standard parameter writes use the trusted request UPID for both lookup and
		// response identity because some drives echo a stale UPID in the reply.
		// Use cmdCode if byte29 is status (0x00), otherwise use byte29
		writeCmdCode := byte29
		if byte29 == 0x00 && (cmdCode == CommandCode.WriteRAM || cmdCode == CommandCode.WriteROM) {
			writeCmdCode = cmdCode
		}
		param := protocol_common.ParameterID(upidForRegistry)
		if registry, ok := LookupResponseRegistryCmdAware(writeCmdCode, param); ok {
			return registry(status, value, upidForRegistry, rtcCounter, 0x00, writeCmdCode), nil
		}
		// Fallback: no typed registry entry exists for this UPID/cmd pair.
		// Return a generic RTCSetParamResponse instead of failing parse.
		return NewRTCSetParamResponseWithCmdCode(
			status,
			value,
			upidForRegistry,
			rtcCounter,
			0x00,
			writeCmdCode,
		), nil
	}

	// WriteRAMAndROM (0x14) - extended command, don't use UPID registry
	// Handle both cases: drive echoes command code (0x14) or returns status (0x00)
	if byte29 == CommandCode.WriteRAMAndROM ||
		(byte29 == 0x00 && cmdCode == CommandCode.WriteRAMAndROM) {
		writeCmdCode := byte29
		if byte29 == 0x00 {
			writeCmdCode = cmdCode
		}
		return NewRTCSetParamResponseWithCmdCode(
			status,
			value,
			upidForRegistry,
			rtcCounter,
			0x00,
			writeCmdCode,
		), nil
	}

	// Otherwise, this is a CT or special command where byte29 contains the status
	// For these commands, byte29 is the status (0x00=OK, 0x02=Busy, 0xC0-0xD4=Errors)
	// We need the original cmdCode (from request) for registry lookup
	// Use request UPID for registry lookup (special commands don't have UPID registries, but we keep consistency)
	param := protocol_common.ParameterID(upidForRegistry)

	// For special/CT commands registered by cmdCode, use the provided cmdCode
	if cmdCode != 0 && (IsSpecialCommand(cmdCode) || IsCTCommand(cmdCode)) {
		if registry, ok := LookupResponseRegistryByCmd(cmdCode); ok {
			// byte29 is the status, cmdCode is the original command from request
			return registry(status, value, responseUPID, rtcCounter, byte29, cmdCode), nil
		}
	}

	// Fallback: try lookup treating byte29 as status (for unregistered commands)
	if registry, ok := LookupResponseRegistryCmdAware(byte29, param); ok {
		return registry(status, value, responseUPID, rtcCounter, byte29, byte29), nil
	}

	// Couldn't find registry - return error with the original cmdCode if available
	errCmdCode := byte29
	if cmdCode != 0 {
		errCmdCode = cmdCode
	}
	return nil, &MissingResponseRegistryError{
		Op:      "ReadRTCResponse",
		UPID:    upidForRegistry,
		CmdCode: errCmdCode,
		Counter: rtcCounter,
	}
}
