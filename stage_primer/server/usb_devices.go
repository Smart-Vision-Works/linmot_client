package server

import (
	"bufio"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// getMockUSBDevices returns fake USB devices for testing
func getMockUSBDevices() []USBDevice {
	return []USBDevice{
		{
			Bus:          "001",
			Device:       "004",
			IDVendor:     "2e8a",
			IDProduct:    "0005",
			Manufacturer: "Teknic Inc.",
			Product:      "Stage Primer v1.0",
			Serial:       "SP001",
		},
		{
			Bus:          "001",
			Device:       "005",
			IDVendor:     "2e8a",
			IDProduct:    "0005",
			Manufacturer: "Teknic Inc.",
			Product:      "Stage Primer v1.0",
			Serial:       "SP002",
		},
	}
}

// getUSBDevicesFunc is a function variable that allows mocking USB device enumeration in tests.
var getUSBDevicesFunc = getUSBDevicesReal

func getUSBDevicesReal() ([]USBDevice, error) {
	cmd := exec.Command("lsusb", "-v")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "lsusb command not found or failed to execute")
	}

	outputStr := string(output)
	devices, err := parseLSUSBVerboseOutput(outputStr)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// parseLSUSBVerboseOutput parses the output of lsusb -v to extract device information including serial numbers
func parseLSUSBVerboseOutput(output string) ([]USBDevice, error) {
	var allDevices []USBDevice
	var currentDevice *USBDevice

	scanner := bufio.NewScanner(strings.NewReader(output))

	// Regular expressions for parsing different parts
	busDeviceRegex := regexp.MustCompile(`(?i)Bus\s+(\d+)\s+Device\s+(\d+):\s+ID\s+([0-9a-f]+):([0-9a-f]+)(?:\s+(.+))?`)
	manufacturerRegex := regexp.MustCompile(`(?i)iManufacturer\s+\d+\s+(.+)`)
	productRegex := regexp.MustCompile(`(?i)iProduct\s+\d+\s+(.+)`)
	serialRegex := regexp.MustCompile(`(?i)iSerial\s+\d+\s+(.+)`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check for new device
		if matches := busDeviceRegex.FindStringSubmatch(line); matches != nil {
			// Save previous device if it exists and has required fields
			if currentDevice != nil && currentDevice.IDVendor != "" {
				allDevices = append(allDevices, *currentDevice)
			}

			// Start new device
			description := strings.TrimSpace(matches[5])
			currentDevice = &USBDevice{
				Bus:          matches[1],
				Device:       matches[2],
				IDVendor:     matches[3],
				IDProduct:    matches[4],
				Manufacturer: "",
				Product:      description,
				Serial:       "",
			}

			// Try to extract manufacturer and product from the initial description
			if description != "" {
				if strings.Contains(description, ", Inc.") || strings.Contains(description, ", Ltd.") {
					if idx := strings.Index(description, ", Inc. "); idx != -1 {
						currentDevice.Manufacturer = description[:idx+6]
						currentDevice.Product = description[idx+7:]
					} else if idx := strings.Index(description, ", Ltd. "); idx != -1 {
						currentDevice.Manufacturer = description[:idx+6]
						currentDevice.Product = description[idx+7:]
					}
				} else if strings.Contains(description, " ") {
					parts := strings.SplitN(description, " ", 2)
					if len(parts) == 2 {
						currentDevice.Manufacturer = parts[0]
						currentDevice.Product = parts[1]
					}
				}
			}
		} else if currentDevice != nil {
			// Parse additional device information from subsequent lines
			if matches := manufacturerRegex.FindStringSubmatch(line); matches != nil {
				currentDevice.Manufacturer = strings.TrimSpace(matches[1])
			} else if matches := productRegex.FindStringSubmatch(line); matches != nil {
				currentDevice.Product = strings.TrimSpace(matches[1])
			} else if matches := serialRegex.FindStringSubmatch(line); matches != nil {
				currentDevice.Serial = strings.TrimSpace(matches[1])
			}
		}
	}

	// Add the last device
	if currentDevice != nil && currentDevice.IDVendor != "" {
		allDevices = append(allDevices, *currentDevice)
	}

	// Filter out system-level devices
	var devices []USBDevice
	for _, device := range allDevices {
		if isUserDevice(device) {
			devices = append(devices, device)
		}
	}

	return devices, nil
}

// isUserDevice determines if a USB device is a user-connectable device (not system-level)
func isUserDevice(device USBDevice) bool {
	// Filter out Linux USB host controllers (xhci-hcd, ehci-hcd, etc.)
	if device.IDVendor == "1d6b" {
		return false
	}

	// Filter out devices without serial numbers (typically system hubs and controllers)
	if device.Serial == "" {
		return false
	}

	return true
}
