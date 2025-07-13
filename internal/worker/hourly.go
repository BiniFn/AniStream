package worker

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/client/hianime"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	hourlyInterval = time.Hour
)

func HourlyTask(
	ctx context.Context,
	scraper *hianime.HianimeScraper,
	repo *repository.Queries,
	redis *cache.RedisClient,
) {
	ticker := time.NewTicker(hourlyInterval)
	defer ticker.Stop()

	log.Println("⏰ Bootstrapping hourly cron job...")
	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 HourlyTask shutting down...")
			return
		case <-ticker.C:
			log.Println("🔄 Running hourly task...")
			if err := scrapeRecentlyUpdated(ctx, scraper, repo, redis); err != nil {
				log.Printf("🚨 Error in hourly task: %v", err)
			} else {
				log.Println("✅ Hourly task completed successfully")
			}
		}
	}
}

func scrapeRecentlyUpdated(
	ctx context.Context,
	scraper *hianime.HianimeScraper,
	repo *repository.Queries,
	redis *cache.RedisClient,
) error {
	listing, err := scraper.GetRecentlyUpdatedAnime(ctx, 1)
	if err != nil {
		return err
	}
	items := listing.Items
	now := time.Now()

	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for idx := len(items) - 1; idx >= 0; idx-- {
		scraped := items[idx]
		offset := (len(items) - 1) - idx

		wg.Add(1)
		go func(scraped hianime.ScrapedAnimeInfoDto, offset int) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()

			dbAnime, err := repo.GetAnimeByHiAnimeId(ctx, scraped.HiAnimeID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) && !errors.Is(err, pgx.ErrNoRows) {
				log.Printf("❌ RU DB lookup error for %s: %v", scraped.HiAnimeID, err)
				return
			}

			hasExisting := dbAnime.ID != "" && dbAnime.HiAnimeID == scraped.HiAnimeID
			needsFetch := errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) || dbAnime.LastEpisode != int32(scraped.LastEpisode)
			if !needsFetch {
				return
			}

			info, err := retryFetchDetail(ctx, scraper, scraped.HiAnimeID)
			if err != nil {
				log.Printf("⚠️ Error fetching detail for HiAnimeID %s: %v", scraped.HiAnimeID, err)
				return
			}

			updatedAt := now.Add(time.Duration(offset) * updateSpacing)

			if hasExisting {
				params := repository.UpdateAnimeParams{
					ID:          dbAnime.ID,
					Ename:       dbAnime.Ename,
					Jname:       dbAnime.Jname,
					ImageUrl:    scraped.PosterURL,
					Genre:       info.Genre,
					HiAnimeID:   scraped.HiAnimeID,
					MalID:       pgtype.Int4{Int32: int32(info.MalID), Valid: info.MalID > 0},
					AnilistID:   pgtype.Int4{Int32: int32(info.AnilistID), Valid: info.AnilistID > 0},
					LastEpisode: int32(scraped.LastEpisode),
					UpdatedAt:   pgtype.Timestamp{Time: updatedAt, Valid: true},
				}
				if err := repo.UpdateAnime(ctx, params); err != nil {
					log.Printf("❌ RU update failed for %s: %v", scraped.HiAnimeID, err)
				}
				if err := redis.Del(ctx, "anime_episodes:"+dbAnime.ID); err != nil {
					log.Printf("⚠️ Failed to delete cache for anime episodes of %s: %v", dbAnime.ID, err)
				}
				log.Printf("✅ Updated existing anime %s with new data", scraped.HiAnimeID)
			} else {
				params := repository.InsertAnimeParams{
					Ename:       info.EName,
					Jname:       info.JName,
					ImageUrl:    info.PosterURL,
					Genre:       info.Genre,
					HiAnimeID:   scraped.HiAnimeID,
					MalID:       pgtype.Int4{Int32: int32(info.MalID), Valid: info.MalID > 0},
					AnilistID:   pgtype.Int4{Int32: int32(info.AnilistID), Valid: info.AnilistID > 0},
					LastEpisode: int32(scraped.LastEpisode),
					CreatedAt:   pgtype.Timestamp{Time: updatedAt, Valid: true},
					UpdatedAt:   pgtype.Timestamp{Time: updatedAt, Valid: true},
				}
				if err := repo.InsertAnime(ctx, params); err != nil {
					log.Printf("❌ RU insert failed for %s: %v", scraped.HiAnimeID, err)
				}
			}
		}(scraped, offset)
	}

	wg.Wait()
	log.Printf("✅ Finished scraping %d recently updated anime", len(items))
	return nil
}
