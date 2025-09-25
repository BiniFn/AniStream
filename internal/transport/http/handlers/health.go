package handlers

import (
	"encoding/json"
	"maps"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type HealthResponse struct {
	Status    string            `json:"status" example:"healthy"`
	Timestamp string            `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version" example:"1.0.0"`
}

type HealthCheckResponse struct {
	Status    string `json:"status" example:"healthy"`
	Service   string `json:"service" example:"api"`
	Timestamp string `json:"timestamp" example:"2023-01-01T00:00:00Z"`
}

// Healthz godoc
// @Summary		Comprehensive health check
// @Description	Check the health of all services and external APIs
// @Tags			health
// @Accept			json
// @Produce		json
// @Success		200	{object}	HealthResponse
// @Failure		503	{object}	HealthResponse
// @Router			/health/z [get]
func (h *Handler) healthz(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dbStatus := "healthy"
	if err := h.deps.Db.Ping(ctx); err != nil {
		dbStatus = "unhealthy"
		h.logger(r).Error("Database health check failed", "error", err)
	}

	redisStatus := "healthy"
	if h.deps.Cache != nil {
		if err := h.deps.Cache.Set(ctx, "health_check", "ping", time.Second); err != nil {
			redisStatus = "unhealthy"
			h.logger(r).Error("Redis health check failed", "error", err)
		} else {
			h.deps.Cache.Del(ctx, "health_check")
		}
	} else {
		redisStatus = "not configured"
	}

	externalAPIs := map[string]string{
		"myanimelist": "healthy",
		"anilist":     "healthy",
		"shikimori":   "healthy",
		"hianime":     "healthy",
	}

	if h.deps.MAL != nil {
		if _, err := h.deps.MAL.GetAnimeMetadata(ctx, 16498); err != nil {
			externalAPIs["myanimelist"] = "unhealthy"
			h.logger(r).Error("MyAnimeList API health check failed", "error", err)
		}
	}

	if h.deps.Anilist != nil {
		if _, err := h.deps.Anilist.GetTrendingAnime(ctx); err != nil {
			externalAPIs["anilist"] = "unhealthy"
			h.logger(r).Error("AniList API health check failed", "error", err)
		}
	}

	if h.deps.Shiki != nil {
		if _, err := h.deps.Shiki.GetAnimeFranchise(ctx, 16498); err != nil {
			externalAPIs["shikimori"] = "unhealthy"
			h.logger(r).Error("Shikimori API health check failed", "error", err)
		}
	}

	if h.deps.Scraper != nil {
		if _, err := h.deps.Scraper.GetAZList(ctx, 1); err != nil {
			externalAPIs["hianime"] = "unhealthy"
			h.logger(r).Error("HiAnime scraper health check failed", "error", err)
		}
	}

	overallStatus := "healthy"
	if dbStatus == "unhealthy" || redisStatus == "unhealthy" {
		overallStatus = "unhealthy"
	}

	for _, status := range externalAPIs {
		if status == "unhealthy" {
			overallStatus = "degraded"
			break
		}
	}

	status := http.StatusOK
	switch overallStatus {
	case "unhealthy":
		status = http.StatusServiceUnavailable
	case "degraded":
		status = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	healthResponse := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Services: map[string]string{
			"database": dbStatus,
			"redis":    redisStatus,
		},
		Version: "1.0.0",
	}

	maps.Copy(healthResponse.Services, externalAPIs)

	_ = json.NewEncoder(w).Encode(healthResponse)
}

// Health godoc
// @Summary		Simple health check
// @Description	Check if the API service is running
// @Tags			health
// @Accept			json
// @Produce		json
// @Success		200	{object}	HealthCheckResponse
// @Router			/health [get]
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthCheckResponse{
		Status:    "healthy",
		Service:   "api",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) HealthRoutes() {
	h.r.Route("/health", func(r chi.Router) {
		r.Get("/", h.health)
		r.Get("/z", h.healthz)
	})
}

