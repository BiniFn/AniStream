package anime

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/coeeter/aniways/internal/models"
)

func (s *AnimeService) GetAnimeEpisodes(ctx context.Context, id string) ([]models.EpisodeDto, error) {
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

func (s *AnimeService) GetServersForEpisode(ctx context.Context, id, episodeID string) (models.ServerDto, error) {
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

func (s *AnimeService) GetEpisodeStream(ctx context.Context, id, episodeID, streamType string) (models.EpisodeSourceDto, error) {
	var cached models.EpisodeSourceDto

	key := fmt.Sprintf("source:%s:%s:%s", id, episodeID, streamType)
	if _, err := s.redis.GetOrFill(ctx, key, &cached, 24*time.Hour, func(ctx context.Context) (any, error) {
		data, err := s.scraper.GetStreamingData(ctx, episodeID, streamType, "megaplay")
		if err != nil {
			log.Printf("failed to fetch episode stream for anime ID %s episode %s type %s: %v", id, episodeID, streamType, err)
			return models.EpisodeSourceDto{}, err
		}

		encoder := base64.StdEncoding
		p := encoder.EncodeToString([]byte(data.Source))
		s := "megaplay"

		return models.EpisodeSourceDto{
			URL:    fmt.Sprintf("/proxy?p=%s&s=%s", p, s),
			RawURL: data.Source,
		}, nil
	}); err != nil {
		log.Printf("failed to get episode stream for anime ID %s episode %s type %s from cache: %v", id, episodeID, streamType, err)
		return cached, err
	}

	return cached, nil
}

func (s *AnimeService) GetStreamingData(ctx context.Context, serverID, streamType, serverName string) (models.StreamingDataDto, error) {
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
