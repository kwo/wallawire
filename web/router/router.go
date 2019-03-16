package router

import (
	"compress/flate"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"wallawire/logging"
	"wallawire/web/accesslog"
	wallaware "wallawire/web/middleware"
)

const (
	cacheNoCache    = "no-cache"
	hAcceptEncoding = "Accept-Encoding"
	hCacheControl   = "Cache-Control"
	hContentLength  = "Content-Length"
	hVary           = "Vary"
)

type IdGenerator interface {
	NewID() string
}

type Options struct {
	Authenticator   []func(http.Handler) http.Handler
	AuthorizerUsers func(http.Handler) http.Handler
	ChangePassword  http.HandlerFunc
	ChangeUsername  http.HandlerFunc
	ChangeProfile   http.HandlerFunc
	IdGenerator     IdGenerator // because composing middleware in router
	Login           http.HandlerFunc
	Logout          http.HandlerFunc
	Notifier        http.HandlerFunc
	Static          http.HandlerFunc
	Status          http.HandlerFunc
	Whoami          http.HandlerFunc
}

func Router(opts Options) (http.Handler, error) {

	// logger := log.With().Str("component", "router").Logger()

	rMain := chi.NewRouter()

	// global middleware
	rMain.Use(wallaware.CorrelationID(opts.IdGenerator))
	rMain.Use(accesslog.AccessHandler(accessLogger()))
	rMain.Use(noCache)
	rMain.Use(middleware.SetHeader(hVary, hAcceptEncoding))
	// TODO security middleware
	// TODO rMain.Use(compression)

	// api router
	rMain.Route("/api/*", func(rApi chi.Router) {
		// group for routes requiring authentication
		rApi.Group(func(rAuth chi.Router) {
			rAuth.Use(opts.Authenticator...)
			rAuth.Use(opts.AuthorizerUsers)
			rAuth.Group(func(rTimeout chi.Router) {
				rTimeout.Use(middleware.Timeout(time.Second * 60))
				rTimeout.Post("/changepassword", opts.ChangePassword)
				rTimeout.Post("/changeusername", opts.ChangeUsername)
				rTimeout.Post("/changeprofile", opts.ChangeProfile)
				rTimeout.Get("/whoami", opts.Whoami)
			})
			rAuth.Get("/inbox", opts.Notifier) // no timeout
		})

		rApi.Post("/login", opts.Login)
		rApi.Get("/logout", opts.Logout)
		rApi.Post("/logout", opts.Logout)
		rApi.Get("/status", opts.Status)
		rApi.NotFound(sendMessageHandler(http.StatusNotFound))
		rApi.MethodNotAllowed(sendMessageHandler(http.StatusMethodNotAllowed))
	})

	// static router
	rMain.Route("/*", func(rStatic chi.Router) {
		rStatic.Get("/*", opts.Static)
		rStatic.MethodNotAllowed(sendMessageHandler(http.StatusMethodNotAllowed))
	})

	return rMain, nil

}

func accessLogger() func(r *http.Request, status, size int, duration time.Duration) {
	return func(r *http.Request, status, size int, duration time.Duration) {
		ctx := r.Context()
		logger := logging.New(ctx, "accesslog")
		logger.Debug().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration.Truncate(time.Microsecond)).
			Msg("request")
	}
}

func compression(next http.Handler) http.Handler {
	additionalCompressionMimeTypes := []string{
		"text/html",
		"text/css",
		"text/plain",
		"text/javascript",
		"text/xml",
		"application/javascript",
		"application/x-javascript",
		"application/json",
		"application/atom+xml",
		"application/rss+xml",
		"application/xml",
	}
	return middleware.Compress(flate.BestCompression, additionalCompressionMimeTypes...)(next)
}

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(hCacheControl, cacheNoCache)
		h.ServeHTTP(w, r)
	})
}

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
