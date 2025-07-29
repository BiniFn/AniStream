package oauth

import (
	"context"

	"github.com/coeeter/aniways/internal/cache"
)

type AnilistProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
	redis        *cache.RedisClient
}

func NewAnilistProvider(clientID, clientSecret, redirectURL string, redis *cache.RedisClient) *AnilistProvider {
	return &AnilistProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		redis:        redis,
	}
}

func (a *AnilistProvider) Name() string {
	return "anilist"
}

func (a *AnilistProvider) AuthURL(ctx context.Context, state string) (string, error) {
	panic("unimplemented")
}

func (a *AnilistProvider) ExchangeToken(ctx context.Context, state string, code string) (TokenResponse, error) {
	panic("unimplemented")
}

var _ Provider = (*AnilistProvider)(nil)
