package entity

// Reply is a message sent to the User as a result of Otomo interpreting an Instruction.
type Reply struct {
	body        string
	attachments []string // TODO
}

func (r *Reply) Body() string { return r.body }

func NewReply(body string, attachments []string) *Reply {
	return &Reply{
		body:        body,
		attachments: attachments,
	}
}
