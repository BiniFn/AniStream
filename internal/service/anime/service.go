package anime

import (
	"context"

	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
)

type Service struct {
	repo *repository.Queries
}

func New(repo *repository.Queries) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetRecentlyUpdatedAnimes(
	ctx context.Context,
	page, size int,
) (models.Pagination[models.AnimeDto], error) {
	offset := int32((page - 1) * size)
	limit := int32(size)
	rows, err := s.repo.GetRecentlyUpdatedAnimes(ctx, repository.GetRecentlyUpdatedAnimesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return models.Pagination[models.AnimeDto]{}, err
	}

	total, err := s.repo.GetRecentlyUpdatedAnimesCount(ctx)
	if err != nil {
		return models.Pagination[models.AnimeDto]{}, err
	}

	items := make([]models.AnimeDto, len(rows))
	for i, a := range rows {
		items[i] = models.AnimeDto{}.FromRepository(a)
	}

	totalPages := int((total + int64(size) - 1) / int64(size))
	return models.Pagination[models.AnimeDto]{
		PageInfo: models.PageInfo{
			CurrentPage: page,
			TotalPages:  totalPages,
			HasNextPage: page < totalPages,
			HasPrevPage: page > 1,
		},
		Items: items,
	}, nil
}
