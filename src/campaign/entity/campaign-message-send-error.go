package campaign_entity

import (
	cmn_model "github.com/Astervia/wacraft-core/src/common/model"
	"github.com/google/uuid"
)

type CampaignMessageSendError struct {
	ErrorData         string    `json:"error_data,omitempty" gorm:"type:string; default:null"` // Error message.
	CampaignMessageID uuid.UUID `json:"campaign_message_id" gorm:"type:uuid; not null"`

	CampaignMessage *CampaignMessage `json:"campaign_message,omitempty" gorm:"foreignKey:CampaignMessageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	cmn_model.Audit
}
