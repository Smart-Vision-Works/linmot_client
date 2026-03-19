package protocol_motion_control

import "sync"

// MCCounter manages Motion Control command counter allocation.
// The counter range is 1-4 and wraps to 1 after reaching 4.
// This is independent from the RTC Command Count (1-14 range).
//
// Reference: LINUDP_LIBRARY_ANALYSIS.md - Motion Control Counter System
type MCCounter struct {
	mu      sync.Mutex
	current uint8 // Current counter value (1-4)
	next    uint8 // Next counter to allocate (1-4)
}

// NewMCCounter creates a new MC counter starting at 1.
func NewMCCounter() *MCCounter {
	return &MCCounter{
		current: 0, // 0 indicates no counter has been used yet
		next:    1, // First counter to allocate is 1
	}
}

// Allocate reserves the next counter value.
// Returns a value in range 1-4.
// Call Commit() to finalize the allocation.
func (c *MCCounter) Allocate() uint8 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.next
}

// Get returns the current counter value without allocating a new one.
// Returns 0 if no counter has been committed yet.
func (c *MCCounter) Get() uint8 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.current
}

// Commit finalizes the allocation of the next counter.
// Advances the counter (1 → 2 → 3 → 4 → 1).
func (c *MCCounter) Commit() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.current = c.next

	// Advance to next counter (wrap at 4)
	c.next++
	if c.next > MCCounterMax {
		c.next = MCCounterMin
	}
}

// Next returns the next counter value and advances atomically.
// Returns a value in range 1-4, wrapping from 4 to 1.
// This replaces the non-atomic Allocate()+Commit() pattern.
func (c *MCCounter) Next() uint8 {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := c.next
	c.current = c.next
	c.next++
	if c.next > MCCounterMax {
		c.next = MCCounterMin
	}
	return result
}

// Reset resets the counter to initial state (next = 1, current = 0).
// This is useful for testing or reconnection scenarios.
func (c *MCCounter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.current = 0
	c.next = 1
}
