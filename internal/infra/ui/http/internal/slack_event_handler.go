package internal

import (
	"net/http"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/app/usecase"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/infra/tool"
	"github.com/mackee/tanukirpc"
)

type slackEventRequest struct {
	ChannelID string `json:"channel"`
	Message   string `json:"message"`
}

type slackEventResponse struct {
	Status string `json:"status"`
}

func slackEventHandler(ctx tanukirpc.Context[*registry], req *slackEventRequest) (*slackEventResponse, error) {
	otomo, err := chat.NewOtomo(ctx.Registry().Brain)
	if err != nil {
		return nil, tanukirpc.WrapErrorWithStatus(http.StatusInternalServerError, err)
	}
	if p := config.Config.LLM.SystemPrompt; p != "" {
		otomo.SetSystemPrompt(chat.SystemPrompt(p))
	}

	cid, err := core.NewChannelID(req.ChannelID)
	if err != nil {
		return nil, tanukirpc.WrapErrorWithStatus(http.StatusBadRequest, err)
	}

	// Register WebSearchTool and WebFetchTool
	tools := []reasoning.Tool{
		tool.NewWebSearchTool(config.Config.Tool.WebSearch),
		tool.NewWebFetchTool(config.Config.Tool.WebFetch),
	}
	uc := usecase.NewReplyToUser(ctx.Registry().Slack, tools)
	if err := uc.Run(ctx, otomo, cid, core.PromptBody(req.Message)); err != nil {
		return nil, tanukirpc.WrapErrorWithStatus(http.StatusInternalServerError, err)
	}

	res := slackEventResponse{
		Status: "ok",
	}

	return &res, nil
}
