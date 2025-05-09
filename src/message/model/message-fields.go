package message_model

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	"github.com/google/uuid"
)

type MessageFields struct {
	SenderData         *SenderData   `json:"sender_data,omitempty" gorm:"type:jsonb;default:null"`   // Specific data that allows to send message.
	ReceiverData       *ReceiverData `json:"receiver_data,omitempty" gorm:"type:jsonb;default:null"` // Specific data about the product. For example, the webhook data received.
	ProductData        *ProductData  `json:"product_data,omitempty" gorm:"type:jsonb;default:null"`  // Specific data about the product. For example, the webhook data received.
	FromId             uuid.UUID     `json:"from_id,omitempty" gorm:"type:uuid;default:null"`        // Null if manager sent the message.
	ToId               uuid.UUID     `json:"to_id,omitempty" gorm:"type:uuid;default:null"`          // Null if manager received the message.
	MessagingProductId uuid.UUID     `json:"messaging_product_id,omitempty" gorm:"type:uuid;not null"`

	common_model.AuditWithDeleted
}

func (m *MessageFields) TableName() string {
	return "messages"
}
