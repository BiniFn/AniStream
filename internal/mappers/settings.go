package mappers

import (
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
)

func SettingsFromRepository(r repository.Setting) models.SettingsResponse {
	return models.SettingsResponse{
		UserID:            r.UserID,
		AutoNextEpisode:   r.AutoNextEpisode,
		AutoPlayEpisode:   r.AutoPlayEpisode,
		AutoResumeEpisode: r.AutoResumeEpisode,
		IncognitoMode:     r.IncognitoMode,
	}
}
