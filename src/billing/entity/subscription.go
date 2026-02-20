package billing_entity

import (
	"time"

	billing_model "github.com/Astervia/wacraft-core/src/billing/model"
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	user_entity "github.com/Astervia/wacraft-core/src/user/entity"
	workspace_entity "github.com/Astervia/wacraft-core/src/workspace/entity"
	"github.com/google/uuid"
)

type Subscription struct {
	PlanID               uuid.UUID                        `json:"plan_id" gorm:"type:uuid;not null"`
	Scope                billing_model.Scope              `json:"scope" gorm:"type:varchar(20);not null"`  // "user" or "workspace"
	UserID               uuid.UUID                        `json:"user_id" gorm:"type:uuid;not null"`       // Who purchased
	WorkspaceID          *uuid.UUID                       `json:"workspace_id,omitempty" gorm:"type:uuid"` // Set when scope=workspace
	ThroughputOverride   *int                             `json:"throughput_override,omitempty"`           // Admin override for custom plans
	StartsAt             time.Time                        `json:"starts_at" gorm:"not null"`
	ExpiresAt            time.Time                        `json:"expires_at" gorm:"not null;index"`
	CancelledAt          *time.Time                       `json:"cancelled_at,omitempty"`
	PaymentProvider      string                           `json:"payment_provider" gorm:"type:varchar(50);not null;default:'manual'"`
	PaymentExternalID    *string                          `json:"payment_external_id,omitempty"`
	PaymentMode          billing_model.PaymentMode        `json:"payment_mode" gorm:"type:varchar(20);not null;default:'payment'"` // "payment" (one-time) or "subscription" (recurring)
	StripeSubscriptionID *string                          `json:"stripe_subscription_id,omitempty"`                                // Stripe subscription ID for recurring plans
	CancelAtPeriodEnd    bool                             `json:"cancel_at_period_end" gorm:"not null;default:false"`              // True when cancellation is pending (active until ExpiresAt)
	Status               billing_model.SubscriptionStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`        // "pending", "active", "cancelled"

	Plan      *Plan                       `json:"plan,omitempty" gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User      *user_entity.User           `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Workspace *workspace_entity.Workspace `json:"workspace,omitempty" gorm:"foreignKey:WorkspaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	common_model.Audit
}

// IsActive checks if the subscription is currently active.
func (s *Subscription) IsActive() bool {
	now := time.Now()
	return s.Status == billing_model.SubscriptionStatusActive && s.CancelledAt == nil && now.After(s.StartsAt) && now.Before(s.ExpiresAt)
}

// EffectiveThroughput returns the throughput limit considering any override.
func (s *Subscription) EffectiveThroughput(planLimit int) int {
	if s.ThroughputOverride != nil {
		return *s.ThroughputOverride
	}
	return planLimit
}
