package server

import (
	"encoding/json"
	"fmt"
	"os"

	"stage_primer_config"
)

// validateConfig performs additional validation on the config beyond JSON binding
func validateConfig(cfg *config.Config) error {
	// Ensure clearcores array is initialized
	if cfg.ClearCores == nil {
		cfg.ClearCores = make([]config.ClearCoreConfig, 0)
		return nil
	}

	// Validate each ClearCore configuration
	usbIDs := make(map[string]bool)
	for i, cc := range cfg.ClearCores {
		// Check for duplicate USB IDs (only if USB ID is not empty)
		if cc.USBID != "" {
			if usbIDs[cc.USBID] {
				return fmt.Errorf("clearcore[%d]: duplicate usb_id '%s'", i, cc.USBID)
			}
			usbIDs[cc.USBID] = true
		}

		// Validate network settings if DHCP is disabled
		if !cc.DHCP {
			if cc.IPAddress == "" {
				return fmt.Errorf("clearcore[%d]: ip_address is required when dhcp is false", i)
			}
			if cc.Gateway == "" {
				return fmt.Errorf("clearcore[%d]: gateway is required when dhcp is false", i)
			}
			if cc.Subnet == "" {
				return fmt.Errorf("clearcore[%d]: subnet is required when dhcp is false", i)
			}
			if cc.DNS == "" {
				return fmt.Errorf("clearcore[%d]: dns is required when dhcp is false", i)
			}
		}

		// Validate LinMot configurations
		linmotIPs := make(map[string]bool)
		for j, lm := range cc.LinMots {
			if lm.IP == "" {
				return fmt.Errorf("clearcore[%d].linmot[%d]: ip is required", i, j)
			}
			if linmotIPs[lm.IP] {
				return fmt.Errorf("clearcore[%d].linmot[%d]: duplicate ip '%s'", i, j, lm.IP)
			}
			linmotIPs[lm.IP] = true
		}
	}

	return nil
}

// saveConfigToFile writes the config to file with proper formatting
func saveConfigToFile(filePath string, cfg *config.Config) error {
	// Convert to JSON with indentation
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %v", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
