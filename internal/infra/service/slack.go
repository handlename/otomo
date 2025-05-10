package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/samber/lo"
	"github.com/slack-go/slack"
)

var _ service.Messenger = (*Slack)(nil)

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

// FetchThread implements service.Messenger.
func (s *Slack) FetchThread(ctx context.Context, channelID string, threadID string) (entity.Thread, error) {
	t := entity.NewThread(entity.ThreadID(threadID))
	more := true
	next := ""

	for more {
		msgs := []slack.Message{}
		var err error
		msgs, more, next, err = s.fetchThread(ctx, channelID, threadID, next)
		if err != nil {
			return nil, failure.Wrap(err)
		}

		t.AddMessages(lo.Map(msgs, func(m slack.Message, _ int) entity.ThreadMessage {
			return entity.NewThreadMessage(entity.ThreadMessageID(m.Timestamp), m.Text)
		})...)
	}

	return t, nil
}

func (s *Slack) fetchThread(ctx context.Context, channelID, threadID, cursor string) ([]slack.Message, bool, string, error) {
	params := &slack.GetConversationRepliesParameters{
		ChannelID: channelID,
		Timestamp: threadID,
	}
	if cursor != "" {
		params.Cursor = cursor
	}

	msgs, more, next, err := s.client.GetConversationRepliesContext(ctx, params)
	if err != nil {
		return nil, false, "", failure.Wrap(err,
			failure.Message("failed to fetch conversation replies"),
			failure.Context{
				"channelID": channelID,
				"threadID":  threadID,
				"cursor":    cursor,
			})
	}

	return msgs, more, next, nil
}
