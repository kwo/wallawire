package router

import (
	"compress/flate"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/hlog"

	"wallawire/ctxutil"
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
	Static          http.HandlerFunc
	Status          http.HandlerFunc
	Whoami          http.HandlerFunc
}

func Router(opts Options) (http.Handler, error) {

	// logger := log.With().Str("component", "router").Logger()

	h := chi.NewRouter()

	// global middleware
	h.Use(wallaware.CorrelationID(opts.IdGenerator))
	h.Use(hlog.AccessHandler(accessLogger()))
	h.Use(middleware.Timeout(time.Second * 60))
	// h.Use(middleware.CloseNotify) built-in to go1.8
	// TODO security middleware
	h.Use(noCache)
	// TODO h.Use(compression)
	h.Use(middleware.SetHeader(hVary, hAcceptEncoding))

	// api router
	h.Route("/api/*", func(r chi.Router) {
		// group for routes requiring authentication
		r.Group(func(p chi.Router) {
			p.Use(opts.Authenticator...)
			p.Use(opts.AuthorizerUsers)
			p.Post("/changepassword", opts.ChangePassword)
			p.Post("/changeusername", opts.ChangeUsername)
			p.Post("/changeprofile", opts.ChangeProfile)
			p.Get("/whoami", opts.Whoami)
		})

		r.Post("/login", opts.Login)
		r.Get("/logout", opts.Logout)
		r.Post("/logout", opts.Logout)
		r.Get("/status", opts.Status)
		r.NotFound(sendMessageHandler(http.StatusNotFound))
		r.MethodNotAllowed(sendMessageHandler(http.StatusMethodNotAllowed))
	})

	// static router
	h.Route("/*", func(r chi.Router) {
		r.Get("/*", opts.Static)
		r.MethodNotAllowed(sendMessageHandler(http.StatusMethodNotAllowed))
	})

	return h, nil

}

func accessLogger() func(r *http.Request, status, size int, duration time.Duration) {
	return func(r *http.Request, status, size int, duration time.Duration) {
		ctx := r.Context()
		logger := ctxutil.NewLogger("accesslog", "", ctx)
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
