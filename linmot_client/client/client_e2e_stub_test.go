//go:build !linmot_instrumentation
// +build !linmot_instrumentation

package client

import "testing"

// logResyncCountIfUDP is a no-op stub when built without -tags linmot_instrumentation.
// The real implementation is in client_e2e_instrumentation_test.go.
func logResyncCountIfUDP(t *testing.T, client *Client) {
	// No-op: real implementation is in tagged file
}

// dumpRouterTraceOnFailure is a no-op stub when built without -tags linmot_instrumentation.
// The real implementation is in client_e2e_instrumentation_test.go.
func dumpRouterTraceOnFailure(t *testing.T, client *Client) {
	// No-op: real implementation is in tagged file
}
