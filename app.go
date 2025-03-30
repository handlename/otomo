package otomo

import (
	"context"
	"fmt"

	"github.com/fujiwara/ridge"
	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/errorcode"
	ihttp "github.com/handlename/otomo/internal/infra/ui/http"
	"github.com/morikuni/failure/v2"
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
