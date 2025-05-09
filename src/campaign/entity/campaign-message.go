package campaign_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	message_entity "github.com/Astervia/wacraft-core/src/message/entity"
	message_model "github.com/Astervia/wacraft-core/src/message/model"
	"github.com/google/uuid"
)

type CampaignMessage struct {
	SenderData *message_model.SenderData `json:"sender_data" gorm:"type:jsonb; not null"` // Specific data that allows to send message.
	MessageId  uuid.UUID                 `json:"message_id,omitempty" gorm:"type:uuid;default:null"`
	CampaignId uuid.UUID                 `json:"campaign_id,omitempty" gorm:"type:uuid; not null"`

	Message  *message_entity.Message `json:"message,omitempty" gorm:"foreignKey:MessageId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Campaign *Campaign               `json:"campaign,omitempty" gorm:"foreignKey:CampaignId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	common_model.Audit
}
