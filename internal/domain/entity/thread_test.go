package entity_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func Test_Thread_MessagesOrdered(t *testing.T) {
	tests := []struct {
		name     string
		input    []entity.ThreadMessage
		expected []entity.ThreadMessage
	}{
		{
			name:     "empty input",
			input:    []entity.ThreadMessage{},
			expected: []entity.ThreadMessage{},
		},
		{
			name: "orderd input",
			input: []entity.ThreadMessage{
				entity.NewThreadMessage(entity.ThreadMessageID("1"), "mes1"),
				entity.NewThreadMessage(entity.ThreadMessageID("2"), "mes2"),
			},
			expected: []entity.ThreadMessage{
				entity.NewThreadMessage(entity.ThreadMessageID("1"), "mes1"),
				entity.NewThreadMessage(entity.ThreadMessageID("2"), "mes2"),
			},
		},
		{
			name: "unordered input",
			input: []entity.ThreadMessage{
				entity.NewThreadMessage(entity.ThreadMessageID("2"), "mes2"),
				entity.NewThreadMessage(entity.ThreadMessageID("1"), "mes1"),
			},
			expected: []entity.ThreadMessage{
				entity.NewThreadMessage(entity.ThreadMessageID("1"), "mes1"),
				entity.NewThreadMessage(entity.ThreadMessageID("2"), "mes2"),
			},
		},
		{
			name: "duplicated input",
			input: []entity.ThreadMessage{
				entity.NewThreadMessage(entity.ThreadMessageID("1"), "mes1"),
				entity.NewThreadMessage(entity.ThreadMessageID("1"), "mes1"),
				entity.NewThreadMessage(entity.ThreadMessageID("2"), "mes2"),
			},
			expected: []entity.ThreadMessage{
				entity.NewThreadMessage(entity.ThreadMessageID("1"), "mes1"),
				entity.NewThreadMessage(entity.ThreadMessageID("2"), "mes2"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thread := entity.NewThread(entity.ThreadID("1234"))
			thread.AddMessages(tt.input...)

			got := thread.Messages()
			assert.Equal(t, tt.expected, got)
		})
	}
}
