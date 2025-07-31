package library

import (
	"context"
	"errors"

	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/coeeter/aniways/internal/utils"
)

type LibraryService struct {
	repo      *repository.Queries
	refresher *anime.MetadataRefresher
}

func NewLibraryService(repo *repository.Queries, refresher *anime.MetadataRefresher) *LibraryService {
	return &LibraryService{
		repo:      repo,
		refresher: refresher,
	}
}

type GetLibraryParams struct {
	UserID       string
	Status       string
	Page         int
	ItemsPerPage int
}

func isValidStatus(status string) bool {
	switch status {
	case string(repository.LibraryStatusWatching),
		string(repository.LibraryStatusCompleted),
		string(repository.LibraryStatusDropped),
		string(repository.LibraryStatusPaused),
		string(repository.LibraryStatusPlanning):
		return true
	default:
		return false
	}
}

var ErrInvalidStatus = errors.New("invalid status")

func (s *LibraryService) GetLibrary(ctx context.Context, params GetLibraryParams) (models.Pagination[LibraryDto], error) {
	limit, offset, err := utils.ValidatePaginationParams(params.Page, params.ItemsPerPage)
	if err != nil {
		return models.Pagination[LibraryDto]{}, err
	}

	if !isValidStatus(params.Status) {
		return models.Pagination[LibraryDto]{}, ErrInvalidStatus
	}

	rows, err := s.repo.GetLibrary(ctx, repository.GetLibraryParams{
		UserID: params.UserID,
		Status: repository.LibraryStatus(params.Status),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return models.Pagination[LibraryDto]{}, err
	}

	for _, r := range rows {
		s.refresher.Enqueue(r.Anime.MalID.Int32)
	}

	total, err := s.repo.GetLibraryCount(ctx, repository.GetLibraryCountParams{
		UserID: params.UserID,
		Status: repository.LibraryStatus(params.Status),
	})
	if err != nil {
		return models.Pagination[LibraryDto]{}, err
	}

	out := make([]LibraryDto, 0, len(rows))
	for _, item := range rows {
		out = append(out, LibraryDto{}.FromRepository(item.Library, item.Anime))
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(params.Page, pageSize, total)
	return models.Pagination[LibraryDto]{
		Items:    out,
		PageInfo: pageInfo,
	}, nil
}

func (s *LibraryService) GetLibraryByAnimeID(ctx context.Context, userID, animeID string) (LibraryDto, error) {
	row, err := s.repo.GetLibraryOfUserByAnimeID(ctx, repository.GetLibraryOfUserByAnimeIDParams{
		UserID:  userID,
		AnimeID: animeID,
	})
	if err != nil {
		return LibraryDto{}, err
	}

	return LibraryDto{}.FromRepository(row.Library, row.Anime), nil
}

var ErrInvalidWatchedEpisodes = errors.New("invalid watched episodes")

func (s *LibraryService) SaveLibrary(ctx context.Context, userID, animeID, status string, watchedEpisodes int32) (LibraryDto, error) {
	if !isValidStatus(status) {
		return LibraryDto{}, ErrInvalidStatus
	}

	if watchedEpisodes < 0 {
		return LibraryDto{}, ErrInvalidWatchedEpisodes
	}

	err := s.repo.UpsertLibrary(ctx, repository.UpsertLibraryParams{
		UserID:          userID,
		AnimeID:         animeID,
		Status:          repository.LibraryStatus(status),
		WatchedEpisodes: watchedEpisodes,
	})
	if err != nil {
		return LibraryDto{}, err
	}

	lib, err := s.GetLibraryByAnimeID(ctx, userID, animeID)
	if err != nil {
		return LibraryDto{}, err
	}

	return lib, nil
}

func (s *LibraryService) DeleteLibrary(ctx context.Context, userID, animeID string) error {
	err := s.repo.DeleteLibrary(ctx, repository.DeleteLibraryParams{
		UserID:  userID,
		AnimeID: animeID,
	})
	return err
}

type GetContinueWatchingAnimeParams struct {
	UserID             string
	Page, ItemsPerPage int
}

func (s *LibraryService) GetContinueWatching(ctx context.Context, params GetContinueWatchingAnimeParams) (models.Pagination[LibraryDto], error) {
	limit, offset, err := utils.ValidatePaginationParams(params.Page, params.ItemsPerPage)
	if err != nil {
		return models.Pagination[LibraryDto]{}, err
	}

	rows, err := s.repo.GetContinueWatchingAnime(ctx, repository.GetContinueWatchingAnimeParams{
		UserID: params.UserID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return models.Pagination[LibraryDto]{}, err
	}

	total, err := s.repo.GetContinueWatchingAnimeCount(ctx, params.UserID)
	if err != nil {
		return models.Pagination[LibraryDto]{}, err
	}

	out := make([]LibraryDto, 0, len(rows))
	for _, item := range rows {
		out = append(out, LibraryDto{}.FromRepository(item.Library, item.Anime))
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(params.Page, pageSize, total)
	return models.Pagination[LibraryDto]{
		Items:    out,
		PageInfo: pageInfo,
	}, nil
}

type GetPlanToWatchAnimeParams struct {
	UserID             string
	Page, ItemsPerPage int
}

func (s *LibraryService) GetPlanToWatch(ctx context.Context, params GetPlanToWatchAnimeParams) (models.Pagination[LibraryDto], error) {
	limit, offset, err := utils.ValidatePaginationParams(params.Page, params.ItemsPerPage)
	if err != nil {
		return models.Pagination[LibraryDto]{}, err
	}

	rows, err := s.repo.GetPlanToWatchAnime(ctx, repository.GetPlanToWatchAnimeParams{
		UserID: params.UserID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return models.Pagination[LibraryDto]{}, err
	}

	total, err := s.repo.GetPlanToWatchAnimeCount(ctx, params.UserID)
	if err != nil {
		return models.Pagination[LibraryDto]{}, err
	}

	out := make([]LibraryDto, 0, len(rows))
	for _, item := range rows {
		out = append(out, LibraryDto{}.FromRepository(item.Library, item.Anime))
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(params.Page, pageSize, total)
	return models.Pagination[LibraryDto]{
		Items:    out,
		PageInfo: pageInfo,
	}, nil
}
