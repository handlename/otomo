//go:generate go run ../../../tools/gen-vo -file=types.go
package core

import (
	"fmt"
	"strings"
)

// @vo
type UserID struct {
	value string
}

// NewUserID creates a new UserID with validation.
func NewUserID(value string) (UserID, error) {
	if value == "" {
		return UserID{}, fmt.Errorf("user ID cannot be empty")
	}
	return UserID{value: value}, nil
}

type MessageBody string
type PromptBody string

// @vo
type ChannelID struct {
	value string
}

// NewChannelID creates a new ChannelID with validation.
func NewChannelID(value string) (ChannelID, error) {
	if value == "" {
		return ChannelID{}, fmt.Errorf("channel ID cannot be empty")
	}
	if !strings.HasPrefix(value, "C") {
		return ChannelID{}, fmt.Errorf("channel ID must start with 'C': %s", value)
	}
	return ChannelID{value: value}, nil
}

// @vo
type MessageID struct {
	value string
}

// NewMessageID creates a new MessageID with validation.
func NewMessageID(value string) (MessageID, error) {
	if value == "" {
		return MessageID{}, fmt.Errorf("message ID cannot be empty")
	}
	for _, r := range value {
		if (r < '0' || r > '9') && r != '.' {
			return MessageID{}, fmt.Errorf("message ID must contain only digits or periods: %s", value)
		}
	}
	return MessageID{value: value}, nil
}
