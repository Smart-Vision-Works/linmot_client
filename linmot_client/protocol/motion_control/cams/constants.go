package protocol_cams

// ============================================================================
// CAM Sub ID Constants
// ============================================================================

// SubID represents a CAM command identifier.
type SubID uint8

// SubIDs groups CAM command Sub ID constants for Master ID 0x05.
// Note: Limited documentation available for this Master ID group.
// Most cam functionality appears to be in Encoder CAM (Master ID 0x06).
var SubIDs = struct {
	// Placeholder for future cam commands if documented
	Reserved SubID // 0x0 - Reserved
}{
	Reserved: 0x0,
}
