package protocol_command_tables

import (
	"encoding/binary"
	"fmt"
)

// WireEntry represents a command table entry in wire format (on-the-wire representation).
// This is a pure wire-format type with no domain-specific metadata like YAML tags.
//
// Note: The ID field is not populated by DecodeEntry. It is only used when encoding
// from domain entries via BuildCTEntry, where it's set from the domain Entry.ID.
type WireEntry struct {
	ID             uint8
	Name           string
	Type           string
	SequencedEntry *uint8
	Par1           *int64
	Par2           *int64
	Par3           *int64
	Par4           *int64
	Par5           *int64
	Par6           *int64
	Par7           *int64
	Par8           *int64
}

// TypeSpec defines encoding/validation for a LinMot-Talk command type.
type TypeSpec struct {
	Name           string
	Header         uint16
	ValidateParams func(WireEntry) error
	EncodeParams   func(WireEntry, []byte) error // writes into bytes[6:38]
}

// ValidateWireEntry validates a WireEntry's type and parameters.
// Returns an error if the type is unsupported or if parameters are invalid.
func ValidateWireEntry(be WireEntry) error {
	spec, ok := getTypeSpec(be.Type)
	if !ok {
		return fmt.Errorf("unsupported type %q", be.Type)
	}
	return spec.ValidateParams(be)
}

// getTypeSpec returns a supported type spec by name (internal version).
func getTypeSpec(name string) (TypeSpec, bool) {
	switch name {
	case "VAI_GoToPos":
		return TypeSpec{
			Name:   name,
			Header: 0x0100, // MoveAbs
			ValidateParams: func(be WireEntry) error {
				if be.Par1 == nil || be.Par2 == nil || be.Par3 == nil || be.Par4 == nil {
					return fmt.Errorf("par1..par4 required")
				}
				return nil
			},
			EncodeParams: func(be WireEntry, buf []byte) error {
				// par1: S32 pos @ 6..9; par2: U32 vel @ 10..13; par3: U32 acc @ 14..17; par4: U32 dec @ 18..21
				binary.LittleEndian.PutUint32(buf[0:4], uint32(int32(*be.Par1)))
				binary.LittleEndian.PutUint32(buf[4:8], uint32(*be.Par2))
				binary.LittleEndian.PutUint32(buf[8:12], uint32(*be.Par3))
				binary.LittleEndian.PutUint32(buf[12:16], uint32(*be.Par4))
				for i := 16; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "MoveAbs":
		return getTypeSpec("VAI_GoToPos")
	case "MoveRel":
		return TypeSpec{
			Name:   name,
			Header: 0x0110,
			ValidateParams: func(be WireEntry) error {
				if be.Par1 == nil || be.Par2 == nil || be.Par3 == nil {
					return fmt.Errorf("par1..par3 required")
				}
				return nil
			},
			EncodeParams: func(be WireEntry, buf []byte) error {
				// par1: S32 delta @ 6..9; par2: U32 vel @ 10..13; par3: U32 acc @ 14..17; rest 0
				binary.LittleEndian.PutUint32(buf[0:4], uint32(int32(*be.Par1)))
				binary.LittleEndian.PutUint32(buf[4:8], uint32(*be.Par2))
				binary.LittleEndian.PutUint32(buf[8:12], uint32(*be.Par3))
				for i := 12; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "Home":
		return TypeSpec{
			Name:   name,
			Header: 0x0090,
			ValidateParams: func(be WireEntry) error {
				if be.Par1 == nil {
					return fmt.Errorf("par1 (home position) required")
				}
				return nil
			},
			EncodeParams: func(be WireEntry, buf []byte) error {
				binary.LittleEndian.PutUint32(buf[0:4], uint32(int32(*be.Par1)))
				for i := 4; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "Stop":
		// VAI Stop: header 0x0170, Par1 U32 decel @ 6..9
		return TypeSpec{
			Name:   name,
			Header: 0x0170,
			ValidateParams: func(be WireEntry) error {
				if be.Par1 == nil || *be.Par1 < 0 {
					return fmt.Errorf("par1 (decel U32) required and >=0")
				}
				return nil
			},
			EncodeParams: func(be WireEntry, buf []byte) error {
				binary.LittleEndian.PutUint32(buf[0:4], uint32(*be.Par1)) // decel
				for i := 4; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "Delay":
		return TypeSpec{
			Name:   name,
			Header: 0x2100,
			ValidateParams: func(be WireEntry) error {
				if be.Par1 == nil || *be.Par1 < 0 {
					return fmt.Errorf("par1 (time) required and >=0")
				}
				return nil
			},
			EncodeParams: func(be WireEntry, buf []byte) error {
				binary.LittleEndian.PutUint32(buf[0:4], uint32(*be.Par1))
				for i := 4; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "WaitDemandVelLT":
		// Wait Until Demand Velocity < threshold: header 0x2290, Par1 S32 @ 6..9
		return TypeSpec{
			Name:   name,
			Header: 0x2290,
			ValidateParams: func(be WireEntry) error {
				if be.Par1 == nil {
					return fmt.Errorf("par1 (S32 threshold) required")
				}
				return nil
			},
			EncodeParams: func(be WireEntry, buf []byte) error {
				binary.LittleEndian.PutUint32(buf[0:4], uint32(int32(*be.Par1)))
				for i := 4; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "NoOp":
		// No Operation: header 0x0000, no params
		return TypeSpec{
			Name:           name,
			Header:         0x0000,
			ValidateParams: func(WireEntry) error { return nil },
			EncodeParams: func(be WireEntry, buf []byte) error {
				for i := 0; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "WaitRising":
		return TypeSpec{
			Name:           name,
			Header:         0x2130,
			ValidateParams: func(WireEntry) error { return nil },
			EncodeParams: func(be WireEntry, buf []byte) error {
				for i := 0; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "WaitFalling":
		return TypeSpec{
			Name:           name,
			Header:         0x2140,
			ValidateParams: func(WireEntry) error { return nil },
			EncodeParams: func(be WireEntry, buf []byte) error {
				for i := 0; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "WaitMotionFinished":
		// Wait until Motion Finished [211xh]: header 0x2110, no params
		return TypeSpec{
			Name:           name,
			Header:         0x2110,
			ValidateParams: func(WireEntry) error { return nil },
			EncodeParams: func(be WireEntry, buf []byte) error {
				for i := 0; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "InfiniteMotionPos":
		// Predef VAI Infinite Motion +: header 0x02E0, no params
		return TypeSpec{Name: name, Header: 0x02E0, ValidateParams: func(WireEntry) error { return nil }, EncodeParams: func(be WireEntry, buf []byte) error {
			for i := 0; i < 32; i++ {
				buf[i] = 0x00
			}
			return nil
		}}, true
	case "InfiniteMotionNeg":
		// Predef VAI Infinite Motion -: header 0x02F0, no params
		return TypeSpec{Name: name, Header: 0x02F0, ValidateParams: func(WireEntry) error { return nil }, EncodeParams: func(be WireEntry, buf []byte) error {
			for i := 0; i < 32; i++ {
				buf[i] = 0x00
			}
			return nil
		}}, true
	case "InfiniteMotionPos_DecEqAcc":
		// Dec=Acc Infinite Motion +: header 0x0CE0, no params
		return TypeSpec{Name: name, Header: 0x0CE0, ValidateParams: func(WireEntry) error { return nil }, EncodeParams: func(be WireEntry, buf []byte) error {
			for i := 0; i < 32; i++ {
				buf[i] = 0x00
			}
			return nil
		}}, true
	case "InfiniteMotionNeg_DecEqAcc":
		// Dec=Acc Infinite Motion -: header 0x0CF0, no params
		return TypeSpec{Name: name, Header: 0x0CF0, ValidateParams: func(WireEntry) error { return nil }, EncodeParams: func(be WireEntry, buf []byte) error {
			for i := 0; i < 32; i++ {
				buf[i] = 0x00
			}
			return nil
		}}, true
	case "SetDO":
		return TypeSpec{
			Name:   name,
			Header: 0x0030,
			ValidateParams: func(be WireEntry) error {
				if be.Par1 == nil || be.Par2 == nil {
					return fmt.Errorf("par1 (mask U16) and par2 (value U16) required")
				}
				if *be.Par1 < 0 || *be.Par2 < 0 || *be.Par1 > 0xFFFF || *be.Par2 > 0xFFFF {
					return fmt.Errorf("mask/value must be 0..65535")
				}
				return nil
			},
			EncodeParams: func(be WireEntry, buf []byte) error {
				binary.LittleEndian.PutUint16(buf[0:2], uint16(*be.Par1)) // mask @ 6..7
				binary.LittleEndian.PutUint16(buf[2:4], uint16(*be.Par2)) // value @ 8..9
				for i := 4; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "SetX6":
		return TypeSpec{
			Name:   name,
			Header: 0x0040,
			ValidateParams: func(be WireEntry) error {
				if be.Par1 == nil || be.Par2 == nil {
					return fmt.Errorf("par1 (mask U16) and par2 (value U16) required")
				}
				if *be.Par1 < 0 || *be.Par2 < 0 || *be.Par1 > 0xFFFF || *be.Par2 > 0xFFFF {
					return fmt.Errorf("mask/value must be 0..65535")
				}
				return nil
			},
			EncodeParams: func(be WireEntry, buf []byte) error {
				binary.LittleEndian.PutUint16(buf[0:2], uint16(*be.Par1)) // mask @ 6..7
				binary.LittleEndian.PutUint16(buf[2:4], uint16(*be.Par2)) // value @ 8..9
				for i := 4; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "ClearX6":
		return TypeSpec{
			Name:   name,
			Header: 0x0040,
			ValidateParams: func(be WireEntry) error {
				if be.Par1 == nil {
					return fmt.Errorf("par1 (mask U16) required")
				}
				if *be.Par1 < 0 || *be.Par1 > 0xFFFF {
					return fmt.Errorf("mask must be 0..65535")
				}
				return nil
			},
			EncodeParams: func(be WireEntry, buf []byte) error {
				binary.LittleEndian.PutUint16(buf[0:2], uint16(*be.Par1)) // mask
				binary.LittleEndian.PutUint16(buf[2:4], 0)                // value=0 clears
				for i := 4; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "ClearDO":
		// Same header as SetDO with value forced to 0 for masked bits.
		return TypeSpec{
			Name:   name,
			Header: 0x0030,
			ValidateParams: func(be WireEntry) error {
				if be.Par1 == nil {
					return fmt.Errorf("par1 (mask U16) required")
				}
				// par2 ignored; force 0
				if *be.Par1 < 0 || *be.Par1 > 0xFFFF {
					return fmt.Errorf("mask must be 0..65535")
				}
				return nil
			},
			EncodeParams: func(be WireEntry, buf []byte) error {
				binary.LittleEndian.PutUint16(buf[0:2], uint16(*be.Par1)) // mask
				binary.LittleEndian.PutUint16(buf[2:4], 0)                // value=0 clears masked bits
				for i := 4; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	case "VAI_IncrementDemandPos":
		// Alias for MoveRel
		return getTypeSpec("MoveRel")
	case "VAI_Stop":
		return getTypeSpec("Stop")
	case "VAI_InfiniteMotionPlus":
		return getTypeSpec("InfiniteMotionPos")
	case "VAI_InfiniteMotionMinus":
		return getTypeSpec("InfiniteMotionNeg")
	case "VAI_DecEqAcc_InfiniteMotionPlus":
		return getTypeSpec("InfiniteMotionPos_DecEqAcc")
	case "VAI_DecEqAcc_InfiniteMotionMinus":
		return getTypeSpec("InfiniteMotionNeg_DecEqAcc")
	case "VAI_WaitDemandVelocityLT":
		return getTypeSpec("WaitDemandVelLT")
	case "VAI_WriteX4Mask":
		return getTypeSpec("SetDO")
	case "VAI_WriteX6Mask":
		return getTypeSpec("SetX6")
	case "ParamWrite":
		return TypeSpec{
			Name:   name,
			Header: 0x0020,
			ValidateParams: func(be WireEntry) error {
				if be.Par1 == nil || be.Par2 == nil {
					return fmt.Errorf("par1 (UPID U16) and par2 (U32 value) required")
				}
				if *be.Par1 < 0 || *be.Par1 > 0xFFFF {
					return fmt.Errorf("UPID must be 0..65535")
				}
				return nil
			},
			EncodeParams: func(be WireEntry, buf []byte) error {
				binary.LittleEndian.PutUint16(buf[0:2], uint16(*be.Par1)) // UPID @ 6..7
				binary.LittleEndian.PutUint32(buf[2:6], uint32(*be.Par2)) // value @ 8..11
				for i := 6; i < 32; i++ {
					buf[i] = 0x00
				}
				return nil
			},
		}, true
	// Conditional jump family — thresholds S32 @6..9, true ID @12..13, false ID @14..15
	case "IfDemandPosLT":
		return TypeSpec{
			Name:           name,
			Header:         0x2580,
			ValidateParams: getIfThresholdValidator(),
			EncodeParams:   getIfThresholdEncoder(),
		}, true
	case "IfDemandPosGT":
		return TypeSpec{Name: name, Header: 0x2590, ValidateParams: getIfThresholdValidator(), EncodeParams: getIfThresholdEncoder()}, true
	case "IfActualPosLT":
		return TypeSpec{Name: name, Header: 0x25A0, ValidateParams: getIfThresholdValidator(), EncodeParams: getIfThresholdEncoder()}, true
	case "IfActualPosGT":
		return TypeSpec{Name: name, Header: 0x25B0, ValidateParams: getIfThresholdValidator(), EncodeParams: getIfThresholdEncoder()}, true
	case "IfDiffPosLT":
		return TypeSpec{Name: name, Header: 0x25C0, ValidateParams: getIfThresholdValidator(), EncodeParams: getIfThresholdEncoder()}, true
	case "IfDiffPosGT":
		return TypeSpec{Name: name, Header: 0x25D0, ValidateParams: getIfThresholdValidator(), EncodeParams: getIfThresholdEncoder()}, true
	case "IfCurrentLT":
		return TypeSpec{Name: name, Header: 0x25E0, ValidateParams: getIfThresholdValidator(), EncodeParams: getIfThresholdEncoder()}, true
	case "IfCurrentGT":
		return TypeSpec{Name: name, Header: 0x25F0, ValidateParams: getIfThresholdValidator(), EncodeParams: getIfThresholdEncoder()}, true
	case "IfAnalogX44LT":
		return TypeSpec{Name: name, Header: 0x2600, ValidateParams: getIfThresholdValidator(), EncodeParams: getIfThresholdEncoder()}, true
	// Masked equality family — mask U16 @6..7, value U16 @8..9, true ID @10..11, false ID @12..13
	case "IfMaskedX4Eq":
		return TypeSpec{Name: name, Header: 0x2620, ValidateParams: getIfMaskedEqValidator(), EncodeParams: getIfMaskedEqEncoder()}, true
	case "IfMaskedX6Eq":
		return TypeSpec{Name: name, Header: 0x2630, ValidateParams: getIfMaskedEqValidator(), EncodeParams: getIfMaskedEqEncoder()}, true
	case "IfMaskedStatusEq":
		return TypeSpec{Name: name, Header: 0x2640, ValidateParams: getIfMaskedEqValidator(), EncodeParams: getIfMaskedEqEncoder()}, true
	case "IfMaskedWarnEq":
		return TypeSpec{Name: name, Header: 0x2650, ValidateParams: getIfMaskedEqValidator(), EncodeParams: getIfMaskedEqEncoder()}, true
	default:
		return TypeSpec{}, false
	}
}

// BuildCTEntry encodes a 64-byte CT entry from a WireEntry using the spec.
// It enforces A701h, linked ID, header, parameters, name, and reserved zeros.
func BuildCTEntry(be WireEntry) ([]byte, error) {
	spec, ok := getTypeSpec(be.Type)
	if !ok {
		return nil, fmt.Errorf("unsupported type %q", be.Type)
	}
	// 64-byte buffer
	out := make([]byte, 64)
	// 0..1 Version ID A701h (little endian bytes 0x01 0xA7)
	binary.LittleEndian.PutUint16(out[0:2], 0xA701)
	// 2..3 Linked entry (FFFFh = not linked)
	linked := uint16(0xFFFF)
	if be.SequencedEntry != nil {
		linked = uint16(*be.SequencedEntry)
	}
	binary.LittleEndian.PutUint16(out[2:4], linked)
	// 4..5 Motion header
	binary.LittleEndian.PutUint16(out[4:6], spec.Header)
	// 6..37 params
	params := out[6:38]
	for i := range params {
		params[i] = 0
	}
	if err := spec.EncodeParams(be, params); err != nil {
		return nil, err
	}
	// 38..53 name (NUL-terminated, zero pad)
	nameBytes := []byte(be.Name)
	if len(nameBytes) > 15 {
		nameBytes = nameBytes[:15]
	}
	copy(out[38:38+len(nameBytes)], nameBytes)
	// Always write a trailing NUL at 38+len(name), then zero pad remainder
	if 38+len(nameBytes) < 54 {
		out[38+len(nameBytes)] = 0x00
	}
	for i := 38 + len(nameBytes) + 1; i < 54; i++ {
		out[i] = 0x00
	}
	// 54..63 reserved zero (already zero)
	return out, nil
}

func getIfThresholdValidator() func(WireEntry) error {
	return func(be WireEntry) error {
		if be.Par1 == nil || be.Par2 == nil || be.Par3 == nil {
			return fmt.Errorf("par1 (S32 threshold), par2 (trueID), par3 (falseID) required")
		}
		if *be.Par2 < 1 || *be.Par2 > 255 || *be.Par3 < 1 || *be.Par3 > 255 {
			return fmt.Errorf("true/false IDs must be 1..255")
		}
		return nil
	}
}

func getIfThresholdEncoder() func(WireEntry, []byte) error {
	return func(be WireEntry, buf []byte) error {
		binary.LittleEndian.PutUint32(buf[0:4], uint32(int32(*be.Par1))) // threshold @6..9
		buf[4], buf[5] = 0x00, 0x00                                      // 10..11 zero
		binary.LittleEndian.PutUint16(buf[6:8], uint16(*be.Par2))        // true @12..13
		binary.LittleEndian.PutUint16(buf[8:10], uint16(*be.Par3))       // false @14..15
		for i := 10; i < 32; i++ {
			buf[i] = 0x00
		}
		return nil
	}
}

func getIfMaskedEqValidator() func(WireEntry) error {
	return func(be WireEntry) error {
		if be.Par1 == nil || be.Par2 == nil || be.Par3 == nil || be.Par4 == nil {
			return fmt.Errorf("par1 (mask U16), par2 (value U16), par3 (trueID), par4 (falseID) required")
		}
		if *be.Par1 < 0 || *be.Par1 > 0xFFFF || *be.Par2 < 0 || *be.Par2 > 0xFFFF {
			return fmt.Errorf("mask/value must be 0..65535")
		}
		if *be.Par3 < 1 || *be.Par3 > 255 || *be.Par4 < 1 || *be.Par4 > 255 {
			return fmt.Errorf("true/false IDs must be 1..255")
		}
		return nil
	}
}

func getIfMaskedEqEncoder() func(WireEntry, []byte) error {
	return func(be WireEntry, buf []byte) error {
		binary.LittleEndian.PutUint16(buf[0:2], uint16(*be.Par1)) // mask @6..7
		binary.LittleEndian.PutUint16(buf[2:4], uint16(*be.Par2)) // value @8..9
		binary.LittleEndian.PutUint16(buf[4:6], uint16(*be.Par3)) // true @10..11
		binary.LittleEndian.PutUint16(buf[6:8], uint16(*be.Par4)) // false @12..13
		for i := 8; i < 32; i++ {
			buf[i] = 0x00
		}
		return nil
	}
}

// headerToType maps motion header values to their Type string names.
// This is the reverse lookup for getTypeSpec.
var headerToType = map[uint16]string{
	0x0100: "VAI_GoToPos", // Also "MoveAbs"
	0x0110: "MoveRel",
	0x0090: "Home",
	0x0170: "Stop",
	0x2100: "Delay",
	0x2290: "WaitDemandVelLT",
	0x0000: "NoOp",
	0x2110: "WaitMotionFinished",
	0x2130: "WaitRising",
	0x2140: "WaitFalling",
	0x02E0: "InfiniteMotionPos",
	0x02F0: "InfiniteMotionNeg",
	0x0CE0: "InfiniteMotionPos_DecEqAcc",
	0x0CF0: "InfiniteMotionNeg_DecEqAcc",
	0x0030: "SetDO", // Also "ClearDO" (distinguished by value)
	0x0040: "SetX6", // Also "ClearX6" (distinguished by value)
	0x0020: "ParamWrite",
	0x2580: "IfDemandPosLT",
	0x2590: "IfDemandPosGT",
	0x25A0: "IfActualPosLT",
	0x25B0: "IfActualPosGT",
	0x25C0: "IfDiffPosLT",
	0x25D0: "IfDiffPosGT",
	0x25E0: "IfCurrentLT",
	0x25F0: "IfCurrentGT",
	0x2600: "IfAnalogX44LT",
	0x2620: "IfMaskedX4Eq",
	0x2630: "IfMaskedX6Eq",
	0x2640: "IfMaskedStatusEq",
	0x2650: "IfMaskedWarnEq",
}

// getTypeFromHeader returns the Type string for a given motion header.
// Returns empty string if header is unknown.
func getTypeFromHeader(header uint16) string {
	if typ, ok := headerToType[header]; ok {
		return typ
	}
	return ""
}

// DecodeEntry decodes a 64-byte binary CT entry payload into a WireEntry struct.
// This is the inverse of BuildCTEntry.
func DecodeEntry(data []byte) (*WireEntry, error) {
	// Verify size
	if len(data) < 64 {
		return nil, fmt.Errorf("entry data must be at least 64 bytes, got %d", len(data))
	}

	// Verify A701h version header
	version := binary.LittleEndian.Uint16(data[0:2])
	if version != 0xA701 {
		return nil, fmt.Errorf("invalid version header: expected 0xA701, got 0x%04X", version)
	}

	// Extract linked entry ID (bytes 2-3)
	linkedID := binary.LittleEndian.Uint16(data[2:4])
	var sequencedEntry *uint8
	if linkedID != 0xFFFF {
		id := uint8(linkedID)
		sequencedEntry = &id
	}

	// Extract motion header (bytes 4-5)
	header := binary.LittleEndian.Uint16(data[4:6])
	typ := getTypeFromHeader(header)
	if typ == "" {
		return nil, fmt.Errorf("unknown motion header: 0x%04X", header)
	}

	// Handle special cases where header alone doesn't determine type
	// SetDO vs ClearDO: Check if value is 0
	if header == 0x0030 {
		value := binary.LittleEndian.Uint16(data[8:10])
		if value == 0 {
			typ = "ClearDO"
		} else {
			typ = "SetDO"
		}
	}
	// SetX6 vs ClearX6: Check if value is 0
	if header == 0x0040 {
		value := binary.LittleEndian.Uint16(data[8:10])
		if value == 0 {
			typ = "ClearX6"
		} else {
			typ = "SetX6"
		}
	}

	// Extract parameters based on type
	params := data[6:38]
	entry := &WireEntry{
		Type:           typ,
		SequencedEntry: sequencedEntry,
	}

	// Decode parameters based on type
	if err := decodeParams(entry, typ, params); err != nil {
		return nil, fmt.Errorf("failed to decode parameters for type %q: %w", typ, err)
	}

	// Extract name (bytes 38-53, NUL-terminated)
	nameField := data[38:54]
	nulIdx := -1
	for i, b := range nameField {
		if b == 0 {
			nulIdx = i
			break
		}
	}
	if nulIdx == -1 {
		return nil, fmt.Errorf("name field not NUL-terminated")
	}
	entry.Name = string(nameField[:nulIdx])

	return entry, nil
}

// decodeParams decodes the parameter bytes based on the entry type.
func decodeParams(entry *WireEntry, typ string, params []byte) error {
	switch typ {
	case "VAI_GoToPos", "MoveAbs":
		// par1: S32 pos @ 6..9; par2: U32 vel @ 10..13; par3: U32 acc @ 14..17; par4: U32 dec @ 18..21
		pos := int64(int32(binary.LittleEndian.Uint32(params[0:4])))
		vel := int64(binary.LittleEndian.Uint32(params[4:8]))
		acc := int64(binary.LittleEndian.Uint32(params[8:12]))
		dec := int64(binary.LittleEndian.Uint32(params[12:16]))
		entry.Par1 = &pos
		entry.Par2 = &vel
		entry.Par3 = &acc
		entry.Par4 = &dec
		return nil

	case "MoveRel":
		// par1: S32 delta @ 6..9; par2: U32 vel @ 10..13; par3: U32 acc @ 14..17
		delta := int64(int32(binary.LittleEndian.Uint32(params[0:4])))
		vel := int64(binary.LittleEndian.Uint32(params[4:8]))
		acc := int64(binary.LittleEndian.Uint32(params[8:12]))
		entry.Par1 = &delta
		entry.Par2 = &vel
		entry.Par3 = &acc
		return nil

	case "Home":
		// par1: S32 home position @ 6..9
		pos := int64(int32(binary.LittleEndian.Uint32(params[0:4])))
		entry.Par1 = &pos
		return nil

	case "Stop":
		// par1: U32 decel @ 6..9
		decel := int64(binary.LittleEndian.Uint32(params[0:4]))
		entry.Par1 = &decel
		return nil

	case "Delay":
		// par1: U32 time @ 6..9
		time := int64(binary.LittleEndian.Uint32(params[0:4]))
		entry.Par1 = &time
		return nil

	case "WaitDemandVelLT":
		// par1: S32 threshold @ 6..9
		threshold := int64(int32(binary.LittleEndian.Uint32(params[0:4])))
		entry.Par1 = &threshold
		return nil

	case "NoOp", "WaitRising", "WaitFalling", "WaitMotionFinished", "InfiniteMotionPos", "InfiniteMotionNeg",
		"InfiniteMotionPos_DecEqAcc", "InfiniteMotionNeg_DecEqAcc":
		// No parameters
		return nil

	case "SetDO":
		// par1: U16 mask @ 6..7; par2: U16 value @ 8..9
		mask := int64(binary.LittleEndian.Uint16(params[0:2]))
		value := int64(binary.LittleEndian.Uint16(params[2:4]))
		entry.Par1 = &mask
		entry.Par2 = &value
		return nil

	case "ClearDO":
		// par1: U16 mask @ 6..7; value is always 0 for ClearDO, so we don't store Par2
		mask := int64(binary.LittleEndian.Uint16(params[0:2]))
		entry.Par1 = &mask
		// Par2 is nil for ClearDO (value is forced to 0)
		return nil

	case "SetX6":
		// par1: U16 mask @ 6..7; par2: U16 value @ 8..9
		mask := int64(binary.LittleEndian.Uint16(params[0:2]))
		value := int64(binary.LittleEndian.Uint16(params[2:4]))
		entry.Par1 = &mask
		entry.Par2 = &value
		return nil

	case "ClearX6":
		// par1: U16 mask @ 6..7; value is always 0 for ClearX6, so we don't store Par2
		mask := int64(binary.LittleEndian.Uint16(params[0:2]))
		entry.Par1 = &mask
		// Par2 is nil for ClearX6 (value is forced to 0)
		return nil

	case "ParamWrite":
		// par1: U16 UPID @ 6..7; par2: U32 value @ 8..11
		upid := int64(binary.LittleEndian.Uint16(params[0:2]))
		value := int64(binary.LittleEndian.Uint32(params[2:6]))
		entry.Par1 = &upid
		entry.Par2 = &value
		return nil

	case "IfDemandPosLT", "IfDemandPosGT", "IfActualPosLT", "IfActualPosGT",
		"IfDiffPosLT", "IfDiffPosGT", "IfCurrentLT", "IfCurrentGT", "IfAnalogX44LT":
		// par1: S32 threshold @ 6..9; par2: U16 trueID @ 12..13; par3: U16 falseID @ 14..15
		threshold := int64(int32(binary.LittleEndian.Uint32(params[0:4])))
		trueID := int64(binary.LittleEndian.Uint16(params[6:8]))
		falseID := int64(binary.LittleEndian.Uint16(params[8:10]))
		entry.Par1 = &threshold
		entry.Par2 = &trueID
		entry.Par3 = &falseID
		return nil

	case "IfMaskedX4Eq", "IfMaskedX6Eq", "IfMaskedStatusEq", "IfMaskedWarnEq":
		// par1: U16 mask @ 6..7; par2: U16 value @ 8..9; par3: U16 trueID @ 10..11; par4: U16 falseID @ 12..13
		mask := int64(binary.LittleEndian.Uint16(params[0:2]))
		value := int64(binary.LittleEndian.Uint16(params[2:4]))
		trueID := int64(binary.LittleEndian.Uint16(params[4:6]))
		falseID := int64(binary.LittleEndian.Uint16(params[6:8]))
		entry.Par1 = &mask
		entry.Par2 = &value
		entry.Par3 = &trueID
		entry.Par4 = &falseID
		return nil

	default:
		return fmt.Errorf("unsupported type for decoding: %q", typ)
	}
}
