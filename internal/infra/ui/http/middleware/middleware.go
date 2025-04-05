package middleware

import (
	"net/http"

	"github.com/samber/lo"
)

type Middleware interface {
	Wrap(next http.Handler) http.Handler
}

func Wrap(handler func(w http.ResponseWriter, r *http.Request), middlewares ...Middleware) http.Handler {
	var h http.Handler = http.HandlerFunc(handler)
	lo.ForEach(middlewares, func(m Middleware, _ int) {
		h = m.Wrap(h)
	})

	return h
}
