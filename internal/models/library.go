package models

import "time"

// Library management models

// LibraryRequest represents the request body for adding/updating anime in library
type LibraryRequest struct {
	Status          LibraryStatus `json:"status" validate:"required,oneof=watching completed on_hold dropped plan_to_watch" example:"watching"`
	WatchedEpisodes int32         `json:"watchedEpisodes" validate:"min=0" example:"12"`
}

// ImportJobResponse represents the response for library import job creation
type ImportJobResponse struct {
	ID string `json:"id" example:"V1StGXR8Z5jdHi6BmyT23"`
}

// Library response models

// LibraryResponse represents a library entry in HTTP responses
type LibraryResponse struct {
	ID              string        `json:"id" example:"V1StGXR8Z5jdHi6B"`
	UserID          string        `json:"userId" example:"V1StGXR8Z5jdHi6B"`
	AnimeID         string        `json:"animeId" example:"V1StGXR8Z5jdHi6B"`
	Status          string        `json:"status" example:"watching"`
	WatchedEpisodes int32         `json:"watchedEpisodes" example:"12"`
	CreatedAt       time.Time     `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt       time.Time     `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
	Anime           AnimeResponse `json:"anime"`
}


// LibraryImportJobResponse represents a library import job status
type LibraryImportJobResponse struct {
	ID          string    `json:"id" example:"V1StGXR8Z5jdHi6B"`
	UserID      string    `json:"userId" example:"V1StGXR8Z5jdHi6B"`
	Status      string    `json:"status" example:"pending"`
	CreatedAt   time.Time `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
	CompletedAt time.Time `json:"completedAt" example:"2023-01-01T00:00:00Z"`
}

// Pagination type alias
type LibraryListResponse = Pagination[LibraryResponse]