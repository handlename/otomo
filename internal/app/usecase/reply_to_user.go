package usecase

import (
	"context"

	"github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/entity"
	repo "github.com/handlename/otomo/internal/domain/repository"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
)

type ReplyToUser struct {
	repoSession repo.Session
	messenger   service.Messenger
}

func (u *ReplyToUser) Run(ctx context.Context, otomo entity.Otomo, userPrompt string) error {
	c := entity.NewContext()
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

func NewReplyToUser(repoSession repo.Session, messenger service.Messenger) *ReplyToUser {
	return &ReplyToUser{
		repoSession: repoSession,
		messenger:   messenger,
	}
}
