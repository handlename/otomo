package service

import (
	"context"
	"errors"
	"sync"

	"github.com/handlename/otomo/internal/domain/event"
)

var _ event.Publisher = (*EventPublisher)(nil)

type EventPublisher struct {
	handlers map[event.Kind][]event.Subscriber
	mutex    sync.RWMutex
}

func NewEventPublisher() *EventPublisher {
	return &EventPublisher{
		handlers: make(map[event.Kind][]event.Subscriber),
	}
}

// Publish implements event.Publisher.
func (p *EventPublisher) Publish(ctx context.Context, event event.Event) error {
	p.mutex.RLock()
	handlers := p.handlers[event.Kind()]
	p.mutex.RUnlock()

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
func (p *EventPublisher) Subscribe(kind event.Kind, handler event.Subscriber) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, exists := p.handlers[kind]; !exists {
		p.handlers[kind] = []event.Subscriber{}
	}
	p.handlers[kind] = append(p.handlers[kind], handler)
}
