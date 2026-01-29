package workspace_model

type CreateInvitationRequest struct {
	Email    string   `json:"email" validate:"required,email"`
	Policies []Policy `json:"policies" validate:"required,min=1"`
}

type AcceptInvitationRequest struct {
	Token    string  `json:"token" validate:"required"`
	Name     *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`     // Required if user doesn't exist
	Password *string `json:"password,omitempty" validate:"omitempty,min=8,max=72"` // Required if user doesn't exist
}

type InvitationResponse struct {
	ID          string   `json:"id"`
	WorkspaceID string   `json:"workspace_id"`
	Email       string   `json:"email"`
	Policies    []Policy `json:"policies"`
	ExpiresAt   string   `json:"expires_at"`
	AcceptedAt  *string  `json:"accepted_at,omitempty"`
	InvitedBy   string   `json:"invited_by"`
}
