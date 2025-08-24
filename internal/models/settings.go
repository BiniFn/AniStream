package models


// Settings management models

// SettingsRequest represents the request body for updating user settings
type SettingsRequest struct {
	AutoNextEpisode   bool `json:"autoNextEpisode" example:"true"`
	AutoPlayEpisode   bool `json:"autoPlayEpisode" example:"false"`
	AutoResumeEpisode bool `json:"autoResumeEpisode" example:"true"`
	IncognitoMode     bool `json:"incognitoMode" example:"false"`
}

// Settings response models

// SettingsResponse represents user settings in HTTP responses
type SettingsResponse struct {
	UserID            string `json:"userId" example:"V1StGXR8Z5jdHi6B"`
	AutoNextEpisode   bool   `json:"autoNextEpisode" example:"true"`
	AutoPlayEpisode   bool   `json:"autoPlayEpisode" example:"false"`
	AutoResumeEpisode bool   `json:"autoResumeEpisode" example:"true"`
	IncognitoMode     bool   `json:"incognitoMode" example:"false"`
}

