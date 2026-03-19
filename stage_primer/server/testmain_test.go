package server

import (
	"fmt"
	"os"
	"testing"

	"primer/linmot"
)

// serverPanicFactory is installed as the LinMot factory default during all tests
// in the server package. Any handler test that forgets to install its own mock
// factory will panic immediately instead of silently attempting a real UDP connection.
type serverPanicFactory struct{}

func (f *serverPanicFactory) CreateClient(ip string) (linmot.LinMotClient, error) {
	panic(fmt.Sprintf(
		"BUG: server test attempted to create a real LinMot connection to %q — "+
			"call linmot.SetClientFactory(mockFactory) in the test setup",
		ip,
	))
}

func (f *serverPanicFactory) Close() {}

// TestMain installs the panic factory before any test in this package runs
// and sets it as the reset target so ResetClientFactory never falls back to
// a real UDP factory between tests.
func TestMain(m *testing.M) {
	linmot.SetClientFactory(&serverPanicFactory{})
	linmot.SetDefaultFactory(&serverPanicFactory{})

	// Safeguard against real USB enumeration in tests
	getUSBDevicesFunc = func() ([]USBDevice, error) {
		panic("BUG: test attempted real USB enumeration — use server.SetMockMode(true) or mock getUSBDevicesFunc")
	}

	os.Exit(m.Run())
}
