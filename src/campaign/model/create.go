package campaign_model

import (
	message_model "github.com/Astervia/wacraft-core/src/message/model"
	"github.com/google/uuid"
)

type CreateCampaign struct {
	Name               string     `json:"name"`
	MessagingProductID *uuid.UUID `json:"messaging_product_id" validate:"required"`
}

type CreateCampaignMessage struct {
	SenderData *message_model.SenderData `json:"sender_data"` // Specific data that allows to send message.
	CampaignID uuid.UUID                 `json:"campaign_id"`
}
