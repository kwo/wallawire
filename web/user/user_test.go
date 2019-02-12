package user_test

import (
	"context"
	"flag"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/rs/zerolog"

	"wallawire/model"
	"wallawire/web/auth"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Verbose() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
	os.Exit(m.Run())
}

const (
	ignoreValue    = "XXX"
	hContentLength = "Content-Length"
	hContentType   = "Content-Type"
	hCookie        = "Cookie"
	hDate          = "Date"
	hSetCookie     = "Set-Cookie"
	mimeTypeJson   = "application/json"
	mimeTypeText   = "text/plain; charset=utf-8"
	testPassword   = "secret"
)

func sendMessage(w http.ResponseWriter, statusCode int) {
	msg := []byte(http.StatusText(statusCode) + "\n")
	w.Header().Set(hContentLength, strconv.Itoa(len(msg)))
	w.WriteHeader(statusCode)
	w.Write(msg)
}

func getCookieString(user *model.SessionToken, password string) string {
	r, err := auth.MakeJWT(user, password)
	if err != nil {
		panic(err)
	}
	c := &http.Cookie{
		Name:    auth.CookieName,
		Value:   r,
		Expires: user.Expires,
		Path:    "/",
		Secure:  true,
	}
	return c.String()
}

type UserServiceMock struct {
	ChangePasswordResponse model.ChangePasswordResponse
	ChangeUsernameResponse model.ChangeUsernameResponse
	ChangeProfileResponse  model.ChangeProfileResponse
}

func (z *UserServiceMock) ChangePassword(ctx context.Context, req model.ChangePasswordRequest) model.ChangePasswordResponse {
	return z.ChangePasswordResponse
}

func (z *UserServiceMock) ChangeUsername(ctx context.Context, req model.ChangeUsernameRequest) model.ChangeUsernameResponse {
	return z.ChangeUsernameResponse
}

func (z *UserServiceMock) ChangeProfile(ctx context.Context, req model.ChangeProfileRequest) model.ChangeProfileResponse {
	return z.ChangeProfileResponse
}
