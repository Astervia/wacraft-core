package messaging_product_model

import (
	common_model "github.com/Astervia/omni-core/src/common/model"
	database_model "github.com/Astervia/omni-core/src/database/model"
	"github.com/google/uuid"
)

type QueryContactPaginated struct {
	database_model.Paginate
	QueryContact
}

type QueryWhatsAppContactPaginated struct {
	ContactID uuid.UUID `json:"contact_id,omitempty" query:"contact_id"`
	UnrequiredWhatsAppProductDetails

	common_model.UnrequiredId
	database_model.Paginate
	DateOrder
	DateWhere
}
