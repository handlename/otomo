package reasoning

import "fmt"

type AnswerBody string

// ToolCall represents a single tool execution request from the LLM.
type ToolCall struct {
	id        string
	name      string
	inputJSON string
}

func NewToolCall(id, name, inputJSON string) (ToolCall, error) {
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

func (tc ToolCall) ID() string        { return tc.id }
func (tc ToolCall) Name() string      { return tc.name }
func (tc ToolCall) InputJSON() string { return tc.inputJSON }

// Answer is a value object representing the outcome of Brain reasoning.
type Answer struct {
	body      AnswerBody
	toolCalls []ToolCall
}

func (ans *Answer) Body() AnswerBody      { return ans.body }
func (ans *Answer) ToolCalls() []ToolCall { return ans.toolCalls }
func (ans *Answer) HasToolCalls() bool    { return len(ans.toolCalls) > 0 }

func NewAnswer(body AnswerBody, toolCalls []ToolCall) (*Answer, error) {
	if len(toolCalls) == 0 && body == "" {
		return nil, fmt.Errorf("answer body is required when there are no tool calls")
	}
	return &Answer{
		body:      body,
		toolCalls: toolCalls,
	}, nil
}
