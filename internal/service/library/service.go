package library

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/coeeter/aniways/internal/mappers"
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/coeeter/aniways/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func (s *LibraryService) GetLibrary(ctx context.Context, params GetLibraryParams) (models.LibraryListResponse, error) {
	limit, offset, err := utils.ValidatePaginationParams(params.Page, params.ItemsPerPage)
	if err != nil {
		return models.LibraryListResponse{}, err
	}

	if !isValidStatus(params.Status) {
		return models.LibraryListResponse{}, ErrInvalidStatus
	}

	rows, err := s.repo.GetLibrary(ctx, repository.GetLibraryParams{
		UserID: params.UserID,
		Status: repository.LibraryStatus(params.Status),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return models.LibraryListResponse{}, err
	}

	for _, r := range rows {
		s.refresher.Enqueue(r.Anime.MalID.Int32)
	}

	total, err := s.repo.GetLibraryCount(ctx, repository.GetLibraryCountParams{
		UserID: params.UserID,
		Status: repository.LibraryStatus(params.Status),
	})
	if err != nil {
		return models.LibraryListResponse{}, err
	}

	out := make([]models.LibraryResponse, 0, len(rows))
	for _, item := range rows {
		out = append(out, mappers.LibraryFromRepository(item.Library, item.Anime))
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(params.Page, pageSize, total)
	return models.LibraryListResponse{
		Items:    out,
		PageInfo: pageInfo,
	}, nil
}

var ErrLibraryNotFound = errors.New("library not found")

func (s *LibraryService) GetLibraryByAnimeID(ctx context.Context, userID, animeID string) (models.LibraryResponse, error) {
	row, err := s.repo.GetLibraryOfUserByAnimeID(ctx, repository.GetLibraryOfUserByAnimeIDParams{
		UserID:  userID,
		AnimeID: animeID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return models.LibraryResponse{}, ErrLibraryNotFound
	}
	if err != nil {
		return models.LibraryResponse{}, err
	}

	return mappers.LibraryFromRepository(row.Library, row.Anime), nil
}

var ErrInvalidWatchedEpisodes = errors.New("invalid watched episodes")

func (s *LibraryService) CreateLibrary(ctx context.Context, userID, animeID, status string, watchedEpisodes int32) (models.LibraryResponse, error) {
	if !isValidStatus(status) {
		return models.LibraryResponse{}, ErrInvalidStatus
	}

	if watchedEpisodes < 0 {
		return models.LibraryResponse{}, ErrInvalidWatchedEpisodes
	}

	err := s.repo.InsertLibrary(ctx, repository.InsertLibraryParams{
		UserID:          userID,
		AnimeID:         animeID,
		Status:          repository.LibraryStatus(status),
		WatchedEpisodes: watchedEpisodes,
	})
	if err != nil {
		return models.LibraryResponse{}, err
	}

	s.queueSync(ctx, userID, animeID, repository.LibraryActionsAddEntry, SyncPayload{
		Status:          &status,
		WatchedEpisodes: &watchedEpisodes,
	})

	lib, err := s.GetLibraryByAnimeID(ctx, userID, animeID)
	if err != nil {
		return models.LibraryResponse{}, err
	}

	return lib, nil
}

func (s *LibraryService) UpdateLibrary(ctx context.Context, userID, animeID, status string, watchedEpisodes int32) (models.LibraryResponse, error) {
	if !isValidStatus(status) {
		return models.LibraryResponse{}, ErrInvalidStatus
	}

	if watchedEpisodes < 0 {
		return models.LibraryResponse{}, ErrInvalidWatchedEpisodes
	}

	old, err := s.GetLibraryByAnimeID(ctx, userID, animeID)
	if err != nil {
		return models.LibraryResponse{}, err
	}

	err = s.repo.UpdateLibrary(ctx, repository.UpdateLibraryParams{
		UserID:          userID,
		AnimeID:         animeID,
		Status:          repository.LibraryStatus(status),
		WatchedEpisodes: watchedEpisodes,
	})
	if err != nil {
		return models.LibraryResponse{}, err
	}

	lib, err := s.GetLibraryByAnimeID(ctx, userID, animeID)
	if err != nil {
		return models.LibraryResponse{}, err
	}

	if old.Status != models.LibraryStatus(status) {
		s.queueSync(ctx, userID, animeID, repository.LibraryActionsUpdateStatus, SyncPayload{
			Status: &status,
		})
	}

	if old.WatchedEpisodes != watchedEpisodes {
		s.queueSync(ctx, userID, animeID, repository.LibraryActionsUpdateProgress, SyncPayload{
			WatchedEpisodes: &watchedEpisodes,
		})
	}

	return lib, nil
}

func (s *LibraryService) DeleteLibrary(ctx context.Context, userID, animeID string) error {
	err := s.repo.DeleteLibrary(ctx, repository.DeleteLibraryParams{
		UserID:  userID,
		AnimeID: animeID,
	})
	s.queueSync(ctx, userID, animeID, repository.LibraryActionsDeleteEntry, SyncPayload{})
	return err
}

type GetContinueWatchingAnimeParams struct {
	UserID             string
	Page, ItemsPerPage int
}

func (s *LibraryService) GetContinueWatching(ctx context.Context, params GetContinueWatchingAnimeParams) (models.LibraryListResponse, error) {
	limit, offset, err := utils.ValidatePaginationParams(params.Page, params.ItemsPerPage)
	if err != nil {
		return models.LibraryListResponse{}, err
	}

	rows, err := s.repo.GetContinueWatchingAnime(ctx, repository.GetContinueWatchingAnimeParams{
		UserID: params.UserID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return models.LibraryListResponse{}, err
	}

	total, err := s.repo.GetContinueWatchingAnimeCount(ctx, params.UserID)
	if err != nil {
		return models.LibraryListResponse{}, err
	}

	out := make([]models.LibraryResponse, 0, len(rows))
	for _, item := range rows {
		out = append(out, mappers.LibraryFromRepository(item.Library, item.Anime))
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(params.Page, pageSize, total)
	return models.LibraryListResponse{
		Items:    out,
		PageInfo: pageInfo,
	}, nil
}

type GetPlanToWatchAnimeParams struct {
	UserID             string
	Page, ItemsPerPage int
}

func (s *LibraryService) GetPlanToWatch(ctx context.Context, params GetPlanToWatchAnimeParams) (models.LibraryListResponse, error) {
	limit, offset, err := utils.ValidatePaginationParams(params.Page, params.ItemsPerPage)
	if err != nil {
		return models.LibraryListResponse{}, err
	}

	rows, err := s.repo.GetPlanToWatchAnime(ctx, repository.GetPlanToWatchAnimeParams{
		UserID: params.UserID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return models.LibraryListResponse{}, err
	}

	total, err := s.repo.GetPlanToWatchAnimeCount(ctx, params.UserID)
	if err != nil {
		return models.LibraryListResponse{}, err
	}

	out := make([]models.LibraryResponse, 0, len(rows))
	for _, item := range rows {
		out = append(out, mappers.LibraryFromRepository(item.Library, item.Anime))
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(params.Page, pageSize, total)
	return models.LibraryListResponse{
		Items:    out,
		PageInfo: pageInfo,
	}, nil
}

func (s *LibraryService) GetLibraryStats(ctx context.Context, userID string) (models.LibraryStatsResponse, error) {
	var stats models.LibraryStatsResponse
	var err error

	watchingCount, err := s.repo.GetLibraryCount(ctx, repository.GetLibraryCountParams{
		UserID: userID,
		Status: repository.LibraryStatusWatching,
	})
	if err != nil {
		return stats, err
	}
	stats.Watching = watchingCount

	planningCount, err := s.repo.GetPlanToWatchAnimeCount(ctx, userID)
	if err != nil {
		return stats, err
	}
	stats.Planning = planningCount

	completedCount, err := s.repo.GetLibraryCount(ctx, repository.GetLibraryCountParams{
		UserID: userID,
		Status: repository.LibraryStatusCompleted,
	})
	if err != nil {
		return stats, err
	}
	stats.Completed = completedCount

	return stats, nil
}

type SyncPayload struct {
	Status          *string `json:"status,omitempty"`
	WatchedEpisodes *int32  `json:"watched_episodes,omitempty"`
}

func (s *LibraryService) queueSync(ctx context.Context, userID, animeID string, action repository.LibraryActions, payload SyncPayload) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	providers := []repository.Provider{
		repository.ProviderMyanimelist,
		repository.ProviderAnilist,
	}

	for _, p := range providers {
		_ = s.repo.UpsertLibrarySync(ctx, repository.UpsertLibrarySyncParams{
			UserID:   userID,
			AnimeID:  animeID,
			Provider: p,
			Action:   action,
			Payload:  data,
		})
	}
}

var ErrInvalidProvider = errors.New("invalid provider")

func (s *LibraryService) ImportLibrary(ctx context.Context, userID, provider string) (string, error) {
	if provider != string(repository.ProviderMyanimelist) && provider != string(repository.ProviderAnilist) {
		return "", ErrInvalidProvider
	}

	return s.repo.CreateLibraryImportJob(ctx, repository.CreateLibraryImportJobParams{
		UserID:   userID,
		Provider: repository.Provider(provider),
	})
}

var ErrJobNotFound = errors.New("job not found")

func (s *LibraryService) GetImportLibraryStatus(ctx context.Context, jobID string) (models.LibraryImportJobResponse, error) {
	status, err := s.repo.GetLibraryImportJob(ctx, jobID)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.LibraryImportJobResponse{}, ErrJobNotFound
	}
	if err != nil {
		return models.LibraryImportJobResponse{}, err
	}

	return mappers.LibraryImportJobFromRepository(status), nil
}

func (s *LibraryService) ClearLibrary(ctx context.Context, userID string) error {
	return s.repo.ClearLibrary(ctx, userID)
}

var ErrVariationMismatch = errors.New("variations must have the same MAL ID")

func (s *LibraryService) SwitchLibraryVariation(ctx context.Context, userID, currentAnimeID, variationID string) (models.LibraryResponse, error) {
	// Get both animes to validate they share the same MAL ID
	currentAnime, err := s.repo.GetAnimeById(ctx, currentAnimeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.LibraryResponse{}, anime.ErrAnimeNotFound
		}
		return models.LibraryResponse{}, err
	}

	variationAnime, err := s.repo.GetAnimeById(ctx, variationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.LibraryResponse{}, anime.ErrAnimeNotFound
		}
		return models.LibraryResponse{}, err
	}

	// Validate both have MAL IDs and they match
	if !currentAnime.MalID.Valid || !variationAnime.MalID.Valid {
		return models.LibraryResponse{}, ErrVariationMismatch
	}
	if currentAnime.MalID.Int32 != variationAnime.MalID.Int32 {
		return models.LibraryResponse{}, ErrVariationMismatch
	}

	malID := currentAnime.MalID.Int32

	// Get current library entry to preserve status and episodes
	currentLib, err := s.GetLibraryByAnimeID(ctx, userID, currentAnimeID)
	if err != nil {
		return models.LibraryResponse{}, err
	}

	// Get all animes with the same MAL ID
	allVariations, err := s.repo.GetAnimeByMalId(ctx, pgtype.Int4{Int32: malID, Valid: true})
	if err != nil {
		return models.LibraryResponse{}, err
	}

	// Delete ALL library entries for this user that have the same MAL ID
	// This ensures only one library entry per MAL ID at any time
	for _, variation := range allVariations {
		// Try to get library entry for this variation
		_, err := s.GetLibraryByAnimeID(ctx, userID, variation.ID)
		if err == nil {
			// Entry exists, delete it
			err = s.repo.DeleteLibrary(ctx, repository.DeleteLibraryParams{
				UserID:  userID,
				AnimeID: variation.ID,
			})
			if err != nil {
				return models.LibraryResponse{}, err
			}
		} else if !errors.Is(err, ErrLibraryNotFound) {
			// Some other error occurred
			return models.LibraryResponse{}, err
		}
		// If ErrLibraryNotFound, just continue - nothing to delete
	}

	// Insert new entry with variation ID, preserving status and episodes
	err = s.repo.InsertLibrary(ctx, repository.InsertLibraryParams{
		UserID:          userID,
		AnimeID:         variationID,
		Status:          repository.LibraryStatus(currentLib.Status),
		WatchedEpisodes: currentLib.WatchedEpisodes,
	})
	if err != nil {
		return models.LibraryResponse{}, err
	}

	// Return the new library entry
	return s.GetLibraryByAnimeID(ctx, userID, variationID)
}
