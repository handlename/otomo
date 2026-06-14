package reasoning

import (
	"fmt"

	"github.com/handlename/otomo/internal/domain/core"
	"github.com/samber/lo"
)

// Context is an entity that accumulates necessary information (prompts, history) for reasoning.
type Context struct {
	systemPrompt *core.Prompt
	userPrompt   *core.Prompt
	messages     []core.Message
}

func NewContext() *Context {
	return &Context{
		systemPrompt: core.NewPrompt("", "", []*core.Prompt{}),
		userPrompt:   core.NewPrompt("", "", []*core.Prompt{}),
		messages:     []core.Message{},
	}
}

func (c *Context) GetUserPrompt() *core.Prompt {
	return c.userPrompt
}

func (c *Context) SetSystemPrompt(body string) {
	c.systemPrompt = core.NewPrompt(core.PromptTagSystem, body, []*core.Prompt{})
}

func (c *Context) SetUserPrompt(body string) {
	c.userPrompt = core.NewPrompt(core.PromptTagUser, body, []*core.Prompt{})
}

func (c *Context) SetMessages(messages []core.Message) {
	c.messages = messages
}

func (c *Context) Prompt() *core.Prompt {
	return core.NewPrompt(
		"",
		"",
		[]*core.Prompt{
			c.systemPrompt,
			core.NewPrompt("thread", "", lo.Map(c.messages, func(msg core.Message, _ int) *core.Prompt {
				var tag core.PromptTag
				if msg.User != "" {
					tag = core.PromptTag(fmt.Sprintf("message user=%s", msg.User))
				} else {
					tag = core.PromptTag(fmt.Sprintf("message role=%s", msg.Role))
				}
				return core.NewPrompt(tag, msg.Body, nil)
			})),
			c.userPrompt,
		},
	)
}
