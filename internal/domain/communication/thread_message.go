package communication

import "fmt"

type ThreadMessageID string

// ThreadMessage is an entity representing a message within a Thread.
type ThreadMessage interface {
	ID() ThreadMessageID
	User() string
	Body() string
	String() string
}

type threadMessage struct {
	id   ThreadMessageID
	user string
	body string
}

func NewThreadMessage(id ThreadMessageID, user, body string) ThreadMessage {
	return &threadMessage{
		id:   id,
		user: user,
		body: body,
	}
}

func (t *threadMessage) String() string {
	return fmt.Sprintf("id=%s user=%s body=%q", t.ID(), t.User(), t.Body())
}

func (t *threadMessage) User() string {
	return t.user
}

func (t *threadMessage) Body() string {
	return t.body
}

func (t *threadMessage) ID() ThreadMessageID {
	return t.id
}
