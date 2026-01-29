package workspace_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	user_entity "github.com/Astervia/wacraft-core/src/user/entity"
	"github.com/google/uuid"
)

type Workspace struct {
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"not null;uniqueIndex"`
	Description *string   `json:"description,omitempty"`
	CreatedBy   uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`

	Creator *user_entity.User `json:"creator,omitempty" gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	common_model.Audit
}
