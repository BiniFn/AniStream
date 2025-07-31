package anime

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/coeeter/aniways/internal/infra/cache"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrAnimeNotFound = errors.New("anime not found")

func (s *AnimeService) GetAnimeByID(
	ctx context.Context,
	id string,
) (*AnimeWithMetadataDto, error) {
	a, err := s.repo.GetAnimeById(ctx, id)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAnimeNotFound
	}
	if err != nil {
		return nil, err
	}

	m, err := s.repo.GetAnimeMetadataByMalId(ctx, a.MalID.Int32)
	switch {
	case err == nil && time.Since(m.UpdatedAt.Time) < s.refresher.ttl:
		dto := AnimeWithMetadataDto{}.FromRepository(a, m)
		return &dto, nil

	case err != nil && !errors.Is(err, sql.ErrNoRows) && !errors.Is(err, pgx.ErrNoRows):
		return nil, err
	}

	// if error and no existing metadata return error
	if err := s.refresher.RefreshBlocking(ctx, a.MalID.Int32); err != nil && m.MalID == 0 {
		return nil, err
	}

	m, err = s.repo.GetAnimeMetadataByMalId(ctx, a.MalID.Int32)
	if err != nil {
		return nil, err
	}

	dto := AnimeWithMetadataDto{}.FromRepository(a, m)

	return &dto, nil
}

func (s *AnimeService) GetAnimeTrailer(ctx context.Context, id string) (*TrailerDto, error) {
	a, err := s.GetAnimeByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAnimeNotFound
	}
	if err != nil {
		return nil, err
	}

	if a.Metadata.TrailerEmbedURL != "" {
		return &TrailerDto{Trailer: a.Metadata.TrailerEmbedURL}, nil
	}

	t, err := s.malClient.GetTrailer(ctx, int(a.MalID))
	if err != nil || t == "" {
		return nil, fmt.Errorf("failed to fetch trailer for MAL ID %d: %v", a.MalID, err)
	}

	a.Metadata.TrailerEmbedURL = t
	params := repository.UpdateAnimeMetadataTrailerParams{
		TrailerEmbedUrl: pgtype.Text{String: t, Valid: true},
		MalID:           a.MalID,
	}
	if err := s.repo.UpdateAnimeMetadataTrailer(ctx, params); err != nil {
		return nil, fmt.Errorf("failed to update metadata for MAL ID %d: %v", a.MalID, err)
	}

	return &TrailerDto{Trailer: t}, nil
}

func (s *AnimeService) GetAnimeBanner(ctx context.Context, id string) (BannerDto, error) {
	banner, err := cache.GetOrFill(ctx, s.redis, fmt.Sprintf("anime_banner:%s", id), 30*24*time.Hour, func(ctx context.Context) (string, error) {
		a, err := s.repo.GetAnimeById(ctx, id)
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return "", ErrAnimeNotFound
		}
		if err != nil {
			return "", fmt.Errorf("failed to fetch anime by ID %s: %v", id, err)
		}

		anime, err := s.anilistClient.GetAnimeDetails(ctx, int(a.MalID.Int32))
		if err != nil {
			return "", fmt.Errorf("failed to fetch anime details from Anilist for MAL ID %d: %v", a.MalID.Int32, err)
		}

		return anime.Media.BannerImage, nil
	})

	if err != nil {
		return BannerDto{}, err
	}

	return BannerDto{URL: banner}, nil
}

func (s *AnimeService) GetAnimeRelations(ctx context.Context, id string) (RelationsDto, error) {
	key := fmt.Sprintf("anime_relations:%s", id)

	return cache.GetOrFill(ctx, s.redis, key, 7*24*time.Hour, func(ctx context.Context) (RelationsDto, error) {
		a, err := s.repo.GetAnimeById(ctx, id)
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return RelationsDto{}, ErrAnimeNotFound
		}
		if err != nil {
			return RelationsDto{}, fmt.Errorf("failed to fetch anime by ID %s: %v", id, err)
		}

		if !a.MalID.Valid || a.MalID.Int32 <= 0 {
			return RelationsDto{}, fmt.Errorf("invalid MAL ID for anime ID %s", id)
		}

		malID := int(a.MalID.Int32)
		fr, err := s.shikimoriClient.GetAnimeFranchise(ctx, malID)
		if err != nil {
			return RelationsDto{}, fmt.Errorf("failed to fetch franchise for MAL ID %d: %v", malID, err)
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
			return RelationsDto{}, fmt.Errorf("failed to fetch related animes: %w", err)
		}

		dtoMap := make(map[int32]AnimeDto, len(rows))
		for _, r := range rows {
			dtoMap[r.MalID.Int32] = AnimeDto{}.FromRepository(r)
			s.refresher.Enqueue(r.MalID.Int32)
		}

		sliceDto := func(ids []int) []AnimeDto {
			out := make([]AnimeDto, 0, len(ids))
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

		var relations RelationsDto

		if len(watchIDs) > 1 && slices.Contains(watchIDs, malID) {
			relations = RelationsDto{
				WatchOrder: sliceDto(watchIDs),
				Related:    sliceDto(reverse(relatedIDs)),
			}
		} else {
			relations = RelationsDto{
				WatchOrder: []AnimeDto{},
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
}
