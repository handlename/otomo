package usecase

import (
	"context"
	"fmt"

	appservice "github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
)

type ReplyToUser struct {
	messenger appservice.Messenger
	tools     []reasoning.Tool
}

func NewReplyToUser(messenger appservice.Messenger, tools []reasoning.Tool) *ReplyToUser {
	return &ReplyToUser{
		messenger: messenger,
		tools:     tools,
	}
}

func (u *ReplyToUser) Run(ctx context.Context, otomo *chat.Otomo, channelID core.ChannelID, userPrompt core.PromptBody) error {
	turns := 0

	c := reasoning.NewContext()
	c.SetUserPrompt(userPrompt)
	c.SetTools(u.tools)

	for {
		if err := ctx.Err(); err != nil {
			return failure.Wrap(err)
		}

		if turns >= reasoning.MaxToolTurns {
			return failure.New(errorcode.ErrInternal, failure.Message("too many tool execution turns"))
		}
		turns++

		ans, err := otomo.Think(ctx, c)
		if err != nil {
			return failure.Wrap(err,
				failure.WithCode(errorcode.ErrInternal),
				failure.Message("failed to think"),
			)
		}

		if !ans.HasToolCalls() {
			reply, err := chat.NewReply(chat.ReplyBody(ans.Body()), []chat.Attachment{})
			if err != nil {
				return failure.Wrap(err, failure.WithCode(errorcode.ErrInternal))
			}
			if err := u.messenger.PostMessage(ctx, channelID, core.MessageID{}, reply.Body()); err != nil {
				return failure.Wrap(err,
					failure.WithCode(errorcode.ErrInternal),
					failure.Message("failed to send reply"),
				)
			}
			return nil
		}

		var results []reasoning.ToolResult
		for _, tc := range ans.ToolCalls() {
			tool, ok := u.findTool(tc.Name())
			if !ok {
				tr, err := reasoning.NewToolResult(
					tc.ID(),
					fmt.Sprintf("error: tool '%s' not found", tc.Name().Value()),
					reasoning.IsError(true),
				)
				if err != nil {
					return failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create tool result"))
				}
				results = append(results, tr)
				continue
			}

			out, err := tool.Execute(ctx, tc.InputJSON())
			if err != nil {
				tr, err := reasoning.NewToolResult(
					tc.ID(),
					fmt.Sprintf("error executing tool: %v", err),
					reasoning.IsError(true),
				)
				if err != nil {
					return failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create tool result"))
				}
				results = append(results, tr)
			} else {
				tr, err := reasoning.NewToolResult(
					tc.ID(),
					out,
					reasoning.IsError(false),
				)
				if err != nil {
					return failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create tool result"))
				}
				results = append(results, tr)
			}
		}

		if err := c.AddToolUseResponse(string(ans.Body()), ans.ToolCalls()); err != nil {
			return failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to update context with tool calls"))
		}
		if err := c.AddToolResults(results); err != nil {
			return failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to update context with tool results"))
		}
	}
}

func (u *ReplyToUser) findTool(name reasoning.ToolName) (reasoning.Tool, bool) {
	for _, t := range u.tools {
		if t.Name().Equals(name) {
			return t, true
		}
	}
	return nil, false
}
