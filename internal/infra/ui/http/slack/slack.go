package slack

import (
	"context"
	"net/http"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/app/usecase"
	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/infra/brain"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/handlename/otomo/internal/infra/ui/http/middleware"
	"github.com/samber/lo"
)

func New(ctx context.Context, prefix string) http.Handler {
	slack := service.NewSlack(config.Config.Slack.BotToken, config.Config.Slack.SigningSecret)
	brainThinker := lo.Must(brain.NewGeneral(ctx))
	brain := entity.NewBrain(brainThinker)
	otomo := entity.NewOtomo(brain)
	if p := config.Config.LLM.SystemPrompt; p != "" {
		otomo.SetSystemPrompt(p)
	}

	publisher := service.NewEventPublisher()
	usecase.NewAckInstruction(slack).Subscribe(publisher)
	usecase.NewReply(otomo, slack).Subscribe(publisher)

	reg := NewRegistry(ctx, publisher, slack)
	mids := []middleware.Middleware{
		middleware.NewRegistry(reg),
		middleware.NewSlackRetryIgnorere(),
		middleware.NewSlackEventVerifier(slack),
		middleware.NewAccesslog(),
		middleware.NewRecover(),
	}

	mux := http.NewServeMux()
	mux.Handle("POST "+prefix+"/event", middleware.Wrap(eventHandler, mids...))
	return mux
}
