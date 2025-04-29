package entity

import (
	"context"
)

const BrainBasePromptUserPromptPlaceholder = "{{userPrompt}}"

type Brain interface {
	// Think returns the answer to the instruction.
	Think(context.Context, Context, *Instruction) (*Answer, error)
}
