package entity

import (
	"bytes"
	"context"
	_ "embed"
	"strings"
	"text/template"

	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type Brain interface {
	// Think returns the answer to the instruction.
	Think(context.Context, Context) (*Answer, error)

	// AddTool adds a tool to the brain.
	AddTool(context.Context, Tool) error

	// SelectTool selects the tool to use for the given instruction.
	SelectTool(context.Context, Context) (Tool, error)

	// UseTool executes the given tool with the provided parameters.
	UseTool(context.Context, Tool, ToolParams) (ToolAnswer, error)
}

type BrainThinker interface {
	Think(context.Context, Context) (*Answer, error)
}

type brain struct {
	tools   []Tool
	thinker BrainThinker
}

// AddTool implements Brain.
func (b *brain) AddTool(_ context.Context, tool Tool) error {
	b.tools = append(b.tools, tool)
	return nil
}

//go:embed brain_tool_select_system_prompt.tmpl.md
var brainToolSelectSystemPrompt string
var brainToolSelectSystemPromptTmpl = template.Must(template.New("brain_tool_select_system_prompt").Parse(brainToolSelectSystemPrompt))

// SelectTool implements Brain.
func (b *brain) SelectTool(ctx context.Context, c Context) (Tool, error) {
	if len(b.tools) == 0 {
		return nil, nil
	}

	// Context for tool selection
	sc := NewContext()
	var buf bytes.Buffer
	if err := brainToolSelectSystemPromptTmpl.Execute(&buf, map[string]any{
		"Tools": b.tools,
	}); err != nil {
		return nil, failure.Wrap(err)
	}

	sc.SetSystemPrompt(buf.String())
	sc.SetUserPrompt(c.GetUserPrompt().String())
	ans, err := b.Think(ctx, sc)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	candidate := strings.TrimSpace(ans.Body())
	if candidate == "" {
		return nil, nil
	}

	tool, found := lo.Find(b.tools, func(item Tool) bool {
		return item.Name() == candidate
	})
	if !found {
		log.Debug().Str("candidate", candidate).Msg("tool selected, but not exists")
		return nil, nil
	}

	return tool, nil
}

// Think implements Brain.
func (b *brain) Think(ctx context.Context, c Context) (*Answer, error) {
	return b.thinker.Think(ctx, c)
}

// UseTool implements Brain.
func (b *brain) UseTool(context.Context, Tool, ToolParams) (ToolAnswer, error) {
	panic("unimplemented")
}

func NewBrain(thinker BrainThinker) Brain {
	return &brain{
		tools:   []Tool{},
		thinker: thinker,
	}
}
