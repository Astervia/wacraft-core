package webhook_model

// CreateWebhook represents the model for creating a new webhook
type CreateWebhook struct {
	Url           string `json:"url" validate:"required,url"`             // Webhook URL, required and must be a valid URL
	Authorization string `json:"authorization,omitempty"`                 // Optional authorization token
	Event         Event  `json:"event" validate:"required,webhook_event"` // Event type associated with the webhook, required
	HttpMethod    string `json:"http_method" validate:"required,oneof=GET POST PUT DELETE PATCH"`
	Timeout       *int   `json:"timeout,omitempty" validate:"omitempty,gte=1,lte=60"`
}
