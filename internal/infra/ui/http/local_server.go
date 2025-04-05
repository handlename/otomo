package http

import (
	"context"

	"connectrpc.com/connect"
	"github.com/handlename/otomo/internal/app/usecase"
	"github.com/handlename/otomo/internal/domain/entity"
	drepo "github.com/handlename/otomo/internal/domain/repository"
	"github.com/handlename/otomo/internal/infra/service"
	entityv1 "github.com/handlename/otomo/internal/proto/entity/v1"
	servicev1 "github.com/handlename/otomo/internal/proto/service/v1"
	"github.com/handlename/otomo/internal/proto/service/v1/servicev1connect"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

var _ servicev1connect.LocalHandler = (*LocalHandler)(nil)

type LocalHandler struct {
	RepoSession drepo.Session
	RepoBrain drepo.Brain
}

// GetReply implements servicev1connect.LocalHandler.
func (h *LocalHandler) GetReply(ctx context.Context, req *connect.Request[servicev1.LocalGetReplyRequest]) (*connect.Response[servicev1.LocalGetReplyResponse], error) {
	msgr := &service.NopSlack{}

	brain, err := h.RepoBrain.New(ctx)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to new brain repository"))
	}

	otomo := entity.NewOtomo(brain)
	inst := entity.NewInstruction("dummy", req.Msg.GetPrompt().GetMessage())
	log.Debug().
		Any("req", req).
		Any("inst", inst).
		Str("message", req.Msg.GetPrompt().GetMessage()).
		Msg("check inst")

	uc := usecase.NewReplyToUser(h.RepoSession, msgr)
	if err := uc.Run(ctx, otomo, inst); err != nil {
		return nil, failure.Wrap(err)
	}

	res := connect.NewResponse(&servicev1.LocalGetReplyResponse{
		Answer: &entityv1.Answer{
			Message: msgr.Memory,
			Cost:    0,
		},
	})

	return res, nil
}
