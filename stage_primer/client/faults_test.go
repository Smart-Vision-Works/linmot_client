package client

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pb "primer/proto"
)

func TestFormatFaultNotification_IncludesProbeAndWarningDetails(t *testing.T) {
	text := formatFaultNotification(&pb.FaultNotification{
		ErrorText:   "Position Lag Error",
		WarningWord: 0x0004,
		WarningText: "Drive temperature high",
		ProbeError:  "probe timeout",
		StatusWord:  0x12AB,
		StateVar:    0x0200,
	})

	assert.Contains(t, text, `error="Position Lag Error"`)
	assert.Contains(t, text, `warning="Drive temperature high"`)
	assert.Contains(t, text, "warning_word=0x4")
	assert.Contains(t, text, `probe_error="probe timeout"`)
	assert.Contains(t, text, "status_word=0x12AB")
	assert.Contains(t, text, "state_var=0x200")
}

func TestFormatFaultNotification_UsesErrorCodeWhenTextMissing(t *testing.T) {
	text := formatFaultNotification(&pb.FaultNotification{
		ErrorCode: 0x20,
	})

	assert.Equal(t, "error_code=0x20", text)
}

func TestFormatFaultNotification_UnknownWhenNoFields(t *testing.T) {
	assert.Equal(t, "Unknown fault", formatFaultNotification(&pb.FaultNotification{}))
	assert.Equal(t, "Unknown fault", formatFaultNotification(nil))
}
