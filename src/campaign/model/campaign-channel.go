package campaign_model

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	synch_contract "github.com/Astervia/wacraft-core/src/synch/contract"
	websocket_model "github.com/Astervia/wacraft-core/src/websocket/model"
	"github.com/pterm/pterm"
)

type CampaignChannel struct {
	cancel    *context.CancelFunc
	Sending   bool
	SendingMu *sync.Mutex

	cancelMu *sync.Mutex

	// Distributed primitives (nil in memory-only mode).
	cache       synch_contract.DistributedCache
	pubsub      synch_contract.PubSub
	cancelSub   synch_contract.Subscription
	progressSub synch_contract.Subscription
	campaignID  string

	websocket_model.Channel[websocket_model.ClientID, CampaignResults, string]
}

func (c *CampaignChannel) AddCancel(cancel *context.CancelFunc) {
	c.cancelMu.Lock()
	defer c.cancelMu.Unlock()
	c.cancel = cancel
}

func (c *CampaignChannel) Cancel() error {
	c.cancelMu.Lock()
	defer c.cancelMu.Unlock()

	// If distributed pub/sub is available, publish cancel signal
	// so the executing instance (which may be different) receives it.
	if c.pubsub != nil && c.campaignID != "" {
		if err := c.pubsub.Publish("campaign:"+c.campaignID+":cancel", []byte("cancel")); err != nil {
			return err
		}
		return nil
	}

	if c.cancel == nil {
		return errors.New("cancel is nil")
	}

	(*c.cancel)()
	return nil
}

// SetSending updates the Sending flag. When a distributed cache is available,
// the flag is also stored/deleted in Redis so other instances can see it.
func (c *CampaignChannel) SetSending(sending bool) {
	c.SendingMu.Lock()
	c.Sending = sending
	c.SendingMu.Unlock()

	if c.cache != nil && c.campaignID != "" {
		key := "campaign:" + c.campaignID + ":sending"
		if sending {
			c.cache.Set(key, []byte("1"), 3600_000_000_000) // 1 hour
		} else {
			c.cache.Delete(key)
		}
	}
}

// IsSending checks the Sending flag. When a distributed cache is available,
// it checks Redis first (another instance may have started the campaign).
func (c *CampaignChannel) IsSending() bool {
	if c.cache != nil && c.campaignID != "" {
		key := "campaign:" + c.campaignID + ":sending"
		_, found, err := c.cache.Get(key)
		if err == nil && found {
			return true
		}
	}

	c.SendingMu.Lock()
	defer c.SendingMu.Unlock()
	return c.Sending
}

// SubscribeCancel starts listening for cancel signals from other instances.
// When a signal arrives, the local context.CancelFunc is invoked.
// Call this after AddCancel when starting a campaign.
func (c *CampaignChannel) SubscribeCancel() error {
	if c.pubsub == nil || c.campaignID == "" {
		return nil
	}

	sub, err := c.pubsub.Subscribe("campaign:" + c.campaignID + ":cancel")
	if err != nil {
		return err
	}
	c.cancelSub = sub

	go func() {
		for range sub.Channel() {
			c.cancelMu.Lock()
			if c.cancel != nil {
				(*c.cancel)()
			}
			c.cancelMu.Unlock()
		}
	}()

	return nil
}

// UnsubscribeCancel stops listening for cancel signals. Call when campaign ends.
func (c *CampaignChannel) UnsubscribeCancel() {
	if c.cancelSub != nil {
		c.cancelSub.Unsubscribe()
		c.cancelSub = nil
	}
}

// BroadcastProgress sends campaign progress to all local WebSocket clients.
// When a PubSub backend is configured the data is published to the distributed
// channel so every instance can deliver it to its own local clients.
func (c *CampaignChannel) BroadcastProgress(data CampaignResults) {
	if c.pubsub != nil && c.campaignID != "" {
		b, err := json.Marshal(data)
		if err != nil {
			pterm.DefaultLogger.Error("campaign broadcast progress marshal error: " + err.Error())
			return
		}
		if err := c.pubsub.Publish("campaign:"+c.campaignID+":progress", b); err != nil {
			pterm.DefaultLogger.Error("campaign broadcast progress publish error: " + err.Error())
		}
		return
	}
	// Memory-only mode: local broadcast.
	c.BroadcastJsonMultithread(data)
}

// subscribeProgress subscribes to cross-instance progress events for this campaign.
// Messages received are broadcast to all local WebSocket clients.
func (c *CampaignChannel) subscribeProgress() {
	if c.pubsub == nil || c.campaignID == "" {
		return
	}
	sub, err := c.pubsub.Subscribe("campaign:" + c.campaignID + ":progress")
	if err != nil {
		pterm.DefaultLogger.Error("campaign subscribe progress error: " + err.Error())
		return
	}
	c.progressSub = sub

	go func() {
		for msg := range sub.Channel() {
			var data CampaignResults
			if err := json.Unmarshal(msg, &data); err != nil {
				pterm.DefaultLogger.Error("campaign progress unmarshal error: " + err.Error())
				continue
			}
			c.BroadcastJsonMultithread(data)
		}
	}()
}

// UnsubscribeProgress stops listening for progress events. Call when all clients disconnect.
func (c *CampaignChannel) UnsubscribeProgress() {
	if c.progressSub != nil {
		c.progressSub.Unsubscribe()
		c.progressSub = nil
	}
}

func CreateCampaignChannel(
	cancel *context.CancelFunc,
) *CampaignChannel {
	var cancelMu sync.Mutex
	var sendingMu sync.Mutex
	channel := websocket_model.CreateChannel[websocket_model.ClientID, CampaignResults, string]()
	return &CampaignChannel{
		cancel:    cancel,
		Sending:   false,
		cancelMu:  &cancelMu,
		SendingMu: &sendingMu,
		Channel:   *channel,
	}
}

// CreateCampaignChannelWithDistributed creates a CampaignChannel backed by
// distributed cache (for Sending flag) and pub/sub (for Cancel and Progress).
func CreateCampaignChannelWithDistributed(
	cancel *context.CancelFunc,
	cache synch_contract.DistributedCache,
	pubsub synch_contract.PubSub,
	campaignID string,
) *CampaignChannel {
	ch := CreateCampaignChannel(cancel)
	ch.cache = cache
	ch.pubsub = pubsub
	ch.campaignID = campaignID
	ch.subscribeProgress()
	return ch
}
