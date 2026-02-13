package billing_model

import "github.com/google/uuid"

// CreateManualSubscription is for admin-created subscriptions (custom plans, manual payment).
type CreateManualSubscription struct {
	PlanID             uuid.UUID  `json:"plan_id" validate:"required"`
	Scope              Scope      `json:"scope" validate:"required,oneof=user workspace"`
	UserID             uuid.UUID  `json:"user_id" validate:"required"`
	WorkspaceID        *uuid.UUID `json:"workspace_id,omitempty"`
	ThroughputOverride *int       `json:"throughput_override,omitempty"`
}
