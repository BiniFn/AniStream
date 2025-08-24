package models

// Library management models

// LibraryRequest represents the request body for adding/updating anime in library
type LibraryRequest struct {
	Status          string `json:"status"`
	WatchedEpisodes int32  `json:"watchedEpisodes"`
}

// ImportJobResponse represents the response for library import job creation
type ImportJobResponse struct {
	ID string `json:"id"`
}