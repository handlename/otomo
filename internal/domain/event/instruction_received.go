package event

import "time"

const KindInstructionReceived Kind = "instruction_received"

type InstructionReceivedData struct {
	MessageID      string
	ThreadID       string
	RawInstruction string
	SentAt         time.Time
}

type InstructionReceived struct {
	baseEvent
}

func NewInstructionReceived(data InstructionReceivedData) *InstructionReceived {
	return &InstructionReceived{
		baseEvent: newBaseEvent(KindInstructionReceived, data),
	}
}
