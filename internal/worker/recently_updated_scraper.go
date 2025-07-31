package worker

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/client/hianime"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func hourlyTask(
	ctx context.Context,
	scraper *hianime.HianimeScraper,
	repo *repository.Queries,
	redis *cache.RedisClient,
	log *slog.Logger,
) {
	log.Info("Running hourly task")
	if err := scrapeRecentlyUpdated(ctx, scraper, repo, redis, log); err != nil {
		log.Error("Error in hourly task", "err", err)
	} else {
		log.Info("Hourly task completed successfully")
	}
}

func scrapeRecentlyUpdated(
	ctx context.Context,
	scraper *hianime.HianimeScraper,
	repo *repository.Queries,
	redis *cache.RedisClient,
	log *slog.Logger,
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

			child := log.With("hi_id", scraped.HiAnimeID)

			dbAnime, err := repo.GetAnimeByHiAnimeId(ctx, scraped.HiAnimeID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) && !errors.Is(err, pgx.ErrNoRows) {
				child.Error("db lookup failed", "err", err)
				return
			}

			hasExisting := dbAnime.ID != "" && dbAnime.HiAnimeID == scraped.HiAnimeID
			needsFetch := errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) || dbAnime.LastEpisode != int32(scraped.LastEpisode)
			if !needsFetch {
				return
			}

			info, err := retryFetchDetail(ctx, scraper, scraped.HiAnimeID)
			if err != nil {
				child.Warn("detail fetch failed", "err", err)
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
					child.Error("update failed", "err", err)
				}
				if err := redis.Del(ctx, "anime_episodes:"+dbAnime.ID); err != nil {
					child.Warn("cache delete failed", "err", err)
				}
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
					child.Error("insert failed", "err", err)
				}
			}
		}(scraped, offset)
	}

	wg.Wait()
	log.Info("recently‑updated page processed", "items", len(items))
	return nil
}
