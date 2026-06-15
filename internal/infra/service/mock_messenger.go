package service

import (
	"context"

	aservice "github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/rs/zerolog/log"
)

var _ aservice.Messenger = (*MockMessenger)(nil)

type UploadFileCall struct {
	ChannelID string
	ThreadTS  string
	Filename  string
	Content   string
}

// MockMessenger is a mock implementation of service.Messenger
type MockMessenger struct {
	PostMessageFunc func(ctx context.Context, channelID, messageID, message string) error
	AddReactionFunc func(ctx context.Context, channelID, messageID string, emoji string) error
	FetchThreadFunc func(ctx context.Context, channelID string, threadID string) (*chat.Thread, error)
	UploadFileFunc  func(ctx context.Context, channelID, threadTS, filename, content string) error

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
	UploadFileHistory []UploadFileCall
}

// FetchThread implements service.Messenger.
func (m *MockMessenger) FetchThread(ctx context.Context, channelID string, threadID string) (*chat.Thread, error) {
	if m.FetchThreadFunc == nil {
		log.Warn().Msg("FetchThreadFunc is empty! you may set the func")
		return chat.NewThread(""), nil
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

// UploadFile implements service.Messenger.
func (m *MockMessenger) UploadFile(ctx context.Context, channelID, threadTS, filename, content string) error {
	m.UploadFileHistory = append(m.UploadFileHistory, UploadFileCall{
		ChannelID: channelID,
		ThreadTS:  threadTS,
		Filename:  filename,
		Content:   content,
	})
	if m.UploadFileFunc != nil {
		return m.UploadFileFunc(ctx, channelID, threadTS, filename, content)
	}
	return nil
}

