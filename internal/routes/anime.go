package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/coeeter/aniways/internal/hianime"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/go-chi/chi/v5"
)

func mountAnimeRoutes(r *chi.Mux, repo *repository.Queries) {
	scraper := hianime.NewHianimeScraper()
	r.Route("/anime", func(r chi.Router) {
		r.Get("/recently-updated", func(w http.ResponseWriter, r *http.Request) {
			page := r.URL.Query().Get("page")
			if page == "" {
				page = "1" // Default to page 1 if not specified
			}

			ctx := r.Context()
			pageNum, err := strconv.Atoi(page)
			if err != nil {
				http.Error(w, "Invalid page number", http.StatusBadRequest)
				return
			}
			animeList, err := scraper.GetRecentlyUpdatedAnime(ctx, pageNum)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(animeList); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/{hiAnimeID}", func(w http.ResponseWriter, r *http.Request) {
			hiAnimeID := chi.URLParam(r, "hiAnimeID")
			if hiAnimeID == "" {
				http.Error(w, "HiAnime ID is required", http.StatusBadRequest)
				return
			}

			ctx := r.Context()
			animeInfo, err := scraper.GetAnimeInfoByHiAnimeID(ctx, hiAnimeID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(animeInfo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/{hiAnimeID}/episodes", func(w http.ResponseWriter, r *http.Request) {
			hiAnimeID := chi.URLParam(r, "hiAnimeID")
			if hiAnimeID == "" {
				http.Error(w, "HiAnime ID is required", http.StatusBadRequest)
				return
			}

			ctx := r.Context()
			episodes, err := scraper.GetAnimeEpisodes(ctx, hiAnimeID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(episodes); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/{hiAnimeID}/servers/{episodeID}", func(w http.ResponseWriter, r *http.Request) {
			hiAnimeID := chi.URLParam(r, "hiAnimeID")
			episodeID := chi.URLParam(r, "episodeID")
			if episodeID == "" {
				http.Error(w, "Episode ID is required", http.StatusBadRequest)
				return
			}
			if hiAnimeID == "" {
				http.Error(w, "HiAnime ID is required", http.StatusBadRequest)
				return
			}

			ctx := r.Context()
			servers, err := scraper.GetEpisodeServers(ctx, hiAnimeID, episodeID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(servers); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	})
}
