package chat

import (
	"fmt"
	"strings"
	"time"

	"github.com/handlename/otomo/internal/domain/core"
	"github.com/morikuni/failure/v2"
)

const KindInstructionReceived core.Kind = "instruction_received"

// InstructionReceivedData is a value object containing the data for the InstructionReceived event.
type InstructionReceivedData struct {
	channelID      string
	messageID      string
	threadID       string
	rawInstruction string
	sentAt         time.Time
}

// NewInstructionReceivedData creates and validates InstructionReceivedData using inline validation.
func NewInstructionReceivedData(channelID, messageID, threadID, rawInstruction string, sentAt time.Time) (*InstructionReceivedData, error) {
	if channelID == "" || !strings.HasPrefix(channelID, "C") {
		return nil, fmt.Errorf("channel ID is required and must start with 'C'")
	}
	if messageID == "" {
		return nil, fmt.Errorf("message ID is required")
	}
	for _, char := range messageID {
		if (char < '0' || char > '9') && char != '.' {
			return nil, fmt.Errorf("message ID must be numeric")
		}
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

func (d *InstructionReceivedData) ChannelID() string      { return d.channelID }
func (d *InstructionReceivedData) MessageID() string      { return d.messageID }
func (d *InstructionReceivedData) ThreadID() string       { return d.threadID }
func (d *InstructionReceivedData) RawInstruction() string { return d.rawInstruction }
func (d *InstructionReceivedData) SentAt() time.Time      { return d.sentAt }

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
