package http

import (
	"context"
	"net/http"

	"github.com/handlename/otomo/internal/app/usecase"
	"github.com/handlename/otomo/internal/domain/entity"
	drepo "github.com/handlename/otomo/internal/domain/repository"
	irepo "github.com/handlename/otomo/internal/infra/repository"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/mackee/tanukirpc"
)

type localRegistry struct {
	RepoSession drepo.Session
	RepoBrain   drepo.Brain
}

func setupLocal(ctx context.Context) *tanukirpc.Router[*localRegistry] {
	reg := localRegistry{
		RepoSession: &irepo.VolatileSession{},
		RepoBrain:   irepo.NewGeneralBrain(ctx),
	}
	r := tanukirpc.NewRouter(&reg)
	r.Post("/reply", tanukirpc.NewHandler(localReply))
	return r
}

type replyRequest struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

type replyResponse struct {
	Message string `json:"message"`
}

func localReply(ctx tanukirpc.Context[*localRegistry], req *replyRequest) (*replyResponse, error)  {
	msgr := &service.NopSlack{}

	brain, err := ctx.Registry().RepoBrain.New(ctx)
	if err != nil {
		return nil, tanukirpc.WrapErrorWithStatus(http.StatusInternalServerError, err)
	}

	otomo := entity.NewOtomo(brain)
	inst := entity.NewInstruction("dummy", req.Message)
	uc := usecase.NewReplyToUser(ctx.Registry().RepoSession, msgr)
	if err := uc.Run(ctx, otomo, inst); err != nil {
		return nil, tanukirpc.WrapErrorWithStatus(http.StatusInternalServerError, err)
	}

	res := replyResponse{
		Message: msgr.Memory,
	}

	return &res, nil
}
