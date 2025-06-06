package service

import (
	"context"

	"github.com/handlename/otomo/internal/domain/entity"
)

// MockAnswer implements entity.Answer interface
type MockAnswer struct {
	body string
}

// Body returns the body of the response
func (r *MockAnswer) Body() string {
	return r.body
}

var _ entity.Brain = (*MockBrain)(nil)

// MockBrain is a mock implementation for entity.Brain
type MockBrain struct {
	ThinkFunc func(ctx context.Context, c entity.Context) (*entity.Answer, error)
}

// Think mocks the Brain's Think method
func (b *MockBrain) Think(ctx context.Context, ectx entity.Context) (*entity.Answer, error) {
	if b.ThinkFunc != nil {
		return b.ThinkFunc(ctx, ectx)
	}

	return entity.NewAnswer("mock response"), nil
}
