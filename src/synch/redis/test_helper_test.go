package synch_redis

import (
	"os"
	"testing"
	"time"
)

func testRedisClient(t *testing.T) *Client {
	t.Helper()

	url := os.Getenv("REDIS_URL")
	if url == "" {
		t.Skip("REDIS_URL not set, skipping Redis integration test")
	}

	client, err := NewClient(Config{
		URL:       url,
		DB:        15, // Use DB 15 for tests to avoid conflicts
		KeyPrefix: "test:",
		LockTTL:   5 * time.Second,
		CacheTTL:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("failed to create Redis client: %v", err)
	}

	if err := client.PingWithTimeout(2 * time.Second); err != nil {
		t.Fatalf("failed to ping Redis: %v", err)
	}

	// Flush test DB before each test
	client.Redis().FlushDB(t.Context())

	t.Cleanup(func() {
		client.Redis().FlushDB(t.Context())
		client.Close()
	})

	return client
}
