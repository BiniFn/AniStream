package anime

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/coeeter/aniways/internal/client/myanimelist"
	"github.com/coeeter/aniways/internal/repository"
	"golang.org/x/time/rate"
)

const (
	defaultWorkerCount = 20
	defaultQueueSize   = 1000
	defaultTTL         = 30 * 24 * time.Hour // 30 days
)

var defaultMALRate = rate.Every(time.Minute / 60) // ~60 req/min

type MetadataRefresher struct {
	repo      *repository.Queries
	malClient *myanimelist.Client
	ttl       time.Duration
	limiter   *rate.Limiter
	queue     chan int32
	inFlight  map[int32]struct{}
	mu        sync.Mutex
}

func NewRefresher(repo *repository.Queries, malClient *myanimelist.Client) *MetadataRefresher {
	m := &MetadataRefresher{
		repo:      repo,
		malClient: malClient,
		ttl:       defaultTTL,
		limiter:   rate.NewLimiter(defaultMALRate, 1),
		queue:     make(chan int32, defaultQueueSize),
		inFlight:  make(map[int32]struct{}),
	}
	for range defaultWorkerCount {
		go m.worker()
	}
	return m
}

func (m *MetadataRefresher) Enqueue(malID int32) {
	if malID <= 0 {
		log.Printf("invalid MAL ID %d, skipping metadata refresh", malID)
		return
	}

	m.mu.Lock()
	if _, busy := m.inFlight[malID]; busy {
		m.mu.Unlock()
		return
	}
	m.inFlight[malID] = struct{}{}
	m.mu.Unlock()

	select {
	case m.queue <- malID:
	default:
		m.clearInFlight(malID)
		log.Printf("queue full, dropping metadata refresh for MAL ID %d", malID)
	}
}

func (m *MetadataRefresher) worker() {
	for malID := range m.queue {
		row, err := m.repo.GetAnimeMetadataByMalId(context.Background(), malID)
		if err == nil && time.Since(row.UpdatedAt.Time) < m.ttl {
			m.clearInFlight(malID)
			continue
		}

		if err := m.limiter.Wait(context.Background()); err != nil {
			log.Printf("rate limiter error: %v", err)
			m.clearInFlight(malID)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		dto, err := m.malClient.GetAnimeMetadata(ctx, int(malID))
		cancel()
		if err != nil {
			log.Printf("MAL fetch failed for %d: %v", malID, err)
			m.clearInFlight(malID)
			continue
		}

		meta := dto.ToRepository()
		params := repository.InsertAnimeMetadataParams{
			MalID:              meta.MalID,
			Description:        meta.Description,
			MainPictureUrl:     meta.MainPictureUrl,
			MediaType:          meta.MediaType,
			Rating:             meta.Rating,
			AiringStatus:       meta.AiringStatus,
			AvgEpisodeDuration: meta.AvgEpisodeDuration,
			TotalEpisodes:      meta.TotalEpisodes,
			Studio:             meta.Studio,
			Rank:               meta.Rank,
			Mean:               meta.Mean,
			Scoringusers:       meta.Scoringusers,
			Popularity:         meta.Popularity,
			AiringStartDate:    meta.AiringStartDate,
			AiringEndDate:      meta.AiringEndDate,
			Source:             meta.Source,
			SeasonYear:         meta.SeasonYear,
			Season:             meta.Season,
			TrailerEmbedUrl:    meta.TrailerEmbedUrl,
		}
		if err := m.repo.InsertAnimeMetadata(ctx, params); err != nil {
			log.Printf("metadata upsert failed for %d: %v", malID, err)
		}

		m.clearInFlight(malID)
	}
}

func (m *MetadataRefresher) clearInFlight(malID int32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.inFlight, malID)
}

func (m *MetadataRefresher) RefreshBlocking(ctx context.Context, malID int32) error {
	if malID <= 0 {
		return fmt.Errorf("invalid MAL ID %d", malID)
	}
	if err := m.limiter.Wait(ctx); err != nil {
		return err
	}

	dto, err := m.malClient.GetAnimeMetadata(ctx, int(malID))
	if err != nil {
		return err
	}
	meta := dto.ToRepository()
	params := repository.InsertAnimeMetadataParams{
		MalID:              meta.MalID,
		Description:        meta.Description,
		MainPictureUrl:     meta.MainPictureUrl,
		MediaType:          meta.MediaType,
		Rating:             meta.Rating,
		AiringStatus:       meta.AiringStatus,
		AvgEpisodeDuration: meta.AvgEpisodeDuration,
		TotalEpisodes:      meta.TotalEpisodes,
		Studio:             meta.Studio,
		Rank:               meta.Rank,
		Mean:               meta.Mean,
		Scoringusers:       meta.Scoringusers,
		Popularity:         meta.Popularity,
		AiringStartDate:    meta.AiringStartDate,
		AiringEndDate:      meta.AiringEndDate,
		Source:             meta.Source,
		SeasonYear:         meta.SeasonYear,
		Season:             meta.Season,
		TrailerEmbedUrl:    meta.TrailerEmbedUrl,
	}

	return m.repo.InsertAnimeMetadata(ctx, params)
}

func (m *MetadataRefresher) Close() {
	close(m.queue) // stop workers
	m.mu.Lock()
	defer m.mu.Unlock()
	for malID := range m.inFlight {
		delete(m.inFlight, malID) // clear any remaining in-flight entries
	}
}
