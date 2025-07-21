package anime

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coeeter/aniways/internal/client/anilist/graphql"
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *AnimeService) GetRecentlyUpdatedAnimes(
	ctx context.Context,
	page, size int,
) (models.Pagination[models.AnimeDto], error) {
	limit, offset, err := utils.ValidatePaginationParams(page, size)
	if err != nil {
		return models.Pagination[models.AnimeDto]{}, err
	}

	rows, err := s.repo.GetRecentlyUpdatedAnimes(ctx, repository.GetRecentlyUpdatedAnimesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return models.Pagination[models.AnimeDto]{}, err
	}

	for _, r := range rows {
		s.refresher.Enqueue(r.MalID.Int32)
	}

	total, err := s.repo.GetRecentlyUpdatedAnimesCount(ctx)
	if err != nil {
		return models.Pagination[models.AnimeDto]{}, err
	}

	items := make([]models.AnimeDto, len(rows))
	for i, a := range rows {
		items[i] = models.AnimeDto{}.FromRepository(a)
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(page, pageSize, total)
	return models.Pagination[models.AnimeDto]{
		PageInfo: pageInfo,
		Items:    items,
	}, nil
}

func (s *AnimeService) GetAnimeGenres(ctx context.Context) ([]string, error) {
	var genres []string

	_, err := s.redis.GetOrFill(ctx, "anime_genres", &genres, 30*24*time.Hour, func(ctx context.Context) (any, error) {
		return s.repo.GetAllGenres(ctx)
	})

	if err != nil {
		log.Printf("failed to get anime genres from cache: %v", err)
		return nil, err
	}

	return genres, nil
}

func (s *AnimeService) SearchAnimes(ctx context.Context, query, genre string, page, size int) (models.Pagination[models.AnimeDto], error) {
	limit, offset, err := utils.ValidatePaginationParams(page, size)
	if err != nil {
		return models.Pagination[models.AnimeDto]{}, err
	}

	rows, err := s.repo.SearchAnimes(ctx, repository.SearchAnimesParams{
		Query:  query,
		Genre:  genre,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return models.Pagination[models.AnimeDto]{}, err
	}

	for _, r := range rows {
		s.refresher.Enqueue(r.MalID.Int32)
	}

	total, err := s.repo.SearchAnimesCount(ctx, repository.SearchAnimesCountParams{
		Query: query,
	})
	if err != nil {
		return models.Pagination[models.AnimeDto]{}, err
	}

	items := make([]models.AnimeDto, len(rows))
	for i, a := range rows {
		items[i] = models.AnimeDto{}.FromSearch(a)
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(page, pageSize, total)
	return models.Pagination[models.AnimeDto]{
		PageInfo: pageInfo,
		Items:    items,
	}, nil
}

func (s *AnimeService) GetRandomAnime(ctx context.Context) (models.AnimeDto, error) {
	data, err := s.repo.GetRandomAnime(ctx)
	if err != nil {
		log.Printf("failed to fetch random anime: %v", err)
		return models.AnimeDto{}, err
	}

	s.refresher.Enqueue(data.MalID.Int32)

	return models.AnimeDto{}.FromRepository(data), nil
}

func (s *AnimeService) GetRandomAnimeByGenre(ctx context.Context, genre string) (models.AnimeDto, error) {
	if genre == "" {
		return models.AnimeDto{}, fmt.Errorf("genre is required")
	}

	data, err := s.repo.GetRandomAnimeByGenre(ctx, pgtype.Text{String: genre, Valid: true})
	if err != nil {
		log.Printf("failed to fetch random anime by genre %s: %v", genre, err)
		return models.AnimeDto{}, err
	}

	s.refresher.Enqueue(data.MalID.Int32)

	return models.AnimeDto{}.FromRepository(data), nil
}

func (s *AnimeService) GetAnimesByGenre(ctx context.Context, genre string, page, size int) (models.Pagination[models.AnimeDto], error) {
	limit, offset, err := utils.ValidatePaginationParams(page, size)
	if err != nil {
		return models.Pagination[models.AnimeDto]{}, err
	}

	rows, err := s.repo.GetAnimeByGenre(ctx, repository.GetAnimeByGenreParams{
		Genre:  pgtype.Text{String: genre, Valid: true},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return models.Pagination[models.AnimeDto]{}, err
	}

	for _, r := range rows {
		s.refresher.Enqueue(r.MalID.Int32)
	}

	total, err := s.repo.GetAnimeByGenreCount(ctx, pgtype.Text{String: genre, Valid: true})
	if err != nil {
		return models.Pagination[models.AnimeDto]{}, err
	}

	items := make([]models.AnimeDto, len(rows))
	for i, a := range rows {
		items[i] = models.AnimeDto{}.FromRepository(a)
	}

	pageSize := int64(limit)
	pageInfo := utils.PageInfo(page, pageSize, total)
	return models.Pagination[models.AnimeDto]{
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

func (s *AnimeService) GetSeasonalAnimes(ctx context.Context) ([]models.SeasonalAnimeDto, error) {
	var cachedAnimes []models.SeasonalAnimeDto

	_, err := s.redis.GetOrFill(ctx, "seasonal_animes", &cachedAnimes, 30*24*time.Hour, func(ctx context.Context) (any, error) {
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
			log.Printf("failed to fetch seasonal animes for year %d season %s: %v", year, season, err)
			return nil, err
		}

		_, dtoMap, err := fetchAnimeDtosFromAnilist(ctx, s, animes.Page.Media, func(m graphql.GetSeasonalAnimePageMedia) int {
			return m.IdMal
		})
		if err != nil {
			log.Printf("failed to map seasonal animes: %v", err)
			return nil, err
		}

		seasonalAnimes := make([]models.SeasonalAnimeDto, 0, len(animes.Page.Media))
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
			seasonalAnimes = append(seasonalAnimes, models.SeasonalAnimeDto{
				ID:             dbAnime.ID,
				BannerImageURL: a.BannerImage,
				Description:    a.Description,
				Type:           string(a.Type),
				StartDate:      startDate.UnixMilli(),
				Episodes:       int32(a.Episodes),
				Anime:          models.AnimeDto{}.FromRepository(dbAnime),
			})
		}

		return seasonalAnimes, nil
	})

	if err != nil {
		log.Printf("failed to fetch seasonal animes: %v", err)
		return nil, err
	}

	return cachedAnimes, nil
}

func (s *AnimeService) GetTrendingAnimes(ctx context.Context) ([]models.AnimeDto, error) {
	var cachedAnimes []models.AnimeDto

	_, err := s.redis.GetOrFill(ctx, "trending_animes", &cachedAnimes, 24*time.Hour, func(ctx context.Context) (any, error) {
		animes, err := s.anilistClient.GetTrendingAnime(ctx)
		if err != nil {
			log.Printf("failed to fetch trending animes: %v", err)
			return nil, err
		}

		_, dtoMap, err := fetchAnimeDtosFromAnilist(ctx, s, animes.Page.Media, func(m graphql.GetTrendingAnimePageMedia) int {
			return m.IdMal
		})
		if err != nil {
			log.Printf("failed to map trending animes: %v", err)
			return nil, err
		}

		trendingAnimes := make([]models.AnimeDto, 0, len(animes.Page.Media))
		for _, a := range animes.Page.Media {
			id := int32(a.IdMal)
			dbAnime, ok := dtoMap[id]
			if !ok {
				continue
			}
			trendingAnimes = append(trendingAnimes, models.AnimeDto{}.FromRepository(dbAnime))
		}

		return trendingAnimes, nil
	})

	if err != nil {
		log.Printf("failed to get trending animes from cache: %v", err)
		return nil, err
	}

	return cachedAnimes, nil
}

func (s *AnimeService) GetPopularAnimes(ctx context.Context) ([]models.AnimeDto, error) {
	var cachedAnimes []models.AnimeDto
	_, err := s.redis.GetOrFill(ctx, "popular_animes", &cachedAnimes, 24*time.Hour, func(ctx context.Context) (any, error) {
		animes, err := s.anilistClient.GetPopularAnime(ctx)
		if err != nil {
			log.Printf("failed to fetch popular animes: %v", err)
			return nil, err
		}

		_, dtoMap, err := fetchAnimeDtosFromAnilist(ctx, s, animes.Page.Media, func(m graphql.GetPopularAnimePageMedia) int {
			return m.IdMal
		})
		if err != nil {
			log.Printf("failed to map popular animes: %v", err)
			return nil, err
		}

		popularAnimes := make([]models.AnimeDto, 0, len(animes.Page.Media))
		for _, a := range animes.Page.Media {
			id := int32(a.IdMal)
			dbAnime, ok := dtoMap[id]
			if !ok {
				continue
			}
			popularAnimes = append(popularAnimes, models.AnimeDto{}.FromRepository(dbAnime))
		}

		return popularAnimes, nil
	})

	if err != nil {
		log.Printf("failed to get popular animes from cache: %v", err)
		return nil, err
	}

	return cachedAnimes, nil
}
