package database_model

import "gorm.io/gorm"

type Paginable interface {
	PaginateQuery(db **gorm.DB) error
}
