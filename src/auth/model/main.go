package auth_model

type TokenRequest struct {
	GrantType    GrantType `json:"grant_type" example:"password"` // password | refresh_token
	Username     string    `json:"username,omitempty" example:"user@mail.com"`
	Password     string    `json:"password,omitempty" example:"123456"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token" example:"eyJhbGciOi..."`  // JWT
	RefreshToken string    `json:"refresh_token" example:"eyJhbGciOi..."` // JWT
	TokenType    TokenType `json:"token_type" example:"bearer"`           // Always "bearer"
	ExpiresIn    int       `json:"expires_in" example:"3600"`             // seconds
}
