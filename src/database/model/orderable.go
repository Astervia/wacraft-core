package database_model

import "gorm.io/gorm"

type Orderable interface {
	OrderQuery(db **gorm.DB, prefix string) error
}
