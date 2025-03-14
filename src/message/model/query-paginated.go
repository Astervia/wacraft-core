package message_model

import (
	common_model "github.com/Astervia/omni-core/src/common/model"
	database_model "github.com/Astervia/omni-core/src/database/model"
	"github.com/google/uuid"
)

type QueryPaginated struct {
	FromId             uuid.UUID `json:"from_id,omitempty" query:"from_id"`
	ToId               uuid.UUID `json:"to_id,omitempty" query:"to_id"`
	MessagingProductId uuid.UUID `json:"messaging_product_id,omitempty" query:"messaging_product_id"`

	common_model.UnrequiredId
	database_model.Paginate
	database_model.DateOrderWithDeletedAt
	database_model.DateWhereWithDeletedAt
}

type Query struct {
	FromId             uuid.UUID `json:"from_id,omitempty" query:"from_id"`
	ToId               uuid.UUID `json:"to_id,omitempty" query:"to_id"`
	MessagingProductId uuid.UUID `json:"messaging_product_id,omitempty" query:"messaging_product_id"`

	common_model.UnrequiredId
	database_model.DateOrderWithDeletedAt
	database_model.DateWhereWithDeletedAt
}
