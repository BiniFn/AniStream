package worker

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/coeeter/aniways/internal/infra/client/anilist"
	"github.com/coeeter/aniways/internal/infra/client/myanimelist"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type libraryImportJobPayload struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Provider string `json:"provider"`
	Status   string `json:"status"`
}

func startLibraryImportJobListener(
	ctx context.Context,
	db *pgxpool.Pool,
	repo *repository.Queries,
	malClient *myanimelist.Client,
	aniClient *anilist.Client,
	log *slog.Logger,
) error {
	conn, err := db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "LISTEN library_import_jobs")
	if err != nil {
		return err
	}

	log.Info("library import listener started")

	for {
		notification, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Error("WaitForNotification error", "err", err)
			continue
		}

		var payload libraryImportJobPayload
		if err := json.Unmarshal([]byte(notification.Payload), &payload); err != nil {
			log.Error("invalid library import job payload", "err", err)
			continue
		}

		go handleLibraryImportJob(ctx, repo, malClient, aniClient, log, payload)
	}
}

func handleLibraryImportJob(
	ctx context.Context,
	repo *repository.Queries,
	malClient *myanimelist.Client,
	aniClient *anilist.Client,
	log *slog.Logger,
	payload libraryImportJobPayload,
) {
	log.Info("processing library import job",
		"id", payload.ID,
		"user_id", payload.UserID,
		"status", payload.Status,
		"provider", payload.Provider,
	)

	token, err := repo.GetToken(ctx, repository.GetTokenParams{
		UserID:   payload.UserID,
		Provider: repository.Provider(payload.Provider),
	})
	if err != nil {
		log.Error("failed to get token", "err", err)
		return
	}

	err = repo.UpdateLibraryImportJob(ctx, repository.UpdateLibraryImportJobParams{
		ID:     payload.ID,
		Status: repository.LibraryImportStatusInProgress,
	})
	if err != nil {
		log.Error("failed to update library import job", "err", err)
		return
	}

	finalStatus := repository.LibraryImportStatusCompleted
	switch payload.Provider {
	case string(repository.ProviderAnilist):
		err = importFromAnilist(ctx, repo, aniClient, token.Token, payload, log)
	case string(repository.ProviderMyanimelist):
		err = importFromMal(ctx, repo, malClient, token.Token, payload, log)
	default:
		log.Warn("unsupported provider", "provider", payload.Provider)
		finalStatus = repository.LibraryImportStatusFailed
	}

	errMsg := pgtype.Text{}
	if err != nil {
		finalStatus = repository.LibraryImportStatusFailed
		errMsg = pgtype.Text{
			String: err.Error(),
			Valid:  true,
		}
	}

	err = repo.UpdateLibraryImportJob(ctx, repository.UpdateLibraryImportJobParams{
		ID:           payload.ID,
		Status:       finalStatus,
		ErrorMessage: errMsg,
	})
	if err != nil {
		log.Error("failed to update library import job", "err", err)
	}
}

func importFromMal(
	ctx context.Context,
	repo *repository.Queries,
	malClient *myanimelist.Client,
	token string,
	payload libraryImportJobPayload,
	log *slog.Logger,
) error {
	page := 1
	itemsPerPage := 100

	for {
		list, err := malClient.GetAnimeList(ctx, myanimelist.GetAnimeListParams{
			Token:        token,
			Page:         page,
			ItemsPerPage: itemsPerPage,
		})
		if err != nil {
			return err
		}

		if len(list.Data) == 0 {
			break
		}

		page++

		for _, item := range list.Data {
			malID := pgtype.Int4{
				Int32: int32(item.Node.MalID),
				Valid: item.Node.MalID != 0,
			}
			anime, err := repo.GetAnimeByMalId(ctx, malID)
			if err != nil {
				log.Error("failed to get anime by mal id", "mal_id", malID, "err", err)
				continue
			}

			inLibraryAlready, err := repo.IsAnimeInLibrary(ctx, repository.IsAnimeInLibraryParams{
				UserID:  payload.UserID,
				AnimeID: anime.ID,
			})
			if err != nil {
				log.Error("failed to check if anime is in library", "anime_id", anime.ID, "err", err)
				continue
			}

			status := myanimelist.MalListStatus(item.ListStatus.Status)
			watchedEpisodes := int32(item.ListStatus.EpisodesWatched)
			updatedAt, err := time.Parse(time.RFC3339, item.ListStatus.UpdatedAt)
			if err != nil {
				updatedAt = time.Now()
			}

			if !inLibraryAlready {
				err = repo.InsertLibrary(ctx, repository.InsertLibraryParams{
					UserID:          payload.UserID,
					AnimeID:         anime.ID,
					Status:          repository.LibraryStatus(status.ToRepository()),
					WatchedEpisodes: watchedEpisodes,
					UpdatedAt:       updatedAt,
				})
				if err != nil {
					log.Error("failed to insert library entry", "anime_id", anime.ID, "err", err)
				}
				continue
			}

			err = repo.UpdateLibrary(ctx, repository.UpdateLibraryParams{
				UserID:          payload.UserID,
				AnimeID:         anime.ID,
				Status:          repository.LibraryStatus(status.ToRepository()),
				WatchedEpisodes: watchedEpisodes,
				UpdatedAt: pgtype.Timestamp{
					Time:  updatedAt,
					Valid: true,
				},
			})
			if err != nil {
				log.Error("failed to update library entry", "anime_id", anime.ID, "err", err)
				continue
			}
		}
	}

	return nil
}

func importFromAnilist(
	ctx context.Context,
	repo *repository.Queries,
	aniClient *anilist.Client,
	token string,
	payload libraryImportJobPayload,
	log *slog.Logger,
) error {
	page := 1
	itemsPerPage := 100

	for {
		list, err := aniClient.GetUserAnimeList(ctx, anilist.GetUserAnimeListParams{
			Token:        token,
			Page:         page,
			ItemsPerPage: itemsPerPage,
		})
		if err != nil {
			return err
		}

		if len(list.Page.MediaList) == 0 {
			break
		}

		page++

		for _, item := range list.Page.MediaList {
			malID := pgtype.Int4{
				Int32: int32(item.Media.GetIdMal()),
				Valid: item.Media.GetIdMal() != 0,
			}
			anime, err := repo.GetAnimeByMalId(ctx, malID)
			if err != nil {
				log.Error("failed to get anime by mal id", "mal_id", malID, "err", err)
				continue
			}

			inLibraryAlready, err := repo.IsAnimeInLibrary(ctx, repository.IsAnimeInLibraryParams{
				UserID:  payload.UserID,
				AnimeID: anime.ID,
			})
			if err != nil {
				log.Error("failed to check if anime is in library", "anime_id", anime.ID, "err", err)
				continue
			}

			status := aniClient.ConvertToRepoStatus(item.GetStatus())
			watchedEpisodes := item.GetProgress()

			if !inLibraryAlready {
				err = repo.InsertLibrary(ctx, repository.InsertLibraryParams{
					UserID:          payload.UserID,
					AnimeID:         anime.ID,
					Status:          repository.LibraryStatus(status),
					WatchedEpisodes: int32(watchedEpisodes),
				})
				if err != nil {
					log.Error("failed to insert library entry", "anime_id", anime.ID, "err", err)
				}
				continue
			}

			err = repo.UpdateLibrary(ctx, repository.UpdateLibraryParams{
				UserID:          payload.UserID,
				AnimeID:         anime.ID,
				Status:          repository.LibraryStatus(status),
				WatchedEpisodes: int32(watchedEpisodes),
			})
			if err != nil {
				log.Error("failed to update library entry", "anime_id", anime.ID, "err", err)
				continue
			}
		}
	}

	return nil
}
