package service

import (
	"context"

	"github.com/handlename/otomo/internal/app/service"
)

var _ service.Messenger = (*NopSlack)(nil)

type NopSlack struct {
	Memory string
}

// AddReaction implements service.Messenger.
func (n *NopSlack) AddReaction(ctx context.Context, channelID string, messageID string, emoji string) error {
	panic("unimplemented")
}

// PostMessage implements service.Messenger.
func (n *NopSlack) PostMessage(ctx context.Context, channelID string, messageID string, msg string) error {
	n.Memory += msg
	return nil
}
