package protocol_motion_control

import (
	"sync"
	"testing"
)

func TestMCCounter_ConcurrentNext_UniqueFirstCycle(t *testing.T) {
	counter := NewMCCounter()

	const numGoroutines = 4 // One for each value in range [1..4]

	var wg sync.WaitGroup
	results := make(chan uint8, numGoroutines)

	// Spawn 4 goroutines, each calling Next() once
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := counter.Next()
			results <- result
		}()
	}

	wg.Wait()
	close(results)

	// Collect all results
	seen := make(map[uint8]bool)
	var values []uint8
	for result := range results {
		values = append(values, result)
		if seen[result] {
			t.Errorf("Duplicate value found: %d", result)
		}
		seen[result] = true
		if result < MCCounterMin || result > MCCounterMax {
			t.Errorf("Value out of range: %d (should be %d-%d)", result, MCCounterMin, MCCounterMax)
		}
	}

	// Assert: exactly 4 unique values: {1, 2, 3, 4}
	if len(seen) != numGoroutines {
		t.Errorf("Expected %d unique values, got %d. Values: %v", numGoroutines, len(seen), values)
	}

	// Assert: all values in range [1..4]
	for i := MCCounterMin; i <= MCCounterMax; i++ {
		if !seen[i] {
			t.Errorf("Value %d was not seen in first cycle", i)
		}
	}
}

func TestMCCounter_ConcurrentNext_MultipleCycles(t *testing.T) {
	counter := NewMCCounter()

	const numGoroutines = 20
	const opsPerGoroutine = 10

	var wg sync.WaitGroup
	results := make(chan uint8, numGoroutines*opsPerGoroutine)

	// Spawn many goroutines calling Next() multiple times
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				result := counter.Next()
				results <- result
			}
		}()
	}

	wg.Wait()
	close(results)

	// Collect all results and verify they're in range
	count := 0
	for result := range results {
		count++
		if result < MCCounterMin || result > MCCounterMax {
			t.Errorf("Value out of range: %d (should be %d-%d)", result, MCCounterMin, MCCounterMax)
		}
	}

	if count != numGoroutines*opsPerGoroutine {
		t.Errorf("Expected %d results, got %d", numGoroutines*opsPerGoroutine, count)
	}
}
