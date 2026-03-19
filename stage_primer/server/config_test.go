package server

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stage_primer_config"
)

// TestConfigLoading tests basic config file loading
func TestConfigLoading(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a basic config
	cfg := config.Config{
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

	// Create temporary config file
	tempFile, cleanup := CreateTempConfigFile(t, cfg)
	defer cleanup()

	// Test loading the config
	loadedConfig, err := config.LoadConfigFromFile(tempFile)
	require.NoError(t, err)
	assert.Equal(t, cfg, loadedConfig)
}

// TestConfigValidation tests basic config validation
func TestConfigValidation(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		cfg := config.Config{
			ClearCores: []config.ClearCoreConfig{
				{
					USBID:     "valid-id",
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

		err := validateConfig(&cfg)
		assert.NoError(t, err)
	})

	t.Run("InvalidConfig", func(t *testing.T) {
		cfg := config.Config{
			ClearCores: []config.ClearCoreConfig{
				{
					USBID:  "test-id",
					DHCP:   false,
					// Missing required IP address
					Gateway: "192.168.1.1",
					Subnet:  "255.255.255.0",
					DNS:     "8.8.8.8",
				},
			},
		}

		err := validateConfig(&cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ip_address is required")
	})
}
