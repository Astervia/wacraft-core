package user_entity

import (
	"time"

	common_model "github.com/Astervia/wacraft-core/src/common/model"
	"github.com/google/uuid"
)

type PasswordResetToken struct {
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	Token     string     `json:"-" gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	UsedAt    *time.Time `json:"used_at,omitempty"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	common_model.Audit
}

func (p *PasswordResetToken) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

func (p *PasswordResetToken) IsUsed() bool {
	return p.UsedAt != nil
}

func (p *PasswordResetToken) IsValid() bool {
	return !p.IsExpired() && !p.IsUsed()
}
