package webhook_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	"github.com/google/uuid"
)

type WebhookLog struct {
	Payload          any `json:"payload,omitempty" gorm:"serializer:json;type:jsonb"`
	HttpResponseCode int `json:"http_response_code,omitempty"`
	ResponseData     any `json:"response_data,omitempty" gorm:"serializer:json;type:jsonb"`

	WebhookID uuid.UUID `json:"webhook_id" gorm:"type:uuid;not null"`

	Webhook *Webhook `json:"webhook,omitempty" gorm:"foreignKey:WebhookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Delivery tracking
	DeliveryID    *uuid.UUID `json:"delivery_id,omitempty" gorm:"type:uuid;index"`
	AttemptNumber int        `json:"attempt_number,omitempty" gorm:"default:1"`
	DurationMs    int64      `json:"duration_ms,omitempty"`

	// Request details
	SignatureSent  bool              `json:"signature_sent,omitempty" gorm:"default:false"`
	IdempotencyKey string            `json:"idempotency_key,omitempty" gorm:"index"`
	RequestHeaders map[string]string `json:"request_headers,omitempty" gorm:"serializer:json;type:jsonb"`
	RequestUrl     string            `json:"request_url,omitempty"`

	common_model.Audit
}
