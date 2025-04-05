package slack

import (
	"context"

	"github.com/handlename/otomo/internal/infra/service"
)

type registry struct {
	Slack *service.Slack
}

type registryKey struct{}

func NewRegistry(ctx context.Context, slack *service.Slack) *registry {
	return &registry{
		Slack: slack,
	}
}
