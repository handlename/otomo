package event

import (
	"time"

	"github.com/morikuni/failure/v2"
)

const KindInstructionReceived Kind = "instruction_received"

type InstructionReceivedData struct {
	ChannelID      string    `validate:"required,startswith=C"`
	MessageID      string    `validate:"required,numeric"`
	ThreadID       string    // empty if instruction is not in thread
	RawInstruction string    `validate:"required"`
	SentAt         time.Time `validate:"required"`
}

func (d InstructionReceivedData) Validate() error {
	return validate.Struct(d)
}

type InstructionReceived struct {
	baseEvent
}

func NewInstructionReceived(data InstructionReceivedData) (*InstructionReceived, error) {
	if err := data.Validate(); err != nil {
		return nil, failure.Wrap(err)
	}

	return &InstructionReceived{
		baseEvent: newBaseEvent(KindInstructionReceived, data),
	}, nil
}
