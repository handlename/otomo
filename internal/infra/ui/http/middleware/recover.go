package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

var _ Middleware = (*Recover)(nil)

type Recover struct{}

func NewRecover() *Recover {
	return &Recover{}
}

// Wrap implements Middleware.
func (a *Recover) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				e, ok := err.(error)
				if !ok {
					e = failure.New("panic in middleware")
				}

				log.Err(e).
					Any("code", failure.CodeOf(e)).
					Str("message", failure.MessageOf(e).String()).
					Str("stack", strings.Join(
						lo.Map(failure.CallStackOf(e).Frames(), func(frame failure.Frame, _ int) string {
							return fmt.Sprintf("%s.%s %s:%d", frame.PkgPath(), frame.Func(), frame.File(), frame.Line())
						}),
						"\n"),
					).
					Msg("error recovered")

				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Internal server error")
			}
		}()

		next.ServeHTTP(w, req)
	})
}
