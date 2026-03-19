package linmot

// Recovery state registry.
//
// When a drive enters recovery (after ROM writes in Setup, or after flash save
// in deploy), it is marked as "in recovery". While in this state the fault
// monitor skips the drive entirely — no GetStatus polls, no fault broadcasts,
// no wasted packets to a drive we know is unresponsive.
//
// The owner of the recovery (Setup, deploy) is responsible for clearing the
// state once waitForDriveRecovery confirms the drive is alive again.

import (
	"fmt"
	"sync"
)

var globalRecoveryState = &recoveryStateRegistry{
	recovering: make(map[string]bool),
}

type recoveryStateRegistry struct {
	mu         sync.RWMutex
	recovering map[string]bool
}

// EnterRecoveryState marks a drive as recovering. The fault monitor will skip
// this IP entirely until ExitRecoveryState is called.
func EnterRecoveryState(ip string) {
	globalRecoveryState.mu.Lock()
	defer globalRecoveryState.mu.Unlock()
	globalRecoveryState.recovering[ip] = true
	fmt.Printf("[LinMot %s] entered recovery state — fault monitor paused\n", ip)
}

// ExitRecoveryState clears the recovery flag. The fault monitor resumes
// polling this IP on its next tick.
func ExitRecoveryState(ip string) {
	globalRecoveryState.mu.Lock()
	defer globalRecoveryState.mu.Unlock()
	delete(globalRecoveryState.recovering, ip)
	fmt.Printf("[LinMot %s] exited recovery state — fault monitor resumed\n", ip)
}

// IsInRecovery returns true if the drive is currently recovering and should
// not be polled by the fault monitor.
func IsInRecovery(ip string) bool {
	globalRecoveryState.mu.RLock()
	defer globalRecoveryState.mu.RUnlock()
	return globalRecoveryState.recovering[ip]
}
