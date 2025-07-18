package user_model

import common_model "github.com/Astervia/wacraft-core/src/common/model"

type Update struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type UpdateWithId struct {
	Role *Role `json:"role,omitempty"`

	common_model.RequiredId
	Update
}

type UpdateWithPassword struct {
	Update
	Password string `json:"password,omitempty"`
}
