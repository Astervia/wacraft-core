package common_model

import "github.com/google/uuid"

// Represents an optional UUID.
type UnrequiredID struct {
	ID uuid.UUID `json:"id,omitempty" validate:"omitempty"` // The unique identifier
}
