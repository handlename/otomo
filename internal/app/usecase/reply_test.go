package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/domain/event"
	vo "github.com/handlename/otomo/internal/domain/valueobject"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Reply_Run(t *testing.T) {
	ctx := t.Context()

	// Arrange

	mockBrain := &service.MockBrain{}
	mockOtomo := lo.Must(entity.NewOtomo(mockBrain))
	mockMessenger := &service.MockMessenger{}
	uc := NewReply(mockOtomo, mockMessenger)

	// Act

	input := ReplyInput{
		EventData: event.InstructionReceivedData{
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

func Test_Reply_Run_Error(t *testing.T) {
	ctx := t.Context()

	// Arrange

	// Create mock brain with error
	mockError := assert.AnError
	mockBrain := &service.MockBrain{
		ThinkFunc: func(ctx context.Context, ectx entity.Context, prompt vo.Prompt) (*entity.Answer, error) {
			return nil, mockError
		},
	}

	mockOtomo := lo.Must(entity.NewOtomo(mockBrain))
	mockMessenger := &service.MockMessenger{}
	uc := NewReply(mockOtomo, mockMessenger)

	// Act

	input := ReplyInput{
		EventData: event.InstructionReceivedData{
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

func Test_Reply_buildPrompt(t *testing.T) {
	// Store the original bot user ID
	originalBotUserID := config.Config.Slack.BotUserID
	// Set a known bot user ID for testing
	config.Config.Slack.BotUserID = "U12345678"
	// Restore the original bot user ID after the test
	defer func() { config.Config.Slack.BotUserID = originalBotUserID }()

	r := Reply{}

	// Create test cases
	tests := []struct {
		name     string
		raw      string
		expected vo.Prompt
	}{
		{
			name:     "should trim spaces",
			raw:      "  hello world  ",
			expected: vo.NewPlainPrompt("hello world"),
		},
		{
			name:     "should remove bot user ID",
			raw:      "<U12345678> hello world",
			expected: vo.NewPlainPrompt(" hello world"),
		},
		{
			name:     "should trim spaces and remove bot user ID",
			raw:      "  <U12345678> hello world  ",
			expected: vo.NewPlainPrompt(" hello world"),
		},
		{
			name:     "should handle message with no spaces to trim or bot user ID",
			raw:      "hello world",
			expected: vo.NewPlainPrompt("hello world"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.buildPrompt(tt.raw)
			assert.Equal(t, tt.expected, got)
		})
	}
}
