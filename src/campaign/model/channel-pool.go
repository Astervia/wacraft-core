package campaign_model

import (
	"context"
	"sync"

	websocket_model "github.com/Astervia/wacraft-core/src/websocket/model"
	"github.com/google/uuid"
)

type ChannelPool struct {
	mu       *sync.Mutex
	channels map[uuid.UUID]CampaignChannel
}

func (cp *ChannelPool) AddUser(
	client websocket_model.Client[websocket_model.ClientId],
	key string,
	campaignId uuid.UUID,
	cancel *context.CancelFunc,
) *CampaignChannel {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	// Check if the channel already exists
	if channel, ok := cp.channels[campaignId]; ok {
		channel.Clients[key] = client

		return &channel
	}

	channel := CreateCampaignChannel(cancel)
	channel.AppendClient(client, key)

	cp.channels[campaignId] = *channel
	return channel
}

func (cp *ChannelPool) RemoveUser(
	clientKey string,
	campaignId uuid.UUID,
) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	// Check if the channel exists
	channel, ok := cp.channels[campaignId]
	if !ok {
		return
	}

	channel.RemoveClient(clientKey)

	if len(channel.Clients) == 0 {
		delete(cp.channels, campaignId)
	}
}

func CreateChannelPool() *ChannelPool {
	var mu sync.Mutex
	return &ChannelPool{
		mu:       &mu,
		channels: make(map[uuid.UUID]CampaignChannel),
	}
}
