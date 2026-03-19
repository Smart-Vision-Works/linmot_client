//go:build linmot_instrumentation
// +build linmot_instrumentation

package client

import (
	"fmt"
	"testing"
)

// logResyncCountIfUDP logs the RTC resync count for UDP clients (hardware tests only).
// This is the real implementation that overrides the no-op stub in client_e2e_test.go
// when built with -tags linmot_instrumentation.
func logResyncCountIfUDP(t *testing.T, client *Client) {
	if *linmotMode == "udp" {
		t.Cleanup(func() {
			t.Logf("RTC resync count: %d", client.GetRTCResyncCountForTesting())
		})
	}
}

// dumpRouterTraceOnFailure logs the router trace ring buffer when a test fails.
func dumpRouterTraceOnFailure(t *testing.T, client *Client) {
	if *linmotMode != "udp" {
		return
	}
	t.Cleanup(func() {
		if !t.Failed() {
			return
		}
		reason := fmt.Sprintf("test failure: %s", t.Name())
		dump := client.requestManager.RouterTraceDumpForTesting(reason)
		if dump == "" {
			t.Logf("Router trace dump unavailable (debug disabled or empty).")
			return
		}
		t.Logf("Router trace dump:\n%s", dump)
	})
}
