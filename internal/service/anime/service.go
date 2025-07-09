package anime

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
)

type Service struct {
	repo      *repository.Queries
	refresher *MetadataRefresher
}

func New(repo *repository.Queries, refresher *MetadataRefresher) *Service {
	return &Service{repo: repo, refresher: refresher}
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

	for _, r := range rows {
		go s.refresher.MaybeRefresh(ctx, r.MalID.Int32)
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

func (s *Service) GetAnimeByID(
	ctx context.Context,
	id string,
) (*models.AnimeWithMetadataDto, error) {
	a, err := s.repo.GetAnimeById(ctx, id)
	if err != nil {
		return nil, err
	}

	m, err := s.repo.GetAnimeMetadataByMalId(ctx, a.MalID.Int32)
	if err == sql.ErrNoRows || time.Since(m.UpdatedAt.Time) > s.refresher.ttl {
		if err := s.refresher.RefreshBlocking(ctx, a.MalID.Int32); err != nil {
			log.Printf("blocking refresh failed: %v", err)
		}
		m, err = s.repo.GetAnimeMetadataByMalId(ctx, a.MalID.Int32)
		if err != nil {
			return nil, err
		}
	}

	dto := models.AnimeWithMetadataDto{}.FromRepository(a, m)

	return &dto, nil
}
