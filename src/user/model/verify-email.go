package user_model

type VerifyEmailResponse struct {
	Message string `json:"message" example:"Email verified successfully"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResendVerificationResponse struct {
	Message string `json:"message" example:"If your email is registered, you will receive a verification link"`
}
