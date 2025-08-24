package common_model

import "github.com/google/uuid"

// Represents a required UUID.
type RequiredID struct {
	ID uuid.UUID `json:"id"` // The unique identifier.
}
