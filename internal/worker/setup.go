package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/client/hianime"
	"github.com/coeeter/aniways/internal/repository"
)

type Manager struct {
	repo    *repository.Queries
	scraper *hianime.HianimeScraper
	redis   *cache.RedisClient
}

func NewManager(repo *repository.Queries, scraper *hianime.HianimeScraper, redis *cache.RedisClient) *Manager {
	return &Manager{
		repo:    repo,
		scraper: scraper,
		redis:   redis,
	}
}

func (m *Manager) Bootstrap(ctx context.Context) error {
	count, err := m.repo.GetCountOfAnimes(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get count of animes: %w", err)
	}

	if count == 0 {
		log.Println("⚙️  no anime in DB—running initial scraper (this will block)…")
		if err := FullSeed(ctx, m.scraper, m.repo); err != nil {
			return fmt.Errorf("🚨 FullSeed failed: %w", err)
		}
		log.Println("✅ initial scrape complete, starting HTTP server")
	} else {
		log.Printf("ℹ️  DB already has %d anime, skipping initial scrape", count)
	}

	return nil
}

func (m *Manager) StartBackground(ctx context.Context) {
	go HourlyTask(ctx, m.scraper, m.repo, m.redis)
}
