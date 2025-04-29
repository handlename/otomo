package internal

import (
	"net/http"

	"github.com/handlename/otomo/internal/app/usecase"
	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/mackee/tanukirpc"
)

type slackEventRequest struct {
	ChannelID string `json:"role"`
	Message   string `json:"message"`
}

type slackEventResponse struct {
	Status string `json:"status"`
}

func slackEventHandler(ctx tanukirpc.Context[*registry], req *slackEventRequest) (*slackEventResponse, error)  {
	brain, err := ctx.Registry().RepoBrain.New(ctx)
	if err != nil {
		return nil, tanukirpc.WrapErrorWithStatus(http.StatusInternalServerError, err)
	}

	otomo , err := entity.NewOtomo(brain)
	if err != nil {
		return nil, tanukirpc.WrapErrorWithStatus(http.StatusInternalServerError, err)
	}

	inst := entity.NewInstruction("dummy", req.Message)
	uc := usecase.NewReplyToUser(ctx.Registry().RepoSession, ctx.Registry().Slack)
	if err := uc.Run(ctx, otomo, inst); err != nil {
		return nil, tanukirpc.WrapErrorWithStatus(http.StatusInternalServerError, err)
	}

	res := slackEventResponse{
		Status: "ok",
	}

	return &res, nil
}
