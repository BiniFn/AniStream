package utils

import (
	"fmt"

	"github.com/coeeter/aniways/internal/models"
)

const maxPageSize = 100

func ValidatePaginationParams(page, size int) (limit, offset int32, err error) {
	if page < 1 || size < 1 {
		return 0, 0, fmt.Errorf("invalid pagination: page=%d size=%d", page, size)
	}
	if size > maxPageSize {
		return 0, 0, fmt.Errorf("size too large: max=%d got=%d", maxPageSize, size)
	}
	return int32(size), int32((page - 1) * size), nil
}

func PageInfo(page int, pageSize, total int64) models.PageInfo {
	totalPages := int((total + pageSize - 1) / pageSize)
	return models.PageInfo{
		CurrentPage: page,
		TotalPages:  totalPages,
		HasNextPage: page < totalPages,
		HasPrevPage: page > 1,
	}
}
