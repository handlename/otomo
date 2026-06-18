package core_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/core"
	"github.com/samber/lo"
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
			user:    lo.Must(core.NewUserID("U1234")),
			body:    "hello",
			wantErr: false,
		},
		{
			name:    "valid system message without user",
			role:    core.RoleSystem,
			user:    core.UserID{},
			body:    "system init",
			wantErr: false,
		},
		{
			name:    "invalid role",
			role:    core.MessageRole("invalid"),
			user:    lo.Must(core.NewUserID("U1234")),
			body:    "hello",
			wantErr: true,
		},
		{
			name:    "empty body",
			role:    core.RoleUser,
			user:    lo.Must(core.NewUserID("U1234")),
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

func TestNewUserID(t *testing.T) {
	_, err := core.NewUserID("")
	assert.Error(t, err)

	uid, err := core.NewUserID("U12345")
	assert.NoError(t, err)
	assert.Equal(t, "U12345", uid.Value())
}

func TestNewChannelID(t *testing.T) {
	_, err := core.NewChannelID("")
	assert.Error(t, err)

	_, err = core.NewChannelID("invalid")
	assert.Error(t, err)

	cid, err := core.NewChannelID("C12345")
	assert.NoError(t, err)
	assert.Equal(t, "C12345", cid.Value())
}

func TestNewMessageID(t *testing.T) {
	_, err := core.NewMessageID("")
	assert.Error(t, err)

	_, err = core.NewMessageID("abc.123")
	assert.Error(t, err)

	mid, err := core.NewMessageID("123.456")
	assert.NoError(t, err)
	assert.Equal(t, "123.456", mid.Value())
}
