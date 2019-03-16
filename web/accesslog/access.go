package accesslog

import (
	"net/http"
	"time"
)

// AccessHandler returns a handler that call f after each request.
func AccessHandler(f func(r *http.Request, status, size int, duration time.Duration)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lw := WrapWriter(w)
			next.ServeHTTP(lw, r)
			f(r, lw.Status(), lw.BytesWritten(), time.Since(start))
		})
	}
}
