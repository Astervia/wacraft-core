package synch_redis

import (
	"sync"
	"testing"
	"time"
)

func TestRedisCounter_CrossInstance(t *testing.T) {
	client := testRedisClient(t)

	counterA := NewRedisCounter(client)
	counterB := NewRedisCounter(client)

	counterA.Increment("shared", 5)
	counterB.Increment("shared", 5)

	// Both instances should see the aggregated value
	valA, err := counterA.Get("shared")
	if err != nil {
		t.Fatalf("Get from A error: %v", err)
	}
	if valA != 10 {
		t.Errorf("Get from A = %d, want 10", valA)
	}

	valB, err := counterB.Get("shared")
	if err != nil {
		t.Fatalf("Get from B error: %v", err)
	}
	if valB != 10 {
		t.Errorf("Get from B = %d, want 10", valB)
	}
}

func TestRedisCounter_ConcurrentIncrement(t *testing.T) {
	client := testRedisClient(t)
	counter := NewRedisCounter(client)

	var wg sync.WaitGroup
	total := 100

	wg.Add(total)
	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()
			counter.Increment("concurrent", 1)
		}()
	}

	wg.Wait()

	val, _ := counter.Get("concurrent")
	if val != int64(total) {
		t.Errorf("Get() = %d, want %d", val, total)
	}
}

func TestRedisCounter_TTL(t *testing.T) {
	client := testRedisClient(t)
	counter := NewRedisCounter(client)

	counter.Increment("ttl-key", 42)
	counter.SetTTL("ttl-key", 1*time.Second) // Redis minimum EXPIRE is 1s

	val, _ := counter.Get("ttl-key")
	if val != 42 {
		t.Errorf("before TTL: Get() = %d, want 42", val)
	}

	time.Sleep(1100 * time.Millisecond)

	val, _ = counter.Get("ttl-key")
	if val != 0 {
		t.Errorf("after TTL: Get() = %d, want 0", val)
	}
}

func TestRedisCounter_Delete(t *testing.T) {
	client := testRedisClient(t)
	counter := NewRedisCounter(client)

	counter.Increment("del-key", 10)
	counter.Delete("del-key")

	val, _ := counter.Get("del-key")
	if val != 0 {
		t.Errorf("after Delete: Get() = %d, want 0", val)
	}
}

func TestRedisCounter_GetNonExistent(t *testing.T) {
	client := testRedisClient(t)
	counter := NewRedisCounter(client)

	val, err := counter.Get("nonexistent")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if val != 0 {
		t.Errorf("Get(nonexistent) = %d, want 0", val)
	}
}
