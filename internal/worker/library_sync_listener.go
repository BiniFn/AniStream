package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/coeeter/aniways/internal/infra/client/anilist"
	"github.com/coeeter/aniways/internal/infra/client/myanimelist"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type librarySyncPayload struct {
	UserID   string          `json:"user_id"`
	AnimeID  string          `json:"anime_id"`
	Provider string          `json:"provider"`
	Action   string          `json:"action"`
	Payload  json.RawMessage `json:"payload"`
}

type SyncData struct {
	Status          *string `json:"status"`
	WatchedEpisodes *int32  `json:"watched_episodes"`
}

func startLibrarySyncListener(
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

	_, err = conn.Exec(ctx, "LISTEN library_sync")
	if err != nil {
		return err
	}

	log.Info("library sync listener started")

	for {
		notification, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Error("WaitForNotification error", "err", err)
			continue
		}

		var payload librarySyncPayload
		if err := json.Unmarshal([]byte(notification.Payload), &payload); err != nil {
			log.Error("Invalid library sync payload", "err", err)
			continue
		}

		go handleLibrarySync(ctx, repo, malClient, aniClient, log, payload)
	}
}

func handleLibrarySync(
	ctx context.Context,
	repo *repository.Queries,
	malClient *myanimelist.Client,
	aniClient *anilist.Client,
	log *slog.Logger,
	payload librarySyncPayload,
) {
	log.Info("Processing library sync",
		"user_id", payload.UserID,
		"anime_id", payload.AnimeID,
		"provider", payload.Provider,
		"action", payload.Action,
	)

	var syncData SyncData
	if err := json.Unmarshal(payload.Payload, &syncData); err != nil {
		log.Error("Failed to parse sync payload", "err", err)
		return
	}

	status := ""
	if syncData.Status != nil {
		status = *syncData.Status
	}
	episodes := 0
	if syncData.WatchedEpisodes != nil {
		episodes = int(*syncData.WatchedEpisodes)
	}

	tokens, err := repo.GetAllOauthTokensOfUser(ctx, payload.UserID)
	if err != nil {
		log.Error("Failed to get oauth tokens", "err", err)
		return
	}
	if len(tokens) == 0 {
		log.Error("No oauth tokens found")
		return
	}

	anime, err := repo.GetAnimeById(ctx, payload.AnimeID)
	if err != nil {
		log.Error("Failed to get anime", "err", err)
		return
	}

	for _, token := range tokens {
		ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		var err error
		switch token.Provider {
		case repository.ProviderMyanimelist:
			err = handleMalProvider(ctx, malClient, anime, token.Token, payload.Action, status, episodes)

		case repository.ProviderAnilist:
			err = handleAniProvider(ctx, aniClient, anime, token.Token, payload.Action, status, episodes)

		default:
			log.Warn("Unsupported provider", "provider", token.Provider)
			continue
		}

		finalStatus := repository.LibrarySyncStatusSuccess
		if err != nil {
			finalStatus = repository.LibrarySyncStatusFailed
			log.Error("Failed to handle provider", "provider", token.Provider, "err", err)
		}

		repo.UpdateLibrarySyncStatus(ctx, repository.UpdateLibrarySyncStatusParams{
			Status:   finalStatus,
			UserID:   payload.UserID,
			AnimeID:  payload.AnimeID,
			Provider: repository.Provider(payload.Provider),
			Action:   repository.LibraryActions(payload.Action),
		})
	}
}

func handleMalProvider(
	ctx context.Context,
	malClient *myanimelist.Client,
	anime repository.Anime,
	token string,
	action string,
	status string,
	episodes int,
) error {
	switch action {
	case string(repository.LibraryActionsAddEntry):
		return malClient.UpdateAnimeList(ctx, myanimelist.UpdateAnimeListParams{
			Token:           token,
			AnimeID:         int(anime.MalID.Int32),
			Status:          status,
			WatchedEpisodes: episodes,
		})

	case string(repository.LibraryActionsUpdateProgress):
		return malClient.UpdateAnimeList(ctx, myanimelist.UpdateAnimeListParams{
			Token:           token,
			AnimeID:         int(anime.MalID.Int32),
			WatchedEpisodes: episodes,
		})

	case string(repository.LibraryActionsUpdateStatus):
		return malClient.UpdateAnimeList(ctx, myanimelist.UpdateAnimeListParams{
			Token:   token,
			AnimeID: int(anime.MalID.Int32),
			Status:  status,
		})

	case string(repository.LibraryActionsDeleteEntry):
		return malClient.DeleteAnimeList(ctx, myanimelist.DeleteAnimeListParams{
			Token:   token,
			AnimeID: int(anime.MalID.Int32),
		})
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}
}

func handleAniProvider(
	ctx context.Context,
	aniClient *anilist.Client,
	anime repository.Anime,
	token string,
	action string,
	status string,
	episodes int,
) error {
	switch action {
	case string(repository.LibraryActionsAddEntry):
		return aniClient.InsertAnimeToList(ctx, anilist.InsertAnimeToListParams{
			Token:           token,
			MalID:           int(anime.MalID.Int32),
			Status:          status,
			WatchedEpisodes: episodes,
		})

	case string(repository.LibraryActionsUpdateProgress):
		return aniClient.UpdateAnimeEntryProgress(ctx, anilist.UpdateAnimeEntryProgressParams{
			Token:           token,
			MalID:           int(anime.MalID.Int32),
			WatchedEpisodes: episodes,
		})

	case string(repository.LibraryActionsUpdateStatus):
		return aniClient.UpdateAnimeEntryStatus(ctx, anilist.UpdateAnimeEntryStatusParams{
			Token:  token,
			MalID:  int(anime.MalID.Int32),
			Status: status,
		})

	case string(repository.LibraryActionsDeleteEntry):
		return aniClient.DeleteAnimeList(ctx, anilist.DeleteAnimeListParams{
			Token: token,
			MalID: int(anime.MalID.Int32),
		})
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}
}
