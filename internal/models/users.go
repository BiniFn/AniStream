package models

import "time"

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30" example:"johndoe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=6,max=128" example:"password123"`
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30" example:"johndoe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required" example:"oldpassword123"`
	NewPassword string `json:"newPassword" validate:"required,min=6,max=128" example:"newpassword123"`
}

type DeleteUserRequest struct {
	Password string `json:"password" validate:"required" example:"password123"`
}

type UserResponse struct {
	ID             string    `json:"id" example:"V1StGXR8Z5jdHi6B"`
	Username       string    `json:"username" example:"johndoe"`
	Email          string    `json:"email" example:"john@example.com"`
	ProfilePicture string    `json:"profilePicture,omitempty" example:"https://res.cloudinary.com/example/avatar.jpg"`
	CreatedAt      time.Time `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt      time.Time `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
}
