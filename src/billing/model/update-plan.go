package billing_model

// UpdatePlan represents the fields that can be updated on a billing plan.
// Pointer fields allow distinguishing between "not provided" and "set to zero value".
type UpdatePlan struct {
	Name            *string `json:"name,omitempty"`
	Slug            *string `json:"slug,omitempty"`
	Description     *string `json:"description,omitempty"`
	ThroughputLimit *int    `json:"throughput_limit,omitempty"` // <= 0 means unlimited throughput
	WindowSeconds   *int    `json:"window_seconds,omitempty" validate:"omitempty,gt=0"`
	DurationDays    *int    `json:"duration_days,omitempty" validate:"omitempty,gt=0"`
	PriceCents      *int64  `json:"price_cents,omitempty" validate:"omitempty,gte=0"`
	Currency        *string `json:"currency,omitempty"`
	IsDefault       *bool   `json:"is_default,omitempty"`
	IsCustom        *bool   `json:"is_custom,omitempty"`
	Active          *bool   `json:"active,omitempty"`
}
