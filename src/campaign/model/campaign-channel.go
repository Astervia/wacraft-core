package campaign_model

import (
	"context"
	"errors"
	"sync"

	websocket_model "github.com/Astervia/omni-core/src/websocket/model"
)

type CampaignChannel struct {
	cancel    *context.CancelFunc
	Sending   bool
	SendingMu *sync.Mutex

	cancelMu *sync.Mutex
	websocket_model.Channel[websocket_model.ClientId, CampaignResults, string]
}

func (c *CampaignChannel) AddCancel(cancel *context.CancelFunc) {
	c.cancelMu.Lock()
	defer c.cancelMu.Unlock()
	c.cancel = cancel
}

func (c *CampaignChannel) Cancel() error {
	c.cancelMu.Lock()
	defer c.cancelMu.Unlock()

	if c.cancel == nil {
		return errors.New("cancel is nil")
	}

	(*c.cancel)()
	return nil
}

func CreateCampaignChannel(
	cancel *context.CancelFunc,
) *CampaignChannel {
	var cancelMu sync.Mutex
	var sendingMu sync.Mutex
	channel := websocket_model.CreateChannel[websocket_model.ClientId, CampaignResults, string]()
	return &CampaignChannel{
		cancel:    cancel,
		Sending:   false,
		cancelMu:  &cancelMu,
		SendingMu: &sendingMu,
		Channel:   *channel,
	}
}
