package usecase

import (
	"context"

	appservice "github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/communication"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
)

type ReplyToUser struct {
	messenger appservice.Messenger
}

func (u *ReplyToUser) Run(ctx context.Context, otomo communication.Otomo, userPrompt string) error {
	c := reasoning.NewContext()
	c.SetUserPrompt(userPrompt)

	rep, err := otomo.Think(ctx, c)
	if err != nil {
		return failure.Wrap(err,
			failure.WithCode(errorcode.ErrInternal),
			failure.Message("failed to think"),
		)
	}

	if err := u.messenger.PostMessage(ctx, "", "", rep.Body()); err != nil {
		return failure.Wrap(err,
			failure.WithCode(errorcode.ErrInternal),
			failure.Message("failed to send reply"),
		)
	}

	return nil
}

func NewReplyToUser(messenger appservice.Messenger) *ReplyToUser {
	return &ReplyToUser{
		messenger: messenger,
	}
}
