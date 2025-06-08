package internal

import (
	"context"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/entity"
	drepo "github.com/handlename/otomo/internal/domain/repository"
	"github.com/handlename/otomo/internal/infra/brain"
	irepo "github.com/handlename/otomo/internal/infra/repository"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/morikuni/failure/v2"
)

type registry struct {
	Slack       *service.Slack
	RepoSession drepo.Session
	Brain       entity.Brain
}

func NewRegistry(ctx context.Context) (*registry, error) {
	brainThinker, err := brain.NewGeneral(ctx)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	return &registry{
		Slack:       service.NewSlack(config.Config.Slack.BotToken, config.Config.Slack.SigningSecret),
		RepoSession: &irepo.VolatileSession{}, // TODO: replace
		Brain:       entity.NewBrain(brainThinker),
	}, nil
}
