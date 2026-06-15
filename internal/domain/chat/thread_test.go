package chat_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Thread_MessagesOrdered(t *testing.T) {
	tests := []struct {
		name     string
		input    []*chat.ThreadMessage
		expected []*chat.ThreadMessage
	}{
		{
			name:     "empty input",
			input:    []*chat.ThreadMessage{},
			expected: []*chat.ThreadMessage{},
		},
		{
			name: "ordered input",
			input: []*chat.ThreadMessage{
				lo.Must(chat.NewThreadMessage(chat.ThreadMessageID("1"), core.UserID("alice"), core.MessageBody("mes1"))),
				lo.Must(chat.NewThreadMessage(chat.ThreadMessageID("2"), core.UserID("bob"), core.MessageBody("mes2"))),
			},
			expected: []*chat.ThreadMessage{
				lo.Must(chat.NewThreadMessage(chat.ThreadMessageID("1"), core.UserID("alice"), core.MessageBody("mes1"))),
				lo.Must(chat.NewThreadMessage(chat.ThreadMessageID("2"), core.UserID("bob"), core.MessageBody("mes2"))),
			},
		},
		{
			name: "unordered input",
			input: []*chat.ThreadMessage{
				lo.Must(chat.NewThreadMessage(chat.ThreadMessageID("2"), core.UserID("bob"), core.MessageBody("mes2"))),
				lo.Must(chat.NewThreadMessage(chat.ThreadMessageID("1"), core.UserID("alice"), core.MessageBody("mes1"))),
			},
			expected: []*chat.ThreadMessage{
				lo.Must(chat.NewThreadMessage(chat.ThreadMessageID("1"), core.UserID("alice"), core.MessageBody("mes1"))),
				lo.Must(chat.NewThreadMessage(chat.ThreadMessageID("2"), core.UserID("bob"), core.MessageBody("mes2"))),
			},
		},
		{
			name: "duplicated input",
			input: []*chat.ThreadMessage{
				lo.Must(chat.NewThreadMessage(chat.ThreadMessageID("1"), core.UserID("alice"), core.MessageBody("mes1"))),
				lo.Must(chat.NewThreadMessage(chat.ThreadMessageID("1"), core.UserID("alice"), core.MessageBody("mes1"))),
				lo.Must(chat.NewThreadMessage(chat.ThreadMessageID("2"), core.UserID("bob"), core.MessageBody("mes2"))),
			},
			expected: []*chat.ThreadMessage{
				lo.Must(chat.NewThreadMessage(chat.ThreadMessageID("1"), core.UserID("alice"), core.MessageBody("mes1"))),
				lo.Must(chat.NewThreadMessage(chat.ThreadMessageID("2"), core.UserID("bob"), core.MessageBody("mes2"))),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thread, err := chat.NewThread(chat.ThreadID("1234"))
			require.NoError(t, err)
			thread.AddMessages(tt.input...)

			got := thread.Messages()
			assert.Equal(t, tt.expected, got)
		})
	}
}
