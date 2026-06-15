package chat_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/stretchr/testify/assert"
)

func TestNewReply_Validation(t *testing.T) {
	tests := []struct {
		name        string
		body        chat.ReplyBody
		attachments []chat.Attachment
		expectErr   bool
	}{
		{
			name:        "valid reply",
			body:        chat.ReplyBody("hello"),
			attachments: []chat.Attachment{},
			expectErr:   false,
		},
		{
			name:        "empty body",
			body:        chat.ReplyBody(""),
			attachments: []chat.Attachment{},
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chat.NewReply(tt.body, tt.attachments)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.body, got.Body())
				assert.Equal(t, tt.attachments, got.Attachments())
			}
		})
	}
}
