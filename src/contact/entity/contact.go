package contact_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	workspace_entity "github.com/Astervia/wacraft-core/src/workspace/entity"
	"github.com/google/uuid"
)

type Contact struct {
	Name        *string    `json:"name,omitempty"`
	Email       *string    `json:"email,omitempty"`
	PhotoPath   *string    `json:"photo_path,omitempty"`
	WorkspaceID *uuid.UUID `json:"workspace_id,omitempty" gorm:"type:uuid;index"`

	Workspace *workspace_entity.Workspace `json:"workspace,omitempty" gorm:"foreignKey:WorkspaceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	common_model.Audit
}
