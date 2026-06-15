package service

import (
	"context"

	"github.com/handlename/otomo/internal/domain/core"
)

// Publisher is an application service interface for publishing domain events.
type Publisher interface {
	Subscribe(core.EventKind, Subscriber)
	Publish(context.Context, core.Event) error
}

// Subscriber is a callback function that handles a published domain event.
type Subscriber func(context.Context, core.Event) error
