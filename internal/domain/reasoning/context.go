package reasoning

import (
	"fmt"
	"slices"

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
	return slices.Clone(m.toolCalls)
}

func (m *ContextMessage) ToolResults() []ToolResult {
	return slices.Clone(m.toolResults)
}

// NewContextMessage creates a new ContextMessage after validating its arguments.
func NewContextMessage(
	role string,
	user core.UserID,
	content string,
	toolCalls []ToolCall,
	toolResults []ToolResult,
) (*ContextMessage, error) {
	r := core.MessageRole(role)
	switch r {
	case core.RoleSystem:
		if len(toolCalls) > 0 || len(toolResults) > 0 {
			return nil, fmt.Errorf("system message cannot contain tool calls or tool results")
		}
	case core.RoleUser:
		if len(toolCalls) > 0 {
			return nil, fmt.Errorf("user message cannot contain tool calls")
		}
	case core.RoleAssistant:
		if len(toolResults) > 0 {
			return nil, fmt.Errorf("assistant message cannot contain tool results")
		}
	default:
		return nil, fmt.Errorf("invalid message role: %s", role)
	}

	if content == "" && len(toolCalls) == 0 && len(toolResults) == 0 {
		return nil, fmt.Errorf("content cannot be empty unless tool calls or tool results are present")
	}

	return &ContextMessage{
		role:        r,
		user:        user,
		content:     content,
		toolCalls:   slices.Clone(toolCalls),
		toolResults: slices.Clone(toolResults),
	}, nil
}

// ToolResult represents the execution output of a tool call.
type ToolResult struct {
	toolUseID ToolCallID
	output    string
	isError   bool
}

func NewToolResult(toolUseID ToolCallID, output string, isError bool) (ToolResult, error) {
	if toolUseID.Value() == "" {
		return ToolResult{}, fmt.Errorf("tool use ID cannot be empty")
	}
	return ToolResult{
		toolUseID: toolUseID,
		output:    output,
		isError:   isError,
	}, nil
}

func (tr ToolResult) ToolUseID() ToolCallID {
	return tr.toolUseID
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

func (c *Context) SetMessages(messages []*core.Message) error {
	var msgs []*ContextMessage
	for _, msg := range messages {
		if msg != nil {
			cm, err := NewContextMessage(string(msg.Role()), msg.User(), string(msg.Body()), nil, nil)
			if err != nil {
				return err
			}
			msgs = append(msgs, cm)
		}
	}
	c.messages = msgs
	return nil
}

func (c *Context) Messages() []*ContextMessage {
	return slices.Clone(c.messages)
}

func (c *Context) Tools() []Tool {
	return slices.Clone(c.tools)
}

func (c *Context) SetTools(tools []Tool) {
	c.tools = slices.Clone(tools)
}

func (c *Context) AddToolUseResponse(content string, toolCalls []ToolCall) error {
	msg, err := NewContextMessage(string(core.RoleAssistant), core.UserID{}, content, toolCalls, nil)
	if err != nil {
		return err
	}
	c.messages = append(c.messages, msg)
	return nil
}

func (c *Context) AddToolResults(results []ToolResult) error {
	msg, err := NewContextMessage(string(core.RoleUser), core.UserID{}, "", nil, results)
	if err != nil {
		return err
	}
	c.messages = append(c.messages, msg)
	return nil
}

// Prompt returns the prompt representation of the context.
// Warning: This method does not serialize tool calls or tool results in the legacy XML string.
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
