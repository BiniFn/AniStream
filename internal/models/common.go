package models

type ErrorResponse struct {
	Error string `json:"error" validate:"required" example:"Invalid request"`
}

type ValidationErrorResponse struct {
	Error   string            `json:"error" validate:"required" example:"Validation failed"`
	Details map[string]string `json:"details,omitempty" example:"email:invalid email format"`
}
