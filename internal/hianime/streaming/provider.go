package streaming

import "context"

type Provider interface {
	Supports(serverName string) bool
	FetchSources(ctx context.Context, serverID, streamType, serverName string) (ScrapedUnencryptedSources, error)
}
