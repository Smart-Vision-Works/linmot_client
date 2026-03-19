package protocol_common

// PacketWritable represents events that can serialize themselves to UDP packets without a counter.
type PacketWritable interface {
	WritePacket() ([]byte, error)
}
