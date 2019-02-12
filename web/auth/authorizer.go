package auth

import (
	"net/http"

	"wallawire/ctxutil"
)

// NewAuthorizer returns middleware to forbid users without ALL the specified roles
func NewAuthorizer(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			logger := ctxutil.NewLogger("auth", "Authorizer", ctx)
			user := ctxutil.TokenFromContext(r.Context()) // guarenteed to always be not nil

			for _, role := range roles {
				if !user.HasRole(role) {
					logger.Info().Msg("forbidden")
					sendMessage(w, http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)

		})
	}
}
