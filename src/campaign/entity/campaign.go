package campaign_entity

import (
	"time"

	common_model "github.com/Astervia/wacraft-core/src/common/model"
	messaging_product_entity "github.com/Astervia/wacraft-core/src/messaging-product/entity"
	workspace_entity "github.com/Astervia/wacraft-core/src/workspace/entity"
	"github.com/google/uuid"
)

// Allows user to create messages to be sent to a list of contacts at the same time.
type Campaign struct {
	Name               string     `json:"name,omitempty" gorm:"not null"`
	MessagingProductID *uuid.UUID `json:"messaging_product_id,omitempty" gorm:"type:uuid;not null"`
	WorkspaceID        *uuid.UUID `json:"workspace_id,omitempty" gorm:"type:uuid;index"`

	// Status tracks the campaign lifecycle: draft | scheduled | running | completed | failed | cancelled.
	Status string `json:"status,omitempty" gorm:"not null;default:'draft'"`
	// ScheduledAt is the UTC time at which the scheduler worker will start sending the campaign.
	// Null means the campaign is not scheduled.
	ScheduledAt *time.Time `json:"scheduled_at,omitempty" gorm:"type:timestamptz"`

	MessagingProduct *messaging_product_entity.MessagingProduct `json:"messaging_product,omitempty" gorm:"foreignKey:MessagingProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Null if manager received the message.
	Workspace        *workspace_entity.Workspace                `json:"workspace,omitempty" gorm:"foreignKey:WorkspaceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	common_model.Audit
}
