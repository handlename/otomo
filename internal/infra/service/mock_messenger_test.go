package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/samber/lo"
)

func TestMockMessenger_UploadFile(t *testing.T) {
	mock := &service.MockMessenger{}

	ctx := context.Background()
	err := mock.UploadFile(ctx, lo.Must(core.NewChannelID("Cchan-id")), chat.ThreadID("ts-1234"), "test.txt", "test content")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mock.UploadFileHistory) != 1 {
		t.Fatalf("expected 1 upload call, got %d", len(mock.UploadFileHistory))
	}

	call := mock.UploadFileHistory[0]
	if call.ChannelID != lo.Must(core.NewChannelID("Cchan-id")) {
		t.Errorf("expected ChannelID to be 'chan-id', got '%s'", call.ChannelID)
	}
	if call.ThreadID != chat.ThreadID("ts-1234") {
		t.Errorf("expected ThreadID to be 'ts-1234', got '%s'", call.ThreadID)
	}
	if call.Filename != "test.txt" {
		t.Errorf("expected Filename to be 'test.txt', got '%s'", call.Filename)
	}
	if call.Content != "test content" {
		t.Errorf("expected Content to be 'test content', got '%s'", call.Content)
	}
}

func TestMockMessenger_UploadFileFunc(t *testing.T) {
	customErr := errors.New("custom upload error")
	var called bool
	mock := &service.MockMessenger{
		UploadFileFunc: func(ctx context.Context, channelID core.ChannelID, threadID chat.ThreadID, filename, content string) error {
			called = true
			if channelID != lo.Must(core.NewChannelID("Cchan-id")) || threadID != chat.ThreadID("ts-1234") || filename != "test.txt" || content != "test content" {
				t.Errorf("unexpected parameters in mock func")
			}
			return customErr
		},
	}

	ctx := context.Background()
	err := mock.UploadFile(ctx, lo.Must(core.NewChannelID("Cchan-id")), chat.ThreadID("ts-1234"), "test.txt", "test content")
	if !errors.Is(err, customErr) {
		t.Fatalf("expected custom error, got %v", err)
	}

	if !called {
		t.Error("expected UploadFileFunc to be called")
	}
}

