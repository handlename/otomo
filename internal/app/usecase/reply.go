package usecase

import (
	"context"

	"github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/domain/event"
	"github.com/handlename/otomo/internal/domain/repository"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

type ReplyInput struct {
	EventData event.InstructionReceivedData
}

type ReplyOutput struct{}

type Reply struct {
	otomo           entity.Otomo
	slack           service.Messenger
	repoInstruction repository.Instruction
}

func NewReply(otomo entity.Otomo, slack service.Messenger, repoInstruction repository.Instruction) *Reply {
	return &Reply{
		otomo:           otomo,
		slack:           slack,
		repoInstruction: repoInstruction,
	}
}

func (r *Reply) Run(ctx context.Context, input ReplyInput) (*ReplyOutput, error) {
	rep, err := r.otomo.Think(ctx,
		*entity.NewContext(""),
		r.repoInstruction.NewFromInstructionReceivedData(ctx, input.EventData),
	)
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
