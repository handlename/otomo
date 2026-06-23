package usecase

import (
	"context"

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

	ans, err := ExecuteToolLoop(ctx, otomo, c, u.tools)
	if err != nil {
		return err
	}

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
