package synch

import (
	"testing"
)

func TestFactory_MemoryBackend(t *testing.T) {
	f := NewFactory(BackendMemory, nil)

	if f.Backend() != BackendMemory {
		t.Errorf("Backend() = %q, want %q", f.Backend(), BackendMemory)
	}

	if f.RedisClient() != nil {
		t.Error("RedisClient() should be nil for memory backend")
	}
}

func TestFactory_MemoryLock(t *testing.T) {
	f := NewFactory(BackendMemory, nil)

	lock := NewLock[string](f)
	if lock == nil {
		t.Fatal("NewLock() returned nil")
	}

	// Verify it works
	if err := lock.Lock("test"); err != nil {
		t.Fatalf("Lock() error: %v", err)
	}
	if err := lock.Unlock("test"); err != nil {
		t.Fatalf("Unlock() error: %v", err)
	}
}

func TestFactory_MemoryLockIntKey(t *testing.T) {
	f := NewFactory(BackendMemory, nil)

	lock := NewLock[int](f)
	if lock == nil {
		t.Fatal("NewLock[int]() returned nil")
	}

	if err := lock.Lock(42); err != nil {
		t.Fatalf("Lock() error: %v", err)
	}
	if err := lock.Unlock(42); err != nil {
		t.Fatalf("Unlock() error: %v", err)
	}
}

func TestFactory_MemoryPubSub(t *testing.T) {
	f := NewFactory(BackendMemory, nil)

	ps := f.NewPubSub()
	if ps == nil {
		t.Fatal("NewPubSub() returned nil")
	}

	sub, err := ps.Subscribe("test-channel")
	if err != nil {
		t.Fatalf("Subscribe() error: %v", err)
	}
	defer sub.Unsubscribe()

	if err := ps.Publish("test-channel", []byte("msg")); err != nil {
		t.Fatalf("Publish() error: %v", err)
	}

	got := <-sub.Channel()
	if string(got) != "msg" {
		t.Errorf("received %q, want %q", got, "msg")
	}
}

func TestFactory_MemoryCounter(t *testing.T) {
	f := NewFactory(BackendMemory, nil)

	counter := f.NewCounter()
	if counter == nil {
		t.Fatal("NewCounter() returned nil")
	}

	result, err := counter.Increment("k", 5)
	if err != nil {
		t.Fatalf("Increment() error: %v", err)
	}
	if result != 5 {
		t.Errorf("Increment() = %d, want 5", result)
	}

	val, err := counter.Get("k")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if val != 5 {
		t.Errorf("Get() = %d, want 5", val)
	}
}

func TestFactory_MemoryCache(t *testing.T) {
	f := NewFactory(BackendMemory, nil)

	cache := f.NewCache()
	if cache == nil {
		t.Fatal("NewCache() returned nil")
	}

	if err := cache.Set("k", []byte("v"), 5000000000); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	data, ok, err := cache.Get("k")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if !ok {
		t.Fatal("Get() should find the key")
	}
	if string(data) != "v" {
		t.Errorf("Get() = %q, want %q", data, "v")
	}
}

func TestFactory_BackendConstants(t *testing.T) {
	if BackendMemory != "memory" {
		t.Errorf("BackendMemory = %q, want %q", BackendMemory, "memory")
	}
	if BackendRedis != "redis" {
		t.Errorf("BackendRedis = %q, want %q", BackendRedis, "redis")
	}
}
