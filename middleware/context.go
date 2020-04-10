package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const requestUserID contextKey = "request_user_id"

func setRequestCtx(ctx context.Context, r *http.Request, key, value interface{}) *http.Request {
	return r.WithContext(context.WithValue(ctx, key, value))
}

func contextValue(ctx context.Context, key interface{}) interface{} {
	return ctx.Value(key)
}

// CtxReqUserID retrieves the authenticated userID
func CtxReqUserID(ctx context.Context) string {
	if v := contextValue(ctx, requestUserID); v != nil {
		return v.(string)
	}
	return ""
}
