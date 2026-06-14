package core

// MessageRole represents the role of the speaker of a message (e.g. system, user, assistant).
type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

// Message is a value object that represents a single message in a chat history.
// It contains the role of the speaker, the user identifier, and the message body.
type Message struct {
	Role MessageRole
	User string
	Body string
}
