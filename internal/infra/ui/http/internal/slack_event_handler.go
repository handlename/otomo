package internal

import (
	"net/http"

	"github.com/handlename/otomo/config"
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
	otomo:= entity.NewOtomo(ctx.Registry().Brain)
	if p := config.Config.LLM.SystemPrompt; p != "" {
		otomo.SetSystemPrompt(p)
	}

	uc := usecase.NewReplyToUser( ctx.Registry().Slack)
	if err := uc.Run(ctx, otomo, req.Message); err != nil {
		return nil, tanukirpc.WrapErrorWithStatus(http.StatusInternalServerError, err)
	}

	res := slackEventResponse{
		Status: "ok",
	}

	return &res, nil
}
