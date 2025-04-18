package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/handlename/otomo/config"
	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/domain/event"
	"github.com/handlename/otomo/internal/domain/repository"
	vo "github.com/handlename/otomo/internal/domain/valueobject"
)

var _ repository.Instruction = (*SlackInstruction)(nil)

type SlackInstruction struct{}

func NewSlackInstruction() *SlackInstruction {
	return &SlackInstruction{}
}

// NewFromInstructionReceivedData implements repository.Instruction.
func (s *SlackInstruction) NewFromInstructionReceivedData(_ context.Context, data event.InstructionReceivedData) *entity.Instruction {
	body := data.RawInstruction
	body = strings.TrimSpace(body)
	body = strings.TrimPrefix(body, fmt.Sprintf("<%s>", config.Config.Slack.BotUserID))

	return entity.NewInstruction(vo.InstructionID(data.ThreadID), body)
}
