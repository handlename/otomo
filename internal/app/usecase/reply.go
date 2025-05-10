package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/domain/event"
	vo "github.com/handlename/otomo/internal/domain/valueobject"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type ReplyInput struct {
	EventData event.InstructionReceivedData
}

type ReplyOutput struct{}

type Reply struct {
	otomo entity.Otomo
	slack service.Messenger
}

func NewReply(otomo entity.Otomo, slack service.Messenger) *Reply {
	return &Reply{
		otomo: otomo,
		slack: slack,
	}
}

func (r *Reply) Run(ctx context.Context, input ReplyInput) (*ReplyOutput, error) {
	c := entity.NewContext()
	c.SetUserPrompt(input.EventData.RawInstruction)

	thread, err := r.slack.FetchThread(ctx, input.EventData.ChannelID, input.EventData.ThreadID)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	thread.Messages()
	log.Debug().Strs("messages", lo.Map(thread.Messages(), func(m entity.ThreadMessage, _ int) string {
		return m.String()
	})).Msg("fetched thread")
	c.SetThread(thread)

	rep, err := r.otomo.Think(ctx, c)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	err = r.slack.PostMessage(ctx, input.EventData.ChannelID, input.EventData.MessageID, rep.Body())
	return &ReplyOutput{}, err
}

func (r *Reply) Subscribe(publisher event.Publisher) {
	publisher.Subscribe(event.KindInstructionReceived, func(ctx context.Context, eev event.Event) error {
		ev, ok := eev.(*event.InstructionReceived)
		if !ok {
			log.Error().Msg("failed to assert event")
			return nil
		}

		input := ReplyInput{
			EventData: ev.Data().(event.InstructionReceivedData),
		}
		_, err := r.Run(ctx, input)
		return err
	})
}

func (r *Reply) buildPrompt(raw string) vo.Prompt {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, fmt.Sprintf("<%s>", config.Config.Slack.BotUserID))
	return vo.NewPlainPrompt(raw)
}
