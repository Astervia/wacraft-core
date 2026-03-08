package billing_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	"github.com/google/uuid"
)

// PlanPrice holds a currency-specific price for a plan.
// A plan may have multiple PlanPrice rows — one per currency.
// Exactly one per plan should have IsDefault=true; that one is used
// when no currency is specified at checkout.
type PlanPrice struct {
	PlanID          uuid.UUID `json:"plan_id" gorm:"type:uuid;not null;uniqueIndex:idx_plan_price_currency;index"`
	Currency        string    `json:"currency" gorm:"type:varchar(10);not null;uniqueIndex:idx_plan_price_currency"`
	PriceCents      int64     `json:"price_cents" gorm:"not null;default:0"`
	IsDefault       bool      `json:"is_default" gorm:"not null;default:false"`
	StripePriceID   *string   `json:"stripe_price_id,omitempty"`   // Cached Stripe Price ID for subscription checkouts
	StripeProductID *string   `json:"stripe_product_id,omitempty"` // Cached Stripe Product ID for subscription checkouts

	Plan *Plan `json:"plan,omitempty" gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	common_model.Audit
}
