package testutil

import (
	"context"
	"sync"

	appservice "github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/core"
)

var _ appservice.Publisher = (*MockEventPublisher)(nil)

type MockEventPublisher struct {
	mu          sync.RWMutex
	subscribers map[core.Kind][]appservice.Subscriber
	Published   []core.Event
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		subscribers: make(map[core.Kind][]appservice.Subscriber),
		Published:   []core.Event{},
	}
}

func (m *MockEventPublisher) Publish(ctx context.Context, ev core.Event) error {
	m.mu.Lock()
	m.Published = append(m.Published, ev)
	m.mu.Unlock()
	return nil
}

func (m *MockEventPublisher) Subscribe(kind core.Kind, sub appservice.Subscriber) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscribers[kind] = append(m.subscribers[kind], sub)
}
