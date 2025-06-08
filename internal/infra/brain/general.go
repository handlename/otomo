package brain

import (
	"context"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/morikuni/failure/v2"
)

var _ entity.BrainThinker = (*General)(nil)

type General struct {
	client *service.Bedrock
}

func NewGeneral(ctx context.Context) (entity.BrainThinker, error) {
	client, err := service.NewBedrock(ctx, config.Config.Bedrock.ModelID)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to create bedrock client"))
	}

	brain := &General{
		client: client,
	}

	return brain, nil
}

// Think implements entity.BrainThinker.
func (g *General) Think(ctx context.Context, c entity.Context) (*entity.Answer, error) {
	res, err := g.client.Invoke(ctx, c.Prompt().String())
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to invoke bedrock"))
	}

	ans := entity.NewAnswer(res)
	return ans, nil
}
