package validators

import (
	"strings"

	user_model "github.com/Astervia/wacraft-core/src/user/model"
	"github.com/go-playground/validator/v10"
)

func searchableUserColumnValidation(fl validator.FieldLevel) bool {
	input := fl.Field().String()

	msgType := user_model.SearchableUserColumn(strings.ToLower(input))
	return msgType.IsValid()
}

func RegisterSearchableUserColumnValidator(v *validator.Validate) error {
	return v.RegisterValidation("searchable_user_column", searchableUserColumnValidation)
}
