package status_model

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	"github.com/google/uuid"
)

type StatusFields struct {
	ProductData *ProductData `json:"product_data,omitempty" gorm:"type:jsonb;not null"` // Specific data about the product. For example, the webhook data received.
	MessageId   uuid.UUID    `json:"message_id,omitempty,omitzero" gorm:"type:uuid;not null"`

	common_model.Audit
}

func (s *StatusFields) TableName() string {
	return "statuses"
}
