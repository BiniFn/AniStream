package anime

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coeeter/aniways/internal/infra/cache"
	"github.com/jackc/pgx/v5"
)

func (s *AnimeService) GetAnimeEpisodes(ctx context.Context, id string) ([]EpisodeDto, error) {
	return cache.GetOrFill(ctx, s.redis, fmt.Sprintf("anime_episodes:%s", id), 7*24*time.Hour, func(ctx context.Context) ([]EpisodeDto, error) {
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

		episodeDtos := make([]EpisodeDto, len(episodes))
		for i, ep := range episodes {
			episodeDtos[i] = EpisodeDto{}.FromScraper(ep)
		}
		return episodeDtos, nil
	})
}

func (s *AnimeService) GetEpisodeServers(ctx context.Context, id, episodeID string) ([]EpisodeServerDto, error) {
	key := fmt.Sprintf("episode_servers:%s:%s", id, episodeID)
	return cache.GetOrFill(ctx, s.redis, key, 24*time.Hour, func(ctx context.Context) ([]EpisodeServerDto, error) {
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

		serverDtos := make([]EpisodeServerDto, len(servers))
		for i, server := range servers {
			serverDtos[i] = EpisodeServerDto{}.FromScraper(server)
		}

		return serverDtos, nil
	})
}

func (s *AnimeService) GetEpisodeStream(ctx context.Context, id, serverID, serverName, streamType string) (StreamingDataDto, error) {
	key := fmt.Sprintf("episode_stream:%s:%s:%s:%s", id, serverID, serverName, streamType)
	return cache.GetOrFill(ctx, s.redis, key, 24*time.Hour, func(ctx context.Context) (StreamingDataDto, error) {
		_, err := s.repo.GetAnimeById(ctx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			return StreamingDataDto{}, ErrAnimeNotFound
		}
		if err != nil {
			return StreamingDataDto{}, fmt.Errorf("failed to fetch anime by ID %s: %v", id, err)
		}

		streamData, err := s.scraper.GetStreamData(ctx, serverID, streamType, serverName)
		if err != nil {
			return StreamingDataDto{}, fmt.Errorf("failed to fetch episode stream for anime ID %s server %s: %v", id, serverID, err)
		}

		return StreamingDataDto{}.FromScraper(streamData), nil
	})
}
