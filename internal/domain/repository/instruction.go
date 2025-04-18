package repository

import (
	"context"

	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/domain/event"
)

type Instruction interface {
	NewFromInstructionReceivedData(context.Context, event.InstructionReceivedData) *entity.Instruction
}
