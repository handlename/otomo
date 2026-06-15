package usecase

import (
	"context"
	"sync"

	appservice "github.com/handlename/otomo/internal/app/service"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/domain/reasoning"
)

var _ appservice.Publisher = (*mockEventPublisher)(nil)
var _ appservice.Messenger = (*mockMessenger)(nil)
var _ reasoning.BrainThinker = (*mockBrain)(nil)

type mockBrain struct {
	ThinkFunc func(context.Context, *reasoning.Context) (*reasoning.Answer, error)
}

func (m *mockBrain) Think(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
	if m.ThinkFunc != nil {
		return m.ThinkFunc(ctx, c)
	}
	return reasoning.NewAnswer(reasoning.AnswerBody("mock response"))
}

type mockEventPublisher struct {
	mu        sync.RWMutex
	Published []core.Event
}

func (m *mockEventPublisher) Publish(ctx context.Context, ev core.Event) error {
	m.mu.Lock()
	m.Published = append(m.Published, ev)
	m.mu.Unlock()
	return nil
}

func (m *mockEventPublisher) Subscribe(kind core.EventKind, sub appservice.Subscriber) {
	// No-op for tests
}

type mockMessenger struct {
	FetchThreadFunc func(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID) (*chat.Thread, error)
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
	UploadFileHistory []struct {
		ChannelID string
		ThreadTS  string
		Filename  string
		Content   string
	}
}

func (m *mockMessenger) PostMessage(ctx context.Context, channelID core.ChannelID, messageID core.MessageID, msg chat.ReplyBody) error {
	m.History = append(m.History, struct {
		ChannelID string
		MessageID string
		Message   string
	}{
		ChannelID: string(channelID),
		MessageID: string(messageID),
		Message:   string(msg),
	})
	return nil
}

func (m *mockMessenger) AddReaction(ctx context.Context, channelID core.ChannelID, messageID core.MessageID, emoji string) error {
	m.ReactionHistory = append(m.ReactionHistory, struct {
		ChannelID string
		MessageID string
		Emoji     string
	}{
		ChannelID: string(channelID),
		MessageID: string(messageID),
		Emoji:     emoji,
	})
	return nil
}

func (m *mockMessenger) FetchThread(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID) (*chat.Thread, error) {
	if m.FetchThreadFunc != nil {
		return m.FetchThreadFunc(ctx, channelID, threadID)
	}
	t, _ := chat.NewThread(threadID)
	return t, nil
}

func (m *mockMessenger) UploadFile(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID, filename string, content string) error {
	m.UploadFileHistory = append(m.UploadFileHistory, struct {
		ChannelID string
		ThreadTS  string
		Filename  string
		Content   string
	}{
		ChannelID: string(channelID),
		ThreadTS:  string(threadID),
		Filename:  filename,
		Content:   content,
	})
	return nil
}
