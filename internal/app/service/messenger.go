package service

import (
	"context"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
)

type Messenger interface {
	PostMessage(ctx context.Context, channelID core.ChannelID, messageID core.MessageID, msg chat.ReplyBody) error
	AddReaction(ctx context.Context, channelID core.ChannelID, messageID core.MessageID, emoji string) error
	FetchThread(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID) (*chat.Thread, error)
	UploadFile(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID, filename string, content string) error
}
