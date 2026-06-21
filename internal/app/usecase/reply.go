package usecase

import (
	"context"
	"fmt"

	"github.com/handlename/otomo/config"
	appservice "github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

type ReplyInput struct {
	EventData *chat.InstructionReceivedData
}

type ReplyOutput struct{}

type Reply struct {
	otomo *chat.Otomo
	slack appservice.Messenger
	tools []reasoning.Tool
}

func NewReply(otomo *chat.Otomo, slack appservice.Messenger, tools []reasoning.Tool) *Reply {
	return &Reply{
		otomo: otomo,
		slack: slack,
		tools: tools,
	}
}

func (r *Reply) Run(ctx context.Context, input ReplyInput) (*ReplyOutput, error) {
	c := reasoning.NewContext()
	c.SetUserPrompt(core.PromptBody(input.EventData.RawInstruction()))
	c.SetTools(r.tools)

	if input.EventData.ThreadID().Value() != input.EventData.MessageID().Value() {
		thread, err := r.slack.FetchThread(ctx, input.EventData.ChannelID(), input.EventData.ThreadID())
		if err != nil {
			r.handleError(ctx, input.EventData, err)
			return nil, failure.Wrap(err)
		}

		msgs := make([]*core.Message, len(thread.Messages()))
		for i, m := range thread.Messages() {
			role := core.RoleUser
			// For now, treat all history as User unless it is from the bot
			// Optional: mapping logic can be refined later if needed.
			msg, err := core.NewMessage(role, m.User(), m.Body())
			if err != nil {
				r.handleError(ctx, input.EventData, err)
				return nil, failure.Wrap(err)
			}
			msgs[i] = msg
		}
		if err := c.SetMessages(msgs); err != nil {
			r.handleError(ctx, input.EventData, err)
			return nil, failure.Wrap(err)
		}
	}

	turns := 0
	for {
		if err := ctx.Err(); err != nil {
			r.handleError(ctx, input.EventData, err)
			return nil, failure.Wrap(err)
		}

		if !reasoning.ShouldContinueToUseTool(turns) {
			err := failure.New(errorcode.ErrInternal, failure.Message("too many tool execution turns"))
			r.handleError(ctx, input.EventData, err)
			return nil, err
		}
		turns++

		ans, err := r.otomo.Think(ctx, c)
		if err != nil {
			wrappedErr := failure.Wrap(err,
				failure.WithCode(errorcode.ErrInternal),
				failure.Message("failed to think"),
			)
			r.handleError(ctx, input.EventData, wrappedErr)
			return nil, wrappedErr
		}

		if !ans.HasToolCalls() {
			reply, err := chat.NewReply(chat.ReplyBody(ans.Body()), []chat.Attachment{})
			if err != nil {
				wrappedErr := failure.Wrap(err, failure.WithCode(errorcode.ErrInternal))
				r.handleError(ctx, input.EventData, wrappedErr)
				return nil, wrappedErr
			}
			err = r.slack.PostMessage(ctx, input.EventData.ChannelID(), input.EventData.MessageID(), reply.Body())
			if err != nil {
				wrappedErr := failure.Wrap(err,
					failure.WithCode(errorcode.ErrInternal),
					failure.Message("failed to send reply"),
				)
				r.handleError(ctx, input.EventData, wrappedErr)
				return nil, wrappedErr
			}
			return &ReplyOutput{}, nil
		}

		var results []reasoning.ToolResult
		for _, tc := range ans.ToolCalls() {
			tool, ok := r.findTool(tc.Name())
			if !ok {
				tr, err := reasoning.NewToolResult(
					tc.ID(),
					fmt.Sprintf("error: tool '%s' not found", tc.Name().Value()),
					reasoning.ToolResultError,
				)
				if err != nil {
					wrappedErr := failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create tool result"))
					r.handleError(ctx, input.EventData, wrappedErr)
					return nil, wrappedErr
				}
				results = append(results, tr)
				continue
			}

			out, err := tool.Execute(ctx, tc.InputJSON())
			if err != nil {
				tr, err := reasoning.NewToolResult(
					tc.ID(),
					fmt.Sprintf("error executing tool: %v", err),
					reasoning.ToolResultError,
				)
				if err != nil {
					wrappedErr := failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create tool result"))
					r.handleError(ctx, input.EventData, wrappedErr)
					return nil, wrappedErr
				}
				results = append(results, tr)
			} else {
				tr, err := reasoning.NewToolResult(
					tc.ID(),
					out,
					reasoning.ToolResultSuccess,
				)
				if err != nil {
					wrappedErr := failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create tool result"))
					r.handleError(ctx, input.EventData, wrappedErr)
					return nil, wrappedErr
				}
				results = append(results, tr)
			}
		}

		if err := c.AddToolUseResponse(string(ans.Body()), ans.ToolCalls()); err != nil {
			wrappedErr := failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to update context with tool calls"))
			r.handleError(ctx, input.EventData, wrappedErr)
			return nil, wrappedErr
		}
		if err := c.AddToolResults(results); err != nil {
			wrappedErr := failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to update context with tool results"))
			r.handleError(ctx, input.EventData, wrappedErr)
			return nil, wrappedErr
		}
	}
}

func (r *Reply) findTool(name reasoning.ToolName) (reasoning.Tool, bool) {
	for _, t := range r.tools {
		if t.Name().Equals(name) {
			return t, true
		}
	}
	return nil, false
}

func (r *Reply) handleError(ctx context.Context, data *chat.InstructionReceivedData, targetErr error) {
	cfg := config.Config.Slack.ErrorFeedback

	threadTS := data.ThreadID()

	if cfg.GetEnableReaction() {
		emoji := cfg.GetReactionEmoji()
		if err := r.slack.AddReaction(ctx, data.ChannelID(), data.MessageID(), emoji); err != nil {
			log.Error().Err(err).Msg("failed to add error reaction emoji")
		}
	}

	if cfg.GetEnablePostSnippet() {
		errContent := targetErr.Error()
		if err := r.slack.UploadFile(ctx, data.ChannelID(), threadTS, "error.txt", errContent); err != nil {
			log.Error().Err(err).Msg("failed to upload error snippet")
		}
	}
}

func (r *Reply) Subscribe(publisher appservice.Publisher) {
	publisher.Subscribe(chat.KindInstructionReceived, func(ctx context.Context, eev core.Event) error {
		ev, ok := eev.(*chat.InstructionReceived)
		if !ok {
			log.Error().Msg("failed to assert event")
			return nil
		}

		input := ReplyInput{
			EventData: ev.Data().(*chat.InstructionReceivedData),
		}
		_, err := r.Run(ctx, input)
		return err
	})
}
