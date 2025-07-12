package api

import (
	"net/http"

	"github.com/coeeter/aniways/internal/api/handlers"
	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/myanimelist"
	"github.com/coeeter/aniways/internal/repository"
	animeSvc "github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func MountGlobalRoutes(r *chi.Mux, env *config.Env, repo *repository.Queries, redis *cache.RedisClient) {
	r.Route("/anime", func(r chi.Router) {
		malClient := myanimelist.NewClient(myanimelist.ClientConfig{
			ClientID:     env.MyAnimeListClientID,
			ClientSecret: env.MyAnimeListClientSecret,
		})
		refresher := animeSvc.NewRefresher(repo, malClient)
		svc := animeSvc.New(repo, refresher, malClient, redis)

		handlers.MountAnimeRoutes(r, svc)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("AniWays API"))
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}
