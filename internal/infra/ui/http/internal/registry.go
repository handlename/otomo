package internal

import (
	"context"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/infra/brain"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/morikuni/failure/v2"
)

type registry struct {
	Slack *service.Slack
	Brain *reasoning.Brain
}

func NewRegistry(ctx context.Context) (*registry, error) {
	brainThinker, err := brain.NewGeneral(ctx)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	b, err := reasoning.NewBrain(brainThinker)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	return &registry{
		Slack: service.NewSlack(config.Config.Slack.BotToken, config.Config.Slack.SigningSecret),
		Brain: b,
	}, nil
}
