package server

import (
	"testing"
)

func TestParseLSUSBVerboseOutput(t *testing.T) {
	sampleOutput := `Bus 002 Device 001: ID 1d6b:0003 Linux 6.12.32-v8 xhci-hcd xHCI Host Controller
Device Descriptor:
  bLength                18
  bDescriptorType         1
  bcdUSB               3.00
  bDeviceClass            9
  bDeviceSubClass         0
  bDeviceProtocol         3
  bMaxPacketSize0         9
  idVendor           0x1d6b
  idProduct          0x0003
  bcdDevice            6.12
  iManufacturer           3 Linux 6.12.32-v8 xhci-hcd
  iProduct                2 xHCI Host Controller
  iSerial                 1 0000:01:00.0

Bus 001 Device 004: ID 2890:8022 Teknic, Inc. Teknic ClearCore
Device Descriptor:
  bLength                18
  bDescriptorType         1
  bcdUSB               2.00
  bDeviceClass            2
  bDeviceSubClass         0
  bDeviceProtocol         0
  bMaxPacketSize0        64
  idVendor           0x2890
  idProduct          0x8022
  bcdDevice            1.00
  iManufacturer           1 Teknic, Inc.
  iProduct                2 Teknic ClearCore
  iSerial                 3 85293209534D394852202020FF161134

Bus 001 Device 002: ID 2109:3431  USB2.0 Hub
Device Descriptor:
  bLength                18
  bDescriptorType         1
  bcdUSB               2.10
  bDeviceClass            9
  bDeviceSubClass         0
  bDeviceProtocol         1
  bMaxPacketSize0        64
  idVendor           0x2109
  idProduct          0x3431
  bcdDevice            4.21
  iManufacturer           0
  iProduct                1 USB2.0 Hub
  iSerial                 0

Bus 001 Device 001: ID 1d6b:0002 Linux 6.12.32-v8 xhci-hcd xHCI Host Controller
Device Descriptor:
  bLength                18
  bDescriptorType         1
  bcdUSB               2.00
  bDeviceClass            9
  bDeviceSubClass         0
  bDeviceProtocol         1
  bMaxPacketSize0        64
  idVendor           0x1d6b
  idProduct          0x0002
  bcdDevice            6.12
  iManufacturer           3 Linux 6.12.32-v8 xhci-hcd
  iProduct                2 xHCI Host Controller
  iSerial                 1 0000:01:00.0`

	devices, err := parseLSUSBVerboseOutput(sampleOutput)
	if err != nil {
		t.Fatalf("parseLSUSBVerboseOutput returned error: %v", err)
	}

	expectedCount := 1 // Filtered out Linux host controllers and devices without serials
	if len(devices) != expectedCount {
		t.Fatalf("Expected %d devices, got %d", expectedCount, len(devices))
	}

	// Check the ClearCore device (Bus 001 Device 004) - should be the only device in filtered list
	clearCore := devices[0]
	if clearCore.Bus != "001" {
		t.Errorf("Expected Bus 001, got %s", clearCore.Bus)
	}
	if clearCore.Device != "004" {
		t.Errorf("Expected Device 004, got %s", clearCore.Device)
	}
	if clearCore.IDVendor != "2890" {
		t.Errorf("Expected IDVendor 2890, got %s", clearCore.IDVendor)
	}
	if clearCore.IDProduct != "8022" {
		t.Errorf("Expected IDProduct 8022, got %s", clearCore.IDProduct)
	}
	if clearCore.Manufacturer != "Teknic, Inc." {
		t.Errorf("Expected Manufacturer 'Teknic, Inc.', got '%s'", clearCore.Manufacturer)
	}
	if clearCore.Product != "Teknic ClearCore" {
		t.Errorf("Expected Product 'Teknic ClearCore', got '%s'", clearCore.Product)
	}
	if clearCore.Serial != "85293209534D394852202020FF161134" {
		t.Errorf("Expected Serial '85293209534D394852202020FF161134', got '%s'", clearCore.Serial)
	}

	t.Logf("Parsed devices: %+v", devices)
}
