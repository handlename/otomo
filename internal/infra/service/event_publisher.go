package service

import (
	"errors"
	"sync"

	"github.com/handlename/otomo/internal/domain/event"
)

var _ event.Publisher = (*EventPublisher)(nil)

type EventPublisher struct {
	handlers map[event.Kind][]event.Handler
	mutex    sync.RWMutex
}

func NewEventPublisher() *EventPublisher {
	return &EventPublisher{
		handlers: make(map[event.Kind][]event.Handler),
	}
}

// Publish implements event.Publisher.
func (p *EventPublisher) Publish(event event.Event) error {
	p.mutex.RLock()
	handlers := p.handlers[event.Kind()]
	p.mutex.Unlock()

	errs := []error{}
	for _, handler := range handlers {
		if err := handler(event); err != nil {
			errs = append(errs, err)
		}
	}

	if 0 < len(errs) {
		return errors.Join(errs...)
	}

	return nil
}

// Subscribe implements event.Publisher.
func (p *EventPublisher) Subscribe(kind event.Kind, handler event.Handler) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	if _, exists := p.handlers[kind]; !exists {
		p.handlers[kind] = []event.Handler{}
	}
	p.handlers[kind] = append(p.handlers[kind], handler)
}
