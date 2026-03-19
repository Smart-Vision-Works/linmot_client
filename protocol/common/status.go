package protocol_common

import (
	"encoding/binary"
	"fmt"
	"time"
)

// Compile-time interface checks
var (
	_ PacketWritable = (*StatusRequest)(nil)
	_ PacketWritable = (*StatusResponse)(nil)
	_ PacketWritable = (*MonitoringStatusRequest)(nil)
	_ PacketWritable = (*MonitoringStatusResponse)(nil)
)

// Status represents the current status of a LinMot drive.
// Fields are populated based on which bits are set in the response definition.
//
// Per LinUDP V2 Protocol Specification (Doc: 0185-1108-E_1V9_MA_LinUDP_V2):
//
//	Bit 0: Status Word        (2 bytes)
//	Bit 1: State Var          (2 bytes)
//	Bit 2: Actual Position    (4 bytes)
//	Bit 3: Demand Position    (4 bytes)
//	Bit 4: Current            (2 bytes)
//	Bit 5: Warn Word          (2 bytes)
//	Bit 6: Error Code         (2 bytes)
//	Bit 7: Monitoring Channel (16 bytes = 4x 32-bit values)
//	Bit 8: RTC Reply          (8 bytes) - handled separately in RTC response parsing
type Status struct {
	// Response definition bits that were present in the response
	ResponseBits uint32

	// Bit 0: Status Word (2 bytes) - Drive status flags
	// See "User Manual Motion Control Software" for bit meanings
	StatusWord uint16

	// Bit 1: State Var (2 bytes) - MainState (high byte) + SubState (low byte)
	// Contains all relevant flags and information for clean handshaking
	StateVar uint16

	// Bit 2: Actual Position (4 bytes) - Current motor position
	// Resolution: 0.1 µm (divide by 10000 for mm)
	ActualPosition int32

	// Bit 3: Demand Position (4 bytes) - Target motor position
	// Resolution: 0.1 µm (divide by 10000 for mm)
	DemandPosition int32

	// Bit 4: Current (2 bytes) - Motor current
	// Resolution: 1 mA
	Current int16

	// Bit 5: Warn Word (2 bytes) - Warning flags
	// See "User Manual Motion Control Software" for bit meanings
	WarnWord uint16

	// Bit 6: Error Code (2 bytes) - Error code if drive is in error state
	// See "User Manual Motion Control Software" for error code meanings
	ErrorCode uint16

	// Bit 7: Monitoring Channel (16 bytes = 4x 32-bit values)
	// Values from 4 UPIDs configured in LinUDP Intf parameters:
	//   Channel 1 UPID (0x20A8), Channel 2 UPID (0x20A9),
	//   Channel 3 UPID (0x20AA), Channel 4 UPID (0x20AB)
	MonitoringChannel [4]int32
}

// Response definition bit masks for dynamic parsing
// Based on LinUDP V2 protocol specification (Doc: 0185-1108-E_1V9_MA_LinUDP_V2, November 2019)
//
// OFFICIAL LinUDP V2 Response Definition (Section 4.4.1):
//
//	Bit 0: Status Word        (2 bytes)  - Drive status flags
//	Bit 1: State Var          (2 bytes)  - MainState + SubState
//	Bit 2: Actual Position    (4 bytes)  - Current motor position (0.1 µm resolution)
//	Bit 3: Demand Position    (4 bytes)  - Target motor position (0.1 µm resolution)
//	Bit 4: Current            (2 bytes)  - Motor current (1 mA resolution)
//	Bit 5: Warn Word          (2 bytes)  - Warning flags
//	Bit 6: Error Code         (2 bytes)  - Error code if in error state
//	Bit 7: Monitoring Channel (16 bytes) - 4x 32-bit values from configured UPIDs
//	Bit 8: RTC Reply          (8 bytes)  - Real-Time Configuration response data
//	Bits 9-31: Reserved for future expansions
//
// Example: Response definition 0x1FF means bits 0-8 are all set:
//
//	Data = 2+2+4+4+2+2+2+16+8 = 42 bytes + 8 byte header = 50 bytes total
const (
	RespBitStatusWord        uint32 = 1 << 0 // Bit 0:  2 bytes - Status Word (drive status flags)
	RespBitStateVar          uint32 = 1 << 1 // Bit 1:  2 bytes - State Var (MainState + SubState)
	RespBitActualPosition    uint32 = 1 << 2 // Bit 2:  4 bytes - Actual Position (0.1 µm resolution)
	RespBitDemandPosition    uint32 = 1 << 3 // Bit 3:  4 bytes - Demand Position (0.1 µm resolution)
	RespBitCurrent           uint32 = 1 << 4 // Bit 4:  2 bytes - Current (1 mA resolution)
	RespBitWarnWord          uint32 = 1 << 5 // Bit 5:  2 bytes - Warn Word (warning flags)
	RespBitErrorCode         uint32 = 1 << 6 // Bit 6:  2 bytes - Error Code
	RespBitMonitoringChannel uint32 = 1 << 7 // Bit 7: 16 bytes - Monitoring Channel (4x 32-bit UPID values)
	RespBitRTCReplyData      uint32 = 1 << 8 // Bit 8:  8 bytes - RTC reply data (counter, status, UPID, value)
	// Bits 9-31 reserved for future expansions per LinUDP V2 spec
)

// BlockSizes maps response definition bits to their data sizes in bytes
// Per LinUDP V2 protocol specification Section 4.4.1
var BlockSizes = map[uint32]int{
	RespBitStatusWord:        2,  // Bit 0: Status Word
	RespBitStateVar:          2,  // Bit 1: State Var
	RespBitActualPosition:    4,  // Bit 2: Actual Position
	RespBitDemandPosition:    4,  // Bit 3: Demand Position
	RespBitCurrent:           2,  // Bit 4: Current
	RespBitWarnWord:          2,  // Bit 5: Warn Word
	RespBitErrorCode:         2,  // Bit 6: Error Code
	RespBitMonitoringChannel: 16, // Bit 7: Monitoring Channel (4x 32-bit values = 16 bytes)
	RespBitRTCReplyData:      8,  // Bit 8: RTC reply data
}

// CalculateExpectedDataSize calculates the expected data size based on response definition bits.
func CalculateExpectedDataSize(repBits uint32) int {
	size := 0
	for bit := uint32(0); bit < 32; bit++ {
		mask := uint32(1) << bit
		if repBits&mask != 0 {
			if blockSize, ok := BlockSizes[mask]; ok {
				size += blockSize
			}
		}
	}
	return size
}

func NewErrorStatus(errorCode uint16) *Status {
	return &Status{ErrorCode: errorCode}
}

// ReadStatusDynamic parses status data dynamically based on the response definition bits.
// This handles variable-length responses from the drive correctly.
//
// Per LinUDP V2 Protocol Specification (Doc: 0185-1108-E_1V9_MA_LinUDP_V2):
// Response data order follows the bit order in the response definition:
//
//	Bit 0: Status Word        (2 bytes)
//	Bit 1: State Var          (2 bytes)
//	Bit 2: Actual Position    (4 bytes)
//	Bit 3: Demand Position    (4 bytes)
//	Bit 4: Current            (2 bytes)
//	Bit 5: Warn Word          (2 bytes)
//	Bit 6: Error Code         (2 bytes)
//	Bit 7: Monitoring Channel (16 bytes)
//	Bit 8: RTC Reply          (8 bytes) - handled separately in RTC response parsing
func ReadStatusDynamic(data []byte) (*Status, error) {
	if len(data) < PacketHeaderSize {
		return nil, NewPacketTooShortError("ReadStatusDynamic", len(data), PacketHeaderSize, data)
	}

	// Read response definition from header (bytes 4-7, little-endian)
	repBits := binary.LittleEndian.Uint32(data[4:8])

	status := &Status{
		ResponseBits: repBits,
	}

	// Start parsing after 8-byte header
	offset := PacketHeaderSize

	// Parse each block in order based on which bits are set
	// Order MUST match the bit order per LinUDP V2 specification

	// Bit 0: Status Word (2 bytes) - Drive status flags
	if repBits&RespBitStatusWord != 0 {
		if err := ensureFieldBytes(data, offset, 2, "StatusWord"); err != nil {
			return status, err
		}
		status.StatusWord = binary.LittleEndian.Uint16(data[offset : offset+2])
		offset += 2
	}

	// Bit 1: State Var (2 bytes) - MainState + SubState
	if repBits&RespBitStateVar != 0 {
		if err := ensureFieldBytes(data, offset, 2, "StateVar"); err != nil {
			return status, err
		}
		status.StateVar = binary.LittleEndian.Uint16(data[offset : offset+2])
		offset += 2
	}

	// Bit 2: Actual Position (4 bytes) - Resolution: 0.1 µm
	if repBits&RespBitActualPosition != 0 {
		if err := ensureFieldBytes(data, offset, 4, "ActualPosition"); err != nil {
			return status, err
		}
		status.ActualPosition = int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
		offset += 4
	}

	// Bit 3: Demand Position (4 bytes) - Resolution: 0.1 µm
	if repBits&RespBitDemandPosition != 0 {
		if err := ensureFieldBytes(data, offset, 4, "DemandPosition"); err != nil {
			return status, err
		}
		status.DemandPosition = int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
		offset += 4
	}

	// Bit 4: Current (2 bytes) - Resolution: 1 mA
	if repBits&RespBitCurrent != 0 {
		if err := ensureFieldBytes(data, offset, 2, "Current"); err != nil {
			return status, err
		}
		status.Current = int16(binary.LittleEndian.Uint16(data[offset : offset+2]))
		offset += 2
	}

	// Bit 5: Warn Word (2 bytes) - Warning flags
	if repBits&RespBitWarnWord != 0 {
		if err := ensureFieldBytes(data, offset, 2, "WarnWord"); err != nil {
			return status, err
		}
		status.WarnWord = binary.LittleEndian.Uint16(data[offset : offset+2])
		offset += 2
	}

	// Bit 6: Error Code (2 bytes)
	if repBits&RespBitErrorCode != 0 {
		if err := ensureFieldBytes(data, offset, 2, "ErrorCode"); err != nil {
			return status, err
		}
		status.ErrorCode = binary.LittleEndian.Uint16(data[offset : offset+2])
		offset += 2
	}

	// Bit 7: Monitoring Channel (16 bytes = 4x 32-bit values)
	// These are values from 4 configurable UPIDs (0x20A8-0x20AB)
	if repBits&RespBitMonitoringChannel != 0 {
		if err := ensureFieldBytes(data, offset, 16, "MonitoringChannel"); err != nil {
			return status, err
		}
		for i := 0; i < 4; i++ {
			status.MonitoringChannel[i] = int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
			offset += 4
		}
	}

	// Bit 8: RTC Reply (8 bytes) is handled separately in RTC response parsing
	// We don't parse it here as it goes through ReadRTCResponse instead

	return status, nil
}

func ensureFieldBytes(data []byte, offset, length int, field string) error {
	if offset+length > len(data) {
		return fmt.Errorf("truncated packet: expected %s at offset %d, got %d bytes", field, offset, len(data))
	}
	return nil
}

// ReadStatus extracts Status from packet data using fixed offsets (legacy compatibility).
// For new code, prefer ReadStatusDynamic which handles variable-length responses.
func ReadStatus(data []byte, includeMonitoring bool) (*Status, error) {
	// Use dynamic parsing for robustness
	return ReadStatusDynamic(data)
}

func (s *Status) ActualPositionMM() float64 {
	return float64(s.ActualPosition) / float64(Factor.Position)
}

func (s *Status) DemandPositionMM() float64 {
	return float64(s.DemandPosition) / float64(Factor.Position)
}

func (s *Status) HasError() bool {
	return s.ErrorCode != 0
}

// HasMonitoringData returns true if any monitoring channel has non-zero data.
func (s *Status) HasMonitoringData() bool {
	return s.MonitoringChannel[0] != 0 || s.MonitoringChannel[1] != 0 ||
		s.MonitoringChannel[2] != 0 || s.MonitoringChannel[3] != 0
}

// StatusRequest requests drive status.
type StatusRequest struct{}

// NewStatusRequest creates a new status request.
func NewStatusRequest() *StatusRequest {
	return &StatusRequest{}
}

// ReadStatusRequest parses a status request (8 bytes, client → server).
func ReadStatusRequest(data []byte) (*StatusRequest, error) {
	if len(data) < PacketHeaderSize {
		return nil, NewPacketTooShortError("ReadStatusRequest", len(data), PacketHeaderSize, data)
	}
	return &StatusRequest{}, nil
}

// OperationTimeout returns the timeout duration for status requests.
func (*StatusRequest) OperationTimeout() time.Duration {
	return DefaultOperationTimeout
}

// WritePacket serializes a StatusRequest to a UDP packet (8 bytes).
func (request *StatusRequest) WritePacket() ([]byte, error) {
	packet := make([]byte, PacketHeaderSize)
	binary.LittleEndian.PutUint32(packet[0:4], 0x00000000)
	binary.LittleEndian.PutUint32(packet[4:8], ResponseFlags.Standard)
	return packet, nil
}

// ConnectivityProbeRequest is a status request that uses extended repBits (e.g., 0x1FF)
// for connectivity validation, matching C# library behavior.
type ConnectivityProbeRequest struct {
	repBits uint32
}

// NewConnectivityProbeRequest creates a connectivity probe request with custom repBits.
func NewConnectivityProbeRequest(repBits uint32) *ConnectivityProbeRequest {
	return &ConnectivityProbeRequest{repBits: repBits}
}

// OperationTimeout returns the timeout duration for connectivity probe requests.
func (*ConnectivityProbeRequest) OperationTimeout() time.Duration {
	return DefaultOperationTimeout
}

// WritePacket serializes a ConnectivityProbeRequest to a UDP packet (8 bytes).
func (request *ConnectivityProbeRequest) WritePacket() ([]byte, error) {
	packet := make([]byte, PacketHeaderSize)
	binary.LittleEndian.PutUint32(packet[0:4], 0x00000000) // reqBits = 0
	binary.LittleEndian.PutUint32(packet[4:8], request.repBits)
	return packet, nil
}

// StatusResponse indicates status data was received.
type StatusResponse struct {
	status *Status
}

// NewStatusResponse creates a new StatusResponse with the given values.
// This is primarily for testing and mock implementations.
func NewStatusResponse(status *Status) *StatusResponse {
	return &StatusResponse{
		status: status,
	}
}

// ReadStatusResponse parses a status response dynamically based on response definition bits.
// Handles variable-length responses (28, 32, 42, 48+ bytes) correctly.
func ReadStatusResponse(data []byte) (*StatusResponse, error) {
	if len(data) < PacketHeaderSize {
		return nil, NewPacketTooShortError("ReadStatusResponse", len(data), PacketHeaderSize, data)
	}

	// Use dynamic parsing to handle variable-length responses
	status, err := ReadStatusDynamic(data)
	if err != nil {
		return nil, err
	}

	return NewStatusResponse(status), nil
}

func (r *StatusResponse) Status() *Status {
	return r.status
}

// WritePacket serializes a StatusResponse to a UDP packet using dynamic format.
// Per LinUDP V2 spec: bits 0-6 = StatusWord, StateVar, ActPos, DemPos, Current, WarnWord, ErrorCode
func (r *StatusResponse) WritePacket() ([]byte, error) {
	// Response definition per LinUDP V2 protocol (bits 0-6)
	repBits := RespBitStatusWord | RespBitStateVar | RespBitActualPosition | RespBitDemandPosition |
		RespBitCurrent | RespBitWarnWord | RespBitErrorCode

	// Calculate packet size: header (8) + data based on bits
	// Bits 0-6: 2+2+4+4+2+2+2 = 18 bytes
	dataSize := 0
	if repBits&RespBitStatusWord != 0 {
		dataSize += 2 // Bit 0
	}
	if repBits&RespBitStateVar != 0 {
		dataSize += 2 // Bit 1
	}
	if repBits&RespBitActualPosition != 0 {
		dataSize += 4 // Bit 2
	}
	if repBits&RespBitDemandPosition != 0 {
		dataSize += 4 // Bit 3
	}
	if repBits&RespBitCurrent != 0 {
		dataSize += 2 // Bit 4
	}
	if repBits&RespBitWarnWord != 0 {
		dataSize += 2 // Bit 5
	}
	if repBits&RespBitErrorCode != 0 {
		dataSize += 2 // Bit 6
	}

	packet := make([]byte, PacketHeaderSize+dataSize)

	// Header
	binary.LittleEndian.PutUint32(packet[0:4], 0x00000000)
	binary.LittleEndian.PutUint32(packet[4:8], repBits)

	// Write data blocks in bit order per LinUDP V2 spec
	offset := PacketHeaderSize

	if repBits&RespBitStatusWord != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], r.status.StatusWord)
		offset += 2
	}
	if repBits&RespBitStateVar != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], r.status.StateVar)
		offset += 2
	}
	if repBits&RespBitActualPosition != 0 {
		binary.LittleEndian.PutUint32(packet[offset:offset+4], uint32(r.status.ActualPosition))
		offset += 4
	}
	if repBits&RespBitDemandPosition != 0 {
		binary.LittleEndian.PutUint32(packet[offset:offset+4], uint32(r.status.DemandPosition))
		offset += 4
	}
	if repBits&RespBitCurrent != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], uint16(r.status.Current))
		offset += 2
	}
	if repBits&RespBitWarnWord != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], r.status.WarnWord)
		offset += 2
	}
	if repBits&RespBitErrorCode != 0 {
		binary.LittleEndian.PutUint16(packet[offset:offset+2], r.status.ErrorCode)
		offset += 2
	}

	return packet, nil
}

// MonitoringStatusRequest requests drive status with monitoring channel data.
type MonitoringStatusRequest struct{}

// NewMonitoringStatusRequest creates a new monitoring status request.
func NewMonitoringStatusRequest() *MonitoringStatusRequest {
	return &MonitoringStatusRequest{}
}

// ReadMonitoringStatusRequest parses a monitoring status request (8 bytes, client → server).
func ReadMonitoringStatusRequest(data []byte) (*MonitoringStatusRequest, error) {
	if len(data) < PacketHeaderSize {
		return nil, NewPacketTooShortError("ReadMonitoringStatusRequest", len(data), PacketHeaderSize, data)
	}
	return &MonitoringStatusRequest{}, nil
}

// OperationTimeout returns the timeout duration for monitoring status requests.
func (*MonitoringStatusRequest) OperationTimeout() time.Duration {
	return DefaultOperationTimeout
}

// WritePacket serializes a MonitoringStatusRequest to a UDP packet (8 bytes).
func (request *MonitoringStatusRequest) WritePacket() ([]byte, error) {
	packet := make([]byte, PacketHeaderSize)
	binary.LittleEndian.PutUint32(packet[0:4], 0x00000000)                           // No request modules
	binary.LittleEndian.PutUint32(packet[4:8], ResponseFlags.StandardWithMonitoring) // Includes monitoring (bit 7)
	return packet, nil
}

// MonitoringStatusResponse indicates status data with monitoring channel was received.
type MonitoringStatusResponse struct {
	status *Status
}

// NewMonitoringStatusResponse creates a new MonitoringStatusResponse with the given values.
// This is primarily for testing and mock implementations.
func NewMonitoringStatusResponse(status *Status) *MonitoringStatusResponse {
	return &MonitoringStatusResponse{
		status: status,
	}
}

// ReadMonitoringStatusResponse parses a monitoring status response dynamically.
// Handles variable-length responses based on response definition bits.
func ReadMonitoringStatusResponse(data []byte) (*MonitoringStatusResponse, error) {
	if len(data) < PacketHeaderSize {
		return nil, NewPacketTooShortError("ReadMonitoringStatusResponse", len(data), PacketHeaderSize, data)
	}

	// Use dynamic parsing - it handles monitoring channels automatically based on bits 16-23
	status, err := ReadStatusDynamic(data)
	if err != nil {
		return nil, err
	}

	return NewMonitoringStatusResponse(status), nil
}

func (r *MonitoringStatusResponse) Status() *Status {
	return r.status
}

// WritePacket serializes a MonitoringStatusResponse to a UDP packet (44 bytes).
func (r *MonitoringStatusResponse) WritePacket() ([]byte, error) {
	packet := make([]byte, MinStatusPacketSize+16) // 28 + 16 = 44 bytes

	// Header
	binary.LittleEndian.PutUint32(packet[0:4], 0x00000000)
	binary.LittleEndian.PutUint32(packet[4:8], ResponseFlags.StandardWithMonitoring) // Includes bit 7

	// Status data (20 bytes, offsets 8-27)
	binary.LittleEndian.PutUint16(packet[8:10], r.status.StatusWord)
	binary.LittleEndian.PutUint16(packet[10:12], r.status.StateVar)
	binary.LittleEndian.PutUint32(packet[12:16], uint32(r.status.ActualPosition))
	binary.LittleEndian.PutUint32(packet[16:20], uint32(r.status.DemandPosition))
	binary.LittleEndian.PutUint16(packet[20:22], uint16(r.status.Current))
	binary.LittleEndian.PutUint16(packet[22:24], r.status.WarnWord)
	binary.LittleEndian.PutUint16(packet[24:26], r.status.ErrorCode)

	// Monitoring channel data (16 bytes, offsets 26-41)
	binary.LittleEndian.PutUint32(packet[26:30], uint32(r.status.MonitoringChannel[0]))
	binary.LittleEndian.PutUint32(packet[30:34], uint32(r.status.MonitoringChannel[1]))
	binary.LittleEndian.PutUint32(packet[34:38], uint32(r.status.MonitoringChannel[2]))
	binary.LittleEndian.PutUint32(packet[38:42], uint32(r.status.MonitoringChannel[3]))

	// Padding (2 bytes, offsets 42-43)
	packet[42] = 0
	packet[43] = 0

	return packet, nil
}
