package phone_config_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	workspace_entity "github.com/Astervia/wacraft-core/src/workspace/entity"
	"github.com/google/uuid"
)

// PhoneConfig stores WhatsApp Business Account configuration per workspace.
type PhoneConfig struct {
	Name               string     `json:"name" gorm:"not null"` // Friendly name for this configuration
	WorkspaceID        *uuid.UUID `json:"workspace_id,omitempty" gorm:"type:uuid;index"`
	WabaID             string     `json:"waba_id" gorm:"not null;index:idx_active_waba_id,where:is_active = true,unique"` // Phone Number ID from Meta (unique when active)
	WabaAccountID      string     `json:"waba_account_id" gorm:"not null"`                                                // WhatsApp Business Account ID
	DisplayPhone       string     `json:"display_phone" gorm:"not null"`                                                  // Display phone number (e.g., +1234567890)
	AccessToken        string     `json:"access_token" gorm:"not null"`
	MetaAppSecret      string     `json:"meta_app_secret" gorm:"not null"`
	WebhookVerifyToken string     `json:"webhook_verify_token" gorm:"not null"`
	IsActive           bool       `json:"is_active" gorm:"default:true"`

	Workspace *workspace_entity.Workspace `json:"workspace,omitempty" gorm:"foreignKey:WorkspaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	common_model.Audit
}

func (PhoneConfig) TableName() string {
	return "phone_configs"
}
