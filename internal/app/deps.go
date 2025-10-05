package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/infra/cache"
	"github.com/coeeter/aniways/internal/infra/client/anilist"
	"github.com/coeeter/aniways/internal/infra/client/hianime"
	"github.com/coeeter/aniways/internal/infra/client/jikan"
	"github.com/coeeter/aniways/internal/infra/client/myanimelist"
	"github.com/coeeter/aniways/internal/infra/client/shikimori"
	"github.com/coeeter/aniways/internal/infra/database"
	"github.com/coeeter/aniways/internal/infra/email"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/service/auth/oauth"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Deps struct {
	Env         *config.Env
	Log         *slog.Logger
	Db          *pgxpool.Pool
	Repo        *repository.Queries
	Cache       *cache.RedisClient
	Scraper     *hianime.HianimeScraper
	MAL         *myanimelist.Client
	Jikan       *jikan.Client
	Anilist     *anilist.Client
	Shiki       *shikimori.Client
	Cld         *cloudinary.Cloudinary
	EmailClient email.EmailClient
	Providers   []oauth.Provider
}

func InitDeps(ctx context.Context, svcName string) (*Deps, error) {
	env, err := config.LoadEnv()
	rootLogger := NewLogger(svcName) // if env is loaded wrongly just use json handler (prod way)
	if err != nil {
		return &Deps{Log: rootLogger}, err
	}

	dbLog := rootLogger.With("component", "database")
	db, err := database.New(env, dbLog)
	if err != nil {
		return &Deps{Log: rootLogger}, err
	}

	cacheLog := rootLogger.With("component", "cache")
	cache, err := cache.NewRedisClient(ctx, env.AppEnv, env.RedisAddr, env.RedisPassword, cacheLog)
	if err != nil {
		return &Deps{Log: rootLogger}, err
	}

	repo := repository.New(db)
	scraper := hianime.NewHianimeScraper()
	malClient := myanimelist.NewClient(env.MyAnimeListClientID)
	jikanClient := jikan.NewClient()
	anilistClient := anilist.New()
	shiki := shikimori.NewClient(cache)
	emailClient := email.NewClient(env.ResendAPIKey, env.ResendFromEmail)

	cld, err := cloudinary.NewFromParams(env.CloudinaryName, env.CloudinaryAPIKey, env.CloudinaryAPISecret)
	if err != nil {
		return &Deps{Log: rootLogger}, err
	}

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

	return &Deps{
		Env:         env,
		Log:         rootLogger,
		Db:          db,
		Repo:        repo,
		Cache:       cache,
		Scraper:     scraper,
		MAL:         malClient,
		Jikan:       jikanClient,
		Anilist:     anilistClient,
		Shiki:       shiki,
		Cld:         cld,
		EmailClient: emailClient,
		Providers:   []oauth.Provider{malOauthProvider, anilistOauthProvider},
	}, nil
}

func (d *Deps) Close() error {
	d.Db.Close()
	d.Log.Info("Closed db connections", "component", "database")

	if err := d.Cache.Close(); err != nil {
		d.Log.Error("Failed to close cache", "error", err)
		return err
	}
	d.Log.Info("Closed cache connections", "component", "cache")

	return nil
}
