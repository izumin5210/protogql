package loader

import (
	user_pb "apis/go/user"
	"context"
	"net/http"
)

type Loaders struct {
	UserClient user_pb.UserServiceClient
}

type loaderCtxKey struct{}

func (l *Loaders) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loaderCtxKey{}, l)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loaderCtxKey{}).(*Loaders)
}
