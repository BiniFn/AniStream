package anime

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/coeeter/aniways/internal/infra/cache"
	"github.com/jackc/pgx/v5"
)

func (s *AnimeService) GetAnimeEpisodes(ctx context.Context, id string) ([]EpisodeDto, error) {
	return cache.GetOrFill(ctx, s.redis, fmt.Sprintf("anime_episodes:%s", id), 7*24*time.Hour, func(ctx context.Context) ([]EpisodeDto, error) {
		a, err := s.repo.GetAnimeById(ctx, id)
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
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
	return cache.GetOrFill(ctx, s.redis, fmt.Sprintf("episode_servers:%s:%s", id, episodeID), 24*time.Hour, func(ctx context.Context) ([]EpisodeServerDto, error) {
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

func (s *AnimeService) GetEpisodeStream(ctx context.Context, id, episodeID, streamType string) (EpisodeSourceDto, error) {
	key := fmt.Sprintf("source:%s:%s:%s", id, episodeID, streamType)
	return cache.GetOrFill(ctx, s.redis, key, 24*time.Hour, func(ctx context.Context) (EpisodeSourceDto, error) {
		_, err := s.repo.GetAnimeById(ctx, id)
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return EpisodeSourceDto{}, ErrAnimeNotFound
		}
		if err != nil {
			return EpisodeSourceDto{}, fmt.Errorf("failed to fetch anime by ID %s: %v", id, err)
		}

		data, err := s.scraper.GetEpisodeStream(ctx, episodeID, streamType)
		if err != nil {
			return EpisodeSourceDto{}, fmt.Errorf("failed to fetch episode stream for anime ID %s episode %s: %v", id, episodeID, err)
		}

		encoder := base64.StdEncoding
		p := encoder.EncodeToString([]byte(data))
		s := "megaplay"

		return EpisodeSourceDto{
			URL:    fmt.Sprintf("/proxy?p=%s&s=%s", p, s),
			RawURL: data,
		}, nil
	})
}

func (s *AnimeService) GetStreamMetadata(ctx context.Context, id, episodeID, streamType string) (StreamingMetadataDto, error) {
	return cache.GetOrFill(ctx, s.redis, fmt.Sprintf("stream-metadata:%s:%s:%s", id, episodeID, streamType), 24*time.Hour, func(ctx context.Context) (StreamingMetadataDto, error) {
		a, err := s.repo.GetAnimeById(ctx, id)
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return StreamingMetadataDto{}, ErrAnimeNotFound
		}
		if err != nil {
			return StreamingMetadataDto{}, fmt.Errorf("failed to fetch anime by ID %s: %v", id, err)
		}

		data, err := s.scraper.GetStreamMetadata(ctx, a.HiAnimeID, episodeID, streamType)
		if err != nil {
			return StreamingMetadataDto{}, fmt.Errorf("failed to fetch stream metadata for anime ID %s episode %s type %s: %v", id, episodeID, streamType, err)
		}

		return StreamingMetadataDto{}.FromScraper(data), nil
	})
}
