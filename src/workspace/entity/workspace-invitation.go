package workspace_entity

import (
	"time"

	common_model "github.com/Astervia/wacraft-core/src/common/model"
	user_entity "github.com/Astervia/wacraft-core/src/user/entity"
	workspace_model "github.com/Astervia/wacraft-core/src/workspace/model"
	"github.com/google/uuid"
)

type WorkspaceInvitation struct {
	WorkspaceID uuid.UUID                `json:"workspace_id" gorm:"type:uuid;not null;index"`
	Email       string                   `json:"email" gorm:"not null;index"`
	Token       string                   `json:"-" gorm:"not null;uniqueIndex"`
	Policies    []workspace_model.Policy `json:"policies" gorm:"serializer:json;type:jsonb"`
	ExpiresAt   time.Time                `json:"expires_at" gorm:"not null"`
	AcceptedAt  *time.Time               `json:"accepted_at,omitempty"`
	InvitedBy   uuid.UUID                `json:"invited_by" gorm:"type:uuid;not null"`

	Workspace *Workspace        `json:"workspace,omitempty" gorm:"foreignKey:WorkspaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Inviter   *user_entity.User `json:"inviter,omitempty" gorm:"foreignKey:InvitedBy;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	common_model.Audit
}

func (i *WorkspaceInvitation) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

func (i *WorkspaceInvitation) IsAccepted() bool {
	return i.AcceptedAt != nil
}

func (i *WorkspaceInvitation) IsValid() bool {
	return !i.IsExpired() && !i.IsAccepted()
}
