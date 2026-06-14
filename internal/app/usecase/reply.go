package usecase

import (
	"context"

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
		return nil, failure.Wrap(err)
	}

	err = r.slack.PostMessage(ctx, input.EventData.ChannelID, input.EventData.MessageID, rep.Body())
	return &ReplyOutput{}, err
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
