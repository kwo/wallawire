package logging

import (
	"context"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"wallawire/model"
)

func New(ctx context.Context, components ...string) *zerolog.Logger {

	fields := map[string]interface{}{}

	fields["component"] = strings.Join(components, ":")

	if ctx != nil {
		correlationID := model.CorrelationIDFromContext(ctx)
		if len(correlationID) != 0 {
			fields[model.CorrelationIDKey] = correlationID
		}
	}

	return NewLoggerFromFields(fields)
}

func NewLoggerFromFields(fields map[string]interface{}) *zerolog.Logger {
	logger := log.With().Fields(fields).Logger()
	return &logger
}
