package ctxutil

import (
	"context"

	"wallawire/model"
)

const (
	CorrelationIDKey = "correlationID"
	UserKey          = "user"
)

func CorrelationIDFromContext(ctx context.Context) string {
	if value := ctx.Value(CorrelationIDKey); value != nil {
		if uuid, ok := value.(string); ok {
			return uuid
		}
	}
	return ""
}

func TokenFromContext(ctx context.Context) model.SessionToken {
	if value := ctx.Value(UserKey); value != nil {
		if user, ok := value.(model.SessionToken); ok {
			return user
		}
	}
	return model.SessionToken{}
}
