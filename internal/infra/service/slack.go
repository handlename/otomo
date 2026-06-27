package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/handlename/otomo/internal/infra/trace"
	"go.opentelemetry.io/otel"
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

func (s *Slack) PostMessage(ctx context.Context, channelID core.ChannelID, messageID core.MessageID, msg chat.ReplyBody) error {
	ctx, span := otel.Tracer("otomo").Start(ctx, "Slack PostMessage")
	defer span.End()

	block := slack.NewMarkdownBlock("", string(msg))
	_, _, err := s.client.PostMessageContext(
		ctx,
		channelID.Value(),
		slack.MsgOptionTS(messageID.Value()),
		slack.MsgOptionText(string(msg), false), // fallback text
		slack.MsgOptionBlocks(block),
	)
	if err != nil {
		err = failure.Wrap(err)
		trace.RecordError(span, err)
		return err
	}
	return nil
}

func (s *Slack) AddReaction(ctx context.Context, channelID core.ChannelID, messageID core.MessageID, emoji string) error {
	ctx, span := otel.Tracer("otomo").Start(ctx, "Slack AddReaction")
	defer span.End()

	err := s.client.AddReactionContext(ctx, emoji, slack.ItemRef{
		Channel:   channelID.Value(),
		Timestamp: messageID.Value(),
	})
	if err != nil {
		err = failure.Wrap(err)
		trace.RecordError(span, err)
		return err
	}
	return nil
}

// FetchThread implements service.Messenger.
func (s *Slack) FetchThread(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID) (*chat.Thread, error) {
	ctx, span := otel.Tracer("otomo").Start(ctx, "Slack FetchThread")
	defer span.End()

	t, err := chat.NewThread(threadID)
	if err != nil {
		err = failure.Wrap(err)
		trace.RecordError(span, err)
		return nil, err
	}
	more := true
	next := ""

	for more {
		var msgs []slack.Message
		msgs, more, next, err = s.fetchThread(ctx, channelID.Value(), threadID.Value(), next)
		if err != nil {
			err = failure.Wrap(err)
			trace.RecordError(span, err)
			return nil, err
		}

		threadMsgs := make([]*chat.ThreadMessage, 0, len(msgs))
		for _, m := range msgs {
			body := m.Text
			body = strings.TrimPrefix(body, fmt.Sprintf("<%s>", config.Config.Slack.BotUserID))
			body = strings.TrimSpace(body)

			user := m.User
			if user == "" {
				user = m.BotID
			}
			if user == "" {
				user = "unknown"
			}

			u, err := core.NewUserID(user)
			if err != nil {
				log.Warn().Err(err).Msg("failed to create UserID from slack user")
				continue
			}

			tmID, err := chat.NewThreadMessageID(m.Timestamp)
			if err != nil {
				log.Warn().Err(err).Msg("failed to create ThreadMessageID from slack message timestamp")
				continue
			}
			tm, err := chat.NewThreadMessage(tmID, u, core.MessageBody(body))
			if err != nil {
				log.Warn().Err(err).Msg("failed to create thread message from slack message")
				continue
			}
			threadMsgs = append(threadMsgs, tm)
		}
		t.AddMessages(threadMsgs...)
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

// UploadFile implements service.Messenger.
func (s *Slack) UploadFile(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID, filename, content string) error {
	ctx, span := otel.Tracer("otomo").Start(ctx, "Slack UploadFile")
	defer span.End()

	_, err := s.client.UploadFileV2Context(ctx, slack.UploadFileV2Parameters{
		Channel:         channelID.Value(),
		ThreadTimestamp: threadID.Value(),
		Filename:        filename,
		Content:         content,
		FileSize:        len(content),
	})
	if err != nil {
		err = failure.Wrap(err)
		trace.RecordError(span, err)
		return err
	}
	return nil
}
