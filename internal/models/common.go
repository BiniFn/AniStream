package models

// ErrorResponse represents a standard API error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// ValidationErrorResponse represents validation error with field details
type ValidationErrorResponse struct {
	Error   string            `json:"error" example:"Validation failed"`
	Details map[string]string `json:"details,omitempty" example:"email:invalid email format"`
}