package middleware

import (
	"context"
	"errors"
	"net/http"
)

type registryKey struct{}

var _ Middleware = (*Registry[struct{}])(nil)

type Registry[R any] struct {
	reg R
}

func NewRegistry[R any](reg R) *Registry[R] {
	return &Registry[R]{
		reg: reg,
	}
}

// Wrap implements Middleware.
func (r *Registry[R]) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), registryKey{}, r.reg)
		req = req.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func GetRegistry[R any](ctx context.Context) (R, error) {
	reg, ok := ctx.Value(registryKey{}).(R)
	if !ok {
		var zero R
		return zero, errors.New("invalid registry")
	}

	return reg, nil
}
