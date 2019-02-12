package auth

import (
	"net/http"
	"time"

	"wallawire/ctxutil"
)

func Logout() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := ctxutil.NewLogger("auth", "LogoutHandler", ctx)
		user := ctxutil.TokenFromContext(ctx)
		logger.Info().Str("UserID", user.ID).Str("SessionID", user.SessionID).Msg("logout")
		cookie := &http.Cookie{
			Name:    CookieName,
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
			Secure:  true,
		}
		http.SetCookie(w, cookie)
		sendMessage(w, http.StatusOK)
	})
}
