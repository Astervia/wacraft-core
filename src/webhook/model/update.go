package webhook_model

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
)

// UpdateWebhook represents the model for updating an existing webhook
type UpdateWebhook struct {
	Url           string `json:"url,omitempty" validate:"omitempty,url"`             // Optional updated URL, validated to be a valid URL
	Authorization string `json:"authorization,omitempty"`                            // Optional updated authorization token
	Event         Event  `json:"event,omitempty" validate:"omitempty,webhook_event"` // Optional updated event associated with the webhook
	HttpMethod    string `json:"http_method,omitempty" validate:"omitempty,oneof=GET POST PUT DELETE PATCH"`
	Timeout       *int   `json:"timeout,omitempty" validate:"omitempty,gte=1,lte=60"`

	// New fields for enhanced webhook functionality
	MaxRetries    *int              `json:"max_retries,omitempty" validate:"omitempty,gte=0,lte=10"`         // Max retry attempts
	RetryDelayMs  *int              `json:"retry_delay_ms,omitempty" validate:"omitempty,gte=100,lte=60000"` // Base retry delay in ms
	IsActive      *bool             `json:"is_active,omitempty"`                                             // Enable/disable webhook
	CustomHeaders map[string]string `json:"custom_headers,omitempty"`                                        // Custom headers
	EventFilter   *EventFilter      `json:"event_filter,omitempty"`                                          // Event filter

	common_model.RequiredID
}
