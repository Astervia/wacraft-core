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

	channel := CreateCampaignChannel(cancel)
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
