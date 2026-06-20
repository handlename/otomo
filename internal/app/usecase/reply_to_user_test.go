package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTool struct {
	name        reasoning.ToolName
	description string
	inputSchema string
	executeFunc func(ctx context.Context, inputJSON string) (string, error)
}

func (m mockTool) Name() reasoning.ToolName { return m.name }
func (m mockTool) Description() string      { return m.description }
func (m mockTool) InputSchema() string      { return m.inputSchema }
func (m mockTool) Execute(ctx context.Context, inputJSON string) (string, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, inputJSON)
	}
	return "", nil
}

func TestReplyToUser_Run(t *testing.T) {
	ctx := t.Context()

	t.Run("no tool calls", func(t *testing.T) {
		mockBrain, err := reasoning.NewBrain(&mockBrain{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				return reasoning.NewAnswer(reasoning.AnswerBody("hello from brain"), nil)
			},
		})
		require.NoError(t, err)

		mockOtomo, err := chat.NewOtomo(mockBrain)
		require.NoError(t, err)

		messenger := &mockMessenger{}
		uc := NewReplyToUser(messenger, []reasoning.Tool{})

		err = uc.Run(ctx, mockOtomo, lo.Must(core.NewChannelID("C1")), core.PromptBody("hello"))
		require.NoError(t, err)

		require.Len(t, messenger.History, 1)
		assert.Equal(t, "C1", messenger.History[0].ChannelID)
		assert.Equal(t, "hello from brain", messenger.History[0].Message)
	})

	t.Run("with tool calls and loop", func(t *testing.T) {
		toolCallCount := 0
		mockBrain, err := reasoning.NewBrain(&mockBrain{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				if toolCallCount == 0 {
					toolCallCount++
					tc, err := reasoning.NewToolCall(
						lo.Must(reasoning.NewToolCallID("call-1")),
						lo.Must(reasoning.NewToolName("mock_tool")),
						`{"param": "val"}`,
					)
					if err != nil {
						return nil, err
					}
					return reasoning.NewAnswer(reasoning.AnswerBody("calling tool..."), []reasoning.ToolCall{tc})
				}

				// Verify context has tool results
				messages := c.Messages()
				require.GreaterOrEqual(t, len(messages), 2)

				// Message 1 (assistant): thinking with tool call
				assert.Equal(t, "assistant", messages[len(messages)-2].Role())
				assert.Equal(t, core.MessageBody("calling tool..."), messages[len(messages)-2].Content())
				assert.Len(t, messages[len(messages)-2].ToolCalls(), 1)
				assert.Equal(t, "mock_tool", messages[len(messages)-2].ToolCalls()[0].Name().Value())

				// Message 2 (user): tool results
				assert.Equal(t, "user", messages[len(messages)-1].Role())
				assert.Len(t, messages[len(messages)-1].ToolResults(), 1)
				assert.Equal(t, "call-1", messages[len(messages)-1].ToolResults()[0].ToolUseID().Value())
				assert.Equal(t, `{"result": "ok"}`, messages[len(messages)-1].ToolResults()[0].Output())

				return reasoning.NewAnswer(reasoning.AnswerBody("final response"), nil)
			},
		})
		require.NoError(t, err)

		mockOtomo, err := chat.NewOtomo(mockBrain)
		require.NoError(t, err)

		toolExecuted := false
		mTool := mockTool{
			name: lo.Must(reasoning.NewToolName("mock_tool")),
			executeFunc: func(ctx context.Context, inputJSON string) (string, error) {
				assert.Equal(t, `{"param": "val"}`, inputJSON)
				toolExecuted = true
				return `{"result": "ok"}`, nil
			},
		}

		messenger := &mockMessenger{}
		uc := NewReplyToUser(messenger, []reasoning.Tool{mTool})

		err = uc.Run(ctx, mockOtomo, lo.Must(core.NewChannelID("C1")), core.PromptBody("hello"))
		require.NoError(t, err)

		assert.True(t, toolExecuted)
		require.Len(t, messenger.History, 1)
		assert.Equal(t, "C1", messenger.History[0].ChannelID)
		assert.Equal(t, "final response", messenger.History[0].Message)
	})

	t.Run("tool not found error", func(t *testing.T) {
		toolCallCount := 0
		mockBrain, err := reasoning.NewBrain(&mockBrain{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				if toolCallCount == 0 {
					toolCallCount++
					tc, err := reasoning.NewToolCall(
						lo.Must(reasoning.NewToolCallID("call-1")),
						lo.Must(reasoning.NewToolName("missing_tool")),
						`{}`,
					)
					if err != nil {
						return nil, err
					}
					return reasoning.NewAnswer(reasoning.AnswerBody("calling tool..."), []reasoning.ToolCall{tc})
				}

				messages := c.Messages()
				require.GreaterOrEqual(t, len(messages), 1)
				lastMsg := messages[len(messages)-1]
				require.Len(t, lastMsg.ToolResults(), 1)
				assert.Contains(t, lastMsg.ToolResults()[0].Output(), "error: tool 'missing_tool' not found")
				assert.Equal(t, reasoning.ToolResultError, lastMsg.ToolResults()[0].Status())

				return reasoning.NewAnswer(reasoning.AnswerBody("final response"), nil)
			},
		})
		require.NoError(t, err)

		mockOtomo, err := chat.NewOtomo(mockBrain)
		require.NoError(t, err)

		messenger := &mockMessenger{}
		uc := NewReplyToUser(messenger, []reasoning.Tool{})

		err = uc.Run(ctx, mockOtomo, lo.Must(core.NewChannelID("C1")), core.PromptBody("hello"))
		require.NoError(t, err)

		require.Len(t, messenger.History, 1)
		assert.Equal(t, "final response", messenger.History[0].Message)
	})

	t.Run("tool execution error", func(t *testing.T) {
		toolCallCount := 0
		mockBrain, err := reasoning.NewBrain(&mockBrain{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				if toolCallCount == 0 {
					toolCallCount++
					tc, err := reasoning.NewToolCall(
						lo.Must(reasoning.NewToolCallID("call-1")),
						lo.Must(reasoning.NewToolName("error_tool")),
						`{}`,
					)
					if err != nil {
						return nil, err
					}
					return reasoning.NewAnswer(reasoning.AnswerBody("calling tool..."), []reasoning.ToolCall{tc})
				}

				messages := c.Messages()
				require.GreaterOrEqual(t, len(messages), 1)
				lastMsg := messages[len(messages)-1]
				require.Len(t, lastMsg.ToolResults(), 1)
				assert.Contains(t, lastMsg.ToolResults()[0].Output(), "error executing tool: execute error")
				assert.Equal(t, reasoning.ToolResultError, lastMsg.ToolResults()[0].Status())

				return reasoning.NewAnswer(reasoning.AnswerBody("final response"), nil)
			},
		})
		require.NoError(t, err)

		mockOtomo, err := chat.NewOtomo(mockBrain)
		require.NoError(t, err)

		mTool := mockTool{
			name: lo.Must(reasoning.NewToolName("error_tool")),
			executeFunc: func(ctx context.Context, inputJSON string) (string, error) {
				return "", fmt.Errorf("execute error")
			},
		}

		messenger := &mockMessenger{}
		uc := NewReplyToUser(messenger, []reasoning.Tool{mTool})

		err = uc.Run(ctx, mockOtomo, lo.Must(core.NewChannelID("C1")), core.PromptBody("hello"))
		require.NoError(t, err)

		require.Len(t, messenger.History, 1)
		assert.Equal(t, "final response", messenger.History[0].Message)
	})

	t.Run("max turns exceeded error", func(t *testing.T) {
		thinkCount := 0
		mockBrain, err := reasoning.NewBrain(&mockBrain{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				thinkCount++
				tc, err := reasoning.NewToolCall(
					lo.Must(reasoning.NewToolCallID(fmt.Sprintf("call-%d", thinkCount))),
					lo.Must(reasoning.NewToolName("mock_tool")),
					`{"param": "val"}`,
				)
				if err != nil {
					return nil, err
				}
				return reasoning.NewAnswer(reasoning.AnswerBody("calling tool..."), []reasoning.ToolCall{tc})
			},
		})
		require.NoError(t, err)

		mockOtomo, err := chat.NewOtomo(mockBrain)
		require.NoError(t, err)

		mTool := mockTool{
			name: lo.Must(reasoning.NewToolName("mock_tool")),
			executeFunc: func(ctx context.Context, inputJSON string) (string, error) {
				return `{"result": "ok"}`, nil
			},
		}

		messenger := &mockMessenger{}
		uc := NewReplyToUser(messenger, []reasoning.Tool{mTool})

		err = uc.Run(ctx, mockOtomo, lo.Must(core.NewChannelID("C1")), core.PromptBody("hello"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "too many tool execution turns")
		assert.Equal(t, 5, thinkCount)
	})
}
