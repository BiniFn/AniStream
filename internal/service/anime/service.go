package anime

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/client/anilist"
	"github.com/coeeter/aniways/internal/client/hianime"
	"github.com/coeeter/aniways/internal/client/myanimelist"
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo          *repository.Queries
	refresher     *MetadataRefresher
	scraper       *hianime.HianimeScraper
	malClient     *myanimelist.Client
	anilistClient *anilist.Client
	redis         *cache.RedisClient
}

func New(
	repo *repository.Queries,
	refresher *MetadataRefresher,
	malClient *myanimelist.Client,
	anilistClient *anilist.Client,
	redis *cache.RedisClient,
) *Service {
	return &Service{
		repo:          repo,
		refresher:     refresher,
		malClient:     malClient,
		anilistClient: anilistClient,
		scraper:       hianime.NewHianimeScraper(),
		redis:         redis,
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
	key := "anime_genres"
	var genres []string
	if ok, err := s.redis.Get(ctx, key, &genres); err != nil {
		log.Printf("failed to get genres from cache: %v", err)
	} else if ok {
		log.Printf("found %d genres in cache", len(genres))
		return genres, nil
	}

	genres, err := s.repo.GetAllGenres(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.redis.Set(ctx, key, genres, 30*24*time.Hour); err != nil {
		log.Printf("failed to cache genres: %v", err)
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
	key := fmt.Sprintf("anime_episodes:%s", id)
	var cachedEpisodes []models.EpisodeDto
	if ok, err := s.redis.Get(ctx, key, &cachedEpisodes); err != nil {
		log.Printf("failed to get episodes from cache: %v", err)
	} else if ok {
		log.Printf("found %d episodes in cache for anime ID %s", len(cachedEpisodes), id)
		return cachedEpisodes, nil
	}

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

	if err := s.redis.Set(ctx, key, episodeDtos, 7*24*time.Hour); err != nil {
		log.Printf("failed to cache episodes for anime ID %s: %v", id, err)
	} else {
		log.Printf("cached %d episodes for anime ID %s", len(episodeDtos), id)
	}

	return episodeDtos, nil
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
		go s.refresher.MaybeRefresh(context.Background(), r.MalID.Int32)
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
	key := fmt.Sprintf("anime_servers:%s:episode:%s", id, episodeID)
	var cachedServers models.ServerDto
	if ok, err := s.redis.Get(ctx, key, &cachedServers); err != nil {
		log.Printf("failed to get servers from cache: %v", err)
	} else if ok {
		log.Printf("found %d servers in cache for anime ID %s episode %s", len(cachedServers.Sub)+len(cachedServers.Dub)+len(cachedServers.Raw), id, episodeID)
		return cachedServers, nil
	}

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

	if err := s.redis.Set(ctx, key, serverDto, 24*time.Hour); err != nil {
		log.Printf("failed to cache servers for anime ID %s episode %s: %v", id, episodeID, err)
	} else {
		log.Printf("cached %d servers for anime ID %s episode %s", len(serverDto.Sub)+len(serverDto.Dub)+len(serverDto.Raw), id, episodeID)
	}

	return serverDto, nil
}

func (s *Service) GetStreamingData(ctx context.Context, serverID, streamType, serverName string) (models.StreamingDataDto, error) {
	if serverID == "" || streamType == "" || serverName == "" {
		return models.StreamingDataDto{}, fmt.Errorf("serverID, streamType, and serverName are required")
	}

	key := fmt.Sprintf("streaming_data:%s:%s:%s", serverID, streamType, serverName)
	var cachedData models.StreamingDataDto
	if ok, err := s.redis.Get(ctx, key, &cachedData); err != nil {
		log.Printf("failed to get streaming data from cache: %v", err)
	} else if ok {
		log.Printf("found streaming data in cache for server %s type %s name %s", serverID, streamType, serverName)
		return cachedData, nil
	}

	data, err := s.scraper.GetStreamingData(ctx, serverID, streamType, serverName)
	if err != nil {
		log.Printf("failed to fetch streaming data for server %s type %s name %s: %v", serverID, streamType, serverName, err)
		return models.StreamingDataDto{}, err
	}

	dto := models.StreamingDataDto{}.FromScraper(data)

	if err := s.redis.Set(ctx, key, dto, 24*time.Hour); err != nil {
		log.Printf("failed to cache streaming data for server %s type %s name %s: %v", serverID, streamType, serverName, err)
	} else {
		log.Printf("cached streaming data for server %s type %s name %s", serverID, streamType, serverName)
	}

	return dto, nil
}

func (s *Service) GetRandomAnime(ctx context.Context) (models.AnimeDto, error) {
	data, err := s.repo.GetRandomAnime(ctx)
	if err != nil {
		log.Printf("failed to fetch random anime: %v", err)
		return models.AnimeDto{}, err
	}

	go s.refresher.MaybeRefresh(context.Background(), data.MalID.Int32)

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

	go s.refresher.MaybeRefresh(context.Background(), data.MalID.Int32)

	return models.AnimeDto{}.FromRepository(data), nil
}

func (s *Service) GetSeasonalAnimes(ctx context.Context) ([]models.SeasonalAnimeDto, error) {
	key := "seasonal_animes"
	var cachedAnimes []models.SeasonalAnimeDto
	if ok, err := s.redis.Get(ctx, key, &cachedAnimes); err != nil {
		log.Printf("failed to get seasonal animes from cache: %v", err)
	} else if ok {
		log.Printf("found %d seasonal animes in cache", len(cachedAnimes))
		return cachedAnimes, nil
	}

	year := time.Now().Year()
	month := time.Now().Month()

	season := "WINTER"
	switch month {
	case time.February, time.March, time.April:
		season = "WINTER"
	case time.May, time.June, time.July:
		season = "SPRING"
	case time.August, time.September, time.October:
		season = "SUMMER"
	case time.November, time.December, time.January:
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
		go s.refresher.MaybeRefresh(context.Background(), d.MalID.Int32)
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

	if err := s.redis.Set(ctx, key, seasonalAnimes, 30*24*time.Hour); err != nil {
		log.Printf("failed to cache seasonal animes: %v", err)
	} else {
		log.Printf("cached %d seasonal animes", len(seasonalAnimes))
	}

	return seasonalAnimes, nil
}

func (s *Service) GetAnimeBanner(ctx context.Context, id string) (string, error) {
	key := fmt.Sprintf("anime_banner:%s", id)

	var cachedBanner string
	if ok, err := s.redis.Get(ctx, key, &cachedBanner); err != nil {
		log.Printf("failed to get banner from cache: %v", err)
	} else if ok {
		log.Printf("found banner in cache for anime ID %s: %s", id, cachedBanner)
		return cachedBanner, nil
	}

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

	if err := s.redis.Set(ctx, key, bannerURL, 30*24*time.Hour); err != nil {
		log.Printf("failed to cache banner for anime ID %s: %v", id, err)
	} else {
		log.Printf("cached banner for anime ID %s: %s", id, bannerURL)
	}

	return bannerURL, nil
}
