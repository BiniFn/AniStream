package models

// User management request models

// UpdateUserRequest represents the request body for updating user information
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UpdatePasswordRequest represents the request body for updating user password
type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// DeleteUserRequest represents the request body for deleting user account
type DeleteUserRequest struct {
	Password string `json:"password"`
}