package streaming

import (
	"context"
	"fmt"
	"net/http"
)

type Fetcher struct {
	providers []Provider
}

func NewFetcher(httpClient *http.Client) *Fetcher {
	return &Fetcher{
		providers: []Provider{
			NewHianimeProvider(httpClient),
			NewMegaplayProvider(httpClient),
		},
	}
}

func (f *Fetcher) Fetch(
	ctx context.Context,
	serverName, serverID, streamType string,
) (ScrapedUnencryptedSources, error) {
	for _, p := range f.providers {
		if p.Supports(serverName) {
			return p.FetchSources(ctx, serverID, streamType, serverName)
		}
	}
	return ScrapedUnencryptedSources{}, fmt.Errorf("no provider for %q", serverName)
}
