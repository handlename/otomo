package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/handlename/otomo/internal/domain/communication"
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
		ThinkFunc: func(context.Context, reasoning.Context) (*reasoning.Answer, error) {
			return reasoning.NewAnswer("mock response"), nil
		},
	})
	mockOtomo := communication.NewOtomo(mockBrain)
	mockMessenger := &service.MockMessenger{
		FetchThreadFunc: func(ctx context.Context, channelID string, threadID string) (communication.Thread, error) {
			return communication.NewThread(""), nil
		},
	}
	uc := NewReply(mockOtomo, mockMessenger)

	// Act

	input := ReplyInput{
		EventData: communication.InstructionReceivedData{
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
	mockBrain := reasoning.NewBrain(&brain.Mock{
		ThinkFunc: func(ctx context.Context, c reasoning.Context) (*reasoning.Answer, error) {
			return nil, mockError
		},
	})

	mockOtomo := (communication.NewOtomo(mockBrain))
	mockMessenger := &service.MockMessenger{}
	uc := NewReply(mockOtomo, mockMessenger)

	// Act

	input := ReplyInput{
		EventData: communication.InstructionReceivedData{
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
