package service

import (
	"context"

	"github.com/handlename/otomo/internal/domain/entity"
)

type Messenger interface {
	PostMessage(ctx context.Context, channelID, messageID, msg string) error
	AddReaction(ctx context.Context, channelID, messageID string, emoji string) error
	FetchThread(ctx context.Context, channelID, threadID string) (entity.Thread, error)
}
