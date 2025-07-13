package anime

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/coeeter/aniways/internal/client/myanimelist"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/semaphore"
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
	sem       *semaphore.Weighted
	queue     chan int32
	inFlight  map[int32]struct{}
	mu        sync.Mutex
}

// New creates a MetadataRefresher with sensible defaults:
//   - 20 concurrent workers
//   - queue size 1000
//   - 30-day TTL
//   - ~60 MAL API calls per minute
func NewRefresher(repo *repository.Queries, malClient *myanimelist.Client) *MetadataRefresher {
	m := &MetadataRefresher{
		repo:      repo,
		malClient: malClient,
		ttl:       defaultTTL,
		limiter:   rate.NewLimiter(defaultMALRate, 1),
		sem:       semaphore.NewWeighted(int64(defaultWorkerCount)),
		queue:     make(chan int32, defaultQueueSize),
		inFlight:  make(map[int32]struct{}),
	}
	for i := 0; i < defaultWorkerCount; i++ {
		go m.worker()
	}
	return m
}

func (m *MetadataRefresher) MaybeRefresh(ctx context.Context, malID int32) {
	row, err := m.repo.GetAnimeMetadataByMalId(ctx, malID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) && !errors.Is(err, pgx.ErrNoRows) {
			log.Printf("metadata lookup error for MAL %d: %v", malID, err)
			return
		}
	} else {
		if time.Since(row.UpdatedAt.Time) < m.ttl {
			return
		}
	}

	m.mu.Lock()
	if _, busy := m.inFlight[malID]; busy {
		m.mu.Unlock()
		return
	}
	m.inFlight[malID] = struct{}{}
	m.mu.Unlock()

	m.queue <- malID
}

func (m *MetadataRefresher) worker() {
	for malID := range m.queue {
		row, err := m.repo.GetAnimeMetadataByMalId(context.Background(), malID)
		if err == nil && time.Since(row.UpdatedAt.Time) < m.ttl {
			m.mu.Lock()
			delete(m.inFlight, malID)
			m.mu.Unlock()
			continue
		}

		if err := m.limiter.Wait(context.Background()); err != nil {
			log.Printf("rate limiter error: %v", err)
		}

		if err := m.sem.Acquire(context.Background(), 1); err != nil {
			continue
		}

		go func(id int32) {
			defer m.sem.Release(1)
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			dto, err := m.malClient.GetAnimeMetadata(ctx, int(id))
			if err != nil {
				log.Printf("MAL fetch failed for %d: %v", id, err)
			} else {
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
					log.Printf("metadata upsert failed for %d: %v", id, err)
				}
			}

			// clear in-flight flag
			m.mu.Lock()
			delete(m.inFlight, id)
			m.mu.Unlock()
		}(malID)
	}
}

// RefreshBlocking performs one immediate fetch/upsert under rate- and concurrency- limits.
// Use this in GetAnimeByID when you need fresh data before returning.
func (m *MetadataRefresher) RefreshBlocking(ctx context.Context, malID int32) error {
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

	err = m.repo.InsertAnimeMetadata(ctx, params)
	return err
}
