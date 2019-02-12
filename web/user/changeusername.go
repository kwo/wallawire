package user

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"wallawire/ctxutil"
	"wallawire/model"
	"wallawire/web/auth"
)

type ChangeUsernameService interface {
	ChangeUsername(context.Context, model.ChangeUsernameRequest) model.ChangeUsernameResponse
}

func ChangeUsername(userService ChangeUsernameService, tokenPassword string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		logger := ctxutil.NewLogger("ChangeUsername", "", ctx)
		logger.Debug().Msg("invoked")

		sessionToken := ctxutil.TokenFromContext(r.Context())
		if len(sessionToken.ID) == 0 {
			msg := "cannot retrieve user from context"
			logger.Error().Msg(msg)
			sendJsonMessage(w, http.StatusUnauthorized, msg)
			return
		}

		if r.Header.Get(hContentType) != mimeTypeJson {
			msg := "bad or missing content type"
			logger.Debug().Str(hContentType, r.Header.Get(hContentType)).Msg(msg)
			sendJsonMessage(w, http.StatusBadRequest, msg)
			return
		}

		body, errBody := ioutil.ReadAll(r.Body)
		if errBody != nil {
			msg := "cannot read request"
			logger.Debug().Err(errBody).Msg(msg)
			sendJsonMessage(w, http.StatusBadRequest, msg)
			return
		}
		defer r.Body.Close()

		var req model.ChangeUsernameRequest
		if err := json.Unmarshal(body, &req); err != nil {
			msg := "bad json payload"
			logger.Debug().Err(err).Msg(msg)
			sendJsonMessage(w, http.StatusBadRequest, msg)
			return
		}

		req.UserID = sessionToken.ID
		rsp := userService.ChangeUsername(ctx, req)

		if rsp.Code == http.StatusOK {

			rsp.SessionToken.Issued = time.Now().Truncate(time.Minute)
			rsp.SessionToken.Expires = rsp.SessionToken.Issued.Add(model.LoginTimeout)

			token, errToken := auth.MakeJWT(rsp.SessionToken, tokenPassword)
			if errToken != nil {
				msg := "cannot create JWT"
				logger.Error().Err(errToken).Msg(msg)
				sendJsonMessage(w, http.StatusInternalServerError, msg)
				return
			}

			cookie := &http.Cookie{
				Name:    auth.CookieName,
				Value:   token,
				Path:    "/",
				Expires: rsp.SessionToken.Expires,
				Secure:  true,
			}
			http.SetCookie(w, cookie)

		}

		sendJsonMessage(w, rsp.Code, rsp.Message)

	})
}
