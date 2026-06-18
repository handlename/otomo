package tool_test

import (
	"context"
	"testing"

	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/infra/tool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDummyTool(t *testing.T) {
	dt := tool.NewDummyTool()
	
	// Test Name matching custom types
	expectedName, err := reasoning.NewToolName("dummy_tool")
	require.NoError(t, err)
	assert.Equal(t, expectedName, dt.Name())
	assert.Contains(t, dt.Description(), "verification")
	assert.Contains(t, dt.InputSchema(), "properties")

	ctx := context.Background()
	
	// Test success path (characters count)
	out, err := dt.Execute(ctx, `{"text":"hello"}`)
	require.NoError(t, err)
	assert.Equal(t, `{"length": 5}`, out)

	// Test multi-byte string character count
	out, err = dt.Execute(ctx, `{"text":"おとも"}`)
	require.NoError(t, err)
	assert.Equal(t, `{"length": 3}`, out)

	// Test invalid JSON syntax
	_, err = dt.Execute(ctx, `{"invalid"`)
	assert.Error(t, err)
}
