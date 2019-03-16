package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"

	"wallawire/logging"
	"wallawire/model"
)

func NewAuthenticator(password string) []func(http.Handler) http.Handler {
	tokenAuthenticator := jwtauth.New(jwt.SigningMethodHS256.Name, []byte(password), nil)
	return []func(next http.Handler) http.Handler{
		jwtauth.Verifier(tokenAuthenticator),
		jwtauth.Authenticator,
		tokenToUser(),
	}
}

// tokenToUser converts the JWT token into a session user object and add it to the request context.
// The user object has all fields including id and roles but not the password_hash field filled.
func tokenToUser() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			logger := logging.New(ctx, "auth", "Authenticator")

			_, claims, err := jwtauth.FromContext(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("Cannot retrieve jwt token from context")
				sendMessage(w, http.StatusInternalServerError)
				return
			}

			user := &model.SessionToken{}

			if value, ok := claims.Get("sessionid"); ok {
				if id, ok := value.(string); ok {
					user.SessionID = id
				}
			}

			if value, ok := claims.Get("id"); ok {
				if id, ok := value.(string); ok {
					user.ID = id
				}
			}

			if value, ok := claims.Get("username"); ok {
				if username, ok := value.(string); ok {
					user.Username = username
				}
			}

			if value, ok := claims.Get("name"); ok {
				if name, ok := value.(string); ok {
					user.Name = name
				}
			}

			if value, ok := claims.Get("roles"); ok {
				if roles, ok := value.(string); ok {
					user.Roles = strings.Split(roles, ",")
				}
			}

			if value, ok := claims.Get("iat"); ok {
				if iat, ok := value.(int64); ok {
					user.Issued = time.Unix(iat, 0)
				}
			}

			if value, ok := claims.Get("exp"); ok {
				if exp, ok := value.(int64); ok {
					user.Expires = time.Unix(exp, 0)
				}
			}

			ctx = context.WithValue(r.Context(), model.UserKey, *user)
			logger.Info().Str("UserID", user.ID).Str("SessionID", user.SessionID).Msg("authenticated")

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
