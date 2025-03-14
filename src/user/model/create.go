package user_model

type Create struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     *Role  `json:"role,omitempty"`
}
