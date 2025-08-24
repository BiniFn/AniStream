package anime

import (
	"context"
	"fmt"
	"time"

	"github.com/coeeter/aniways/internal/infra/cache"
	"github.com/coeeter/aniways/internal/infra/client/anilist/graphql"
	"github.com/coeeter/aniways/internal/mappers"
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *AnimeService) GetRecentlyUpdatedAnimes(
	ctx context.Context,
	page, size int,
) (models.AnimeListResponse, error) {
	limit, offset, err := utils.ValidatePaginationParams(page, size)
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	rows, err := s.repo.GetRecentlyUpdatedAnimes(ctx, repository.GetRecentlyUpdatedAnimesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	for _, r := range rows {
		s.refresher.Enqueue(r.MalID.Int32)
	}

	total, err := s.repo.GetRecentlyUpdatedAnimesCount(ctx)
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	items := make([]models.AnimeResponse, len(rows))
	for i, a := range rows {
		items[i] = mappers.AnimeFromRepository(a)
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(page, pageSize, total)
	return models.AnimeListResponse{
		PageInfo: pageInfo,
		Items:    items,
	}, nil
}

func (s *AnimeService) GetAnimeGenres(ctx context.Context) ([]string, error) {
	return cache.GetOrFill(ctx, s.redis, "anime_genres", 30*24*time.Hour, func(ctx context.Context) ([]string, error) {
		return s.repo.GetAllGenres(ctx)
	})
}

func (s *AnimeService) SearchAnimes(ctx context.Context, query, genre string, page, size int) (models.AnimeListResponse, error) {
	limit, offset, err := utils.ValidatePaginationParams(page, size)
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	rows, err := s.repo.SearchAnimes(ctx, repository.SearchAnimesParams{
		Query:  query,
		Genre:  genre,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	for _, r := range rows {
		s.refresher.Enqueue(r.MalID.Int32)
	}

	total, err := s.repo.SearchAnimesCount(ctx, repository.SearchAnimesCountParams{
		Query: query,
	})
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	items := make([]models.AnimeResponse, len(rows))
	for i, a := range rows {
		items[i] = mappers.AnimeFromSearch(a)
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(page, pageSize, total)
	return models.AnimeListResponse{
		PageInfo: pageInfo,
		Items:    items,
	}, nil
}

func (s *AnimeService) convertStringToSeason(season string) repository.Season {
	switch season {
	case "winter":
		return repository.SeasonWinter
	case "spring":
		return repository.SeasonSpring
	case "summer":
		return repository.SeasonSummer
	case "fall":
		return repository.SeasonFall
	default:
		return repository.SeasonUnknown
	}
}

func (s *AnimeService) GetAnimeBySeasonAndYear(
	ctx context.Context,
	season string,
	year int32,
	page, size int,
) (models.AnimeListResponse, error) {
	limit, offset, err := utils.ValidatePaginationParams(page, size)
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	rows, err := s.repo.GetAnimeBySeasonAndYear(ctx, repository.GetAnimeBySeasonAndYearParams{
		Season:     s.convertStringToSeason(season),
		SeasonYear: year,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	total, err := s.repo.GetAnimeBySeasonAndYearCount(ctx, repository.GetAnimeBySeasonAndYearCountParams{
		Season:     s.convertStringToSeason(season),
		SeasonYear: year,
	})
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	for _, r := range rows {
		s.refresher.Enqueue(r.MalID.Int32)
	}

	items := make([]models.AnimeResponse, len(rows))
	for i, a := range rows {
		items[i] = mappers.AnimeFromRepository(a)
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(page, pageSize, total)
	return models.AnimeListResponse{
		PageInfo: pageInfo,
		Items:    items,
	}, nil
}

func (s *AnimeService) GetAnimeByYear(
	ctx context.Context,
	year int32,
	page, size int,
) (models.AnimeListResponse, error) {
	limit, offset, err := utils.ValidatePaginationParams(page, size)
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	rows, err := s.repo.GetAnimeByYear(ctx, repository.GetAnimeByYearParams{
		SeasonYear: year,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	total, err := s.repo.GetAnimeByYearCount(ctx, year)
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	for _, r := range rows {
		s.refresher.Enqueue(r.MalID.Int32)
	}

	items := make([]models.AnimeResponse, len(rows))
	for i, a := range rows {
		items[i] = mappers.AnimeFromRepository(a)
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(page, pageSize, total)
	return models.AnimeListResponse{
		PageInfo: pageInfo,
		Items:    items,
	}, nil
}

func (s *AnimeService) GetAnimeBySeason(
	ctx context.Context,
	season string,
	page, size int,
) (models.AnimeListResponse, error) {
	limit, offset, err := utils.ValidatePaginationParams(page, size)
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	rows, err := s.repo.GetAnimeBySeason(ctx, repository.GetAnimeBySeasonParams{
		Season: s.convertStringToSeason(season),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	total, err := s.repo.GetAnimeBySeasonCount(ctx, s.convertStringToSeason(season))
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	for _, r := range rows {
		s.refresher.Enqueue(r.MalID.Int32)
	}

	items := make([]models.AnimeResponse, len(rows))
	for i, a := range rows {
		items[i] = mappers.AnimeFromRepository(a)
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(page, pageSize, total)
	return models.AnimeListResponse{
		PageInfo: pageInfo,
		Items:    items,
	}, nil
}

func (s *AnimeService) GetRandomAnime(ctx context.Context) (models.AnimeResponse, error) {
	data, err := s.repo.GetRandomAnime(ctx)
	if err != nil {
		return models.AnimeResponse{}, fmt.Errorf("failed to fetch random anime: %w", err)
	}

	s.refresher.Enqueue(data.MalID.Int32)

	return mappers.AnimeFromRepository(data), nil
}

func (s *AnimeService) GetRandomAnimeByGenre(ctx context.Context, genre string) (models.AnimeResponse, error) {
	if genre == "" {
		return models.AnimeResponse{}, fmt.Errorf("genre is required")
	}

	data, err := s.repo.GetRandomAnimeByGenre(ctx, pgtype.Text{String: genre, Valid: true})
	if err != nil {
		return models.AnimeResponse{}, fmt.Errorf("failed to fetch random anime by genre %s: %w", genre, err)
	}

	s.refresher.Enqueue(data.MalID.Int32)

	return mappers.AnimeFromRepository(data), nil
}

func (s *AnimeService) GetAnimesByGenre(ctx context.Context, genre string, page, size int) (models.AnimeListResponse, error) {
	limit, offset, err := utils.ValidatePaginationParams(page, size)
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	rows, err := s.repo.GetAnimeByGenre(ctx, repository.GetAnimeByGenreParams{
		Genre:  pgtype.Text{String: genre, Valid: true},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	for _, r := range rows {
		s.refresher.Enqueue(r.MalID.Int32)
	}

	total, err := s.repo.GetAnimeByGenreCount(ctx, pgtype.Text{String: genre, Valid: true})
	if err != nil {
		return models.AnimeListResponse{}, err
	}

	items := make([]models.AnimeResponse, len(rows))
	for i, a := range rows {
		items[i] = mappers.AnimeFromRepository(a)
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(page, pageSize, total)
	return models.AnimeListResponse{
		PageInfo: pageInfo,
		Items:    items,
	}, nil
}

func fetchAnimeDtosFromAnilist[T any](
	ctx context.Context,
	srv *AnimeService,
	media []T,
	getMalID func(T) int,
) ([]repository.Anime, map[int32]repository.Anime, error) {
	if len(media) == 0 {
		return nil, nil, nil
	}

	ids := make([]int32, 0, len(media))
	for _, m := range media {
		if id := getMalID(m); id > 0 {
			ids = append(ids, int32(id))
		}
	}

	if len(ids) == 0 {
		return nil, nil, nil
	}

	rows, err := srv.repo.GetAnimesByMalIds(ctx, ids)
	if err != nil {
		return nil, nil, err
	}

	rowMap := make(map[int32]repository.Anime, len(rows))
	for _, r := range rows {
		if r.MalID.Valid {
			rowMap[r.MalID.Int32] = r
			srv.refresher.Enqueue(r.MalID.Int32)
		}
	}

	return rows, rowMap, nil
}

func (s *AnimeService) GetSeasonalAnimes(ctx context.Context) (models.SeasonalAnimeListResponse, error) {
	return cache.GetOrFill(ctx, s.redis, "seasonal_animes", 30*24*time.Hour, func(ctx context.Context) (models.SeasonalAnimeListResponse, error) {
		now := time.Now()
		ref := now.AddDate(0, -1, 0)
		year := ref.Year()

		var season string
		switch ref.Month() {
		case time.January, time.February, time.March:
			season = "WINTER"
		case time.April, time.May, time.June:
			season = "SPRING"
		case time.July, time.August, time.September:
			season = "SUMMER"
		default:
			season = "FALL"
		}
		animes, err := s.anilistClient.GetSeasonalMedia(ctx, year, season)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch seasonal animes for year %d season %s: %w", year, season, err)
		}

		_, dtoMap, err := fetchAnimeDtosFromAnilist(ctx, s, animes.Page.Media, func(m graphql.GetSeasonalAnimePageMedia) int {
			return m.IdMal
		})
		if err != nil {
			return nil, fmt.Errorf("failed to map seasonal animes: %w", err)
		}

		seasonalAnimes := make([]models.SeasonalAnimeResponse, 0, len(animes.Page.Media))
		for _, a := range animes.Page.Media {
			id := int32(a.IdMal)
			dbAnime, ok := dtoMap[id]
			if !ok {
				continue
			}
			startDate := time.Date(
				a.StartDate.Year,
				time.Month(a.StartDate.Month),
				a.StartDate.Day,
				0, 0, 0, 0,
				time.UTC,
			)
			seasonalAnimes = append(seasonalAnimes, models.SeasonalAnimeResponse{
				ID:             dbAnime.ID,
				BannerImageURL: a.BannerImage,
				Description:    a.Description,
				Type:           string(a.Type),
				StartDate:      startDate.UnixMilli(),
				Episodes:       int32(a.Episodes),
				Anime:          mappers.AnimeFromRepository(dbAnime),
			})
		}

		return seasonalAnimes, nil
	})
}

func (s *AnimeService) GetTrendingAnimes(ctx context.Context) (models.TrendingAnimeListResponse, error) {
	return cache.GetOrFill(ctx, s.redis, "trending_animes", 24*time.Hour, func(ctx context.Context) (models.TrendingAnimeListResponse, error) {
		animes, err := s.anilistClient.GetTrendingAnime(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch trending animes: %w", err)
		}

		_, dtoMap, err := fetchAnimeDtosFromAnilist(ctx, s, animes.Page.Media, func(m graphql.GetTrendingAnimePageMedia) int {
			return m.IdMal
		})
		if err != nil {
			return nil, fmt.Errorf("failed to map trending animes: %w", err)
		}

		trendingAnimes := make([]models.AnimeResponse, 0, len(animes.Page.Media))
		for _, a := range animes.Page.Media {
			id := int32(a.IdMal)
			dbAnime, ok := dtoMap[id]
			if !ok {
				continue
			}
			trendingAnimes = append(trendingAnimes, mappers.AnimeFromRepository(dbAnime))
		}

		return trendingAnimes, nil
	})
}

func (s *AnimeService) GetPopularAnimes(ctx context.Context) (models.PopularAnimeListResponse, error) {
	return cache.GetOrFill(ctx, s.redis, "popular_animes", 24*time.Hour, func(ctx context.Context) (models.PopularAnimeListResponse, error) {
		animes, err := s.anilistClient.GetPopularAnime(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch popular animes: %w", err)
		}

		_, dtoMap, err := fetchAnimeDtosFromAnilist(ctx, s, animes.Page.Media, func(m graphql.GetPopularAnimePageMedia) int {
			return m.IdMal
		})
		if err != nil {
			return nil, fmt.Errorf("failed to map popular animes: %w", err)
		}

		popularAnimes := make([]models.AnimeResponse, 0, len(animes.Page.Media))
		for _, a := range animes.Page.Media {
			id := int32(a.IdMal)
			dbAnime, ok := dtoMap[id]
			if !ok {
				continue
			}
			popularAnimes = append(popularAnimes, mappers.AnimeFromRepository(dbAnime))
		}

		return popularAnimes, nil
	})
}
