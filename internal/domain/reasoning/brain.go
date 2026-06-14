package reasoning

import (
	"context"
)

// BrainThinker is an entity interface for making inferences.
type BrainThinker interface {
	Think(context.Context, *Context) (*Answer, error)
}

// Brain represents the model reasoning capability.
type Brain struct {
	thinker BrainThinker
}

func (b *Brain) Think(ctx context.Context, c *Context) (*Answer, error) {
	return b.thinker.Think(ctx, c)
}

func NewBrain(thinker BrainThinker) *Brain {
	return &Brain{
		thinker: thinker,
	}
}
