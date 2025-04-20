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

// MockBrain is a mock implementation for entity.Brain
type MockBrain struct {
	ThinkFunc func(ctx context.Context, ectx entity.Context, instruction *entity.Instruction) (*entity.Answer, error)
}

// Think mocks the Brain's Think method
func (b *MockBrain) Think(ctx context.Context, ectx entity.Context, instruction *entity.Instruction) (*entity.Answer, error) {
	if b.ThinkFunc != nil {
		return b.ThinkFunc(ctx, ectx, instruction)
	}
	
	return entity.NewAnswer("mock response"), nil
}