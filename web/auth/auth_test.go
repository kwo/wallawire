package auth_test

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

func sendMessageHandler(statusCode int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sendMessage(w, statusCode)
	})
}

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
	LoginResponse model.LoginResponse
}

func (z *UserServiceMock) Login(ctx context.Context, req model.LoginRequest) model.LoginResponse {
	return z.LoginResponse
}
