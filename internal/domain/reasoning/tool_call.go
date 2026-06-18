package reasoning

import "fmt"

type ToolCallID string
type ToolName string

// ToolCall represents a single tool execution request from the LLM.
type ToolCall struct {
	id        ToolCallID
	name      ToolName
	inputJSON string
}

// NewToolCall creates a new ToolCall value object.
func NewToolCall(id ToolCallID, name ToolName, inputJSON string) (ToolCall, error) {
	if id == "" {
		return ToolCall{}, fmt.Errorf("tool call id is required")
	}
	if name == "" {
		return ToolCall{}, fmt.Errorf("tool call name is required")
	}
	return ToolCall{
		id:        id,
		name:      name,
		inputJSON: inputJSON,
	}, nil
}

// ID returns the tool call ID.
func (tc ToolCall) ID() ToolCallID { return tc.id }

// Name returns the tool name.
func (tc ToolCall) Name() ToolName { return tc.name }

// InputJSON returns the raw tool input JSON string.
func (tc ToolCall) InputJSON() string { return tc.inputJSON }
