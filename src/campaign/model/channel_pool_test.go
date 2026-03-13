package campaign_model

import (
	"testing"

	synch_service "github.com/Astervia/wacraft-core/src/synch/service"
	websocket_model "github.com/Astervia/wacraft-core/src/websocket/model"
	"github.com/google/uuid"
)

// fakeClient creates a Client with a nil connection for testing pool logic.
func fakeClient(key string) websocket_model.Client[websocket_model.ClientID] {
	return websocket_model.Client[websocket_model.ClientID]{
		Data: websocket_model.ClientID{
			UserID: uuid.New(),
			ConnID: 0,
		},
	}
}

func TestChannelPool_CreateMemory(t *testing.T) {
	pool := CreateChannelPool()
	if pool == nil {
		t.Fatal("expected non-nil pool")
	}
}

func TestChannelPool_CreateDistributed(t *testing.T) {
	cache := synch_service.NewMemoryCache()
	pubsub := synch_service.NewMemoryPubSub()
	pool := CreateChannelPoolWithDistributed(cache, pubsub)
	if pool == nil {
		t.Fatal("expected non-nil pool")
	}
	if pool.cache != cache {
		t.Fatal("expected cache to be set")
	}
	if pool.pubsub != pubsub {
		t.Fatal("expected pubsub to be set")
	}
}

func TestChannelPool_AddUser_MemoryChannel(t *testing.T) {
	pool := CreateChannelPool()
	campaignID := uuid.New()

	ch := pool.AddUser(fakeClient("u1"), "u1", campaignID, nil)
	if ch == nil {
		t.Fatal("expected non-nil channel")
	}
	if ch.cache != nil || ch.pubsub != nil {
		t.Fatal("memory pool should produce channels without distributed primitives")
	}
}

func TestChannelPool_AddUser_DistributedChannel(t *testing.T) {
	cache := synch_service.NewMemoryCache()
	pubsub := synch_service.NewMemoryPubSub()
	pool := CreateChannelPoolWithDistributed(cache, pubsub)
	campaignID := uuid.New()

	ch := pool.AddUser(fakeClient("u1"), "u1", campaignID, nil)
	if ch == nil {
		t.Fatal("expected non-nil channel")
	}
	if ch.cache != cache {
		t.Fatal("distributed pool should produce channels with cache")
	}
	if ch.pubsub != pubsub {
		t.Fatal("distributed pool should produce channels with pubsub")
	}
}

func TestChannelPool_AddUser_SameCampaignReturnsExisting(t *testing.T) {
	pool := CreateChannelPool()
	campaignID := uuid.New()

	ch1 := pool.AddUser(fakeClient("u1"), "u1", campaignID, nil)
	ch2 := pool.AddUser(fakeClient("u2"), "u2", campaignID, nil)

	// Both calls should return a channel for the same campaign
	if ch1 == nil || ch2 == nil {
		t.Fatal("expected non-nil channels")
	}

	// The second user should be in the channel's clients
	pool.mu.Lock()
	channel := pool.channels[campaignID]
	pool.mu.Unlock()
	if len(channel.Clients) != 2 {
		t.Fatalf("expected 2 clients, got %d", len(channel.Clients))
	}
}

func TestChannelPool_RemoveUser(t *testing.T) {
	pool := CreateChannelPool()
	campaignID := uuid.New()

	pool.AddUser(fakeClient("u1"), "u1", campaignID, nil)
	pool.AddUser(fakeClient("u2"), "u2", campaignID, nil)

	pool.RemoveUser("u1", campaignID)

	pool.mu.Lock()
	channel, ok := pool.channels[campaignID]
	pool.mu.Unlock()
	if !ok {
		t.Fatal("channel should still exist with one user")
	}
	if len(channel.Clients) != 1 {
		t.Fatalf("expected 1 client, got %d", len(channel.Clients))
	}
}

func TestChannelPool_RemoveUser_DeletesEmptyChannel(t *testing.T) {
	pool := CreateChannelPool()
	campaignID := uuid.New()

	pool.AddUser(fakeClient("u1"), "u1", campaignID, nil)
	pool.RemoveUser("u1", campaignID)

	pool.mu.Lock()
	_, ok := pool.channels[campaignID]
	pool.mu.Unlock()
	if ok {
		t.Fatal("empty channel should have been removed")
	}
}

func TestChannelPool_RemoveUser_NonExistentCampaign(t *testing.T) {
	pool := CreateChannelPool()
	// Should not panic
	pool.RemoveUser("u1", uuid.New())
}

func TestChannelPool_DistributedSendingCrossChannel(t *testing.T) {
	cache := synch_service.NewMemoryCache()
	pubsub := synch_service.NewMemoryPubSub()
	pool := CreateChannelPoolWithDistributed(cache, pubsub)
	campaignID := uuid.New()

	ch := pool.AddUser(fakeClient("u1"), "u1", campaignID, nil)

	ch.SetSending(true)

	// A second pool (simulating another instance) sharing the same cache
	pool2 := CreateChannelPoolWithDistributed(cache, pubsub)
	ch2 := pool2.AddUser(fakeClient("u2"), "u2", campaignID, nil)

	if !ch2.IsSending() {
		t.Fatal("expected IsSending=true via shared cache on second pool")
	}
}
