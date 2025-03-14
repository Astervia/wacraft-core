package message_entity

import (
	message_model "github.com/Astervia/omni-core/src/message/model"
	messaging_product_entity "github.com/Astervia/omni-core/src/messaging-product/entity"
	status_model "github.com/Astervia/omni-core/src/status/model"
)

type Message struct {
	From             *messaging_product_entity.MessagingProductContact `json:"from,omitempty" gorm:"foreignKey:FromId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`                          // Null if manager sent the message.
	To               *messaging_product_entity.MessagingProductContact `json:"to,omitempty" gorm:"foreignKey:ToId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`                              // Null if manager received the message.
	MessagingProduct *messaging_product_entity.MessagingProduct        `json:"messaging_product,omitempty" gorm:"foreignKey:MessagingProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Null if manager received the message.
	Statuses         []*status_model.StatusFields                      `json:"statuses,omitempty" gorm:"foreignKey:MessageId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	message_model.MessageFields
}
