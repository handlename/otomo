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
	c := reasoning.NewContext()
	c.SetUserPrompt(userPrompt)
	c.SetTools(u.tools)

	for {
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
				tr, _ := reasoning.NewToolResult(
					tc.ID(),
					fmt.Sprintf("error: tool '%s' not found", tc.Name().Value()),
					true,
				)
				results = append(results, tr)
				continue
			}

			out, err := tool.Execute(ctx, tc.InputJSON())
			if err != nil {
				tr, _ := reasoning.NewToolResult(
					tc.ID(),
					fmt.Sprintf("error executing tool: %v", err),
					true,
				)
				results = append(results, tr)
			} else {
				tr, _ := reasoning.NewToolResult(
					tc.ID(),
					out,
					false,
				)
				results = append(results, tr)
			}
		}

		c.AddToolUseResponse(string(ans.Body()), ans.ToolCalls())
		c.AddToolResults(results)
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
