package http

import (
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/client/anilist"
	"github.com/coeeter/aniways/internal/client/myanimelist"
	"github.com/coeeter/aniways/internal/client/shikimori"
	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/email"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/service/auth/oauth"
)

type Dependencies struct {
	Env         *config.Env
	Repo        *repository.Queries
	Cache       *cache.RedisClient
	MAL         *myanimelist.Client
	Anilist     *anilist.Client
	Shiki       *shikimori.Client
	Cld         *cloudinary.Cloudinary
	EmailClient email.EmailClient
	Providers   []oauth.Provider
}

func BuildDeps(
	env *config.Env,
	repo *repository.Queries,
	cache *cache.RedisClient,
) (*Dependencies, error) {
	malClient := myanimelist.NewClient(myanimelist.ClientConfig{
		ClientID:     env.MyAnimeListClientID,
		ClientSecret: env.MyAnimeListClientSecret,
	})

	anilistClient := anilist.New()

	shiki := shikimori.NewClient(cache)

	cld, err := cloudinary.NewFromParams(env.CloudinaryName, env.CloudinaryAPIKey, env.CloudinaryAPISecret)
	if err != nil {
		return nil, err
	}

	emailClient := email.NewClient(env.ResendAPIKey, env.ResendFromEmail)

	malOauthProvider := oauth.NewMALProvider(
		env.MyAnimeListClientID,
		env.MyAnimeListClientSecret,
		fmt.Sprintf("%s/auth/oauth/myanimelist/callback", env.ApiURL),
		repo,
		cache,
	)

	anilistOauthProvider := oauth.NewAnilistProvider(
		env.AnilistClientID,
		env.AnilistClientSecret,
		fmt.Sprintf("%s/auth/oauth/anilist/callback", env.ApiURL),
		repo,
	)

	return &Dependencies{
		Env:         env,
		Repo:        repo,
		Cache:       cache,
		MAL:         malClient,
		Anilist:     anilistClient,
		Shiki:       shiki,
		Cld:         cld,
		EmailClient: emailClient,
		Providers:   []oauth.Provider{malOauthProvider, anilistOauthProvider},
	}, nil
}
