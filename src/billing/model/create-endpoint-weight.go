package billing_model

// CreateEndpointWeight represents the data needed to create a custom endpoint weight.
type CreateEndpointWeight struct {
	Method      string  `json:"method" validate:"required"`
	PathPattern string  `json:"path_pattern" validate:"required"`
	Weight      int     `json:"weight" validate:"required,gt=0"`
	Description *string `json:"description,omitempty"`
}
