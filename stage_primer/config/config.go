package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const (
	// DefaultConfigPath is the default path to the clear_core.json configuration file
	DefaultConfigPath = "/opt/svw/stage_primer/settings/clear_core.json"
)

// LinMotConfig represents configuration for a single LinMot drive
type LinMotConfig struct {
	IP string `json:"ip"`
}

// ClearCoreConfig represents configuration for a single ClearCore device
type ClearCoreConfig struct {
	USBID                 string         `json:"usb_id"`
	DHCP                  bool           `json:"dhcp,omitempty"`
	IPAddress             string         `json:"ip_address,omitempty"`
	Gateway               string         `json:"gateway,omitempty"`
	Subnet                string         `json:"subnet,omitempty"`
	DNS                   string         `json:"dns,omitempty"`
	RetransmissionTimeout uint8          `json:"retransmission_timeout,omitempty"`
	RetransmissionCount   uint8          `json:"retransmission_count,omitempty"`
	LinMots               []LinMotConfig `json:"linmots,omitempty"`
}

// Config represents the root configuration structure
type Config struct {
	ClearCores []ClearCoreConfig `json:"clearcores"`
}

// GetClearCoreConfig returns the configuration for a specific ClearCore device
func (c *Config) GetClearCoreConfig(connectionID string) (ClearCoreConfig, bool) {
	if c.ClearCores == nil {
		return ClearCoreConfig{}, false
	}
	for _, ccConfig := range c.ClearCores {
		if ccConfig.USBID == connectionID {
			return ccConfig, true
		}
	}
	return ClearCoreConfig{}, false
}

// GetAllLinMots returns all LinMot configurations from all ClearCores
func (c *Config) GetAllLinMots() []LinMotConfig {
	var allLinMots []LinMotConfig

	for _, ccConfig := range c.ClearCores {
		allLinMots = append(allLinMots, ccConfig.LinMots...)
	}

	return allLinMots
}

// ConfigStore holds the current configuration in memory so callers can read it
// without hitting the filesystem on every access. The zero value is not usable;
// create one with NewConfigStore.
type ConfigStore struct {
	mu  sync.RWMutex
	cfg Config
}

// NewConfigStore returns a ConfigStore pre-loaded with cfg.
func NewConfigStore(cfg Config) *ConfigStore {
	return &ConfigStore{cfg: cfg}
}

// Get returns the current in-memory configuration.
func (s *ConfigStore) Get() (Config, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg, nil
}

// Set replaces the stored configuration.
func (s *ConfigStore) Set(cfg Config) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cfg = cfg
}

// LoadConfigFromFile loads configuration from JSON file
func LoadConfigFromFile(configPath string) (Config, error) {
	var config Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Config file %s not found, using empty config\n", configPath)
		config.ClearCores = make([]ClearCoreConfig, 0)
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file %s: %v", configPath, err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse config file %s: %v", configPath, err)
	}

	if config.ClearCores == nil {
		config.ClearCores = make([]ClearCoreConfig, 0)
	}

	fmt.Printf("Configuration loaded from %s\n", configPath)
	return config, nil
}
