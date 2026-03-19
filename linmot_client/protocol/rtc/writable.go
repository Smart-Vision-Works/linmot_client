package protocol_rtc

// RTCPacketWritable represents RTC request events that can serialize themselves to UDP packets with a counter.
type RTCPacketWritable interface {
	WriteRtcPacket(counter uint8) ([]byte, error)
}
