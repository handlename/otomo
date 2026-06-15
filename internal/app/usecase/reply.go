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
)

type ReplyInput struct {
	EventData chat.InstructionReceivedData
}

type ReplyOutput struct{}

type Reply struct {
	otomo *chat.Otomo
	slack appservice.Messenger
}

func NewReply(otomo *chat.Otomo, slack appservice.Messenger) *Reply {
	return &Reply{
		otomo: otomo,
		slack: slack,
	}
}

func (r *Reply) Run(ctx context.Context, input ReplyInput) (*ReplyOutput, error) {
	c := reasoning.NewContext()
	c.SetUserPrompt(input.EventData.RawInstruction)

	if input.EventData.ThreadID != "" {
		thread, err := r.slack.FetchThread(ctx, input.EventData.ChannelID, input.EventData.ThreadID)
		if err != nil {
			r.handleError(ctx, input.EventData, err)
			return nil, failure.Wrap(err)
		}

		msgs := make([]core.Message, len(thread.Messages()))
		for i, m := range thread.Messages() {
			role := core.RoleUser
			// For now, treat all history as User unless it is from the bot
			// Optional: mapping logic can be refined later if needed.
			msgs[i] = core.Message{
				Role: role,
				User: m.User(),
				Body: m.Body(),
			}
		}
		c.SetMessages(msgs)
	}

	rep, err := r.otomo.Think(ctx, c)
	if err != nil {
		r.handleError(ctx, input.EventData, err)
		return nil, failure.Wrap(err)
	}

	err = r.slack.PostMessage(ctx, input.EventData.ChannelID, input.EventData.MessageID, rep.Body())
	if err != nil {
		r.handleError(ctx, input.EventData, err)
		return nil, failure.Wrap(err)
	}
	return &ReplyOutput{}, nil
}

func (r *Reply) handleError(ctx context.Context, data chat.InstructionReceivedData, targetErr error) {
	cfg := config.Config.Slack.ErrorFeedback

	var threadTS string
	if data.ThreadID != "" {
		threadTS = data.ThreadID
	} else {
		threadTS = data.MessageID
	}

	if cfg.GetEnableReaction() {
		emoji := cfg.GetReactionEmoji()
		if err := r.slack.AddReaction(ctx, data.ChannelID, data.MessageID, emoji); err != nil {
			log.Error().Err(err).Msg("failed to add error reaction emoji")
		}
	}

	if cfg.GetEnablePostSnippet() {
		errContent := targetErr.Error()
		if err := r.slack.UploadFile(ctx, data.ChannelID, threadTS, "error.txt", errContent); err != nil {
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
			EventData: ev.Data().(chat.InstructionReceivedData),
		}
		_, err := r.Run(ctx, input)
		return err
	})
}
