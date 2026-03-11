package synch_service

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMemoryLock_LockUnlock(t *testing.T) {
	lock := NewMemoryLock[string]()

	if err := lock.Lock("key1"); err != nil {
		t.Fatalf("Lock() error: %v", err)
	}
	if err := lock.Unlock("key1"); err != nil {
		t.Fatalf("Unlock() error: %v", err)
	}
}

func TestMemoryLock_ConcurrentSameKey(t *testing.T) {
	lock := NewMemoryLock[string]()
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lock.Lock("shared")
			// Non-atomic increment; mutex must protect it
			c := atomic.LoadInt64(&counter)
			time.Sleep(time.Microsecond) // Increase chance of race detection
			atomic.StoreInt64(&counter, c+1)
			lock.Unlock("shared")
		}()
	}

	wg.Wait()

	if counter != 100 {
		t.Errorf("counter = %d, want 100 (mutual exclusion failed)", counter)
	}
}

func TestMemoryLock_DifferentKeysParallel(t *testing.T) {
	lock := NewMemoryLock[string]()
	var ready sync.WaitGroup
	ready.Add(2)

	acquired := make(chan string, 2)

	go func() {
		lock.Lock("a")
		acquired <- "a"
		ready.Done()
		ready.Wait() // Wait for both to acquire
		lock.Unlock("a")
	}()

	go func() {
		lock.Lock("b")
		acquired <- "b"
		ready.Done()
		ready.Wait() // Wait for both to acquire
		lock.Unlock("b")
	}()

	// Both should acquire within a reasonable time
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()

	for i := 0; i < 2; i++ {
		select {
		case <-acquired:
		case <-timer.C:
			t.Fatal("timeout: different keys should not block each other")
		}
	}
}

func TestMemoryLock_SecondLockBlocks(t *testing.T) {
	lock := NewMemoryLock[string]()
	lock.Lock("key")

	blocked := make(chan struct{})
	go func() {
		lock.Lock("key")
		close(blocked)
		lock.Unlock("key")
	}()

	select {
	case <-blocked:
		t.Fatal("second Lock() should block while first holds it")
	case <-time.After(50 * time.Millisecond):
		// Expected: still blocked
	}

	lock.Unlock("key")

	select {
	case <-blocked:
		// Expected: unblocked after first unlock
	case <-time.After(2 * time.Second):
		t.Fatal("second Lock() should have acquired after first Unlock()")
	}
}

func TestMemoryLock_IntKey(t *testing.T) {
	lock := NewMemoryLock[int]()

	if err := lock.Lock(42); err != nil {
		t.Fatalf("Lock(int) error: %v", err)
	}
	if err := lock.Unlock(42); err != nil {
		t.Fatalf("Unlock(int) error: %v", err)
	}
}
