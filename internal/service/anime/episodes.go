package anime

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"
)

func (s *AnimeService) GetAnimeEpisodes(ctx context.Context, id string) ([]EpisodeDto, error) {
	var cachedEpisodes []EpisodeDto

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

		episodeDtos := make([]EpisodeDto, len(episodes))
		for i, ep := range episodes {
			episodeDtos[i] = EpisodeDto{}.FromScraper(ep)
		}
		return episodeDtos, nil
	})

	if err != nil {
		log.Printf("failed to get anime episodes from cache: %v", err)
		return nil, err
	}

	return cachedEpisodes, nil
}

func (s *AnimeService) GetEpisodeLangs(ctx context.Context, id, episodeID string) ([]string, error) {
	var langs []string

	if id == "" || episodeID == "" {
		return langs, fmt.Errorf("id and episodeID are required")
	}

	var cachedLangs []string
	key := fmt.Sprintf("episode_langs:%s:%s", id, episodeID)
	_, err := s.redis.GetOrFill(ctx, key, &cachedLangs, 24*time.Hour, func(ctx context.Context) (any, error) {
		a, err := s.repo.GetAnimeById(ctx, id)
		if err != nil {
			log.Printf("failed to fetch anime by ID %s: %v", id, err)
			return nil, err
		}

		langs, err := s.scraper.GetEpisodeLangs(ctx, a.HiAnimeID, episodeID)
		if err != nil {
			log.Printf("failed to fetch episode langs for anime ID %s episode %s: %v", id, episodeID, err)
			return nil, err
		}

		return langs, nil
	})

	if err != nil {
		log.Printf("failed to get episode langs for anime ID %s episode %s from cache: %v", id, episodeID, err)
		return langs, err
	}

	return cachedLangs, nil
}

func (s *AnimeService) GetEpisodeStream(ctx context.Context, id, episodeID, streamType string) (EpisodeSourceDto, error) {
	var cached EpisodeSourceDto

	key := fmt.Sprintf("source:%s:%s:%s", id, episodeID, streamType)
	if _, err := s.redis.GetOrFill(ctx, key, &cached, 24*time.Hour, func(ctx context.Context) (any, error) {
		data, err := s.scraper.GetEpisodeStream(ctx, episodeID, streamType)
		if err != nil {
			log.Printf("failed to fetch episode stream for anime ID %s episode %s type %s: %v", id, episodeID, streamType, err)
			return EpisodeSourceDto{}, err
		}

		encoder := base64.StdEncoding
		p := encoder.EncodeToString([]byte(data))
		s := "megaplay"

		return EpisodeSourceDto{
			URL:    fmt.Sprintf("/proxy?p=%s&s=%s", p, s),
			RawURL: data,
		}, nil
	}); err != nil {
		log.Printf("failed to get episode stream for anime ID %s episode %s type %s from cache: %v", id, episodeID, streamType, err)
		return cached, err
	}

	return cached, nil
}

func (s *AnimeService) GetStreamMetadata(ctx context.Context, id, episodeID, streamType string) (StreamingMetadataDto, error) {
	var cached StreamingMetadataDto

	key := fmt.Sprintf("metadata:%s:%s:%s", id, episodeID, streamType)
	if _, err := s.redis.GetOrFill(ctx, key, &cached, 24*time.Hour, func(ctx context.Context) (any, error) {
		a, err := s.repo.GetAnimeById(ctx, id)
		if err != nil {
			log.Printf("failed to fetch anime metadata for ID %s: %v", id, err)
			return StreamingMetadataDto{}, err
		}

		data, err := s.scraper.GetStreamMetadata(ctx, a.HiAnimeID, episodeID, streamType)
		if err != nil {
			log.Printf("failed to fetch stream metadata for anime ID %s episode %s type %s: %v", id, episodeID, streamType, err)
			return StreamingMetadataDto{}, err
		}

		return StreamingMetadataDto{}.FromScraper(data), nil
	}); err != nil {
		log.Printf("failed to get stream metadata for anime ID %s episode %s type %s from cache: %v", id, episodeID, streamType, err)
		return cached, err
	}

	return cached, nil
}
