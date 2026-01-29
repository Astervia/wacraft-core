package workspace_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	workspace_model "github.com/Astervia/wacraft-core/src/workspace/model"
	"github.com/google/uuid"
)

type WorkspaceMemberPolicy struct {
	WorkspaceMemberID uuid.UUID              `json:"workspace_member_id" gorm:"type:uuid;not null;index"`
	Policy            workspace_model.Policy `json:"policy" gorm:"type:varchar(50);not null"`

	WorkspaceMember *WorkspaceMember `json:"workspace_member,omitempty" gorm:"foreignKey:WorkspaceMemberID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	common_model.Audit
}
