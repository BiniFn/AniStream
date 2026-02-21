package models

type SettingsRequest struct {
	AutoNextEpisode   bool `json:"autoNextEpisode" example:"true"`
	AutoPlayEpisode   bool `json:"autoPlayEpisode" example:"false"`
	AutoResumeEpisode bool `json:"autoResumeEpisode" example:"true"`
	IncognitoMode     bool `json:"incognitoMode" example:"false"`
	ThemeId           int  `json:"themeId" example:"1"`
}

type Theme struct {
	ID          int    `json:"id" validate:"required" example:"1"`
	Name        string `json:"name" validate:"required" example:"Default"`
	Description string `json:"description" validate:"required" example:"The default AniStream theme."`
	ClassName   string `json:"className" validate:"required" example:"theme-default"`
}

type SettingsResponse struct {
	UserID            string `json:"userId" validate:"required" example:"V1StGXR8Z5jdHi6B"`
	AutoNextEpisode   bool   `json:"autoNextEpisode" validate:"required" example:"true"`
	AutoPlayEpisode   bool   `json:"autoPlayEpisode" validate:"required" example:"false"`
	AutoResumeEpisode bool   `json:"autoResumeEpisode" validate:"required" example:"true"`
	IncognitoMode     bool   `json:"incognitoMode" validate:"required" example:"false"`
	Theme             Theme  `json:"theme" validate:"required"`
}
