package repository

import (
	"context"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/entity"
	drepo "github.com/handlename/otomo/internal/domain/repository"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/morikuni/failure/v2"
)

type GeneralBrain struct{}

// New implements repository.Brain.
func (g *GeneralBrain) New(ctx context.Context) (entity.Brain, error) {
	client, err := service.NewBedrock(ctx, config.Config.Bedrock.ModelID)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to create bedrock client"))
	}

	brain := &generalBrain{
		client: client,
	}

	return brain, nil
}

func NewGeneralBrain(ctx context.Context) drepo.Brain {
	return &GeneralBrain{}
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
