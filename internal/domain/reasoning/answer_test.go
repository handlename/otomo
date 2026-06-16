package reasoning_test

import (
	"testing"

	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/stretchr/testify/assert"
)

func TestNewAnswer(t *testing.T) {
	tests := []struct {
		name        string
		body        reasoning.AnswerBody
		expectError bool
	}{
		{
			name:        "valid body",
			body:        "This is a valid answer",
			expectError: false,
		},
		{
			name:        "empty body should return error",
			body:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := reasoning.NewAnswer(tt.body)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, ans)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, ans)
				assert.Equal(t, tt.body, ans.Body())
			}
		})
	}
}
