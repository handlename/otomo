package reasoning_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/stretchr/testify/assert"
)

func TestNewAnswer(t *testing.T) {
	tc, err := reasoning.NewToolCall("call_123", "dummy_tool", `{"text":"hello"}`)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		body        reasoning.AnswerBody
		toolCalls   []reasoning.ToolCall
		expectError bool
	}{
		{
			name:        "valid body with no tool calls",
			body:        "This is a valid answer",
			toolCalls:   nil,
			expectError: false,
		},
		{
			name:        "empty body with tool calls should be valid",
			body:        "",
			toolCalls:   []reasoning.ToolCall{tc},
			expectError: false,
		},
		{
			name:        "empty body and empty tool calls should return error",
			body:        "",
			toolCalls:   nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := reasoning.NewAnswer(tt.body, tt.toolCalls)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, ans)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, ans)
				assert.Equal(t, tt.body, ans.Body())
				assert.Equal(t, tt.toolCalls, ans.ToolCalls())
			}
		})
	}
}
