package handlers

import (
	"net/http"
	"strconv"

	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/transport/http/middleware"
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
		r.Get("/genres/previews", h.genrePreviews)
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

// @Summary Get anime catalog with optional library information
// @Description Get anime catalog. When authenticated, includes library status. Supports library-only filtering and search.
// @Tags Anime Listings
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param itemsPerPage query int false "Number of items per page (default: 30, max: 100)"
// @Param search query string false "Search anime by name"
// @Param genres query []string false "Filter by genres (repeat for multiple)" collectionFormat(multi)
// @Param genresMode query string false "Genre matching mode: 'any' (default) or 'all'" Enums(any,all)
// @Param seasons query []string false "Filter by seasons (repeat for multiple)" Enums(winter,spring,summer,fall,unknown) collectionFormat(multi)
// @Param years query []int false "Filter by specific years (repeat for multiple)" collectionFormat(multi)
// @Param yearMin query int false "Filter by minimum year (inclusive)"
// @Param yearMax query int false "Filter by maximum year (inclusive)"
// @Param sortBy query string false "Sort field" Enums(ename,jname,season,year,relevance,updated_at,anime_updated_at,library_updated_at)
// @Param sortOrder query string false "Sort order: 'asc' or 'desc' (default: 'desc')" Enums(asc,desc)
// @Param inLibraryOnly query bool false "Only show anime in user's library (requires authentication)"
// @Param status query string false "Filter by library status (requires authentication)" Enums(watching,completed,planning,dropped,paused)
// @Success 200 {object} models.AnimeWithLibraryListResponse "Anime catalog with optional library information"
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Authentication required for library features"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /anime/listings [get]
func (h *Handler) catalog(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*models.GetAnimeCatalogParams)

	var userID *string
	if input.InLibraryOnly != nil && *input.InLibraryOnly {
		user := middleware.GetUser(r)
		if user == nil {
			h.jsonError(w, http.StatusUnauthorized, "authentication required for library access")
			return
		}
		userID = &user.ID
	} else if input.Status != nil {
		user := middleware.GetUser(r)
		if user == nil {
			h.jsonError(w, http.StatusUnauthorized, "authentication required for status filtering")
			return
		}
		userID = &user.ID
	}

	resp, err := h.animeService.GetAnimeCatalog(r.Context(), input, userID)
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

// @Summary Get genre previews
// @Description Get 6 preview image URLs for each genre
// @Tags Anime Listings
// @Accept json
// @Produce json
// @Success 200 {array} models.GenrePreview
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/listings/genres/previews [get]
func (h *Handler) genrePreviews(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	resp, err := h.animeService.GetGenrePreviews(r.Context())
	if err != nil {
		log.Error("failed to fetch genre previews", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch genre previews")
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
