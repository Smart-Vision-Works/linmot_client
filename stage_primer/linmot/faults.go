// Package linmot implements the LinMot fault lifecycle and auto-recovery policy.
// It manages fault detection, classification, and retry state per drive.
package linmot

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	pkgerrors "github.com/pkg/errors"

	"github.com/Smart-Vision-Works/staged_robot/client"
	protocol_control_word "github.com/Smart-Vision-Works/staged_robot/protocol/control_word"

	clearcore "stage_primer_config"
)

const (
	faultPollInterval = time.Second
	faultPollTimeout  = 750 * time.Millisecond

	// Auto-recovery policy (intentional): at most 3 acknowledge attempts per 60s per IP.
	// After that budget is exhausted, faults are escalated immediately.
	autoRecoveryMaxFaults = 3
	autoRecoveryWindow    = 60 * time.Second
)

// FaultLevel classifies how a fault should be handled.
type FaultLevel int

const (
	FaultLevelWarning  FaultLevel = iota // Warning only — log internally, do NOT broadcast
	FaultLevelNonFatal                   // Non-fatal error — attempt auto-recovery within budget
	FaultLevelFatal                      // Fatal error — escalate immediately
)

type FaultListener func(ip string, err error)

// FaultEscalationListener is called when a fault is escalated (fatal or budget exceeded).
// The ip is the LinMot IP and isFatal indicates whether it was a fatal error.
type FaultEscalationListener func(ip string, fault *client.DriveFaultError, isFatal bool)

var (
	listenersMu sync.RWMutex
	listeners   = map[uint64]FaultListener{}
	nextID      uint64

	escalationListenersMu sync.RWMutex
	escalationListeners   = map[uint64]FaultEscalationListener{}
	escalationNextID      uint64
)

// FaultLifecycleState models the active incident state for a single LinMot IP.
type FaultLifecycleState int

const (
	FaultLifecycleStateHealthy FaultLifecycleState = iota
	FaultLifecycleStateRecovering
	FaultLifecycleStateEscalated
)

func (s FaultLifecycleState) String() string {
	switch s {
	case FaultLifecycleStateHealthy:
		return "Healthy"
	case FaultLifecycleStateRecovering:
		return "Recovering"
	case FaultLifecycleStateEscalated:
		return "Escalated"
	default:
		return "Unknown"
	}
}

type faultLifecycleRecord struct {
	State                       FaultLifecycleState
	LastFaultCode               uint16
	LastFaultLevel              FaultLevel
	LastFaultSignature          string
	LastFaultBroadcastSignature string
	LastFaultBroadcastTime      time.Time
	LastRecoverAttemptTime      time.Time
	LastEscalationSignature     string
	LastEscalationTime          time.Time
	ConsecutiveRecoverFailures  int
}

type faultLifecycleTracker struct {
	mu      sync.Mutex
	records map[string]faultLifecycleRecord // LinMot IP -> active incident lifecycle
}

var globalFaultLifecycle = &faultLifecycleTracker{
	records: make(map[string]faultLifecycleRecord),
}

// faultBudget tracks auto-recovery attempts per LinMot IP using a sliding window.
type faultBudget struct {
	mu      sync.Mutex
	windows map[string][]time.Time // IP -> timestamps of recent recovery attempts
}

var globalFaultBudget = &faultBudget{
	windows: make(map[string][]time.Time),
}

func faultSignature(fault *client.DriveFaultError, level FaultLevel) string {
	if fault == nil {
		return fmt.Sprintf("level=%d:nil", level)
	}
	// Only include stable fault-identifying fields in the signature.
	// StatusWord and StateVar contain volatile bits (counter nibble, motion-active, in-target)
	// that change between polls even when the fault persists. Use only ErrorCode,
	// WarningWord, and the main state byte (bits 8-15) of StateVar which identify the fault.
	mainStateByte := (fault.StateVar >> 8) & 0xFF
	return fmt.Sprintf(
		"level=%d:mainstate=0x%02X:error=0x%04X:warn=0x%04X",
		level, mainStateByte, fault.ErrorCode, fault.WarningWord,
	)
}

func (t *faultLifecycleTracker) markHealthy(ip string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	rec := t.records[ip]
	rec.State = FaultLifecycleStateHealthy
	rec.LastFaultCode = 0
	rec.LastFaultLevel = FaultLevelWarning
	rec.LastFaultSignature = ""
	rec.LastFaultBroadcastSignature = ""
	rec.LastFaultBroadcastTime = time.Time{}
	rec.LastEscalationSignature = ""
	rec.LastEscalationTime = time.Time{}
	rec.ConsecutiveRecoverFailures = 0
	t.records[ip] = rec
}

func (t *faultLifecycleTracker) beginRecovery(ip string, fault *client.DriveFaultError, level FaultLevel) {
	t.mu.Lock()
	defer t.mu.Unlock()

	rec := t.records[ip]
	rec.State = FaultLifecycleStateRecovering
	rec.LastFaultCode = fault.ErrorCode
	rec.LastFaultLevel = level
	rec.LastFaultSignature = faultSignature(fault, level)
	rec.LastRecoverAttemptTime = time.Now()
	t.records[ip] = rec
}

func (t *faultLifecycleTracker) markRecoveryFailed(ip string, fault *client.DriveFaultError, level FaultLevel) {
	t.mu.Lock()
	defer t.mu.Unlock()

	rec := t.records[ip]
	rec.State = FaultLifecycleStateRecovering
	rec.LastFaultCode = fault.ErrorCode
	rec.LastFaultLevel = level
	rec.LastFaultSignature = faultSignature(fault, level)
	rec.ConsecutiveRecoverFailures++
	t.records[ip] = rec
}

func (t *faultLifecycleTracker) shouldBroadcastFault(ip string, fault *client.DriveFaultError, level FaultLevel) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	rec := t.records[ip]
	signature := faultSignature(fault, level)
	if rec.LastFaultBroadcastSignature == signature {
		return false
	}
	rec.LastFaultBroadcastSignature = signature
	rec.LastFaultBroadcastTime = time.Now()
	t.records[ip] = rec
	return true
}

func (t *faultLifecycleTracker) shouldBroadcastEscalation(ip string, fault *client.DriveFaultError, level FaultLevel) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	rec := t.records[ip]
	signature := faultSignature(fault, level)
	if rec.LastEscalationSignature == signature {
		return false
	}
	rec.LastEscalationSignature = signature
	rec.LastEscalationTime = time.Now()
	t.records[ip] = rec
	return true
}

func (t *faultLifecycleTracker) markEscalated(ip string, fault *client.DriveFaultError, level FaultLevel) {
	t.mu.Lock()
	defer t.mu.Unlock()

	rec := t.records[ip]
	rec.State = FaultLifecycleStateEscalated
	rec.LastFaultCode = fault.ErrorCode
	rec.LastFaultLevel = level
	rec.LastFaultSignature = faultSignature(fault, level)
	t.records[ip] = rec
}

func (t *faultLifecycleTracker) get(ip string) faultLifecycleRecord {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.records[ip]
}

func (t *faultLifecycleTracker) resetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.records = make(map[string]faultLifecycleRecord)
}

// FaultLifecycleSnapshot captures lifecycle metadata for one LinMot IP.
type FaultLifecycleSnapshot struct {
	State                      FaultLifecycleState
	LastFaultCode              uint16
	LastFaultLevel             FaultLevel
	LastRecoverAttemptTime     time.Time
	LastEscalationTime         time.Time
	ConsecutiveRecoverFailures int
}

func getFaultLifecycleSnapshot(ip string) FaultLifecycleSnapshot {
	rec := globalFaultLifecycle.get(ip)
	return FaultLifecycleSnapshot{
		State:                      rec.State,
		LastFaultCode:              rec.LastFaultCode,
		LastFaultLevel:             rec.LastFaultLevel,
		LastRecoverAttemptTime:     rec.LastRecoverAttemptTime,
		LastEscalationTime:         rec.LastEscalationTime,
		ConsecutiveRecoverFailures: rec.ConsecutiveRecoverFailures,
	}
}

func resetFaultLifecycleStateForTests() {
	globalFaultLifecycle.resetAll()
}

// tryConsume returns true if an auto-recovery attempt is within budget.
// If within budget, records the attempt. If exceeded, returns false.
func (b *faultBudget) tryConsume(ip string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-autoRecoveryWindow)

	// Prune expired entries
	times := b.windows[ip]
	pruned := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			pruned = append(pruned, t)
		}
	}

	if len(pruned) >= autoRecoveryMaxFaults {
		b.windows[ip] = pruned
		return false // Budget exceeded
	}

	b.windows[ip] = append(pruned, now)
	return true
}

// reset clears the budget for a given IP.
func (b *faultBudget) reset(ip string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.windows, ip)
}

// ResetFaultBudget resets the auto-recovery budget for the given LinMot IP.
func ResetFaultBudget(ip string) {
	globalFaultBudget.reset(ip)
}

// AddFaultListener registers a listener for fault events.
// It returns a cleanup function that unregisters the listener.
func AddFaultListener(l FaultListener) func() {
	listenersMu.Lock()
	id := nextID
	nextID++
	listeners[id] = l
	listenersMu.Unlock()

	return func() {
		listenersMu.Lock()
		delete(listeners, id)
		listenersMu.Unlock()
	}
}

// AddEscalationListener registers a listener for escalated faults (fatal or budget-exceeded).
// Returns a cleanup function.
func AddEscalationListener(l FaultEscalationListener) func() {
	escalationListenersMu.Lock()
	id := escalationNextID
	escalationNextID++
	escalationListeners[id] = l
	escalationListenersMu.Unlock()

	return func() {
		escalationListenersMu.Lock()
		delete(escalationListeners, id)
		escalationListenersMu.Unlock()
	}
}

func broadcastFault(ip string, err error) {
	listenersMu.RLock()
	defer listenersMu.RUnlock()
	for _, l := range listeners {
		l(ip, err)
	}
}

// BroadcastClearCoreFault is used by the ClearCore fault bridge to inject
// ClearCore device faults into the linmot fault notification pipeline.
func BroadcastClearCoreFault(deviceID string, err error) {
	broadcastFault(deviceID, err)
}

func broadcastEscalation(ip string, fault *client.DriveFaultError, isFatal bool) {
	escalationListenersMu.RLock()
	defer escalationListenersMu.RUnlock()
	for _, l := range escalationListeners {
		l(ip, fault, isFatal)
	}
}

// classifyFault determines the fault level from a DriveFaultError.
func classifyFault(fault *client.DriveFaultError) FaultLevel {
	if fault == nil {
		return FaultLevelWarning
	}

	// Fatal errors cannot be recovered
	helper := protocol_control_word.NewStatusWordHelper(fault.StatusWord, fault.StateVar)
	if helper.IsFatalError() {
		return FaultLevelFatal
	}

	// Non-zero error code = non-fatal error (recoverable)
	if fault.ErrorCode != 0 {
		return FaultLevelNonFatal
	}

	// Warning only (WarnWord != 0, ErrorCode == 0)
	return FaultLevelWarning
}

// MonitorFaults polls LinMot drives for faults on a fixed interval.
// Logs faults and continues polling until context is canceled.
func MonitorFaults(ctx context.Context, config clearcore.Config) error {
	return MonitorFaultsWithConfigProvider(ctx, func() (clearcore.Config, error) {
		return config, nil
	})
}

// MonitorFaultsWithConfigProvider polls LinMot drives for faults on a fixed
// interval. Configuration is loaded on each tick so runtime config updates are
// applied without requiring process restart.
func MonitorFaultsWithConfigProvider(ctx context.Context, loadConfig func() (clearcore.Config, error)) error {
	ticker := time.NewTicker(faultPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			config, err := loadConfig()
			if err != nil {
				fmt.Printf("[LinMot] fault poll config load failed: %v\n", err)
				continue
			}
			for _, lm := range config.GetAllLinMots() {
				ip := strings.TrimSpace(lm.IP)
				if ip == "" {
					continue
				}
				if IsInRecovery(ip) {
					continue
				}
				if err := checkDriveFaultOnce(ctx, ip); err != nil {
					handleDetectedFault(ctx, ip, err)
					continue
				}
				// A clean poll means the incident ended; reset lifecycle to Healthy.
				globalFaultLifecycle.markHealthy(ip)
			}
		}
	}
}

// handleDetectedFault is the Stage Primer policy layer for LinMot incidents.
// The protocol primitive (AcknowledgeError) stays in gsail-go; policy decisions
// (retry budget, lifecycle transitions, and escalation dedupe) stay here.
func handleDetectedFault(ctx context.Context, ip string, faultErr error) {
	var fault *client.DriveFaultError
	if !errors.As(faultErr, &fault) {
		// Not a structured drive fault (e.g., network error, context deadline).
		// During deployments, the fault monitor is paused via recovery state, so
		// any connectivity error that reaches here is a genuine issue — broadcast it.
		logFault(ip, faultErr)
		broadcastFault(ip, faultErr)
		return
	}

	level := classifyFault(fault)

	switch level {
	case FaultLevelWarning:
		// Warning only — log internally, do NOT broadcast to UI
		logFaultInternal(ip, fault, "warning-only")
		globalFaultLifecycle.markHealthy(ip)
		return

	case FaultLevelNonFatal:
		// Attempt auto-recovery within budget
		if globalFaultBudget.tryConsume(ip) {
			globalFaultLifecycle.beginRecovery(ip, fault, level)
			logFaultInternal(ip, fault, "auto-recovery attempt")
			if tryAutoRecover(ctx, ip) {
				logFaultInternal(ip, fault, "auto-recovery succeeded")
				globalFaultLifecycle.markHealthy(ip)
				return
			}
			logFaultInternal(ip, fault, "auto-recovery failed")
			globalFaultLifecycle.markRecoveryFailed(ip, fault, level)
			if globalFaultLifecycle.shouldBroadcastFault(ip, fault, level) {
				logFault(ip, faultErr)
				broadcastFault(ip, faultErr)
			} else {
				logFaultInternal(ip, fault, "duplicate non-fatal fault notification suppressed during recovery")
			}
			return
		}

		logFaultInternal(ip, fault, "budget exceeded — escalating")
		globalFaultLifecycle.markEscalated(ip, fault, level)
		if globalFaultLifecycle.shouldBroadcastFault(ip, fault, level) {
			logFault(ip, faultErr)
			broadcastFault(ip, faultErr)
		} else {
			logFaultInternal(ip, fault, "duplicate fault notification suppressed on escalation")
		}
		if globalFaultLifecycle.shouldBroadcastEscalation(ip, fault, level) {
			broadcastEscalation(ip, fault, false)
		} else {
			logFaultInternal(ip, fault, "duplicate escalation suppressed")
		}

	case FaultLevelFatal:
		globalFaultLifecycle.markEscalated(ip, fault, level)
		if globalFaultLifecycle.shouldBroadcastFault(ip, fault, level) {
			logFault(ip, faultErr)
			broadcastFault(ip, faultErr)
		} else {
			logFaultInternal(ip, fault, "persistent escalated incident unchanged; suppressing duplicate fault notification")
		}
		if globalFaultLifecycle.shouldBroadcastEscalation(ip, fault, level) {
			broadcastEscalation(ip, fault, true)
		} else {
			logFaultInternal(ip, fault, "persistent escalated incident unchanged; suppressing duplicate escalation")
		}
	}
}

// tryAutoRecover attempts to acknowledge the error on a LinMot drive.
func tryAutoRecover(ctx context.Context, ip string) bool {
	ackCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	linmotClient, err := globalClientFactory.CreateClient(ip)
	if err != nil {
		fmt.Printf("[LinMot %s] auto-recovery: failed to create client: %v\n", ip, err)
		return false
	}

	_, err = linmotClient.AcknowledgeError(ackCtx)
	if err != nil {
		fmt.Printf("[LinMot %s] auto-recovery: AcknowledgeError failed: %v\n", ip, err)
		return false
	}

	return true
}

func checkDriveFaultOnce(ctx context.Context, ip string) error {
	pollCtx, cancel := context.WithTimeout(ctx, faultPollTimeout)
	defer cancel()

	linmotClient, err := globalClientFactory.CreateClient(ip)
	if err != nil {
		return pkgerrors.Wrapf(err, "failed to create LinMot client for %s", ip)
	}

	return linmotClient.CheckDriveFault(pollCtx)
}

// logFaultInternal logs a fault at debug level (not forwarded to UI).
func logFaultInternal(ip string, fault *client.DriveFaultError, action string) {
	fmt.Printf("[LinMot %s] %s: error_code=0x%04X warning_word=0x%04X\n",
		ip, action, fault.ErrorCode, fault.WarningWord)
}

func logFault(ip string, err error) {
	var fault *client.DriveFaultError
	if errors.As(err, &fault) {
		line := fmt.Sprintf("[LinMot %s] fault detected status_word=0x%04X state_var=0x%04X error_code=0x%04X error_text=%q warning_word=0x%04X warning_text=%q",
			ip, fault.StatusWord, fault.StateVar, fault.ErrorCode, fault.ErrorText, fault.WarningWord, fault.WarningText)

		var probeErr *client.DriveFaultProbeError
		if errors.As(err, &probeErr) && probeErr.ProbeErr != nil {
			line = fmt.Sprintf("%s probe_err=%v", line, probeErr.ProbeErr)
		}

		fmt.Println(line)
		return
	}

	fmt.Printf("[LinMot %s] fault check failed: %v\n", ip, err)
}

// HumanReadableFault generates a human-readable fault description that includes
// the stage identity (robot/stage index) derived from the config.
func HumanReadableFault(ip string, fault *client.DriveFaultError, cfg clearcore.Config) string {
	identity := resolveStageIdentity(ip, cfg)

	if fault == nil {
		return fmt.Sprintf("%s: unknown fault", identity)
	}

	level := classifyFault(fault)
	var levelStr string
	switch level {
	case FaultLevelFatal:
		levelStr = "FATAL"
	case FaultLevelNonFatal:
		levelStr = "ERROR"
	default:
		levelStr = "WARNING"
	}

	msg := fmt.Sprintf("%s [%s]", identity, levelStr)

	if fault.ErrorText != "" {
		msg += fmt.Sprintf(": %s", fault.ErrorText)
	} else if fault.ErrorCode != 0 {
		msg += fmt.Sprintf(": error code 0x%04X", fault.ErrorCode)
	}

	if fault.WarningText != "" {
		if fault.ErrorText != "" || fault.ErrorCode != 0 {
			msg += fmt.Sprintf("; warning: %s", fault.WarningText)
		} else {
			msg += fmt.Sprintf(": %s", fault.WarningText)
		}
	}

	return msg
}

// resolveStageIdentity maps a LinMot IP to a human-readable stage identity.
func resolveStageIdentity(ip string, cfg clearcore.Config) string {
	for robotIdx, cc := range cfg.ClearCores {
		for stageIdx, lm := range cc.LinMots {
			if strings.TrimSpace(lm.IP) == strings.TrimSpace(ip) {
				return fmt.Sprintf("Robot %d / Stage %d (%s)", robotIdx+1, stageIdx+1, ip)
			}
		}
	}
	return fmt.Sprintf("LinMot %s", ip)
}

// FindClearCoreForLinMot returns the USBID of the ClearCore that owns the given LinMot IP.
func FindClearCoreForLinMot(ip string, cfg clearcore.Config) (string, bool) {
	for _, cc := range cfg.ClearCores {
		for _, lm := range cc.LinMots {
			if strings.TrimSpace(lm.IP) == strings.TrimSpace(ip) {
				return cc.USBID, true
			}
		}
	}
	return "", false
}
