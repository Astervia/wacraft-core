package messaging_product_model

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	database_model "github.com/Astervia/wacraft-core/src/database/model"
)

type QueryPaginated struct {
	Name MessagingProductName `json:"name,omitempty" validate:"omitempty,oneof=WhatsApp"`

	common_model.UnrequiredId
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}
