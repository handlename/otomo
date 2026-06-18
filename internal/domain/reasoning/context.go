package reasoning

import (
	"fmt"

	"github.com/handlename/otomo/internal/domain/core"
)

// ContextMessage represents a single message in the reasoning context.
type ContextMessage struct {
	role        core.MessageRole
	user        core.UserID
	content     string
	toolCalls   []ToolCall
	toolResults []ToolResult
}

func (m *ContextMessage) Role() string {
	return string(m.role)
}

func (m *ContextMessage) User() core.UserID {
	return m.user
}

func (m *ContextMessage) Content() string {
	return m.content
}

func (m *ContextMessage) ToolCalls() []ToolCall {
	return m.toolCalls
}

func (m *ContextMessage) ToolResults() []ToolResult {
	return m.toolResults
}

// ToolResult represents the execution output of a tool call.
type ToolResult struct {
	toolCallID ToolCallID
	output     string
	isError    bool
}

func NewToolResult(toolCallID ToolCallID, output string, isError bool) ToolResult {
	return ToolResult{
		toolCallID: toolCallID,
		output:     output,
		isError:    isError,
	}
}

func (tr ToolResult) ToolCallID() ToolCallID {
	return tr.toolCallID
}

func (tr ToolResult) Output() string {
	return tr.output
}

func (tr ToolResult) IsError() bool {
	return tr.isError
}

// Context is an entity that accumulates necessary information (prompts, history) for reasoning.
type Context struct {
	systemPrompt     *core.Prompt
	systemPromptBody string
	userPrompt       *core.Prompt
	messages         []*ContextMessage
	tools            []Tool
}

func NewContext() *Context {
	systemPrompt, _ := core.NewPrompt("", "", []*core.Prompt{})
	userPrompt, _ := core.NewPrompt("", "", []*core.Prompt{})
	return &Context{
		systemPrompt: systemPrompt,
		userPrompt:   userPrompt,
		messages:     []*ContextMessage{},
		tools:        []Tool{},
	}
}

func (c *Context) GetUserPrompt() *core.Prompt {
	return c.userPrompt
}

func (c *Context) SetSystemPrompt(body core.PromptBody) {
	c.systemPromptBody = string(body)
	c.systemPrompt, _ = core.NewPrompt(core.PromptTagSystem, body, []*core.Prompt{})
}

func (c *Context) SystemPromptBody() string {
	return c.systemPromptBody
}

func (c *Context) SetUserPrompt(body core.PromptBody) {
	c.userPrompt, _ = core.NewPrompt(core.PromptTagUser, body, []*core.Prompt{})
}

func (c *Context) SetMessages(messages []*core.Message) {
	var msgs []*ContextMessage
	for _, msg := range messages {
		if msg != nil {
			msgs = append(msgs, &ContextMessage{
				role:    msg.Role(),
				user:    msg.User(),
				content: string(msg.Body()),
			})
		}
	}
	c.messages = msgs
}

func (c *Context) Messages() []*ContextMessage {
	return c.messages
}

func (c *Context) Tools() []Tool {
	return c.tools
}

func (c *Context) SetTools(tools []Tool) {
	c.tools = tools
}

func (c *Context) AddToolUseResponse(content string, toolCalls []ToolCall) {
	c.messages = append(c.messages, &ContextMessage{
		role:      core.RoleAssistant,
		content:   content,
		toolCalls: toolCalls,
	})
}

func (c *Context) AddToolResults(results []ToolResult) {
	c.messages = append(c.messages, &ContextMessage{
		role:        core.RoleUser,
		toolResults: results,
	})
}

func (c *Context) Prompt() *core.Prompt {
	var prompts []*core.Prompt
	for _, msg := range c.messages {
		if msg == nil {
			continue
		}
		var tag core.PromptTag
		if msg.User().Value() != "" {
			tag = core.PromptTag(fmt.Sprintf("message user=%s", msg.User().Value()))
		} else {
			tag = core.PromptTag(fmt.Sprintf("message role=%s", string(msg.Role())))
		}
		p, _ := core.NewPrompt(tag, core.PromptBody(msg.Content()), nil)
		prompts = append(prompts, p)
	}

	threadPrompt, _ := core.NewPrompt("thread", "", prompts)

	prompt, _ := core.NewPrompt(
		"",
		"",
		[]*core.Prompt{
			c.systemPrompt,
			threadPrompt,
			c.userPrompt,
		},
	)
	return prompt
}
