package synch_redis

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRedisLock_LockUnlock(t *testing.T) {
	client := testRedisClient(t)
	lock := NewRedisLock[string](client)

	if err := lock.Lock("key1"); err != nil {
		t.Fatalf("Lock() error: %v", err)
	}

	// Key should exist in Redis
	exists := client.Redis().Exists(t.Context(), client.PrefixKey("lock:key1")).Val()
	if exists != 1 {
		t.Error("lock key should exist in Redis after Lock()")
	}

	if err := lock.Unlock("key1"); err != nil {
		t.Fatalf("Unlock() error: %v", err)
	}

	// Key should be gone after unlock
	exists = client.Redis().Exists(t.Context(), client.PrefixKey("lock:key1")).Val()
	if exists != 0 {
		t.Error("lock key should be deleted after Unlock()")
	}
}

func TestRedisLock_MutualExclusion(t *testing.T) {
	client := testRedisClient(t)

	// Two lock instances simulating two app instances
	lockA := NewRedisLock[string](client)
	lockB := NewRedisLock[string](client)

	var counter int64
	var wg sync.WaitGroup

	iterations := 50
	wg.Add(iterations * 2)

	increment := func(lock *RedisLock[string]) {
		defer wg.Done()
		if err := lock.Lock("shared"); err != nil {
			t.Errorf("Lock() error: %v", err)
			return
		}
		// Non-atomic read-modify-write; the lock must protect it
		c := atomic.LoadInt64(&counter)
		time.Sleep(time.Microsecond)
		atomic.StoreInt64(&counter, c+1)
		lock.Unlock("shared")
	}

	for i := 0; i < iterations; i++ {
		go increment(lockA)
		go increment(lockB)
	}

	wg.Wait()

	if counter != int64(iterations*2) {
		t.Errorf("counter = %d, want %d (mutual exclusion failed)", counter, iterations*2)
	}
}

func TestRedisLock_TTLExpiry(t *testing.T) {
	client := testRedisClient(t)

	// Lock with very short TTL
	lock := NewRedisLock[string](client)
	lock.ttl = 100 * time.Millisecond

	if err := lock.Lock("expire-key"); err != nil {
		t.Fatalf("Lock() error: %v", err)
	}

	// Don't unlock — simulate crash
	time.Sleep(150 * time.Millisecond)

	// Another lock instance should be able to acquire
	lock2 := NewRedisLock[string](client)
	lock2.ttl = 5 * time.Second

	done := make(chan error, 1)
	go func() {
		done <- lock2.Lock("expire-key")
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("second Lock() error: %v", err)
		}
		lock2.Unlock("expire-key")
	case <-time.After(2 * time.Second):
		t.Fatal("second Lock() should have acquired after TTL expiry")
	}
}

func TestRedisLock_UnlockOnlyOwner(t *testing.T) {
	client := testRedisClient(t)

	lockA := NewRedisLock[string](client)
	lockB := NewRedisLock[string](client)

	lockA.Lock("owner-test")

	// lockB tries to unlock a key it doesn't own
	lockB.Unlock("owner-test")

	// Key should still exist (lockA still holds it)
	exists := client.Redis().Exists(t.Context(), client.PrefixKey("lock:owner-test")).Val()
	if exists != 1 {
		t.Error("lock should still exist after non-owner unlock attempt")
	}

	lockA.Unlock("owner-test")
}

func TestRedisLock_ConcurrentHighContention(t *testing.T) {
	client := testRedisClient(t)

	locks := make([]*RedisLock[string], 5)
	for i := range locks {
		locks[i] = NewRedisLock[string](client)
	}

	var counter int64
	var wg sync.WaitGroup
	total := 50

	wg.Add(total)
	for i := 0; i < total; i++ {
		lock := locks[i%len(locks)]
		go func() {
			defer wg.Done()
			lock.Lock("contention")
			c := atomic.LoadInt64(&counter)
			time.Sleep(time.Microsecond)
			atomic.StoreInt64(&counter, c+1)
			lock.Unlock("contention")
		}()
	}

	wg.Wait()

	if counter != int64(total) {
		t.Errorf("counter = %d, want %d", counter, total)
	}
}
