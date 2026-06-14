package reasoning

import (
	"context"
	"fmt"

	"github.com/handlename/otomo/internal/domain/core"
	"github.com/samber/lo"
)

type ContextRefresher func(context.Context, Context) error

// Context is an entity that accumulates necessary information (prompts, history) for reasoning.
type Context interface {
	SetMessages([]core.Message)
	SetSystemPrompt(string)
	SetUserPrompt(string)
	GetUserPrompt() core.Prompt
	Prompt() core.Prompt
}

func NewContext() Context {
	return &ct{
		systemPrompt: core.NewPrompt("", "", []core.Prompt{}),
		userPrompt:   core.NewPrompt("", "", []core.Prompt{}),
		messages:     []core.Message{},
	}
}

type ct struct {
	systemPrompt core.Prompt
	userPrompt   core.Prompt
	messages     []core.Message
}

func (c *ct) GetUserPrompt() core.Prompt {
	return c.userPrompt
}

func (c *ct) SetSystemPrompt(body string) {
	c.systemPrompt = core.NewPrompt(core.PromptTagSystem, body, []core.Prompt{})
}

func (c *ct) SetUserPrompt(body string) {
	c.userPrompt = core.NewPrompt(core.PromptTagUser, body, []core.Prompt{})
}

func (c *ct) SetMessages(messages []core.Message) {
	c.messages = messages
}

func (c *ct) Prompt() core.Prompt {
	return core.NewPrompt(
		"",
		"",
		[]core.Prompt{
			c.systemPrompt,
			core.NewPrompt("thread", "", lo.Map(c.messages, func(msg core.Message, _ int) core.Prompt {
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
