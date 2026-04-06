package campaign_model

import (
	"context"
	"sync"

	synch_contract "github.com/Astervia/wacraft-core/src/synch/contract"
	websocket_model "github.com/Astervia/wacraft-core/src/websocket/model"
	"github.com/google/uuid"
)

type ChannelPool struct {
	mu         *sync.Mutex
	channels   map[uuid.UUID]CampaignChannel
	workerRefs map[uuid.UUID]int // counts scheduler/worker holds per campaign

	// Distributed primitives (nil in memory-only mode).
	cache  synch_contract.DistributedCache
	pubsub synch_contract.PubSub
}

func (cp *ChannelPool) AddUser(
	client websocket_model.Client[websocket_model.ClientID],
	key string,
	campaignID uuid.UUID,
	cancel *context.CancelFunc,
) *CampaignChannel {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	// Check if the channel already exists
	if channel, ok := cp.channels[campaignID]; ok {
		channel.Clients[key] = client

		return &channel
	}

	var channel *CampaignChannel
	if cp.cache != nil || cp.pubsub != nil {
		channel = CreateCampaignChannelWithDistributed(cancel, cp.cache, cp.pubsub, campaignID.String())
	} else {
		channel = CreateCampaignChannel(cancel)
	}
	channel.AppendClient(client, key)

	cp.channels[campaignID] = *channel
	return channel
}

func (cp *ChannelPool) RemoveUser(
	clientKey string,
	campaignID uuid.UUID,
) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	// Check if the channel exists
	channel, ok := cp.channels[campaignID]
	if !ok {
		return
	}

	channel.RemoveClient(clientKey)

	// Only clean up when no WebSocket clients AND no worker holds remain.
	if len(channel.Clients) == 0 && cp.workerRefs[campaignID] == 0 {
		channel.UnsubscribeProgress()
		channel.UnsubscribeCancel()
		delete(cp.channels, campaignID)
		delete(cp.workerRefs, campaignID)
	}
}

// GetOrCreateChannel returns the CampaignChannel for campaignID, creating it if
// it does not exist. It increments the worker reference count so the channel is
// not removed when all WebSocket clients disconnect while the worker is still
// sending. Call ReleaseChannel when the worker is done.
func (cp *ChannelPool) GetOrCreateChannel(campaignID uuid.UUID) *CampaignChannel {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.workerRefs[campaignID]++

	if channel, ok := cp.channels[campaignID]; ok {
		return &channel
	}

	var channel *CampaignChannel
	if cp.cache != nil || cp.pubsub != nil {
		channel = CreateCampaignChannelWithDistributed(nil, cp.cache, cp.pubsub, campaignID.String())
	} else {
		channel = CreateCampaignChannel(nil)
	}

	cp.channels[campaignID] = *channel
	return channel
}

// ReleaseChannel decrements the worker reference count for campaignID. If the
// reference count reaches zero and there are no connected WebSocket clients, the
// channel is removed from the pool.
func (cp *ChannelPool) ReleaseChannel(campaignID uuid.UUID) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if cp.workerRefs[campaignID] > 0 {
		cp.workerRefs[campaignID]--
	}

	channel, ok := cp.channels[campaignID]
	if !ok {
		return
	}

	if len(channel.Clients) == 0 && cp.workerRefs[campaignID] == 0 {
		channel.UnsubscribeProgress()
		channel.UnsubscribeCancel()
		delete(cp.channels, campaignID)
		delete(cp.workerRefs, campaignID)
	}
}

func CreateChannelPool() *ChannelPool {
	var mu sync.Mutex
	return &ChannelPool{
		mu:         &mu,
		channels:   make(map[uuid.UUID]CampaignChannel),
		workerRefs: make(map[uuid.UUID]int),
	}
}

// CreateChannelPoolWithDistributed creates a ChannelPool that passes
// distributed cache and pub/sub to every CampaignChannel it creates.
func CreateChannelPoolWithDistributed(
	cache synch_contract.DistributedCache,
	pubsub synch_contract.PubSub,
) *ChannelPool {
	pool := CreateChannelPool()
	pool.cache = cache
	pool.pubsub = pubsub
	return pool
}
