package valueobject

import (
	"bytes"
	"fmt"
)

// type Prompt string

type PromptTag string

type Prompt interface {
	Tag() PromptTag
	String() string
}

type prompt struct {
	tag      PromptTag
	body     string
	children []Prompt
}

// Tag implements Prompt.
func (p *prompt) Tag() PromptTag {
	return p.tag
}

// String implements Prompt.
func (p *prompt) String() string {
	buf := bytes.NewBuffer([]byte{})
	if p.tag != "" {
		fmt.Fprintf(buf, "<%s>\n", p.tag)
	}

	if p.body != "" {
		fmt.Fprintln(buf, p.body)
	}

	for _, c := range p.children {
		fmt.Fprint(buf, c.String())
	}

	if p.tag != "" {
		fmt.Fprintf(buf, "</%s>\n", p.tag)
	}

	return buf.String()
}

func NewPlainPrompt(body string) Prompt {
	return NewPrompt("", body, nil)
}

func NewPrompt(tag PromptTag, body string, children []Prompt) Prompt {
	return &prompt{
		tag:      tag,
		body:     body,
		children: children,
	}
}
