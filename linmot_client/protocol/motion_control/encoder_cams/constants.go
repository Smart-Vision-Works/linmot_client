package protocol_encoder_cams

// ============================================================================
// Encoder CAM Sub ID Constants
// ============================================================================

// SubID represents an Encoder CAM command identifier.
type SubID uint8

// SubIDs groups all Encoder CAM command Sub ID constants.
// Reference: LinMot_MotionCtrl.txt, Section 4.3.63-4.3.68
var SubIDs = struct {
	SetupOnRisingWithDelay              SubID // 0x9 - Setup Encoder Cam On Rising Trigger Event With Delay Counts (069xh)
	SetupOnRisingWithDelayTargetPosLen  SubID // 0xA - Setup Encoder Cam On Rising Trigger Event With Delay Counts, Target Pos and Length (06Axh)
	SetupOnFallingWithDelay             SubID // 0xB - Setup Encoder Cam On Falling Trigger Event With Delay Counts (06Bxh)
	SetupOnFallingWithDelayTargetPosLen SubID // 0xC - Setup Encoder Cam On Falling Trigger Event With Delay Counts, Target Pos and Length (06Cxh)
	SetupOnRisingWithDelayAmpScaleLen   SubID // 0xD - Setup Encoder Cam On Rising Trigger Event With Delay Counts, Amplitude scale and Length (06Dxh)
	SetupOnFallingWithDelayAmpScaleLen  SubID // 0xE - Setup Encoder Cam On Falling Trigger Event With Delay Counts, Amplitude scale and Length (06Exh)
}{
	SetupOnRisingWithDelay:              0x9,
	SetupOnRisingWithDelayTargetPosLen:  0xA,
	SetupOnFallingWithDelay:             0xB,
	SetupOnFallingWithDelayTargetPosLen: 0xC,
	SetupOnRisingWithDelayAmpScaleLen:   0xD,
	SetupOnFallingWithDelayAmpScaleLen:  0xE,
}
