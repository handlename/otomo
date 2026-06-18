package tool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/handlename/otomo/internal/domain/reasoning"
)

var _ reasoning.Tool = (*DummyTool)(nil)

type DummyTool struct{}

func NewDummyTool() *DummyTool {
	return &DummyTool{}
}

func (t *DummyTool) Name() reasoning.ToolName {
	name, _ := reasoning.NewToolName("dummy_tool")
	return name
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
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", fmt.Errorf("failed to unmarshal inputs: %w", err)
	}

	length := len([]rune(input.Text))
	return fmt.Sprintf(`{"length": %d}`, length), nil
}
