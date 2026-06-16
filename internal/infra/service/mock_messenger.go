package service

import (
	"context"

	aservice "github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/rs/zerolog/log"
)

var _ aservice.Messenger = (*MockMessenger)(nil)

type UploadFileCall struct {
	ChannelID core.ChannelID
	ThreadID  chat.ThreadID
	Filename  string
	Content   string
}

// MockMessenger is a mock implementation of service.Messenger
type MockMessenger struct {
	PostMessageFunc func(ctx context.Context, channelID core.ChannelID, messageID core.MessageID, message chat.ReplyBody) error
	AddReactionFunc func(ctx context.Context, channelID core.ChannelID, messageID core.MessageID, emoji string) error
	FetchThreadFunc func(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID) (*chat.Thread, error)
	UploadFileFunc  func(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID, filename, content string) error

	History []struct {
		ChannelID core.ChannelID
		MessageID core.MessageID
		Message   chat.ReplyBody
	}
	ReactionHistory []struct {
		ChannelID core.ChannelID
		MessageID core.MessageID
		Emoji     string
	}
	UploadFileHistory []UploadFileCall
}

// FetchThread implements service.Messenger.
func (m *MockMessenger) FetchThread(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID) (*chat.Thread, error) {
	if m.FetchThreadFunc == nil {
		log.Warn().Msg("FetchThreadFunc is empty! you may set the func")
		t, _ := chat.NewThread(threadID)
		return t, nil
	}

	return m.FetchThreadFunc(ctx, channelID, threadID)
}

// PostMessage implements the Messenger interface
func (m *MockMessenger) PostMessage(ctx context.Context, channelID core.ChannelID, messageID core.MessageID, message chat.ReplyBody) error {
	m.History = append(m.History, struct {
		ChannelID core.ChannelID
		MessageID core.MessageID
		Message   chat.ReplyBody
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
func (m *MockMessenger) AddReaction(ctx context.Context, channelID core.ChannelID, messageID core.MessageID, emoji string) error {
	m.ReactionHistory = append(m.ReactionHistory, struct {
		ChannelID core.ChannelID
		MessageID core.MessageID
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
func (m *MockMessenger) UploadFile(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID, filename, content string) error {
	m.UploadFileHistory = append(m.UploadFileHistory, UploadFileCall{
		ChannelID: channelID,
		ThreadID:  threadID,
		Filename:  filename,
		Content:   content,
	})
	if m.UploadFileFunc != nil {
		return m.UploadFileFunc(ctx, channelID, threadID, filename, content)
	}
	return nil
}

