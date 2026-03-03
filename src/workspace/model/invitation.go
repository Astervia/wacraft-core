package workspace_model

type CreateInvitationRequest struct {
	Email    string   `json:"email" validate:"required,email"`
	Policies []Policy `json:"policies" validate:"required,min=1"`
}

type ClaimInvitationRequest struct {
	Token string `json:"token" validate:"required"`
}

type ClaimInvitationResponse struct {
	Message     string `json:"message"`
	WorkspaceID string `json:"workspace_id"`
}

type InvitationResponse struct {
	ID          string   `json:"id"`
	WorkspaceID string   `json:"workspace_id"`
	Email       string   `json:"email"`
	Token       string   `json:"token"`
	Policies    []Policy `json:"policies"`
	ExpiresAt   string   `json:"expires_at"`
	AcceptedAt  *string  `json:"accepted_at,omitempty"`
	InvitedBy   string   `json:"invited_by"`
}
