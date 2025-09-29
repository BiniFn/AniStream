package handlers

import (
	"net/http"

	"github.com/coeeter/aniways/internal/template"
)

func (h *Handler) serveAdminPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(template.AdminPanelTemplate))
}