package messaging_product_model

import (
	common_model "github.com/Astervia/omni-core/src/common/model"
	database_model "github.com/Astervia/omni-core/src/database/model"
	"github.com/google/uuid"
)

type QueryContactPaginated struct {
	ContactID          uuid.UUID `json:"contact_id,omitempty" query:"contact_id"`
	MessagingProductID uuid.UUID `json:"messaging_product_id,omitempty" query:"messaging_product_id"`
	Blocked            bool      `json:"blocked,omitempty" query:"blocked"`

	common_model.UnrequiredId
	database_model.Paginate
	DateOrder
	DateWhere
}

type QueryWhatsAppContactPaginated struct {
	ContactID uuid.UUID `json:"contact_id,omitempty" query:"contact_id"`
	UnrequiredWhatsAppProductDetails

	common_model.UnrequiredId
	database_model.Paginate
	DateOrder
	DateWhere
}
