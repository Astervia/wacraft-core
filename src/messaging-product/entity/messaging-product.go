package messaging_product_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	messaging_product_model "github.com/Astervia/wacraft-core/src/messaging-product/model"
	phone_config_entity "github.com/Astervia/wacraft-core/src/phone-config/entity"
	workspace_entity "github.com/Astervia/wacraft-core/src/workspace/entity"
	"github.com/google/uuid"
)

// This entity must provide a way to easily access diferent message clusters through database and realtime channels
type MessagingProduct struct {
	Name          messaging_product_model.MessagingProductName `json:"name,omitempty" gorm:"not null" validate:"omitempty,oneof=WhatsApp"` // Add type:enum('WhatsApp'); when it becomes supported by GORM and PostgreSQL
	WorkspaceID   *uuid.UUID                                   `json:"workspace_id,omitempty" gorm:"type:uuid;index"`
	PhoneConfigID *uuid.UUID                                   `json:"phone_config_id,omitempty" gorm:"type:uuid;index"`

	Workspace   *workspace_entity.Workspace      `json:"workspace,omitempty" gorm:"foreignKey:WorkspaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PhoneConfig *phone_config_entity.PhoneConfig `json:"phone_config,omitempty" gorm:"foreignKey:PhoneConfigID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	common_model.Audit
}
