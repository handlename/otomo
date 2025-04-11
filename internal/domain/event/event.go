package event

import (
	"time"

	"github.com/google/uuid"
)

type ID string
type Kind string

type Event interface {
	ID() ID
	Kind() Kind
	OccuredAt() time.Time
	Data() any
}

var _ Event = (*baseEvent)(nil)

type baseEvent struct {
	id      ID
	kind    Kind
	ts      time.Time
	payload any
}

func newBaseEvent(kind Kind, payload any) baseEvent {
	return baseEvent{
		id:      ID(uuid.New().String()),
		kind:    kind,
		ts:      time.Now(),
		payload: payload,
	}
}

// Data implements Event.
func (e *baseEvent) Data() any {
	return e.payload
}

// ID implements Event.
func (e *baseEvent) ID() ID {
	return e.id
}

// Kind implements Event.
func (e *baseEvent) Kind() Kind {
	return e.kind
}

// OccuredAt implements Event.
func (e *baseEvent) OccuredAt() time.Time {
	return e.ts
}
