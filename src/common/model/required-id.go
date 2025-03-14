package common_model

import "github.com/google/uuid"

// Represents a required UUID.
type RequiredId struct {
	Id uuid.UUID `json:"id"` // The unique identifier.
}
