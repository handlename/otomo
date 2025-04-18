package service

import "context"

type Messenger interface {
	PostMessage(ctx context.Context, channelID, messageID, msg string) error
	AddReaction(ctx context.Context, channelID, messageID string, emoji string) error
}
