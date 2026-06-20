package brain

import (
	"context"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/morikuni/failure/v2"
)

var _ reasoning.BrainThinker = (*General)(nil)

type General struct {
	client *service.Bedrock
}

func NewGeneral(ctx context.Context) (reasoning.BrainThinker, error) {
	client, err := service.NewBedrock(ctx, config.Config.LLM.ModelID)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to create bedrock client"))
	}

	brain := &General{
		client: client,
	}

	return brain, nil
}

// Think implements reasoning.BrainThinker.
func (g *General) Think(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
	ans, err := g.client.InvokeWithTools(ctx, c.SystemPromptBody(), c.Messages(), c.Tools())
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to invoke bedrock with tools"))
	}

	return ans, nil
}
