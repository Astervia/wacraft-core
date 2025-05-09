package messaging_product_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	messaging_product_model "github.com/Astervia/wacraft-core/src/messaging-product/model"
)

// This entity must provide a way to easily access diferent message clusters through database and realtime channels
type MessagingProduct struct {
	Name messaging_product_model.MessagingProductName `json:"name,omitempty" gorm:"not null" validate:"omitempty,oneof=WhatsApp"` // Add type:enum('WhatsApp'); when it becomes supported by GORM and PostgreSQL

	common_model.Audit
}
