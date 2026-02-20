package phone_config_model

// CreatePhoneConfig is used to create a new phone configuration
type CreatePhoneConfig struct {
	Name               string `json:"name" validate:"required"`
	WabaID             string `json:"waba_id" validate:"required"`         // Phone Number ID from Meta
	WabaAccountID      string `json:"waba_account_id" validate:"required"` // WhatsApp Business Account ID
	DisplayPhone       string `json:"display_phone" validate:"required"`
	AccessToken        string `json:"access_token" validate:"required"`
	MetaAppSecret      string `json:"meta_app_secret" validate:"required"`
	WebhookVerifyToken string `json:"webhook_verify_token" validate:"required"`
	IsActive           *bool  `json:"is_active,omitempty"`
}

// UpdatePhoneConfig is used to update an existing phone configuration
type UpdatePhoneConfig struct {
	Name               *string `json:"name,omitempty"`
	WabaID             *string `json:"waba_id,omitempty"`
	WabaAccountID      *string `json:"waba_account_id,omitempty"`
	DisplayPhone       *string `json:"display_phone,omitempty"`
	AccessToken        *string `json:"access_token,omitempty"`
	MetaAppSecret      *string `json:"meta_app_secret,omitempty"`
	WebhookVerifyToken *string `json:"webhook_verify_token,omitempty"`
	IsActive           *bool   `json:"is_active,omitempty"`
}
