package protocol_rtc

import (
	"encoding/binary"
	"fmt"
	"time"

	protocol_common "gsail-go/linmot/protocol/common"
)

// Compile-time interface checks
var (
	_ protocol_common.Request = (*RTCGetParamRequest)(nil)
	_ RTCPacketWritable       = (*RTCGetParamRequest)(nil)

	_ protocol_common.Request = (*RTCSetParamRequest)(nil)
	_ RTCPacketWritable       = (*RTCSetParamRequest)(nil)
)

// RTCGetParamRequest reads a parameter from ROM or RAM.
type RTCGetParamRequest struct {
	upid    uint16
	cmdCode uint8 // ReadROM or ReadRAM
}

// NewRTCGetParamRequest creates a new RTC parameter get request with the specified parameters.
func NewRTCGetParamRequest(upid uint16, cmdCode uint8) *RTCGetParamRequest {
	return &RTCGetParamRequest{
		upid:    upid,
		cmdCode: cmdCode,
	}
}

// ReadRTCGetParamRequest parses an RTC parameter get request (16 bytes, client → server).
// Returns the event and the counter extracted from the packet.
func ReadRTCGetParamRequest(data []byte) (any, uint8, error) {
	if len(data) < PacketSize {
		return nil, 0, protocol_common.NewPacketTooShortError("ReadRTCGetParamRequest", len(data), PacketSize, data)
	}

	rtcCounter := data[8]
	cmdCode := data[9]
	upid := binary.LittleEndian.Uint16(data[10:12])

	return &RTCGetParamRequest{
		upid:    upid,
		cmdCode: cmdCode,
	}, rtcCounter, nil
}

func (request *RTCGetParamRequest) UPID() uint16   { return request.upid }
func (request *RTCGetParamRequest) CmdCode() uint8 { return request.cmdCode }

// OperationTimeout returns the timeout duration for RTC parameter get requests.
func (*RTCGetParamRequest) OperationTimeout() time.Duration {
	return protocol_common.DefaultOperationTimeout
}

// WriteRtcPacket serializes an RTCGetParamRequest to a UDP packet (16 bytes) with the provided counter.
func (request *RTCGetParamRequest) WriteRtcPacket(counter uint8) ([]byte, error) {
	if err := ValidateRTCRead(request.upid, request.cmdCode); err != nil {
		return nil, err
	}

	packet := make([]byte, PacketSize)
	binary.LittleEndian.PutUint32(packet[0:4], protocol_common.RequestFlags.RTCCommand)
	binary.LittleEndian.PutUint32(packet[4:8], protocol_common.ResponseFlags.RTCReply)
	packet[8] = counter
	packet[9] = request.cmdCode
	binary.LittleEndian.PutUint16(packet[10:12], request.upid)
	binary.LittleEndian.PutUint32(packet[12:16], 0x00000000)
	return packet, nil
}

// RTCSetParamRequest writes a parameter to ROM or RAM.
type RTCSetParamRequest struct {
	upid    uint16
	value   int32
	cmdCode uint8 // WriteROM or WriteRAM
}

// NewRTCSetParamRequest creates a new RTC parameter set request with the specified parameters.
func NewRTCSetParamRequest(upid uint16, value int32, cmdCode uint8) *RTCSetParamRequest {
	return &RTCSetParamRequest{
		upid:    upid,
		value:   value,
		cmdCode: cmdCode,
	}
}

// NewWriteModeRequest creates a new RTC parameter set request for run mode (upid 0x1450).
func NewWriteModeRequestRAM(mode protocol_common.RunMode) (*RTCSetParamRequest, error) {
	upid := protocol_common.PUID.RunMode
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	return &RTCSetParamRequest{
		upid:    uint16(upid),
		value:   int32(mode),
		cmdCode: CommandCode.WriteRAM,
	}, nil
}

// NewWriteModeRequest creates a new RTC parameter set request for run mode (upid 0x1450).
func NewWriteModeRequestROM(mode protocol_common.RunMode) (*RTCSetParamRequest, error) {
	upid := protocol_common.PUID.RunMode
	if !protocol_common.IsValidPUID(upid) {
		return nil, fmt.Errorf("invalid upid 0x%04X", upid)
	}
	return &RTCSetParamRequest{
		upid:    uint16(upid),
		value:   int32(mode),
		cmdCode: CommandCode.WriteROM,
	}, nil
}

// NewSaveAllCurvesToFlashRequest creates a new RTC parameter set request to save all curves from RAM to Flash.
// This command does not require a upid or value (both are 0).
// Reference: LinUDP V2 specification, Command Code 0x40 (Curve Service: Save all Curves from RAM to Flash)
func NewSaveAllCurvesToFlashRequest() *RTCSetParamRequest {
	return &RTCSetParamRequest{
		upid:    0,
		value:   0,
		cmdCode: CommandCode.SaveAllCurvesToFlash,
	}
}

// ReadRTCSetParamRequest parses an RTC parameter set request (16 bytes, client → server).
// Returns the event and the counter extracted from the packet.
func ReadRTCSetParamRequest(data []byte) (any, uint8, error) {
	if len(data) < PacketSize {
		return nil, 0, protocol_common.NewPacketTooShortError("ReadRTCSetParamRequest", len(data), PacketSize, data)
	}

	rtcCounter := data[8]
	cmdCode := data[9]
	upid := binary.LittleEndian.Uint16(data[10:12])
	value := int32(binary.LittleEndian.Uint32(data[12:16]))

	return &RTCSetParamRequest{
		upid:    upid,
		value:   value,
		cmdCode: cmdCode,
	}, rtcCounter, nil
}

func (request *RTCSetParamRequest) UPID() uint16   { return request.upid }
func (request *RTCSetParamRequest) Value() int32   { return request.value }
func (request *RTCSetParamRequest) CmdCode() uint8 { return request.cmdCode }

// OperationTimeout returns the timeout duration for RTC parameter set requests.
// Flash save operations require extended timeouts because they write to non-volatile flash memory.
func (request *RTCSetParamRequest) OperationTimeout() time.Duration {
	if request.cmdCode == CommandCode.SaveAllCurvesToFlash ||
		request.cmdCode == 0x80 { // SaveCommandTable
		return protocol_common.FlashOperationTimeout
	}
	return protocol_common.DefaultOperationTimeout
}

// WriteRtcPacket serializes an RTCSetParamRequest to a UDP packet (16 bytes) with the provided counter.
func (request *RTCSetParamRequest) WriteRtcPacket(counter uint8) ([]byte, error) {
	if err := ValidateRTCWrite(request.upid, request.value, request.cmdCode); err != nil {
		return nil, err
	}

	packet := make([]byte, PacketSize)
	binary.LittleEndian.PutUint32(packet[0:4], protocol_common.RequestFlags.RTCCommand)
	binary.LittleEndian.PutUint32(packet[4:8], protocol_common.ResponseFlags.RTCReply)
	packet[8] = counter
	packet[9] = request.cmdCode
	binary.LittleEndian.PutUint16(packet[10:12], request.upid)
	binary.LittleEndian.PutUint32(packet[12:16], uint32(request.value))
	return packet, nil
}
