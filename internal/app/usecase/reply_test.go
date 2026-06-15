package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/infra/brain"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Reply_Run(t *testing.T) {
	ctx := t.Context()

	// Arrange

	mockBrain := reasoning.NewBrain(&brain.Mock{
		ThinkFunc: func(context.Context, *reasoning.Context) (*reasoning.Answer, error) {
			return reasoning.NewAnswer("mock response"), nil
		},
	})
	mockOtomo := chat.NewOtomo(mockBrain)
	mockMessenger := &service.MockMessenger{
		FetchThreadFunc: func(ctx context.Context, channelID string, threadID string) (*chat.Thread, error) {
			return chat.NewThread(""), nil
		},
	}
	uc := NewReply(mockOtomo, mockMessenger)

	// Act

	input := ReplyInput{
		EventData: chat.InstructionReceivedData{
			ChannelID:      "test-channel",
			MessageID:      "test-message",
			ThreadID:       "test-thread",
			RawInstruction: "test instruction",
			SentAt:         time.Now(),
		},
	}
	output, err := uc.Run(ctx, input)

	// Assert

	require.NoError(t, err)

	expect := &ReplyOutput{}
	assert.Equal(t, expect, output)

	// Verify messenger was called with correct args
	require.Equal(t, 1, len(mockMessenger.History))
	assert.Equal(t, "test-channel", mockMessenger.History[0].ChannelID)
	assert.Equal(t, "test-message", mockMessenger.History[0].MessageID)
	assert.Equal(t, "mock response", mockMessenger.History[0].Message)
}

func Test_Reply_Run_WithThread(t *testing.T) {
	ctx := t.Context()

	// Arrange
	var receivedPrompt string
	mockBrain := reasoning.NewBrain(&brain.Mock{
		ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
			receivedPrompt = c.Prompt().String()
			return reasoning.NewAnswer("mock response"), nil
		},
	})
	mockOtomo := chat.NewOtomo(mockBrain)
	mockMessenger := &service.MockMessenger{
		FetchThreadFunc: func(ctx context.Context, channelID string, threadID string) (*chat.Thread, error) {
			tld := chat.NewThread(chat.ThreadID(threadID))
			tld.AddMessages(
				chat.NewThreadMessage(chat.ThreadMessageID("1"), "alice", "hello"),
				chat.NewThreadMessage(chat.ThreadMessageID("2"), "bob", "world"),
			)
			return tld, nil
		},
	}
	uc := NewReply(mockOtomo, mockMessenger)

	// Act
	input := ReplyInput{
		EventData: chat.InstructionReceivedData{
			ChannelID:      "test-channel",
			MessageID:      "test-message",
			ThreadID:       "test-thread",
			RawInstruction: "test instruction",
			SentAt:         time.Now(),
		},
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
	mockBrain := reasoning.NewBrain(&brain.Mock{
		ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
			return nil, mockError
		},
	})

	mockOtomo := (chat.NewOtomo(mockBrain))
	mockMessenger := &service.MockMessenger{}
	uc := NewReply(mockOtomo, mockMessenger)

	// Act

	input := ReplyInput{
		EventData: chat.InstructionReceivedData{
			ChannelID:      "test-channel",
			MessageID:      "test-message",
			ThreadID:       "test-thread",
			RawInstruction: "test instruction",
			SentAt:         time.Now(),
		},
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

		mockMessenger := &service.MockMessenger{}
		mockOtomo := chat.NewOtomo(reasoning.NewBrain(&brain.Mock{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				return nil, errors.New("thinking error")
			},
		}))
		uc := NewReply(mockOtomo, mockMessenger)

		_, err := uc.Run(t.Context(), ReplyInput{
			EventData: chat.InstructionReceivedData{
				ChannelID:      "C12345",
				MessageID:      "1234567890.123456",
				RawInstruction: "hello",
			},
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

		mockMessenger := &service.MockMessenger{}
		mockOtomo := chat.NewOtomo(reasoning.NewBrain(&brain.Mock{
			ThinkFunc: func(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
				return nil, errors.New("thinking error detail")
			},
		}))
		uc := NewReply(mockOtomo, mockMessenger)

		_, err := uc.Run(t.Context(), ReplyInput{
			EventData: chat.InstructionReceivedData{
				ChannelID:      "C12345",
				MessageID:      "1234567890.123456",
				RawInstruction: "hello",
			},
		})
		assert.Error(t, err)

		assert.Equal(t, 1, len(mockMessenger.ReactionHistory))
		require.Equal(t, 1, len(mockMessenger.UploadFileHistory))
		assert.Equal(t, "C12345", mockMessenger.UploadFileHistory[0].ChannelID)
		assert.Equal(t, "1234567890.123456", mockMessenger.UploadFileHistory[0].ThreadTS)
		assert.Contains(t, mockMessenger.UploadFileHistory[0].Content, "thinking error detail")
	})
}
