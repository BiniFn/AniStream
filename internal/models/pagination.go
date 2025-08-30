package models

type PageInfo struct {
	CurrentPage int  `json:"currentPage" validate:"required"`
	TotalPages  int  `json:"totalPages" validate:"required"`
	HasNextPage bool `json:"hasNextPage" validate:"required"`
	HasPrevPage bool `json:"hasPrevPage" validate:"required"`
}

type Pagination[T any] struct {
	PageInfo PageInfo `json:"pageInfo" validate:"required"`
	Items    []T      `json:"items" validate:"required"`
}
