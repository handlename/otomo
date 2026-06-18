package tool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/samber/lo"
)

var _ reasoning.Tool = (*DummyTool)(nil)

type DummyTool struct {
	name reasoning.ToolName
}

func NewDummyTool() *DummyTool {
	return &DummyTool{
		name: lo.Must(reasoning.NewToolName("dummy_tool")),
	}
}

func (t *DummyTool) Name() reasoning.ToolName {
	return t.name
}

func (t *DummyTool) Description() string {
	return "A dummy tool for verification. It returns the character length of the input text parameter."
}

func (t *DummyTool) InputSchema() string {
	return `{
		"type": "object",
		"properties": {
			"text": {
				"type": "string",
				"description": "The text to count characters."
			}
		},
		"required": ["text"]
	}`
}

func (t *DummyTool) Execute(ctx context.Context, inputJSON string) (string, error) {
	var input struct {
		Text *string `json:"text"`
	}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", failure.Wrap(err, failure.WithCode(errorcode.ErrInvalidArgument), failure.Message("failed to unmarshal inputs"))
	}

	if input.Text == nil {
		return "", failure.New(errorcode.ErrInvalidArgument, failure.Message("text is required"))
	}

	length := len([]rune(*input.Text))
	return fmt.Sprintf(`{"length": %d}`, length), nil
}
