package tool

import (
	"context"

	"github.com/handlename/otomo/internal/domain/entity"
)

// Mock is a mock implementation of Tool for testing
type Mock struct {
	name        string
	description string
}

// NewMock creates a new Mock tool with the given name and description
func NewMock(name, description string) *Mock {
	return &Mock{
		name:        name,
		description: description,
	}
}

// Name returns the tool name
func (m *Mock) Name() string {
	return m.name
}

// Execute executes the mock tool (returns nil for testing purposes)
func (m *Mock) Execute(ctx context.Context, params entity.ToolParams) (*entity.ToolAnswer, error) {
	return nil, nil
}

// Description returns the tool description
func (m *Mock) Description() string {
	return m.description
}
