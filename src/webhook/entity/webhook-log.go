package webhook_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	"github.com/google/uuid"
)

type WebhookLog struct {
	Payload          interface{} `json:"payload,omitempty" gorm:"serializer:json;type:jsonb"`
	HttpResponseCode int         `json:"http_response_code,omitempty"`
	ResponseData     interface{} `json:"response_data,omitempty" gorm:"serializer:json;type:jsonb"`

	WebhookID uuid.UUID `json:"webhook_id" gorm:"type:uuid;not null"`

	Webhook *Webhook `json:"webhook,omitempty" gorm:"foreignKey:WebhookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	common_model.Audit
}
