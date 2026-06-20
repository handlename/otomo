package reasoning

import "context"

// Tool represents a capability that the bot can run to interact with external resources.
type Tool interface {
	Name() ToolName
	Description() string
	InputSchema() string // JSON Schema definition of the input parameters
	Execute(ctx context.Context, inputJSON string) (string, error)
}
