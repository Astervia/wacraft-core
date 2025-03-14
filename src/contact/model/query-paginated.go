package contact_model

import (
	common_model "github.com/Astervia/omni-core/src/common/model"
	database_model "github.com/Astervia/omni-core/src/database/model"
)

type QueryPaginated struct {
	common_model.UnrequiredId
	CreateContact
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}
