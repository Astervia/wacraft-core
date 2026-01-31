package user_model

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message" example:"If your email is registered, you will receive a password reset link"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type ResetPasswordResponse struct {
	Message string `json:"message" example:"Password reset successful"`
}
