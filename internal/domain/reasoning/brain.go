package reasoning

import (
	"context"
)

// Brain is an entity interface that represents the model reasoning capability.
type Brain interface {
	Think(context.Context, Context) (*Answer, error)
}

// BrainThinker is an entity interface for making inferences.
type BrainThinker interface {
	Think(context.Context, Context) (*Answer, error)
}

type brain struct {
	thinker BrainThinker
}

func (b *brain) Think(ctx context.Context, c Context) (*Answer, error) {
	return b.thinker.Think(ctx, c)
}

func NewBrain(thinker BrainThinker) Brain {
	return &brain{
		thinker: thinker,
	}
}
