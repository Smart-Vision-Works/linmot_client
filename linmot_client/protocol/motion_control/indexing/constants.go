package protocol_indexing

// ============================================================================
// Position Indexing Sub ID Constants
// ============================================================================

// SubID represents a Position Indexing command identifier.
type SubID uint8

// SubIDs groups all Position Indexing command Sub ID constants.
// Reference: LinMot_MotionCtrl.txt, Section 4.3.69-4.3.72
var SubIDs = struct {
	StartVAIEncoderIndexing       SubID // 0x0 - Start VAI Encoder Position Indexing (070xh)
	StartPredefVAIEncoderIndexing SubID // 0x1 - Start Predef VAI Encoder Position Indexing (071xh)
	StopIndexingVAIGoToPos        SubID // 0xE - Stop Position Indexing And VAI Go To Pos (07Exh)
	StopIndexingPredefVAIGoToPos  SubID // 0xF - Stop Position Indexing And Predefined VAI Go To Pos (07Fxh)
}{
	StartVAIEncoderIndexing:       0x0,
	StartPredefVAIEncoderIndexing: 0x1,
	StopIndexingVAIGoToPos:        0xE,
	StopIndexingPredefVAIGoToPos:  0xF,
}
