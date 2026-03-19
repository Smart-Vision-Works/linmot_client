package transport

// Server TransportServer receives requests and sends responses for server-side communication.
// Used by mock drive implementations for testing.
//
// Handles only byte I/O; protocol details are managed by the protocol package.
type Server interface {

	// SendPacket SendResponse sends a response to the client.
	SendPacket(response []byte)

	// RecvPacket ReceiveRequest receives a request from the client.
	// Returns nil if the transport is closed.
	RecvPacket() []byte

	// Close closes the transport.
	Close()
}
