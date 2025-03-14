package status_model

import (
	common_model "github.com/Astervia/omni-core/src/common/model"
	database_model "github.com/Astervia/omni-core/src/database/model"
	"github.com/google/uuid"
)

type QueryPaginated struct {
	MessageId uuid.UUID `json:"message_id,omitempty" query:"message_id"`

	common_model.UnrequiredId
	database_model.Paginate
	database_model.DateOrderWithDeletedAt
	database_model.DateWhereWithDeletedAt
}

type Query struct {
	MessageId uuid.UUID `json:"message_id,omitempty" query:"message_id"`

	common_model.UnrequiredId
	database_model.DateOrderWithDeletedAt
	database_model.DateWhereWithDeletedAt
}
