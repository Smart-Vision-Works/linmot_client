package client_common

import (
	"sync"

	protocol_rtc "gsail-go/linmot/protocol/rtc"
)

// RtcCounter manages the RTC command counter for LinUDP V2 protocol communication.
type RtcCounter struct {
	mu  sync.Mutex
	cnt uint8
}

// NewRTCCounter creates a new RTC counter initialized to 1.
func NewRTCCounter() *RtcCounter {
	return &RtcCounter{
		cnt: 1,
	}
}

// Allocate returns the current counter value without incrementing.
// Invariant: Allocate() always returns a value in range 1-14, never 0.
func (r *RtcCounter) Allocate() uint8 {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cnt == 0 {
		r.cnt = 1
	}
	return r.cnt
}

// Commit increments the counter to the next value.
func (r *RtcCounter) Commit() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cnt++
	if r.cnt == 0 || r.cnt > protocol_rtc.CounterMax {
		r.cnt = 1
	}
}

// Next returns the current counter value and advances to the next value atomically.
// Returns a value in range 1-14, wrapping from 14 to 1.
// This replaces the non-atomic Allocate()+Commit() pattern.
// Invariant: Next() always returns a value in range 1-14, never 0.
func (r *RtcCounter) Next() uint8 {
	r.mu.Lock()
	defer r.mu.Unlock()
	current := r.cnt
	// Normalize current to ensure it's never 0 before returning
	if current == 0 {
		current = 1
	}
	r.cnt++
	if r.cnt == 0 || r.cnt > protocol_rtc.CounterMax {
		r.cnt = 1
	}
	return current
}

// Set sets the counter value, clamped to valid range (1-14).
// Intended for internal seeding; not for tests.
func (r *RtcCounter) Set(value uint8) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if value < 1 || value > protocol_rtc.CounterMax {
		if value == 0 {
			value = 1
		} else {
			value = protocol_rtc.CounterMax
		}
	}
	r.cnt = value
}

// SetForTesting sets the counter value for testing purposes only.
// This should only be used in test code.
func (r *RtcCounter) SetForTesting(value uint8) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if value < 1 || value > protocol_rtc.CounterMax {
		// Clamp to valid range
		if value == 0 {
			value = 1
		} else {
			value = protocol_rtc.CounterMax
		}
	}
	r.cnt = value
}
