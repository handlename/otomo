package chat

import (
	"fmt"
	"slices"
)

type ReplyBody string
type Attachment string

// Reply is a value object sent to the User as a result of Otomo interpreting an Instruction.
type Reply struct {
	body        ReplyBody
	attachments []Attachment
}

func (r *Reply) Body() ReplyBody { return r.body }

func (r *Reply) Attachments() []Attachment {
	return slices.Clone(r.attachments)
}

func NewReply(body ReplyBody, attachments []Attachment) (*Reply, error) {
	if body == "" {
		return nil, fmt.Errorf("reply body is required")
	}
	return &Reply{
		body:        body,
		attachments: attachments,
	}, nil
}
