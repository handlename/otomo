package entity_test

import (
	"context"
	"testing"

	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/infra/brain"
	"github.com/handlename/otomo/internal/infra/tool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBrain_SelectTool(t *testing.T) {
	tests := []struct {
		name          string
		tools         []entity.Tool
		userPrompt    string
		expectedTool  entity.Tool
		expectedError bool
	}{
		{
			name:          "no tools available",
			tools:         []entity.Tool{},
			userPrompt:    "test prompt",
			expectedTool:  nil,
			expectedError: false,
		},
		{
			name: "tools available but no match",
			tools: []entity.Tool{
				tool.NewMock("tool1", "description for tool1"),
				tool.NewMock("tool2", "description for tool2"),
			},
			userPrompt:    "test prompt",
			expectedTool:  nil,
			expectedError: false,
		},
	}

	// Test case with successful tool selection using mock thinker
	t.Run("successful tool selection", func(t *testing.T) {
		tool1 := tool.NewMock("web", "web search tool")
		tool2 := tool.NewMock("calc", "calculator tool")

		// Mock thinker that returns "web" as the selected tool
		mockThinker := &brain.Mock{
			ThinkFunc: func(ctx context.Context, c entity.Context) (*entity.Answer, error) {
				return entity.NewAnswer("web"), nil
			},
		}
		b := entity.NewBrain(mockThinker)

		// Add tools to brain
		ctx := context.Background()
		err := b.AddTool(ctx, tool1)
		require.NoError(t, err)
		err = b.AddTool(ctx, tool2)
		require.NoError(t, err)

		// Create context with user prompt
		c := entity.NewContext()
		c.SetUserPrompt("search something")

		// Execute SelectTool
		result, err := b.SelectTool(ctx, c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, tool1, result)
	})

	// Test case with empty response from thinker
	t.Run("empty response from thinker", func(t *testing.T) {
		tool1 := tool.NewMock("web", "web search tool")
		tool2 := tool.NewMock("calc", "calculator tool")

		// Mock thinker that returns empty string (no tool selected)
		mockThinker := &brain.Mock{
			ThinkFunc: func(ctx context.Context, c entity.Context) (*entity.Answer, error) {
				return entity.NewAnswer(""), nil
			},
		}
		b := entity.NewBrain(mockThinker)

		// Add tools to brain
		ctx := context.Background()
		err := b.AddTool(ctx, tool1)
		require.NoError(t, err)
		err = b.AddTool(ctx, tool2)
		require.NoError(t, err)

		// Create context with user prompt
		c := entity.NewContext()
		c.SetUserPrompt("unclear request")

		// Execute SelectTool
		result, err := b.SelectTool(ctx, c)

		// Assertions
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	// Test case with response containing whitespace
	t.Run("response with whitespace trimming", func(t *testing.T) {
		tool1 := tool.NewMock("web", "web search tool")
		tool2 := tool.NewMock("calc", "calculator tool")

		// Mock thinker that returns tool name with surrounding whitespace
		mockThinker := &brain.Mock{
			ThinkFunc: func(ctx context.Context, c entity.Context) (*entity.Answer, error) {
				return entity.NewAnswer("  web  "), nil
			},
		}
		b := entity.NewBrain(mockThinker)

		// Add tools to brain
		ctx := context.Background()
		err := b.AddTool(ctx, tool1)
		require.NoError(t, err)
		err = b.AddTool(ctx, tool2)
		require.NoError(t, err)

		// Create context with user prompt
		c := entity.NewContext()
		c.SetUserPrompt("search something")

		// Execute SelectTool
		result, err := b.SelectTool(ctx, c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, tool1, result)
	})

	// Test case with non-existent tool name response
	t.Run("non-existent tool name response", func(t *testing.T) {
		tool1 := tool.NewMock("web", "web search tool")
		tool2 := tool.NewMock("calc", "calculator tool")

		// Mock thinker that returns a tool name that doesn't exist
		mockThinker := &brain.Mock{
			ThinkFunc: func(ctx context.Context, c entity.Context) (*entity.Answer, error) {
				return entity.NewAnswer("nonexistent"), nil
			},
		}
		b := entity.NewBrain(mockThinker)

		// Add tools to brain
		ctx := context.Background()
		err := b.AddTool(ctx, tool1)
		require.NoError(t, err)
		err = b.AddTool(ctx, tool2)
		require.NoError(t, err)

		// Create context with user prompt
		c := entity.NewContext()
		c.SetUserPrompt("some request")

		// Execute SelectTool
		result, err := b.SelectTool(ctx, c)

		// Assertions
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create brain with Straw thinker
			thinker := &brain.Straw{}
			b := entity.NewBrain(thinker)

			// Add tools to brain
			ctx := context.Background()
			for _, tool := range tt.tools {
				err := b.AddTool(ctx, tool)
				require.NoError(t, err)
			}

			// Create context with user prompt
			c := entity.NewContext()
			c.SetUserPrompt(tt.userPrompt)

			// Execute SelectTool
			result, err := b.SelectTool(ctx, c)

			// Assertions
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedTool, result)
		})
	}
}
