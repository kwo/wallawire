package auth

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"wallawire/ctxutil"
	"wallawire/model"
)

type UserService interface {
	Login(context.Context, model.LoginRequest) model.LoginResponse
}

func Login(userService UserService, tokenPassword string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := ctxutil.NewLogger("auth", "LoginHandler", ctx)
		logger.Debug().Msg("invoked")

		if r.Header.Get(hContentType) != mimeTypeJson {
			msg := "bad or missing content type"
			logger.Debug().Str(hContentType, r.Header.Get(hContentType)).Msg(msg)
			sendMessageText(w, http.StatusBadRequest, msg)
			return
		}

		body, errBody := ioutil.ReadAll(r.Body)
		if errBody != nil {
			msg := "cannot read request body"
			logger.Debug().Err(errBody).Msg(msg)
			sendMessageText(w, http.StatusBadRequest, msg)
			return
		}
		defer r.Body.Close()

		var req model.LoginRequest
		if err := json.Unmarshal(body, &req); err != nil {
			msg := "cannot unmarshal json"
			logger.Debug().Err(err).Msg(msg)
			sendMessageText(w, http.StatusBadRequest, msg)
			return
		}

		rsp := userService.Login(ctx, req)
		if rsp.Code != http.StatusOK {
			sendMessageText(w, rsp.Code, rsp.Message)
			return
		}

		// OK

		token, errToken := MakeJWT(rsp.SessionToken, tokenPassword)
		if errToken != nil {
			msg := "cannot create JWT"
			logger.Error().Err(errToken).Msg(msg)
			sendMessageText(w, http.StatusInternalServerError, msg)
			return
		}

		cookie := &http.Cookie{
			Name:    CookieName,
			Value:   token,
			Path:    "/",
			Expires: rsp.SessionToken.Expires,
			Secure:  true,
		}
		http.SetCookie(w, cookie)

		sendMessage(w, http.StatusOK)

	})
}

func MakeJWT(user *model.SessionToken, tokenPassword string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sessionid": user.SessionID,
		"id":        user.ID,
		"username":  user.Username,
		"name":      user.Name,
		"roles":     strings.Join(user.Roles, ","),
		"iat":       user.Issued.Unix(),
		"exp":       user.Expires.Unix(),
	})

	return token.SignedString([]byte(tokenPassword))

}
