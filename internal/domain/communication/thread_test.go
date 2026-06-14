package communication_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/communication"
	"github.com/stretchr/testify/assert"
)

func Test_Thread_MessagesOrdered(t *testing.T) {
	tests := []struct {
		name     string
		input    []communication.ThreadMessage
		expected []communication.ThreadMessage
	}{
		{
			name:     "empty input",
			input:    []communication.ThreadMessage{},
			expected: []communication.ThreadMessage{},
		},
		{
			name: "orderd input",
			input: []communication.ThreadMessage{
				communication.NewThreadMessage(communication.ThreadMessageID("1"), "alice", "mes1"),
				communication.NewThreadMessage(communication.ThreadMessageID("2"), "bob", "mes2"),
			},
			expected: []communication.ThreadMessage{
				communication.NewThreadMessage(communication.ThreadMessageID("1"), "alice", "mes1"),
				communication.NewThreadMessage(communication.ThreadMessageID("2"), "bob", "mes2"),
			},
		},
		{
			name: "unordered input",
			input: []communication.ThreadMessage{
				communication.NewThreadMessage(communication.ThreadMessageID("2"), "bob", "mes2"),
				communication.NewThreadMessage(communication.ThreadMessageID("1"), "alice", "mes1"),
			},
			expected: []communication.ThreadMessage{
				communication.NewThreadMessage(communication.ThreadMessageID("1"), "alice", "mes1"),
				communication.NewThreadMessage(communication.ThreadMessageID("2"), "bob", "mes2"),
			},
		},
		{
			name: "duplicated input",
			input: []communication.ThreadMessage{
				communication.NewThreadMessage(communication.ThreadMessageID("1"), "alice", "mes1"),
				communication.NewThreadMessage(communication.ThreadMessageID("1"), "alice", "mes1"),
				communication.NewThreadMessage(communication.ThreadMessageID("2"), "bob", "mes2"),
			},
			expected: []communication.ThreadMessage{
				communication.NewThreadMessage(communication.ThreadMessageID("1"), "alice", "mes1"),
				communication.NewThreadMessage(communication.ThreadMessageID("2"), "bob", "mes2"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thread := communication.NewThread(communication.ThreadID("1234"))
			thread.AddMessages(tt.input...)

			got := thread.Messages()
			assert.Equal(t, tt.expected, got)
		})
	}
}
