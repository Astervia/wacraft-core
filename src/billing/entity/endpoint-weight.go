package billing_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
)

type EndpointWeight struct {
	Method      string `json:"method" gorm:"not null"`                      // HTTP method (GET, POST, etc.)
	PathPattern string `json:"path_pattern" gorm:"not null"`                // Route pattern (e.g. "/message", "/contact")
	Weight      int    `json:"weight" gorm:"not null;default:1"`            // Cost of this endpoint
	Description *string `json:"description,omitempty"`                      // Optional description

	common_model.Audit
}
