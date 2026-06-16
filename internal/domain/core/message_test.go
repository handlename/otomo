package core_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/core"
	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	tests := []struct {
		name    string
		role    core.MessageRole
		user    core.UserID
		body    core.MessageBody
		wantErr bool
	}{
		{
			name:    "valid user message",
			role:    core.RoleUser,
			user:    "U1234",
			body:    "hello",
			wantErr: false,
		},
		{
			name:    "valid system message without user",
			role:    core.RoleSystem,
			user:    "",
			body:    "system init",
			wantErr: false,
		},
		{
			name:    "invalid role",
			role:    core.MessageRole("invalid"),
			user:    "U1234",
			body:    "hello",
			wantErr: true,
		},
		{
			name:    "empty body",
			role:    core.RoleUser,
			user:    "U1234",
			body:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := core.NewMessage(tt.role, tt.user, tt.body)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, msg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, msg)
				assert.Equal(t, tt.role, msg.Role())
				assert.Equal(t, tt.user, msg.User())
				assert.Equal(t, tt.body, msg.Body())
			}
		})
	}
}
