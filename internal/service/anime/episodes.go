package anime

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coeeter/aniways/internal/infra/cache"
	"github.com/coeeter/aniways/internal/mappers"
	"github.com/coeeter/aniways/internal/models"
	"github.com/jackc/pgx/v5"
)

func (s *AnimeService) GetAnimeEpisodes(ctx context.Context, id string) (models.EpisodeListResponse, error) {
	return cache.GetOrFill(ctx, s.redis, fmt.Sprintf("anime_episodes:%s", id), 7*24*time.Hour, func(ctx context.Context) (models.EpisodeListResponse, error) {
		a, err := s.repo.GetAnimeById(ctx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAnimeNotFound
		}
		if err != nil {
			return nil, fmt.Errorf("failed to fetch anime by ID %s: %v", id, err)
		}

		episodes, err := s.scraper.GetAnimeEpisodes(ctx, a.HiAnimeID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch episodes for anime ID %s: %v", id, err)
		}

		if len(episodes) == 0 {
			return nil, fmt.Errorf("no episodes found for anime ID %s", id)
		}

		episodeResponses := make([]models.EpisodeResponse, len(episodes))
		for i, ep := range episodes {
			episodeResponses[i] = mappers.EpisodeFromScraper(ep)
		}
		return episodeResponses, nil
	})
}

func (s *AnimeService) GetEpisodeServers(ctx context.Context, id, episodeID string) (models.EpisodeServerListResponse, error) {
	key := fmt.Sprintf("episode_servers:%s:%s", id, episodeID)
	return cache.GetOrFill(ctx, s.redis, key, 24*time.Hour, func(ctx context.Context) (models.EpisodeServerListResponse, error) {
		a, err := s.repo.GetAnimeById(ctx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAnimeNotFound
		}
		if err != nil {
			return nil, fmt.Errorf("failed to fetch anime by ID %s: %v", id, err)
		}

		servers, err := s.scraper.GetEpisodeServers(ctx, a.HiAnimeID, episodeID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch episode servers for anime ID %s episode %s: %v", id, episodeID, err)
		}

		serverResponses := make([]models.EpisodeServerResponse, len(servers))
		for i, server := range servers {
			serverResponses[i] = mappers.EpisodeServerFromScraper(server)
		}

		return serverResponses, nil
	})
}

func (s *AnimeService) GetEpisodeStream(ctx context.Context, id, serverID, serverName, streamType string) (models.StreamingDataResponse, error) {
	key := fmt.Sprintf("episode_stream:%s:%s:%s:%s", id, serverID, serverName, streamType)
	return cache.GetOrFill(ctx, s.redis, key, 24*time.Hour, func(ctx context.Context) (models.StreamingDataResponse, error) {
		_, err := s.repo.GetAnimeById(ctx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			return models.StreamingDataResponse{}, ErrAnimeNotFound
		}
		if err != nil {
			return models.StreamingDataResponse{}, fmt.Errorf("failed to fetch anime by ID %s: %v", id, err)
		}

		streamData, err := s.scraper.GetStreamData(ctx, serverID, streamType, serverName)
		if err != nil {
			return models.StreamingDataResponse{}, fmt.Errorf("failed to fetch episode stream for anime ID %s server %s: %v", id, serverID, err)
		}

		return mappers.StreamingDataFromScraper(streamData), nil
	})
}
