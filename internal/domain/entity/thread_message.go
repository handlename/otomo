package entity

type ThreadMessageID string

// ThreadMessage is a message in Thread
type ThreadMessage interface {
	ID() ThreadMessageID
	Body() string
}

func NewThreadMessage(id ThreadMessageID, body string) ThreadMessage {
	return &threadMessage{
		id:   id,
		body: body,
	}
}

type threadMessage struct {
	id   ThreadMessageID
	body string
}

// Body implements ThreadMessage.
func (t *threadMessage) Body() string {
	return t.body
}

// ID implements ThreadMessage.
func (t *threadMessage) ID() ThreadMessageID {
	return t.id
}
