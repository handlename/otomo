package service

import (
	"context"
	"errors"
	"sync"

	appservice "github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/core"
)

var _ appservice.Publisher = (*EventPublisher)(nil)

type EventPublisher struct {
	subscribers map[core.EventKind][]appservice.Subscriber
	mu          sync.RWMutex
}

func NewEventPublisher() *EventPublisher {
	return &EventPublisher{
		subscribers: make(map[core.EventKind][]appservice.Subscriber),
	}
}

func (p *EventPublisher) Publish(ctx context.Context, ev core.Event) error {
	p.mu.RLock()
	subs, ok := p.subscribers[ev.Kind()]
	p.mu.RUnlock()
	if !ok {
		return nil
	}
	var errs []error
	for _, sub := range subs {
		if err := sub(ctx, ev); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (p *EventPublisher) Subscribe(kind core.EventKind, sub appservice.Subscriber) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.subscribers[kind] = append(p.subscribers[kind], sub)
}
