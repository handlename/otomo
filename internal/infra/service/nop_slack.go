package service

import (
	"context"

	"github.com/handlename/otomo/internal/app/service"
)

var _ service.Messenger = (*NopSlack)(nil)

type NopSlack struct {
	Memory string
}

// Send implements service.Messenger.
func (n *NopSlack) Send(_ context.Context, msg string) error {
	n.Memory += msg
	return nil
}
