package database_model

import "gorm.io/gorm"

type Whereable interface {
	Where(db **gorm.DB, prefix string) error
}
