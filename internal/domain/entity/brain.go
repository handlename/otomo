package entity

import (
	"context"
)

const BrainSystemPromptUserPromptPlaceholder = "{{userPrompt}}"

type Brain interface {
	// Think returns the answer to the instruction.
	Think(context.Context, Context) (*Answer, error)
}

type BrainThinker interface {
	Think(context.Context, Context) (*Answer, error)
}

type brain struct {
	thinker BrainThinker
}

// Think implements Brain.
func (b *brain) Think(ctx context.Context, c Context) (*Answer, error) {
	return b.thinker.Think(ctx, c)
}

func NewBrain(thinker BrainThinker) Brain {
	return &brain{
		thinker: thinker,
	}
}
