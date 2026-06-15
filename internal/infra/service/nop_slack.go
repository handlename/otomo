package service

import (
	"context"

	"github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
)

var _ service.Messenger = (*NopSlack)(nil)

type NopSlack struct {
	Memory string
}

// FetchThread implements service.Messenger.
func (n *NopSlack) FetchThread(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID) (*chat.Thread, error) {
	panic("unimplemented")
}

// AddReaction implements service.Messenger.
func (n *NopSlack) AddReaction(ctx context.Context, channelID core.ChannelID, messageID core.MessageID, emoji string) error {
	panic("unimplemented")
}

// PostMessage implements service.Messenger.
func (n *NopSlack) PostMessage(ctx context.Context, channelID core.ChannelID, messageID core.MessageID, msg chat.ReplyBody) error {
	n.Memory += string(msg)
	return nil
}

// UploadFile implements service.Messenger.
func (n *NopSlack) UploadFile(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID, filename, content string) error {
	return nil
}

