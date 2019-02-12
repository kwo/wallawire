package ctxutil

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger(component, fn string, ctx context.Context) zerolog.Logger {

	fields := map[string]interface{}{}

	fields["component"] = component

	if len(fn) != 0 {
		fields["fn"] = fn
	}

	if ctx != nil {

		correlationID := CorrelationIDFromContext(ctx)
		if len(correlationID) != 0 {
			fields[CorrelationIDKey] = correlationID
		}

	}

	return NewLoggerFromFields(fields)
}

func NewLoggerFromFields(fields map[string]interface{}) zerolog.Logger {
	return log.With().Fields(fields).Logger()
}
