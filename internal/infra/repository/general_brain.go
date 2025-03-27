package repository

import (
	"context"

	"github.com/handlename/otomo/internal/domain/entity"
	drepo "github.com/handlename/otomo/internal/domain/repository"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/morikuni/failure/v2"
)

type GeneralBrain struct {
	botToken string
	appToken string
}

// New implements repository.Brain.
func (g *GeneralBrain) New(ctx context.Context) (entity.Brain, error) {
	client, err := service.NewBedrock(ctx)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to create bedrock client"))
	}

	brain := &generalBrain{
		client: client,
	}

	return brain, nil
}

func NewGeneralBrain(ctx context.Context, botToken, appToken string) drepo.Brain {
	brain := &GeneralBrain{
		botToken: botToken,
		appToken: appToken,
	}

	return brain
}

var _ entity.Brain = (*generalBrain)(nil)

type generalBrain struct {
	client *service.Bedrock
}

// Think implements entity.Brain.
func (g *generalBrain) Think(ctx context.Context, context entity.Context, ins *entity.Instruction) (*entity.Answer, error) {
	res, err := g.client.Invoke(ctx, ins.Body())
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to invoke bedrock"))
	}

	ans := entity.NewAnswer(res)
	return ans, nil
}
