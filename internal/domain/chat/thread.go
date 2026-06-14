package chat

import (
	"slices"
	"strings"

	"github.com/samber/lo"
)

type ThreadID string

// Thread is an entity representing a sequence of messages in a single conversation context.
type Thread interface {
	ID() ThreadID
	Messages() []ThreadMessage
	AddMessage(ThreadMessage)
	AddMessages(...ThreadMessage)
}

type thread struct {
	id       ThreadID
	messages []ThreadMessage
}

func NewThread(id ThreadID) Thread {
	return &thread{
		id:       id,
		messages: []ThreadMessage{},
	}
}

func (t *thread) AddMessage(msg ThreadMessage) {
	t.AddMessages(msg)
}

func (t *thread) AddMessages(msgs ...ThreadMessage) {
	t.messages = append(t.messages, msgs...)
	t.sortMessages()
}

func (t *thread) ID() ThreadID {
	return t.id
}

func (t *thread) Messages() []ThreadMessage {
	return t.messages
}

func (t *thread) sortMessages() {
	t.messages = lo.UniqBy(t.messages, func(msg ThreadMessage) ThreadMessageID {
		return msg.ID()
	})
	slices.SortFunc(t.messages, func(a, b ThreadMessage) int {
		return strings.Compare(string(a.ID()), string(b.ID()))
	})
}
