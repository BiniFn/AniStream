package models

// LibraryStatus represents the status of an anime in a user's library
type LibraryStatus string

const (
	LibraryStatusWatching    LibraryStatus = "watching"
	LibraryStatusCompleted   LibraryStatus = "completed"
	LibraryStatusOnHold      LibraryStatus = "on_hold"
	LibraryStatusDropped     LibraryStatus = "dropped"
	LibraryStatusPlanToWatch LibraryStatus = "plan_to_watch"
)

// IsValid checks if the library status is valid
func (s LibraryStatus) IsValid() bool {
	switch s {
	case LibraryStatusWatching, LibraryStatusCompleted, LibraryStatusOnHold, LibraryStatusDropped, LibraryStatusPlanToWatch:
		return true
	default:
		return false
	}
}

// OAuthProvider represents supported OAuth providers
type OAuthProvider string

const (
	OAuthProviderAnilist OAuthProvider = "anilist"
	OAuthProviderMAL     OAuthProvider = "mal"
)

// IsValid checks if the OAuth provider is valid
func (p OAuthProvider) IsValid() bool {
	switch p {
	case OAuthProviderAnilist, OAuthProviderMAL:
		return true
	default:
		return false
	}
}