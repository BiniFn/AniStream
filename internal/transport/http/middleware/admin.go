package middleware

import (
	"net/http"

	"github.com/coeeter/aniways/internal/utils"
)

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		adminKey := authHeader[7:]
		if !utils.ValidateAdminKey(adminKey) {
			http.Error(w, "Invalid admin key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

