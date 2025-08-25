package validators

import (
	"strings"

	campaign_model "github.com/Astervia/wacraft-core/src/campaign/model"
	"github.com/go-playground/validator/v10"
)

func jsonMessageKeyValidation(fl validator.FieldLevel) bool {
	input := fl.Field().String()

	msgType := campaign_model.SearchableCampaignColumn(strings.ToLower(input))
	return msgType.IsValid()
}

func RegisterJsonMessageKeyValidator(v *validator.Validate) error {
	return v.RegisterValidation("json_message_key", jsonMessageKeyValidation)
}
