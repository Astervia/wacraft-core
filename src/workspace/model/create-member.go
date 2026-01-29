package workspace_model

import "github.com/google/uuid"

type CreateMember struct {
	UserID   uuid.UUID `json:"user_id" validate:"required"`
	Policies []Policy  `json:"policies" validate:"required,min=1"`
}

type UpdateMemberPolicies struct {
	Policies []Policy `json:"policies" validate:"required"`
}
