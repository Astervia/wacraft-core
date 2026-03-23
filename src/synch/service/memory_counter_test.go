package synch_service

import (
	"sync"
	"testing"
	"time"
)

func TestMemoryCounter_Increment(t *testing.T) {
	c := NewMemoryCounter()

	for i := 0; i < 3; i++ {
		c.Increment("k", 1)
	}

	val, err := c.Get("k")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if val != 3 {
		t.Errorf("Get() = %d, want 3", val)
	}
}

func TestMemoryCounter_IncrementDelta(t *testing.T) {
	c := NewMemoryCounter()

	result, err := c.Increment("k", 10)
	if err != nil {
		t.Fatalf("Increment() error: %v", err)
	}
	if result != 10 {
		t.Errorf("Increment() returned %d, want 10", result)
	}

	result, err = c.Increment("k", 5)
	if err != nil {
		t.Fatalf("Increment() error: %v", err)
	}
	if result != 15 {
		t.Errorf("Increment() returned %d, want 15", result)
	}
}

func TestMemoryCounter_ConcurrentIncrement(t *testing.T) {
	c := NewMemoryCounter()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Increment("k", 1)
		}()
	}

	wg.Wait()

	val, _ := c.Get("k")
	if val != 100 {
		t.Errorf("Get() = %d, want 100", val)
	}
}

func TestMemoryCounter_TTLExpiry(t *testing.T) {
	c := NewMemoryCounter()

	c.Increment("k", 5)
	c.SetTTL("k", 50*time.Millisecond)

	// Before expiry
	val, _ := c.Get("k")
	if val != 5 {
		t.Errorf("before TTL: Get() = %d, want 5", val)
	}

	time.Sleep(60 * time.Millisecond)

	val, _ = c.Get("k")
	if val != 0 {
		t.Errorf("after TTL: Get() = %d, want 0", val)
	}
}

func TestMemoryCounter_Delete(t *testing.T) {
	c := NewMemoryCounter()

	c.Increment("k", 10)
	c.Delete("k")

	val, _ := c.Get("k")
	if val != 0 {
		t.Errorf("after Delete: Get() = %d, want 0", val)
	}
}

func TestMemoryCounter_GetNonExistent(t *testing.T) {
	c := NewMemoryCounter()

	val, err := c.Get("nonexistent")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if val != 0 {
		t.Errorf("Get(nonexistent) = %d, want 0", val)
	}
}

func TestMemoryCounter_SetTTLNonExistent(t *testing.T) {
	c := NewMemoryCounter()

	// Should not error on non-existent key
	if err := c.SetTTL("nonexistent", time.Second); err != nil {
		t.Fatalf("SetTTL() error: %v", err)
	}
}

func TestMemoryCounter_DeleteNonExistent(t *testing.T) {
	c := NewMemoryCounter()

	if err := c.Delete("nonexistent"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
}

func TestMemoryCounter_IncrementAfterTTLExpiry(t *testing.T) {
	c := NewMemoryCounter()

	c.Increment("k", 10)
	c.SetTTL("k", 50*time.Millisecond)

	time.Sleep(60 * time.Millisecond)

	// Increment after expiry should reset and start from delta
	result, _ := c.Increment("k", 3)
	if result != 3 {
		t.Errorf("Increment after TTL = %d, want 3", result)
	}
}

// TestMemoryCounter_ConcurrentIncrementAndSetTTL reproduces the cold-start data race
// that caused spurious 429s on the first request. On first use, one goroutine does
// the initial Increment (and then calls SetTTL), while concurrent goroutines also
// call Increment on the same key. Without the per-entry mutex, hasTTL and expiresAt
// could be partially written while another goroutine reads them, resetting the counter.
func TestMemoryCounter_ConcurrentIncrementAndSetTTL(t *testing.T) {
	const goroutines = 50
	const ttl = 5 * time.Second

	c := NewMemoryCounter()
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()
			val, _ := c.Increment("k", 1)
			// Simulate what ThroughputCounter does: the goroutine that gets
			// val==1 (first increment) sets the TTL, racing all others.
			if val == 1 {
				c.SetTTL("k", ttl)
			}
		}(i)
	}

	wg.Wait()

	// All goroutines incremented by 1. Counter must be exactly goroutines.
	// A race resetting the counter mid-flight would make this less.
	val, _ := c.Get("k")
	if val != goroutines {
		t.Errorf("counter = %d after %d concurrent increments, want %d (counter was reset by a race)", val, goroutines, goroutines)
	}
}

// TestMemoryCounter_ConcurrentSetTTLAndIncrement checks the reverse ordering:
// SetTTL is called just before a burst of Increments begins (simulates a window
// boundary where TTL is already set on entry creation and goroutines pile in).
func TestMemoryCounter_ConcurrentSetTTLAndIncrement(t *testing.T) {
	const goroutines = 50

	c := NewMemoryCounter()
	c.Increment("k", 1)
	c.SetTTL("k", 5*time.Second)

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			c.Increment("k", 1)
		}()
	}

	wg.Wait()

	val, _ := c.Get("k")
	if val != goroutines+1 {
		t.Errorf("counter = %d, want %d (one seed + %d concurrent)", val, goroutines+1, goroutines)
	}
}
