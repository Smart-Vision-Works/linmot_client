package transport

import "time"

// Default UDP ports for LinMot communication
const (
	DefaultDrivePort  = 49360 // Drive listens on this port
	DefaultMasterPort = 41136 // Master uses ports in range 41100-41199
)

// DefaultTimeout is the default I/O timeout for UDP communication
const DefaultTimeout = 1000 * time.Millisecond
