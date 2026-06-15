package chat

import "fmt"

type ThreadMessageID string

// ThreadMessage is an entity representing a message within a Thread.
type ThreadMessage struct {
	id   ThreadMessageID
	user string
	body string
}

func NewThreadMessage(id ThreadMessageID, user, body string) (*ThreadMessage, error) {
	if id == "" {
		return nil, fmt.Errorf("thread message ID is required")
	}
	if user == "" {
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

func (t *ThreadMessage) User() string {
	return t.user
}

func (t *ThreadMessage) Body() string {
	return t.body
}

func (t *ThreadMessage) ID() ThreadMessageID {
	return t.id
}
