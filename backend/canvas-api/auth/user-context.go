package auth

import (
	"context"
)

// Context key type to avoid context key collisions
type contextKey string

const userIDKey contextKey = "userID"

// SetUserIDInContext adds the userID to the context.
func SetUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext retrieves the userID from the context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
