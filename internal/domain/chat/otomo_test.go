package chat_test

import (
	"context"
	"testing"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/reasoning"
)

type dummyThinker struct{}

func (d *dummyThinker) Think(ctx context.Context, c *reasoning.Context) (*reasoning.Answer, error) {
	return reasoning.NewAnswer("dummy reply"), nil
}

func TestOtomo_Think(t *testing.T) {
	brain := reasoning.NewBrain(&dummyThinker{})
	o := chat.NewOtomo(brain)
	
	ctx := context.Background()
	c := reasoning.NewContext()
	c.SetUserPrompt("hello")

	reply, err := o.Think(ctx, c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reply.Body() != "dummy reply" {
		t.Errorf("expected 'dummy reply', got %q", reply.Body())
	}
}
