package service

import (
	"context"

	aservice "github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/rs/zerolog/log"
)

var _ aservice.Messenger = (*MockMessenger)(nil)

// MockMessenger is a mock implementation of service.Messenger
type MockMessenger struct {
	PostMessageFunc func(ctx context.Context, channelID, messageID, message string) error
	AddReactionFunc func(ctx context.Context, channelID, messageID string, emoji string) error
	FetchThreadFunc func(ctx context.Context, channelID string, threadID string) (entity.Thread, error)

	History []struct {
		ChannelID string
		MessageID string
		Message   string
	}
	ReactionHistory []struct {
		ChannelID string
		MessageID string
		Emoji     string
	}
}

// FetchThread implements service.Messenger.
func (m *MockMessenger) FetchThread(ctx context.Context, channelID string, threadID string) (entity.Thread, error) {
	if m.FetchThreadFunc == nil {
		log.Warn().Msg("FetchThreadFunc is empty! you may set the func")
		return entity.NewThread(""), nil
	}

	return m.FetchThreadFunc(ctx, channelID, threadID)
}

// PostMessage implements the Messenger interface
func (m *MockMessenger) PostMessage(ctx context.Context, channelID, messageID, message string) error {
	m.History = append(m.History, struct {
		ChannelID string
		MessageID string
		Message   string
	}{
		ChannelID: channelID,
		MessageID: messageID,
		Message:   message,
	})

	if m.PostMessageFunc != nil {
		return m.PostMessageFunc(ctx, channelID, messageID, message)
	}
	log.Warn().Msg("PostMessageFunc is empty! you may set the func")

	return nil
}

// AddReaction implements the Messenger interface
func (m *MockMessenger) AddReaction(ctx context.Context, channelID, messageID string, emoji string) error {
	m.ReactionHistory = append(m.ReactionHistory, struct {
		ChannelID string
		MessageID string
		Emoji     string
	}{
		ChannelID: channelID,
		MessageID: messageID,
		Emoji:     emoji,
	})

	if m.AddReactionFunc != nil {
		return m.AddReactionFunc(ctx, channelID, messageID, emoji)
	}
	log.Warn().Msg("AddReactionFunc is empty! you may set the func")

	return nil
}
