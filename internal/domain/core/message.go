package core

import "fmt"

// MessageRole represents the role of the speaker of a message (e.g. system, user, assistant).
type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

// Message is a value object that represents a single message in a chat history.
type Message struct {
	role MessageRole
	user string
	body string
}

// NewMessage creates a new Message with validation.
func NewMessage(role MessageRole, user string, body string) (*Message, error) {
	if role != RoleSystem && role != RoleUser && role != RoleAssistant {
		return nil, fmt.Errorf("invalid message role: %s", role)
	}
	if body == "" {
		return nil, fmt.Errorf("message body cannot be empty")
	}
	return &Message{
		role: role,
		user: user,
		body: body,
	}, nil
}

func (m *Message) Role() MessageRole { return m.role }
func (m *Message) User() string      { return m.user }
func (m *Message) Body() string      { return m.body }
