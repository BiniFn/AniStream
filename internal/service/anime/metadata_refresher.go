package anime

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/coeeter/aniways/internal/infra/client/myanimelist"
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

func logger() *slog.Logger {
	return slog.Default().With("component", "metadata_refresher")
}

func (m *MetadataRefresher) Enqueue(malID int32) {
	log := logger()

	if malID <= 0 {
		log.Error("invalid MAL ID", "mal_id", malID)
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
		log.Warn("queue full, dropping metadata refresh for MAL ID", "mal_id", malID)
	}
}

func (m *MetadataRefresher) worker() {
	log := logger()

	for malID := range m.queue {
		row, err := m.repo.GetAnimeMetadataByMalId(context.Background(), malID)
		if err == nil && time.Since(row.UpdatedAt.Time) < m.ttl {
			m.clearInFlight(malID)
			continue
		}

		if err := m.limiter.Wait(context.Background()); err != nil {
			log.Error("rate limiter error", "error", err)
			m.clearInFlight(malID)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		dto, err := m.malClient.GetAnimeMetadata(ctx, int(malID))
		if err != nil {
			log.Error("MAL fetch failed", "mal_id", malID, "error", err)
			m.clearInFlight(malID)
			cancel()
			continue
		}

		params := dto.ToUpsertParams()
		if err := m.repo.UpsertAnimeMetadata(ctx, params); err != nil {
			log.Error("metadata upsert failed", "mal_id", malID, "error", err)
			cancel()
		}

		m.clearInFlight(malID)
		cancel()
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
	params := dto.ToUpsertParams()

	return m.repo.UpsertAnimeMetadata(ctx, params)
}
