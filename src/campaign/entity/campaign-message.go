package campaign_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	message_entity "github.com/Astervia/wacraft-core/src/message/entity"
	message_model "github.com/Astervia/wacraft-core/src/message/model"
	"github.com/google/uuid"
)

type CampaignMessage struct {
	SenderData *message_model.SenderData `json:"sender_data" gorm:"type:jsonb; not null"` // Specific data that allows to send message.
	MessageID  uuid.UUID                 `json:"message_id,omitempty" gorm:"type:uuid;default:null"`
	CampaignID uuid.UUID                 `json:"campaign_id,omitempty" gorm:"type:uuid; not null"`

	Message  *message_entity.Message `json:"message,omitempty" gorm:"foreignKey:MessageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Campaign *Campaign               `json:"campaign,omitempty" gorm:"foreignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	common_model.Audit
}
