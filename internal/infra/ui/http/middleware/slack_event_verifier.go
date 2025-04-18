package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/handlename/otomo/internal/errorcode"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

var _ Middleware = (*SlackEventVerifier)(nil)

type SlackEventVerifier struct {
	slack *service.Slack
}

func NewSlackEventVerifier(slack *service.Slack) *SlackEventVerifier {
	return &SlackEventVerifier{
		slack: slack,
	}
}

// Wrap implements Middleware.
func (s *SlackEventVerifier) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Read the body but save it for future use
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read body")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal Server Error")
			return
		}

		// Close the original body as we've read it
		req.Body.Close()

		if err := s.slack.Verify(req.Header, body); err != nil {
			if failure.Is(err, errorcode.ErrInvalidArgument) {
				log.Debug().Err(err).Msg("invalid signing")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "invilid signing")
				return
			}

			log.Error().Err(err).Msg("failed to verify slack event")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal Server Error")
			return
		}

		log.Debug().Msg("slack event verification succeeded")

		// Restore the body so it can be read in next handlers
		req.Body = io.NopCloser(bytes.NewBuffer(body))

		next.ServeHTTP(w, req)
	})
}
