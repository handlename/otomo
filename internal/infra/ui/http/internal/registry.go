package internal

import (
	"context"

	"github.com/handlename/otomo/config"
	drepo "github.com/handlename/otomo/internal/domain/repository"
	irepo "github.com/handlename/otomo/internal/infra/repository"
	"github.com/handlename/otomo/internal/infra/service"
)

type registry struct {
	Slack       *service.Slack
	RepoSession drepo.Session
	RepoBrain   drepo.Brain
}

func NewRegistry(ctx context.Context) *registry {
	return &registry{
		Slack:       service.NewSlack(config.Config.Slack.BotToken, config.Config.Slack.SigningSecret),
		RepoSession: &irepo.VolatileSession{}, // TODO: replace
		RepoBrain:   irepo.NewGeneralBrain(ctx),
	}
}
