package middleware

import (
	"context"
	"net/http"

	"wallawire/ctxutil"
)

type IdGenerator interface {
	NewID() string
}

func CorrelationID(idgen IdGenerator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			correlationID := ctxutil.CorrelationIDFromContext(ctx)
			if len(correlationID) == 0 {
				correlationID = idgen.NewID()
				ctx = context.WithValue(ctx, ctxutil.CorrelationIDKey, correlationID)
			}

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

		})
	}
}
