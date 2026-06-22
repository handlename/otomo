package terminal

import (
	"context"
	"errors"
	"testing"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockBrainThinker struct {
	ThinkFunc func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error)
}

func (m *mockBrainThinker) Think(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
	return m.ThinkFunc(ctx, c)
}

type mockTool struct {
	name        string
	executeFunc func(ctx context.Context, inputJSON string) (string, error)
}

func (m *mockTool) Name() reasoning.ToolName {
	tn, _ := reasoning.NewToolName(m.name)
	return tn
}

func (m *mockTool) Description() string {
	return "mock tool"
}

func (m *mockTool) InputSchema() string {
	return "{}"
}

func (m *mockTool) Execute(ctx context.Context, inputJSON string) (string, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, inputJSON)
	}
	return "", nil
}

func TestFindTool(t *testing.T) {
	toolName1, _ := reasoning.NewToolName("tool1")
	tool1 := &mockTool{name: "tool1"}
	tool2 := &mockTool{name: "tool2"}
	tools := []reasoning.Tool{tool1, tool2}

	t.Run("found", func(t *testing.T) {
		res, ok := findTool(tools, toolName1)
		assert.True(t, ok)
		assert.Equal(t, tool1, res)
	})

	t.Run("not found", func(t *testing.T) {
		unknown, _ := reasoning.NewToolName("unknown")
		res, ok := findTool(tools, unknown)
		assert.False(t, ok)
		assert.Nil(t, res)
	})
}

func TestExecuteToolLoop(t *testing.T) {
	t.Run("no tool calls", func(t *testing.T) {
		thinker := &mockBrainThinker{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				return reasoning.NewAnswer("hello from otomo", nil)
			},
		}
		brain, err := reasoning.NewBrain(thinker)
		require.NoError(t, err)
		otomo, err := chat.NewOtomo(brain)
		require.NoError(t, err)

		c := reasoning.NewContext()
		ans, err := executeToolLoop(context.Background(), otomo, c, nil)
		require.NoError(t, err)
		assert.Equal(t, reasoning.AnswerBody("hello from otomo"), ans.Body())
	})

	t.Run("one tool call success", func(t *testing.T) {
		callCount := 0
		thinker := &mockBrainThinker{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				callCount++
				if callCount == 1 {
					tcId, _ := reasoning.NewToolCallID("call-1")
					tcName, _ := reasoning.NewToolName("search")
					tc, _ := reasoning.NewToolCall(tcId, tcName, `{"query":"test"}`)
					return reasoning.NewAnswer("", []reasoning.ToolCall{tc})
				}
				// Verify the context has the tool result
				messages := c.Messages()
				require.Len(t, messages, 3) // 1 user prompt, 1 assistant tool call, 1 user tool result
				return reasoning.NewAnswer("tool executed successfully", nil)
			},
		}
		brain, err := reasoning.NewBrain(thinker)
		require.NoError(t, err)
		otomo, err := chat.NewOtomo(brain)
		require.NoError(t, err)

		searchTool := &mockTool{
			name: "search",
			executeFunc: func(ctx context.Context, inputJSON string) (string, error) {
				return "search result", nil
			},
		}

		c := reasoning.NewContext()
		c.SetUserPrompt("use tool")
		ans, err := executeToolLoop(context.Background(), otomo, c, []reasoning.Tool{searchTool})
		require.NoError(t, err)
		assert.Equal(t, reasoning.AnswerBody("tool executed successfully"), ans.Body())
	})

	t.Run("one tool call execution error", func(t *testing.T) {
		callCount := 0
		thinker := &mockBrainThinker{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				callCount++
				if callCount == 1 {
					tcId, _ := reasoning.NewToolCallID("call-1")
					tcName, _ := reasoning.NewToolName("search")
					tc, _ := reasoning.NewToolCall(tcId, tcName, `{"query":"test"}`)
					return reasoning.NewAnswer("", []reasoning.ToolCall{tc})
				}
				messages := c.Messages()
				require.Len(t, messages, 3)
				// The result should indicate error status
				assert.Equal(t, reasoning.ToolResultError, messages[2].ToolResults()[0].Status())
				assert.Equal(t, "some tool error", messages[2].ToolResults()[0].Output())
				return reasoning.NewAnswer("handled error", nil)
			},
		}
		brain, err := reasoning.NewBrain(thinker)
		require.NoError(t, err)
		otomo, err := chat.NewOtomo(brain)
		require.NoError(t, err)

		searchTool := &mockTool{
			name: "search",
			executeFunc: func(ctx context.Context, inputJSON string) (string, error) {
				return "", errors.New("some tool error")
			},
		}

		c := reasoning.NewContext()
		c.SetUserPrompt("use tool")
		ans, err := executeToolLoop(context.Background(), otomo, c, []reasoning.Tool{searchTool})
		require.NoError(t, err)
		assert.Equal(t, reasoning.AnswerBody("handled error"), ans.Body())
	})

	t.Run("tool not found error", func(t *testing.T) {
		callCount := 0
		thinker := &mockBrainThinker{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				callCount++
				if callCount == 1 {
					tcId, _ := reasoning.NewToolCallID("call-1")
					tcName, _ := reasoning.NewToolName("missing")
					tc, _ := reasoning.NewToolCall(tcId, tcName, `{}`)
					return reasoning.NewAnswer("", []reasoning.ToolCall{tc})
				}
				messages := c.Messages()
				require.Len(t, messages, 3)
				assert.Equal(t, reasoning.ToolResultError, messages[2].ToolResults()[0].Status())
				assert.Equal(t, "tool not found", messages[2].ToolResults()[0].Output())
				return reasoning.NewAnswer("handled missing tool", nil)
			},
		}
		brain, err := reasoning.NewBrain(thinker)
		require.NoError(t, err)
		otomo, err := chat.NewOtomo(brain)
		require.NoError(t, err)

		c := reasoning.NewContext()
		c.SetUserPrompt("use missing tool")
		ans, err := executeToolLoop(context.Background(), otomo, c, nil)
		require.NoError(t, err)
		assert.Equal(t, reasoning.AnswerBody("handled missing tool"), ans.Body())
	})

	t.Run("too many tool execution turns error", func(t *testing.T) {
		thinker := &mockBrainThinker{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				tcId, _ := reasoning.NewToolCallID("call-loop")
				tcName, _ := reasoning.NewToolName("loop")
				tc, _ := reasoning.NewToolCall(tcId, tcName, `{}`)
				return reasoning.NewAnswer("", []reasoning.ToolCall{tc})
			},
		}
		brain, err := reasoning.NewBrain(thinker)
		require.NoError(t, err)
		otomo, err := chat.NewOtomo(brain)
		require.NoError(t, err)

		loopTool := &mockTool{name: "loop"}

		c := reasoning.NewContext()
		_, err = executeToolLoop(context.Background(), otomo, c, []reasoning.Tool{loopTool})
		require.Error(t, err)
		assert.True(t, failure.Is(err, errorcode.ErrInternal))
		assert.Contains(t, err.Error(), "too many tool execution turns")
	})
}
