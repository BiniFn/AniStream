package models

// ErrorResponse represents a standard API error response
type ErrorResponse struct {
	Error string `json:"error"`
}