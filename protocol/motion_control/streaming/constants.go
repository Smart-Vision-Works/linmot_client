package protocol_streaming

// ============================================================================
// Streaming Sub ID Constants
// ============================================================================

// SubID represents a Streaming command identifier.
type SubID uint8

// SubIDs groups all Streaming command Sub ID constants.
// Reference: LinMot_MotionCtrl.txt, Section 4.3.38-4.3.43
var SubIDs = struct {
	PStreamSlaveTimestamp               SubID // 0x0 - P Stream With Slave Generated Time Stamp (030xh)
	PVStreamSlaveTimestamp              SubID // 0x1 - PV Stream With Slave Generated Time Stamp (031xh)
	PStreamSlaveTimestampConfigPeriod   SubID // 0x2 - P Stream With Slave Generated Time Stamp and Configured Period Time (032xh)
	PVStreamSlaveTimestampConfigPeriod  SubID // 0x3 - PV Stream With Slave Generated Time Stamp and Configured Period Time (033xh)
	PVAStreamSlaveTimestamp             SubID // 0x4 - PVA Stream With Slave Generated Time Stamp (034xh)
	PVStreamSlaveTimestampConfigPeriod2 SubID // 0x5 - PV Stream With Slave Generated Time Stamp and Configured Period Time (035xh)
	StopStreaming                       SubID // 0xF - Stop Streaming (03Fxh)
}{
	PStreamSlaveTimestamp:               0x0,
	PVStreamSlaveTimestamp:              0x1,
	PStreamSlaveTimestampConfigPeriod:   0x2,
	PVStreamSlaveTimestampConfigPeriod:  0x3,
	PVAStreamSlaveTimestamp:             0x4,
	PVStreamSlaveTimestampConfigPeriod2: 0x5,
	StopStreaming:                       0xF,
}
