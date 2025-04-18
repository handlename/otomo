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
	client        *slack.Client
}

func NewSlack(token string, signingSecret string) *Slack {
	return &Slack{
		signingSecret: signingSecret,
		client:        slack.New(token),
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

func (s *Slack) PostMessage(ctx context.Context, channelID, messageID, msg string) error {
	_, _, err := s.client.PostMessage(
		channelID,
		slack.MsgOptionTS(messageID),
		slack.MsgOptionText(msg, false),
	)
	return err
}

func (s *Slack) AddReaction(ctx context.Context, channelID, messageID string, emoji string) error {
	return s.client.AddReaction(emoji, slack.ItemRef{
		Channel:   channelID,
		Timestamp: messageID,
	})
}
