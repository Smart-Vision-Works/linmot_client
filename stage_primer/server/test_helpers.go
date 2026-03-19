package server

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"stage_primer_config"
)

// TestServerConfig holds configuration for creating test servers
type TestServerConfig struct {
	Config config.Config
}

// CreateTestServer creates a test server with a temporary config file
func CreateTestServer(t *testing.T, config config.Config) (*Server, func()) {
	gin.SetMode(gin.TestMode)

	tempFile, cleanup := CreateTempConfigFile(t, config)

	server, err := NewServer(tempFile)
	require.NoError(t, err, "Failed to create test server")

	return server, cleanup
}


// CreateTempConfigFile creates a temporary config file for testing
func CreateTempConfigFile(t *testing.T, config config.Config) (string, func()) {
	tempDir, err := os.MkdirTemp("", "stage_primer_test_")
	require.NoError(t, err)

	tempFile := filepath.Join(tempDir, "config.json")
	err = saveConfigToFile(tempFile, &config)
	require.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempFile, cleanup
}

// CreateDefaultTestConfig creates a standard test configuration
func CreateDefaultTestConfig() config.Config {
	return config.Config{
		ClearCores: []config.ClearCoreConfig{
			{
				USBID:  "test-usb-id",
				DHCP:   true,
				LinMots: []config.LinMotConfig{
					{IP: "127.0.0.1"},
				},
			},
		},
	}
}

// CreateTestConfigWithMultipleLinMots creates a config with multiple LinMot drives
func CreateTestConfigWithMultipleLinMots() config.Config {
	return config.Config{
		ClearCores: []config.ClearCoreConfig{
			{
				USBID: "test-usb-id",
				DHCP:  true,
				LinMots: []config.LinMotConfig{
					{IP: "127.0.0.1"},
					{IP: "127.0.0.2"},
				},
			},
		},
	}
}

// CreateTestConfigWithStaticIP creates a config with static IP settings
func CreateTestConfigWithStaticIP() config.Config {
	return config.Config{
		ClearCores: []config.ClearCoreConfig{
			{
				USBID:     "test-usb-id",
				DHCP:      false,
				IPAddress: "192.168.1.100",
				Gateway:   "192.168.1.1",
				Subnet:    "255.255.255.0",
				DNS:       "8.8.8.8",
				LinMots: []config.LinMotConfig{
					{IP: "192.168.1.200"},
				},
			},
		},
	}
}

// CreateEmptyTestConfig creates an empty configuration
func CreateEmptyTestConfig() config.Config {
	return config.Config{
		ClearCores: []config.ClearCoreConfig{},
	}
}

// AssertConfigEqual asserts that two configs are equal
func AssertConfigEqual(t *testing.T, expected, actual config.Config) {
	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)

	actualJSON, err := json.Marshal(actual)
	require.NoError(t, err)

	require.JSONEq(t, string(expectedJSON), string(actualJSON))
}