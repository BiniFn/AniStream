package http

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/client/anilist"
	"github.com/coeeter/aniways/internal/client/myanimelist"
	"github.com/coeeter/aniways/internal/client/shikimori"
	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/repository"
)

type Dependencies struct {
	Env     *config.Env
	Repo    *repository.Queries
	Cache   *cache.RedisClient
	MAL     *myanimelist.Client
	Anilist *anilist.Client
	Shiki   *shikimori.Client
	Cld     *cloudinary.Cloudinary
}

func BuildDeps(
	env *config.Env,
	repo *repository.Queries,
	cache *cache.RedisClient,
) (*Dependencies, error) {
	mal := myanimelist.NewClient(myanimelist.ClientConfig{
		ClientID:     env.MyAnimeListClientID,
		ClientSecret: env.MyAnimeListClientSecret,
	})

	anilist := anilist.New()

	shiki := shikimori.NewClient(cache)

	cld, err := cloudinary.NewFromParams(env.CloudinaryName, env.CloudinaryAPIKey, env.CloudinaryAPISecret)
	if err != nil {
		return nil, err
	}

	return &Dependencies{
		Env:     env,
		Repo:    repo,
		Cache:   cache,
		MAL:     mal,
		Anilist: anilist,
		Shiki:   shiki,
		Cld:     cld,
	}, nil
}
