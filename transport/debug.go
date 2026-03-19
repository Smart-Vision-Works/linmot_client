package transport

// DebugSetter allows transport implementations to receive debug enablement.
type DebugSetter interface {
	SetDebug(enabled bool)
}

// UDPInfoProvider exposes UDP socket details for debug logging.
type UDPInfoProvider interface {
	UDPInfo() UDPInfo
}

// UDPInfo describes configured and bound UDP socket info.
type UDPInfo struct {
	LocalAddr  string
	RemoteAddr string
	MasterPort int
	DrivePort  int
}
