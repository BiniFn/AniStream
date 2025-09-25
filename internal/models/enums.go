package models

type LibraryStatus string

const (
	LibraryStatusPlanning  LibraryStatus = "planning"
	LibraryStatusWatching  LibraryStatus = "watching"
	LibraryStatusCompleted LibraryStatus = "completed"
	LibraryStatusDropped   LibraryStatus = "dropped"
	LibraryStatusPaused    LibraryStatus = "paused"
)

func (s LibraryStatus) IsValid() bool {
	switch s {
	case LibraryStatusWatching, LibraryStatusCompleted, LibraryStatusPaused, LibraryStatusDropped, LibraryStatusPlanning:
		return true
	default:
		return false
	}
}

type LibraryImportStatus string

const (
	LibraryImportStatusPending    LibraryImportStatus = "pending"
	LibraryImportStatusInProgress LibraryImportStatus = "in_progress"
	LibraryImportStatusCompleted  LibraryImportStatus = "completed"
	LibraryImportStatusFailed     LibraryImportStatus = "failed"
)

func (s LibraryImportStatus) IsValid() bool {
	switch s {
	case LibraryImportStatusPending, LibraryImportStatusInProgress, LibraryImportStatusCompleted, LibraryImportStatusFailed:
		return true
	default:
		return false
	}
}

type OAuthProvider string

const (
	OAuthProviderAnilist OAuthProvider = "anilist"
	OAuthProviderMAL     OAuthProvider = "myanimelist"
)

func (p OAuthProvider) IsValid() bool {
	switch p {
	case OAuthProviderAnilist, OAuthProviderMAL:
		return true
	default:
		return false
	}
}
