package valueobject_test

import (
	"testing"

	vo "github.com/handlename/otomo/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func Test_Prompt_String(t *testing.T) {
	tests := []struct {
		name     string
		prompt   vo.Prompt
		expected string
	}{
		{
			name:     "empty prompt",
			prompt:   vo.NewPromptWithChildren(nil, []vo.Prompt{}),
			expected: "",
		},
		{
			name: "complex prompt",
			prompt: vo.NewPromptWithChildren(
				nil,
				[]vo.Prompt{
					vo.NewPlainPrompt("plain0-1"),
					vo.NewPromptWithChildren(
						vo.NewTaggedPrompt("tag1"),
						[]vo.Prompt{
							vo.NewPlainPrompt("plain1-1"),
						},
					),
					vo.NewPromptWithChildren(
						vo.NewTaggedPrompt("tag2"),
						[]vo.Prompt{
							vo.NewPlainPrompt("plain2-1"),
							vo.NewPlainPrompt("plain2-2"),
							vo.NewPromptWithChildren(
								vo.NewTaggedPrompt("tag2-1"),
								[]vo.Prompt{
									vo.NewPlainPrompt("plain2-1-1"),
								},
							),
						},
					),
					vo.NewPlainPrompt("plain0-2"),
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
