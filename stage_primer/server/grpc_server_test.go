package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"primer/linmot"
	pb "primer/proto"
	config "stage_primer_config"
)

func TestGRPCServer_LoadConfigFromStore(t *testing.T) {
	// Install mock factory: SetConfig now runs Setup for new LinMot IPs
	mockFactory := linmot.NewMockClientFactory()
	linmot.SetClientFactory(mockFactory)
	defer func() {
		mockFactory.Close()
		linmot.ResetClientFactory()
	}()

	initial := config.Config{
		ClearCores: []config.ClearCoreConfig{
			{
				USBID: "initial",
				LinMots: []config.LinMotConfig{
					{IP: "10.0.0.1"},
				},
			},
		},
	}

	updated := config.Config{
		ClearCores: []config.ClearCoreConfig{
			{
				USBID:     "updated",
				IPAddress: "10.0.0.2",
				Gateway:   "10.0.0.1",
				Subnet:    "255.255.248.0",
				DNS:       "8.8.8.8",
				LinMots: []config.LinMotConfig{
					{IP: "10.0.0.2"},
				},
			},
		},
	}

	configPath, cleanup := CreateTempConfigFile(t, initial)
	defer cleanup()

	grpcSrv := NewGRPCServer(initial)
	grpcSrv.SetConfigPath(configPath)

	// Initial load returns the startup config without any disk read.
	loaded, err := grpcSrv.loadConfig()
	require.NoError(t, err)
	assert.Equal(t, initial, loaded)

	// SetConfig writes to disk, updates the in-memory store, and runs Setup
	// for any new LinMot IPs.
	_, err = grpcSrv.SetConfig(context.Background(), &pb.SetConfigRequest{
		Clearcores: mapConfigToProto(updated),
	})
	require.NoError(t, err)

	// loadConfig now reflects the update from the store — no disk read.
	reloaded, err := grpcSrv.loadConfig()
	require.NoError(t, err)
	assert.Equal(t, updated, reloaded)
}

func TestGRPCServer_SetVacuum(t *testing.T) {
	// Implementation would require more mock infrastructure here or just test the dispatch logic
}

func TestGRPCServer_GetUSBDevices(t *testing.T) {
	s := NewGRPCServer(config.Config{})
	s.SetMockMode(true)

	resp, err := s.GetUSBDevices(nil, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Devices)
}

func TestGRPCServer_Config(t *testing.T) {
	initial := config.Config{
		ClearCores: []config.ClearCoreConfig{
			{
				USBID:     "initial",
				DHCP:      false,
				IPAddress: "10.0.0.40",
				Gateway:   "10.0.0.1",
				Subnet:    "255.255.255.0",
				DNS:       "8.8.8.8",
			},
		},
	}
	configPath, cleanup := CreateTempConfigFile(t, initial)
	defer cleanup()

	s := NewGRPCServer(initial)
	s.SetConfigPath(configPath)

	// Test GetConfig
	resp, err := s.GetConfig(nil, nil)
	require.NoError(t, err)
	assert.Len(t, resp.Clearcores, 1)
	assert.Equal(t, "initial", resp.Clearcores[0].UsbId)

	// Test SetConfig
	newCCs := resp.Clearcores
	newCCs = append(newCCs, &pb.ClearCoreConfig{
		UsbId:     "new",
		Dhcp:      false,
		IpAddress: "10.0.0.50",
		Gateway:   "10.0.0.1",
		Subnet:    "255.255.255.0",
		Dns:       "8.8.8.8",
	})
	_, err = s.SetConfig(nil, &pb.SetConfigRequest{Clearcores: newCCs})
	require.NoError(t, err)

	// Verify change reflected in GetConfig
	resp2, err := s.GetConfig(nil, nil)
	require.NoError(t, err)
	assert.Len(t, resp2.Clearcores, 2)
}
