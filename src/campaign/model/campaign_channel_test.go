package campaign_model

import (
	"context"
	"sync"
	"testing"
	"time"

	synch_service "github.com/Astervia/wacraft-core/src/synch/service"
)

// ─── SetSending / IsSending (memory-only) ────────────────────────────────────

func TestCampaignChannel_SetSending_Memory(t *testing.T) {
	ch := CreateCampaignChannel(nil)

	if ch.IsSending() {
		t.Fatal("expected IsSending=false on new channel")
	}

	ch.SetSending(true)
	if !ch.IsSending() {
		t.Fatal("expected IsSending=true after SetSending(true)")
	}

	ch.SetSending(false)
	if ch.IsSending() {
		t.Fatal("expected IsSending=false after SetSending(false)")
	}
}

// ─── SetSending / IsSending (distributed cache) ─────────────────────────────

func TestCampaignChannel_SetSending_Distributed(t *testing.T) {
	cache := synch_service.NewMemoryCache()
	pubsub := synch_service.NewMemoryPubSub()

	chA := CreateCampaignChannelWithDistributed(nil, cache, pubsub, "campaign-123")
	chB := CreateCampaignChannelWithDistributed(nil, cache, pubsub, "campaign-123")

	// Instance A sets sending
	chA.SetSending(true)

	// Instance B should see it via the shared cache
	if !chB.IsSending() {
		t.Fatal("expected IsSending=true on instance B after A set it")
	}

	// Instance A clears sending
	chA.SetSending(false)

	if chB.IsSending() {
		t.Fatal("expected IsSending=false on instance B after A cleared it")
	}
}

func TestCampaignChannel_IsSending_DifferentCampaigns(t *testing.T) {
	cache := synch_service.NewMemoryCache()
	pubsub := synch_service.NewMemoryPubSub()

	chA := CreateCampaignChannelWithDistributed(nil, cache, pubsub, "campaign-A")
	chB := CreateCampaignChannelWithDistributed(nil, cache, pubsub, "campaign-B")

	chA.SetSending(true)

	if chB.IsSending() {
		t.Fatal("campaign-B should not be sending when only campaign-A is")
	}
}

// ─── Cancel (memory-only) ───────────────────────────────────────────────────

func TestCampaignChannel_Cancel_Memory(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := CreateCampaignChannel(&cancel)

	if err := ch.Cancel(); err != nil {
		t.Fatalf("Cancel: %v", err)
	}

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled")
	}
}

func TestCampaignChannel_Cancel_NilCancel(t *testing.T) {
	ch := CreateCampaignChannel(nil)
	err := ch.Cancel()
	if err == nil {
		t.Fatal("expected error when cancel is nil")
	}
}

// ─── Cancel via PubSub (distributed) ────────────────────────────────────────

func TestCampaignChannel_Cancel_Distributed(t *testing.T) {
	cache := synch_service.NewMemoryCache()
	pubsub := synch_service.NewMemoryPubSub()

	// Instance A is executing the campaign
	ctx, cancel := context.WithCancel(context.Background())
	chA := CreateCampaignChannelWithDistributed(&cancel, cache, pubsub, "campaign-456")
	if err := chA.SubscribeCancel(); err != nil {
		t.Fatalf("SubscribeCancel: %v", err)
	}
	defer chA.UnsubscribeCancel()

	// Instance B requests cancellation
	chB := CreateCampaignChannelWithDistributed(nil, cache, pubsub, "campaign-456")
	if err := chB.Cancel(); err != nil {
		t.Fatalf("Cancel from B: %v", err)
	}

	select {
	case <-ctx.Done():
		// expected — A's context was cancelled via PubSub
	case <-time.After(2 * time.Second):
		t.Fatal("context on instance A was not cancelled by instance B's PubSub cancel")
	}
}

func TestCampaignChannel_Cancel_DistributedDifferentCampaigns(t *testing.T) {
	cache := synch_service.NewMemoryCache()
	pubsub := synch_service.NewMemoryPubSub()

	ctx, cancel := context.WithCancel(context.Background())
	chA := CreateCampaignChannelWithDistributed(&cancel, cache, pubsub, "campaign-X")
	if err := chA.SubscribeCancel(); err != nil {
		t.Fatalf("SubscribeCancel: %v", err)
	}
	defer chA.UnsubscribeCancel()

	// Cancel a different campaign
	chOther := CreateCampaignChannelWithDistributed(nil, cache, pubsub, "campaign-Y")
	chOther.Cancel()

	select {
	case <-ctx.Done():
		t.Fatal("campaign-X should not have been cancelled by campaign-Y cancel")
	case <-time.After(100 * time.Millisecond):
		// expected — context still alive
	}
}

// ─── SubscribeCancel / UnsubscribeCancel ────────────────────────────────────

func TestCampaignChannel_SubscribeCancel_NoopWithoutPubSub(t *testing.T) {
	ch := CreateCampaignChannel(nil)
	if err := ch.SubscribeCancel(); err != nil {
		t.Fatalf("expected nil error for memory-only channel, got: %v", err)
	}
	// UnsubscribeCancel should be safe too
	ch.UnsubscribeCancel()
}

func TestCampaignChannel_UnsubscribeCancel_StopsListening(t *testing.T) {
	cache := synch_service.NewMemoryCache()
	pubsub := synch_service.NewMemoryPubSub()

	ctx, cancel := context.WithCancel(context.Background())
	ch := CreateCampaignChannelWithDistributed(&cancel, cache, pubsub, "campaign-unsub")
	ch.SubscribeCancel()

	// Unsubscribe before any cancel signal
	ch.UnsubscribeCancel()

	// Now publish cancel — it should NOT cancel the context
	pubsub.Publish("campaign:campaign-unsub:cancel", []byte("cancel"))

	select {
	case <-ctx.Done():
		t.Fatal("context should not have been cancelled after unsubscribe")
	case <-time.After(100 * time.Millisecond):
		// expected
	}
}

// ─── Concurrent SetSending ──────────────────────────────────────────────────

func TestCampaignChannel_ConcurrentSetSending(t *testing.T) {
	cache := synch_service.NewMemoryCache()
	pubsub := synch_service.NewMemoryPubSub()
	ch := CreateCampaignChannelWithDistributed(nil, cache, pubsub, "concurrent-test")

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Go(func() {
			ch.SetSending(true)
			ch.IsSending()
			ch.SetSending(false)
		})
	}
	wg.Wait()

	if ch.IsSending() {
		t.Fatal("expected IsSending=false after all goroutines finished")
	}
}
