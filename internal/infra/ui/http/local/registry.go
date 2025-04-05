package local

import (
	"context"

	drepo "github.com/handlename/otomo/internal/domain/repository"
	irepo "github.com/handlename/otomo/internal/infra/repository"
)

type registry struct {
	RepoSession drepo.Session
	RepoBrain   drepo.Brain
}

func NewRegistry(ctx context.Context) *registry {
	return &registry{
		RepoSession: &irepo.VolatileSession{},
		RepoBrain:   irepo.NewGeneralBrain(ctx),
	}
}
