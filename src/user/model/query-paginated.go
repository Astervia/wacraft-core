package user_model

import (
	common_model "github.com/Astervia/omni-core/src/common/model"
	database_model "github.com/Astervia/omni-core/src/database/model"
)

type QueryPaginated struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Role  *Role  `json:"role,omitempty"`

	common_model.UnrequiredId
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}
