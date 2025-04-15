package slack

import (
	"context"

	"github.com/handlename/otomo/internal/infra/service"
)

type registry struct {
	Publisher *service.EventPublisher
	Slack     *service.Slack
}

type registryKey struct{}

func NewRegistry(ctx context.Context, publisher *service.EventPublisher, slack *service.Slack) *registry {
	return &registry{
		Publisher: publisher,
		Slack:     slack,
	}
}
