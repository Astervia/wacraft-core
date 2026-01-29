package user_entity

import (
	"time"

	common_model "github.com/Astervia/wacraft-core/src/common/model"
	"github.com/google/uuid"
)

type EmailVerification struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Token     string    `json:"-" gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	Verified  bool      `json:"verified" gorm:"default:false"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	common_model.Audit
}
