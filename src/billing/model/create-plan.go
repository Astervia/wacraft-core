package billing_model

type CreatePlan struct {
	Name            string           `json:"name" validate:"required"`
	Slug            string           `json:"slug" validate:"required"`
	Description     *string          `json:"description,omitempty"`
	ThroughputLimit int              `json:"throughput_limit"` // <= 0 means unlimited throughput
	WindowSeconds   int              `json:"window_seconds" validate:"required,gt=0"`
	DurationDays    int              `json:"duration_days" validate:"required,gt=0"`
	IsDefault       bool             `json:"is_default"`
	IsCustom        bool             `json:"is_custom"`
	Active          bool             `json:"active"`
	Prices          []CreatePlanPrice `json:"prices,omitempty" validate:"omitempty,dive"` // Optional: create prices atomically with the plan
}
