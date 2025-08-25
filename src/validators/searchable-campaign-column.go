package validators

import (
	"strings"

	campaign_model "github.com/Astervia/wacraft-core/src/campaign/model"
	"github.com/go-playground/validator/v10"
)

func searchableCampaignColumnValidation(fl validator.FieldLevel) bool {
	input := fl.Field().String()

	msgType := campaign_model.SearchableCampaignColumn(strings.ToLower(input))
	return msgType.IsValid()
}

func RegisterSearchableCampaignColumnValidator(v *validator.Validate) error {
	return v.RegisterValidation("searchable_campaign_column", searchableCampaignColumnValidation)
}
