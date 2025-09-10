package models

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
}

type ResetPasswordRequest struct {
	Password string `json:"password" validate:"required,min=6" example:"newpassword123"`
}
