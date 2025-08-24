package middleware

import (
	"net/http"

	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/utils"
)

// RequireUser middleware ensures that the user is authenticated.
func RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isUserAuthenticated(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetUser retrieves the authenticated user from the request context.
func GetUser(r *http.Request) *models.UserResponse {
	user, _ := utils.CtxValue[models.UserResponse](r.Context())
	return &user
}

func isUserAuthenticated(r *http.Request) bool {
	user, ok := utils.CtxValue[models.UserResponse](r.Context())
	if !ok {
		return false
	}
	return user.ID != ""
}
