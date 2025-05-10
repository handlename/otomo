package repository

import (
	"context"
	"fmt"

	"github.com/handlename/otomo/internal/domain/entity"
	drepo "github.com/handlename/otomo/internal/domain/repository"
)

type StrawBrain struct{}

// New implements repository.Brain.
func (g *StrawBrain) New(ctx context.Context) (entity.Brain, error) {
	return &strawBrain{}, nil
}

func NewStrawBrain(ctx context.Context) drepo.Brain {
	return &StrawBrain{}
}

var _ entity.Brain = (*strawBrain)(nil)

type strawBrain struct{}

// Think implements entity.Brain.
func (g *strawBrain) Think(_ context.Context, c entity.Context) (*entity.Answer, error) {
	ans := entity.NewAnswer(fmt.Sprintf(`Did you say "%s" ?`, c.GetUserPrompt().String()))
	return ans, nil
}
