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
