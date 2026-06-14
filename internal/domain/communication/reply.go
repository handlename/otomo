package communication

// Reply is a value object sent to the User as a result of Otomo interpreting an Instruction.
type Reply interface {
	Body() string
}

type reply struct {
	body        string
	attachments []string
}

func (r *reply) Body() string { return r.body }

func NewReply(body string, attachments []string) Reply {
	return &reply{
		body:        body,
		attachments: attachments,
	}
}
