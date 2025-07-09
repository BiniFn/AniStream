package worker

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/coeeter/aniways/internal/hianime"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	hourlyInterval = time.Hour
)

func HourlyTask(
	ctx context.Context,
	scraper *hianime.HianimeScraper,
	repo *repository.Queries,
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
			if err := scrapeRecentlyUpdated(ctx, scraper, repo); err != nil {
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
			if err != nil && err != sql.ErrNoRows {
				log.Printf("❌ RU DB lookup error for %s: %v", scraped.HiAnimeID, err)
				return
			}

			needsFetch := err == sql.ErrNoRows || dbAnime.LastEpisode != int32(scraped.LastEpisode)
			if !needsFetch {
				return
			}

			info, err := retryFetchDetail(ctx, scraper, scraped.HiAnimeID)
			if err != nil {
				log.Printf("⚠️ Error fetching detail for HiAnimeID %s: %v", scraped.HiAnimeID, err)
				return
			}

			updatedAt := now.Add(time.Duration(offset) * updateSpacing)

			if err == sql.ErrNoRows {
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
			} else {
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
			}
		}(scraped, offset)
	}

	wg.Wait()
	log.Printf("✅ Finished scraping %d recently updated anime", len(items))
	return nil
}
