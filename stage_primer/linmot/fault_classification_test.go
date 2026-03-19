package linmot

import (
	"testing"
	"time"

	"github.com/Smart-Vision-Works/staged_robot/client"

	config "stage_primer_config"

	"github.com/stretchr/testify/assert"
)

func TestClassifyFault_NilIsWarning(t *testing.T) {
	assert.Equal(t, FaultLevelWarning, classifyFault(nil))
}

func TestClassifyFault_WarningOnly(t *testing.T) {
	fault := &client.DriveFaultError{
		StatusWord:  0x0000,
		WarningWord: 0x0001,
		ErrorCode:   0,
	}
	assert.Equal(t, FaultLevelWarning, classifyFault(fault))
}

func TestClassifyFault_NonFatal(t *testing.T) {
	fault := &client.DriveFaultError{
		StatusWord: 0x0000, // bit 12 clear = not fatal
		ErrorCode:  0x0042,
	}
	assert.Equal(t, FaultLevelNonFatal, classifyFault(fault))
}

func TestClassifyFault_Fatal(t *testing.T) {
	fault := &client.DriveFaultError{
		StatusWord: 0x1000, // bit 12 set = fatal
		ErrorCode:  0x0042,
	}
	assert.Equal(t, FaultLevelFatal, classifyFault(fault))
}

func TestFaultBudget_WithinBudget(t *testing.T) {
	b := &faultBudget{windows: make(map[string][]time.Time)}
	ip := "10.0.0.1"

	for i := 0; i < autoRecoveryMaxFaults; i++ {
		assert.True(t, b.tryConsume(ip), "attempt %d should be within budget", i)
	}
}

func TestFaultBudget_Exceeded(t *testing.T) {
	b := &faultBudget{windows: make(map[string][]time.Time)}
	ip := "10.0.0.1"

	for i := 0; i < autoRecoveryMaxFaults; i++ {
		b.tryConsume(ip)
	}
	assert.False(t, b.tryConsume(ip), "should exceed budget")
}

func TestFaultBudget_Reset(t *testing.T) {
	b := &faultBudget{windows: make(map[string][]time.Time)}
	ip := "10.0.0.1"

	for i := 0; i < autoRecoveryMaxFaults; i++ {
		b.tryConsume(ip)
	}
	assert.False(t, b.tryConsume(ip))

	b.reset(ip)
	assert.True(t, b.tryConsume(ip), "should be within budget after reset")
}

func TestFaultBudget_IsolatedPerIP(t *testing.T) {
	b := &faultBudget{windows: make(map[string][]time.Time)}

	for i := 0; i < autoRecoveryMaxFaults; i++ {
		b.tryConsume("10.0.0.1")
	}
	assert.False(t, b.tryConsume("10.0.0.1"))
	assert.True(t, b.tryConsume("10.0.0.2"), "different IP should have its own budget")
}

func TestEscalationListener(t *testing.T) {
	var gotIP string
	var gotFatal bool

	remove := AddEscalationListener(func(ip string, fault *client.DriveFaultError, isFatal bool) {
		gotIP = ip
		gotFatal = isFatal
	})
	t.Cleanup(remove)

	fault := &client.DriveFaultError{StatusWord: 0x1000, ErrorCode: 0x01}
	broadcastEscalation("10.0.0.5", fault, true)

	assert.Equal(t, "10.0.0.5", gotIP)
	assert.True(t, gotFatal)
}

func TestEscalationListener_Unregister(t *testing.T) {
	calls := 0
	remove := AddEscalationListener(func(ip string, fault *client.DriveFaultError, isFatal bool) {
		calls++
	})
	t.Cleanup(remove)

	broadcastEscalation("10.0.0.1", nil, false)
	assert.Equal(t, 1, calls)

	remove()
	broadcastEscalation("10.0.0.1", nil, false)
	assert.Equal(t, 1, calls)
}

func TestHumanReadableFault_WithConfig(t *testing.T) {
	cfg := config.Config{
		ClearCores: []config.ClearCoreConfig{
			{
				USBID: "USB001",
				LinMots: []config.LinMotConfig{
					{IP: "10.0.0.1"},
					{IP: "10.0.0.2"},
				},
			},
		},
	}

	fault := &client.DriveFaultError{
		StatusWord: 0x0000,
		ErrorCode:  0x0042,
		ErrorText:  "Motor overheated",
	}

	msg := HumanReadableFault("10.0.0.2", fault, cfg)
	assert.Contains(t, msg, "Robot 1 / Stage 2")
	assert.Contains(t, msg, "10.0.0.2")
	assert.Contains(t, msg, "ERROR")
	assert.Contains(t, msg, "Motor overheated")
}

func TestHumanReadableFault_IncludesWarningTextAlongsideErrorText(t *testing.T) {
	cfg := config.Config{}
	fault := &client.DriveFaultError{
		ErrorCode:   0x0042,
		ErrorText:   "Motor overheated",
		WarningText: "Drive temperature high",
	}

	msg := HumanReadableFault("10.0.0.2", fault, cfg)
	assert.Contains(t, msg, "Motor overheated")
	assert.Contains(t, msg, "warning: Drive temperature high")
}

func TestHumanReadableFault_UnknownIP(t *testing.T) {
	cfg := config.Config{}
	fault := &client.DriveFaultError{ErrorCode: 0x01, ErrorText: "test"}
	msg := HumanReadableFault("99.99.99.99", fault, cfg)
	assert.Contains(t, msg, "LinMot 99.99.99.99")
}

func TestHumanReadableFault_Fatal(t *testing.T) {
	cfg := config.Config{}
	fault := &client.DriveFaultError{StatusWord: 0x1000, ErrorCode: 0x01, ErrorText: "Fatal motor error"}
	msg := HumanReadableFault("10.0.0.1", fault, cfg)
	assert.Contains(t, msg, "FATAL")
}

func TestFindClearCoreForLinMot(t *testing.T) {
	cfg := config.Config{
		ClearCores: []config.ClearCoreConfig{
			{
				USBID:   "USB-AAA",
				LinMots: []config.LinMotConfig{{IP: "10.0.0.1"}},
			},
			{
				USBID:   "USB-BBB",
				LinMots: []config.LinMotConfig{{IP: "10.0.0.2"}, {IP: "10.0.0.3"}},
			},
		},
	}

	usbID, found := FindClearCoreForLinMot("10.0.0.3", cfg)
	assert.True(t, found)
	assert.Equal(t, "USB-BBB", usbID)

	_, found = FindClearCoreForLinMot("99.99.99.99", cfg)
	assert.False(t, found)
}

// --- Production warn word regression guards (2026-03-12 log analysis) ---
// These values were observed in production after command table deployment and must
// always be classified as FaultLevelWarning (non-escalated). See _AGENTS/stage_primer_log_analysis_2026_03_12.md.

func TestClassifyFault_MotorNotHomedIsWarning(t *testing.T) {
	// warning_word=0x0080 → "Warning: Motor Not Homed Yet"
	// Expected after every MC restart triggered by command table deployment.
	fault := &client.DriveFaultError{
		StatusWord:  0x0000,
		WarningWord: 0x0080,
		ErrorCode:   0,
	}
	assert.Equal(t, FaultLevelWarning, classifyFault(fault))
}

func TestClassifyFault_SupplyVoltageLowIsWarning(t *testing.T) {
	// warning_word=0x0004 → "Motor Supply Voltage Reached Low Warn Limit"
	// Transient on power cycle — must remain FaultLevelWarning.
	fault := &client.DriveFaultError{
		StatusWord:  0x0000,
		WarningWord: 0x0004,
		ErrorCode:   0,
	}
	assert.Equal(t, FaultLevelWarning, classifyFault(fault))
}

func TestClassifyFault_MotorNotHomedAndVoltageLowIsWarning(t *testing.T) {
	// warning_word=0x0084 = 0x0080 | 0x0004 (both warnings simultaneously)
	// Observed at 21:00:44 in production. Must remain FaultLevelWarning.
	fault := &client.DriveFaultError{
		StatusWord:  0x0000,
		WarningWord: 0x0084,
		ErrorCode:   0,
	}
	assert.Equal(t, FaultLevelWarning, classifyFault(fault))
}

func TestBroadcastClearCoreFault(t *testing.T) {
	var gotIP string
	remove := AddFaultListener(func(ip string, err error) {
		gotIP = ip
	})
	defer remove()

	BroadcastClearCoreFault("USB-123", assert.AnError)
	assert.Equal(t, "USB-123", gotIP)
}
