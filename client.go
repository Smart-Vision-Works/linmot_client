package linmot

import (
	"gsail-go/linmot/client"
	"gsail-go/linmot/transport"
)

// NewClient creates a new UDP client for communicating with a LinMot drive using default settings.
// This is a convenience wrapper that uses default ports and timeout.
func NewClient(driveIP string) (*client.Client, error) {
	return client.NewUDPClient(
		driveIP,
		transport.DefaultDrivePort,
		transport.DefaultMasterPort,
		"", // bind to all interfaces
		transport.DefaultTimeout,
	)
}
