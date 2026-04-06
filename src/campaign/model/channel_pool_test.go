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

// ─── GetOrCreateChannel / ReleaseChannel tests ───────────────────────────────

func TestChannelPool_GetOrCreateChannel_CreatesChannel(t *testing.T) {
	pool := CreateChannelPool()
	campaignID := uuid.New()

	ch := pool.GetOrCreateChannel(campaignID)
	if ch == nil {
		t.Fatal("expected non-nil channel")
	}

	pool.mu.Lock()
	_, exists := pool.channels[campaignID]
	refs := pool.workerRefs[campaignID]
	pool.mu.Unlock()

	if !exists {
		t.Fatal("channel should be in pool after GetOrCreateChannel")
	}
	if refs != 1 {
		t.Fatalf("workerRefs: got %d, want 1", refs)
	}
}

func TestChannelPool_GetOrCreateChannel_ReturnsExisting(t *testing.T) {
	pool := CreateChannelPool()
	campaignID := uuid.New()

	pool.AddUser(fakeClient("u1"), "u1", campaignID, nil)
	ch := pool.GetOrCreateChannel(campaignID)

	if ch == nil {
		t.Fatal("expected non-nil channel")
	}
	pool.mu.Lock()
	refs := pool.workerRefs[campaignID]
	pool.mu.Unlock()
	if refs != 1 {
		t.Fatalf("workerRefs after GetOrCreate on existing: got %d, want 1", refs)
	}
}

func TestChannelPool_ReleaseChannel_DecrementsRef(t *testing.T) {
	pool := CreateChannelPool()
	campaignID := uuid.New()

	pool.GetOrCreateChannel(campaignID)
	pool.ReleaseChannel(campaignID)

	pool.mu.Lock()
	refs := pool.workerRefs[campaignID]
	_, channelExists := pool.channels[campaignID]
	pool.mu.Unlock()

	if refs != 0 {
		t.Fatalf("workerRefs after release: got %d, want 0", refs)
	}
	if channelExists {
		t.Fatal("channel should be removed after release with no clients")
	}
}

func TestChannelPool_ReleaseChannel_KeepsAliveWhileClientsPresent(t *testing.T) {
	pool := CreateChannelPool()
	campaignID := uuid.New()

	pool.AddUser(fakeClient("u1"), "u1", campaignID, nil)
	pool.GetOrCreateChannel(campaignID)
	pool.ReleaseChannel(campaignID)

	pool.mu.Lock()
	_, channelExists := pool.channels[campaignID]
	pool.mu.Unlock()

	if !channelExists {
		t.Fatal("channel should remain in pool while WebSocket clients are connected")
	}
}

func TestChannelPool_RemoveUser_KeepsChannelWhileWorkerHoldsRef(t *testing.T) {
	pool := CreateChannelPool()
	campaignID := uuid.New()

	pool.AddUser(fakeClient("u1"), "u1", campaignID, nil)
	pool.GetOrCreateChannel(campaignID) // worker holds ref

	pool.RemoveUser("u1", campaignID) // last client disconnects

	pool.mu.Lock()
	_, channelExists := pool.channels[campaignID]
	pool.mu.Unlock()

	if !channelExists {
		t.Fatal("channel should remain in pool while worker holds a reference")
	}
}

func TestChannelPool_FullLifecycle_WorkerAndClient(t *testing.T) {
	pool := CreateChannelPool()
	campaignID := uuid.New()

	// Scheduler acquires hold first.
	pool.GetOrCreateChannel(campaignID)

	// Client connects while scheduler is running.
	pool.AddUser(fakeClient("u1"), "u1", campaignID, nil)

	// Client disconnects.
	pool.RemoveUser("u1", campaignID)

	// Channel still alive (worker holds ref).
	pool.mu.Lock()
	_, alive := pool.channels[campaignID]
	pool.mu.Unlock()
	if !alive {
		t.Fatal("channel should be alive while worker still holds ref")
	}

	// Scheduler finishes and releases.
	pool.ReleaseChannel(campaignID)

	// Channel should now be cleaned up.
	pool.mu.Lock()
	_, alive = pool.channels[campaignID]
	pool.mu.Unlock()
	if alive {
		t.Fatal("channel should be removed after worker releases and no clients remain")
	}
}

// TestChannelPool_WorkerBroadcastReachesWebSocketClient verifies that the
// Clients map is shared between the channel returned by GetOrCreateChannel and
// the channel returned by a subsequent AddUser call.
//
// In memory mode, CampaignChannel is stored as a value in the pool map; however,
// map fields inside the struct are reference types and therefore shared across
// copies. This ensures that BroadcastProgress on the scheduler's channel
// iterates the same Clients map that AddUser populated.
func TestChannelPool_WorkerBroadcastReachesWebSocketClient(t *testing.T) {
	pool := CreateChannelPool()
	campaignID := uuid.New()

	// Scheduler acquires hold before any client connects.
	workerCh := pool.GetOrCreateChannel(campaignID)

	// A WebSocket client joins.
	pool.AddUser(fakeClient("u1"), "u1", campaignID, nil)

	// The Clients map in the worker channel must contain the client added via
	// AddUser, because both copies of the struct share the same underlying map.
	if len(workerCh.Clients) != 1 {
		t.Fatalf("worker channel Clients: got %d entries, want 1 — client added via AddUser should be visible to BroadcastProgress", len(workerCh.Clients))
	}
}
