package slack

import (
	"context"
	"net/http"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/handlename/otomo/internal/infra/ui/http/middleware"
)

func New(ctx context.Context, prefix string) http.Handler {
	publisher := service.NewEventPublisher()
	slack := service.NewSlack(config.Config.Slack.SigningSecret)
	reg := NewRegistry(ctx, publisher, slack)
	mids := []middleware.Middleware{
		middleware.NewRegistry(reg),
		middleware.NewAccesslog(),
	}
	mux := http.NewServeMux()
	mux.Handle("POST "+prefix+"/event", middleware.Wrap(eventHandler, mids...))
	return mux
}
