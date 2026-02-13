package billing_entity

import (
	"time"

	billing_model "github.com/Astervia/wacraft-core/src/billing/model"
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	"github.com/google/uuid"
)

type UsageLog struct {
	Scope           billing_model.Scope `json:"scope" gorm:"type:varchar(20);not null"`
	UserID          *uuid.UUID          `json:"user_id,omitempty" gorm:"type:uuid;index"`
	WorkspaceID     *uuid.UUID          `json:"workspace_id,omitempty" gorm:"type:uuid;index"`
	WindowStart     time.Time           `json:"window_start" gorm:"not null;index"`
	WindowEnd       time.Time           `json:"window_end" gorm:"not null"`
	WeightedCount   int64               `json:"weighted_count" gorm:"not null"`
	ThroughputLimit int64               `json:"throughput_limit" gorm:"not null"`

	common_model.Audit
}
