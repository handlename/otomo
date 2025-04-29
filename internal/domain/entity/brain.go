package entity

import (
	"context"
	"strings"
	"text/template"

	"github.com/morikuni/failure/v2"
)

const BrainBasePromptUserPromptPlaceholder = "{{userPrompt}}"

type Brain interface {
	Think(context.Context, Context, *Instruction) (*Answer, error)

	// SetBasePrompt sets the base prompt for the brain.
	// The prompt must contains placeholder `{{userPrompt}}`.
	SetBasePrompt(prompt string) error
}

var _ Brain = (*BaseBrain)(nil)

type BaseBrain struct {
	basePrompt string
}

// Think implements Brain.
func (b *BaseBrain) Think(ctx context.Context, ctxCtx Context, inst *Instruction) (*Answer, error) {
	panic("this must be implemented by a struct that embeds BaseBrain")
}

func NewBaseBrain() *BaseBrain {
	return &BaseBrain{}
}

// SetBasePrompt implements entity.Brain.
func (b *BaseBrain) SetBasePrompt(prompt string) error {
	tmpl, err := template.New("base_prompt").Parse(prompt)
	if err != nil {
		return failure.Wrap(err, failure.Message("failed to parse prompt as template"))
	}

	// Validate template by executing it
	var buf strings.Builder
	err = tmpl.Execute(&buf, map[string]string{"UserPrompt": ""})
	if err != nil {
		return failure.Wrap(err, failure.Message("failed to execute prompt as template"))
	}

	b.basePrompt = prompt
	return nil
}
