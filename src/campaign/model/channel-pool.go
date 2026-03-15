package campaign_model

import (
	"context"
	"sync"

	synch_contract "github.com/Astervia/wacraft-core/src/synch/contract"
	websocket_model "github.com/Astervia/wacraft-core/src/websocket/model"
	"github.com/google/uuid"
)

type ChannelPool struct {
	mu       *sync.Mutex
	channels map[uuid.UUID]CampaignChannel

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

	if len(channel.Clients) == 0 {
		channel.UnsubscribeProgress()
		channel.UnsubscribeCancel()
		delete(cp.channels, campaignID)
	}
}

func CreateChannelPool() *ChannelPool {
	var mu sync.Mutex
	return &ChannelPool{
		mu:       &mu,
		channels: make(map[uuid.UUID]CampaignChannel),
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
