package middleware

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

var _ Middleware = (*SlackRetryIgnorere)(nil)

type SlackRetryIgnorere struct {
}

func NewSlackRetryIgnorere() *SlackRetryIgnorere {
	return &SlackRetryIgnorere{}
}

// Wrap implements Middleware.
func (s *SlackRetryIgnorere) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if r := req.Header.Get("x-slack-retry-num"); r != "" {
			log.Debug().
				Str("retry-num", r).
				Str("retry-reason", req.Header.Get("x-slack-retry-reason")).
				Msg("slack retry ignored")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "ignored")
			return
		}

		next.ServeHTTP(w, req)
	})
}
