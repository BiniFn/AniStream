package models

import "time"

type LibraryRequest struct {
	Status          LibraryStatus `json:"status" validate:"required" example:"watching"`
	WatchedEpisodes int32         `json:"watchedEpisodes" validate:"min=0" example:"12"`
}

type ImportJobResponse struct {
	ID string `json:"id" validate:"required" example:"V1StGXR8Z5jdHi6BmyT23"`
}

type LibraryResponse struct {
	ID              string        `json:"id" validate:"required" example:"V1StGXR8Z5jdHi6B"`
	UserID          string        `json:"userId" validate:"required" example:"V1StGXR8Z5jdHi6B"`
	AnimeID         string        `json:"animeId" validate:"required" example:"V1StGXR8Z5jdHi6B"`
	Status          LibraryStatus `json:"status" validate:"required" example:"watching"`
	WatchedEpisodes int32         `json:"watchedEpisodes" validate:"required" example:"12"`
	CreatedAt       time.Time     `json:"createdAt" validate:"required" example:"2023-01-01T00:00:00Z"`
	UpdatedAt       time.Time     `json:"updatedAt" validate:"required" example:"2023-01-01T00:00:00Z"`
	Anime           AnimeResponse `json:"anime" validate:"required"`
}

type LibraryImportJobResponse struct {
	ID          string        `json:"id" validate:"required" example:"V1StGXR8Z5jdHi6B"`
	UserID      string        `json:"userId" validate:"required" example:"V1StGXR8Z5jdHi6B"`
	Status      LibraryStatus `json:"status" validate:"required" example:"pending"`
	CreatedAt   time.Time     `json:"createdAt" validate:"required" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time     `json:"updatedAt" validate:"required" example:"2023-01-01T00:00:00Z"`
	CompletedAt time.Time     `json:"completedAt" validate:"required" example:"2023-01-01T00:00:00Z"`
}

type LibraryListResponse = Pagination[LibraryResponse]
