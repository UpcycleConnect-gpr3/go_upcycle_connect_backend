package auth_middleware

import (
	"context"
	"go-upcycle_connect-backend/utils/jwt"
	"go-upcycle_connect-backend/utils/log"
	"net/http"
)

type contextKey string

const userIdKey contextKey = "userId"

func IsAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := jwt.Auth(w, r)
		if userId == "" {
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func GetUserId(ctx context.Context) string {

	userId, ok := ctx.Value(userIdKey).(string)

	if !ok {
		log.Info("No user id from context")
		return ""
	}

	return userId
}
