package auth_model

type GrantType string

var (
	Password     GrantType = "password"
	RefreshToken GrantType = "refresh_token"
)
