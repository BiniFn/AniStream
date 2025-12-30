package models

import "time"

type DesktopReleaseResponse struct {
	ID           string    `json:"id"`
	Version      string    `json:"version"`
	Platform     string    `json:"platform"`
	DownloadURL  string    `json:"downloadUrl"`
	FileName     string    `json:"fileName"`
	FileSize     int64     `json:"fileSize"`
	ReleaseNotes *string   `json:"releaseNotes,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

type DesktopPlatformRelease struct {
	Platform    string `json:"platform"`
	DownloadURL string `json:"downloadUrl"`
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"`
}

type DesktopVersionResponse struct {
	Version      string                   `json:"version"`
	ReleaseNotes *string                  `json:"releaseNotes,omitempty"`
	CreatedAt    time.Time                `json:"createdAt"`
	Platforms    []DesktopPlatformRelease `json:"platforms"`
}

type CreateDesktopReleaseRequest struct {
	Version      string `json:"version" validate:"required"`
	Platform     string `json:"platform" validate:"required,oneof=darwin-arm64 darwin-x64 win32-x64 win32-arm64 linux-x64 linux-arm64"`
	DownloadURL  string `json:"downloadUrl" validate:"required,url"`
	FileName     string `json:"fileName" validate:"required"`
	FileSize     int64  `json:"fileSize" validate:"required,gt=0"`
	ReleaseNotes string `json:"releaseNotes"`
}
