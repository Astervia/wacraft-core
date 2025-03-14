package common_model

import (
	"time"

	"github.com/google/uuid"
)

type Audit struct {
	Id        uuid.UUID `json:"id,omitempty" gorm:"primaryKey;type:uuid;default:gen_random_uuid();not null"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}

type AuditWithDeleted struct {
	Audit

	DeletedAt time.Time `json:"deleted_at,omitempty,omitzero" gorm:"autoDeleteTime, index"`
}
