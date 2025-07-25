package middleware

import (
	"net/http"

	"github.com/coeeter/aniways/internal/ctxutil"
	"github.com/coeeter/aniways/internal/service/users"
)

func RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isUserAuthenticated(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isUserAuthenticated(r *http.Request) bool {
	user, ok := ctxutil.Get[users.User](r.Context())
	if !ok {
		return false
	}
	return user.ID != ""
}
