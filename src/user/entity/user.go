package user_entity

import (
	common_model "github.com/Astervia/omni-core/src/common/model"
	crypto_service "github.com/Astervia/omni-core/src/crypto/service"
	user_model "github.com/Astervia/omni-core/src/user/model"
	"gorm.io/gorm"
)

type User struct {
	Name     string           `json:"name,omitempty" gorm:"not null"`
	Email    string           `json:"email,omitempty" gorm:"not null;unique"`
	Password string           `json:"password,omitempty" gorm:"not null"`
	Role     *user_model.Role `json:"role,omitempty" gorm:"type:varchar(20);not null;default:'user'"`

	common_model.Audit
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Hash passord before saving
	var err error = nil
	u.Password, err = crypto_service.HashPassword(u.Password)
	return err
}
