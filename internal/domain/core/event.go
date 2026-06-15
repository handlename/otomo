package core

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EventID string
type EventKind string

// Event is an entity interface representing something that has occurred in the domain.
type Event interface {
	ID() EventID
	Kind() EventKind
	OccuredAt() time.Time
	Data() any
	String() string
}

var _ Event = (*BaseEvent)(nil)

// BaseEvent is a value object helper that provides standard fields for Event implementations.
type BaseEvent struct {
	id      EventID
	kind    EventKind
	ts      time.Time
	payload any
}

func NewBaseEvent(kind EventKind, payload any) (BaseEvent, error) {
	if kind == "" {
		return BaseEvent{}, fmt.Errorf("event kind is required")
	}
	return BaseEvent{
		id:      EventID(uuid.New().String()),
		kind:    kind,
		ts:      time.Now(),
		payload: payload,
	}, nil
}

// Data implements Event.
func (e *BaseEvent) Data() any {
	return e.payload
}

// ID implements Event.
func (e *BaseEvent) ID() EventID {
	return e.id
}

// Kind implements Event.
func (e *BaseEvent) Kind() EventKind {
	return e.kind
}

// OccuredAt implements Event.
func (e *BaseEvent) OccuredAt() time.Time {
	return e.ts
}

// String implements Event.
func (e *BaseEvent) String() string {
	return fmt.Sprintf("[id:%s kind:%s occured_at:%s data:%+v]",
		e.ID(),
		e.Kind(),
		e.OccuredAt(),
		e.Data(),
	)
}
