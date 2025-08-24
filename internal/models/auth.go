package models

// Authentication request models

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ForgetPasswordRequest represents the request body for password reset
type ForgetPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest represents the request body for resetting password
type ResetPasswordRequest struct {
	Password string `json:"password"`
}