package contact_model

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	database_model "github.com/Astervia/wacraft-core/src/database/model"
)

type QueryPaginated struct {
	common_model.UnrequiredID
	CreateContact
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}
