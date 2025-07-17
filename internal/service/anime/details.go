package anime

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

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
		params := repository.UpdateAnimeMetadataTrailerParams{
			TrailerEmbedUrl: pgtype.Text{String: t, Valid: len(t) > 0},
			MalID:           a.MalID,
		}
		if err := s.repo.UpdateAnimeMetadataTrailer(ctx, params); err != nil {
			log.Printf("failed to update metadata for MAL ID %d: %v", a.MalID, err)
			return nil, err
		}
		log.Printf("updated trailer for MAL ID %d: %s", a.MalID, a.Metadata.TrailerEmbedURL)
	}
	return &models.TrailerDto{Trailer: a.Metadata.TrailerEmbedURL}, nil
}

func (s *Service) GetAnimeBanner(ctx context.Context, id string) (string, error) {
	var cachedBanner string
	_, err := s.redis.GetOrFill(ctx, fmt.Sprintf("anime_banner:%s", id), &cachedBanner, 30*24*time.Hour, func(ctx context.Context) (any, error) {
		a, err := s.repo.GetAnimeById(ctx, id)
		if err != nil {
			log.Printf("failed to fetch anime by ID %s: %v", id, err)
			return "", err
		}

		anime, err := s.anilistClient.GetAnimeDetails(ctx, int(a.MalID.Int32))
		if err != nil {
			log.Printf("failed to fetch anime details from Anilist for MAL ID %d: %v", a.MalID.Int32, err)
			return "", err
		}

		bannerURL := anime.Media.GetBannerImage()
		if bannerURL == "" {
			log.Printf("no banner image found for anime ID %s", id)
			return "", fmt.Errorf("no banner image found for anime ID %s", id)
		}

		return bannerURL, nil
	})

	if err != nil {
		log.Printf("failed to get anime banner from cache: %v", err)
		return "", err
	}

	return cachedBanner, nil
}

func (s *Service) GetAnimeRelations(ctx context.Context, animeID string) (models.RelationsDto, error) {
	cachedKey := fmt.Sprintf("anime_relations:%s", animeID)
	var cachedRelations models.RelationsDto

	_, err := s.redis.GetOrFill(ctx, cachedKey, &cachedRelations, 7*24*time.Hour, func(ctx context.Context) (any, error) {
		anime, err := s.repo.GetAnimeById(ctx, animeID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
				return models.RelationsDto{}, nil
			}
			log.Printf("failed to fetch anime by ID %s: %v", animeID, err)
			return models.RelationsDto{}, err
		}

		if !anime.MalID.Valid || anime.MalID.Int32 <= 0 {
			log.Printf("invalid MAL ID for anime ID %s: %v", animeID, anime.MalID)
			return models.RelationsDto{}, fmt.Errorf("invalid MAL ID for anime ID %s", animeID)
		}

		malID := int(anime.MalID.Int32)
		fr, err := s.shikimoriClient.GetAnimeFranchise(ctx, malID)
		if err != nil {
			log.Printf("failed to fetch franchise for MAL ID %d: %v", malID, err)
			return models.RelationsDto{}, err
		}

		watchIDs := deriveWatchOrder(fr, malID)
		relatedIDs := deriveRelated(fr, malID, watchIDs)
		fullIDs := deriveFullFranchise(fr)

		rows, err := s.repo.GetAnimesByMalIds(ctx, func(ids []int) []int32 {
			out := make([]int32, 0, len(ids))
			for _, id := range ids {
				out = append(out, int32(id))
			}
			return out
		}(fullIDs))
		if err != nil {
			log.Printf("failed to fetch related animes by MAL IDs %v: %v", fullIDs, err)
			return models.RelationsDto{}, fmt.Errorf("failed to fetch related animes: %w", err)
		}

		dtoMap := make(map[int32]models.AnimeDto, len(rows))
		for _, r := range rows {
			dtoMap[r.MalID.Int32] = models.AnimeDto{}.FromRepository(r)
			s.refresher.Enqueue(r.MalID.Int32)
		}

		sliceDto := func(ids []int) []models.AnimeDto {
			out := make([]models.AnimeDto, 0, len(ids))
			for _, id := range ids {
				if dto, ok := dtoMap[int32(id)]; ok {
					out = append(out, dto)
				}
			}
			return out
		}

		reverse := func(ids []int) []int {
			out := make([]int, len(ids))
			for i, id := range ids {
				out[len(ids)-1-i] = id
			}
			return out
		}

		var relations models.RelationsDto

		if len(watchIDs) > 1 && slices.Contains(watchIDs, malID) {
			relations = models.RelationsDto{
				WatchOrder: sliceDto(watchIDs),
				Related:    sliceDto(reverse(relatedIDs)),
			}
		} else {
			relations = models.RelationsDto{
				WatchOrder: []models.AnimeDto{},
				Related: sliceDto(func(ids []int) []int {
					if len(ids) == 1 {
						return []int{}
					}
					return reverse(ids)
				}(fullIDs)),
			}
		}

		return relations, nil
	})

	if err != nil {
		log.Printf("failed to get anime relations from cache: %v", err)
		return models.RelationsDto{}, err
	}

	return cachedRelations, nil
}
