package webhook_entity

import (
	"time"

	common_model "github.com/Astervia/wacraft-core/src/common/model"
	"github.com/google/uuid"
)

// DeliveryStatus represents the status of a webhook delivery
type DeliveryStatus string

const (
	DeliveryStatusPending    DeliveryStatus = "pending"
	DeliveryStatusAttempted  DeliveryStatus = "attempted" // Has been attempted but not yet successful
	DeliveryStatusSucceeded  DeliveryStatus = "succeeded"
	DeliveryStatusFailed     DeliveryStatus = "failed"      // Max retries exhausted
	DeliveryStatusDeadLetter DeliveryStatus = "dead_letter" // Moved to dead letter queue
)

// WebhookDelivery represents a webhook delivery in the queue
type WebhookDelivery struct {
	WebhookID uuid.UUID `json:"webhook_id,omitempty" gorm:"type:uuid;not null;index"`
	Webhook   *Webhook  `json:"webhook,omitempty" gorm:"foreignKey:WebhookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	IdempotencyKey string         `json:"idempotency_key,omitempty" gorm:"not null;uniqueIndex"`
	Payload        any            `json:"payload,omitempty" gorm:"serializer:json;type:jsonb"`
	Status         DeliveryStatus `json:"status,omitempty" gorm:"default:'pending';index"`

	// Retry management
	AttemptCount  int        `json:"attempt_count,omitempty" gorm:"default:0"`
	MaxAttempts   int        `json:"max_attempts,omitempty" gorm:"default:3"`
	NextAttemptAt *time.Time `json:"next_attempt_at,omitempty" gorm:"index"`

	// Last attempt details
	LastAttemptAt    *time.Time `json:"last_attempt_at,omitempty"`
	LastHttpCode     *int       `json:"last_http_code,omitempty"`
	LastError        *string    `json:"last_error,omitempty"`
	LastResponseBody *string    `json:"last_response_body,omitempty" gorm:"type:text"`

	// Event metadata
	EventType      string    `json:"event_type,omitempty" gorm:"not null;index"`
	EventTimestamp time.Time `json:"event_timestamp,omitempty" gorm:"not null"`

	common_model.Audit
}

// TableName specifies the table name for GORM
func (WebhookDelivery) TableName() string {
	return "webhook_deliveries"
}
