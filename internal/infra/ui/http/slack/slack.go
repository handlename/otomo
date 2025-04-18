package slack

import (
	"context"
	"net/http"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/app/usecase"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/handlename/otomo/internal/infra/ui/http/middleware"
)

func New(ctx context.Context, prefix string) http.Handler {
	slack := service.NewSlack(config.Config.Slack.BotToken, config.Config.Slack.SigningSecret)

	publisher := service.NewEventPublisher()
	usecase.NewAckInstruction(slack).Subscribe(publisher)

	reg := NewRegistry(ctx, publisher, slack)
	mids := []middleware.Middleware{
		middleware.NewRegistry(reg),
		middleware.NewSlackEventVerifier(slack),
		middleware.NewAccesslog(),
	}

	mux := http.NewServeMux()
	mux.Handle("POST "+prefix+"/event", middleware.Wrap(eventHandler, mids...))
	return mux
}
