package chat_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/stretchr/testify/assert"
)

func TestNewThreadMessage_Validation(t *testing.T) {
	tests := []struct {
		name      string
		id        chat.ThreadMessageID
		user      core.UserID
		body      core.MessageBody
		expectErr bool
	}{
		{
			name:      "valid thread message",
			id:        chat.ThreadMessageID("1"),
			user:      core.UserID("alice"),
			body:      core.MessageBody("hello"),
			expectErr: false,
		},
		{
			name:      "empty ID",
			id:        chat.ThreadMessageID(""),
			user:      core.UserID("alice"),
			body:      core.MessageBody("hello"),
			expectErr: true,
		},
		{
			name:      "empty user",
			id:        chat.ThreadMessageID("1"),
			user:      core.UserID(""),
			body:      core.MessageBody("hello"),
			expectErr: true,
		},
		{
			name:      "empty body",
			id:        chat.ThreadMessageID("1"),
			user:      core.UserID("alice"),
			body:      core.MessageBody(""),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chat.NewThreadMessage(tt.id, tt.user, tt.body)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.id, got.ID())
				assert.Equal(t, tt.user, got.User())
				assert.Equal(t, tt.body, got.Body())
			}
		})
	}
}
