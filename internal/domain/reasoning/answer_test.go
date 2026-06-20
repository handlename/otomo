package reasoning_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustToolCallID(v string) reasoning.ToolCallID {
	id, err := reasoning.NewToolCallID(v)
	if err != nil {
		panic(err)
	}
	return id
}

func mustToolName(v string) reasoning.ToolName {
	name, err := reasoning.NewToolName(v)
	if err != nil {
		panic(err)
	}
	return name
}

func TestNewAnswer(t *testing.T) {
	tc, err := reasoning.NewToolCall(
		mustToolCallID("call_123"),
		mustToolName("dummy_tool"),
		`{"text":"hello"}`,
	)
	require.NoError(t, err)

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

func TestAnswer_Immutability(t *testing.T) {
	tc1, err := reasoning.NewToolCall(
		mustToolCallID("call_1"),
		mustToolName("tool_1"),
		`{}`,
	)
	require.NoError(t, err)

	tc2, err := reasoning.NewToolCall(
		mustToolCallID("call_2"),
		mustToolName("tool_2"),
		`{}`,
	)
	require.NoError(t, err)

	toolCalls := []reasoning.ToolCall{tc1}
	ans, err := reasoning.NewAnswer("body", toolCalls)
	require.NoError(t, err)

	// Mutate the original slice used to construct the Answer.
	toolCalls[0] = tc2

	// Verify that the internal slice of Answer was defensively copied.
	assert.Equal(t, tc1, ans.ToolCalls()[0], "Answer constructor should defensively copy the input toolCalls slice")

	// Mutate the slice returned by the ToolCalls() getter.
	ret := ans.ToolCalls()
	ret[0] = tc2

	// Verify that the internal slice of Answer was not affected.
	assert.Equal(t, tc1, ans.ToolCalls()[0], "Answer.ToolCalls() getter should return a defensively copied slice")
}
