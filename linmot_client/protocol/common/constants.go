package protocol_common

import (
	"fmt"
	"reflect"
)

// ============================================================================
// Protocol Packet Constants
// ============================================================================

const (
	PacketHeaderSize    = 8  // Request flags (4) + Response flags (4)
	statusDataSize      = 20 // Standard status data size
	MinStatusPacketSize = 28 // Header (8) + status data (20)
)

// ============================================================================
// Protocol Flags
// ============================================================================

// RequestFlags groups all request flag constants.
// Usage: RequestFlags.RTCCommand, RequestFlags.ControlWord, RequestFlags.MotionControl
// Reference: LinUDP V2 protocol specification, Section 4.4.1
// Reference: C# library (LinUDP v2.1.1.0) - canonical production SDK
var RequestFlags = struct {
	// ControlWord indicates a control word is present (2 bytes).
	// Bit 0: Control Word present
	ControlWord uint32

	// MotionControl indicates an Motion Control command is present (32 bytes).
	// Bit 1: Motion Control present
	MotionControl uint32

	// RTCCommand indicates an RTC (Real-Time Configuration) command is present (8 bytes).
	// Bit 2: RTC command present
	// FIXED: Was incorrectly 0x00000100 (response bit 8), now correctly 0x00000004 (request bit 2)
	RTCCommand uint32
}{
	ControlWord:   0x00000001, // Bit 0
	MotionControl: 0x00000002, // Bit 1
	RTCCommand:    0x00000004, // Bit 2 (FIXED from 0x00000100)
}

// ResponseFlags groups all response flag constants.
// Usage: ResponseFlags.Standard, ResponseFlags.RTCReply
// Reference: LinUDP V2 protocol specification, Section 4.4.1
// Reference: C# library (LinUDP v2.1.1.0) - canonical production SDK
var ResponseFlags = struct {
	// Standard indicates standard status data is present (bits 0-6).
	// Includes: StatusWord, StateVar, ActualPosition, DemandPosition, Current, WarnWord, ErrorCode
	Standard uint32

	// StandardWithMonitoring includes standard status data plus monitoring channel (bits 0-7).
	StandardWithMonitoring uint32

	// RTCReply indicates standard status data plus RTC reply is present (bits 0-8).
	// Includes all standard fields plus RTC response data
	RTCReply uint32

	// MonitoringChannel is the bit flag for monitoring channel (bit 7).
	MonitoringChannel uint32

	// RTCResponseBit is the bit flag for RTC response data (bit 8).
	// Used in response definition to indicate RTC reply data is present
	RTCResponseBit uint32

	// All requests all available response data from the drive.
	// The drive will respond with whatever bits it supports.
	All uint32
}{
	Standard:               0x0000007F, // Bits 0-6
	StandardWithMonitoring: 0x000000FF, // Bits 0-7
	RTCReply:               0x000001FF, // Bits 0-8 (0x7F + 0x80 + 0x100)
	MonitoringChannel:      0x00000080, // Bit 7
	RTCResponseBit:         0x00000100, // Bit 8
	All:                    0xFFFFFFFF, // All bits - request everything
}

// ============================================================================
// Parameter ID Type (upid namespace)
// ============================================================================
// ParameterID represents a Universal Parameter ID (upid) address.
// Using a type creates a namespace effect and improves type safety.
// Reference: C library uses uint16 for upid addresses

type ParameterID uint16

// ============================================================================
// Parameter Universal IDs (PUID)
// ============================================================================
// PUID groups all Parameter Universal ID constants for LinMot drive parameters.
// These are memory addresses used in RTC read/write operations.
// Usage: PUID.Position1, PUID.Speed1, PUID.Input46Function
// Reference: C library linUDP.h and LinMot documentation

var PUID = struct {
	// Run Mode
	RunMode ParameterID // Run mode configuration (0x1450)
	// Position Control - Set 1
	Position1     ParameterID // Target position for Position Set 1 (0x145A)
	Speed1        ParameterID // Maximum velocity for Position Set 1 (0x145B)
	Acceleration1 ParameterID // Acceleration for Position Set 1 (0x145C)
	Deceleration1 ParameterID // Deceleration for Position Set 1 (0x145D)
	WaitTime1     ParameterID // Wait time for Position Set 1 (0x147D)

	// Position Control - Set 2
	Position2     ParameterID // Target position for Position Set 2 (0x145F)
	Speed2        ParameterID // Maximum velocity for Position Set 2 (0x1460)
	Acceleration2 ParameterID // Acceleration for Position Set 2 (0x1461)
	Deceleration2 ParameterID // Deceleration for Position Set 2 (0x1462)
	WaitTime2     ParameterID // Wait time for Position Set 2 (0x147E)

	// I/O Configuration - Input Functions
	Input45Function ParameterID // Input 4.5 function configuration (0x1060)
	Input46Function ParameterID // Input 4.6 function configuration (0x1061)
	Input47Function ParameterID // Input 4.7 function configuration (0x1062)
	Input48Function ParameterID // Input 4.8 function configuration (0x1063)

	// I/O Configuration - Output Functions
	Output36Function ParameterID // Output 3.6 configuration (0x1072)
	Output43Function ParameterID // Output 4.3 configuration (0x1070)
	Output44Function ParameterID // Output 4.4 configuration (0x1071)

	// Trigger Configuration
	TriggerMode ParameterID // Trigger mode configuration (0x170C)

	// Easy Steps Configuration
	EasyStepsAutoStart ParameterID // Easy Steps auto start configuration (0x30D4)
	EasyStepsAutoHome  ParameterID // Easy Steps auto home configuration (0x30D5)

	// Easy Steps - Input 4.5
	EasySteps45RisingEdge        ParameterID // Easy Steps 4.5 rising edge action (0x3600)
	EasySteps45IOMotionConfigCmd ParameterID // Easy Steps 4.5 IO motion config curve/CMD ID (0x3620)

	// Easy Steps - Input 4.6
	EasySteps46RisingEdge        ParameterID // Easy Steps 4.6 rising edge action (0x3700)
	EasySteps46IOMotionConfigCmd ParameterID // Easy Steps 4.6 IO motion config curve/CMD ID (0x3720)

	// Easy Steps - Input 4.7
	EasySteps47RisingEdge        ParameterID // Easy Steps 4.7 rising edge action (0x3800)
	EasySteps47IOMotionConfigCmd ParameterID // Easy Steps 4.7 IO motion config curve/CMD ID (0x3820)

	// Easy Steps - Input 4.8
	EasySteps48RisingEdge        ParameterID // Easy Steps 4.8 rising edge action (0x3100)
	EasySteps48IOMotionConfigCmd ParameterID // Easy Steps 4.8 IO motion config curve/CMD ID (0x3120)

	// Monitoring Channel Configuration
	MonitoringChannel1UPID ParameterID // Source UPID for Monitoring Channel 1 (0x20A8)
	MonitoringChannel2UPID ParameterID // Source UPID for Monitoring Channel 2 (0x20A9)
	MonitoringChannel3UPID ParameterID // Source UPID for Monitoring Channel 3 (0x20AA)
	MonitoringChannel4UPID ParameterID // Source UPID for Monitoring Channel 4 (0x20AB)
}{
	// Run Mode
	RunMode: 0x1450,

	// Position Control - Set 1
	Position1:     0x145A,
	Speed1:        0x145B,
	Acceleration1: 0x145C,
	Deceleration1: 0x145D,
	WaitTime1:     0x147D,

	// Position Control - Set 2
	Position2:     0x145F,
	Speed2:        0x1460,
	Acceleration2: 0x1461,
	Deceleration2: 0x1462,
	WaitTime2:     0x147E,

	// I/O Configuration - Input Functions
	Input45Function: 0x1060,
	Input46Function: 0x1061,
	Input47Function: 0x1062,
	Input48Function: 0x1063,

	// I/O Configuration - Output Functions
	Output36Function: 0x1072,
	Output43Function: 0x1070,
	Output44Function: 0x1071,

	// Trigger Configuration
	TriggerMode: 0x170C,

	// Easy Steps Configuration
	EasyStepsAutoStart: 0x30D4,
	EasyStepsAutoHome:  0x30D5,

	// Easy Steps - Input 4.5
	EasySteps45RisingEdge:        0x3600,
	EasySteps45IOMotionConfigCmd: 0x3620,

	// Easy Steps - Input 4.6
	EasySteps46RisingEdge:        0x3700,
	EasySteps46IOMotionConfigCmd: 0x3720,

	// Easy Steps - Input 4.7
	EasySteps47RisingEdge:        0x3800,
	EasySteps47IOMotionConfigCmd: 0x3820,

	// Easy Steps - Input 4.8
	EasySteps48RisingEdge:        0x3100,
	EasySteps48IOMotionConfigCmd: 0x3120,

	// Monitoring Channel Configuration
	MonitoringChannel1UPID: 0x20A8,
	MonitoringChannel2UPID: 0x20A9,
	MonitoringChannel3UPID: 0x20AA,
	MonitoringChannel4UPID: 0x20AB,
}

// Parameter is deprecated. Use PUID instead.
// Maintained for backward compatibility.
var Parameter = PUID

// ============================================================================
// Unit Conversion Factors
// ============================================================================
// Factor groups all unit conversion factor constants.
// Usage: Factor.Position, Factor.Speed, Factor.Acceleration
// Reference: C library linUDP.h lines 82-84
//   #define E1100_FACTOR_POS     10000     // from mm
//   #define E1100_FACTOR_SPEED   1000000   // from m/s
//   #define E1100_FACTOR_ACCEL   100000    // from m/s²
// C# library: Uses same scaling factors when calling SetRAM_ByUPID

var Factor = struct {
	// Position converts millimeters to drive units (0.0001mm resolution).
	// Usage: driveUnits = int32(positionMM * Factor.Position)
	// Example: 10.5mm → 105000 drive units
	// Reference: C library E1100_FACTOR_POS (10000)
	Position int32

	// Speed converts meters per second to drive units (0.000001 m/s resolution).
	// Usage: driveUnits = int32(velocityMS * Factor.Speed)
	// Example: 0.5 m/s → 500000 drive units
	// Reference: C library E1100_FACTOR_SPEED (1000000)
	Speed int32

	// Acceleration converts m/s² to drive units (0.00001 m/s² resolution).
	// Usage: driveUnits = int32(accelMS2 * Factor.Acceleration)
	// Example: 2.5 m/s² → 250000 drive units
	// Reference: C library E1100_FACTOR_ACCEL (100000)
	Acceleration int32
}{
	Position:     10000,
	Speed:        1000000,
	Acceleration: 100000,
}

// ============================================================================
// Drive Run Mode
// ============================================================================
// RunMode groups run mode selection constants.
// Usage: RunMode.Selection, RunMode.MotionCommandInterface
// Reference: C library linUDP.h lines 65-67
//   #define E1100_RUN_MODE_SELECTION        0x1450
//   #define E1100_RUN_MODE_SELECTION_MCI            1
//   #define E1100_RUN_MODE_SELECTION_VAI2POSCON     9
// C# library: Use SetRAM_ByUPID(TargetIP, 0x1450, mode) to set run mode

type RunMode int32

var RunModes = struct {
	// MotionCommandInterface sends real-time motion commands directly from PLC/automation system.
	// Reference: C library E1100_RUN_MODE_SELECTION_MCI (1)
	MotionCommandInterface RunMode

	// TriggeredVAInterpolator performs trapezoidal VA motion triggered by digital input.
	TriggeredVAInterpolator RunMode

	// RiseTriggeredVAIForwardBackwardMotion moves forward on rising edge, backward on falling edge.
	RiseTriggeredVAIForwardBackwardMotion RunMode

	// TriggeredTimeCurves executes predefined time-based motion curve on trigger event.
	TriggeredTimeCurves RunMode

	// CommandTableMode executes stored sequence of motion commands automatically.
	CommandTableMode RunMode

	// TriggeredCommandTable executes command table when triggered by digital inputs/events.
	TriggeredCommandTable RunMode

	// PositionIndexing switches between predefined positions by index value.
	PositionIndexing RunMode

	// Analog controls position/velocity directly via analog voltage input.
	Analog RunMode

	// TriggeredAnalog sets position from analog input, motion begins on trigger signal.
	TriggeredAnalog RunMode

	// VAI2PosContinuous cycles continuously between two predefined positions with VA interpolation.
	// Reference: C library E1100_RUN_MODE_SELECTION_VAI2POSCON (9)
	VAI2PosContinuous RunMode

	// ContinuousCurve runs user-defined curve in continuous repeated loop.
	ContinuousCurve RunMode

	// PCMotionCommandInterface enables control via PC interface for diagnostics/development.
	PCMotionCommandInterface RunMode
}{
	MotionCommandInterface:                0x0001,
	TriggeredVAInterpolator:               0x0002,
	RiseTriggeredVAIForwardBackwardMotion: 0x000D,
	TriggeredTimeCurves:                   0x0007,
	CommandTableMode:                      0x0003,
	TriggeredCommandTable:                 0x000C,
	PositionIndexing:                      0x000A,
	Analog:                                0x0004,
	TriggeredAnalog:                       0x0008,
	VAI2PosContinuous:                     0x0009,
	ContinuousCurve:                       0x0005,
	PCMotionCommandInterface:              0x0010,
}

// ============================================================================
// Easy Steps Auto Start Configuration Values
// ============================================================================
// EasyStepsAutoStart groups Easy Steps auto start configuration values (PUID 0x30D4).
// Usage: EasyStepsAutoStart.Disabled, EasyStepsAutoStart.Enabled

var EasyStepsAutoStart = struct {
	Disabled int32 // Auto start disabled (0x0000)
	Enabled  int32 // Auto start enabled (0x0001)
}{
	Disabled: 0x0000,
	Enabled:  0x0001,
}

// ============================================================================
// Easy Steps Auto Home Configuration Values
// ============================================================================
// EasyStepsAutoHome groups Easy Steps auto home configuration values (PUID 0x30D5).
// Usage: EasyStepsAutoHome.Disabled, EasyStepsAutoHome.Enabled

var EasyStepsAutoHome = struct {
	Disabled int32 // Auto home disabled (0x0000)
	Enabled  int32 // Auto home enabled (0x0001)
}{
	Disabled: 0x0000,
	Enabled:  0x0001,
}

// ============================================================================
// Output Configuration Values
// ============================================================================
// OutputConfig groups output pin configuration values.
// Usage: OutputConfig.None, OutputConfig.Brake, OutputConfig.AlwaysOn

var OutputConfig = struct {
	None                          int32 // No function (0x0000)
	Brake                         int32 // Brake output (0x0008)
	AlwaysOn                      int32 // Always on (0x000F)
	ApplicationOutput             int32 // Application output (0x000D)
	InterfaceOutput               int32 // Interface output (0x000E)
	StatusWordOperationEnabled    int32 // Status Word: Operation Enabled (0x0020)
	StatusWordEnableOperation     int32 // Status Word: Enable Operation (0x0022)
	StatusWordError               int32 // Status Word: Error (0x0023)
	StatusWordWarning             int32 // Status Word: Warning (0x0027)
	StatusWordSpecialMotionActive int32 // Status Word: Special Motion Active (0x0029)
	StatusWordInTargetPosition    int32 // Status Word: In Target Position (0x002A)
	StatusWordHomed               int32 // Status Word: Homed (0x002B)
	StatusWordFatalError          int32 // Status Word: Fatal Error (0x002C)
	StatusWordMotionActive        int32 // Status Word: Motion Active (0x002D)
	StatusWordRangeIndicator1     int32 // Status Word: Range Indicator 1 (0x002E)
	StatusWordRangeIndicator2     int32 // Status Word: Range Indicator 2 (0x002F)
}{
	None:                          0x0000,
	Brake:                         0x0008,
	AlwaysOn:                      0x000F,
	ApplicationOutput:             0x000D,
	InterfaceOutput:               0x000E,
	StatusWordOperationEnabled:    0x0020,
	StatusWordEnableOperation:     0x0022,
	StatusWordError:               0x0023,
	StatusWordWarning:             0x0027,
	StatusWordSpecialMotionActive: 0x0029,
	StatusWordInTargetPosition:    0x002A,
	StatusWordHomed:               0x002B,
	StatusWordFatalError:          0x002C,
	StatusWordMotionActive:        0x002D,
	StatusWordRangeIndicator1:     0x002E,
	StatusWordRangeIndicator2:     0x002F,
}

// ============================================================================
// Input Function Configuration Values
// ============================================================================
// InputFunction groups input pin function configuration values.
// Usage: InputFunction.None, InputFunction.Trigger, InputFunction.HomeSwitch

var InputFunction = struct {
	None                        int32 // No function (0x0000)
	Trigger                     int32 // Trigger (0x0001)
	HomeSwitch                  int32 // Home switch (0x0002)
	LimitSwitchNegative         int32 // Limit switch negative (0x0003)
	LimitSwitchPositive         int32 // Limit switch positive (0x0004)
	PTC1                        int32 // PTC 1 (0x0005)
	PTC2                        int32 // PTC 2 (0x0006)
	CtrlWordSwitchOn            int32 // Control Word: Switch On (0x0010)
	CtrlWordVoltageEnabled      int32 // Control Word: Voltage Enabled (0x0011)
	CtrlWordQuickStop           int32 // Control Word: Quick Stop (0x0012)
	CtrlWordEnableOperation     int32 // Control Word: Enable Operation (0x0013)
	CtrlWordAbort               int32 // Control Word: Abort (0x0014)
	CtrlWordFreeze              int32 // Control Word: Freeze (0x0015)
	CtrlWordGoToPosition        int32 // Control Word: Go To Position (0x0016)
	CtrlWordErrorAcknowledge    int32 // Control Word: Error Acknowledge (0x0017)
	CtrlWordJogMoveUp           int32 // Control Word: Jog Move Up (0x0018)
	CtrlWordJogMoveDown         int32 // Control Word: Jog Move Down (0x0019)
	CtrlWordSpecialMode         int32 // Control Word: Special Mode (0x001A)
	CtrlWordHome                int32 // Control Word: Home (0x001B)
	CtrlWordClearanceCheck      int32 // Control Word: Clearance Check (0x001C)
	CtrlWordGoToInitialPosition int32 // Control Word: Go To Initial Position (0x001D)
}{
	None:                        0x0000,
	Trigger:                     0x0001,
	HomeSwitch:                  0x0002,
	LimitSwitchNegative:         0x0003,
	LimitSwitchPositive:         0x0004,
	PTC1:                        0x0005,
	PTC2:                        0x0006,
	CtrlWordSwitchOn:            0x0010,
	CtrlWordVoltageEnabled:      0x0011,
	CtrlWordQuickStop:           0x0012,
	CtrlWordEnableOperation:     0x0013,
	CtrlWordAbort:               0x0014,
	CtrlWordFreeze:              0x0015,
	CtrlWordGoToPosition:        0x0016,
	CtrlWordErrorAcknowledge:    0x0017,
	CtrlWordJogMoveUp:           0x0018,
	CtrlWordJogMoveDown:         0x0019,
	CtrlWordSpecialMode:         0x001A,
	CtrlWordHome:                0x001B,
	CtrlWordClearanceCheck:      0x001C,
	CtrlWordGoToInitialPosition: 0x001D,
}

// ============================================================================
// Trigger Mode Configuration Values
// ============================================================================
// TriggerModeConfig groups trigger mode configuration values.
// Usage: TriggerModeConfig.None, TriggerModeConfig.Direct

var TriggerModeConfig = struct {
	None                int32 // No trigger mode (0x0000)
	Direct              int32 // Direct trigger mode (0x0001)
	Inhibited           int32 // Inhibited trigger mode (0x0002)
	Delayed             int32 // Delayed trigger mode (0x0003)
	InhibitedAndDelayed int32 // Inhibited and delayed trigger mode (0x0004)
}{
	None:                0x0000,
	Direct:              0x0001,
	Inhibited:           0x0002,
	Delayed:             0x0003,
	InhibitedAndDelayed: 0x0004,
}

// ============================================================================
// Easy Steps IO Motion Configuration Values
// ============================================================================
// EasyStepsIOMotion groups Easy Steps IO motion configuration values for rising edge actions.
// Usage: EasyStepsIOMotion.None, EasyStepsIOMotion.GotoAbsPosition

var EasyStepsIOMotion = struct {
	None                                       int32 // No action (0x0000)
	GotoAbsPosition                            int32 // Go to absolute position (0x0001)
	IncrementTargetPosition                    int32 // Increment target position (0x0002)
	IncrementDemandPosition                    int32 // Increment demand position (0x0003)
	GotoAbsPositionFromActualPosition          int32 // Go to absolute position from actual position (0x0004)
	IncrementActualPosition                    int32 // Increment actual position (0x0005)
	GoToAnalogPosition                         int32 // Go to analog position (0x0006)
	IncActualPositionBetweenRiseAndFallingEdge int32 // Increment actual position between rise and falling edge (0x0007)
	StartCurveFromActualPosition               int32 // Start curve from actual position (0x0008)
	GotoAbsPositionWithMaxCurrent              int32 // Go to absolute position with max current (0x0009)
	EvalCommandTableCommand                    int32 // Evaluate command table command (0x000C)
	VAIStop                                    int32 // VAI stop (0x000D)
	VAIInfiniteMotionPositionDirection         int32 // VAI infinite motion position direction (0x000E)
	VAIInfiniteMotionNegativeDirection         int32 // VAI infinite motion negative direction (0x000F)
	MasterHoming                               int32 // Master homing (0x001A)
}{
	None:                                       0x0000,
	GotoAbsPosition:                            0x0001,
	IncrementTargetPosition:                    0x0002,
	IncrementDemandPosition:                    0x0003,
	GotoAbsPositionFromActualPosition:          0x0004,
	IncrementActualPosition:                    0x0005,
	GoToAnalogPosition:                         0x0006,
	IncActualPositionBetweenRiseAndFallingEdge: 0x0007,
	StartCurveFromActualPosition:               0x0008,
	GotoAbsPositionWithMaxCurrent:              0x0009,
	EvalCommandTableCommand:                    0x000C,
	VAIStop:                                    0x000D,
	VAIInfiniteMotionPositionDirection:         0x000E,
	VAIInfiniteMotionNegativeDirection:         0x000F,
	MasterHoming:                               0x001A,
}

// ============================================================================
// Valid Parameter ID Validation Maps
// ============================================================================
// These maps define which PUIDs are valid for RAM/ROM write operations.
// Used for validation in request constructors.

// validPUIDs is a set of all valid PUID values for constant-time lookup
// Automatically built from the PUID struct using reflection at package initialization
// This map is cached and computed only once when the package loads
var validPUIDs map[ParameterID]struct{}

// init builds the validPUIDs map from the PUID struct using reflection
// This runs once at package initialization time and caches the result
func init() {
	validPUIDs = make(map[ParameterID]struct{})

	// Use reflection to iterate over all fields in the PUID struct
	v := reflect.ValueOf(PUID)
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		if fieldValue.CanInterface() { // Check if field is exported
			if paramID, ok := fieldValue.Interface().(ParameterID); ok {
				validPUIDs[paramID] = struct{}{}
			}
		}
	}
}

// IsValidPUID checks if a given UPID is defined in the PUID struct
// Uses constant-time map lookup for efficiency
func IsValidPUID(upid ParameterID) bool {
	_, exists := validPUIDs[upid]
	return exists
}

// ============================================================================
// Parameter Storage Type
// ============================================================================

// ParameterStorageType represents where to store parameter values (RAM or ROM).
// RAM values are temporary and lost on power cycle, ROM values are persistent.
type ParameterStorageType int

// ParameterStorage constants for specifying where to store parameter values.
var ParameterStorage = struct {
	// RAM stores the parameter value in temporary memory (lost on power cycle).
	RAM ParameterStorageType

	// ROM stores the parameter value in non-volatile memory (persistent).
	ROM ParameterStorageType

	// Both stores the parameter value in both RAM and ROM.
	Both ParameterStorageType
}{
	RAM:  0,
	ROM:  1,
	Both: 2,
}

// ============================================================================
// I/O Pin Type
// ============================================================================

// IOPinNumber represents an I/O pin identifier.
// Using a type creates a namespace effect and improves type safety.
type IOPinNumber int

// ============================================================================
// I/O Pin Constants
// ============================================================================

// IOPin groups all I/O pin identifier constants.
// Usage: IOPin.Input45, IOPin.Output43
var IOPin = struct {
	// Input Pins
	Input45 IOPinNumber // Input 4.5
	Input46 IOPinNumber // Input 4.6
	Input47 IOPinNumber // Input 4.7
	Input48 IOPinNumber // Input 4.8

	// Output Pins
	Output36 IOPinNumber // Output 3.6
	Output43 IOPinNumber // Output 4.3
	Output44 IOPinNumber // Output 4.4
}{
	// Input Pins
	Input45: 45,
	Input46: 46,
	Input47: 47,
	Input48: 48,

	// Output Pins
	Output36: 36,
	Output43: 43,
	Output44: 44,
}

// ============================================================================
// I/O Pin Mapping Helpers
// ============================================================================

// GetInputFunctionPUID returns the PUID for an input pin's function configuration.
// inputNumber should be IOPin.Input45, .Input46, .Input47, or .Input48.
// Returns error if inputNumber is not valid.
func GetInputFunctionPUID(inputNumber IOPinNumber) (ParameterID, error) {
	switch inputNumber {
	case IOPin.Input45:
		return PUID.Input45Function, nil
	case IOPin.Input46:
		return PUID.Input46Function, nil
	case IOPin.Input47:
		return PUID.Input47Function, nil
	case IOPin.Input48:
		return PUID.Input48Function, nil
	default:
		return 0, fmt.Errorf("invalid input number %d, must be IOPin.Input45, .Input46, .Input47, or .Input48", inputNumber)
	}
}

// GetOutputFunctionPUID returns the PUID for an output pin's function configuration.
// outputNumber should be IOPin.Output36, .Output43, or .Output44.
// Returns error if outputNumber is not valid.
func GetOutputFunctionPUID(outputNumber IOPinNumber) (ParameterID, error) {
	switch outputNumber {
	case IOPin.Output36:
		return PUID.Output36Function, nil
	case IOPin.Output43:
		return PUID.Output43Function, nil
	case IOPin.Output44:
		return PUID.Output44Function, nil
	default:
		return 0, fmt.Errorf("invalid output number %d, must be IOPin.Output36, .Output43, or .Output44", outputNumber)
	}
}

// GetEasyStepsRisingEdgePUID returns the PUID for an input pin's Easy Steps rising edge configuration.
// inputNumber should be IOPin.Input45, .Input46, .Input47, or .Input48.
// Returns error if inputNumber is not valid.
func GetEasyStepsRisingEdgePUID(inputNumber IOPinNumber) (ParameterID, error) {
	switch inputNumber {
	case IOPin.Input45:
		return PUID.EasySteps45RisingEdge, nil
	case IOPin.Input46:
		return PUID.EasySteps46RisingEdge, nil
	case IOPin.Input47:
		return PUID.EasySteps47RisingEdge, nil
	case IOPin.Input48:
		return PUID.EasySteps48RisingEdge, nil
	default:
		return 0, fmt.Errorf("invalid input number %d, must be IOPin.Input45, .Input46, .Input47, or .Input48", inputNumber)
	}
}

// GetEasyStepsIOMotionConfigCmdPUID returns the PUID for an input pin's Easy Steps IO motion config curve/CMD ID.
// inputNumber should be IOPin.Input45, .Input46, .Input47, or .Input48.
// Returns error if inputNumber is not valid.
func GetEasyStepsIOMotionConfigCmdPUID(inputNumber IOPinNumber) (ParameterID, error) {
	switch inputNumber {
	case IOPin.Input45:
		return PUID.EasySteps45IOMotionConfigCmd, nil
	case IOPin.Input46:
		return PUID.EasySteps46IOMotionConfigCmd, nil
	case IOPin.Input47:
		return PUID.EasySteps47IOMotionConfigCmd, nil
	case IOPin.Input48:
		return PUID.EasySteps48IOMotionConfigCmd, nil
	default:
		return 0, fmt.Errorf("invalid input number %d, must be IOPin.Input45, .Input46, .Input47, or .Input48", inputNumber)
	}
}

// ============================================================================
// Parameter Storage Type Conversion Helpers
// ============================================================================

// ToRTCCommandCode converts a ParameterStorageType to the corresponding RTC command code.
// Returns error if storage type is invalid.
func ToRTCCommandCode(storageType ParameterStorageType) (uint8, error) {
	switch storageType {
	case ParameterStorage.RAM:
		return 0x13, nil // WriteRAM
	case ParameterStorage.ROM:
		return 0x12, nil // WriteROM
	case ParameterStorage.Both:
		return 0x14, nil // WriteRAMAndROM
	default:
		return 0, fmt.Errorf("invalid storage type %d", storageType)
	}
}
