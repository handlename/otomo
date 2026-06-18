package reasoning_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/stretchr/testify/assert"
)

func TestNewToolCall(t *testing.T) {
	tests := []struct {
		name        string
		id          reasoning.ToolCallID
		toolName    reasoning.ToolName
		inputJSON   string
		expectError bool
	}{
		{
			name:        "valid tool call",
			id:          "call_123",
			toolName:    "my_tool",
			inputJSON:   `{"param": "val"}`,
			expectError: false,
		},
		{
			name:        "empty tool call ID returns error",
			id:          "",
			toolName:    "my_tool",
			inputJSON:   `{"param": "val"}`,
			expectError: true,
		},
		{
			name:        "empty tool name returns error",
			id:          "call_123",
			toolName:    "",
			inputJSON:   `{"param": "val"}`,
			expectError: true,
		},
		{
			name:        "empty input JSON is allowed",
			id:          "call_123",
			toolName:    "my_tool",
			inputJSON:   "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc, err := reasoning.NewToolCall(tt.id, tt.toolName, tt.inputJSON)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, reasoning.ToolCall{}, tc)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, tc.ID())
				assert.Equal(t, tt.toolName, tc.Name())
				assert.Equal(t, tt.inputJSON, tc.InputJSON())
			}
		})
	}
}
