package entity

import (
	"context"

	"github.com/pkg/errors"
)

type Otomo interface {
	Think(context.Context, Context, *Instruction) (*Reply, error)
}

var _ Otomo = (*otomo)(nil)

type otomo struct {
	brain Brain
}

func NewOtomo(brain Brain) *otomo {
	return &otomo{
		brain: brain,
	}
}

func (o *otomo) Think(ctx context.Context, context Context, instruction *Instruction) (*Reply, error) {
	ans, err := o.brain.Think(ctx, context, instruction)
	if err != nil {
		return nil, errors.Wrap(err, "failed to think")
	}

	r := NewReply(ans.Body(), []string{})
	return r, nil
}
