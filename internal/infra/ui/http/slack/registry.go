package slack

import "context"

type registry struct {
	Dummy string
}

type registryKey struct{}

func NewRegistry(ctx context.Context) *registry {
	return &registry{
		Dummy: "dummy",
	}
}
