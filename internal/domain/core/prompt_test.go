package core_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/core"
	"github.com/stretchr/testify/assert"
)

func Test_Prompt_String(t *testing.T) {
	tests := []struct {
		name     string
		prompt   core.Prompt
		expected string
	}{
		{
			name:     "empty prompt",
			prompt:   core.NewPrompt("", "", []core.Prompt{}),
			expected: "",
		},
		{
			name: "complex prompt",
			prompt: core.NewPrompt(
				"",
				"",
				[]core.Prompt{
					core.NewPlainPrompt("plain0-1"),
					core.NewPrompt(
						"tag1",
						"",
						[]core.Prompt{
							core.NewPlainPrompt("plain1-1"),
						},
					),
					core.NewPrompt(
						"tag2",
						"",
						[]core.Prompt{
							core.NewPlainPrompt("plain2-1"),
							core.NewPlainPrompt("plain2-2"),
							core.NewPrompt(
								"tag2-1",
								"",
								[]core.Prompt{
									core.NewPlainPrompt("plain2-1-1"),
								},
							),
						},
					),
					core.NewPlainPrompt("plain0-2"),
				},
			),
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
