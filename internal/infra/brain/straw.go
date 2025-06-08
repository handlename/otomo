package brain

import (
	"context"
	"fmt"

	"github.com/handlename/otomo/internal/domain/entity"
)

var _ entity.BrainThinker = (*Straw)(nil)

type Straw struct {
}

// Think implements entity.BrainThinker.
func (s *Straw) Think(_ context.Context, c entity.Context) (*entity.Answer, error) {
	ans := entity.NewAnswer(fmt.Sprintf(`Did you say "%s" ?`, c.GetUserPrompt().String()))
	return ans, nil
}
