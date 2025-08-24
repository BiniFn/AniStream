package models

import "time"

type LibraryRequest struct {
	Status          LibraryStatus `json:"status" validate:"required,oneof=watching completed on_hold dropped plan_to_watch" example:"watching"`
	WatchedEpisodes int32         `json:"watchedEpisodes" validate:"min=0" example:"12"`
}

type ImportJobResponse struct {
	ID string `json:"id" example:"V1StGXR8Z5jdHi6BmyT23"`
}

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

type LibraryImportJobResponse struct {
	ID          string    `json:"id" example:"V1StGXR8Z5jdHi6B"`
	UserID      string    `json:"userId" example:"V1StGXR8Z5jdHi6B"`
	Status      string    `json:"status" example:"pending"`
	CreatedAt   time.Time `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
	CompletedAt time.Time `json:"completedAt" example:"2023-01-01T00:00:00Z"`
}

type LibraryListResponse = Pagination[LibraryResponse]

