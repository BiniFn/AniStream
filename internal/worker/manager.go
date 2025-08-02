package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/coeeter/aniways/internal/infra/cache"
	"github.com/coeeter/aniways/internal/infra/client/anilist"
	"github.com/coeeter/aniways/internal/infra/client/hianime"
	"github.com/coeeter/aniways/internal/infra/client/myanimelist"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/service/auth/oauth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type Manager struct {
	db        *pgxpool.Pool
	repo      *repository.Queries
	scraper   *hianime.HianimeScraper
	malClient *myanimelist.Client
	aniClient *anilist.Client
	redis     *cache.RedisClient
	log       *slog.Logger
}

func NewManager(
	db *pgxpool.Pool,
	repo *repository.Queries,
	scraper *hianime.HianimeScraper,
	malClient *myanimelist.Client,
	aniClient *anilist.Client,
	redis *cache.RedisClient,
	log *slog.Logger,
) *Manager {
	return &Manager{
		db:        db,
		repo:      repo,
		scraper:   scraper,
		malClient: malClient,
		aniClient: aniClient,
		redis:     redis,
		log:       log,
	}
}

func (m *Manager) Bootstrap(ctx context.Context) error {
	count, err := m.repo.GetCountOfAnimes(ctx)
	if err != nil {
		return fmt.Errorf("count animes: %w", err)
	}

	if count == 0 {
		log := m.log.With("job", "full-seed")
		log.Info("no anime in DB — running initial scrape (blocking)")

		if err := fullSeed(ctx, m.scraper, m.repo, log); err != nil {
			return fmt.Errorf("full seed: %w", err)
		}
		log.Info("initial scrape complete")
	} else {
		m.log.Info("database already seeded; skipping initial scrape", "count", count)
	}
	return nil
}

func (m *Manager) StartBackground(ctx context.Context, providers []oauth.Provider) {
	c := cron.New()
	_, err := c.AddFunc("@hourly", func() {
		hourlyTask(ctx, m.scraper, m.repo, m.redis, m.log.With("job", "hourly-scrape"))
	})
	if err != nil {
		m.log.Error("failed to add hourly task", "err", err)
		return
	}

	_, err = c.AddFunc("@daily", func() {
		dailyTask(ctx, m.repo, providers, m.log.With("job", "daily-refresh-token"))
	})
	if err != nil {
		m.log.Error("failed to add daily task", "err", err)
		return
	}

	_, err = c.AddFunc("@every 6h", func() {
		retryFailedLibrarySyncs(ctx, m.repo, m.malClient, m.aniClient, m.log.With("job", "failed-library-sync-cron"))
	})

	m.log.Info("bootstrapping hourly + daily cron job")
	c.Start()

	go func() {
		err := startLibrarySyncListener(
			ctx,
			m.db,
			m.repo,
			m.malClient,
			m.aniClient,
			m.log.With("job", "library-sync"),
		)
		if err != nil {
			m.log.Error("library sync listener stopped", "err", err)
		}
	}()

	go func() {
		err := startLibraryImportJobListener(
			ctx,
			m.db,
			m.repo,
			m.malClient,
			m.aniClient,
			m.log.With("job", "library-import"),
		)
		if err != nil {
			m.log.Error("library import listener stopped", "err", err)
		}
	}()

	go func() {
		<-ctx.Done()
		m.log.Info("Shutting down cron scheduler")
		c.Stop()
	}()
}
