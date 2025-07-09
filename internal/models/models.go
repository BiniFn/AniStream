package models

type PageInfo struct {
	CurrentPage int  `json:"currentPage"`
	TotalPages  int  `json:"totalPages"`
	HasNextPage bool `json:"hasNextPage"`
	HasPrevPage bool `json:"hasPrevPage"`
}

type Pagination[T any] struct {
	PageInfo PageInfo `json:"pageInfo"`
	Items    []T      `json:"items"`
}
