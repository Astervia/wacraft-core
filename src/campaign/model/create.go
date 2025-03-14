package campaign_model

import (
	message_model "github.com/Astervia/omni-core/src/message/model"
	"github.com/google/uuid"
)

type CreateCampaign struct {
	Name               string     `json:"name"`
	MessagingProductId *uuid.UUID `json:"messaging_product_id" validate:"required"`
}

type CreateCampaignMessage struct {
	SenderData *message_model.SenderData `json:"sender_data"` // Specific data that allows to send message.
	CampaignId uuid.UUID                 `json:"campaign_id"`
}
