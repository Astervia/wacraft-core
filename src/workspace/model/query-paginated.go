package workspace_model

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	database_model "github.com/Astervia/wacraft-core/src/database/model"
)

type QueryPaginated struct {
	Name *string `json:"name,omitempty" query:"name"`
	Slug *string `json:"slug,omitempty" query:"slug"`

	common_model.UnrequiredID
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}
