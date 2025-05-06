package usecase

import (
	"context"

	"github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/entity"
	repo "github.com/handlename/otomo/internal/domain/repository"
	vo "github.com/handlename/otomo/internal/domain/valueobject"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/handlename/otomo/internal/infra/repository"
	"github.com/morikuni/failure/v2"
)

type ReplyToUser struct {
	repoSession repo.Session
	messenger   service.Messenger
}

func (u *ReplyToUser) Run(ctx context.Context, otomo entity.Otomo, prompt vo.Prompt) error {
	sctx := repository.SlackContext{}.New()
	// TODO: calls sctx.AddRefresher()

	rep, err := otomo.Think(ctx, sctx, prompt)
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
