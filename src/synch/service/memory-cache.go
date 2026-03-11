package synch_service

import (
	"sync"
	"time"
)

type memoryCacheEntry struct {
	value     []byte
	expiresAt time.Time
}

type MemoryCache struct {
	entries sync.Map
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{}
}

func (c *MemoryCache) Get(key string) ([]byte, bool, error) {
	entry, ok := c.entries.Load(key)
	if !ok {
		return nil, false, nil
	}

	e := entry.(*memoryCacheEntry)
	if time.Now().After(e.expiresAt) {
		c.entries.Delete(key)
		return nil, false, nil
	}

	return e.value, true, nil
}

func (c *MemoryCache) Set(key string, value []byte, ttl time.Duration) error {
	c.entries.Store(key, &memoryCacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	})
	return nil
}

func (c *MemoryCache) Delete(key string) error {
	c.entries.Delete(key)
	return nil
}

func (c *MemoryCache) Invalidate(pattern string) error {
	c.entries.Range(func(key, value any) bool {
		k := key.(string)
		if matchSimplePattern(pattern, k) {
			c.entries.Delete(k)
		}
		return true
	})
	return nil
}

// matchSimplePattern matches a simple glob pattern with only trailing '*'.
// e.g., "prefix*" matches "prefix_anything".
func matchSimplePattern(pattern, s string) bool {
	if len(pattern) == 0 {
		return len(s) == 0
	}

	if pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(s) >= len(prefix) && s[:len(prefix)] == prefix
	}

	return pattern == s
}
