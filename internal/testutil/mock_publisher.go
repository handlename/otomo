package testutil

import (
	"context"
	"errors"

	"github.com/handlename/otomo/internal/domain/event"
)

var _ (event.Publisher) = (*MockEventPublisher)(nil)

type MockEventPublisher struct {
	Handlers map[event.Kind][]event.Subscriber
	History  []event.Event
}

func NewMockPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		Handlers: make(map[event.Kind][]event.Subscriber),
		History:  make([]event.Event, 0),
	}
}

// Publish implements event.Publisher.
func (p *MockEventPublisher) Publish(ctx context.Context, event event.Event) error {
	p.History = append(p.History, event)
	handlers := p.Handlers[event.Kind()]

	errs := []error{}
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			errs = append(errs, err)
		}
	}

	if 0 < len(errs) {
		return errors.Join(errs...)
	}

	return nil
}

// Subscribe implements event.Publisher.
func (p *MockEventPublisher) Subscribe(kind event.Kind, handler event.Subscriber) {
	if _, exists := p.Handlers[kind]; !exists {
		p.Handlers[kind] = []event.Subscriber{}
	}
	p.Handlers[kind] = append(p.Handlers[kind], handler)
}
