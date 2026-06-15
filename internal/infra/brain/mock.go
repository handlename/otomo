package brain

import (
	"context"

	"github.com/handlename/otomo/internal/domain/reasoning"
)

var _ reasoning.BrainThinker = (*Mock)(nil)

type Mock struct {
	ThinkFunc func(context.Context, *reasoning.Context) (*reasoning.Answer, error)
}

// Think implements reasoning.BrainThinker.
func (m *Mock) Think(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
	if m.ThinkFunc != nil {
		return m.ThinkFunc(ctx, c)
	}

	ans, err := reasoning.NewAnswer("mock response")
	if err != nil {
		return nil, err
	}
	return ans, nil
}
