package transport

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"
)

// ProbeResult holds the outcome of a cyclic telegram probe attempt.
type ProbeResult struct {
	Success    bool     // True if at least one response was received
	Attempts   int      // Total number of probe iterations sent
	Responses  int      // Number of successful responses
	FirstError string   // Error message from first failed attempt (if any)
	RawOutput  []string // Raw output lines for logging
	LocalPort  int      // Local UDP port used for the probe
	RemoteAddr string   // Remote address of the drive
}

// SendCyclicTelegramProbe sends C#-style cyclic telegrams (reqDef=6 or reqDef=7) to the drive
// and attempts to elicit UDP responses. This is used for raw connectivity validation.
//
// Parameters:
// - ctx: context for timeout/cancellation
// - transportClient: the transport client to use for communication
// - probes: list of probe configs (reqDef, length)
// - attempts: number of iterations per probe type
// - singleResponse: if true, return after first successful response (no further attempts)
//
// Returns ProbeResult with success status and detailed output.
func SendCyclicTelegramProbe(ctx context.Context, transportClient Client, probes []CyclicProbeConfig, attempts int, singleResponse bool) (*ProbeResult, error) {
	result := &ProbeResult{
		Attempts:  attempts,
		LocalPort: 41136, // Default master port
	}

	// Try each probe configuration
	for _, probeConfig := range probes {
		result.RawOutput = append(result.RawOutput, fmt.Sprintf("Starting probe: reqDef=0x%X, len=%d", probeConfig.ReqDef, probeConfig.Length))

		responses := 0
		for i := 0; i < attempts; i++ {
			// Build probe packet
			probe := make([]byte, probeConfig.Length)
			binary.LittleEndian.PutUint32(probe[0:4], probeConfig.ReqDef) // reqDef
			binary.LittleEndian.PutUint32(probe[4:8], 0x000001FF)         // respDef = 0x1FF (extended response)

			// For reqDef=7, set ControlWord at bytes 8-9
			if probeConfig.ReqDef == 7 && probeConfig.Length >= 10 {
				binary.LittleEndian.PutUint16(probe[8:10], 0x0006) // Safe DS402 control word
			}

			// Send probe via transport
			sendCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
			err := transportClient.SendPacket(sendCtx, probe)
			cancel()

			if err != nil {
				if result.FirstError == "" {
					result.FirstError = fmt.Sprintf("send failed: %v", err)
				}
				result.RawOutput = append(result.RawOutput, fmt.Sprintf("  Iteration %d: send error: %v", i, err))
				continue
			}

			result.RawOutput = append(result.RawOutput, fmt.Sprintf("  Iteration %d: sent %d-byte probe (reqDef=0x%X)", i, probeConfig.Length, probeConfig.ReqDef))

			// Try to receive response with 500ms timeout
			recvCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
			response, err := transportClient.RecvPacket(recvCtx)
			cancel()

			if err != nil {
				result.RawOutput = append(result.RawOutput, fmt.Sprintf("  Iteration %d: recv timeout (no response)", i))
				continue
			}

			// Got a response
			responses++
			result.Responses++
			result.Success = true

			// Parse status word for logging
			var sw uint16
			if len(response) >= 10 {
				sw = binary.LittleEndian.Uint16(response[8:10])
			}
			result.RawOutput = append(result.RawOutput, fmt.Sprintf("  Iteration %d: RESPONSE len=%d sw=0x%04X", i, len(response), sw))

			if singleResponse {
				// Return immediately after first successful response
				result.RawOutput = append(result.RawOutput, fmt.Sprintf("Probe succeeded: %d responses in %d attempts", responses, i+1))
				return result, nil
			}
		}

		if responses > 0 {
			result.RawOutput = append(result.RawOutput, fmt.Sprintf("Probe reqDef=0x%X: %d/%d responses", probeConfig.ReqDef, responses, attempts))
			if singleResponse {
				return result, nil
			}
		}
	}

	// Summary
	if result.Success {
		result.RawOutput = append(result.RawOutput, fmt.Sprintf("SUCCESS: %d/%d total responses", result.Responses, result.Attempts*len(probes)))
	} else {
		result.RawOutput = append(result.RawOutput, fmt.Sprintf("FAILED: 0 responses after %d attempts x %d probe types", result.Attempts, len(probes)))
		if result.FirstError == "" {
			result.FirstError = "no responses received"
		}
	}

	return result, nil
}

// CyclicProbeConfig specifies a single probe configuration.
type CyclicProbeConfig struct {
	ReqDef uint32 // Request definition (e.g., 6 or 7)
	Length int    // Packet length in bytes (typically 48 for reqDef=6, 50 for reqDef=7)
}

// DefaultCyclicProbes returns the standard probe configurations matching C# library behavior.
func DefaultCyclicProbes() []CyclicProbeConfig {
	return []CyclicProbeConfig{
		{ReqDef: 7, Length: 50}, // Steady-state cyclic (DS402 cyclical)
		{ReqDef: 6, Length: 48}, // Standard cyclic
	}
}
