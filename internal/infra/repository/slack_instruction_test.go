package repository

import (
	"context"
	"testing"
	"time"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/event"
	"github.com/stretchr/testify/assert"
)

func Test_SlackInstruction_NewFromInstructionReceivedData(t *testing.T) {
	// Store the original bot user ID
	originalBotUserID := config.Config.Slack.BotUserID
	// Set a known bot user ID for testing
	config.Config.Slack.BotUserID = "U12345678"
	// Restore the original bot user ID after the test
	defer func() { config.Config.Slack.BotUserID = originalBotUserID }()

	// Create test cases
	tests := []struct {
		name           string
		rawInstruction string
		expectedBody   string
	}{
		{
			name:           "should trim spaces",
			rawInstruction: "  hello world  ",
			expectedBody:   "hello world",
		},
		{
			name:           "should remove bot user ID",
			rawInstruction: "<U12345678> hello world",
			expectedBody:   " hello world",
		},
		{
			name:           "should trim spaces and remove bot user ID",
			rawInstruction: "  <U12345678> hello world  ",
			expectedBody:   " hello world",
		},
		{
			name:           "should handle message with no spaces to trim or bot user ID",
			rawInstruction: "hello world",
			expectedBody:   "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare
			data := event.InstructionReceivedData{
				ChannelID:      "C12345678",
				MessageID:      "M12345678",
				ThreadID:       "T12345678",
				RawInstruction: tt.rawInstruction,
				SentAt:         time.Now(),
			}
			repo := NewSlackInstruction()

			// Run
			ctx := context.Background()
			instruction := repo.NewFromInstructionReceivedData(ctx, data)

			// Check
			assert.Equal(t, tt.expectedBody, instruction.Body())
			assert.Equal(t, data.ThreadID, string(instruction.ID()))
		})
	}
}