package models

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

type ValidationErrorResponse struct {
	Error   string            `json:"error" example:"Validation failed"`
	Details map[string]string `json:"details,omitempty" example:"email:invalid email format"`
}

