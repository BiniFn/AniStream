package settings

import (
	"context"
	"database/sql"
	"errors"

	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5"
)

type SettingsService struct {
	repo *repository.Queries
}

func NewSettingsService(repo *repository.Queries) *SettingsService {
	return &SettingsService{
		repo: repo,
	}
}

func (s *SettingsService) GetSettings(ctx context.Context, userID string) (Settings, error) {
	setting, err := s.repo.GetSettingsOfUser(ctx, userID)
	switch {
	case err == nil:
		return Settings{}.FromRepository(setting), nil
	case errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows):
		setting, err := s.repo.SaveSettings(ctx, repository.SaveSettingsParams{
			UserID:            userID,
			AutoNextEpisode:   true,
			AutoPlayEpisode:   true,
			AutoResumeEpisode: false,
			IncognitoMode:     false,
		})
		if err != nil {
			return Settings{}, err
		}
		return Settings{}.FromRepository(setting), nil
	default:
		return Settings{}, err
	}
}

type SaveSettingsParams struct {
	UserID            string
	AutoNextEpisode   bool
	AutoPlayEpisode   bool
	AutoResumeEpisode bool
	IncognitoMode     bool
}

func (s *SettingsService) SaveSettings(ctx context.Context, params SaveSettingsParams) (Settings, error) {
	settings, err := s.repo.SaveSettings(ctx, repository.SaveSettingsParams{
		UserID:            params.UserID,
		AutoNextEpisode:   params.AutoNextEpisode,
		AutoPlayEpisode:   params.AutoPlayEpisode,
		AutoResumeEpisode: params.AutoResumeEpisode,
		IncognitoMode:     params.IncognitoMode,
	})

	if err != nil {
		return Settings{}, err
	}

	return Settings{}.FromRepository(settings), nil
}
