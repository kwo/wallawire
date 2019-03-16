package user

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"wallawire/logging"
	"wallawire/model"
	"wallawire/web/auth"
)

type ChangePasswordService interface {
	ChangePassword(context.Context, model.ChangePasswordRequest) model.ChangePasswordResponse
}

func ChangePassword(userService ChangePasswordService, tokenPassword string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		logger := logging.New(ctx, "ChangePasswordHandler")
		logger.Debug().Msg("invoked")

		sessionToken := model.TokenFromContext(r.Context())
		if len(sessionToken.ID) == 0 {
			msg := "cannot retrieve user from context"
			logger.Error().Msg(msg)
			sendJsonMessage(ctx, w, http.StatusUnauthorized, msg)
			return
		}

		if r.Header.Get(hContentType) != mimeTypeJson {
			msg := "bad or missing content type"
			logger.Debug().Str(hContentType, r.Header.Get(hContentType)).Msg(msg)
			sendJsonMessage(ctx, w, http.StatusBadRequest, msg)
			return
		}

		body, errBody := ioutil.ReadAll(r.Body)
		if errBody != nil {
			msg := "cannot read request"
			logger.Debug().Err(errBody).Msg(msg)
			sendJsonMessage(ctx, w, http.StatusBadRequest, msg)
			return
		}
		defer r.Body.Close()

		var req model.ChangePasswordRequest
		if err := json.Unmarshal(body, &req); err != nil {
			msg := "bad json payload"
			logger.Debug().Err(err).Msg(msg)
			sendJsonMessage(ctx, w, http.StatusBadRequest, msg)
			return
		}

		req.UserID = sessionToken.ID
		rsp := userService.ChangePassword(ctx, req)

		if rsp.Code == http.StatusOK {

			rsp.SessionToken.Issued = time.Now().Truncate(time.Minute)
			rsp.SessionToken.Expires = rsp.SessionToken.Issued.Add(model.LoginTimeout)

			token, errToken := auth.MakeJWT(rsp.SessionToken, tokenPassword)
			if errToken != nil {
				msg := "cannot create JWT"
				logger.Error().Err(errToken).Msg(msg)
				sendJsonMessage(ctx, w, http.StatusInternalServerError, msg)
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

		sendJsonMessage(ctx, w, rsp.Code, rsp.Message)

	})
}
