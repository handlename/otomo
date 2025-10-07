package entity

import "context"

type Tool interface {
	// Name returns the name of the tool.
	Name() string

	// Description returns a description of the tool.
	Description() string

	// Execute executes the tool with the given parameters.
	Execute(context.Context, ToolParams) (*ToolAnswer, error)
}
