package chat

import (
	"fmt"
	"time"

	"github.com/handlename/otomo/internal/domain/core"
	"github.com/morikuni/failure/v2"
)

const KindInstructionReceived core.EventKind = "instruction_received"

type RawInstruction string

// InstructionReceivedData is a value object containing the data for the InstructionReceived event.
type InstructionReceivedData struct {
	channelID      core.ChannelID
	messageID      core.MessageID
	threadID       ThreadID
	rawInstruction RawInstruction
	sentAt         time.Time
}

// NewInstructionReceivedData creates and validates InstructionReceivedData using inline validation.
func NewInstructionReceivedData(channelID core.ChannelID, messageID core.MessageID, threadID ThreadID, rawInstruction RawInstruction, sentAt time.Time) (*InstructionReceivedData, error) {
	if channelID.Value() == "" {
		return nil, fmt.Errorf("channel ID is required")
	}
	if messageID.Value() == "" {
		return nil, fmt.Errorf("message ID is required")
	}
	if threadID.Value() == "" {
		return nil, fmt.Errorf("thread ID is required")
	}
	if rawInstruction == "" {
		return nil, fmt.Errorf("raw instruction is required")
	}
	if sentAt.IsZero() {
		return nil, fmt.Errorf("sent at timestamp is required")
	}

	return &InstructionReceivedData{
		channelID:      channelID,
		messageID:      messageID,
		threadID:       threadID,
		rawInstruction: rawInstruction,
		sentAt:         sentAt,
	}, nil
}

func (d *InstructionReceivedData) ChannelID() core.ChannelID      { return d.channelID }
func (d *InstructionReceivedData) MessageID() core.MessageID      { return d.messageID }
func (d *InstructionReceivedData) ThreadID() ThreadID             { return d.threadID }
func (d *InstructionReceivedData) RawInstruction() RawInstruction { return d.rawInstruction }
func (d *InstructionReceivedData) SentAt() time.Time              { return d.sentAt }

// InstructionReceived is a value object event dispatched when an instruction is received.
type InstructionReceived struct {
	core.BaseEvent
}

func NewInstructionReceived(data *InstructionReceivedData) (*InstructionReceived, error) {
	if data == nil {
		return nil, fmt.Errorf("instruction received data is required")
	}

	base, err := core.NewBaseEvent(KindInstructionReceived, data)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	return &InstructionReceived{
		BaseEvent: base,
	}, nil
}
