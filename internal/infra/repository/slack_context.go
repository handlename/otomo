package repository

import (
	"context"
	"slices"

	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/domain/repository"
	"github.com/morikuni/failure/v2"
)

var _ repository.Context = (*SlackContext)(nil)

type SlackContext struct{}

// New implements repository.Context.
func (s SlackContext) New() entity.Context {
	return &slackContext{}
}

// slackContext is a implementation of entity.Context.
type slackContext struct {
	refreshers []entity.ContextRefresher
	items      []entity.ContextItem
}

// AddRefresher implements entity.Context.
func (c *slackContext) AddRefresher(r entity.ContextRefresher) {
	c.refreshers = append(c.refreshers, r)
}

// Refresh implements entity.Context.
func (c *slackContext) Refresh(ctx context.Context) error {
	for _, r := range c.refreshers {
		if err := r(ctx); err != nil {
			return failure.Wrap(err, failure.Message("failed to refresh context"))
		}
	}

	return nil
}

// Items implements entity.Context.
func (c *slackContext) Items() []entity.ContextItem {
	return slices.Clone(c.items)
}
