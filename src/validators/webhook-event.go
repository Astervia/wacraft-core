package validators

import (
	"strings"

	webhook_model "github.com/Astervia/wacraft-core/src/webhook/model"
	"github.com/go-playground/validator/v10"
)

func webhookEventValidation(fl validator.FieldLevel) bool {
	input := fl.Field().String()

	msgType := webhook_model.Event(strings.ToLower(input))
	return msgType.IsValid()
}

func RegisterWebhookEventValidator(v *validator.Validate) error {
	return v.RegisterValidation("webhook_event", webhookEventValidation)
}
