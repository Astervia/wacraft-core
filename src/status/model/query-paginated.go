package status_model

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	database_model "github.com/Astervia/wacraft-core/src/database/model"
	"github.com/google/uuid"
)

type QueryPaginated struct {
	MessageID uuid.UUID `json:"message_id,omitempty" query:"message_id"`

	common_model.UnrequiredID
	database_model.Paginate
	database_model.DateOrderWithDeletedAt
	database_model.DateWhereWithDeletedAt
}

type Query struct {
	MessageID uuid.UUID `json:"message_id,omitempty" query:"message_id"`

	common_model.UnrequiredID
	database_model.DateOrderWithDeletedAt
	database_model.DateWhereWithDeletedAt
}
