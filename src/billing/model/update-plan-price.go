package billing_model

// UpdatePlanPrice allows changing price or default status for a plan price entry.
// Setting IsDefault=true automatically unsets the previous default for that plan.
type UpdatePlanPrice struct {
	PriceCents *int64 `json:"price_cents,omitempty" validate:"omitempty,gte=0"`
	IsDefault  *bool  `json:"is_default,omitempty"`
}
