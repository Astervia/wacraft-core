package synch_service

import (
	"sync"
	"sync/atomic"
	"time"
)

type memoryCounterEntry struct {
	value     atomic.Int64
	mu        sync.Mutex // protects expiresAt and hasTTL
	expiresAt time.Time
	hasTTL    bool
}

type MemoryCounter struct {
	entries sync.Map
}

func NewMemoryCounter() *MemoryCounter {
	return &MemoryCounter{}
}

func (c *MemoryCounter) Increment(key string, delta int64) (int64, error) {
	entry, _ := c.entries.LoadOrStore(key, &memoryCounterEntry{})
	e := entry.(*memoryCounterEntry)

	e.mu.Lock()
	if e.hasTTL && time.Now().After(e.expiresAt) {
		e.value.Store(0)
		e.hasTTL = false
	}
	e.mu.Unlock()

	return e.value.Add(delta), nil
}

func (c *MemoryCounter) Get(key string) (int64, error) {
	entry, ok := c.entries.Load(key)
	if !ok {
		return 0, nil
	}

	e := entry.(*memoryCounterEntry)
	e.mu.Lock()
	if e.hasTTL && time.Now().After(e.expiresAt) {
		e.mu.Unlock()
		c.entries.Delete(key)
		return 0, nil
	}
	e.mu.Unlock()

	return e.value.Load(), nil
}

func (c *MemoryCounter) SetTTL(key string, ttl time.Duration) error {
	entry, ok := c.entries.Load(key)
	if !ok {
		return nil
	}

	e := entry.(*memoryCounterEntry)
	e.mu.Lock()
	e.expiresAt = time.Now().Add(ttl)
	e.hasTTL = true
	e.mu.Unlock()
	return nil
}

func (c *MemoryCounter) Delete(key string) error {
	c.entries.Delete(key)
	return nil
}
