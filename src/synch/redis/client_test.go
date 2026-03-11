package synch_redis

import (
	"testing"
	"time"
)

func TestNewClient_ParsesValidURL(t *testing.T) {
	client, err := NewClient(Config{
		URL:       "redis://localhost:6379",
		DB:        0,
		KeyPrefix: "test:",
		LockTTL:   30 * time.Second,
		CacheTTL:  5 * time.Minute,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer client.Close()

	if client.rdb == nil {
		t.Fatal("expected redis client to be initialized")
	}
}

func TestNewClient_InvalidURL(t *testing.T) {
	_, err := NewClient(Config{
		URL: "not-a-valid-url://???",
	})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestPrefixKey(t *testing.T) {
	client, err := NewClient(Config{
		URL:       "redis://localhost:6379",
		KeyPrefix: "wacraft:",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer client.Close()

	got := client.PrefixKey("lock:abc")
	want := "wacraft:lock:abc"
	if got != want {
		t.Errorf("PrefixKey() = %q, want %q", got, want)
	}
}

func TestPrefixKey_EmptyPrefix(t *testing.T) {
	client, err := NewClient(Config{
		URL:       "redis://localhost:6379",
		KeyPrefix: "",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer client.Close()

	got := client.PrefixKey("mykey")
	want := "mykey"
	if got != want {
		t.Errorf("PrefixKey() = %q, want %q", got, want)
	}
}

func TestConfig_Accessors(t *testing.T) {
	cfg := Config{
		URL:       "redis://localhost:6379",
		Password:  "secret",
		DB:        2,
		KeyPrefix: "pfx:",
		LockTTL:   10 * time.Second,
		CacheTTL:  1 * time.Minute,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer client.Close()

	got := client.Config()
	if got.URL != cfg.URL || got.DB != cfg.DB || got.KeyPrefix != cfg.KeyPrefix {
		t.Errorf("Config() returned unexpected values: %+v", got)
	}

	if client.Redis() == nil {
		t.Fatal("Redis() returned nil")
	}
}

func TestPingWithTimeout_NoRedis(t *testing.T) {
	// Connect to a port that is almost certainly not running Redis
	client, err := NewClient(Config{
		URL: "redis://localhost:16399",
	})
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}
	defer client.Close()

	err = client.PingWithTimeout(200 * time.Millisecond)
	if err == nil {
		t.Fatal("expected error pinging non-existent Redis")
	}
}
