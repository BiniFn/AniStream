package settings

import (
	"context"
	"errors"

	"github.com/coeeter/aniways/internal/mappers"
	"github.com/coeeter/aniways/internal/models"
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

func (s *SettingsService) GetSettings(ctx context.Context, userID string) (models.SettingsResponse, error) {
	setting, err := s.repo.GetSettingsOfUser(ctx, userID)
	switch {
	case err == nil:
		return mappers.SettingsFromRepository(setting), nil
	case errors.Is(err, pgx.ErrNoRows):
		setting, err := s.repo.SaveSettings(ctx, repository.SaveSettingsParams{
			UserID:            userID,
			AutoNextEpisode:   true,
			AutoPlayEpisode:   true,
			AutoResumeEpisode: false,
			IncognitoMode:     false,
		})
		if err != nil {
			return models.SettingsResponse{}, err
		}
		return mappers.SettingsFromRepository(setting), nil
	default:
		return models.SettingsResponse{}, err
	}
}

type SaveSettingsParams struct {
	UserID            string
	AutoNextEpisode   bool
	AutoPlayEpisode   bool
	AutoResumeEpisode bool
	IncognitoMode     bool
}

func (s *SettingsService) SaveSettings(ctx context.Context, params SaveSettingsParams) (models.SettingsResponse, error) {
	settings, err := s.repo.SaveSettings(ctx, repository.SaveSettingsParams{
		UserID:            params.UserID,
		AutoNextEpisode:   params.AutoNextEpisode,
		AutoPlayEpisode:   params.AutoPlayEpisode,
		AutoResumeEpisode: params.AutoResumeEpisode,
		IncognitoMode:     params.IncognitoMode,
	})

	if err != nil {
		return models.SettingsResponse{}, err
	}

	return mappers.SettingsFromRepository(settings), nil
}
