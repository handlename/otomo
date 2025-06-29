package internal

import (
	"context"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/infra/brain"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/morikuni/failure/v2"
)

type registry struct {
	Slack *service.Slack
	Brain entity.Brain
}

func NewRegistry(ctx context.Context) (*registry, error) {
	brainThinker, err := brain.NewGeneral(ctx)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	return &registry{
		Slack: service.NewSlack(config.Config.Slack.BotToken, config.Config.Slack.SigningSecret),
		Brain: entity.NewBrain(brainThinker),
	}, nil
}
