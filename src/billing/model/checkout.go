package billing_model

import "github.com/google/uuid"

// CheckoutRequest initiates a plan purchase through a payment provider.
type CheckoutRequest struct {
	PlanID      uuid.UUID   `json:"plan_id" validate:"required"`
	Scope       Scope       `json:"scope" validate:"required,oneof=user workspace"`
	WorkspaceID *uuid.UUID  `json:"workspace_id,omitempty"` // Required when scope=workspace
	PaymentMode PaymentMode `json:"payment_mode" validate:"omitempty,oneof=payment subscription"` // Defaults to "payment" if empty
	SuccessURL  string      `json:"success_url" validate:"required,url"`
	CancelURL   string      `json:"cancel_url" validate:"required,url"`
}

// CheckoutResponse returns the payment provider checkout URL.
type CheckoutResponse struct {
	CheckoutURL string `json:"checkout_url"`
	ExternalID  string `json:"external_id"`
}
