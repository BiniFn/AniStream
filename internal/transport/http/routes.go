package http

import (
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/client/anilist"
	"github.com/coeeter/aniways/internal/client/myanimelist"
	"github.com/coeeter/aniways/internal/client/shikimori"
	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/repository"
	animeSvc "github.com/coeeter/aniways/internal/service/anime"
	"github.com/coeeter/aniways/internal/service/users"
	"github.com/coeeter/aniways/internal/transport/http/handlers"
	"github.com/go-chi/chi/v5"
)

func MountGlobalRoutes(r *chi.Mux, env *config.Env, repo *repository.Queries, redis *cache.RedisClient) {
	r.Route("/anime", func(r chi.Router) {
		malClient := myanimelist.NewClient(myanimelist.ClientConfig{
			ClientID:     env.MyAnimeListClientID,
			ClientSecret: env.MyAnimeListClientSecret,
		})
		refresher := animeSvc.NewRefresher(repo, malClient)
		anilistClient := anilist.New()
		shikimoriClient := shikimori.NewClient(redis)
		svc := animeSvc.New(repo, refresher, malClient, anilistClient, shikimoriClient, redis)

		handlers.MountAnimeRoutes(r, svc)
		handlers.MountAnimeListingsRoutes(r, svc)
		handlers.MountAnimeEpisodesRoutes(r, svc)
	})

	r.Route("/users", func(r chi.Router) {
		cld, _ := cloudinary.NewFromParams(env.CloudinaryName, env.CloudinaryAPIKey, env.CloudinaryAPISecret)
		userService := users.NewUserService(repo, cld)
		handlers.MountUsersRoutes(r, userService)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("AniWays API"))
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}
