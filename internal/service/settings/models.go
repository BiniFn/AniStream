package settings

import "github.com/coeeter/aniways/internal/repository"

type Settings struct {
	UserID            string `json:"userId"`
	AutoNextEpisode   bool   `json:"autoNextEpisode"`
	AutoPlayEpisode   bool   `json:"autoPlayEpisode"`
	AutoResumeEpisode bool   `json:"autoResumeEpisode"`
	IncognitoMode     bool   `json:"incognitoMode"`
}

func (s Settings) FromRepository(r repository.Setting) Settings {
	return Settings{
		UserID:            r.UserID,
		AutoNextEpisode:   r.AutoNextEpisode,
		AutoPlayEpisode:   r.AutoPlayEpisode,
		AutoResumeEpisode: r.AutoResumeEpisode,
	}
}
