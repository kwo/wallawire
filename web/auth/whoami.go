package auth

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wallawire/ctxutil"
)

func Whoami() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		logger := ctxutil.NewLogger("auth", "WhoamiHandler", ctx)

		user := ctxutil.TokenFromContext(ctx)
		if len(user.ID) == 0 {
			logger.Warn().Msg("cannot retrieve user from context")
			sendMessage(w, http.StatusUnauthorized)
			return
		}

		payload, errData := json.Marshal(user)
		if errData != nil {
			sendMessage(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set(hContentType, mimeTypeJson)
		w.Header().Set(hContentLength, strconv.Itoa(len(payload)))
		w.Write(payload)

	})
}
