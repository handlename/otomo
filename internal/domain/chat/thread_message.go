//go:generate go run ../../../tools/gen-vo -file=thread_message.go
package chat

import (
	"fmt"

	"github.com/handlename/otomo/internal/domain/core"
)

// @vo
type ThreadMessageID struct {
	value string
}

// NewThreadMessageID creates a new ThreadMessageID with validation.
func NewThreadMessageID(value string) (ThreadMessageID, error) {
	if value == "" {
		return ThreadMessageID{}, fmt.Errorf("thread message ID cannot be empty")
	}
	return ThreadMessageID{value: value}, nil
}

// ThreadMessage is an entity representing a message within a Thread.
type ThreadMessage struct {
	id   ThreadMessageID
	user core.UserID
	body core.MessageBody
}

func NewThreadMessage(id ThreadMessageID, user core.UserID, body core.MessageBody) (*ThreadMessage, error) {
	if id.Value() == "" {
		return nil, fmt.Errorf("thread message ID is required")
	}
	if user.Value() == "" {
		return nil, fmt.Errorf("thread message user is required")
	}
	if body == "" {
		return nil, fmt.Errorf("thread message body is required")
	}

	return &ThreadMessage{
		id:   id,
		user: user,
		body: body,
	}, nil
}

func (t *ThreadMessage) String() string {
	return fmt.Sprintf("id=%s user=%s body=%q", t.ID(), t.User(), t.Body())
}

func (t *ThreadMessage) User() core.UserID {
	return t.user
}

func (t *ThreadMessage) Body() core.MessageBody {
	return t.body
}

func (t *ThreadMessage) ID() ThreadMessageID {
	return t.id
}
