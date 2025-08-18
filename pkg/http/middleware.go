package http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	contextKeyUserGUID = key("context-user-guid")
)

type key string

func userGUIDToContext(ctx context.Context, userGUID uuid.UUID) context.Context {
	return context.WithValue(ctx, contextKeyUserGUID, userGUID)
}

func UserGUIDFromContext(ctx context.Context) uuid.UUID {
	val := ctx.Value(contextKeyUserGUID)

	if userGUID, ok := val.(uuid.UUID); ok {
		return userGUID
	}

	return uuid.Nil
}

type authMiddleware struct{}

type AuthMiddleware interface {
	UserAuthorizationMiddleware(next http.Handler) http.Handler
}

func NewAuthMiddleware() AuthMiddleware {
	return &authMiddleware{}
}

func (m *authMiddleware) UserAuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userGUIDHeader := r.Header.Get("X-User")

		if userGUIDHeader == "" {
			next.ServeHTTP(w, r)

			return
		}

		userGUID, err := uuid.Parse(userGUIDHeader)
		if err != nil {
			next.ServeHTTP(w, r)

			return
		}

		ctx = userGUIDToContext(ctx, userGUID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
