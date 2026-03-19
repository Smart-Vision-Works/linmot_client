//go:build linmot_instrumentation
// +build linmot_instrumentation

package client

// GetRTCResyncCountForTesting returns the total number of RTC resyncs that have occurred.
// This function is only available when built with -tags linmot_instrumentation.
// It is NOT compiled into production binaries or default test builds.
func (c *Client) GetRTCResyncCountForTesting() uint64 {
	return c.requestManager.RTCResyncCount()
}
