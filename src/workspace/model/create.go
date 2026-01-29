package workspace_model

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
)

type CreateWorkspace struct {
	Name        string  `json:"name" validate:"required"`
	Slug        string  `json:"slug" validate:"required"`
	Description *string `json:"description,omitempty"`
}

type UpdateWorkspace struct {
	Name        *string `json:"name,omitempty"`
	Slug        *string `json:"slug,omitempty"`
	Description *string `json:"description,omitempty"`

	common_model.RequiredID
}
