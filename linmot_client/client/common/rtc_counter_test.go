package client_common

import (
	"sync"
	"testing"

	protocol_rtc "github.com/Smart-Vision-Works/staged_robot/protocol/rtc"
)

func TestCounterManagement(t *testing.T) {
	counter := NewRTCCounter()

	t.Run("Counter starts at 1", func(t *testing.T) {
		cnt := counter.Allocate()
		if cnt != 1 {
			t.Errorf("Counter = %d, want 1", cnt)
		}
	})

	t.Run("Allocate does not increment", func(t *testing.T) {
		cnt1 := counter.Allocate()
		cnt2 := counter.Allocate()
		if cnt1 != cnt2 {
			t.Errorf("Allocate should not increment: cnt1 = %d, cnt2 = %d", cnt1, cnt2)
		}
	})

	t.Run("Commit increments counter", func(t *testing.T) {
		cnt1 := counter.Allocate()
		counter.Commit()
		cnt2 := counter.Allocate()
		if cnt2 != cnt1+1 {
			t.Errorf("Commit should increment: cnt1 = %d, cnt2 = %d", cnt1, cnt2)
		}
	})

	t.Run("Counter wraps from 14 to 1", func(t *testing.T) {
		counter.SetForTesting(14)

		cnt := counter.Allocate()
		if cnt != 14 {
			t.Errorf("Counter = %d, want 14", cnt)
		}

		counter.Commit()
		cnt = counter.Allocate()
		if cnt != 1 {
			t.Errorf("Counter after wrap = %d, want 1", cnt)
		}
	})

	t.Run("Counter never becomes 0", func(t *testing.T) {
		counter.SetForTesting(protocol_rtc.CounterMax)

		counter.Commit()
		cnt := counter.Allocate()
		if cnt == 0 {
			t.Error("Counter should never be 0")
		}
		if cnt != 1 {
			t.Errorf("Counter after wrap = %d, want 1", cnt)
		}
	})

	t.Run("Counter range is 1-14", func(t *testing.T) {
		seen := make(map[uint8]bool)
		for i := 0; i < 20; i++ {
			cnt := counter.Allocate()
			if cnt < 1 || cnt > protocol_rtc.CounterMax {
				t.Errorf("Counter out of range: %d (should be 1-%d)", cnt, protocol_rtc.CounterMax)
			}
			seen[cnt] = true
			counter.Commit()
		}

		for i := uint8(1); i <= protocol_rtc.CounterMax; i++ {
			if !seen[i] {
				t.Errorf("Counter value %d was never seen", i)
			}
		}
	})
}

func TestCounterConcurrency(t *testing.T) {
	counter := NewRTCCounter()

	const numGoroutines = 100
	const opsPerGoroutine = 10

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*opsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				cnt := counter.Allocate()
				if cnt < 1 || cnt > protocol_rtc.CounterMax {
					errors <- &counterError{cnt: cnt}
				}
				counter.Commit()
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

type counterError struct {
	cnt uint8
}

func (e *counterError) Error() string {
	return "counter out of range"
}

func TestRtcCounter_ConcurrentNext_UniqueFirstCycle(t *testing.T) {
	counter := NewRTCCounter()

	const numGoroutines = 14 // One for each value in range [1..14]

	var wg sync.WaitGroup
	results := make(chan uint8, numGoroutines)

	// Spawn 14 goroutines, each calling Next() once
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
		if result < 1 || result > protocol_rtc.CounterMax {
			t.Errorf("Value out of range: %d (should be 1-%d)", result, protocol_rtc.CounterMax)
		}
	}

	// Assert: exactly 14 unique values
	if len(seen) != numGoroutines {
		t.Errorf("Expected %d unique values, got %d. Values: %v", numGoroutines, len(seen), values)
	}

	// Assert: all values in range [1..14]
	for i := uint8(1); i <= protocol_rtc.CounterMax; i++ {
		if !seen[i] {
			t.Errorf("Value %d was not seen in first cycle", i)
		}
	}
}

func TestRtcCounter_ConcurrentNext_MultipleCycles(t *testing.T) {
	counter := NewRTCCounter()

	const numGoroutines = 50
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
		if result < 1 || result > protocol_rtc.CounterMax {
			t.Errorf("Value out of range: %d (should be 1-%d)", result, protocol_rtc.CounterMax)
		}
	}

	if count != numGoroutines*opsPerGoroutine {
		t.Errorf("Expected %d results, got %d", numGoroutines*opsPerGoroutine, count)
	}
}
