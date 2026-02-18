package billing_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
)

type Plan struct {
	Name            string  `json:"name" gorm:"not null"`
	Slug            string  `json:"slug" gorm:"not null;uniqueIndex"`
	Description     *string `json:"description,omitempty"`
	ThroughputLimit int     `json:"throughput_limit" gorm:"not null"`          // Weighted requests allowed per window. <= 0 means unlimited.
	WindowSeconds   int     `json:"window_seconds" gorm:"not null;default:60"` // Time window in seconds
	DurationDays    int     `json:"duration_days" gorm:"not null;default:30"`  // Plan validity in days
	PriceCents      int64   `json:"price_cents" gorm:"not null;default:0"`     // Price in smallest currency unit
	Currency        string  `json:"currency" gorm:"not null;default:'usd'"`
	IsDefault       bool    `json:"is_default" gorm:"not null;default:false"` // Fallback free plan
	IsCustom        bool    `json:"is_custom" gorm:"not null;default:false"`  // Admin-created custom plans
	Active          bool    `json:"active" gorm:"not null;default:true"`      // Available for purchase
	StripePriceID   *string `json:"stripe_price_id,omitempty"`   // Cached Stripe Price ID for subscription checkouts
	StripeProductID *string `json:"stripe_product_id,omitempty"` // Cached Stripe Product ID for subscription checkouts

	common_model.Audit
}
