package validators

import (
	"strings"

	status_model "github.com/Astervia/wacraft-core/src/status/model"
	"github.com/go-playground/validator/v10"
)

func searchableStatusColumnValidation(fl validator.FieldLevel) bool {
	input := fl.Field().String()

	statusColumn := status_model.SearchableStatusColumn(strings.ToLower(input))
	return statusColumn.IsValid()
}

func RegisterSearchableStatusColumnValidator(v *validator.Validate) error {
	return v.RegisterValidation("searchable_status_column", searchableStatusColumnValidation)
}
