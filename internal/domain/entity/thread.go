package entity

import (
	"slices"
	"strings"

	"github.com/samber/lo"
)

type ThreadID string

// Thread is a series of ThreadMessages
type Thread interface {
	ID() ThreadID

	// Messages returns slice of ThreadMessages.
	// Each message in the slice is ordered by their id in ascending order and be uniquified.
	Messages() []ThreadMessage

	// AddMessage adds a ThreadMessage to the Thread.
	AddMessage(ThreadMessage)

	// AddMessages adds multiple ThreadMessages to the Thread.
	AddMessages(...ThreadMessage)
}

func NewThread(id ThreadID) Thread {
	return &thread{
		id:       id,
		messages: []ThreadMessage{},
	}
}

type thread struct {
	id       ThreadID
	messages []ThreadMessage
}

// AddMessage implements Thread.
func (t *thread) AddMessage(msg ThreadMessage) {
	t.AddMessages(msg)
}

// AddMessages implements Thread.
func (t *thread) AddMessages(msgs ...ThreadMessage) {
	t.messages = append(t.messages, msgs...)
	t.sortMessages()
}

// ID implements Thread.
func (t *thread) ID() ThreadID {
	return t.id
}

// Messages implements Thread.
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
