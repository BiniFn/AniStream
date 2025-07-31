package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/coeeter/aniways/internal/infra/cache"
	"github.com/coeeter/aniways/internal/infra/client/hianime"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/service/auth/oauth"
	"github.com/robfig/cron/v3"
)

type Manager struct {
	repo    *repository.Queries
	scraper *hianime.HianimeScraper
	redis   *cache.RedisClient
	log     *slog.Logger
}

func NewManager(
	repo *repository.Queries,
	scraper *hianime.HianimeScraper,
	redis *cache.RedisClient,
	log *slog.Logger,
) *Manager {
	return &Manager{
		repo:    repo,
		scraper: scraper,
		redis:   redis,
		log:     log,
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

	m.log.Info("bootstrapping hourly + daily cron job")
	c.Start()

	go func() {
		<-ctx.Done()
		m.log.Info("Shutting down cron scheduler")
		c.Stop()
	}()
}
