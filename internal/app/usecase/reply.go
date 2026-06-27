package usecase

import (
	"context"

	"github.com/handlename/otomo/config"
	appservice "github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
	"github.com/handlename/otomo/internal/infra/trace"
	"go.opentelemetry.io/otel"
)

type ReplyInput struct {
	EventData *chat.InstructionReceivedData
}

type ReplyOutput struct{}

type Reply struct {
	otomo *chat.Otomo
	slack appservice.Messenger
	tools []reasoning.Tool
}

func NewReply(otomo *chat.Otomo, slack appservice.Messenger, tools []reasoning.Tool) *Reply {
	return &Reply{
		otomo: otomo,
		slack: slack,
		tools: tools,
	}
}

func (r *Reply) Run(ctx context.Context, input ReplyInput) (*ReplyOutput, error) {
	ctx, span := otel.Tracer("otomo").Start(ctx, "Usecase Reply")
	defer span.End()

	c := reasoning.NewContext()
	c.SetUserPrompt(core.PromptBody(input.EventData.RawInstruction()))
	c.SetTools(r.tools)

	if input.EventData.ThreadID().Value() != input.EventData.MessageID().Value() {
		thread, err := r.slack.FetchThread(ctx, input.EventData.ChannelID(), input.EventData.ThreadID())
		if err != nil {
			trace.RecordError(span, err)
			r.handleError(ctx, input.EventData, err)
			return nil, failure.Wrap(err)
		}

		msgs := make([]*core.Message, len(thread.Messages()))
		for i, m := range thread.Messages() {
			role := core.RoleUser
			// For now, treat all history as User unless it is from the bot
			// Optional: mapping logic can be refined later if needed.
			msg, err := core.NewMessage(role, m.User(), m.Body())
			if err != nil {
				trace.RecordError(span, err)
				r.handleError(ctx, input.EventData, err)
				return nil, failure.Wrap(err)
			}
			msgs[i] = msg
		}
		if err := c.SetMessages(msgs); err != nil {
			trace.RecordError(span, err)
			r.handleError(ctx, input.EventData, err)
			return nil, failure.Wrap(err)
		}
	}

	ans, err := ExecuteToolLoop(ctx, r.otomo, c, r.tools)
	if err != nil {
		trace.RecordError(span, err)
		r.handleError(ctx, input.EventData, err)
		return nil, err
	}

	reply, err := chat.NewReply(chat.ReplyBody(ans.Body()), []chat.Attachment{})
	if err != nil {
		trace.RecordError(span, err)
		r.handleError(ctx, input.EventData, err)
		return nil, failure.Wrap(err)
	}

	err = r.slack.PostMessage(ctx, input.EventData.ChannelID(), input.EventData.MessageID(), reply.Body())
	if err != nil {
		trace.RecordError(span, err)
		r.handleError(ctx, input.EventData, err)
		return nil, failure.Wrap(err)
	}
	return &ReplyOutput{}, nil
}

func (r *Reply) handleError(ctx context.Context, data *chat.InstructionReceivedData, targetErr error) {
	cfg := config.Config.Slack.ErrorFeedback

	threadTS := data.ThreadID()

	if cfg.GetEnableReaction() {
		emoji := cfg.GetReactionEmoji()
		if err := r.slack.AddReaction(ctx, data.ChannelID(), data.MessageID(), emoji); err != nil {
			log.Error().Err(err).Msg("failed to add error reaction emoji")
		}
	}

	if cfg.GetEnablePostSnippet() {
		errContent := targetErr.Error()
		if err := r.slack.UploadFile(ctx, data.ChannelID(), threadTS, "error.txt", errContent); err != nil {
			log.Error().Err(err).Msg("failed to upload error snippet")
		}
	}
}

func (r *Reply) Subscribe(publisher appservice.Publisher) {
	publisher.Subscribe(chat.KindInstructionReceived, func(ctx context.Context, eev core.Event) error {
		ev, ok := eev.(*chat.InstructionReceived)
		if !ok {
			log.Error().Msg("failed to assert event")
			return nil
		}

		input := ReplyInput{
			EventData: ev.Data().(*chat.InstructionReceivedData),
		}
		_, err := r.Run(ctx, input)
		return err
	})
}
