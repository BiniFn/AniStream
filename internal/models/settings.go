package models

// Settings management models

// SettingsRequest represents the request body for updating user settings
type SettingsRequest struct {
	AutoNextEpisode   bool `json:"autoNextEpisode"`
	AutoPlayEpisode   bool `json:"autoPlayEpisode"`
	AutoResumeEpisode bool `json:"autoResumeEpisode"`
	IncognitoMode     bool `json:"incognitoMode"`
}