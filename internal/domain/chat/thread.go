//go:generate go run ../../../tools/gen-vo -file=thread.go
package chat

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/samber/lo"
)

// @vo
type ThreadID struct {
	value string
}

// NewThreadID creates a new ThreadID with validation.
func NewThreadID(value string) (ThreadID, error) {
	if value == "" {
		return ThreadID{}, fmt.Errorf("thread ID cannot be empty")
	}
	return ThreadID{value: value}, nil
}

// Thread is an entity representing a sequence of messages in a single conversation context.
type Thread struct {
	id       ThreadID
	messages []*ThreadMessage
}

func NewThread(id ThreadID) (*Thread, error) {
	if id.Value() == "" {
		return nil, fmt.Errorf("thread ID is required")
	}
	return &Thread{
		id:       id,
		messages: []*ThreadMessage{},
	}, nil
}

func (t *Thread) AddMessage(msg *ThreadMessage) {
	t.AddMessages(msg)
}

func (t *Thread) AddMessages(msgs ...*ThreadMessage) {
	for _, msg := range msgs {
		if msg == nil {
			continue
		}
		t.messages = append(t.messages, msg)
	}
	t.sortMessages()
}

func (t *Thread) ID() ThreadID {
	return t.id
}

func (t *Thread) Messages() []*ThreadMessage {
	return slices.Clone(t.messages)
}

func (t *Thread) sortMessages() {
	t.messages = lo.UniqBy(t.messages, func(msg *ThreadMessage) ThreadMessageID {
		return msg.ID()
	})
	slices.SortFunc(t.messages, func(a, b *ThreadMessage) int {
		return cmp.Compare(a.id.Value(), b.id.Value())
	})
}
