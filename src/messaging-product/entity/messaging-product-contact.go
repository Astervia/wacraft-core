package messaging_product_entity

import (
	"time"

	common_model "github.com/Astervia/wacraft-core/src/common/model"
	contact_entity "github.com/Astervia/wacraft-core/src/contact/entity"
	messaging_product_model "github.com/Astervia/wacraft-core/src/messaging-product/model"
	"github.com/google/uuid"
)

// Relation of a contact with a messaging product.
type MessagingProductContact struct {
	ProductDetails     *messaging_product_model.ProductDetails `json:"product_details,omitempty" gorm:"serializer:json;type:jsonb"`
	ContactID          uuid.UUID                               `json:"contact_id" gorm:"type:uuid;not null"`
	MessagingProductID uuid.UUID                               `json:"messaging_product_id" gorm:"type:uuid;not null"`
	Blocked            bool                                    `json:"blocked" gorm:"default:false"`
	LastReadAt         time.Time                               `json:"last_read_at,omitempty" gorm:"default:null"` // Timestamp of the last read action.

	Contact          *contact_entity.Contact `json:"contact,omitempty" gorm:"foreignKey:ContactID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MessagingProduct *MessagingProduct       `json:"messaging_product,omitempty" gorm:"foreignKey:MessagingProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// MessagesFrom []message_model.MessageFields `json:"messages_from,omitempty" gorm:"foreignKey:FromID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// MessageTo      []message_model.MessageFields `json:"messages_to,omitempty" gorm:"foreignKey:ToID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	common_model.Audit
}
