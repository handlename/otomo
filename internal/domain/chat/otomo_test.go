package chat_test

import (
	"context"
	"testing"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/stretchr/testify/require"
)

type dummyThinker struct{}

func (d *dummyThinker) Think(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
	return reasoning.NewAnswer("dummy reply")
}

func TestOtomo_Think(t *testing.T) {
	brainThinker := &dummyThinker{}
	brain, err := reasoning.NewBrain(brainThinker)
	require.NoError(t, err)

	o, err := chat.NewOtomo(brain)
	require.NoError(t, err)
	
	ctx := context.Background()
	c := reasoning.NewContext()
	c.SetUserPrompt("hello")

	reply, err := o.Think(ctx, c)
	require.NoError(t, err)

	if reply.Body() != "dummy reply" {
		t.Errorf("expected 'dummy reply', got %q", reply.Body())
	}
}
