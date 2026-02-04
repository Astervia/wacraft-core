package webhook_model

// CreateWebhook represents the model for creating a new webhook
type CreateWebhook struct {
	Url           string `json:"url" validate:"required,url"`             // Webhook URL, required and must be a valid URL
	Authorization string `json:"authorization,omitempty"`                 // Optional authorization token
	Event         Event  `json:"event" validate:"required,webhook_event"` // Event type associated with the webhook, required
	HttpMethod    string `json:"http_method" validate:"required,oneof=GET POST PUT DELETE PATCH"`
	Timeout       *int   `json:"timeout,omitempty" validate:"omitempty,gte=1,lte=60"`

	// New fields for enhanced webhook functionality
	SigningEnabled bool              `json:"signing_enabled,omitempty"`                                       // Enable HMAC-SHA256 signing
	MaxRetries     *int              `json:"max_retries,omitempty" validate:"omitempty,gte=0,lte=10"`         // Max retry attempts (0-10)
	RetryDelayMs   *int              `json:"retry_delay_ms,omitempty" validate:"omitempty,gte=100,lte=60000"` // Base retry delay in ms (100-60000)
	CustomHeaders  map[string]string `json:"custom_headers,omitempty"`                                        // Custom headers to send with requests
	EventFilter    *EventFilter      `json:"event_filter,omitempty"`                                          // Filter to match specific events
}
