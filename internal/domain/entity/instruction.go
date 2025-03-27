package entity

import vo "github.com/handlename/otomo/internal/domain/valueobject"

// Instruction represents a text message containing instructions sent from `user` to `otomo` via Slack
type Instruction struct {
	id   vo.InstructionID
	body string
}

func (i *Instruction) ID() vo.InstructionID { return i.id }
func (i *Instruction) Body() string         { return i.body }

func NewInstruction(id vo.InstructionID, body string) *Instruction {
	return &Instruction{
		id:   id,
		body: body,
	}
}
