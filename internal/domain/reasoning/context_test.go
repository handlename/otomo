package reasoning_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContext_Prompt(t *testing.T) {
	ctx := reasoning.NewContext()
	ctx.SetSystemPrompt("you are a helpful assistant")
	ctx.SetUserPrompt("hello")

	msg1, err := core.NewMessage(core.RoleUser, lo.Must(core.NewUserID("user1")), "hi")
	require.NoError(t, err)

	msg2, err := core.NewMessage(core.RoleAssistant, core.UserID{}, "hello there")
	require.NoError(t, err)

	err = ctx.SetMessages([]*core.Message{msg1, msg2})
	require.NoError(t, err)

	prompt := ctx.Prompt()
	expected := `<system_instruction>
you are a helpful assistant
</system_instruction>
<thread>
<message user=user1>
hi
</message user=user1>
<message role=assistant>
hello there
</message role=assistant>
</thread>
<user_question>
hello
</user_question>
`
	assert.Equal(t, expected, prompt.String())
}

func TestContext_Prompt_EdgeCases(t *testing.T) {
	t.Run("empty prompts and empty messages", func(t *testing.T) {
		ctx := reasoning.NewContext()
		// No prompts, no messages
		prompt := ctx.Prompt()
		expected := `<thread>
</thread>
`
		assert.Equal(t, expected, prompt.String())
	})

	t.Run("nil messages are filtered out", func(t *testing.T) {
		ctx := reasoning.NewContext()
		ctx.SetSystemPrompt("system instruction")

		msg, err := core.NewMessage(core.RoleUser, lo.Must(core.NewUserID("user1")), "hello")
		require.NoError(t, err)

		// Set messages containing a nil pointer
		err = ctx.SetMessages([]*core.Message{nil, msg, nil})
		require.NoError(t, err)

		prompt := ctx.Prompt()
		expected := `<system_instruction>
system instruction
</system_instruction>
<thread>
<message user=user1>
hello
</message user=user1>
</thread>
`
		assert.Equal(t, expected, prompt.String())
	})
}

func TestContext_ToolInteractions(t *testing.T) {
	c := reasoning.NewContext()
	c.SetUserPrompt("test user prompt")

	tc, err := reasoning.NewToolCall(
		mustToolCallID("call-1"),
		mustToolName("dummy_tool"),
		`{"text":"hello"}`,
	)
	require.NoError(t, err)

	err = c.AddToolUseResponse("Thinking...", []reasoning.ToolCall{tc})
	require.NoError(t, err)
	require.Len(t, c.Messages(), 2)
	assert.Equal(t, "user", c.Messages()[0].Role())
	assert.Equal(t, core.MessageBody("test user prompt"), c.Messages()[0].Content())

	assert.Equal(t, "assistant", c.Messages()[1].Role())
	assert.Equal(t, core.MessageBody("Thinking..."), c.Messages()[1].Content())
	assert.Equal(t, tc, c.Messages()[1].ToolCalls()[0])

	result, err := reasoning.NewToolResult(mustToolCallID("call-1"), `{"length":5}`, reasoning.ToolResultSuccess)
	require.NoError(t, err)
	assert.Equal(t, mustToolCallID("call-1"), result.ToolUseID())
	err = c.AddToolResults([]reasoning.ToolResult{result})
	require.NoError(t, err)
	require.Len(t, c.Messages(), 3)
	assert.Equal(t, "user", c.Messages()[2].Role())
	assert.Equal(t, result, c.Messages()[2].ToolResults()[0])
}

func TestNewContextMessage(t *testing.T) {
	tc, err := reasoning.NewToolCall(
		mustToolCallID("call-1"),
		mustToolName("dummy_tool"),
		`{"text":"hello"}`,
	)
	require.NoError(t, err)

	tr, err := reasoning.NewToolResult(mustToolCallID("call-1"), `{"length":5}`, reasoning.ToolResultSuccess)
	require.NoError(t, err)
	assert.Equal(t, mustToolCallID("call-1"), tr.ToolUseID())

	// Validate role validation
	_, err = reasoning.NewContextMessage("invalid", core.UserID{}, core.MessageBody("content"), nil, nil)
	assert.Error(t, err)

	// Content cannot be empty unless tool calls or results are present
	_, err = reasoning.NewContextMessage("user", core.UserID{}, core.MessageBody(""), nil, nil)
	assert.Error(t, err)

	// System message constraints
	_, err = reasoning.NewContextMessage("system", core.UserID{}, core.MessageBody(""), []reasoning.ToolCall{tc}, nil)
	assert.Error(t, err)
	_, err = reasoning.NewContextMessage("system", core.UserID{}, core.MessageBody(""), nil, []reasoning.ToolResult{tr})
	assert.Error(t, err)

	// User message constraints
	_, err = reasoning.NewContextMessage("user", core.UserID{}, core.MessageBody(""), []reasoning.ToolCall{tc}, nil)
	assert.Error(t, err)
	_, err = reasoning.NewContextMessage("user", core.UserID{}, core.MessageBody(""), nil, []reasoning.ToolResult{tr})
	assert.NoError(t, err)

	// Assistant message constraints
	_, err = reasoning.NewContextMessage("assistant", core.UserID{}, core.MessageBody(""), []reasoning.ToolCall{tc}, nil)
	assert.NoError(t, err)
	_, err = reasoning.NewContextMessage("assistant", core.UserID{}, core.MessageBody(""), nil, []reasoning.ToolResult{tr})
	assert.Error(t, err)
}

func TestContextMessage_Immutability(t *testing.T) {
	tc, err := reasoning.NewToolCall(
		mustToolCallID("call-1"),
		mustToolName("dummy_tool"),
		`{"text":"hello"}`,
	)
	require.NoError(t, err)

	toolCalls := []reasoning.ToolCall{tc}
	msg, err := reasoning.NewContextMessage("assistant", core.UserID{}, core.MessageBody("content"), toolCalls, nil)
	require.NoError(t, err)

	// Mutate input slice
	tc2, _ := reasoning.NewToolCall(mustToolCallID("call-2"), mustToolName("another"), `{}`)
	toolCalls[0] = tc2

	// Verify internal copy is not mutated
	assert.Equal(t, tc, msg.ToolCalls()[0])

	// Mutate returned slice
	retCalls := msg.ToolCalls()
	retCalls[0] = tc2
	assert.Equal(t, tc, msg.ToolCalls()[0])
}

func TestNewToolResult(t *testing.T) {
	t.Run("valid tool result", func(t *testing.T) {
		tr, err := reasoning.NewToolResult(mustToolCallID("call-1"), "output", reasoning.ToolResultSuccess)
		require.NoError(t, err)
		assert.Equal(t, mustToolCallID("call-1"), tr.ToolUseID())
		assert.Equal(t, "output", tr.Output())
		assert.Equal(t, reasoning.ToolResultSuccess, tr.Status())
	})

	t.Run("empty tool call ID returns error", func(t *testing.T) {
		_, err := reasoning.NewToolResult(reasoning.ToolCallID{}, "output", reasoning.ToolResultSuccess)
		assert.Error(t, err)
	})
}
