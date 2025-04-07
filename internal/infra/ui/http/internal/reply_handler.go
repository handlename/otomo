package internal

import (
	"net/http"

	"github.com/handlename/otomo/internal/app/usecase"
	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/mackee/tanukirpc"
)

type replyRequest struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

type replyResponse struct {
	Message string `json:"message"`
}

func replyHandler(ctx tanukirpc.Context[*registry], req *replyRequest) (*replyResponse, error)  {
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
