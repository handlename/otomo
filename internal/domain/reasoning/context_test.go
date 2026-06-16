package reasoning_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContext_Prompt(t *testing.T) {
	ctx := reasoning.NewContext()
	ctx.SetSystemPrompt("you are a helpful assistant")
	ctx.SetUserPrompt("hello")

	msg1, err := core.NewMessage(core.RoleUser, "user1", "hi")
	require.NoError(t, err)

	msg2, err := core.NewMessage(core.RoleAssistant, "", "hello there")
	require.NoError(t, err)

	ctx.SetMessages([]*core.Message{msg1, msg2})

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

		msg, err := core.NewMessage(core.RoleUser, "user1", "hello")
		require.NoError(t, err)

		// Set messages containing a nil pointer
		ctx.SetMessages([]*core.Message{nil, msg, nil})

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
