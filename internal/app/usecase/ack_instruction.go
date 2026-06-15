package usecase

import (
	"context"

	appservice "github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/rs/zerolog/log"
)

const ackEmoji = "eyes"

type AckInstructionInput struct {
	ChannelID      string
	MessageID      string
	ThreadID       string
	RawInstruction string
}

type AckInstructionOutput struct{}

type AckInstruction struct {
	slack *service.Slack
}

func NewAckInstruction(slack *service.Slack) *AckInstruction {
	return &AckInstruction{
		slack: slack,
	}
}

func (u *AckInstruction) Run(ctx context.Context, input AckInstructionInput) (*AckInstructionOutput, error) {
	err := u.slack.AddReaction(ctx, input.ChannelID, input.MessageID, ackEmoji)
	return &AckInstructionOutput{}, err
}

func (u *AckInstruction) Subscribe(publisher appservice.Publisher) {
	publisher.Subscribe(chat.KindInstructionReceived, func(ctx context.Context, eev core.Event) error {
		ev, ok := eev.(*chat.InstructionReceived)
		if !ok {
			log.Error().Msg("failed to assert event")
			return nil
		}

		data := ev.Data().(*chat.InstructionReceivedData)
		input := AckInstructionInput{
			ChannelID:      data.ChannelID(),
			MessageID:      data.MessageID(),
			ThreadID:       data.ThreadID(),
			RawInstruction: data.RawInstruction(),
		}
		_, err := u.Run(ctx, input)
		return err
	})
}
