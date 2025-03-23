package repository

import (
	"context"
	"fmt"

	"github.com/handlename/otomo/internal/domain/entity"
	drepo "github.com/handlename/otomo/internal/domain/repository"
)

type StrawBrain struct{}

// New implements repository.Brain.
func (g *StrawBrain) New(context.Context) (entity.Brain, error) {
	return &strawBrain{}, nil
}

func NewStrawBrain(ctx context.Context) drepo.Brain {
	return &StrawBrain{}
}

var _ entity.Brain = (*strawBrain)(nil)

type strawBrain struct{}

// Think implements entity.Brain.
func (g *strawBrain) Think(_ context.Context, _ entity.Context, ins *entity.Instruction) (*entity.Answer, error) {
	ans := entity.NewAnswer(fmt.Sprintf(`Did you say "%s" ?`, ins.Body()))
	return ans, nil
}
