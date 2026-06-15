package chat

import (
	"time"

	"github.com/handlename/otomo/internal/domain/core"
	"github.com/morikuni/failure/v2"
)

const KindInstructionReceived core.Kind = "instruction_received"

type InstructionReceivedData struct {
	ChannelID      string    `validate:"required,startswith=C"`
	MessageID      string    `validate:"required,numeric"`
	ThreadID       string    
	RawInstruction string    `validate:"required"`
	SentAt         time.Time `validate:"required"`
}

func (d InstructionReceivedData) Validate() error {
	return validate.Struct(d)
}

// InstructionReceived is a value object event dispatched when an instruction is received.
type InstructionReceived struct {
	core.BaseEvent
}

func NewInstructionReceived(data InstructionReceivedData) (*InstructionReceived, error) {
	if err := data.Validate(); err != nil {
		return nil, failure.Wrap(err)
	}

	base, err := core.NewBaseEvent(KindInstructionReceived, data)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	return &InstructionReceived{
		BaseEvent: base,
	}, nil
}
