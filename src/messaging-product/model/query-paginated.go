package messaging_product_model

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	database_model "github.com/Astervia/wacraft-core/src/database/model"
	"github.com/google/uuid"
)

type QueryPaginated struct {
	Name        MessagingProductName `json:"name,omitempty" validate:"omitempty,oneof=WhatsApp"`
	WorkspaceID *uuid.UUID           `json:"workspace_id,omitempty" query:"workspace_id"`

	common_model.UnrequiredID
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}
