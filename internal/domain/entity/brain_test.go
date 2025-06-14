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
		name         string
		setupThinker func() entity.BrainThinker
		tools        []entity.Tool
		expectedTool entity.Tool
		wantErr      bool
		errBody      string
	}{
		{
			name: "no tools available",
			setupThinker: func() entity.BrainThinker {
				return &brain.Straw{}
			},
			tools:        []entity.Tool{},
			expectedTool: nil,
		},
		{
			name: "tools available but no match",
			setupThinker: func() entity.BrainThinker {
				return &brain.Straw{}
			},
			tools: []entity.Tool{
				tool.NewMock("tool1", "description for tool1"),
				tool.NewMock("tool2", "description for tool2"),
			},
			expectedTool: nil,
		},
		{
			name: "successful tool selection",
			setupThinker: func() entity.BrainThinker {
				return &brain.Mock{
					ThinkFunc: func(ctx context.Context, c entity.Context) (*entity.Answer, error) {
						return entity.NewAnswer("web"), nil
					},
				}
			},
			tools: []entity.Tool{
				tool.NewMock("web", "web search tool"),
				tool.NewMock("calc", "calculator tool"),
			},
			expectedTool: tool.NewMock("web", "web search tool"),
		},
		{
			name: "empty response from thinker",
			setupThinker: func() entity.BrainThinker {
				return &brain.Mock{
					ThinkFunc: func(ctx context.Context, c entity.Context) (*entity.Answer, error) {
						return entity.NewAnswer(""), nil
					},
				}
			},
			tools: []entity.Tool{
				tool.NewMock("web", "web search tool"),
				tool.NewMock("calc", "calculator tool"),
			},
			expectedTool: nil,
		},
		{
			name: "response with whitespace trimming",
			setupThinker: func() entity.BrainThinker {
				return &brain.Mock{
					ThinkFunc: func(ctx context.Context, c entity.Context) (*entity.Answer, error) {
						return entity.NewAnswer("  web  "), nil
					},
				}
			},
			tools: []entity.Tool{
				tool.NewMock("web", "web search tool"),
				tool.NewMock("calc", "calculator tool"),
			},
			expectedTool: tool.NewMock("web", "web search tool"),
		},
		{
			name: "non-existent tool name response",
			setupThinker: func() entity.BrainThinker {
				return &brain.Mock{
					ThinkFunc: func(ctx context.Context, c entity.Context) (*entity.Answer, error) {
						return entity.NewAnswer("nonexistent"), nil
					},
				}
			},
			tools: []entity.Tool{
				tool.NewMock("web", "web search tool"),
				tool.NewMock("calc", "calculator tool"),
			},
			expectedTool: nil,
		},
		{
			name: "thinker returns error",
			setupThinker: func() entity.BrainThinker {
				return &brain.Mock{
					ThinkFunc: func(ctx context.Context, c entity.Context) (*entity.Answer, error) {
						return nil, assert.AnError
					},
				}
			},
			tools: []entity.Tool{
				tool.NewMock("web", "web search tool"),
				tool.NewMock("calc", "calculator tool"),
			},
			wantErr:      true,
			errBody:      "general error for testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			// Create brain with the configured thinker
			thinker := tt.setupThinker()
			b := entity.NewBrain(thinker)

			// Add tools to brain
			for _, tool := range tt.tools {
				err := b.AddTool(ctx, tool)
				require.NoError(t, err)
			}

			// Execute
			c := entity.NewContext()
			c.SetUserPrompt("dummy")
			got, err := b.SelectTool(ctx, c)

			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errBody)
				return
			}

			assert.NoError(t, err)

			// For tool comparison, we need to compare by name since tools are created separately
			if tt.expectedTool == nil {
				assert.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			assert.Equal(t, tt.expectedTool.Name(), got.Name())
		})
	}
}
