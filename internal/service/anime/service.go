package anime

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/client/anilist"
	"github.com/coeeter/aniways/internal/client/hianime"
	"github.com/coeeter/aniways/internal/client/myanimelist"
	"github.com/coeeter/aniways/internal/client/shikimori"
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo            *repository.Queries
	refresher       *MetadataRefresher
	scraper         *hianime.HianimeScraper
	malClient       *myanimelist.Client
	anilistClient   *anilist.Client
	shikimoriClient *shikimori.Client
	redis           *cache.RedisClient
}

func New(
	repo *repository.Queries,
	refresher *MetadataRefresher,
	malClient *myanimelist.Client,
	anilistClient *anilist.Client,
	shikimoriClient *shikimori.Client,
	redis *cache.RedisClient,
) *Service {
	return &Service{
		repo:            repo,
		refresher:       refresher,
		malClient:       malClient,
		anilistClient:   anilistClient,
		shikimoriClient: shikimoriClient,
		scraper:         hianime.NewHianimeScraper(),
		redis:           redis,
	}
}

func (s *Service) GetRecentlyUpdatedAnimes(
	ctx context.Context,
	page, size int,
) (models.Pagination[models.AnimeDto], error) {
	offset := int32((page - 1) * size)
	limit := int32(size)

	if page < 1 || size < 1 {
		return models.Pagination[models.AnimeDto]{}, fmt.Errorf("invalid pagination parameters: page=%d, size=%d", page, size)
	}

	if size > 100 {
		return models.Pagination[models.AnimeDto]{}, fmt.Errorf("size too large: maximum is 100, got %d", size)
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
	var cachedEpisodes []models.EpisodeDto

	_, err := s.redis.GetOrFill(ctx, fmt.Sprintf("anime_episodes:%s", id), &cachedEpisodes, 7*24*time.Hour, func(ctx context.Context) (any, error) {
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
	})

	if err != nil {
		log.Printf("failed to get anime episodes from cache: %v", err)
		return nil, err
	}

	return cachedEpisodes, nil
}

func (s *Service) SearchAnimes(ctx context.Context, query, genre string, page, size int) (models.Pagination[models.AnimeDto], error) {
	if page < 1 || size < 1 {
		return models.Pagination[models.AnimeDto]{}, fmt.Errorf("invalid pagination parameters: page=%d, size=%d", page, size)
	}

	if size > 100 {
		return models.Pagination[models.AnimeDto]{}, fmt.Errorf("size too large: maximum is 100, got %d", size)
	}

	offset := int32((page - 1) * size)
	limit := int32(size)

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

func (s *Service) GetServersForEpisode(ctx context.Context, id, episodeID string) (models.ServerDto, error) {
	var cachedServers models.ServerDto

	_, err := s.redis.GetOrFill(ctx, fmt.Sprintf("anime_servers:%s:episode:%s", id, episodeID), &cachedServers, 24*time.Hour, func(ctx context.Context) (any, error) {
		a, err := s.repo.GetAnimeById(ctx, id)
		if err != nil {
			log.Printf("failed to fetch anime by ID %s: %v", id, err)
			return models.ServerDto{}, err
		}

		servers, err := s.scraper.GetEpisodeServers(ctx, a.HiAnimeID, episodeID)
		if err != nil {
			log.Printf("failed to fetch servers for anime ID %s episode %s: %v", id, episodeID, err)
			return models.ServerDto{}, err
		}

		if len(servers) == 0 {
			log.Printf("no servers found for anime ID %s episode %s", id, episodeID)
			return models.ServerDto{}, fmt.Errorf("no servers found for anime ID %s episode %s", id, episodeID)
		}

		serverDto := models.ServerDto{}.FromScraper(servers)
		return serverDto, nil
	})

	if err != nil {
		log.Printf("failed to get servers for episode from cache: %v", err)
		return models.ServerDto{}, err
	}

	return cachedServers, nil
}

func (s *Service) GetStreamingData(ctx context.Context, serverID, streamType, serverName string) (models.StreamingDataDto, error) {
	if serverID == "" || streamType == "" || serverName == "" {
		return models.StreamingDataDto{}, fmt.Errorf("serverID, streamType, and serverName are required")
	}

	var cachedData models.StreamingDataDto
	key := fmt.Sprintf("streaming_data:%s:%s:%s", serverID, streamType, serverName)
	_, err := s.redis.GetOrFill(ctx, key, &cachedData, 24*time.Hour, func(ctx context.Context) (any, error) {
		data, err := s.scraper.GetStreamingData(ctx, serverID, streamType, serverName)
		if err != nil {
			log.Printf("failed to fetch streaming data for server %s type %s name %s: %v", serverID, streamType, serverName, err)
			return models.StreamingDataDto{}, err
		}

		dto := models.StreamingDataDto{}.FromScraper(data)
		return dto, nil
	})

	if err != nil {
		log.Printf("failed to get streaming data from cache: %v", err)
		return models.StreamingDataDto{}, err
	}

	return cachedData, nil
}

func (s *Service) GetRandomAnime(ctx context.Context) (models.AnimeDto, error) {
	data, err := s.repo.GetRandomAnime(ctx)
	if err != nil {
		log.Printf("failed to fetch random anime: %v", err)
		return models.AnimeDto{}, err
	}

	s.refresher.Enqueue(data.MalID.Int32)

	return models.AnimeDto{}.FromRepository(data), nil
}

func (s *Service) GetRandomAnimeByGenre(ctx context.Context, genre string) (models.AnimeDto, error) {
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

func (s *Service) GetSeasonalAnimes(ctx context.Context) ([]models.SeasonalAnimeDto, error) {
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

		dbAnimes := make([]repository.Anime, len(animes.Page.Media))
		for i, a := range animes.Page.Media {
			d, err := s.repo.GetAnimeByMalId(ctx, pgtype.Int4{Int32: int32(a.IdMal), Valid: true})
			if err != nil && !errors.Is(err, sql.ErrNoRows) && !errors.Is(err, pgx.ErrNoRows) {
				log.Printf("failed to fetch anime by MAL ID %d: %v", a.IdMal, err)
				continue
			}
			dbAnimes[i] = d
			s.refresher.Enqueue(d.MalID.Int32)
		}

		seasonalAnimes := make([]models.SeasonalAnimeDto, len(animes.Page.Media))
		for i, a := range animes.Page.Media {
			dbAnime := dbAnimes[i]
			startDate := time.Date(
				a.StartDate.Year,
				time.Month(a.StartDate.Month),
				a.StartDate.Day,
				0, 0, 0, 0,
				time.UTC,
			)
			seasonalAnimes[i] = models.SeasonalAnimeDto{
				ID:             dbAnime.ID,
				BannerImageURL: a.BannerImage,
				Description:    a.Description,
				Type:           string(a.Type),
				StartDate:      startDate.UnixMilli(),
				Episodes:       int32(a.Episodes),
				Anime:          models.AnimeDto{}.FromRepository(dbAnime),
			}
		}

		return seasonalAnimes, nil
	})

	if err != nil {
		log.Printf("failed to fetch seasonal animes: %v", err)
		return nil, err
	}

	return cachedAnimes, nil
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

func (s *Service) GetTrendingAnimes(ctx context.Context) ([]models.AnimeDto, error) {
	var cachedAnimes []models.AnimeDto

	_, err := s.redis.GetOrFill(ctx, "trending_animes", &cachedAnimes, 24*time.Hour, func(ctx context.Context) (any, error) {
		animes, err := s.anilistClient.GetTrendingAnime(ctx)
		if err != nil {
			log.Printf("failed to fetch trending animes: %v", err)
			return nil, err
		}

		trendingAnimes := make([]models.AnimeDto, len(animes.Page.Media))
		for i, a := range animes.Page.Media {
			d, err := s.repo.GetAnimeByMalId(ctx, pgtype.Int4{Int32: int32(a.IdMal), Valid: true})
			if err != nil && !errors.Is(err, sql.ErrNoRows) && !errors.Is(err, pgx.ErrNoRows) {
				log.Printf("failed to fetch anime by MAL ID %d: %v", a.IdMal, err)
				continue
			}
			trendingAnimes[i] = models.AnimeDto{}.FromRepository(d)
			s.refresher.Enqueue(d.MalID.Int32)
		}

		return trendingAnimes, nil
	})

	if err != nil {
		log.Printf("failed to get trending animes from cache: %v", err)
		return nil, err
	}

	return cachedAnimes, nil
}

func (s *Service) GetPopularAnimes(ctx context.Context) ([]models.AnimeDto, error) {
	var cachedAnimes []models.AnimeDto
	_, err := s.redis.GetOrFill(ctx, "popular_animes", &cachedAnimes, 24*time.Hour, func(ctx context.Context) (any, error) {
		animes, err := s.anilistClient.GetPopularAnime(ctx)
		if err != nil {
			log.Printf("failed to fetch popular animes: %v", err)
			return nil, err
		}

		popularAnimes := make([]models.AnimeDto, len(animes.Page.Media))
		for i, a := range animes.Page.Media {
			d, err := s.repo.GetAnimeByMalId(ctx, pgtype.Int4{Int32: int32(a.IdMal), Valid: true})
			if err != nil && !errors.Is(err, sql.ErrNoRows) && !errors.Is(err, pgx.ErrNoRows) {
				log.Printf("failed to fetch anime by MAL ID %d: %v", a.IdMal, err)
				continue
			}
			popularAnimes[i] = models.AnimeDto{}.FromRepository(d)
			s.refresher.Enqueue(d.MalID.Int32)
		}

		return popularAnimes, nil
	})

	if err != nil {
		log.Printf("failed to get popular animes from cache: %v", err)
		return nil, err
	}

	return cachedAnimes, nil
}

func (s *Service) GetAnimesByGenre(ctx context.Context, genre string, page, size int) (models.Pagination[models.AnimeDto], error) {
	if page < 1 || size < 1 {
		return models.Pagination[models.AnimeDto]{}, fmt.Errorf("invalid pagination parameters: page=%d, size=%d", page, size)
	}

	if size > 100 {
		return models.Pagination[models.AnimeDto]{}, fmt.Errorf("size too large: maximum is 100, got %d", size)
	}

	offset := int32((page - 1) * size)
	limit := int32(size)

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
				} else {
					log.Printf("no anime found for MAL ID %d", id)
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
