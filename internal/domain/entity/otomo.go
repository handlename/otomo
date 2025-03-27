package entity

import (
	"context"

	"github.com/pkg/errors"
)

type Otomo struct {
	brain Brain
}

func NewOtomo(brain Brain) *Otomo {
	return &Otomo{
		brain: brain,
	}
}

func (o *Otomo) Think(ctx context.Context, context Context, instruction *Instruction) (*Reply, error) {
	ans, err := o.brain.Think(ctx, context, instruction)
	if err != nil {
		return nil, errors.Wrap(err, "failed to think")
	}

	r := NewReply(ans.Body(), []string{})
	return r, nil
}
