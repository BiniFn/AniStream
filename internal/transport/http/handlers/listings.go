package handlers

import (
	"net/http"
	"strconv"

	"github.com/coeeter/aniways/internal/models"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) AnimeListingRoutes() {
	h.r.Route("/anime/listings", func(r chi.Router) {
		r.With(h.catalogInput()).Get("/", h.catalog)
		r.Get("/recently-updated", h.listRecentlyUpdated)
		r.Get("/seasonal", h.seasonalAnimes)
		r.Get("/seasons", h.getBySeason)
		r.Get("/random", h.randomAnime)
		r.Get("/genres", h.listGenres)
		r.Get("/genres/{genre}", h.animeByGenre)
		r.Get("/search", h.searchAnimes)
		r.Get("/trending", h.trendingAnimes)
		r.Get("/popular", h.popularAnimes)
	})
}

func (h *Handler) catalogInput() func(http.Handler) http.Handler {
	errorHandler := httpin.Option.WithErrorHandler(func(w http.ResponseWriter, _ *http.Request, err error) {
		h.jsonError(w, http.StatusBadRequest, err.Error())
	})

	return httpin.NewInput(models.GetAnimeCatalogParams{}, errorHandler)
}

// @Summary Get anime catalog
// @Description Get anime catalog
// @Tags Anime Listings
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param itemsPerPage query int false "Number of items per page"
// @Param search query string false "Search text"
// @Param genres query []string false "Genres (repeat param)" collectionFormat(multi)
// @Param genresMode query string false "Genre match mode" Enums(any,all)
// @Param seasons query []string false "Seasons (repeat param)" Enums(winter,spring,summer,fall,unknown) collectionFormat(multi)
// @Param years query []int false "Years (repeat param)" collectionFormat(multi)
// @Param yearMin query int false "Minimum year (inclusive)"
// @Param yearMax query int false "Maximum year (inclusive)"
// @Param sortBy query string false "Sort field" Enums(ename,jname,season,year,relevance,updated_at)
// @Param sortOrder query string false "Sort order" Enums(asc,desc)
// @Success 200 {object} models.AnimeListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/listings [get]
func (h *Handler) catalog(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*models.GetAnimeCatalogParams)

	resp, err := h.animeService.GetAnimeCatalog(r.Context(), input)
	if err != nil {
		log.Error("failed to fetch anime catalog", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch anime catalog")
		return
	}

	h.jsonOK(w, resp)
}

// @Summary Get recently updated anime
// @Description Get recently updated anime
// @Tags Anime Listings
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param itemsPerPage query int false "Number of items per page"
// @Success 200 {object} models.AnimeListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/listings/recently-updated [get]
func (h *Handler) listRecentlyUpdated(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	page, size, err := h.parsePagination(r, 1, 30)
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.GetRecentlyUpdatedAnimes(r.Context(), page, size)
	if err != nil {
		log.Error("failed to fetch recently updated animes", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch recently updated animes")
		return
	}
	h.jsonOK(w, resp)
}

// @Summary Get seasonal anime
// @Description Get seasonal anime
// @Tags Anime Listings
// @Accept json
// @Produce json
// @Success 200 {object} models.SeasonalAnimeListResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/listings/seasonal [get]
func (h *Handler) seasonalAnimes(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	resp, err := h.animeService.GetSeasonalAnimes(r.Context())
	if err != nil {
		log.Error("failed to fetch seasonal animes", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch seasonal animes")
		return
	}
	h.jsonOK(w, resp)
}

// @Summary Get anime by season and/or year
// @Description Get anime by season and/or year
// @Tags Anime Listings
// @Accept json
// @Produce json
// @Param season query string false "Season filter"
// @Param year query int false "Year filter"
// @Param page query int false "Page number"
// @Param itemsPerPage query int false "Number of items per page"
// @Success 200 {object} models.AnimeListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/listings/seasons [get]
func (h *Handler) getBySeason(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	page, size, err := h.parsePagination(r, 1, 30)
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	season := r.URL.Query().Get("season")
	year := r.URL.Query().Get("year")

	if season == "" && year == "" {
		h.jsonError(w, http.StatusBadRequest, "season and year are required")
		return
	}

	seasonYear, err := strconv.Atoi(year)
	if err != nil && year != "" {
		h.jsonError(w, http.StatusBadRequest, "invalid year")
		return
	}

	var resp models.AnimeListResponse

	if season != "" && year != "" {
		resp, err = h.animeService.GetAnimeBySeasonAndYear(r.Context(), season, int32(seasonYear), page, size)
	} else if season != "" {
		resp, err = h.animeService.GetAnimeBySeason(r.Context(), season, page, size)
	} else if year != "" {
		resp, err = h.animeService.GetAnimeByYear(r.Context(), int32(seasonYear), page, size)
	}

	if err != nil {
		log.Error("failed to fetch anime by season", "season", season, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch anime by season")
		return
	}
	h.jsonOK(w, resp)
}

// @Summary Get random anime
// @Description Get random anime
// @Tags Anime Listings
// @Accept json
// @Produce json
// @Param genre query string false "Optional genre filter"
// @Success 200 {object} models.AnimeResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/listings/random [get]
func (h *Handler) randomAnime(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	genre := r.URL.Query().Get("genre")

	var (
		resp models.AnimeResponse
		err  error
	)

	if genre != "" {
		resp, err = h.animeService.GetRandomAnimeByGenre(r.Context(), genre)
		if err != nil {
			log.Error("failed to fetch random anime by genre", "genre", genre, "err", err)
			h.jsonError(w, http.StatusInternalServerError, "failed to fetch random anime by genre")
			return
		}
	} else {
		resp, err = h.animeService.GetRandomAnime(r.Context())
		if err != nil {
			log.Error("failed to fetch random animes", "err", err)
			h.jsonError(w, http.StatusInternalServerError, "failed to fetch random animes")
			return
		}
	}

	h.jsonOK(w, resp)
}

// @Summary Get list of anime genres
// @Description Get list of anime genres
// @Tags Anime Listings
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/listings/genres [get]
func (h *Handler) listGenres(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	resp, err := h.animeService.GetAnimeGenres(r.Context())
	if err != nil {
		log.Error("failed to fetch anime genres", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch anime genres")
		return
	}
	h.jsonOK(w, resp)
}

// @Summary Search anime
// @Description Search anime
// @Tags Anime Listings
// @Accept json
// @Produce json
// @Param q query string true "Search query (minimum 3 characters)"
// @Param genre query string false "Optional genre filter (minimum 3 characters)"
// @Param page query int false "Page number"
// @Param itemsPerPage query int false "Number of items per page"
// @Success 200 {object} models.AnimeListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/listings/search [get]
func (h *Handler) searchAnimes(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	query := r.URL.Query().Get("q")
	if query == "" {
		h.jsonError(w, http.StatusBadRequest, "search query is required")
		return
	}
	if len(query) < 3 {
		h.jsonError(w, http.StatusBadRequest, "search query must be at least 3 characters long")
		return
	}

	genre := r.URL.Query().Get("genre")
	if genre != "" && len(genre) < 3 {
		h.jsonError(w, http.StatusBadRequest, "genre query must be at least 3 characters long")
		return
	}

	page, size, err := h.parsePagination(r, 1, 30)
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.SearchAnimes(r.Context(), query, genre, page, size)
	if err != nil {
		log.Error("failed to search animes", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to search animes")
		return
	}
	h.jsonOK(w, resp)
}

// @Summary Get anime by genre
// @Description Get anime by genre
// @Tags Anime Listings
// @Accept json
// @Produce json
// @Param genre path string true "Genre name"
// @Param page query int false "Page number"
// @Param itemsPerPage query int false "Number of items per page"
// @Success 200 {object} models.AnimeListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/listings/genres/{genre} [get]
func (h *Handler) animeByGenre(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	genre, err := h.pathParam(r, "genre")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	page, size, err := h.parsePagination(r, 1, 30)
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.GetAnimesByGenre(r.Context(), genre, page, size)
	if err != nil {
		log.Error("failed to fetch animes by genre", "genre", genre, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch animes by genre")
		return
	}
	h.jsonOK(w, resp)
}

// @Summary Get trending anime
// @Description Get trending anime
// @Tags Anime Listings
// @Accept json
// @Produce json
// @Success 200 {object} models.TrendingAnimeListResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/listings/trending [get]
func (h *Handler) trendingAnimes(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	resp, err := h.animeService.GetTrendingAnimes(r.Context())
	if err != nil {
		log.Error("failed to fetch trending animes", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch trending animes")
		return
	}
	h.jsonOK(w, resp)
}

// @Summary Get popular anime
// @Description Get popular anime
// @Tags Anime Listings
// @Accept json
// @Produce json
// @Success 200 {object} models.PopularAnimeListResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/listings/popular [get]
func (h *Handler) popularAnimes(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	resp, err := h.animeService.GetPopularAnimes(r.Context())
	if err != nil {
		log.Error("failed to fetch popular animes", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch popular animes")
		return
	}
	h.jsonOK(w, resp)
}
