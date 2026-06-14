package brain

import (
	"context"
	"fmt"

	"github.com/handlename/otomo/internal/domain/reasoning"
)

var _ reasoning.BrainThinker = (*Straw)(nil)

type Straw struct {
}

// Think implements reasoning.BrainThinker.
func (s *Straw) Think(_ context.Context, c reasoning.Context) (*reasoning.Answer, error) {
	ans := reasoning.NewAnswer(fmt.Sprintf(`Did you say "%s" ?`, c.GetUserPrompt().String()))
	return ans, nil
}
