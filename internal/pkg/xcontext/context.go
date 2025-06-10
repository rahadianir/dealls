package xcontext

import "context"

type contextKey string

const RequestIDKey contextKey = "request.id"
const IPKey contextKey = "ip"
const UserIDKey contextKey = "user.id"

func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return ""
	}

	return userID
}
