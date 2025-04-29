package entity

import (
	"context"

	vo "github.com/handlename/otomo/internal/domain/valueobject"
)

const BrainBasePromptUserPromptPlaceholder = "{{userPrompt}}"

type Brain interface {
	// Think returns the answer to the instruction.
	Think(context.Context, Context, vo.Prompt) (*Answer, error)
}
