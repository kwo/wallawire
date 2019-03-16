package model

import (
	"context"
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

func TokenFromContext(ctx context.Context) SessionToken {
	if value := ctx.Value(UserKey); value != nil {
		if user, ok := value.(SessionToken); ok {
			return user
		}
	}
	return SessionToken{}
}
