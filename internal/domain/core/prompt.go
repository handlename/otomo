package core

import (
	"bytes"
	"fmt"
)

type PromptTag string

const (
	PromptTagSystem PromptTag = "system_instruction"
	PromptTagUser   PromptTag = "user_question"
)

// Prompt is a value object representing structured prompt tokens.
type Prompt struct {
	tag      PromptTag
	body     PromptBody
	children []*Prompt
}

// Tag returns the prompt tag.
func (p *Prompt) Tag() PromptTag {
	return p.tag
}

// String returns the string representation of the prompt.
func (p *Prompt) String() string {
	buf := bytes.NewBuffer([]byte{})
	if p.tag != "" {
		fmt.Fprintf(buf, "<%s>\n", p.tag)
	}

	if p.body != "" {
		fmt.Fprintln(buf, p.body)
	}

	for _, c := range p.children {
		if c != nil {
			fmt.Fprint(buf, c.String())
		}
	}

	if p.tag != "" {
		fmt.Fprintf(buf, "</%s>\n", p.tag)
	}

	return buf.String()
}

func NewPlainPrompt(body PromptBody) (*Prompt, error) {
	return NewPrompt("", body, nil)
}

func NewPrompt(tag PromptTag, body PromptBody, children []*Prompt) (*Prompt, error) {
	return &Prompt{
		tag:      tag,
		body:     body,
		children: children,
	}, nil
}
