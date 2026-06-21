package tool_test

import (
	"testing"

	"github.com/handlename/otomo"
	"github.com/handlename/otomo/internal/infra/tool"
	"github.com/stretchr/testify/assert"
)

func TestUserAgent(t *testing.T) {
	expected := "otomo-bot/" + otomo.Version
	assert.Equal(t, expected, tool.UserAgent())
}
