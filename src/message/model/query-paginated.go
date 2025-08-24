package message_model

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	database_model "github.com/Astervia/wacraft-core/src/database/model"
	"github.com/google/uuid"
)

type QueryPaginated struct {
	FromID             *uuid.UUID `json:"from_id,omitempty" query:"from_id" validate:"omitempty"`
	ToID               *uuid.UUID `json:"to_id,omitempty" query:"to_id" validate:"omitempty"`
	MessagingProductID uuid.UUID  `json:"messaging_product_id,omitempty" query:"messaging_product_id" validate:"omitempty"`

	common_model.UnrequiredID
	database_model.Paginate
	database_model.DateOrderWithDeletedAt
	database_model.DateWhereWithDeletedAt
}

type Query struct {
	FromID             *uuid.UUID `json:"from_id,omitempty" query:"from_id" validate:"omitempty"`
	ToID               *uuid.UUID `json:"to_id,omitempty" query:"to_id" validate:"omitempty"`
	MessagingProductID uuid.UUID  `json:"messaging_product_id,omitempty" query:"messaging_product_id" validate:"omitempty"`

	common_model.UnrequiredID
	database_model.DateOrderWithDeletedAt
	database_model.DateWhereWithDeletedAt
}
