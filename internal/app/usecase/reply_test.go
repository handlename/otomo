package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Reply_Run(t *testing.T) {
	ctx := t.Context()

	// Arrange

	mockBrain, err := reasoning.NewBrain(&mockBrain{
		ThinkFunc: func(context.Context, *reasoning.Context) (*reasoning.Answer, error) {
			return reasoning.NewAnswer(reasoning.AnswerBody("mock response"), nil)
		},
	})
	require.NoError(t, err)

	mockOtomo, err := chat.NewOtomo(mockBrain)
	require.NoError(t, err)

	mockMessenger := &mockMessenger{
		FetchThreadFunc: func(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID) (*chat.Thread, error) {
			return chat.NewThread(threadID)
		},
	}
	uc := NewReply(mockOtomo, mockMessenger, []reasoning.Tool{})

	eventData, err := chat.NewInstructionReceivedData(lo.Must(core.NewChannelID("Ctest-channel")), lo.Must(core.NewMessageID("1234567890.123456")), lo.Must(chat.NewThreadID("test-thread")), chat.RawInstruction("test instruction"), time.Now())
	require.NoError(t, err)

	// Act

	input := ReplyInput{
		EventData: eventData,
	}
	output, err := uc.Run(ctx, input)

	// Assert

	require.NoError(t, err)

	expect := &ReplyOutput{}
	assert.Equal(t, expect, output)

	// Verify messenger was called with correct args
	require.Equal(t, 1, len(mockMessenger.History))
	assert.Equal(t, "Ctest-channel", mockMessenger.History[0].ChannelID)
	assert.Equal(t, "1234567890.123456", mockMessenger.History[0].MessageID)
	assert.Equal(t, "mock response", mockMessenger.History[0].Message)
}

func Test_Reply_Run_WithThread(t *testing.T) {
	ctx := t.Context()

	// Arrange
	var receivedPrompt string
	mockBrain, err := reasoning.NewBrain(&mockBrain{
		ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
			receivedPrompt = c.Prompt().String()
			return reasoning.NewAnswer(reasoning.AnswerBody("mock response"), nil)
		},
	})
	require.NoError(t, err)

	mockOtomo, err := chat.NewOtomo(mockBrain)
	require.NoError(t, err)

	mockMessenger := &mockMessenger{
		FetchThreadFunc: func(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID) (*chat.Thread, error) {
			tld, err := chat.NewThread(threadID)
			if err != nil {
				return nil, err
			}
			msg1, _ := chat.NewThreadMessage(lo.Must(chat.NewThreadMessageID("1")), lo.Must(core.NewUserID("alice")), core.MessageBody("hello"))
			msg2, _ := chat.NewThreadMessage(lo.Must(chat.NewThreadMessageID("2")), lo.Must(core.NewUserID("bob")), core.MessageBody("world"))
			tld.AddMessages(msg1, msg2)
			return tld, nil
		},
	}
	uc := NewReply(mockOtomo, mockMessenger, []reasoning.Tool{})

	eventData, err := chat.NewInstructionReceivedData(lo.Must(core.NewChannelID("Ctest-channel")), lo.Must(core.NewMessageID("1234567890.123456")), lo.Must(chat.NewThreadID("test-thread")), chat.RawInstruction("test instruction"), time.Now())
	require.NoError(t, err)

	// Act
	input := ReplyInput{
		EventData: eventData,
	}
	output, err := uc.Run(ctx, input)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, &ReplyOutput{}, output)

	// Verify the context prompt has system prompt, thread messages and user prompt
	assert.Contains(t, receivedPrompt, "<system_instruction>")
	assert.Contains(t, receivedPrompt, "<thread>")
	assert.Contains(t, receivedPrompt, "<message user=alice>\nhello\n</message user=alice>")
	assert.Contains(t, receivedPrompt, "<message user=bob>\nworld\n</message user=bob>")
	assert.Contains(t, receivedPrompt, "<user_question>\ntest instruction\n</user_question>")
}

func Test_Reply_Run_Error(t *testing.T) {
	ctx := t.Context()

	// Arrange

	// Create mock brain with error
	mockError := assert.AnError
	mockBrain, err := reasoning.NewBrain(&mockBrain{
		ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
			return nil, mockError
		},
	})
	require.NoError(t, err)

	mockOtomo, err := chat.NewOtomo(mockBrain)
	require.NoError(t, err)

	mockMessenger := &mockMessenger{}
	uc := NewReply(mockOtomo, mockMessenger, []reasoning.Tool{})

	eventData, err := chat.NewInstructionReceivedData(lo.Must(core.NewChannelID("Ctest-channel")), lo.Must(core.NewMessageID("1234567890.123456")), lo.Must(chat.NewThreadID("test-thread")), chat.RawInstruction("test instruction"), time.Now())
	require.NoError(t, err)

	// Act

	input := ReplyInput{
		EventData: eventData,
	}
	output, err := uc.Run(ctx, input)

	// Assert

	require.Error(t, err)
	assert.Nil(t, output)

	// Verify messenger was not called
	assert.Equal(t, 0, len(mockMessenger.History))
}

func TestReply_Run_ErrorFeedback(t *testing.T) {
	originalConfig := config.Config
	defer func() { config.Config = originalConfig }()

	t.Run("default error feedback (reaction only)", func(t *testing.T) {
		config.Config.Slack.ErrorFeedback = config.ErrorFeedback{}

		mockMessenger := &mockMessenger{}
		mockBrain, err := reasoning.NewBrain(&mockBrain{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				return nil, errors.New("thinking error")
			},
		})
		require.NoError(t, err)

		mockOtomo, err := chat.NewOtomo(mockBrain)
		require.NoError(t, err)
		uc := NewReply(mockOtomo, mockMessenger, []reasoning.Tool{})

		eventData, err := chat.NewInstructionReceivedData(lo.Must(core.NewChannelID("C12345")), lo.Must(core.NewMessageID("1234567890.123456")), lo.Must(chat.NewThreadID("1234567890.123456")), chat.RawInstruction("hello"), time.Now())
		require.NoError(t, err)

		_, err = uc.Run(t.Context(), ReplyInput{
			EventData: eventData,
		})
		assert.Error(t, err)

		require.Equal(t, 1, len(mockMessenger.ReactionHistory))
		assert.Equal(t, "warning", mockMessenger.ReactionHistory[0].Emoji)
		assert.Equal(t, "C12345", mockMessenger.ReactionHistory[0].ChannelID)
		assert.Equal(t, "1234567890.123456", mockMessenger.ReactionHistory[0].MessageID)

		assert.Equal(t, 0, len(mockMessenger.UploadFileHistory))
	})

	t.Run("snippet posting enabled", func(t *testing.T) {
		enableSnippet := true
		config.Config.Slack.ErrorFeedback = config.ErrorFeedback{
			EnablePostSnippet: &enableSnippet,
		}

		mockMessenger := &mockMessenger{}
		mockBrain, err := reasoning.NewBrain(&mockBrain{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				return nil, errors.New("thinking error detail")
			},
		})
		require.NoError(t, err)

		mockOtomo, err := chat.NewOtomo(mockBrain)
		require.NoError(t, err)
		uc := NewReply(mockOtomo, mockMessenger, []reasoning.Tool{})

		eventData, err := chat.NewInstructionReceivedData(lo.Must(core.NewChannelID("C12345")), lo.Must(core.NewMessageID("1234567890.123456")), lo.Must(chat.NewThreadID("1234567890.123456")), chat.RawInstruction("hello"), time.Now())
		require.NoError(t, err)

		_, err = uc.Run(t.Context(), ReplyInput{
			EventData: eventData,
		})
		assert.Error(t, err)

		assert.Equal(t, 1, len(mockMessenger.ReactionHistory))
		require.Equal(t, 1, len(mockMessenger.UploadFileHistory))
		assert.Equal(t, "C12345", mockMessenger.UploadFileHistory[0].ChannelID)
		assert.Equal(t, "1234567890.123456", mockMessenger.UploadFileHistory[0].ThreadTS)
		assert.Contains(t, mockMessenger.UploadFileHistory[0].Content, "thinking error detail")
	})
}

func Test_Reply_Run_WithToolCallsAndLoop(t *testing.T) {
	ctx := t.Context()

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

	mockMessenger := &mockMessenger{}
	uc := NewReply(mockOtomo, mockMessenger, []reasoning.Tool{mTool})

	eventData, err := chat.NewInstructionReceivedData(
		lo.Must(core.NewChannelID("Ctest-channel")),
		lo.Must(core.NewMessageID("1234567890.123456")),
		lo.Must(chat.NewThreadID("test-thread")),
		chat.RawInstruction("test instruction"),
		time.Now(),
	)
	require.NoError(t, err)

	input := ReplyInput{
		EventData: eventData,
	}
	output, err := uc.Run(ctx, input)
	require.NoError(t, err)

	assert.True(t, toolExecuted)
	assert.Equal(t, &ReplyOutput{}, output)

	// Verify messenger was called with correct args (the final response message)
	require.Equal(t, 1, len(mockMessenger.History))
	assert.Equal(t, "Ctest-channel", mockMessenger.History[0].ChannelID)
	assert.Equal(t, "1234567890.123456", mockMessenger.History[0].MessageID)
	assert.Equal(t, "final response", mockMessenger.History[0].Message)
}
