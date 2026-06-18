package reasoning

import (
	"fmt"
	"slices"
)

type AnswerBody string

// Answer is a value object representing the outcome of Brain reasoning.
type Answer struct {
	body      AnswerBody
	toolCalls []ToolCall
}

func (ans *Answer) Body() AnswerBody      { return ans.body }
func (ans *Answer) ToolCalls() []ToolCall { return slices.Clone(ans.toolCalls) }
func (ans *Answer) HasToolCalls() bool    { return len(ans.toolCalls) > 0 }

func NewAnswer(body AnswerBody, toolCalls []ToolCall) (*Answer, error) {
	if len(toolCalls) == 0 && body == "" {
		return nil, fmt.Errorf("answer body is required when there are no tool calls")
	}
	return &Answer{
		body:      body,
		toolCalls: slices.Clone(toolCalls),
	}, nil
}
