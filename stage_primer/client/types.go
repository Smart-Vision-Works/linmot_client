package client

// CommandTableParams contains the parameters needed for command table deployment
type CommandTableParams struct {
	DefaultSpeed        float64
	DefaultAcceleration float64
	PickTime            float64
	ZDistance            float64
	// LinMotIP is the direct IP address of the LinMot drive to deploy to.
	// This eliminates the dependency on the stage_primer ConfigStore being
	// pre-populated by SyncStagePrimers.
	LinMotIP string
}

// USBDevice represents a USB device
type USBDevice struct {
	Bus          string `json:"bus"`
	Device       string `json:"device"`
	IDVendor     string `json:"idVendor"`
	IDProduct    string `json:"idProduct"`
	Manufacturer string `json:"manufacturer"`
	Product      string `json:"product"`
	Serial       string `json:"serial"`
}

// LinMotConfig represents configuration for a single LinMot drive
type LinMotConfig struct {
	IP string `json:"ip"`
}

// ClearCoreConfig represents configuration for a single ClearCore device
type ClearCoreConfig struct {
	USBID                 string         `json:"usb_id"`
	DHCP                  bool           `json:"dhcp"`
	IPAddress             string         `json:"ip_address,omitempty"`
	Gateway               string         `json:"gateway,omitempty"`
	Subnet                string         `json:"subnet,omitempty"`
	DNS                   string         `json:"dns,omitempty"`
	RetransmissionTimeout uint8          `json:"retransmission_timeout,omitempty"`
	RetransmissionCount   uint8          `json:"retransmission_count,omitempty"`
	LinMots               []LinMotConfig `json:"linmots,omitempty"`
}

// StagePrimerConfig represents the root configuration structure
type StagePrimerConfig struct {
	ClearCores []ClearCoreConfig `json:"clearcores"`
}
