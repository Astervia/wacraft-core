package validators

import (
	"strings"

	webhook_model "github.com/Astervia/wacraft-core/src/webhook/model"
	"github.com/go-playground/validator/v10"
)

func searchableWebhookColumnValidation(fl validator.FieldLevel) bool {
	input := fl.Field().String()

	msgType := webhook_model.SearchableWebhookColumn(strings.ToLower(input))
	return msgType.IsValid()
}

func RegisterSearchableWebhookColumnValidator(v *validator.Validate) error {
	return v.RegisterValidation("searchable_webhook_column", searchableWebhookColumnValidation)
}
