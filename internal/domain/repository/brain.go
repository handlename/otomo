package repository

import (
	"context"

	"github.com/handlename/otomo/internal/domain/entity"
)

type Brain interface {
	New(context.Context) (entity.Brain, error)
}
