package mappers

import (
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
)

func SettingsFromRepository(r repository.GetSettingsOfUserRow) models.SettingsResponse {
	return models.SettingsResponse{
		UserID:            r.Setting.UserID,
		AutoNextEpisode:   r.Setting.AutoNextEpisode,
		AutoPlayEpisode:   r.Setting.AutoPlayEpisode,
		AutoResumeEpisode: r.Setting.AutoResumeEpisode,
		IncognitoMode:     r.Setting.IncognitoMode,
		Theme:             ThemesFromRepository(r.Theme),
	}
}

func SettingsFromSaveRepository(r repository.SaveSettingsRow) models.SettingsResponse {
	return models.SettingsResponse{
		UserID:            r.UserID,
		AutoNextEpisode:   r.AutoNextEpisode,
		AutoPlayEpisode:   r.AutoPlayEpisode,
		AutoResumeEpisode: r.AutoResumeEpisode,
		IncognitoMode:     r.IncognitoMode,
		Theme:             ThemesFromRepository(r.Theme),
	}
}

func ThemesFromRepository(r repository.Theme) models.Theme {
	return models.Theme{
		ID:          int(r.ID),
		Name:        r.Name,
		Description: r.Description.String,
		ClassName:   r.ThemeClass,
	}
}
