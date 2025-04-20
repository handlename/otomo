package service

import (
	"context"
)

// MockMessenger is a mock implementation of service.Messenger
type MockMessenger struct {
	PostMessageFunc func(ctx context.Context, channelID, messageID, message string) error
	AddReactionFunc func(ctx context.Context, channelID, messageID string, emoji string) error
	History         []struct {
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
	
	return nil
}