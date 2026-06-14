package core

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ID string
type Kind string

// Event is an entity interface representing something that has occurred in the domain.
type Event interface {
	ID() ID
	Kind() Kind
	OccuredAt() time.Time
	Data() any
	String() string
}

var _ Event = (*BaseEvent)(nil)

// BaseEvent is a value object helper that provides standard fields for Event implementations.
type BaseEvent struct {
	id      ID
	kind    Kind
	ts      time.Time
	payload any
}

func NewBaseEvent(kind Kind, payload any) BaseEvent {
	return BaseEvent{
		id:      ID(uuid.New().String()),
		kind:    kind,
		ts:      time.Now(),
		payload: payload,
	}
}

// Data implements Event.
func (e *BaseEvent) Data() any {
	return e.payload
}

// ID implements Event.
func (e *BaseEvent) ID() ID {
	return e.id
}

// Kind implements Event.
func (e *BaseEvent) Kind() Kind {
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
