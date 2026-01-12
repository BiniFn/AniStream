package scraper

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync/atomic"
	"time"

	"github.com/coeeter/aniways/internal/infra/cache"
	"github.com/coeeter/aniways/internal/infra/client/hianime"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/sync/errgroup"
)

func HourlyTask(
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

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(maxConcurrency)

	var success, skipped, failed int32

	for i, scraped := range items {
		offset := (len(items) - 1) - i

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			child := log.With("hi_id", scraped.HiAnimeID)

			dbAnime, err := repo.GetAnimeByHiAnimeId(ctx, scraped.HiAnimeID)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				child.Error("db lookup failed", "err", err)
				atomic.AddInt32(&failed, 1)
				return nil
			}

			hasExisting := dbAnime.ID != "" && dbAnime.HiAnimeID == scraped.HiAnimeID
			needsFetch := errors.Is(err, pgx.ErrNoRows) || dbAnime.LastEpisode != int32(scraped.LastEpisode)
			if !needsFetch {
				atomic.AddInt32(&skipped, 1)
				return nil
			}

			info, err := retryFetchDetail(ctx, scraper, scraped.HiAnimeID)
			if err != nil {
				child.Warn("detail fetch failed", "err", err)
				atomic.AddInt32(&failed, 1)
				return nil
			}

			// If mal_id is missing but anilist_id exists, try to find related anime and copy mal_id
			if info.MalID == 0 && info.AnilistID > 0 {
				relatedAnimes, err := repo.GetAnimeByAnilistId(ctx, pgtype.Int4{Int32: int32(info.AnilistID), Valid: true})
				if err == nil {
					// Find anime with same anilist_id that has a mal_id
					for _, related := range relatedAnimes {
						if related.MalID.Valid && related.MalID.Int32 > 0 {
							info.MalID = int(related.MalID.Int32)
							child.Info("copied mal_id from related anime", "mal_id", info.MalID, "anilist_id", info.AnilistID, "from_hi_id", related.HiAnimeID)
							break
						}
					}
				}
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
					Season:      repository.Season(strings.ToLower(info.Season)),
					SeasonYear:  int32(info.SeasonYear),
				}
				if err := repo.UpdateAnime(ctx, params); err != nil {
					child.Error("update failed", "err", err)
					atomic.AddInt32(&failed, 1)
					return nil
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
					Season:      repository.Season(strings.ToLower(info.Season)),
					SeasonYear:  int32(info.SeasonYear),
				}
				if err := repo.InsertAnime(ctx, params); err != nil {
					child.Error("insert failed", "err", err, "hi_anime_id", scraped.HiAnimeID, "mal_id", info.MalID)
					atomic.AddInt32(&failed, 1)
					return nil
				}
			}

			atomic.AddInt32(&success, 1)
			return nil
		})
	}

	_ = g.Wait()

	log.Info("recently-updated page processed",
		"items", len(items),
		"success", success,
		"skipped", skipped,
		"failed", failed,
	)
	return nil
}

