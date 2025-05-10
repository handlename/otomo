package entity

import (
	"context"

	vo "github.com/handlename/otomo/internal/domain/valueobject"
	"github.com/samber/lo"
)

type Context interface {
	SetThread(Thread)
	SetSystemPrompt(string)
	SetUserPrompt(string)
	GetUserPrompt() vo.Prompt

	// Prompt returns all of data in context as valueobject.Prompt
	Prompt() vo.Prompt
}

type ContextRefresher func(context.Context, Context) error

func NewContext() Context {
	return &ct{
		systemPrompt: vo.NewPrompt("", "", []vo.Prompt{}),
		userPrompt:   vo.NewPrompt("", "", []vo.Prompt{}),
		thread:       NewThread(ThreadID("")),
	}
}

type ct struct {
	systemPrompt vo.Prompt
	userPrompt   vo.Prompt
	thread       Thread
}

// GetUserPrompt implements Context.
func (c *ct) GetUserPrompt() vo.Prompt {
	return c.userPrompt
}

// SetSystemPrompt implements Context.
func (c *ct) SetSystemPrompt(body string) {
	c.systemPrompt = vo.NewPrompt(vo.PromptTagSystem, body, []vo.Prompt{})
}

// SetUserPrompt implements Context.
func (c *ct) SetUserPrompt(body string) {
	c.userPrompt = vo.NewPrompt(vo.PromptTagUser, body, []vo.Prompt{})
}

// Prompt implements Context.
func (c *ct) Prompt() vo.Prompt {
	return vo.NewPrompt(
		"",
		"", // TODO
		[]vo.Prompt{
			c.systemPrompt,
			vo.NewPrompt("thread", "", lo.Map(c.thread.Messages(), func(msg ThreadMessage, _ int) vo.Prompt {
				return vo.NewPrompt("message", msg.Body(), nil)
			})),
			c.userPrompt,
		},
	)
}

// SetThread implements Context.
func (c *ct) SetThread(thread Thread) {
	c.thread = thread
}
