package workspace_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	user_entity "github.com/Astervia/wacraft-core/src/user/entity"
	"github.com/google/uuid"
)

type WorkspaceMember struct {
	WorkspaceID uuid.UUID `json:"workspace_id" gorm:"type:uuid;not null;uniqueIndex:idx_workspace_member"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:idx_workspace_member"`

	Workspace *Workspace        `json:"workspace,omitempty" gorm:"foreignKey:WorkspaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User      *user_entity.User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	common_model.Audit
}
