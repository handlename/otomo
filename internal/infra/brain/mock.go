package brain

import (
	"context"

	"github.com/handlename/otomo/internal/domain/entity"
)

var _ entity.BrainThinker = (*Mock)(nil)

type Mock struct {
	ThinkFunc func(context.Context, entity.Context) (*entity.Answer, error)
}

// Think implements entity.BrainThinker.
func (m *Mock) Think(ctx context.Context, c entity.Context) (*entity.Answer, error) {
	if m.ThinkFunc != nil {
		return m.ThinkFunc(ctx, c)
	}

	return entity.NewAnswer("mock response"), nil
}
