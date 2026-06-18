package chat_test

import (
	"testing"
	"time"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestNewInstructionReceivedData_Validation(t *testing.T) {
	validTime := time.Now()
	tests := []struct {
		name           string
		channelID      core.ChannelID
		messageID      core.MessageID
		threadID       chat.ThreadID
		rawInstruction chat.RawInstruction
		sentAt         time.Time
		expectErr      bool
	}{
		{
			name:           "valid data",
			channelID:      lo.Must(core.NewChannelID("C123")),
			messageID:      lo.Must(core.NewMessageID("123.456")),
			threadID:       lo.Must(chat.NewThreadID("T123")),
			rawInstruction: chat.RawInstruction("do something"),
			sentAt:         validTime,
			expectErr:      false,
		},
		{
			name:           "empty channel ID",
			channelID:      core.ChannelID{},
			messageID:      lo.Must(core.NewMessageID("123.456")),
			threadID:       lo.Must(chat.NewThreadID("T123")),
			rawInstruction: chat.RawInstruction("do something"),
			sentAt:         validTime,
			expectErr:      true,
		},
		{
			name:           "empty message ID",
			channelID:      lo.Must(core.NewChannelID("C123")),
			messageID:      core.MessageID{},
			threadID:       lo.Must(chat.NewThreadID("T123")),
			rawInstruction: chat.RawInstruction("do something"),
			sentAt:         validTime,
			expectErr:      true,
		},
		{
			name:           "empty thread ID",
			channelID:      lo.Must(core.NewChannelID("C123")),
			messageID:      lo.Must(core.NewMessageID("123.456")),
			threadID:       chat.ThreadID{},
			rawInstruction: chat.RawInstruction("do something"),
			sentAt:         validTime,
			expectErr:      true,
		},
		{
			name:           "empty raw instruction",
			channelID:      lo.Must(core.NewChannelID("C123")),
			messageID:      lo.Must(core.NewMessageID("123.456")),
			threadID:       lo.Must(chat.NewThreadID("T123")),
			rawInstruction: chat.RawInstruction(""),
			sentAt:         validTime,
			expectErr:      true,
		},
		{
			name:           "zero sent at time",
			channelID:      lo.Must(core.NewChannelID("C123")),
			messageID:      lo.Must(core.NewMessageID("123.456")),
			threadID:       lo.Must(chat.NewThreadID("T123")),
			rawInstruction: chat.RawInstruction("do something"),
			sentAt:         time.Time{},
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chat.NewInstructionReceivedData(tt.channelID, tt.messageID, tt.threadID, tt.rawInstruction, tt.sentAt)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.channelID, got.ChannelID())
				assert.Equal(t, tt.messageID, got.MessageID())
				assert.Equal(t, tt.threadID, got.ThreadID())
				assert.Equal(t, tt.rawInstruction, got.RawInstruction())
				assert.Equal(t, tt.sentAt, got.SentAt())
			}
		})
	}
}
