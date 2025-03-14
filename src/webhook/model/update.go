package webhook_model

import (
	common_model "github.com/Astervia/omni-core/src/common/model"
)

// UpdateWebhook represents the model for updating an existing webhook
type UpdateWebhook struct {
	Url           string `json:"url,omitempty" validate:"url"` // Optional updated URL, validated to be a valid URL
	Authorization string `json:"authorization,omitempty"`      // Optional updated authorization token
	Event         Event  `json:"event,omitempty"`              // Optional updated event associated with the webhook
	HttpMethod    string `json:"http_method,omitempty" validate:"omitempty,oneof=GET POST PUT DELETE PATCH"`
	Timeout       *int   `json:"timeout,omitempty" validate:"omitempty,gte=1,leq=60"`

	common_model.RequiredId
}
