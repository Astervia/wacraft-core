package synch_redis

import (
	"testing"
	"time"
)

func TestRedisCache_CrossInstance(t *testing.T) {
	client := testRedisClient(t)

	cacheA := NewRedisCache(client)
	cacheB := NewRedisCache(client)

	data := []byte(`{"user":"test"}`)
	if err := cacheA.Set("shared-key", data, 5*time.Second); err != nil {
		t.Fatalf("Set from A error: %v", err)
	}

	got, ok, err := cacheB.Get("shared-key")
	if err != nil {
		t.Fatalf("Get from B error: %v", err)
	}
	if !ok {
		t.Fatal("Get from B should find the key set by A")
	}
	if string(got) != string(data) {
		t.Errorf("Get from B = %q, want %q", got, data)
	}
}

func TestRedisCache_TTL(t *testing.T) {
	client := testRedisClient(t)
	cache := NewRedisCache(client)

	cache.Set("ttl-key", []byte("value"), 100*time.Millisecond)

	_, ok, _ := cache.Get("ttl-key")
	if !ok {
		t.Fatal("should be found before TTL")
	}

	time.Sleep(150 * time.Millisecond)

	_, ok, _ = cache.Get("ttl-key")
	if ok {
		t.Fatal("should be expired after TTL")
	}
}

func TestRedisCache_Delete(t *testing.T) {
	client := testRedisClient(t)
	cache := NewRedisCache(client)

	cache.Set("del-key", []byte("value"), 5*time.Second)
	cache.Delete("del-key")

	_, ok, _ := cache.Get("del-key")
	if ok {
		t.Fatal("should be deleted")
	}
}

func TestRedisCache_Miss(t *testing.T) {
	client := testRedisClient(t)
	cache := NewRedisCache(client)

	data, ok, err := cache.Get("nonexistent")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if ok {
		t.Fatal("should be a miss")
	}
	if data != nil {
		t.Errorf("data = %v, want nil", data)
	}
}

func TestRedisCache_Overwrite(t *testing.T) {
	client := testRedisClient(t)
	cache := NewRedisCache(client)

	cache.Set("ow-key", []byte("v1"), 5*time.Second)
	cache.Set("ow-key", []byte("v2"), 5*time.Second)

	got, ok, _ := cache.Get("ow-key")
	if !ok {
		t.Fatal("should be found")
	}
	if string(got) != "v2" {
		t.Errorf("Get() = %q, want %q", got, "v2")
	}
}

func TestRedisCache_Invalidate(t *testing.T) {
	client := testRedisClient(t)
	cache := NewRedisCache(client)

	cache.Set("group:a", []byte("1"), 5*time.Second)
	cache.Set("group:b", []byte("2"), 5*time.Second)
	cache.Set("other:c", []byte("3"), 5*time.Second)

	if err := cache.Invalidate("group:*"); err != nil {
		t.Fatalf("Invalidate() error: %v", err)
	}

	_, ok, _ := cache.Get("group:a")
	if ok {
		t.Error("group:a should be invalidated")
	}

	_, ok, _ = cache.Get("group:b")
	if ok {
		t.Error("group:b should be invalidated")
	}

	_, ok, _ = cache.Get("other:c")
	if !ok {
		t.Error("other:c should NOT be invalidated")
	}
}
