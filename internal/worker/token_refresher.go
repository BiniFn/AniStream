package worker

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/service/auth/oauth"
	"github.com/jackc/pgx/v5"
)

func dailyTask(
	ctx context.Context,
	repo *repository.Queries,
	providers []oauth.Provider,
	log *slog.Logger,
) {
	log.Info("Running daily task")
	if err := refreshAccessTokens(ctx, repo, providers, log); err != nil {
		log.Error("Error in daily task", "err", err)
	} else {
		log.Info("Daily task completed successfully")
	}
}

func refreshAccessTokens(
	ctx context.Context,
	repo *repository.Queries,
	providers []oauth.Provider,
	log *slog.Logger,
) error {
	tokens, err := repo.GetTokensNearToExpiry(ctx)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		log.Info("No tokens near to expiry")
		return nil
	}
	if err != nil {
		return err
	}

	var malProvider *oauth.MALProvider
	for _, provider := range providers {
		if p, ok := provider.(*oauth.MALProvider); ok {
			malProvider = p
			break
		}
	}

	for _, token := range tokens {
		if token.Provider == repository.ProviderAnilist {
			continue
		}

		err := malProvider.RefreshToken(ctx, token.UserID, token.RefreshToken)
		if err != nil {
			return err
		}

		log.Info("Token refreshed", "user_id", token.UserID, "provider", token.Provider)
	}

	return nil
}
