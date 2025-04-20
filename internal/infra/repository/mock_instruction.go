package repository

import (
	"context"

	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/domain/event"
	"github.com/handlename/otomo/internal/domain/repository"
)

var _ repository.Instruction = (*MockInstructionRepository)(nil)

// MockInstructionRepository is a mock implementation of repository.Instruction
type MockInstructionRepository struct {
	NewFromInstructionReceivedDataFunc func(ctx context.Context, data event.InstructionReceivedData) *entity.Instruction
}

// NewFromInstructionReceivedData implements repository.Instruction.
func (r *MockInstructionRepository) NewFromInstructionReceivedData(ctx context.Context, data event.InstructionReceivedData) *entity.Instruction {
	if r.NewFromInstructionReceivedDataFunc != nil {
		return r.NewFromInstructionReceivedDataFunc(ctx, data)
	}

	inst := &entity.Instruction{}
	return inst
}
