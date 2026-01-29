package phone_config_model

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	database_model "github.com/Astervia/wacraft-core/src/database/model"
	"github.com/google/uuid"
)

// QueryPaginated represents the paginated query structure for phone configs
type QueryPaginated struct {
	WorkspaceID  *uuid.UUID `json:"workspace_id,omitempty" query:"workspace_id"`
	Name         string     `json:"name,omitempty" query:"name"`
	WabaID       string     `json:"waba_id,omitempty" query:"waba_id"` // Phone Number ID
	DisplayPhone string     `json:"display_phone,omitempty" query:"display_phone"`
	IsActive     *bool      `json:"is_active,omitempty" query:"is_active"`

	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
	common_model.UnrequiredID
}
