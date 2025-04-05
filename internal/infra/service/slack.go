package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/slack-go/slack"
)

type Slack struct {
	signingSecret string
}

func NewSlack(signingSecret string) *Slack {
	return &Slack{
		signingSecret: signingSecret,
	}
}

func (s *Slack) Verify(header http.Header, body []byte) error {
	verifier, err := slack.NewSecretsVerifier(header, s.signingSecret)
	if err != nil {
		if errors.Is(err, slack.ErrMissingHeaders) {
			return failure.New(errorcode.ErrInvalidArgument, failure.Message("missing headers"))
		}

		return failure.Wrap(err)
	}

	if _, err := verifier.Write(body); err != nil {
		return failure.Wrap(err)
	}

	if err := verifier.Ensure(); err != nil {
		return failure.New(errorcode.ErrInvalidArgument, failure.Message("verify not ensured"))
	}

	return nil
}

func (s *Slack) Send(ctx context.Context, msg string) error {
	return nil
}
