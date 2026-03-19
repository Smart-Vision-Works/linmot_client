package protocol_common

// Response represents events that can be read from packets and written to packets.
type Response interface {
	PacketWritable
}
