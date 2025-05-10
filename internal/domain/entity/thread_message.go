package entity

import "fmt"

type ThreadMessageID string

// ThreadMessage is a message in Thread
type ThreadMessage interface {
	ID() ThreadMessageID

	// User returns user name who post the message.
	User() string

	// Body is the content of the message.
	Body() string

	String() string
}

func NewThreadMessage(id ThreadMessageID, user, body string) ThreadMessage {
	return &threadMessage{
		id:   id,
		user: user,
		body: body,
	}
}

type threadMessage struct {
	id   ThreadMessageID
	user string
	body string
}

// String implements ThreadMessage.
func (t *threadMessage) String() string {
	return fmt.Sprintf("id=%s user=%s body=%q", t.ID(), t.User(), t.Body())
}

// User implements ThreadMessage.
func (t *threadMessage) User() string {
	return t.user
}

// Body implements ThreadMessage.
func (t *threadMessage) Body() string {
	return t.body
}

// ID implements ThreadMessage.
func (t *threadMessage) ID() ThreadMessageID {
	return t.id
}
