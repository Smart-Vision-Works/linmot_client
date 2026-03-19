package transport

import (
	"math/rand"
	"time"
)

// MockTransportConfig configures packet loss simulation behavior for mock transport.
// This is used for testing UDP reliability scenarios.
type MockTransportConfig struct {
	// DropProbability is the probability (0.0-1.0) of dropping packets.
	// 0.0 = no drops, 1.0 = all packets dropped.
	DropProbability float64

	// DelayMin and DelayMax define the random delay range for responses.
	// If both are 0, responses are sent immediately.
	DelayMin time.Duration
	DelayMax time.Duration

	// DuplicateProbability is the probability (0.0-1.0) of duplicating packets.
	// When a packet is duplicated, the original and a copy are both sent.
	DuplicateProbability float64

	// StaleResponseProbability is the probability (0.0-1.0) of returning stale responses.
	// Stale responses have wrong counter values (for RTC) or arrive when no request is pending.
	StaleResponseProbability float64

	// EnablePacketLoss enables all packet loss simulation features.
	// If false, all packets are delivered normally regardless of other settings.
	EnablePacketLoss bool

	// rng is the random number generator for this config.
	// Initialized on first use.
	rng *rand.Rand
}

// DefaultMockTransportConfig returns a config with no packet loss (normal behavior).
func DefaultMockTransportConfig() *MockTransportConfig {
	return &MockTransportConfig{
		EnablePacketLoss: false,
	}
}

// WithPacketLoss returns a config with packet loss enabled.
// This is a convenience method for creating a config with common loss scenarios.
func WithPacketLoss(dropProb, duplicateProb float64) *MockTransportConfig {
	return &MockTransportConfig{
		DropProbability:      dropProb,
		DuplicateProbability: duplicateProb,
		EnablePacketLoss:     true,
	}
}

// shouldDrop returns true if a packet should be dropped based on DropProbability.
func (c *MockTransportConfig) shouldDrop() bool {
	if !c.EnablePacketLoss {
		return false
	}
	return c.random() < c.DropProbability
}

// shouldDuplicate returns true if a packet should be duplicated based on DuplicateProbability.
func (c *MockTransportConfig) shouldDuplicate() bool {
	if !c.EnablePacketLoss {
		return false
	}
	return c.random() < c.DuplicateProbability
}

// randomDelay returns a random delay between DelayMin and DelayMax.
func (c *MockTransportConfig) randomDelay() time.Duration {
	if !c.EnablePacketLoss || c.DelayMin == 0 && c.DelayMax == 0 {
		return 0
	}
	if c.DelayMin == c.DelayMax {
		return c.DelayMin
	}
	diff := c.DelayMax - c.DelayMin
	return c.DelayMin + time.Duration(float64(diff)*c.random())
}

// random returns a random float64 in [0.0, 1.0).
func (c *MockTransportConfig) random() float64 {
	if c.rng == nil {
		c.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return c.rng.Float64()
}

// Reset resets the random number generator for test isolation.
func (c *MockTransportConfig) Reset() {
	c.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}
