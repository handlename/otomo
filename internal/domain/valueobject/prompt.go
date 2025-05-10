package valueobject

import (
	"bytes"
	"fmt"
)

// type Prompt string

type PromptTag string

type Prompt interface {
	Tag() PromptTag
	IsRoot() bool
	String() string
}

type prompt struct {
	root     bool
	tag      PromptTag
	body     string
	children []Prompt
}

// Tag implements Prompt.
func (p *prompt) Tag() PromptTag {
	return p.tag
}

// IsRoot implements Prompt.
func (p *prompt) IsRoot() bool {
	return p.root
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

func NewPromptWithChildren(parent Prompt, children []Prompt) Prompt {
	if parent == nil {
		parent = NewPrompt(nil, "", "")
	}

	return &prompt{
		root:     parent.IsRoot(),
		tag:      parent.Tag(),
		children: children,
	}
}

func NewTaggedPrompt(parent Prompt, tag PromptTag) Prompt {
	return NewPrompt(parent, tag, "")
}

func NewPlainPrompt(parent Prompt, body string) Prompt {
	return NewPrompt(parent, "", body)
}

func NewPrompt(parent Prompt, tag PromptTag, body string) Prompt {
	if parent == nil {
		parent = &prompt{
			root: true,
		}
	}

	return &prompt{
		root: parent.IsRoot(),
		tag:  tag,
		body: body,
	}
}
