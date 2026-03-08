package billing_model

// CreatePlanPrice adds a currency-specific price to a plan.
// IsDefault=true marks this as the price used when no currency is specified at checkout.
// Unsetting the previous default is handled automatically by the server.
type CreatePlanPrice struct {
	Currency   string `json:"currency" validate:"required"`
	PriceCents int64  `json:"price_cents" validate:"gte=0"`
	IsDefault  bool   `json:"is_default"`
}
