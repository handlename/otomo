package chat_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/stretchr/testify/assert"
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
				chat.NewThreadMessage(chat.ThreadMessageID("1"), "alice", "mes1"),
				chat.NewThreadMessage(chat.ThreadMessageID("2"), "bob", "mes2"),
			},
			expected: []*chat.ThreadMessage{
				chat.NewThreadMessage(chat.ThreadMessageID("1"), "alice", "mes1"),
				chat.NewThreadMessage(chat.ThreadMessageID("2"), "bob", "mes2"),
			},
		},
		{
			name: "unordered input",
			input: []*chat.ThreadMessage{
				chat.NewThreadMessage(chat.ThreadMessageID("2"), "bob", "mes2"),
				chat.NewThreadMessage(chat.ThreadMessageID("1"), "alice", "mes1"),
			},
			expected: []*chat.ThreadMessage{
				chat.NewThreadMessage(chat.ThreadMessageID("1"), "alice", "mes1"),
				chat.NewThreadMessage(chat.ThreadMessageID("2"), "bob", "mes2"),
			},
		},
		{
			name: "duplicated input",
			input: []*chat.ThreadMessage{
				chat.NewThreadMessage(chat.ThreadMessageID("1"), "alice", "mes1"),
				chat.NewThreadMessage(chat.ThreadMessageID("1"), "alice", "mes1"),
				chat.NewThreadMessage(chat.ThreadMessageID("2"), "bob", "mes2"),
			},
			expected: []*chat.ThreadMessage{
				chat.NewThreadMessage(chat.ThreadMessageID("1"), "alice", "mes1"),
				chat.NewThreadMessage(chat.ThreadMessageID("2"), "bob", "mes2"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thread := chat.NewThread(chat.ThreadID("1234"))
			thread.AddMessages(tt.input...)

			got := thread.Messages()
			assert.Equal(t, tt.expected, got)
		})
	}
}
