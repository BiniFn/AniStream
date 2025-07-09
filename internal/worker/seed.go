package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/coeeter/aniways/internal/hianime"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	maxConcurrency = 20
	retryCount     = 3
	retryDelay     = 500 * time.Millisecond
	pageDelay      = 1 * time.Second
	updateSpacing  = 20 * time.Millisecond
)

func FullSeed(
	ctx context.Context,
	scraper *hianime.HianimeScraper,
	repo *repository.Queries,
) error {
	log.Println("⚙️  Starting AZ‐list seed")
	if err := scrapeAllAZ(ctx, scraper, repo); err != nil {
		return err
	}

	log.Println("⚙️  Starting reverse Recently‐Updated seed")
	if err := scrapeAllRecentlyUpdated(ctx, scraper, repo); err != nil {
		return err
	}

	log.Println("✅ Full seed complete")
	return nil
}

func scrapeAllAZ(
	ctx context.Context,
	scraper *hianime.HianimeScraper,
	repo *repository.Queries,
) error {
	sem := make(chan struct{}, maxConcurrency)
	for page := 1; ; page++ {
		log.Printf("🔎 AZ page %d", page)
		listing, err := scraper.GetAZList(ctx, page)
		if err != nil {
			return fmt.Errorf("AZ page %d: %w", page, err)
		}

		var (
			wg       sync.WaitGroup
			mu       sync.Mutex
			toInsert = make([]repository.InsertMultipleAnimesParams, 0, len(listing.Items))
		)
		for _, a := range listing.Items {
			wg.Add(1)
			go func(item hianime.ScrapedAnimeInfoDto) {
				defer wg.Done()

				select {
				case sem <- struct{}{}:
				case <-ctx.Done():
					return
				}
				defer func() { <-sem }()

				info, err := retryFetchDetail(ctx, scraper, item.HiAnimeID)
				if err != nil {
					log.Printf("⚠️ AZ detail %s: %v", item.HiAnimeID, err)
					return
				}

				p := repository.InsertMultipleAnimesParams{
					Ename:       item.EName,
					Jname:       item.JName,
					ImageUrl:    item.PosterURL,
					Genre:       info.Genre,
					HiAnimeID:   item.HiAnimeID,
					MalID:       pgtype.Int4{Int32: int32(info.MalID), Valid: info.MalID > 0},
					AnilistID:   pgtype.Int4{Int32: int32(info.AnilistID), Valid: info.AnilistID > 0},
					LastEpisode: int32(item.LastEpisode),
				}
				mu.Lock()
				toInsert = append(toInsert, p)
				mu.Unlock()
			}(a)
		}
		wg.Wait()

		if len(toInsert) > 0 {
			if _, err := repo.InsertMultipleAnimes(ctx, toInsert); err != nil {
				return fmt.Errorf("insert AZ page %d: %w", page, err)
			}
			log.Printf("✅ inserted %d from AZ page %d", len(toInsert), page)
		}

		if !listing.PageInfo.HasNextPage {
			log.Printf("ℹ️  AZ page %d has no next page, stopping", page)
			return nil
		}
		time.Sleep(pageDelay)
	}
}

func retryFetchDetail(
	ctx context.Context,
	scraper *hianime.HianimeScraper,
	hiID string,
) (hianime.ScrapedAnimeInfoDto, error) {
	var lastErr error
	for i := 0; i < retryCount; i++ {
		info, err := scraper.GetAnimeInfoByHiAnimeID(ctx, hiID)
		if err == nil {
			return info, nil
		}
		lastErr = err
		time.Sleep(retryDelay)
	}
	return hianime.ScrapedAnimeInfoDto{}, lastErr
}

func scrapeAllRecentlyUpdated(
	ctx context.Context,
	scraper *hianime.HianimeScraper,
	repo *repository.Queries,
) error {
	first, err := scraper.GetRecentlyUpdatedAnime(ctx, 1)
	if err != nil {
		return fmt.Errorf("RU first page: %w", err)
	}
	total := first.PageInfo.TotalPages
	log.Printf("🔎 will scrape %d RU pages in reverse order", total)

	sem := make(chan struct{}, maxConcurrency)
	start := time.Now()

	for page := total; page >= 1; page-- {
		listing, err := scraper.GetRecentlyUpdatedAnime(ctx, page)
		if err != nil {
			return fmt.Errorf("RU page %d: %w", page, err)
		}
		log.Printf("📄 RU page %d/%d (items: %d)", page, total, len(listing.Items))

		hiIDs := make([]string, len(listing.Items))
		for i, a := range listing.Items {
			hiIDs[i] = a.HiAnimeID
		}
		existingRows, err := repo.GetAnimesByHiAnimeIds(ctx, hiIDs)
		if err != nil {
			return fmt.Errorf("query existing RU page %d: %w", page, err)
		}
		existingMap := make(map[string]repository.Anime, len(existingRows))
		for _, row := range existingRows {
			existingMap[row.HiAnimeID] = row
		}

		var wg sync.WaitGroup
		for idx := len(listing.Items) - 1; idx >= 0; idx-- {
			scraped := listing.Items[idx]
			wg.Add(1)
			go func(scraped hianime.ScrapedAnimeInfoDto, globalIdx int) {
				defer wg.Done()

				select {
				case sem <- struct{}{}:
				case <-ctx.Done():
					return
				}
				defer func() { <-sem }()

				info, err := retryFetchDetail(ctx, scraper, scraped.HiAnimeID)
				if err != nil {
					log.Printf("⚠️ RU detail %s: %v", scraped.HiAnimeID, err)
					return
				}

				updatedAt := start.Add(time.Duration(globalIdx) * updateSpacing)

				if existing, ok := existingMap[scraped.HiAnimeID]; ok {
					params := repository.UpdateAnimeParams{
						ID:          existing.ID,
						Ename:       info.EName,
						Jname:       info.JName,
						ImageUrl:    info.PosterURL,
						Genre:       info.Genre,
						HiAnimeID:   info.HiAnimeID,
						MalID:       pgtype.Int4{Int32: int32(info.MalID), Valid: info.MalID > 0},
						AnilistID:   pgtype.Int4{Int32: int32(info.AnilistID), Valid: info.AnilistID > 0},
						LastEpisode: int32(scraped.LastEpisode),
						UpdatedAt:   pgtype.Timestamp{Time: updatedAt, Valid: true},
					}
					if err := repo.UpdateAnime(ctx, params); err != nil {
						log.Printf("❌ RU update %s: %v", scraped.HiAnimeID, err)
					}
				} else {
					params := repository.InsertAnimeParams{
						Ename:       info.EName,
						Jname:       info.JName,
						ImageUrl:    info.PosterURL,
						Genre:       info.Genre,
						HiAnimeID:   info.HiAnimeID,
						MalID:       pgtype.Int4{Int32: int32(info.MalID), Valid: info.MalID > 0},
						AnilistID:   pgtype.Int4{Int32: int32(info.AnilistID), Valid: info.AnilistID > 0},
						LastEpisode: int32(scraped.LastEpisode),
						CreatedAt:   pgtype.Timestamp{Time: updatedAt, Valid: true},
						UpdatedAt:   pgtype.Timestamp{Time: updatedAt, Valid: true},
					}
					if err := repo.InsertAnime(ctx, params); err != nil {
						log.Printf("❌ RU insert %s: %v", scraped.HiAnimeID, err)
					}
				}

			}(scraped, (total-page)*len(listing.Items)+(len(listing.Items)-1-idx))
		}

		wg.Wait()
		time.Sleep(pageDelay)
	}

	log.Printf("✅ Finished scraping all %d RU pages in %s", total, time.Since(start))
	return nil
}
