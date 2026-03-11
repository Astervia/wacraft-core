package synch_service

import (
	"testing"
	"time"
)

func TestMemoryCache_SetGet(t *testing.T) {
	c := NewMemoryCache()

	data := []byte(`{"id": 1}`)
	if err := c.Set("k", data, 5*time.Second); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	got, ok, err := c.Get("k")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if !ok {
		t.Fatal("Get() returned not found")
	}
	if string(got) != string(data) {
		t.Errorf("Get() = %q, want %q", got, data)
	}
}

func TestMemoryCache_Miss(t *testing.T) {
	c := NewMemoryCache()

	got, ok, err := c.Get("nonexistent")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if ok {
		t.Fatal("Get(nonexistent) should return not found")
	}
	if got != nil {
		t.Errorf("Get(nonexistent) data = %v, want nil", got)
	}
}

func TestMemoryCache_TTLExpiry(t *testing.T) {
	c := NewMemoryCache()

	c.Set("k", []byte("data"), 50*time.Millisecond)

	// Before expiry
	_, ok, _ := c.Get("k")
	if !ok {
		t.Fatal("should be found before TTL")
	}

	time.Sleep(60 * time.Millisecond)

	_, ok, _ = c.Get("k")
	if ok {
		t.Fatal("should be expired after TTL")
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	c := NewMemoryCache()

	c.Set("k", []byte("data"), 5*time.Second)
	c.Delete("k")

	_, ok, _ := c.Get("k")
	if ok {
		t.Fatal("should be deleted")
	}
}

func TestMemoryCache_DeleteNonExistent(t *testing.T) {
	c := NewMemoryCache()

	if err := c.Delete("nonexistent"); err != nil {
		t.Fatalf("Delete(nonexistent) error: %v", err)
	}
}

func TestMemoryCache_Overwrite(t *testing.T) {
	c := NewMemoryCache()

	c.Set("k", []byte("v1"), 5*time.Second)
	c.Set("k", []byte("v2"), 5*time.Second)

	got, ok, _ := c.Get("k")
	if !ok {
		t.Fatal("should be found")
	}
	if string(got) != "v2" {
		t.Errorf("Get() = %q, want %q", got, "v2")
	}
}

func TestMemoryCache_Invalidate_TrailingWildcard(t *testing.T) {
	c := NewMemoryCache()

	c.Set("prefix:a", []byte("1"), 5*time.Second)
	c.Set("prefix:b", []byte("2"), 5*time.Second)
	c.Set("other:c", []byte("3"), 5*time.Second)

	if err := c.Invalidate("prefix:*"); err != nil {
		t.Fatalf("Invalidate() error: %v", err)
	}

	_, ok, _ := c.Get("prefix:a")
	if ok {
		t.Error("prefix:a should be invalidated")
	}

	_, ok, _ = c.Get("prefix:b")
	if ok {
		t.Error("prefix:b should be invalidated")
	}

	_, ok, _ = c.Get("other:c")
	if !ok {
		t.Error("other:c should NOT be invalidated")
	}
}

func TestMemoryCache_Invalidate_ExactMatch(t *testing.T) {
	c := NewMemoryCache()

	c.Set("exact", []byte("1"), 5*time.Second)
	c.Set("exactplus", []byte("2"), 5*time.Second)

	c.Invalidate("exact")

	_, ok, _ := c.Get("exact")
	if ok {
		t.Error("exact should be invalidated")
	}

	_, ok, _ = c.Get("exactplus")
	if !ok {
		t.Error("exactplus should NOT be invalidated by exact match")
	}
}

func TestMemoryCache_Invalidate_EmptyPattern(t *testing.T) {
	c := NewMemoryCache()

	c.Set("k", []byte("v"), 5*time.Second)

	c.Invalidate("")

	_, ok, _ := c.Get("k")
	if !ok {
		t.Error("empty pattern should not invalidate anything")
	}
}
