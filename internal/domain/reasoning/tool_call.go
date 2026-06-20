//go:generate go run ../../../tools/gen-vo -file=tool_call.go
package reasoning

import "fmt"

// @vo
type ToolCallID struct {
	value string
}

func NewToolCallID(value string) (ToolCallID, error) {
	if value == "" {
		return ToolCallID{}, fmt.Errorf("tool call ID cannot be empty")
	}
	return ToolCallID{value: value}, nil
}

// @vo
type ToolName struct {
	value string
}

func NewToolName(value string) (ToolName, error) {
	if value == "" {
		return ToolName{}, fmt.Errorf("tool name cannot be empty")
	}
	return ToolName{value: value}, nil
}

// ToolCall represents a single tool execution request from the LLM.
type ToolCall struct {
	id        ToolCallID
	name      ToolName
	inputJSON string
}

// NewToolCall creates a new ToolCall value object.
func NewToolCall(id ToolCallID, name ToolName, inputJSON string) (ToolCall, error) {
	if id.Value() == "" {
		return ToolCall{}, fmt.Errorf("tool call id is required")
	}
	if name.Value() == "" {
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
