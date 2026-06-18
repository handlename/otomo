package tool_test

import (
	"context"
	"testing"

	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/handlename/otomo/internal/infra/tool"
	"github.com/morikuni/failure/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDummyTool_Metadata(t *testing.T) {
	dt := tool.NewDummyTool()

	expectedName, err := reasoning.NewToolName("dummy_tool")
	require.NoError(t, err)
	assert.Equal(t, expectedName, dt.Name())
	assert.Contains(t, dt.Description(), "verification")
	assert.Contains(t, dt.InputSchema(), "properties")
}

func TestDummyTool_Execute(t *testing.T) {
	tests := []struct {
		name        string
		inputJSON   string
		expectedOut string
		expectErr   bool
		errCode     errorcode.ErrorCode
		errMsg      string
	}{
		{
			name:        "success english text",
			inputJSON:   `{"text":"hello"}`,
			expectedOut: `{"length": 5}`,
			expectErr:   false,
		},
		{
			name:        "success multi-byte text",
			inputJSON:   `{"text":"おとも"}`,
			expectedOut: `{"length": 3}`,
			expectErr:   false,
		},
		{
			name:      "error invalid json syntax",
			inputJSON: `{"invalid"`,
			expectErr: true,
			errCode:   errorcode.ErrInvalidArgument,
			errMsg:    "failed to unmarshal inputs",
		},
		{
			name:      "error invalid parameter type",
			inputJSON: `{"text": 123}`,
			expectErr: true,
			errCode:   errorcode.ErrInvalidArgument,
			errMsg:    "failed to unmarshal inputs",
		},
		{
			name:      "error missing text parameter",
			inputJSON: `{}`,
			expectErr: true,
			errCode:   errorcode.ErrInvalidArgument,
			errMsg:    "text is required",
		},
	}

	dt := tool.NewDummyTool()
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := dt.Execute(ctx, tt.inputJSON)
			if tt.expectErr {
				assert.Error(t, err)
				if tt.errCode != "" {
					assert.True(t, failure.Is(err, tt.errCode), "expected error code %v, got error %v", tt.errCode, err)
				}
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOut, out)
			}
		})
	}
}

