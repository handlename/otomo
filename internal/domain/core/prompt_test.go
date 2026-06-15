package core_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Prompt_String(t *testing.T) {
	emptyPrompt, err := core.NewPrompt("", "", []*core.Prompt{})
	require.NoError(t, err)

	plain0_1, err := core.NewPlainPrompt("plain0-1")
	require.NoError(t, err)
	plain1_1, err := core.NewPlainPrompt("plain1-1")
	require.NoError(t, err)
	plain2_1, err := core.NewPlainPrompt("plain2-1")
	require.NoError(t, err)
	plain2_2, err := core.NewPlainPrompt("plain2-2")
	require.NoError(t, err)
	plain2_1_1, err := core.NewPlainPrompt("plain2-1-1")
	require.NoError(t, err)
	plain0_2, err := core.NewPlainPrompt("plain0-2")
	require.NoError(t, err)

	tag2_1, err := core.NewPrompt("tag2-1", "", []*core.Prompt{plain2_1_1})
	require.NoError(t, err)

	tag2, err := core.NewPrompt("tag2", "", []*core.Prompt{plain2_1, plain2_2, tag2_1})
	require.NoError(t, err)

	tag1, err := core.NewPrompt("tag1", "", []*core.Prompt{plain1_1})
	require.NoError(t, err)

	complexPrompt, err := core.NewPrompt("", "", []*core.Prompt{
		plain0_1,
		tag1,
		tag2,
		plain0_2,
	})
	require.NoError(t, err)

	tests := []struct {
		name     string
		prompt   *core.Prompt
		expected string
	}{
		{
			name:     "empty prompt",
			prompt:   emptyPrompt,
			expected: "",
		},
		{
			name:   "complex prompt",
			prompt: complexPrompt,
			expected: `plain0-1
<tag1>
plain1-1
</tag1>
<tag2>
plain2-1
plain2-2
<tag2-1>
plain2-1-1
</tag2-1>
</tag2>
plain0-2
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.prompt.String())
		})
	}
}
