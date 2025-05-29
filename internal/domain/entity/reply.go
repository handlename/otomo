package entity

// Reply is a message sent to the User as a result of Otomo interpreting an Instruction.
type Reply interface {
	Body() string
}

var _ Reply = (*reply)(nil)

func NewReply(body string, attachments []string) Reply {
	return &reply{
		body:        body,
		attachments: attachments,
	}
}

type reply struct {
	body        string
	attachments []string // TODO
}

func (r *reply) Body() string { return r.body }
