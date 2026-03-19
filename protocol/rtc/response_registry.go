package protocol_rtc

import protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"

// ResponseRegistry creates a typed RTC response from parsed packet data.
// Domain packages (e.g., motion_control, command_tables) can register factories to create
// operation-specific response types (e.g., ReadVelocityResponse, AllocateEntryResponse)
// instead of generic RTCGetParamResponse or RTCSetParamResponse.
//
// This single registry type works for all RTC response types:
// - Parameter get operations (ReadROM, ReadRAM)
// - Parameter set operations (WriteROM, WriteRAM)
// - Command table operations (cmdCode determines operation type)
// - Future: Curve operations, error log access, etc.
type ResponseRegistry func(
	status *protocol_common.Status,
	value int32,
	upid uint16,
	rtcCounter uint8,
	rtcStatus uint8,
	cmdCode uint8,
) protocol_common.Response

var (
	responseRegistryMap = make(map[registryKey]ResponseRegistry)
)

type registryKey struct {
	cmdCode uint8
	hasCmd  bool
	upid    protocol_common.ParameterID
	hasUPID bool
}

func makeRegistryKey(cmdCode uint8, hasCmd bool, upid protocol_common.ParameterID, hasUPID bool) registryKey {
	return registryKey{
		cmdCode: cmdCode,
		hasCmd:  hasCmd,
		upid:    upid,
		hasUPID: hasUPID,
	}
}

// RegisterResponseRegistry registers a registry function for creating typed RTC responses
// based on parameter ID (UPID).
//
// When the parser encounters an RTC response with the given upid, it will call the
// registered registry to create a typed response instead of a generic response.
//
// This works for all RTC response types:
// - Parameter get operations: registry creates typed read responses (e.g., ReadVelocityResponse)
// - Parameter set operations: registry creates typed write responses (e.g., WritePosition1Response)
// - Command table operations: UPID contains entryID, registry creates CT responses
// - Future operations: UPID can contain curveID, logIndex, etc.
//
// This should be called during package initialization (init function) by domain packages
// that want to provide typed responses for specific operations.
//
// If multiple factories are registered for the same upid, the last one wins.
func RegisterResponseRegistry(upid protocol_common.ParameterID, registry ResponseRegistry) {
	responseRegistryMap[makeRegistryKey(0, false, upid, true)] = registry
}

// RegisterResponseRegistryByCmd registers a registry function for command-code keyed responses.
// Use this for CT/special commands that do not encode a unique UPID (often UPID=0).
func RegisterResponseRegistryByCmd(cmdCode uint8, registry ResponseRegistry) {
	responseRegistryMap[makeRegistryKey(cmdCode, true, 0, false)] = registry
}

// RegisterResponseRegistryByCmdAndUPID registers a registry keyed by both command code and UPID.
// Use when a command encodes entry IDs or indices in the UPID field and needs disambiguation.
func RegisterResponseRegistryByCmdAndUPID(cmdCode uint8, upid protocol_common.ParameterID, registry ResponseRegistry) {
	responseRegistryMap[makeRegistryKey(cmdCode, true, upid, true)] = registry
}

// LookupResponseRegistry looks up a registered registry for the given parameter ID.
// Returns the registry and true if found, or nil and false if not registered.
func LookupResponseRegistry(upid protocol_common.ParameterID) (ResponseRegistry, bool) {
	return lookupRegistry(makeRegistryKey(0, false, upid, true))
}

// LookupResponseRegistryByCmd looks up a registry registered for the given command code.
func LookupResponseRegistryByCmd(cmdCode uint8) (ResponseRegistry, bool) {
	return lookupRegistry(makeRegistryKey(cmdCode, true, 0, false))
}

// LookupResponseRegistryByCmdAndUPID looks up a registry registered for the given command code and UPID.
func LookupResponseRegistryByCmdAndUPID(cmdCode uint8, upid protocol_common.ParameterID) (ResponseRegistry, bool) {
	return lookupRegistry(makeRegistryKey(cmdCode, true, upid, true))
}

// LookupResponseRegistryCmdAware attempts cmd+UPID lookup, then cmd-only, then UPID-only.
func LookupResponseRegistryCmdAware(cmdCode uint8, upid protocol_common.ParameterID) (ResponseRegistry, bool) {
	if registry, ok := LookupResponseRegistryByCmdAndUPID(cmdCode, upid); ok {
		return registry, true
	}
	if registry, ok := LookupResponseRegistryByCmd(cmdCode); ok {
		return registry, true
	}
	if registry, ok := LookupResponseRegistry(upid); ok {
		return registry, true
	}
	return nil, false
}

func lookupRegistry(key registryKey) (ResponseRegistry, bool) {
	registry, ok := responseRegistryMap[key]
	return registry, ok
}

// RegisteredCmdResponseRegistryCount returns the number of cmd-code keyed factories.
// Intended primarily for tests and diagnostics.
func RegisteredCmdResponseRegistryCount() int {
	count := 0
	for k := range responseRegistryMap {
		if k.hasCmd {
			count++
		}
	}
	return count
}
