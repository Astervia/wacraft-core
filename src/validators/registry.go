package validators

import (
	"github.com/go-playground/validator/v10"
)

// RegisterAllValidators registers all custom validation functions used in the lib.
func RegisterAllValidators(v *validator.Validate) error {
	registrations := []func(*validator.Validate) error{
		RegisterWebhookEventValidator,
		RegisterSearchableCampaignColumnValidator,
		RegisterSearchableUserColumnValidator,
		RegisterSearchableWebhookColumnValidator,
		RegisterJsonMessageKeyValidator,
	}

	for _, register := range registrations {
		if err := register(v); err != nil {
			return err
		}
	}

	return nil
}
