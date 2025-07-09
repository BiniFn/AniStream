package anime

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/coeeter/aniways/internal/hianime"
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/myanimelist"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo      *repository.Queries
	refresher *MetadataRefresher
	scraper   *hianime.HianimeScraper
	malClient *myanimelist.Client
}

func New(
	repo *repository.Queries,
	refresher *MetadataRefresher,
	malClient *myanimelist.Client,
) *Service {
	return &Service{
		repo:      repo,
		refresher: refresher,
		malClient: malClient,
		scraper:   hianime.NewHianimeScraper(),
	}
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
		go s.refresher.MaybeRefresh(context.Background(), r.MalID.Int32)
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
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) && !errors.Is(err, pgx.ErrNoRows) {
			log.Printf("metadata lookup error for MAL %d: %v", a.MalID.Int32, err)
			return nil, err
		}
	} else {
		// If metadata is fresh, return it
		if time.Since(m.UpdatedAt.Time) < s.refresher.ttl {
			dto := models.AnimeWithMetadataDto{}.FromRepository(a, m)
			return &dto, nil
		}
	}

	if err := s.refresher.RefreshBlocking(ctx, a.MalID.Int32); err != nil {
		log.Printf("blocking refresh failed: %v", err)
	}
	m, err = s.repo.GetAnimeMetadataByMalId(ctx, a.MalID.Int32)
	if err != nil {
		return nil, err
	}

	dto := models.AnimeWithMetadataDto{}.FromRepository(a, m)

	return &dto, nil
}

func (s *Service) GetAnimeGenres(ctx context.Context) ([]string, error) {
	rows, err := s.repo.GetAllGenres(ctx)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *Service) GetAnimeTrailer(ctx context.Context, id string) (*models.TrailerDto, error) {
	a, err := s.GetAnimeByID(ctx, id)
	if err != nil || a == nil {
		return nil, err
	}
	if a.Metadata.TrailerEmbedURL == "" {
		t, err := s.malClient.GetTrailer(ctx, int(a.MalID))
		if err != nil {
			log.Printf("failed to fetch trailer for MAL ID %d: %v", a.MalID, err)
			return nil, err
		}
		if t == "" {
			log.Printf("no trailer found for MAL ID %d", a.MalID)
			return nil, nil
		}
		a.Metadata.TrailerEmbedURL = t
		params := repository.UpdateAnimeMetadataParams{
			TrailerEmbedUrl:    pgtype.Text{String: a.Metadata.TrailerEmbedURL, Valid: true},
			Description:        pgtype.Text{String: a.Metadata.Description, Valid: true},
			MainPictureUrl:     pgtype.Text{String: a.Metadata.MainPictureURL, Valid: true},
			MediaType:          pgtype.Text{String: a.Metadata.MediaType, Valid: true},
			Rating:             repository.Rating(a.Metadata.Rating),
			AiringStatus:       repository.AiringStatus(a.Metadata.AiringStatus),
			AvgEpisodeDuration: pgtype.Int4{Int32: a.Metadata.AvgEpisodeDuration, Valid: true},
			TotalEpisodes:      pgtype.Int4{Int32: a.Metadata.TotalEpisodes, Valid: true},
			Studio:             pgtype.Text{String: a.Metadata.Studio, Valid: true},
			Rank:               pgtype.Int4{Int32: a.Metadata.Rank, Valid: true},
			Mean:               pgtype.Float8{Float64: a.Metadata.Mean, Valid: true},
			Scoringusers:       pgtype.Int4{Int32: a.Metadata.ScoringUsers, Valid: true},
			Popularity:         pgtype.Int4{Int32: a.Metadata.Popularity, Valid: true},
			AiringStartDate:    pgtype.Text{String: a.Metadata.AiringStartDate, Valid: true},
			AiringEndDate:      pgtype.Text{String: a.Metadata.AiringEndDate, Valid: true},
			Source:             pgtype.Text{String: a.Metadata.Source, Valid: true},
			SeasonYear:         pgtype.Int4{Int32: a.Metadata.SeasonYear, Valid: true},
			Season:             repository.Season(a.Metadata.Season),
			MalID:              a.MalID,
		}
		if err := s.repo.UpdateAnimeMetadata(ctx, params); err != nil {
			log.Printf("failed to update metadata for MAL ID %d: %v", a.MalID, err)
			return nil, err
		}
		log.Printf("updated trailer for MAL ID %d: %s", a.MalID, a.Metadata.TrailerEmbedURL)
	}
	return &models.TrailerDto{Trailer: a.Metadata.TrailerEmbedURL}, nil
}

func (s *Service) GetAnimeEpisodes(ctx context.Context, id string) ([]models.EpisodeDto, error) {
	a, err := s.repo.GetAnimeById(ctx, id)
	if err != nil {
		return nil, err
	}

	episodes, err := s.scraper.GetAnimeEpisodes(ctx, a.HiAnimeID)
	if err != nil {
		log.Printf("failed to fetch episodes for anime ID %s: %v", id, err)
		return nil, err
	}

	if len(episodes) == 0 {
		log.Printf("no episodes found for anime ID %s", id)
		return nil, fmt.Errorf("no episodes found for anime ID %s", id)
	}

	episodeDtos := make([]models.EpisodeDto, len(episodes))
	for i, ep := range episodes {
		episodeDtos[i] = models.EpisodeDto{}.FromScraper(ep)
	}
	return episodeDtos, nil
}
