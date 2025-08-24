package models

type LibraryStatus string

const (
	LibraryStatusWatching    LibraryStatus = "watching"
	LibraryStatusCompleted   LibraryStatus = "completed"
	LibraryStatusOnHold      LibraryStatus = "on_hold"
	LibraryStatusDropped     LibraryStatus = "dropped"
	LibraryStatusPlanToWatch LibraryStatus = "plan_to_watch"
)

func (s LibraryStatus) IsValid() bool {
	switch s {
	case LibraryStatusWatching, LibraryStatusCompleted, LibraryStatusOnHold, LibraryStatusDropped, LibraryStatusPlanToWatch:
		return true
	default:
		return false
	}
}

type OAuthProvider string

const (
	OAuthProviderAnilist OAuthProvider = "anilist"
	OAuthProviderMAL     OAuthProvider = "mal"
)

func (p OAuthProvider) IsValid() bool {
	switch p {
	case OAuthProviderAnilist, OAuthProviderMAL:
		return true
	default:
		return false
	}
}

