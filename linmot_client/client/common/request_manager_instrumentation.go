//go:build linmot_instrumentation
// +build linmot_instrumentation

package client_common

import (
	"fmt"
	"strings"
	"time"
)

func (rm *RequestManager) RTCResyncCount() uint64 {
	return rm.rtcResyncCount.Load()
}

// RouterTraceDumpForTesting returns a formatted router trace dump string.
// Only available when built with -tags linmot_instrumentation.
func (rm *RequestManager) RouterTraceDumpForTesting(reason string) string {
	if !rm.debug.Load() {
		return ""
	}
	entries := rm.snapshotRouterTrace()
	if len(entries) == 0 {
		return ""
	}
	var builder strings.Builder
	fmt.Fprintf(&builder, "[ROUTER_TRACE] %s (last %d entries)\n", reason, len(entries))
	for _, entry := range entries {
		fmt.Fprintf(&builder, "[ROUTER_TRACE] %s %s\n", entry.timestamp.Format(time.RFC3339Nano), entry.entry)
	}
	return builder.String()
}
