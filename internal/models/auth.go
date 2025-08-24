package models

// Authentication request models

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// ForgetPasswordRequest represents the request body for password reset
type ForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
}

// ResetPasswordRequest represents the request body for resetting password
type ResetPasswordRequest struct {
	Password string `json:"password" validate:"required,min=6" example:"newpassword123"`
}