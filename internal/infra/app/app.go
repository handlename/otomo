package app

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/fujiwara/ridge"
	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/handlename/otomo/internal/infra/brain"
	"github.com/handlename/otomo/internal/infra/tool"
	ihttp "github.com/handlename/otomo/internal/infra/ui/http"
	"github.com/handlename/otomo/internal/infra/ui/terminal"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

type App struct{}

func (a *App) Init() error {
	if s := config.Config.Slack; s.AppToken == "" || s.BotToken == "" {
		return failure.New(
			errorcode.ErrInvalidArgument,
			failure.Message("both of Slack App token and Slack Bot token are required"),
		)
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	if err := a.Init(); err != nil {
		return failure.Wrap(err, failure.Message("failed to init app"))
	}

	mux := ihttp.NewMux(ctx)
	ridge.RunWithContext(ctx, fmt.Sprintf(":%d", config.Config.Port), "/", mux)

	return nil
}

func (a *App) RunChat(ctx context.Context) error {
	// Redirect console logs to otomo.log to avoid breaking the TUI display
	logFile, err := os.OpenFile("otomo.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Logger = log.Output(logFile).With().Caller().Logger()
		defer logFile.Close()
	} else {
		log.Logger = log.Output(io.Discard)
	}

	brainThinker, err := brain.NewGeneral(ctx)
	if err != nil {
		return failure.Wrap(err, failure.Message("failed to create brain thinker"))
	}
	b, err := reasoning.NewBrain(brainThinker)
	if err != nil {
		return failure.Wrap(err)
	}
	otomo, err := chat.NewOtomo(b)
	if err != nil {
		return failure.Wrap(err)
	}
	if p := config.Config.LLM.SystemPrompt; p != "" {
		otomo.SetSystemPrompt(chat.SystemPrompt(p))
	}

	tools := []reasoning.Tool{
		tool.NewWebSearchTool(config.Config.Tool.WebSearch),
		tool.NewWebFetchTool(config.Config.Tool.WebFetch),
	}

	return terminal.StartChatTUI(ctx, otomo, tools)
}
